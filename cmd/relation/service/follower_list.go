package service

import (
	"context"

	"tiktok/dal/neo4j"
	"tiktok/kitex_gen/relation"
	"tiktok/kitex_gen/user"

	"github.com/cloudwego/kitex/pkg/klog"
)

type FollowerListService struct {
	ctx context.Context
}

// NewFollowerListService new FollowerListService
func NewFollowerListService(ctx context.Context) *FollowerListService {
	return &FollowerListService{ctx: ctx}
}

// FollowerList get list of Follower user
func (s *FollowerListService) FollowerList(req *relation.DouyinRelationFollowerListRequest) ([]*user.User, error) {
	// 获取目标用户的粉丝
	users, err := neo4j.FollowerList(s.ctx, req.ToUserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "neo4j get follower list failed %v", err)
		return nil, err
	}

	// 获取当前用户的关注
	followList, err := neo4j.FollowList(s.ctx, req.UserId)
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
