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
		return resp, nil
	}
	resp = pack.BuildRegisterResp(errno.Success)
	resp.UserId = userId
	return resp, nil
}

// Login implements the UserSrvImpl interface.
func (s *UserSrvImpl) Login(ctx context.Context, req *user.DouyinUserLoginRequest) (resp *user.DouyinUserLoginResponse, err error) {
	// TODO: Your code here...
	return
}

// GetUserInfoById implements the UserSrvImpl interface.
func (s *UserSrvImpl) GetUserInfoById(ctx context.Context, req *user.DouyinUserInfoRequest) (resp *user.DouyinUserInfoResponse, err error) {
	// TODO: Your code here...
	return
}
