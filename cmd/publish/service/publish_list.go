package service

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"tiktok/cmd/publish/dal/db"
	"tiktok/cmd/publish/pack"
	"tiktok/kitex_gen/publish"
)

type PublishListService struct {
	ctx context.Context
}

// NewPublishListService new PublishService
func NewPublishListService(ctx context.Context) *PublishListService {
	return &PublishListService{ctx: ctx}
}

// PublishList get list of video
func (s *PublishListService) PublishList(req *publish.DouyinPublishListRequest) ([]*publish.Video, error) {
	vs, err := db.GetPublishedVideosByUserId(s.ctx, req.ToUserId)
	if err != nil {
		klog.Errorf("db get video failed %v", err)
		return nil, err
	}

	videos, err := pack.Videos(s.ctx, vs, req.UserId)
	if err != nil {
		klog.Errorf("pack video failed %v", err)
		return nil, err
	}

	return videos, nil
}
