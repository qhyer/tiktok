package service

import (
	"context"
	"tiktok/cmd/feed/dal/db"
	"tiktok/cmd/feed/pack"
	"tiktok/kitex_gen/feed"
	"tiktok/pkg/constants"
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
	vs, err := db.MGetVideos(s.ctx, constants.VideoQueryLimit, *req.LatestTime)
	if err != nil {
		return nil, 0, err
	}

	videos, nextTime := pack.Videos(s.ctx, vs, req.UserId)

	return videos, nextTime, nil
}
