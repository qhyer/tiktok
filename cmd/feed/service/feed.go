package service

import (
	"context"

	"tiktok/dal/db"
	"tiktok/dal/pack"
	"tiktok/kitex_gen/feed"
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
	vs, err := db.GetVideosByLatestTime(s.ctx, constants.VideoQueryLimit, *req.LatestTime)
	if err != nil {
		klog.CtxErrorf(s.ctx, "db get video failed %v", err)
		return nil, 0, err
	}

	videos, nextTime := pack.Videos(vs)

	videos, err = minio.SignFeed(s.ctx, videos)
	if err != nil {
		klog.CtxErrorf(s.ctx, "minio sign feed failed %v", err)
		return nil, 0, err
	}

	// TODO add user
	return videos, nextTime, nil
}
