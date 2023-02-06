package service

import (
	"context"

	"tiktok/cmd/rpc"
	"tiktok/dal/redis"
	"tiktok/kitex_gen/feed"
	"tiktok/kitex_gen/publish"

	"github.com/cloudwego/kitex/pkg/klog"
)

type PublishListService struct {
	ctx context.Context
}

// NewPublishListService new PublishService
func NewPublishListService(ctx context.Context) *PublishListService {
	return &PublishListService{ctx: ctx}
}

// PublishList get list of video
func (s *PublishListService) PublishList(req *publish.DouyinPublishListRequest) ([]*feed.Video, error) {
	userId := req.GetUserId()
	toUserId := req.GetToUserId()

	// 查缓存，缓存中没有会读库
	videoIds, err := redis.GetPublishedVideoIdsByUserId(s.ctx, toUserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis get video failed %v", err)
		return nil, err
	}

	if len(videoIds) == 0 {
		return nil, nil
	}

	videoResponse, err := rpc.GetVideosByVideoIdsAndCurrentUserId(s.ctx, &feed.DouyinGetVideosByVideoIdsAndCurrentUserIdRequest{
		UserId:   userId,
		VideoIds: videoIds,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "rpc get video failed %v", err)
		return nil, err
	}

	return videoResponse.GetVideoList(), nil
}
