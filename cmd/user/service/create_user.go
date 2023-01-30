package service

import (
	"context"

	"tiktok/dal/mysql"
	"tiktok/dal/neo4j"
	"tiktok/dal/pack"
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
	users, err := mysql.QueryUser(s.ctx, req.Username)
	if err != nil {
		return 0, err
	}
	if len(users) != 0 {
		return 0, errno.UserAlreadyExistErr
	}

	// 对密码加盐
	passWord, err := HashPassword(req.Password)
	if err != nil {
		klog.CtxErrorf(s.ctx, "hash password failed %v", err)
		return 0, err
	}

	// TODO 延时队列check是否在neo4j中创建用户成功

	// 创建用户
	us, err := mysql.CreateUser(s.ctx, []*mysql.User{{
		UserName: req.Username,
		Password: passWord,
	}})
	if err != nil {
		klog.CtxErrorf(s.ctx, "mysql create user failed %v", err)
		return 0, err
	}

	if len(us) == 0 {
		return 0, errno.UserNotExistErr
	}

	usrs := pack.Users(us)

	// 在neo4j中创建用户节点
	err = neo4j.CreateUser(s.ctx, usrs[0])
	if err != nil {
		return 0, err
	}

	// 获取创建成功后的用户id
	userId := usrs[0].Id
	return userId, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
