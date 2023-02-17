package service

import (
	"context"

	"tiktok/dal/neo4j"
	"tiktok/dal/redis"
	"tiktok/kitex_gen/relation"
	"tiktok/pkg/errno"

	"github.com/cloudwego/kitex/pkg/klog"
)

type RelationActionService struct {
	ctx context.Context
}

// NewRelationActionService new RelationActionService
func NewRelationActionService(ctx context.Context) *RelationActionService {
	return &RelationActionService{ctx: ctx}
}

// Follow user follow action
func (s *RelationActionService) Follow(req *relation.DouyinRelationActionRequest) error {
	userId := req.GetUserId()
	toUserId := req.GetToUserId()

	// 在数据库中关注
	err := neo4j.FollowAction(s.ctx, userId, toUserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "neo4j follow action failed %v", err)
		return errno.DatabaseOperationFailedErr
	}

	// 缓存中 操作用户关注数+1 被关注者粉丝+1 同时修改列表
	err = redis.AddNewFollow(s.ctx, userId, toUserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis add new follow failed %v", err)
		return err
	}

	return nil
}

// Unfollow user unfollow action
func (s *RelationActionService) Unfollow(req *relation.DouyinRelationActionRequest) error {
	userId := req.GetUserId()
	toUserId := req.GetToUserId()

	// 在数据库中取关
	err := neo4j.UnfollowAction(s.ctx, userId, toUserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "neo4j unfollow action failed %v", err)
		return errno.DatabaseOperationFailedErr
	}

	// 缓存中 操作用户关注数-1 被关注者粉丝-1 同时修改列表
	err = redis.Unfollow(s.ctx, userId, toUserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis unfollow failed %v", err)
		return err
	}

	return nil
}
