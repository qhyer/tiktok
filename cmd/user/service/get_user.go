package service

import (
	"context"

	"tiktok/dal/db"
	"tiktok/dal/pack"
	"tiktok/kitex_gen/user"

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
	modelUsers, err := db.MGetUsers(s.ctx, req.ToUserIds)
	if err != nil {
		klog.CtxErrorf(s.ctx, "db get multiple users failed %v", err)
		return nil, err
	}

	// TODO 当前用户和被查询用户的关系

	return pack.Users(modelUsers), nil
}
