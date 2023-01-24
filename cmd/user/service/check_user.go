package service

import (
	"context"
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
		return 0, err
	}
	if len(users) == 0 {
		return 0, errno.AuthorizationFailedErr
	}
	u := users[0]
	if !CheckPasswordHash(req.Password, u.Password) {
		return 0, errno.AuthorizationFailedErr
	}
	return int64(u.Id), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
