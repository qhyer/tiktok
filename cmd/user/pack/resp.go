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
