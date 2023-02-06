package service

import (
	"context"

	"tiktok/cmd/rpc"
	"tiktok/dal/redis"
	"tiktok/kitex_gen/favorite"
	"tiktok/kitex_gen/feed"

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
	userId := req.GetUserId()
	toUserId := req.GetToUserId()

	// 从缓存中读列表，缓存中没有会读库
	fl, err := redis.GetFavoriteVideoIdsByUserId(s.ctx, toUserId)
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

	return feedResponse.GetVideoList(), nil
}

func (s *FavoriteListService) GetUserFavoriteVideoIds(req *favorite.DouyinGetUserFavoriteVideoIdsRequest) ([]int64, error) {
	userId := req.GetUserId()

	// 从缓存中读列表，缓存中没有会读库
	fl, err := redis.GetFavoriteVideoIdsByUserId(s.ctx, userId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "mysql get favorite list failed %v", err)
		return nil, err
	}

	return fl, err
}
