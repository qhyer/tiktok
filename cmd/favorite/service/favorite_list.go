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
	fl, err := mysql.GetFavoriteVideoIdsByUserId(s.ctx, userId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "mysql get favorite list failed %v", err)
		return nil, err
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
