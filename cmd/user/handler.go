package main

import (
	"context"

	"tiktok/cmd/user/service"
	"tiktok/dal/pack"
	"tiktok/kitex_gen/user"
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
		return resp, nil
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
		return resp, nil
	}

	resp = pack.BuildLoginResp(errno.Success)
	resp.UserId = userId
	return resp, nil
}

// GetUserInfoByUserIds implements the UserSrvImpl interface.
func (s *UserSrvImpl) GetUserInfoByUserIds(ctx context.Context, req *user.DouyinUserInfoRequest) (resp *user.DouyinUserInfoResponse, err error) {
	resp = new(user.DouyinUserInfoResponse)

	if len(req.GetToUserIds()) == 0 {
		resp = pack.BuildUserInfoResp(errno.ParamErr)
		return resp, nil
	}

	users, err := service.NewMGetUserService(ctx).MGetUser(req)
	if err != nil {
		resp = pack.BuildUserInfoResp(err)
		return resp, nil
	}

	resp = pack.BuildUserInfoResp(errno.Success)
	resp.User = users
	return resp, nil
}
