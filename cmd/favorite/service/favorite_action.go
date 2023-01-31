package service

import (
	"context"

	"tiktok/dal/mysql"
	"tiktok/kitex_gen/favorite"
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
	userId := req.GetUserId()
	videoId := req.GetVideoId()

	err := mysql.FavoriteAction(s.ctx, &mysql.Favorite{
		UserId:  userId,
		VideoId: videoId,
	})
	return err
}

// CancelFavoriteAction cancel favorite video action
func (s *FavoriteActionService) CancelFavoriteAction(req *favorite.DouyinFavoriteActionRequest) error {
	userId := req.GetUserId()
	videoId := req.GetVideoId()

	err := mysql.CancelFavoriteAction(s.ctx, &mysql.Favorite{
		UserId:  userId,
		VideoId: videoId,
	})
	return err
}
