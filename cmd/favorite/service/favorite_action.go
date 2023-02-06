package service

import (
	"context"

	"tiktok/dal/mysql"
	"tiktok/dal/redis"
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

	// 数据库中创建喜欢
	fav, err := mysql.CreateFavorite(s.ctx, &mysql.Favorite{
		UserId:  userId,
		VideoId: videoId,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "db create favorite failed %v", err)
		return err
	}

	// 缓存中加入视频 喜欢数+1
	err = redis.AddNewFavoriteToFavoriteList(s.ctx, fav)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis update favorite num failed %v", err)
		return err
	}

	return nil
}

// CancelFavorite cancel favorite video action
func (s *FavoriteActionService) CancelFavorite(req *favorite.DouyinFavoriteActionRequest) error {
	userId := req.GetUserId()
	videoId := req.GetVideoId()

	// 数据库中删除喜欢
	fav, err := mysql.DeleteFavorite(s.ctx, &mysql.Favorite{
		UserId:  userId,
		VideoId: videoId,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "db delete favorite failed %v", err)
		return err
	}

	// 缓存中移除视频 喜欢数-1
	err = redis.DeleteFavoriteFromFavoriteList(s.ctx, fav)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis update favorite num failed %v", err)
		return err
	}

	return nil
}
