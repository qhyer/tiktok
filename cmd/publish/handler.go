package main

import (
	"context"

	"tiktok/cmd/publish/service"
	"tiktok/dal/pack"
	"tiktok/kitex_gen/publish"
	"tiktok/pkg/errno"
)

// PublishSrvImpl implements the last service interface defined in the IDL.
type PublishSrvImpl struct{}

// PublishAction implements the PublishSrvImpl interface.
func (s *PublishSrvImpl) PublishAction(ctx context.Context, req *publish.DouyinPublishActionRequest) (resp *publish.DouyinPublishActionResponse, err error) {
	resp = new(publish.DouyinPublishActionResponse)

	err = service.NewPublishActionService(ctx).PublishVideo(req)
	if err != nil {
		resp = pack.BuildPublishActionResp(err)
		return resp, nil
	}

	resp = pack.BuildPublishActionResp(errno.Success)
	return
}

// PublishList implements the PublishSrvImpl interface.
func (s *PublishSrvImpl) PublishList(ctx context.Context, req *publish.DouyinPublishListRequest) (resp *publish.DouyinPublishListResponse, err error) {
	resp = new(publish.DouyinPublishListResponse)

	videos, err := service.NewPublishListService(ctx).PublishList(req)
	if err != nil {
		resp = pack.BuildPublishListResp(err)
		return resp, nil
	}

	resp = pack.BuildPublishListResp(errno.Success)
	resp.VideoList = videos
	return resp, nil
}
