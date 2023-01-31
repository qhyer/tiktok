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
	vs, err := mysql.GetVideosByLatestTime(s.ctx, constants.VideoQueryLimit, *req.LatestTime)
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
		UserId:    req.UserId,
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

	// 查询用户点赞
	favoriteResp, err := rpc.FavoriteList(s.ctx, &favorite.DouyinFavoriteListRequest{
		UserId:   req.UserId,
		ToUserId: req.UserId,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "rpc get user favorite list failed %v", err)
		return nil, 0, err
	}
	favoriteMap := make(map[int64]bool, 0)
	for _, f := range favoriteResp.VideoList {
		favoriteMap[f.Id] = true
	}

	return videos, nextTime, nil
}
