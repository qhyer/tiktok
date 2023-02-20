package service

import (
	"context"

	"tiktok/cmd/rpc"
	"tiktok/dal/redis"
	"tiktok/kitex_gen/relation"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/errno"

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

	userMap := make(map[int64]*user.User, 0)
	// 读取用户
	us, err := redis.MGetUserInfoByUserId(s.ctx, toUserIds)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis get userinfo failed %v", err)
		return nil, err
	}
	for _, u := range us {
		if u == nil {
			continue
		}
		userMap[u.Id] = u
	}

	// 返回所有用户
	users := make([]*user.User, 0, len(toUserIds))
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
	if followResp.GetStatusCode() != errno.SuccessCode {
		klog.CtxErrorf(s.ctx, "rpc get follow list failed %v", followResp.GetStatusMsg())
		return nil, errno.NewErrNo(followResp.GetStatusCode(), followResp.GetStatusMsg())
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
