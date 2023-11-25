package pack

import (
	"sort"

	"github.com/XZ0730/hertz-scaffold/biz/model/manager"
)

func PackParkRecord(resp *manager.GetParkRecordListResp, code int64, msg string, gid string, list []*manager.ParkRecordModel) {
	resp.Code = code
	resp.Msg = msg
	resp.GarageID = gid
	resp.Prlist = make([]*manager.ParkRecordModel, 0)
	resp.Prlist = append(resp.Prlist, list...)
	sort.Slice(resp.Prlist, func(i, j int) bool {
		return resp.Prlist[i].Time < resp.Prlist[j].Time
	})
	resp.Total = int64(len(resp.Prlist))
}
