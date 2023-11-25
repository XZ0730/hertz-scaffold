package service

import (
	"time"

	"github.com/XZ0730/hertz-scaffold/biz/dal/db"
	"github.com/XZ0730/hertz-scaffold/biz/model/manager"
	"github.com/XZ0730/hertz-scaffold/pkg/errno"
	"github.com/cloudwego/kitex/pkg/klog"
	"golang.org/x/sync/errgroup"
)

func (m *ManagerService) GetGarageParkSpace(gid string) (int64, string, []*manager.ParkSpaceModel) {
	psm := make([]*manager.ParkSpaceModel, 0)
	// 查询该车库的所有车位
	ps, err := db.GetParkSpaceByGid(gid)
	if err != nil {
		klog.Error("[manage ps fetch]:", err.Error())
		return errno.GetError.ErrorCode, errno.GetError.ErrorMsg, nil
	}
	var eg errgroup.Group
	for _, val := range ps {
		tmp := val
		eg.Go(func() error {
			psm2 := new(manager.ParkSpaceModel)
			psm2.GarageID = tmp.GarageId
			psm2.ID = tmp.Id
			psm2.ParkState = int64(tmp.ParkState)
			psm2.UID = tmp.Uid
			psm2.ParkType = int64(tmp.ParkType)
			if tmp.ParkState == 1 && tmp.ParkType == 0 {
				// 空闲车位 有车停入
				// 查询车位记录的已经停车的时间
				// 车位记录 条件：离开时间为空 入库时间不为空
				duration, err := db.GetParkDuration(gid, tmp.Id, tmp.VehicleNumber)
				if err != nil {
					klog.Error("[duration]", err.Error())
					psm = append(psm, psm2)
					return err
				}
				psm2.Duration = duration.String()
			}
			psm = append(psm, psm2)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		klog.Error("[manager ps fetch]", err.Error())
		return errno.TimeError.ErrorCode, errno.TimeError.ErrorMsg, psm
	}

	return errno.StatusSuccessCode, errno.StatusSuccessMsg, psm
}

func (m *ManagerService) GetParkRecordByTime(req *manager.GetParkRecordReq) (int64, string, []*manager.ParkRecordModel) {
	prm := make([]*manager.ParkRecordModel, 0)
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		klog.Error("[time error]", err.Error())
		return errno.TimeError.ErrorCode, errno.TimeError.ErrorMsg, nil
	}
	_, err = time.ParseInLocation(time.DateOnly, req.GetStartTime(), loc)
	if err != nil {
		klog.Error("[time error]", err.Error())
		return errno.TimeError.ErrorCode, errno.TimeError.ErrorMsg, nil
	}
	_, err = time.ParseInLocation(time.DateOnly, req.GetEndTime(), loc)
	if err != nil {
		klog.Error("[time error]", err.Error())
		return errno.TimeError.ErrorCode, errno.TimeError.ErrorMsg, nil
	}
	pr, err := db.GetParkRecordByTime(req.GetGarageID(), req.GetStartTime(), req.GetEndTime())
	if err != nil {
		klog.Error("[park_record]", err.Error())
		return errno.GetError.ErrorCode, errno.GetError.ErrorMsg, nil
	}
	// 遍历pr
	var eg errgroup.Group
	for _, val := range pr {
		tmp := val
		eg.Go(func() error {
			entries := tmp.Entries.String()
			if entries >= req.GetStartTime() && entries <= req.GetEndTime() {
				pr := new(manager.ParkRecordModel)
				pr.Gid = tmp.GarageId
				pr.ParkID = tmp.ParkId
				pr.Type = 0
				pr.Time = entries
				pr.VehicleNumber = tmp.VehicleNumber
				prm = append(prm, pr)
			}
			if !tmp.Departure.IsZero() {
				departure := tmp.Departure.String()
				if departure >= req.GetStartTime() && departure <= req.GetEndTime() {
					pr := new(manager.ParkRecordModel)
					pr.Gid = tmp.GarageId
					pr.ParkID = tmp.ParkId
					pr.Type = 1
					pr.Time = departure
					pr.VehicleNumber = tmp.VehicleNumber
					prm = append(prm, pr)
				}
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		klog.Error("[manager pr fetch]", err.Error())
		return errno.TimeError.ErrorCode, errno.TimeError.ErrorMsg, prm
	}
	return errno.StatusSuccessCode, errno.StatusSuccessMsg, prm
}

func (m *ManagerService) CntVolume(gid string) (int64, string, []int64) {

	cnt, err := db.CntVolume(gid)
	if err != nil {
		klog.Error("[cnt_volume]", err.Error())
		return errno.GetError.ErrorCode, errno.GetError.ErrorMsg, nil
	}

	return errno.StatusSuccessCode, errno.StatusSuccessMsg, cnt
}

func (m *ManagerService) GetVehicleAudit(gid string) (int64, string, []*manager.VehicleAuditModel) {
	va_list := make([]*manager.VehicleAuditModel, 0)
	vam, err := db.GetVehicleAudit(gid)
	if err != nil {
		klog.Error("[vehicle_audit]", err.Error())
		return errno.GetError.ErrorCode, errno.GetError.ErrorMsg, nil
	}
	var eg errgroup.Group
	for _, val := range vam {
		tmp := val
		eg.Go(func() error {
			va := new(manager.VehicleAuditModel)
			va.Brand = tmp.Brand
			va.Description = tmp.Description
			va.Image = tmp.VehicleNumImage
			va.ParkID = tmp.ParkId
			va.Model = tmp.Model
			va.VehicleNumber = tmp.VehicleNumber
			va.GarageID = tmp.GarageId
			va.UID = tmp.Uid
			va.ID = tmp.Id
			va_list = append(va_list, va)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		klog.Error("[error]", err.Error())
		return errno.GetError.ErrorCode, errno.GetError.ErrorMsg, nil
	}
	return errno.StatusSuccessCode, errno.StatusSuccessMsg, va_list
}

func (m *ManagerService) AuditVehicle(aid int64, audit int64) (int64, string) {

	va, err := db.JudegeAuditState(aid)
	if err != nil {
		klog.Error("[audit vehicle]", err.Error())
		return errno.GetError.ErrorCode, errno.GetError.ErrorMsg
	}
	if va.IsAudit != 0 {
		klog.Error("[audit vehicle] have dealed")
		return errno.DealError.ErrorCode, errno.DealError.ErrorMsg
	}
	psv, err := db.GetParkSpaceByPidGid(va.ParkId, va.GarageId)
	if err != nil {
		klog.Error("[audit vehicle]", err.Error())
		return errno.GetError.ErrorCode, errno.GetError.ErrorMsg
	}
	if psv.ParkType != 0 {
		audit = 2
	} else if psv.ParkState != 0 {
		// 返回清理车位
		klog.Error("[exist] The parking space is occupied. Please clean the parking space and try again")
		return errno.WaitError.ErrorCode, errno.WaitError.ErrorMsg
	}

	err = db.AuditVehicle(aid, int(audit), va)
	if err != nil {
		klog.Error("[audit vehicle]", err.Error())
		return errno.UpdateError.ErrorCode, errno.UpdateError.ErrorMsg
	}
	return errno.StatusSuccessCode, errno.StatusSuccessMsg
}
