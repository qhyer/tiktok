package service

import (
	"context"

	"tiktok/cmd/rpc"
	"tiktok/dal/pack"
	"tiktok/dal/redis"
	"tiktok/kitex_gen/favorite"
	"tiktok/kitex_gen/feed"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/minio"

	"github.com/cloudwego/kitex/pkg/klog"
)

type GetVideoService struct {
	ctx context.Context
}

// NewGetVideoService new GetVideoService
func NewGetVideoService(ctx context.Context) *GetVideoService {
	return &GetVideoService{ctx: ctx}
}

// GetVideosByVideoIdsAndCurrUserId get videos by video ids and current userid
func (s *GetVideoService) GetVideosByVideoIdsAndCurrUserId(req *feed.DouyinGetVideosByVideoIdsAndCurrentUserIdRequest) ([]*feed.Video, error) {
	userId := req.GetUserId()
	videoIds := req.GetVideoIds()
	if len(videoIds) == 0 {
		return nil, nil
	}

	videoMap := make(map[int64]*feed.Video, 0)
	videos := make([]*feed.Video, 0, len(videoIds))
	// 缓存中查视频详情
	vs, err := redis.MGetVideoInfoByVideoId(s.ctx, videoIds)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis get video info failed %v", err)
		return videos, err
	}
	redisVideos, _ := pack.Videos(vs)
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

	// 给链接签名
	videos, err = minio.SignFeed(s.ctx, videos)
	if err != nil {
		klog.CtxErrorf(s.ctx, "minio sign feed failed %v", err)
		return nil, err
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
		return nil, err
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
		return nil, err
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

	return videos, nil
}
