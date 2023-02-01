package service

import (
	"context"

	rpc2 "tiktok/cmd/rpc"
	"tiktok/dal/mysql"
	"tiktok/dal/pack"
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

	vs, err := mysql.MGetVideosByVideoIds(s.ctx, videoIds)
	if err != nil {
		klog.CtxErrorf(s.ctx, "mysql get video failed %v", err)
		return nil, err
	}

	videos, _ := pack.Videos(vs)
	if len(videos) == 0 {
		return videos, nil
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

	users, err := rpc2.UserInfo(s.ctx, &user.DouyinUserInfoRequest{
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
		if us[i] == nil {
			klog.CtxWarnf(s.ctx, "video author is nil")
			continue
		}
		videos[i].Author = us[i]
	}

	// 查询用户点赞视频
	favoriteResp, err := rpc2.GetUserFavoriteVideoIds(s.ctx, &favorite.DouyinGetUserFavoriteVideoIdsRequest{
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
