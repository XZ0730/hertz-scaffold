package pack

import "github.com/XZ0730/hertz-scaffold/biz/model/manager"

func PackVehicleAudit(resp *manager.GetVehicleAuditResp, code int64, msg string, va_list []*manager.VehicleAuditModel) {
	resp.Code = code
	resp.Msg = msg
	resp.VaList = make([]*manager.VehicleAuditModel, 0)
	resp.VaList = append(resp.VaList, va_list...)
	resp.Total = int64(len(resp.VaList))
}
