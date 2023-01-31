package service

import (
	"context"

	"tiktok/dal/mysql"
	"tiktok/kitex_gen/favorite"
	"tiktok/kitex_gen/feed"
	"tiktok/pkg/rpc"

	"github.com/cloudwego/kitex/pkg/klog"
)

type FavoriteListService struct {
	ctx context.Context
}

// NewFavoriteListService new FavoriteListService
func NewFavoriteListService(ctx context.Context) *FavoriteListService {
	return &FavoriteListService{ctx: ctx}
}

// FavoriteList get user favorite list
func (s *FavoriteListService) FavoriteList(req *favorite.DouyinFavoriteListRequest) ([]*feed.Video, error) {
	userId := req.UserId
	toUserId := req.ToUserId
	fl, err := mysql.GetFavoriteVideoIdsByUserId(s.ctx, toUserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "mysql get favorite list failed %v", err)
		return nil, err
	}

	if len(fl) == 0 {
		return nil, nil
	}

	// rpc通信
	feedResponse, err := rpc.GetVideosByVideoIdsAndCurrentUserId(s.ctx, &feed.DouyinGetVideosByVideoIdsAndCurrentUserIdRequest{
		UserId:   userId,
		VideoIds: fl,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "rpc get video list failed %v", err)
		return nil, err
	}

	return feedResponse.VideoList, nil
}

func (s *FavoriteListService) GetUserFavoriteVideoIds(req *favorite.DouyinGetUserFavoriteVideoIdsRequest) ([]int64, error) {
	fl, err := mysql.GetFavoriteVideoIdsByUserId(s.ctx, req.UserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "mysql get favorite list failed %v", err)
		return nil, err
	}

	return fl, err
}
