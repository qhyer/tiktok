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
	userId := req.GetUserId()
	toUserIds := req.GetToUserIds()
	if len(toUserIds) == 0 {
		return nil, nil
	}

	us, err := neo4j.MGetUserByUserIds(s.ctx, toUserIds)
	if err != nil {
		klog.CtxErrorf(s.ctx, "neo4j get user failed %v", err)
		return nil, err
	}

	// 数据库结果存map 然后返回所有用户
	userMap := make(map[int64]*user.User, 0)
	users := make([]*user.User, 0, len(toUserIds))
	for _, u := range us {
		if u == nil {
			continue
		}
		userMap[u.Id] = u
	}
	for _, u := range toUserIds {
		users = append(users, userMap[u])
	}

	// 获取当前用户与这些用户的关注关系
	followMap := make(map[int64]bool, 0)
	followResp, err := rpc.FollowList(s.ctx, &relation.DouyinRelationFollowListRequest{
		UserId:   userId,
		ToUserId: userId,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "rpc get follow list failed %v", err)
		return nil, err
	}
	followList := followResp.GetUserList()
	for _, u := range followList {
		if u == nil {
			continue
		}
		followMap[u.Id] = true
	}
	for i, u := range users {
		if u == nil {
			continue
		}
		users[i].IsFollow = followMap[u.Id]
	}

	return users, nil
}
