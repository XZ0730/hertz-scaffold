// Code generated by hertz generator.

package base

import (
	"context"

	"github.com/XZ0730/hertz-scaffold/biz/model/base"
	"github.com/XZ0730/hertz-scaffold/biz/pack"
	"github.com/XZ0730/hertz-scaffold/pkg/errno"
	"github.com/XZ0730/hertz-scaffold/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func rootMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _pingMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _authMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _loginMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _login0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _registerMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _register0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _putuserinfoMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		func(c context.Context, ctx *app.RequestContext) {
			token_byte := ctx.GetHeader("token")
			if _, err := utils.CheckToken(string(token_byte)); err != nil {
				resp := base.NewBaseResponse()
				pack.PackBase(resp, errno.AuthorizationFailedErrCode, errno.AuthorizationFailedError.ErrorMsg)
				ctx.JSON(consts.StatusOK, resp)
				ctx.Abort()
				return
			}
			ctx.Next(c)
		},
	}
}

func _getmyvehiclesMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getidleparkspaceMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _postvehicleauditMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _parkMw() []app.HandlerFunc {
	// your code...
	return nil
}
