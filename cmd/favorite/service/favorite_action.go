package service

import (
	"context"

	"tiktok/dal/db"
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
	err := db.FavoriteAction(s.ctx, &db.Favorite{
		UserId:  req.UserId,
		VideoId: req.VideoId,
	})
	return err
}

// CancelFavoriteAction cancel favorite video action
func (s *FavoriteActionService) CancelFavoriteAction(req *favorite.DouyinFavoriteActionRequest) error {
	err := db.CancelFavoriteAction(s.ctx, &db.Favorite{
		UserId:  req.UserId,
		VideoId: req.VideoId,
	})
	return err
}
