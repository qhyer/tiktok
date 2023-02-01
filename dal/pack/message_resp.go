package pack

import (
	"errors"

	"tiktok/kitex_gen/message"
	"tiktok/pkg/errno"
)

func BuildMessageActionResp(err error) *message.DouyinMessageActionResponse {
	if err == nil {
		return messageActionResp(errno.Success)
	}
	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return messageActionResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return messageActionResp(s)
}

func messageActionResp(err errno.ErrNo) *message.DouyinMessageActionResponse {
	return &message.DouyinMessageActionResponse{StatusCode: err.ErrCode, StatusMsg: &err.ErrMsg}
}

func BuildMessageListResp(err error) *message.DouyinMessageListResponse {
	if err == nil {
		return messageListResp(errno.Success)
	}
	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return messageListResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return messageListResp(s)
}

func messageListResp(err errno.ErrNo) *message.DouyinMessageListResponse {
	return &message.DouyinMessageListResponse{StatusCode: err.ErrCode, StatusMsg: &err.ErrMsg}
}
