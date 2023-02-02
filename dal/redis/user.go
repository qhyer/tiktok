package redis

import (
	"context"
	"fmt"

	"tiktok/kitex_gen/user"
	"tiktok/pkg/constants"
)

func GetUserInfoByUserId(ctx context.Context, userId int64) (user.User, error) {
	userKey := fmt.Sprintf(constants.RedisUserKey, userId)

}
