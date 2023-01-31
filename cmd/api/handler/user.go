package handler

import (
	"context"
	"net/http"

	"tiktok/kitex_gen/user"
	"tiktok/pkg/errno"
	"tiktok/pkg/jwt"
	"tiktok/pkg/rpc"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/klog"
)

type RegisterParam struct {
	Username string `query:"username" vd:"len($)>=1&&len($)<=32"`
	Password string `query:"password" vd:"len($)>=6&&len($)<=32"`
}

type RegisterResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
}

// Register 用户注册
func Register(ctx context.Context, c *app.RequestContext) {
	var req RegisterParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		klog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, errno.ParamErr)
		return
	}

	// rpc通信
	registerResponse, err := rpc.Register(ctx, &user.DouyinUserRegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err)
		return
	}

	// 根据传回的userId生成token
	token, err := jwt.GenerateToken(registerResponse.GetUserId())
	if err != nil {
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, RegisterResponse{
		StatusCode: errno.Success.ErrCode,
		StatusMsg:  errno.Success.ErrMsg,
		UserId:     registerResponse.GetUserId(),
		Token:      token,
	})
}

type LoginParam struct {
	Username string `query:"username" vd:"len($)>=1&&len($)<=32"`
	Password string `query:"password" vd:"len($)>=6&&len($)<=32"`
}

type LoginResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
}

// Login 用户登录
func Login(ctx context.Context, c *app.RequestContext) {
	var req LoginParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, errno.ParamErr)
		return
	}

	// rpc通信
	loginResponse, err := rpc.Login(ctx, &user.DouyinUserLoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err)
		return
	}

	// 根据传回的userId生成token
	token, err := jwt.GenerateToken(loginResponse.GetUserId())
	if err != nil {
		hlog.CtxErrorf(ctx, "generate token error %v", err)
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		StatusCode: errno.Success.ErrCode,
		StatusMsg:  errno.Success.ErrMsg,
		UserId:     loginResponse.GetUserId(),
		Token:      token,
	})
}

type GetUserInfoParam struct {
	UserId int64 `query:"user_id" vd:"$>0"`
}

type GetUserInfoResponse struct {
	StatusCode int32      `json:"status_code"`
	StatusMsg  string     `json:"status_msg"`
	User       *user.User `json:"user"`
}

// GetUserInfo 获取用户信息
func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	var req GetUserInfoParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, errno.ParamErr)
		return
	}
	userId := c.GetInt64("UserID")

	// rpc通信
	var userIds []int64
	userIds = append(userIds, req.UserId)
	getUserInfoResponse, err := rpc.UserInfo(ctx, &user.DouyinUserInfoRequest{
		UserId:    userId,
		ToUserIds: userIds,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err)
		return
	}

	// 获取用户信息
	if len(getUserInfoResponse.GetUser()) == 0 {
		hlog.CtxWarnf(ctx, "user not exist error %v", err)
		SendResponse(c, errno.UserNotExistErr)
		return
	}
	usr := getUserInfoResponse.GetUser()[0]

	c.JSON(http.StatusOK, GetUserInfoResponse{
		StatusCode: errno.Success.ErrCode,
		StatusMsg:  errno.Success.ErrMsg,
		User:       usr,
	})
}
