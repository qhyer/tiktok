package service

import (
	"context"

	"tiktok/dal/redis"
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
	userId := req.GetUserId()
	toUserId := req.GetToUserId()

	// 获取目标用户的粉丝
	users, err := redis.GetFollowerListByUserId(s.ctx, toUserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis get follower list failed %v", err)
		return nil, err
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
