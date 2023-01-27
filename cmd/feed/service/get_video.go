package service

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"tiktok/cmd/feed/dal/db"
	"tiktok/cmd/feed/pack"
	"tiktok/kitex_gen/feed"
)

type GetVideoService struct {
	ctx context.Context
}

// NewGetVideoService new GetVideoService
func NewGetVideoService(ctx context.Context) *GetVideoService {
	return &GetVideoService{ctx: ctx}
}

// GetVideosByVideoIdsAndCurrUserId get videos by video ids and current userid
func (s *GetVideoService) GetVideosByVideoIdsAndCurrUserId(req *feed.DouyinGetVideosByVideoIdsAndCurrentUserIdRequest) ([]*feed.Video, error) {
	vs, err := db.MGetVideosByVideoIds(s.ctx, req.VideoIds)
	if err != nil {
		klog.CtxFatalf(s.ctx, "db get video failed %v", err)
		return nil, err
	}

	videos, _, err := pack.Videos(s.ctx, vs, req.UserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "pack video failed %v", err)
		return nil, err
	}

	return videos, nil
}

// IsVideoIdsExist check is video id exist
func (s *GetVideoService) IsVideoIdsExist(req *feed.DouyinIsVideoIdsExistRequest) ([]bool, error) {
	vs, err := db.MGetVideosByVideoIds(s.ctx, req.VideoIds)
	if err != nil {
		klog.CtxFatalf(s.ctx, "db get video failed %v", err)
		return nil, err
	}

	// add id to map
	idMap := make(map[int64]bool, 0)
	for _, v := range vs {
		idMap[v.Id] = true
	}

	// check is video id exist
	res := make([]bool, 0)
	for _, v := range req.VideoIds {
		res = append(res, idMap[v])
	}

	return res, nil
}
