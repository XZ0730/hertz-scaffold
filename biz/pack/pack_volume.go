package pack

import "github.com/XZ0730/hertz-scaffold/biz/model/manager"

func PackCntVolume(resp *manager.CntVolumeResp, code int64, msg string, list []int64) {
	resp.Code = code
	resp.Msg = msg
	resp.Cnt = make([]int64, 0)
	resp.Cnt = append(resp.Cnt, list...)
}
