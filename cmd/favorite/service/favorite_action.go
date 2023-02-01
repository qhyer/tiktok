package service

import (
	"context"

	"tiktok/dal/mysql"
	"tiktok/kitex_gen/favorite"

	"github.com/cloudwego/kitex/pkg/klog"
)

type FavoriteActionService struct {
	ctx context.Context
}

// NewFavoriteActionService new FavoriteActionService
func NewFavoriteActionService(ctx context.Context) *FavoriteActionService {
	return &FavoriteActionService{ctx: ctx}
}

// CreateFavorite user do favorite video action
func (s *FavoriteActionService) CreateFavorite(req *favorite.DouyinFavoriteActionRequest) error {
	userId := req.GetUserId()
	videoId := req.GetVideoId()

	err := mysql.CreateFavorite(s.ctx, &mysql.Favorite{
		UserId:  userId,
		VideoId: videoId,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "db create favorite failed %v", err)
		return err
	}

	return nil
}

// CancelFavorite cancel favorite video action
func (s *FavoriteActionService) CancelFavorite(req *favorite.DouyinFavoriteActionRequest) error {
	userId := req.GetUserId()
	videoId := req.GetVideoId()

	err := mysql.DeleteFavorite(s.ctx, &mysql.Favorite{
		UserId:  userId,
		VideoId: videoId,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "db delete favorite failed %v", err)
		return err
	}

	return nil
}
