package service

import (
	"context"
	"tiktok/cmd/feed/dal/db"
	"tiktok/cmd/feed/pack"
	"tiktok/kitex_gen/feed"
	"tiktok/pkg/constants"

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
		klog.CtxFatalf(s.ctx, "db get video failed %v", err)
		return nil, 0, err
	}

	videos, nextTime, err := pack.Videos(s.ctx, vs, req.UserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "pack video failed %v", err)
		return nil, 0, err
	}

	return videos, nextTime, nil
}
