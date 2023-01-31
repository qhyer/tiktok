package service

import (
	"context"

	"tiktok/dal/mysql"
	"tiktok/dal/pack"
	"tiktok/kitex_gen/favorite"
	"tiktok/kitex_gen/feed"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/constants"
	"tiktok/pkg/minio"
	"tiktok/pkg/rpc"

	"github.com/cloudwego/kitex/pkg/klog"
)

type FeedService struct {
	ctx context.Context
}

// NewFeedService new FeedService
func NewFeedService(ctx context.Context) *FeedService {
	return &FeedService{ctx: ctx}
}

// Feed get list of video
func (s *FeedService) Feed(req *feed.DouyinFeedRequest) ([]*feed.Video, int64, error) {
	latestTime := req.GetLatestTime()
	userId := req.GetUserId()

	vs, err := mysql.GetVideosByLatestTime(s.ctx, constants.VideoQueryLimit, latestTime)
	if err != nil {
		klog.CtxErrorf(s.ctx, "mysql get video failed %v", err)
		return nil, 0, err
	}

	videos, nextTime := pack.Videos(vs)

	if len(videos) == 0 {
		return nil, 0, nil
	}

	// 对链接签名
	videos, err = minio.SignFeed(s.ctx, videos)
	if err != nil {
		klog.CtxErrorf(s.ctx, "minio sign feed failed %v", err)
		return nil, 0, err
	}

	// 查询用户信息
	userIds := make([]int64, 0, len(videos))
	for _, v := range videos {
		userIds = append(userIds, v.Author.Id)
	}

	users, err := rpc.UserInfo(s.ctx, &user.DouyinUserInfoRequest{
		UserId:    userId,
		ToUserIds: userIds,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "rpc get userinfo failed %v", err)
		return nil, 0, err
	}

	// 加入用户信息
	for i := range videos {
		if users.User[i] == nil {
			klog.CtxWarnf(s.ctx, "video author is nil")
			continue
		}
		videos[i].Author = users.User[i]
	}

	// 查询用户点赞视频
	favoriteResp, err := rpc.GetUserFavoriteVideoIds(s.ctx, &favorite.DouyinGetUserFavoriteVideoIdsRequest{
		UserId: userId,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "rpc get user favorite list failed %v", err)
		return nil, 0, err
	}
	favoriteMap := make(map[int64]bool, 0)
	if favoriteResp != nil && favoriteResp.VideoIds != nil {
		for _, f := range favoriteResp.VideoIds {
			favoriteMap[f] = true
		}
	}
	for i, v := range videos {
		videos[i].IsFavorite = favoriteMap[v.Id]
	}

	return videos, nextTime, nil
}
