package service

import (
	"context"

	"tiktok/dal/mysql"
	"tiktok/kitex_gen/feed"
	"tiktok/kitex_gen/publish"
	"tiktok/pkg/rpc"

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
	vs, err := mysql.GetPublishedVideoIdsByUserId(s.ctx, req.ToUserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "mysql get video failed %v", err)
		return nil, err
	}

	videoResponse, _ := rpc.GetVideosByVideoIdsAndCurrentUserId(s.ctx, &feed.DouyinGetVideosByVideoIdsAndCurrentUserIdRequest{
		UserId:   req.UserId,
		VideoIds: vs,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "rpc get video failed %v", err)
	}

	return videoResponse.VideoList, nil
}
