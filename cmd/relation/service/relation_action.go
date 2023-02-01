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

// Follow user follow action
func (s *RelationActionService) Follow(req *relation.DouyinRelationActionRequest) error {
	userId := req.GetUserId()
	toUserId := req.GetToUserId()

	err := neo4j.FollowAction(s.ctx, userId, toUserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "neo4j follow action failed %v", err)
		return err
	}
	return nil
}

// Unfollow user unfollow action
func (s *RelationActionService) Unfollow(req *relation.DouyinRelationActionRequest) error {
	userId := req.GetUserId()
	toUserId := req.GetToUserId()

	err := neo4j.UnfollowAction(s.ctx, userId, toUserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "neo4j unfollow action failed %v", err)
		return err
	}
	return nil
}
