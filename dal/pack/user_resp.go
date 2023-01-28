package pack

import (
	"errors"

	"tiktok/kitex_gen/user"
	"tiktok/pkg/errno"
)

func BuildRegisterResp(err error) *user.DouyinUserRegisterResponse {
	if err == nil {
		return registerResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return registerResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return registerResp(s)
}

func registerResp(err errno.ErrNo) *user.DouyinUserRegisterResponse {
	return &user.DouyinUserRegisterResponse{StatusCode: err.ErrCode, StatusMsg: &err.ErrMsg}
}

func BuildUserInfoResp(err error) *user.DouyinUserInfoResponse {
	if err == nil {
		return userInfoResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return userInfoResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return userInfoResp(s)
}

func userInfoResp(err errno.ErrNo) *user.DouyinUserInfoResponse {
	return &user.DouyinUserInfoResponse{StatusCode: err.ErrCode, StatusMsg: &err.ErrMsg}
}

func BuildLoginResp(err error) *user.DouyinUserLoginResponse {
	if err == nil {
		return loginResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return loginResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return loginResp(s)
}

func loginResp(err errno.ErrNo) *user.DouyinUserLoginResponse {
	return &user.DouyinUserLoginResponse{StatusCode: err.ErrCode, StatusMsg: &err.ErrMsg}
}
