package service

import (
	"context"

	"tiktok/dal/neo4j"
	"tiktok/kitex_gen/relation"
)

type QueryRelationService struct {
	ctx context.Context
}

// NewQueryRelationService new QueryRelationService
func NewQueryRelationService(ctx context.Context) *QueryRelationService {
	return &QueryRelationService{ctx: ctx}
}

// IsFriend query is friend relation
func (s *QueryRelationService) IsFriend(req *relation.DouyinRelationIsFriendRequest) (bool, error) {
	uid1 := req.GetUserId()
	uid2 := req.GetToUserId()

	res, err := neo4j.IsFriend(s.ctx, uid1, uid2)
	if err != nil {
		return false, err
	}

	return res, nil
}
