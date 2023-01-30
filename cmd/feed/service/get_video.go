package service

import (
	"context"

	"tiktok/dal/mysql"
	"tiktok/dal/pack"
	"tiktok/kitex_gen/feed"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/minio"
	"tiktok/pkg/rpc"

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
	if len(req.VideoIds) == 0 {
		return nil, nil
	}

	vs, err := mysql.MGetVideosByVideoIds(s.ctx, req.VideoIds)
	if err != nil {
		klog.CtxErrorf(s.ctx, "mysql get video failed %v", err)
		return nil, err
	}

	if len(vs) == 0 {
		return nil, nil
	}

	videos, _ := pack.Videos(vs)

	// 给链接签名
	videos, err = minio.SignFeed(s.ctx, videos)
	if err != nil {
		klog.CtxErrorf(s.ctx, "minio sign feed failed %v", err)
		return nil, err
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
		return nil, err
	}

	// 加入用户信息
	for i := range videos {
		videos[i].Author = users.User[i]
	}

	return videos, nil
}
