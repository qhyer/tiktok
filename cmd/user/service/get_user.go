package service

import (
	"context"
	"tiktok/cmd/user/dal/db"
	"tiktok/cmd/user/pack"
	"tiktok/kitex_gen/user"
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
		return nil, err
	}

	// TODO 当前用户和被查询用户的关系

	return pack.Users(modelUsers), nil
}
