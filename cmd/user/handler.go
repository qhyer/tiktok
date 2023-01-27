package main

import (
	"context"
	"tiktok/cmd/user/pack"
	"tiktok/cmd/user/service"
	user "tiktok/kitex_gen/user"
	"tiktok/pkg/errno"
)

// UserSrvImpl implements the last service interface defined in the IDL.
type UserSrvImpl struct{}

// Register implements the UserSrvImpl interface.
func (s *UserSrvImpl) Register(ctx context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	resp = new(user.DouyinUserRegisterResponse)

	userId, err := service.NewCreateUserService(ctx).CreateUser(req)
	if err != nil {
		resp = pack.BuildRegisterResp(err)
		return resp, err
	}
	resp = pack.BuildRegisterResp(errno.Success)
	resp.UserId = userId
	return resp, nil
}

// Login implements the UserSrvImpl interface.
func (s *UserSrvImpl) Login(ctx context.Context, req *user.DouyinUserLoginRequest) (resp *user.DouyinUserLoginResponse, err error) {
	resp = new(user.DouyinUserLoginResponse)

	userId, err := service.NewCheckUserService(ctx).CheckUser(req)
	if err != nil {
		resp = pack.BuildLoginResp(err)
		return resp, err
	}
	resp = pack.BuildLoginResp(errno.Success)
	resp.UserId = userId
	return resp, nil
}

// GetUserInfoByUserIds implements the UserSrvImpl interface.
func (s *UserSrvImpl) GetUserInfoByUserIds(ctx context.Context, req *user.DouyinUserInfoRequest) (resp *user.DouyinUserInfoResponse, err error) {
	resp = new(user.DouyinUserInfoResponse)

	if len(req.ToUserIds) == 0 {
		resp = pack.BuildUserInfoResp(errno.ParamErr)
		return resp, err
	}

	users, err := service.NewMGetUserService(ctx).MGetUser(req)
	if err != nil {
		resp = pack.BuildUserInfoResp(err)
		return resp, err
	}
	resp = pack.BuildUserInfoResp(errno.Success)
	resp.User = users
	return resp, nil
}
