package pack

import "github.com/XZ0730/hertz-scaffold/biz/model/base"

func PackBase(resp *base.BaseResponse, code int64, msg string) {
	resp.Code = code
	resp.Message = msg
}
func PackLogin(resp *base.LoginResponse, code int64, msg string, token string) {
	resp.Base = new(base.BaseResponse)
	PackBase(resp.Base, code, msg)
	resp.Token = token
}
