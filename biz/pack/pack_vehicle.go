package pack

import "github.com/XZ0730/hertz-scaffold/biz/model/base"

func PackVehicles(resp *base.GetMyVehicleResponse, code int64, msg string, list []*base.VehicleModel) {
	resp.Code = code
	resp.Msg = msg
	resp.Vlist = list
	resp.Total = int64(len(list))
}
