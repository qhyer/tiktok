package service

import (
	"context"

	"tiktok/dal/redis"
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
	userId := req.GetUserId()
	toUserId := req.GetToUserId()

	// 获取目标用户的关注
	users, err := redis.GetFollowListByUserId(s.ctx, toUserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis get follow list failed %v", err)
		return nil, err
	}

	// 当前用户和目标用户相同
	if userId == toUserId {
		return users, nil
	}

	// 获取当前用户的关注
	followList, err := redis.GetFollowListByUserId(s.ctx, userId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis get follow list failed %v", err)
		return nil, err
	}
	userFollowMap := make(map[int64]bool, 0)
	for _, u := range followList {
		if u == nil {
			continue
		}
		userFollowMap[u.Id] = true
	}
	// 设置粉丝列表中当前用户的关注关系
	for i, u := range users {
		if u == nil {
			continue
		}
		users[i].IsFollow = userFollowMap[u.Id]
	}

	return users, nil
}
