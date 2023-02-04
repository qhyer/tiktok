package redis

import (
	"context"
	"fmt"
	"time"

	"tiktok/dal/mysql"
	"tiktok/dal/pack"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/constants"

	"github.com/redis/go-redis/v9"
	"github.com/u2takey/go-utils/rand"
)

func GetUserInfoByUserId(ctx context.Context, userId int64) (*user.User, error) {
	userKey := fmt.Sprintf(constants.RedisUserKey, userId)

	var us mysql.User
	_, err := RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		err := pipeliner.HMGet(ctx, userKey, "id", "name", "follow_count", "follower_count").Scan(&us)
		if err != nil {
			return err
		}
		pipeliner.Expire(ctx, userKey, constants.UserInfoExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second)
		return nil
	})
	if err != nil {
		return nil, err
	}

	usr := pack.User(&us)
	return usr, nil
}

func MGetUserInfoByUserId(ctx context.Context, userIds []int64) (users []*user.User, notInCacheUserIds []int64) {
	for _, uid := range userIds {
		res, err := GetUserInfoByUserId(ctx, uid)
		if err != nil || res == nil {
			notInCacheUserIds = append(notInCacheUserIds, uid)
		} else {
			users = append(users, res)
		}
	}
	return
}

func MSetUserInfo(ctx context.Context, users []*user.User) error {
	_, err := RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for _, usr := range users {
			if usr == nil {
				continue
			}
			userKey := fmt.Sprintf(constants.RedisUserKey, usr.GetId())
			pipeliner.HMSet(ctx, userKey,
				"id", usr.GetId(),
				"name", usr.GetName(),
				"follow_count", usr.GetFollowCount(),
				"follower_count", usr.GetFollowerCount(),
			)
			pipeliner.Expire(ctx, userKey, constants.UserInfoExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second)
		}
		return nil
	})
	return err
}
