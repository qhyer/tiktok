package service

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"tiktok/cmd/favorite/dal/db"
	"tiktok/cmd/favorite/rpc"
	"tiktok/kitex_gen/favorite"
	"tiktok/kitex_gen/feed"
	"tiktok/pkg/errno"
)

type FavoriteActionService struct {
	ctx context.Context
}

// NewFavoriteActionService new FavoriteActionService
func NewFavoriteActionService(ctx context.Context) *FavoriteActionService {
	return &FavoriteActionService{ctx: ctx}
}

// FavoriteAction user do favorite video action
func (s *FavoriteActionService) FavoriteAction(req *favorite.DouyinFavoriteActionRequest) error {
	// rpc通信
	feedResponse, err := rpc.IsVideoIdsExist(s.ctx, &feed.DouyinIsVideoIdsExistRequest{
		VideoIds: []int64{req.VideoId},
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "rpc check video ids failed %v", err)
		return err
	}

	// videoId 不存在
	if !feedResponse.IsExist[0] {
		return errno.ParamErr
	}

	err = db.FavoriteAction(s.ctx, &db.Favorite{
		UserId:  req.UserId,
		VideoId: req.VideoId,
	})
	return err
}
