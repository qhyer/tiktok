package main

import (
	"context"
	user "tiktok/kitex_gen/user"
)

// UserSrvImpl implements the last service interface defined in the IDL.
type UserSrvImpl struct{}

// Register implements the UserSrvImpl interface.
func (s *UserSrvImpl) Register(ctx context.Context, req *user.DouyinUserRegisterRequest) (resp *user.DouyinUserRegisterResponse, err error) {
	// TODO: Your code here...
	return
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
