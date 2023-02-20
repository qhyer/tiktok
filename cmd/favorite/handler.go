package main

import (
	"context"

	"tiktok/cmd/favorite/service"
	"tiktok/dal/pack"
	"tiktok/kitex_gen/favorite"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
)

// FavoriteSrvImpl implements the last service interface defined in the IDL.
type FavoriteSrvImpl struct{}

// FavoriteAction implements the FavoriteSrvImpl interface.
func (s *FavoriteSrvImpl) FavoriteAction(ctx context.Context, req *favorite.DouyinFavoriteActionRequest) (resp *favorite.DouyinFavoriteActionResponse, err error) {
	resp = new(favorite.DouyinFavoriteActionResponse)

	switch req.GetActionType() {
	case constants.CreateFavoriteAction:
		err = service.NewFavoriteActionService(ctx).CreateFavorite(req)
	case constants.CancelFavoriteAction:
		err = service.NewFavoriteActionService(ctx).CancelFavorite(req)
	default:
		err = errno.ParamErr
	}
	if err != nil {
		resp = pack.BuildFavoriteActionResp(err)
		return resp, nil
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
		return resp, nil
	}

	resp = pack.BuildFavoriteListResp(errno.Success)
	resp.VideoList = videos
	return resp, nil
}

// GetUserFavoriteVideoIds implements the FavoriteSrvImpl interface.
func (s *FavoriteSrvImpl) GetUserFavoriteVideoIds(ctx context.Context, req *favorite.DouyinGetUserFavoriteVideoIdsRequest) (resp *favorite.DouyinGetUserFavoriteVideoIdsResponse, err error) {
	resp = new(favorite.DouyinGetUserFavoriteVideoIdsResponse)

	videoIds, err := service.NewFavoriteListService(ctx).GetUserFavoriteVideoIds(req)
	if err != nil {
		return resp, err
	}

	resp.VideoIds = videoIds
	return resp, nil
}
