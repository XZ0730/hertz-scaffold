package service

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/XZ0730/hertz-scaffold/biz/dal/db"
	"github.com/XZ0730/hertz-scaffold/biz/model/base"
	"github.com/XZ0730/hertz-scaffold/pkg/errno"
	"github.com/XZ0730/hertz-scaffold/pkg/utils"
	"github.com/cloudwego/kitex/pkg/klog"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

func (u *UserService) Login(req *base.LoginRequest) (int64, string, string) {
	//
	ui, err := db.GetUserInfoByUserName(req.Username)
	if err != nil {
		klog.Error("[login]:", err.Error())
		return errno.UserNameError.ErrorCode, errno.UserNameError.ErrorMsg, ""
	}
	if ui.Password != req.Password {
		klog.Error("[login]: pwd error")
		return errno.PWDError.ErrorCode, errno.PWDError.ErrorMsg, ""
	}
	token, err := utils.CreateToken(ui.Id, int64(ui.Role))
	if err != nil {
		klog.Error("[token]", err.Error())
	}
	return errno.StatusSuccessCode, errno.StatusSuccessMsg, token
}

func (u *UserService) Register(req *base.RegisterReq) (int64, string) {
	//

	err := db.CreateUser(req.GetUsername(), req.GetPassword())
	if err != nil {
		klog.Error("[register]:", err.Error())
		return errno.ParamError.ErrorCode, errno.ParamError.ErrorMsg
	}
	return errno.StatusSuccessCode, errno.StatusSuccessMsg
}

func (u *UserService) UpdateUserInfo(uid int64, req *base.UpdateUserInfoReq) (int64, string) {

	err := db.JudgeIDCard(req.GetCardID())
	if err != gorm.ErrRecordNotFound || err == nil {
		klog.Error("[id card] exist")
		return errno.ExistError.ErrorCode, errno.ExistError.ErrorMsg
	}
	if err = db.JudgePhoneNumber(req.GetPhone()); err != gorm.ErrRecordNotFound || err == nil {
		klog.Error("[phone] exist")
		return errno.ExistError.ErrorCode, errno.ExistError.ErrorMsg
	}
	if err != gorm.ErrRecordNotFound || err == nil {
		klog.Error("[id card] exist")
		return errno.ExistError.ErrorCode, errno.ExistError.ErrorMsg
	}

	ui := db.NewUserInfo()
	ui.Id = uid
	ui.Age = int(req.GetAge())
	ui.Name = req.GetName()
	ui.Telephone = req.GetPhone()
	ui.Sex = req.GetSex()
	ui.Role = 1
	ui.CardId = req.GetCardID()
	ui.Telephone = req.GetPhone()
	err = db.UpdateUserInfo(ui)
	if err != nil {
		klog.Error("[update user_info]", err.Error())
		return errno.UpdateError.ErrorCode, errno.UpdateError.ErrorMsg
	}
	return errno.StatusSuccessCode, errno.StatusSuccessMsg
}

// 获取我的车辆信息 -
func (u *UserService) GetMyVehicles(uid int64) (int64, string, []*base.VehicleModel) {
	vm := make([]*base.VehicleModel, 0)

	vehicles, err := db.GetVehicles(uid)
	if err != nil {
		klog.Error("[get vehicle]:", err.Error())
		return errno.GetError.ErrorCode, errno.GetError.ErrorMsg, nil
	}
	var eg errgroup.Group
	for _, val := range vehicles {
		tmp := val
		eg.Go(func() error {
			vo_g := new(base.VehicleModel)

			psk := make([]*base.ParkSpaceKV, 0)
			psv, err := db.GetParkSpaceByVN(uid, tmp.VehicleNumber)
			if err != nil {
				return err
			}
			vo_g.VehicleNumber = tmp.VehicleNumber
			vo_g.Brand = tmp.Brand
			vo_g.Model = tmp.Model
			vo_g.Parklist = make([]*base.ParkSpaceKV, 0)
			for _, v := range psv {
				pskv := new(base.ParkSpaceKV)
				pskv.GarageID = v.GarageId
				pskv.ParkID = v.Id
				psk = append(psk, pskv)
			}
			vo_g.Parklist = append(vo_g.Parklist, psk...)
			vm = append(vm, vo_g)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		klog.Info("[vehicle]get error:", err.Error())
		return errno.GetError.ErrorCode, errno.GetError.ErrorMsg, nil
	}
	return errno.StatusSuccessCode, errno.StatusSuccessMsg, vm
}

// 获取某车库的空闲车位
func (u *UserService) GetIdleParkSpace(gid string) (int64, string, []*base.ParkSpaceKV) {
	psk := make([]*base.ParkSpaceKV, 0)
	ps, err := db.GetIdleParkSpaceByGid(gid)
	if err != nil {
		klog.Error("[idle park]", err.Error())
		return errno.GetError.ErrorCode, errno.GetError.ErrorMsg, nil
	}
	var eg errgroup.Group
	for _, val := range ps {
		tmp := val
		eg.Go(func() error {
			psk2 := new(base.ParkSpaceKV)
			psk2.GarageID = tmp.GarageId
			psk2.ParkID = tmp.Id
			psk = append(psk, psk2)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		klog.Info("[vehicle]get error:", err.Error())
		return errno.GetError.ErrorCode, errno.GetError.ErrorMsg, nil
	}
	return errno.StatusSuccessCode, errno.StatusSuccessMsg, psk
}

// 申请车位 车辆绑定
func (u *UserService) PostVehicleAudit(uid int64, req *base.PostVehicleAuditReq,
	vn_image *multipart.FileHeader) (int64, string) {
	// 判断 是否已经提交过 相同的审核
	if !db.JudgeRepeatAudit(req.GetVehicleNumber(), req.GetGarageID(), req.GetParkID(), uid) {
		klog.Error("[post error] audit exist")
		return errno.ExistError.ErrorCode, errno.ExistError.ErrorMsg
	}
	// 然后判断这个车牌号是否有账号已经注册,如果已经注册 则返回
	// 判断这个车牌号是否已经被注册过了
	vehicle, err := db.GetAuditVehicleByVnum(req.GetVehicleNumber(), req.GetGarageID())
	if err != nil && err != gorm.ErrRecordNotFound {
		klog.Error("[post audit]:", err.Error())
		return errno.GetError.ErrorCode, errno.GetError.ErrorMsg
	} else if vehicle.Uid != uid && vehicle.Id != 0 {
		// 当前车库这个车牌号审核通过记录存在 且 被其他人注册过了
		klog.Error("[post audit]:vehicle has exist owner")
		return errno.ExistError.ErrorCode, errno.ExistError.ErrorMsg
	}

	// 然后判断当前车库这个车位是否有归属 有归属则返回
	ps, err := db.GetParkSpaceByPidGid(req.GetParkID(), req.GetGarageID())
	if err != nil {
		klog.Error("[post audit]", err.Error())
		return errno.GetError.ErrorCode, errno.GetError.ErrorMsg
	} else if ps.ParkType != 0 {
		klog.Error("[post audit]", err.Error())
		return errno.ExistError.ErrorCode, errno.ExistError.ErrorMsg
	}
	// 然后将文件上传 失败返回
	file, err := vn_image.Open()
	if err != nil {
		klog.Error("[post audit]", err.Error())
		return errno.FileError.ErrorCode, errno.FileError.ErrorMsg
	}
	code, url := utils.UploadToQiNiu(file, vn_image, fmt.Sprint(uid))
	if code != 200 {
		klog.Error("[post audit] upload error :", url)
		return errno.UploadError.ErrorCode, errno.UploadError.ErrorMsg
	}
	// 然后创建车辆
	ve := db.NewVehicle()
	ve.Brand = req.GetBrand()
	ve.Model = req.GetModel()
	ve.VehicleNumber = req.GetVehicleNumber()
	ve.Uid = uid
	if err = db.CreateVehicle(ve); err != nil {
		klog.Error("[post audit]:", err.Error())
		return errno.CreateError.ErrorCode, errno.CreateError.ErrorMsg
	}
	// 创建审核记录等待审核
	va := db.NewVehicleAudit()
	va.GarageId = req.GetGarageID()
	va.ParkId = req.GetParkID()
	va.Uid = uid
	va.VehicleNumber = req.GetVehicleNumber()
	va.VehicleNumImage = url
	va.Description = req.GetDescription()
	if err := db.CreateVehicleAudit(va); err != nil {
		klog.Error("[post audit]", err.Error())
		return errno.CreateError.ErrorCode, errno.CreateError.ErrorMsg
	}
	return errno.StatusSuccessCode, errno.StatusSuccessMsg
}

func (u *UserService) Park(uid int64, req *base.ParkRequest) (int64, string, *base.ParkSpaceKV) {

	// 判断当前车库是否已经停放该车辆
	if ok := db.JudgeParkNow(req.GetVehicleNumber(), req.GetGarageID()); ok {
		klog.Error("[park] have park")
		return 797979, "已经入库不能再次入库", nil
	}

	// 入库先判断是否是 违规用户，如果不是则判断是否当前一小时内入库了三次，
	// 如果 入库操作超过了三次则封禁账号
	if ok := db.JudgeViolation(uid); ok {
		klog.Error("[park] violation account")
		return 99999, "账号违规，联系管理员解除", nil
	}
	if ok := db.JudgeParkCnt(req.GetVehicleNumber(), req.GetGarageID()); ok {
		err := db.UpdateViolation(uid)
		if err != nil {
			klog.Error("[violation]", err.Error())
			return errno.UpdateError.ErrorCode, errno.UpdateError.ErrorMsg, nil
		}
		return 99998, "频繁入库,违规封禁，联系管理员解除", nil
	}
	// 如果车主在当前车库存在车位，则允许入库
	ps, err := db.GetParkSpaceByUid(uid, req.GetVehicleNumber(), req.GetGarageID())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 无车位
			// 当前车库车位是否已满
			// 有归属的车位算已占车位
			ok := db.JudgeFull(req.GetGarageID())
			if ok {
				klog.Error("[full] full parkspace")
				return 898989, "车位满了", nil
			}
			ps2, err := db.GetIdleParkSpaceByGid(req.GetGarageID())
			if err != nil {
				klog.Error("[park error]", err.Error())
				return 828282, "车位获取失败", nil
			}
			klog.Info(ps2[0].Id)
			err = db.UpdateParkSpaceState(1, uid, ps2[0].Id, ps2[0].GarageId, req.GetVehicleNumber())
			if err != nil {
				klog.Error("[park error],", err.Error())
				return 909090, "车位更新失败", nil
			}
			err = db.CreateParkRecord(req.GetVehicleNumber(), req.GetGarageID(), ps2[0].Id, time.Now())
			if err != nil {
				klog.Error("[park error]", err.Error())
				return 818181, "入库失败", nil
			}
			return errno.StatusSuccessCode, errno.StatusSuccessMsg, &base.ParkSpaceKV{
				ParkID:   ps2[0].Id,
				GarageID: ps2[0].GarageId,
			}
		}
		klog.Error("[park get],", err.Error())
		return errno.GetError.ErrorCode, errno.GetError.ErrorMsg, nil
	}
	// 有车位
	// 更新车位状态然后返回
	err = db.UpdateParkSpaceState(1, uid, ps.Id, ps.GarageId, ps.VehicleNumber)
	if err != nil {
		klog.Error("[park error],", err.Error())
		return 909090, "车位更新失败", nil
	}
	// 插入入库记录
	err = db.CreateParkRecord(req.GetVehicleNumber(), req.GetGarageID(), ps.Id, time.Now())
	if err != nil {
		klog.Error("[park error]", err.Error())
		return 818181, "入库失败", nil
	}

	return errno.StatusSuccessCode, errno.StatusSuccessMsg, &base.ParkSpaceKV{ParkID: ps.Id, GarageID: ps.GarageId}
}

func (u *UserService) LeavePark(uid int64, req *base.ParkRequest) (int64, string) {

	// 查询车位
	ps, err := db.GetParkSpaceByUid(uid, req.GetVehicleNumber(), req.GetGarageID())
	if err != nil {
		klog.Error("[leave],", err.Error())
		return errno.GetError.ErrorCode, errno.GetError.ErrorMsg
	}

	if ps.ParkType == 2 {
		// 固定车位则只需要将状态位更新
		err = db.UpdateParkSpaceState(0, uid, ps.Id, req.GetGarageID(), req.GetVehicleNumber())
		if err != nil {
			klog.Error("[leave],", err.Error())
			return errno.UpdateError.ErrorCode, errno.UpdateError.ErrorMsg
		}
	} else if ps.ParkType == 0 {
		// 游客车位则需要将字段设为空
		err = db.UpdateParkSpaceState(0, 0, ps.Id, req.GetGarageID(), "无")
		if err != nil {
			klog.Error("[leave],", err.Error())
			return errno.UpdateError.ErrorCode, errno.UpdateError.ErrorMsg
		}
	}

	// 然后更新停车记录
	err = db.UpdateParkRecord(req.GetVehicleNumber(), req.GetGarageID(), ps.Id)
	if err != nil {
		klog.Error("[leave],", err.Error())
		return errno.UpdateError.ErrorCode, errno.UpdateError.ErrorMsg
	}
	return errno.StatusSuccessCode, errno.StatusSuccessMsg
}
