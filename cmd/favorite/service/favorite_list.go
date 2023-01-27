package service

import (
	"context"
	"tiktok/cmd/favorite/dal/db"
	"tiktok/cmd/favorite/pack"
	"tiktok/cmd/favorite/rpc"
	"tiktok/kitex_gen/favorite"
	"tiktok/kitex_gen/feed"
	"tiktok/pkg/errno"

	"github.com/cloudwego/kitex/pkg/klog"
)

type FavoriteListService struct {
	ctx context.Context
}

// NewFavoriteListService new GetVideoService
func NewFavoriteListService(ctx context.Context) *FavoriteListService {
	return &FavoriteListService{ctx: ctx}
}

// FavoriteList get user favorite list
func (s *FavoriteListService) FavoriteList(req *favorite.DouyinFavoriteListRequest) ([]*feed.Video, error) {
	userId := req.UserId
	fl, err := db.FavoriteList(s.ctx, userId)
	if err != nil {
		klog.CtxFatalf(s.ctx, "db get favorite list failed %v", err)
		return nil, err
	}

	favoriteList := pack.Favorites(fl)

	videoIds := make([]int64, 0, len(favoriteList))
	for _, v := range favoriteList {
		videoIds = append(videoIds, v.VideoId)
	}

	// rpc通信
	feedResponse, err := rpc.GetVideosByVideoIdsAndCurrentUserId(s.ctx, &feed.DouyinGetVideosByVideoIdsAndCurrentUserIdRequest{
		UserId:   userId,
		VideoIds: videoIds,
	})
	if err != nil || feedResponse.StatusCode != errno.SuccessCode {
		klog.CtxErrorf(s.ctx, "rpc get video list failed %v", err)
		return nil, err
	}
	return feedResponse.VideoList, nil
}
