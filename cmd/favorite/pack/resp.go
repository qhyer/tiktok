package pack

import (
	"errors"
	"tiktok/kitex_gen/favorite"
	"tiktok/pkg/errno"
)

func BuildFavoriteListResp(err error) *favorite.DouyinFavoriteListResponse {
	if err == nil {
		return favoriteListResp(errno.Success)
	}
	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return favoriteListResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return favoriteListResp(s)
}

func favoriteListResp(err errno.ErrNo) *favorite.DouyinFavoriteListResponse {
	return &favorite.DouyinFavoriteListResponse{StatusCode: err.ErrCode, StatusMsg: &err.ErrMsg}
}
