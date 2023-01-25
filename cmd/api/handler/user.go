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
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
}

func Register(_ context.Context, c *app.RequestContext) {
	var req RegisterParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		SendResponse(c, err)
		return
	}

	// rpc通信
	registerResponse, err := rpc.Register(context.Background(), &user.DouyinUserRegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		SendResponse(c, err)
		return
	}

	// 根据传回的userId生成token
	token, err := util.GenerateToken(registerResponse.UserId)
	if err != nil {
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, RegisterResponse{
		StatusCode: errno.Success.ErrCode,
		StatusMsg:  errno.Success.ErrMsg,
		UserId:     registerResponse.UserId,
		Token:      token,
	})
}

// TODO: 用户名密码合法性校验 https://www.cloudwego.io/zh/docs/hertz/tutorials/basic-feature/binding-and-validate/
type LoginParam struct {
	Username string `query:"username"`
	Password string `query:"password"`
}

type LoginResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     int64  `json:"user_id"`
	Token      string `json:"token"`
}

func Login(_ context.Context, c *app.RequestContext) {
	var req LoginParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		SendResponse(c, err)
		return
	}

	// rpc通信
	loginResponse, err := rpc.Login(context.Background(), &user.DouyinUserLoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		SendResponse(c, err)
		return
	}

	// 根据传回的userId生成token
	token, err := util.GenerateToken(loginResponse.UserId)
	if err != nil {
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		StatusCode: errno.Success.ErrCode,
		StatusMsg:  errno.Success.ErrMsg,
		UserId:     loginResponse.UserId,
		Token:      token,
	})
}

// TODO: userId合法性校验 https://www.cloudwego.io/zh/docs/hertz/tutorials/basic-feature/binding-and-validate/
type GetUserInfoParam struct {
	UserId int64 `query:"user_id"`
}

type User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type GetUserInfoResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	User       User   `json:"user"`
}

func GetUserInfo(_ context.Context, c *app.RequestContext) {
	var req GetUserInfoParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		SendResponse(c, err)
		return
	}
	userId := c.GetInt64("UserID")

	// rpc通信
	var userIds []int64
	userIds = append(userIds, req.UserId)
	getUserInfoResponse, err := rpc.UserInfo(context.Background(), &user.DouyinUserInfoRequest{
		UserId:    userId,
		ToUserIds: userIds,
	})
	if err != nil {
		SendResponse(c, err)
		return
	}
	if len(getUserInfoResponse.User) == 0 {
		SendResponse(c, errno.UserNotExistErr)
		return
	}
	usr := getUserInfoResponse.User[0]

	c.JSON(http.StatusOK, GetUserInfoResponse{
		StatusCode: errno.Success.ErrCode,
		StatusMsg:  errno.Success.ErrMsg,
		User: User{
			Id:            usr.Id,
			Name:          usr.Name,
			FollowerCount: *usr.FollowerCount,
			FollowCount:   *usr.FollowCount,
			IsFollow:      usr.IsFollow,
		},
	})
}