package service

import (
	"context"

	"tiktok/dal/db"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/errno"

	"github.com/cloudwego/kitex/pkg/klog"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserService struct {
	ctx context.Context
}

// NewCreateUserService new CreateUserService
func NewCreateUserService(ctx context.Context) *CreateUserService {
	return &CreateUserService{ctx: ctx}
}

// CreateUser create user info.
func (s *CreateUserService) CreateUser(req *user.DouyinUserRegisterRequest) (int64, error) {
	users, err := db.QueryUser(s.ctx, req.Username)
	if err != nil {
		return 0, err
	}
	if len(users) != 0 {
		return 0, errno.UserAlreadyExistErr
	}

	passWord, err := HashPassword(req.Password)
	if err != nil {
		klog.CtxErrorf(s.ctx, "hash password failed %v", err)
		return 0, err
	}

	err = db.CreateUser(s.ctx, []*db.User{{
		UserName: req.Username,
		Password: passWord,
	}})
	if err != nil {
		klog.CtxErrorf(s.ctx, "db create user failed %v", err)
		return 0, err
	}

	users, err = db.QueryUser(s.ctx, req.Username)
	if err != nil {
		klog.CtxErrorf(s.ctx, "db query user failed %v", err)
		return 0, err
	}
	if len(users) == 0 {
		return 0, errno.UserNotExistErr
	}

	usr := users[0]
	return usr.Id, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
