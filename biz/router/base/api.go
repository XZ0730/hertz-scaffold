// Code generated by hertz generator. DO NOT EDIT.

package base

import (
	base "github.com/XZ0730/hertz-scaffold/biz/handler/base"
	"github.com/cloudwego/hertz/pkg/app/server"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	root.GET("/ping", append(_pingMw(), base.Ping)...)
	{
		_auth := root.Group("/auth", _authMw()...)
		{
			_login := _auth.Group("/login", _loginMw()...)
			_login.POST("/pwd", append(_login0Mw(), base.Login)...)
		}
		{
			_register := _auth.Group("/register", _registerMw()...)
			_register.POST("/pwd", append(_register0Mw(), base.Register)...)
		}
		{
			_user := _auth.Group("/user", _userMw()...)
			_user.PUT("/", append(_putuserinfoMw(), base.PutUserInfo)...)
			_user.POST("/audit", append(_postvehicleauditMw(), base.PostVehicleAudit)...)
			_user.POST("/park", append(_parkMw(), base.Park)...)
			_user.GET("/park_space", append(_getidleparkspaceMw(), base.GetIdleParkSpace)...)
			_user.GET("/vehicles", append(_getmyvehiclesMw(), base.GetMyVehicles)...)
		}
	}
}
