package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"tiktok/cmd/api/rpc"
	"tiktok/cmd/api/util"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/errno"
)

// TODO: 用户名密码合法性校验 https://www.cloudwego.io/zh/docs/hertz/tutorials/basic-feature/binding-and-validate/
type RegisterParam struct {
	Username string `query:"username"`
	Password string `query:"password"`
}

type RegisterResponse struct {
	Response Response
	UserId   int64  `json:"user_id"`
	Token    string `json:"token"`
}

func Register(_ context.Context, ctx *app.RequestContext) {
	var req RegisterParam
	// 参数校验
	err := ctx.BindAndValidate(&req)
	if err != nil {
		SendResponse(ctx, err)
		return
	}
	// rpc通信
	registerResponse, err := rpc.Register(context.Background(), &user.DouyinUserRegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		SendResponse(ctx, err)
		return
	}
	// 根据传回的userId生成token
	token, err := util.GenerateToken(registerResponse.UserId)
	if err != nil {
		SendResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, RegisterResponse{
		Response: Response{
			StatusCode: errno.Success.ErrCode,
			StatusMsg:  errno.Success.ErrMsg,
		},
		UserId: registerResponse.UserId,
		Token:  token,
	})
}
