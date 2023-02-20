package service

import (
	"context"

	"tiktok/cmd/rpc"
	"tiktok/dal/redis"
	"tiktok/kitex_gen/feed"
	"tiktok/kitex_gen/publish"
	"tiktok/pkg/errno"

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

	videoList := make([]*feed.Video, 0)
	if len(videoIds) == 0 {
		return videoList, nil
	}

	videoResponse, err := rpc.GetVideosByVideoIdsAndCurrentUserId(s.ctx, &feed.DouyinGetVideosByVideoIdsAndCurrentUserIdRequest{
		UserId:   userId,
		VideoIds: videoIds,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "rpc get video list failed %v", err)
		return nil, err
	}
	if videoResponse.GetStatusCode() != errno.SuccessCode {
		klog.CtxErrorf(s.ctx, "rpc get video list failed %v", videoResponse.GetStatusMsg())
		return nil, errno.NewErrNo(videoResponse.GetStatusCode(), videoResponse.GetStatusMsg())
	}
	videoList = videoResponse.GetVideoList()

	return videoList, nil
}
