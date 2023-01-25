package service

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"golang.org/x/crypto/bcrypt"
	"tiktok/cmd/user/dal/db"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/errno"
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
	userName := req.Username
	users, err := db.QueryUser(s.ctx, userName)
	if err != nil {
		klog.Warnf("db query user failed %v", err)
		return 0, err
	}
	if len(users) == 0 {
		return 0, errno.AuthorizationFailedErr
	}
	u := users[0]
	if !CheckPasswordHash(req.Password, u.Password) {
		return 0, errno.AuthorizationFailedErr
	}
	return u.Id, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
