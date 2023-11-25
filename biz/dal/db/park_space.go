package db

import (
	"fmt"
	"time"
)

type ParkSpace struct {
	Id            int64
	GarageId      string
	Uid           int64
	VehicleNumber string
	ParkState     int
	ParkType      int
}

type ParkSpaceVO struct {
	Id            int64
	GarageId      string
	Uid           int64
	VehicleNumber string
	ParkState     int
	Address       string
}

type ParkRecord struct {
	Id            int64
	GarageId      string
	ParkId        int64
	VehicleNumber string
	Entries       time.Time
	Departure     time.Time
}

func NewParkRecord() *ParkRecord {
	return &ParkRecord{}
}

func NewParkSpace() *ParkSpace {
	return &ParkSpace{}
}

// 获取某车辆的车位信息
func GetParkSpaceByVN(uid int64, vn string) ([]*ParkSpaceVO, error) {
	ps := make([]*ParkSpaceVO, 0)
	err := _db.Table("park_space").Select("garage_info.address", "park_space.id", "park_space.garage_id", "park_space.uid", "park_space.vehicle_number", "park_space.park_state").
		Joins("JOIN garage_info ON garage_info.id=park_space.garage_id").
		Where("uid=? AND vehicle_number=?", uid, vn).Find(&ps).Error
	return ps, err
}

// 获取车库中没有归属的空闲车位
func GetIdleParkSpaceByGid(garage_id string) ([]*ParkSpace, error) {
	ps := make([]*ParkSpace, 0)
	err := _db.Table("park_space").Where("garage_id=? AND park_state=0 AND park_type=0", garage_id).Find(&ps).Error
	return ps, err
}

func GetParkSpaceByPidGid(pid int64, gid string) (*ParkSpace, error) {
	ps := NewParkSpace()
	err := _db.Table("park_space").Where("id=? AND garage_id=?", pid, gid).First(ps).Error
	return ps, err
}

func GetParkSpaceByGid(gid string) ([]*ParkSpace, error) {
	ps := make([]*ParkSpace, 0)
	err := _db.Table("park_space").Where("garage_id=?", gid).Find(&ps).Error
	return ps, err
}

func GetParkDuration(gid string, pid int64, v_number string) (time.Duration, error) {
	pr := NewParkRecord()
	err := _db.Table("park_record").Where("garage_id=? AND park_id=? AND vehicle_number=? AND departure IS NULL", gid, pid, v_number).First(pr).Error
	if err != nil {
		return time.Second - 1, err
	}
	now := time.Now()
	duration := now.Sub(pr.Entries)
	duration += (time.Hour * 8)
	return duration, err
}

func GetParkRecordByTime(gid string, startTime, endTime string) ([]*ParkRecord, error) {
	pr := make([]*ParkRecord, 0)
	err := _db.Select("park_space.uid", "park_record.garage_id", "park_record.park_id", "park_record.vehicle_number", "park_record.entries", "park_record.departure").
		Joins("JOIN park_space ON park_space.vehicle_number=park_record.vehicle_number").
		Where("(park_record.entries>=? OR park_record.departure>=?) AND (park_record.departure<=? OR park_record.entries<=?) AND park_record.garage_id=?", startTime, startTime, endTime, endTime, gid).
		Find(&pr).Error
	return pr, err
}

func CntVolume(gid string) ([]int64, error) {
	cnts := make([]int64, 0)
	now := time.Now()
	year_month_day := now.Format(time.DateOnly)
	pre := 0
	for i := 0; i < 12; i++ {
		var cnt int64 = 0
		start_time := year_month_day + fmt.Sprintf(" %02d", pre) + ":00:00"
		end_time := year_month_day + fmt.Sprintf(" %02d", pre+2) + ":00:00"
		_db.Table("park_record").
			Where("(entries>=? AND entries<=?) OR (departure>=? AND departure<=?) AND garage_id=?", start_time, end_time, start_time, end_time, gid).
			Count(&cnt)
		cnts = append(cnts, cnt)
		pre += 2
	}

	return cnts, nil
}

func GetParkSpaceByUid(uid int64, vn string, gid string) (*ParkSpace, error) {
	ps := NewParkSpace()
	err := _db.Table("park_space").Where("uid=? AND park_type=2 AND vehicle_number=? AND garage_id=?", uid, vn, gid).First(&ps).Error
	return ps, err
}

func UpdateParkSpaceState(state int64, uid, pid int64, gid string, vn string) error {
	ps := NewParkSpace()
	ps.ParkState = int(state)
	ps.Id = pid
	ps.GarageId = gid
	ps.VehicleNumber = vn
	ps.Uid = uid
	err := _db.Table("park_space").Where("id=? AND garage_id=?", pid, gid).Updates(&ps).Error
	return err
}

func CreateParkRecord(vn, gid string, pid int64, en time.Time) error {
	pr := NewParkRecord()
	pr.Entries = en
	pr.GarageId = gid
	pr.VehicleNumber = vn
	pr.ParkId = pid
	return _db.Table("park_record").Select("garage_id", "vehicle_number", "park_id", "entries").Create(&pr).Error
}

func JudgeFull(gid string) bool {
	var cnt int64
	_db.Table("park_space").Where("garage_id=? AND park_type=0 AND park_state=0", gid).Count(&cnt)
	return cnt == 0
}

func JudgeParkNow(vn string, gid string) bool {
	var cnt int64
	_db.Table("park_record").Where("garage_id=? AND vehicle_number=? AND departure IS NULL", gid, vn).Count(&cnt)
	return cnt == 1
}

func UpdateParkRecord(vn, gid string, pid int64) error {

	return _db.Table("park_record").Where("vehicle_number=? AND garage_id=? AND park_id=?", vn, gid, pid).UpdateColumn("departure", time.Now()).Error
}
