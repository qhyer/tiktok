package pack

import (
	"errors"
	"tiktok/kitex_gen/feed"
	"tiktok/pkg/errno"
)

func BuildFeedResp(err error) *feed.DouyinFeedResponse {
	if err == nil {
		return feedResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return feedResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return feedResp(s)
}

func feedResp(err errno.ErrNo) *feed.DouyinFeedResponse {
	return &feed.DouyinFeedResponse{StatusCode: err.ErrCode, StatusMsg: &err.ErrMsg}
}
