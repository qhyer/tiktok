package service

import (
	"context"

	"tiktok/cmd/rpc"
	"tiktok/dal/pack"
	"tiktok/dal/redis"
	"tiktok/kitex_gen/favorite"
	"tiktok/kitex_gen/feed"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/constants"
	"tiktok/pkg/minio"

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

	// 从缓存中读视频id列表 其中没读到缓存会查库
	videoIds, err := redis.GetVideoIdsByLatestTime(s.ctx, latestTime, constants.VideoQueryLimit)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis get latest video ids failed %v", err)
		return nil, 0, err
	}

	videoMap := make(map[int64]*feed.Video, 0)
	videos := make([]*feed.Video, 0)
	nextTime := latestTime

	// 缓存中查视频详情
	rvs, err := redis.MGetVideoInfoByVideoId(s.ctx, videoIds)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis get video info failed %v", err)
		return videos, 0, err
	}
	redisVideos, nts := pack.Videos(rvs)
	if nts < nextTime {
		nextTime = nts
	}
	for _, v := range redisVideos {
		if v == nil {
			continue
		}
		videoMap[v.Id] = v
	}

	// 合并视频
	for _, i := range videoIds {
		res := videoMap[i]
		if res == nil {
			continue
		}
		videos = append(videos, res)
	}
	if len(videos) == 0 {
		return videos, 0, nil
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
		if v == nil || v.Author == nil {
			continue
		}
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

	us := users.GetUser()
	// 加入用户信息
	for i := range videos {
		res := us[i]
		if res == nil {
			klog.CtxWarnf(s.ctx, "video author is nil")
			continue
		}
		videos[i].Author = res
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
	vids := favoriteResp.GetVideoIds()
	for _, f := range vids {
		favoriteMap[f] = true
	}
	for i, v := range videos {
		if v == nil {
			continue
		}
		videos[i].IsFavorite = favoriteMap[v.Id]
	}

	return videos, nextTime, nil
}
