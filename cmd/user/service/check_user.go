package service

import (
	"context"
	"crypto/hmac"

	"tiktok/dal/mysql"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/errno"

	"github.com/cloudwego/kitex/pkg/klog"
)

type CheckUserService struct {
	ctx context.Context
}

// NewCheckUserService new CheckUserService
func NewCheckUserService(ctx context.Context) *CheckUserService {
	return &CheckUserService{
		ctx: ctx,
	}
}

// CheckUser check user info
func (s *CheckUserService) CheckUser(req *user.DouyinUserLoginRequest) (int64, error) {
	username := req.GetUsername()
	password := req.GetPassword()

	users, err := mysql.QueryUser(s.ctx, username)
	if err != nil {
		klog.CtxErrorf(s.ctx, "mysql query user failed %v", err)
		return 0, err
	}
	if len(users) == 0 {
		return 0, errno.UserNotExistErr
	}
	u := users[0]
	// 校验密码
	if !checkPasswordHash(username, password, u.Password) {
		return 0, errno.AuthorizationFailedErr
	}
	return u.Id, nil
}

func checkPasswordHash(username, password, hash string) bool {
	hashPwd := hashPassword(username, password)
	return hmac.Equal([]byte(hashPwd), []byte(hash))
}
