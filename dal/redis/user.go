package redis

import (
	"context"
	"fmt"
	"time"

	"tiktok/dal/mysql"
	"tiktok/dal/neo4j"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/constants"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
	"github.com/u2takey/go-utils/rand"
)

func MGetUserInfoByUserId(ctx context.Context, userIds []int64) (users []*user.User, err error) {
	res, err := RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for _, uid := range userIds {
			userKey := fmt.Sprintf(constants.RedisUserKey, uid)
			pipeliner.HGetAll(ctx, userKey)
			pipeliner.Expire(ctx, userKey, constants.UserInfoExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second)
		}
		return nil
	})
	if err != nil {
		return users, err
	}
	notInCacheUserIds := make([]int64, 0)

	// 处理getall结果
	for i, l := range res {
		if i%2 == 0 {
			var us mysql.User
			err := MustScan(l.(*redis.MapStringStringCmd), &us)
			if err != nil {
				notInCacheUserIds = append(notInCacheUserIds, userIds[i/2])
				continue
			}
			usr := &user.User{
				Id:            userIds[i/2],
				Name:          us.UserName,
				FollowCount:   &us.FollowCount,
				FollowerCount: &us.FollowerCount,
			}
			users = append(users, usr)
		}
	}

	// 缓存没找到 查库
	if len(notInCacheUserIds) > 0 {
		us, err := neo4j.MGetUserByUserIds(ctx, notInCacheUserIds)
		if err != nil {
			klog.CtxErrorf(ctx, "neo4j get user failed %v", err)
			return users, err
		}

		// 添加查库结果
		users = append(users, us...)

		// 查库结果加入缓存
		err = MSetUserInfo(ctx, us)
		if err != nil {
			klog.CtxErrorf(ctx, "redis set userinfo failed %v", err)
			return users, err
		}

		// 防止缓存穿透 这里没查到的放空用户
		if len(us) != len(notInCacheUserIds) {
			existUserIds := make(map[int64]bool)
			for _, u := range us {
				if u == nil {
					continue
				}
				existUserIds[u.Id] = true
			}
			emptyUsers := make([]*user.User, 0)
			for _, uid := range notInCacheUserIds {
				if !existUserIds[uid] {
					emptyUsers = append(emptyUsers, &user.User{
						Id: uid,
					})
				}
			}
			err = MSetUserInfo(ctx, emptyUsers)
			if err != nil {
				klog.CtxErrorf(ctx, "redis set userinfo failed %v", err)
				return users, err
			}
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
