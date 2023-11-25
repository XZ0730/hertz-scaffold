package db

import "gorm.io/gorm"

// `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
// `vehicle_number` VARCHAR(100) NOT NULL COMMENT '车牌号',
// `brand` VARCHAR(100) NOT NULL COMMENT '车品牌',
// `model` VARCHAR(100) NOT NULL COMMENT '型号',
// `uid` BIGINT NOT NULL COMMENT '车主id',

type Vehicle struct {
	Id            int64
	VehicleNumber string
	Brand         string // 品牌
	Model         string // 型号
	Uid           int64
}
type VehicleAudit struct {
	Id              int64
	VehicleNumber   string
	VehicleNumImage string // 品牌
	GarageId        string // 型号
	ParkId          int64
	Uid             int64
	Description     string
	IsAudit         int
}

type VehicleAuditVO struct {
	Id              int64
	VehicleNumber   string
	VehicleNumImage string
	GarageId        string
	ParkId          int64
	Uid             int64
	Description     string
	IsAudit         int
	Brand           string // 品牌
	Model           string // 型号
}

func NewVehicle() *Vehicle {
	return &Vehicle{}
}
func NewVehicleAudit() *VehicleAudit {
	return &VehicleAudit{}
}

// 发送审核
func PostVehicleAudit(va *VehicleAudit) error {
	return _db.Table("vehicle_audit").Create(&va).Error
}

// 通过车牌号获取车辆
func GetVehicleByVnum(vnum string) (*Vehicle, error) {
	ve := NewVehicle()
	err := _db.Table("vehicle").Where("vehicle_number=?", vnum).First(ve).Error
	return ve, err
}

// 获取当前车库这个车牌审核通过的车辆
func GetAuditVehicleByVnum(vnum, gid string) (*Vehicle, error) {
	ve := NewVehicle()
	err := _db.Table("vehicle_audit").Where("vehicle_number=? AND garage_id=? AND is_audit=1", vnum, gid).First(ve).Error
	return ve, err
}

// 判断重复审核
func JudgeRepeatAudit(vnum string, gid string, pid int64, uid int64) bool {
	va := NewVehicleAudit()
	err := _db.Table("vehicle_audit").Where("garage_id=? AND vehicle_number=? AND park_id=? AND uid=?", gid, vnum, pid, uid).First(va).Error
	return err == gorm.ErrRecordNotFound
}

// 当提交审核前确认 车辆是否登记在车库表时，会先校准
// 校准是否车辆在车辆表中，如果在，则查询车主id是不是当前用户id
// 如果不是则无法提交，提示车牌号已经被其他账号注册，请联系管理员
// 如果车辆不在表中，则直接插入车辆表
func JudgeVehicleOwner(uid int64, vnum string) error {
	ve := NewVehicle()
	err := _db.Table("vehicle").Where("uid=? AND vehicle_number=?", uid, vnum).First(ve).Error
	return err
}

// 通过用户id获取车辆 （通过审核）
func GetVehicles(uid int64) ([]*Vehicle, error) {
	vlist := make([]*Vehicle, 0)
	err := _db.Table("vehicle").Joins("JOIN vehicle_audit ON vehicle_audit.vehicle_number=vehicle.vehicle_number AND vehicle_audit.is_audit=1").
		Where("uid=? ", uid).Find(&vlist).Error
	return vlist, err
}

func CreateVehicle(ve *Vehicle) error {
	return _db.Table("vehicle").Create(ve).Error
}

func CreateVehicleAudit(vea *VehicleAudit) error {
	return _db.Table("vehicle_audit").Create(vea).Error
}

func GetVehicleAudit(gid string) ([]*VehicleAuditVO, error) {
	va_list := make([]*VehicleAuditVO, 0)
	err := _db.Table("vehicle_audit").Select("vehicle_audit.id", "vehicle_audit.vehicle_number", "vehicle.brand", "vehicle.model", "vehicle_audit.garage_id",
		"vehicle_audit.vehicle_num_image", "vehicle_audit.description", "vehicle_audit.park_id", "vehicle_audit.is_audit", "vehicle_audit.uid").
		Joins("JOIN vehicle ON vehicle.vehicle_number=vehicle_audit.vehicle_number").
		Where("vehicle_audit.is_audit=0 AND vehicle_audit.garage_id=?", gid).
		Find(&va_list).Error
	return va_list, err
}

func AuditVehicle(aid int64, is_audit int, va *VehicleAudit) error {
	va.IsAudit = is_audit
	err := _db.Table("vehicle_audit").Where("id=?", aid).Updates(va).Error
	if err != nil {
		return err
	}

	if is_audit == 2 {
		ve := NewVehicle()
		err = _db.Table("vehicle").Unscoped().Where("vehicle_number=?", va.VehicleNumber).Delete(ve).Error
	} else {
		ps := NewParkSpace()
		ps.GarageId = va.GarageId
		ps.ParkType = 2 // 固定有归属
		ps.Uid = va.Uid
		ps.VehicleNumber = va.VehicleNumber
		ps.Id = va.ParkId
		err = _db.Table("park_space").Updates(ps).Error
	}
	return err
}

func JudegeAuditState(aid int64) (*VehicleAudit, error) {
	va := NewVehicleAudit()
	err := _db.Table("vehicle_audit").Where("id=?", aid).First(&va).Error
	return va, err
}
