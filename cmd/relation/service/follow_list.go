package service

import (
	"context"

	"tiktok/dal/neo4j"
	"tiktok/kitex_gen/relation"
	"tiktok/kitex_gen/user"

	"github.com/cloudwego/kitex/pkg/klog"
)

type FollowListService struct {
	ctx context.Context
}

// NewFollowListService new FollowListService
func NewFollowListService(ctx context.Context) *FollowListService {
	return &FollowListService{ctx: ctx}
}

// FollowList get list of follow user
func (s *FollowListService) FollowList(req *relation.DouyinRelationFollowListRequest) ([]*user.User, error) {
	// 获取目标用户的关注
	users, err := neo4j.FollowList(s.ctx, req.ToUserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "neo4j get follow list failed %v", err)
		return nil, err
	}

	// 当前用户和目标用户相同
	if req.UserId == req.ToUserId {
		for i := range users {
			users[i].IsFollow = true
		}
		return users, nil
	}

	// 获取当前用户的关注
	followList, err := neo4j.FollowList(s.ctx, req.UserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "neo4j get follow list failed %v", err)
		return nil, err
	}
	userFollowMap := make(map[int64]bool, 0)
	for _, u := range followList {
		userFollowMap[u.Id] = true
	}

	// 设置粉丝列表中当前用户的关注关系
	for i, u := range users {
		users[i].IsFollow = userFollowMap[u.Id]
	}

	return users, err
}
