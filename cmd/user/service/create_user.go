package service

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"tiktok/dal/mysql"
	"tiktok/dal/pack"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/errno"

	"github.com/cloudwego/kitex/pkg/klog"
	"golang.org/x/crypto/pbkdf2"
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
	username := req.GetUsername()
	password := req.GetPassword()
	users, err := mysql.QueryUser(s.ctx, username)
	if err != nil {
		return 0, err
	}
	if len(users) != 0 {
		return 0, errno.UserAlreadyExistErr
	}

	// 对密码加盐
	password = hashPassword(username, password)

	// 创建用户
	us, err := mysql.CreateUser(s.ctx, []*mysql.User{{
		UserName: username,
		Password: password,
	}})
	if err != nil {
		klog.CtxErrorf(s.ctx, "mysql create user failed %v", err)
		return 0, err
	}

	if len(us) == 0 {
		return 0, errno.UserNotExistErr
	}

	usrs := pack.Users(us)

	// 获取创建成功后的用户id
	userId := usrs[0].Id
	return userId, nil
}

func hashPassword(username, password string) string {
	// 迭代次数
	iterations := 320000

	// 盐为用户名的md5
	salt := fmt.Sprintf("%x", md5.Sum([]byte(username)))

	// pbkdf2加密
	hash := pbkdf2.Key([]byte(password), []byte(salt), iterations, sha256.Size, sha256.New)

	// base64编码
	b64Hash := base64.StdEncoding.EncodeToString(hash)

	// 拼接最终密码
	return fmt.Sprintf("%s$%d$%s$%s", "pbkdf2_sha256", iterations, salt, b64Hash)
}
