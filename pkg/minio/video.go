package minio

import (
	"context"
	"net/url"

	"tiktok/kitex_gen/feed"
	"tiktok/pkg/constants"

	"github.com/cloudwego/kitex/pkg/klog"
)

func SignFeed(ctx context.Context, videos []*feed.Video) ([]*feed.Video, error) {
	for i, v := range videos {
		playUrl := v.PlayUrl
		coverUrl := v.CoverUrl
		reqParams := make(url.Values)
		videoInfo, err := PreSignedGetObject(ctx, constants.VideoBucketName, playUrl, constants.OSSDefaultExpiry, reqParams)
		if err != nil {
			klog.CtxErrorf(ctx, "pre sign get object failed %v", err)
			return nil, err
		}
		coverInfo, err := PreSignedGetObject(ctx, constants.CoverBucketName, coverUrl, constants.OSSDefaultExpiry, reqParams)
		if err != nil {
			klog.CtxErrorf(ctx, "pre sign get object failed %v", err)
			return nil, err
		}
		playUrl = constants.OSSBaseUrl + videoInfo.Path + "?" + videoInfo.RawQuery
		coverUrl = constants.OSSBaseUrl + coverInfo.Path + "?" + coverInfo.RawQuery
		videos[i].PlayUrl = playUrl
		videos[i].CoverUrl = coverUrl
	}
	return videos, nil
}
