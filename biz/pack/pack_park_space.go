package pack

import (
	"github.com/XZ0730/hertz-scaffold/biz/model/base"
	"github.com/XZ0730/hertz-scaffold/biz/model/manager"
)

func PackIdleParkSpace(resp *base.GetIdleParkSpaceResponse, code int64, msg string, list []*base.ParkSpaceKV) {
	resp.Code = code
	resp.Msg = msg
	resp.Parklist = make([]*base.ParkSpaceKV, 0)
	resp.Parklist = append(resp.Parklist, list...)
	resp.Total = int64(len(list))
}

func PackGarageParkSpace(resp *manager.GetParkSpaceResponse, code int64, msg string, gid string, list []*manager.ParkSpaceModel) {
	resp.Code = code
	resp.Msg = msg
	resp.GarageID = gid
	resp.Plist = make([]*manager.ParkSpaceModel, 0)
	resp.Plist = append(resp.Plist, list...)
	resp.Total = int64(len(resp.Plist))
}

func PackPark(resp *base.ParkResponse, code int64, msg string, psv *base.ParkSpaceKV) {
	resp.Code = code
	resp.Msg = msg
	resp.Psv = psv
}
