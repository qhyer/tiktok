package main

import (
	"context"

	"tiktok/cmd/favorite/service"
	"tiktok/dal/pack"
	"tiktok/kitex_gen/favorite"
	"tiktok/pkg/errno"
)

// FavoriteSrvImpl implements the last service interface defined in the IDL.
type FavoriteSrvImpl struct{}

const (
	DoFavoriteAction     = 1
	CancelFavoriteAction = 2
)

// FavoriteAction implements the FavoriteSrvImpl interface.
func (s *FavoriteSrvImpl) FavoriteAction(ctx context.Context, req *favorite.DouyinFavoriteActionRequest) (resp *favorite.DouyinFavoriteActionResponse, err error) {
	resp = new(favorite.DouyinFavoriteActionResponse)

	switch req.ActionType {
	case DoFavoriteAction:
		err = service.NewFavoriteActionService(ctx).FavoriteAction(req)
	case CancelFavoriteAction:
		err = service.NewFavoriteActionService(ctx).CancelFavoriteAction(req)
	default:
		err = errno.ParamErr
	}
	if err != nil {
		resp = pack.BuildFavoriteActionResp(err)
		return resp, err
	}

	resp = pack.BuildFavoriteActionResp(errno.Success)
	return resp, nil
}

// FavoriteList implements the FavoriteSrvImpl interface.
func (s *FavoriteSrvImpl) FavoriteList(ctx context.Context, req *favorite.DouyinFavoriteListRequest) (resp *favorite.DouyinFavoriteListResponse, err error) {
	resp = new(favorite.DouyinFavoriteListResponse)

	videos, err := service.NewFavoriteListService(ctx).FavoriteList(req)
	if err != nil {
		resp = pack.BuildFavoriteListResp(err)
		return resp, err
	}

	resp = pack.BuildFavoriteListResp(errno.Success)
	resp.VideoList = videos
	return resp, nil
}
