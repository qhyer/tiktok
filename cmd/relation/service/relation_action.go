package service

import (
	"context"

	"tiktok/dal/neo4j"
	"tiktok/kitex_gen/relation"

	"github.com/cloudwego/kitex/pkg/klog"
)

type RelationActionService struct {
	ctx context.Context
}

// NewRelationActionService new RelationActionService
func NewRelationActionService(ctx context.Context) *RelationActionService {
	return &RelationActionService{ctx: ctx}
}

// FollowAction user follow action
func (s *RelationActionService) FollowAction(req *relation.DouyinRelationActionRequest) error {
	userId := req.GetUserId()
	toUserId := req.GetToUserId()

	err := neo4j.FollowAction(s.ctx, userId, toUserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "neo4j follow action failed %v", err)
		return err
	}
	return nil
}

// UnFollowAction user unfollow action
func (s *RelationActionService) UnFollowAction(req *relation.DouyinRelationActionRequest) error {
	userId := req.GetUserId()
	toUserId := req.GetToUserId()

	err := neo4j.UnfollowAction(s.ctx, userId, toUserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "neo4j unfollow action failed %v", err)
		return err
	}
	return nil
}
