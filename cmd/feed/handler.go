package main

import (
	"context"
	"tiktok/cmd/feed/pack"
	"tiktok/cmd/feed/service"
	feed "tiktok/kitex_gen/feed"
	"tiktok/pkg/errno"
)

// FeedSrvImpl implements the last service interface defined in the IDL.
type FeedSrvImpl struct{}

// Feed implements the FeedSrvImpl interface.
func (s *FeedSrvImpl) Feed(ctx context.Context, req *feed.DouyinFeedRequest) (resp *feed.DouyinFeedResponse, err error) {
	resp = new(feed.DouyinFeedResponse)

	videos, nextTime, err := service.NewFeedService(ctx).Feed(req)
	if err != nil {
		resp = pack.BuildFeedResp(err)
		return resp, nil
	}

	resp = pack.BuildFeedResp(errno.Success)
	resp.VideoList = videos
	resp.NextTime = &nextTime
	return resp, nil
}
