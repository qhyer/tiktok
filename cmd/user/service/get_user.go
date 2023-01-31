package service

import (
	"context"

	"tiktok/dal/neo4j"
	"tiktok/kitex_gen/relation"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/rpc"

	"github.com/cloudwego/kitex/pkg/klog"
)

type MGetUserService struct {
	ctx context.Context
}

// NewMGetUserService new MGetUserService
func NewMGetUserService(ctx context.Context) *MGetUserService {
	return &MGetUserService{ctx: ctx}
}

// MGetUser multiple get list of user info
func (s *MGetUserService) MGetUser(req *user.DouyinUserInfoRequest) ([]*user.User, error) {
	if len(req.ToUserIds) == 0 {
		return nil, nil
	}

	us, err := neo4j.MGetUserByUserIds(s.ctx, req.ToUserIds)
	if err != nil {
		klog.CtxErrorf(s.ctx, "neo4j get user failed %v", err)
		return nil, err
	}

	// 数据库结果存map 然后返回所有用户
	userMap := make(map[int64]*user.User, 0)
	users := make([]*user.User, 0, len(req.ToUserIds))
	for _, u := range us {
		userMap[u.Id] = u
	}
	for _, u := range req.ToUserIds {
		users = append(users, userMap[u])
	}

	// 获取当前用户与这些用户的关注关系
	followMap := make(map[int64]bool, 0)
	followResp, err := rpc.FollowList(s.ctx, &relation.DouyinRelationFollowListRequest{
		UserId:   req.UserId,
		ToUserId: req.UserId,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "rpc get follow list failed %v", err)
		return nil, err
	}
	if followResp != nil && followResp.UserList != nil {
		for _, u := range followResp.UserList {
			followMap[u.Id] = true
		}
	}
	for i, u := range users {
		users[i].IsFollow = followMap[u.Id]
	}

	return users, nil
}
