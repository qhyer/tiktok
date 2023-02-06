package redis

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"tiktok/dal/neo4j"
	"tiktok/kitex_gen/relation"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/constants"

	"github.com/cloudwego/kitex/pkg/klog"
)

func GetFollowListByUserId(ctx context.Context, userId int64) ([]*user.User, error) {
	follows := make([]*user.User, 0)
	followListKey := fmt.Sprintf(constants.RedisFollowListKey, userId)

	// 判断关注列表是否存在，不存在则读库
	err := updateFollowList(ctx, userId)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update follow list failed %v", err)
		return follows, err
	}

	// 获取关注的用户id
	followIdStrs, err := RDB.SMembers(ctx, followListKey).Result()
	if err != nil {
		klog.CtxErrorf(ctx, "redis get follow list failed %v", err)
		return follows, err
	}

	if followIdStrs == nil {
		return follows, nil
	}

	followIds := make([]int64, 0, len(followIdStrs))
	for _, f := range followIdStrs {
		uid, err := strconv.ParseInt(f, 10, 64)
		if err == nil {
			continue
		}
		followIds = append(followIds, uid)
	}

	// 从缓存中读用户列表
	follows, err = MGetUserInfoByUserId(ctx, followIds)
	if err != nil {
		klog.CtxErrorf(ctx, "redis get userinfo failed %v", err)
		return follows, err
	}
	return follows, nil
}

func GetFollowerListByUserId(ctx context.Context, userId int64) ([]*user.User, error) {
	followers := make([]*user.User, 0)
	followerListKey := fmt.Sprintf(constants.RedisFollowerListKey, userId)

	// 判断粉丝列表是否存在，不存在则读库
	err := updateFollowerList(ctx, userId)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update follower list failed %v", err)
		return followers, err
	}

	// 获取粉丝的用户id
	followerIdStrs, err := RDB.SMembers(ctx, followerListKey).Result()
	if err != nil {
		klog.CtxErrorf(ctx, "redis get follower list failed %v", err)
		return followers, err
	}

	if followerIdStrs == nil {
		return followers, nil
	}

	followerIds := make([]int64, 0, len(followerIdStrs))
	for _, f := range followerIdStrs {
		uid, err := strconv.ParseInt(f, 10, 64)
		if err == nil {
			continue
		}
		followerIds = append(followerIds, uid)
	}

	// 从缓存中读用户列表
	followers, err = MGetUserInfoByUserId(ctx, followerIds)
	if err != nil {
		klog.CtxErrorf(ctx, "redis get userinfo failed %v", err)
		return followers, err
	}
	return followers, nil
}

func GetFriendListByUserId(ctx context.Context, userId int64) ([]*relation.FriendUser, error) {
	friends := make([]*relation.FriendUser, 0)
	friendListKey := fmt.Sprintf(constants.RedisFriendListKey, userId)

	// 判断好友列表是否存在，不存在则读库
	err := updateFriendList(ctx, userId)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update friend list failed %v", err)
		return friends, err
	}

	// 获取好友的用户id
	friendIdStrs, err := RDB.SMembers(ctx, friendListKey).Result()
	if err != nil {
		klog.CtxErrorf(ctx, "redis get friend list failed %v", err)
		return friends, err
	}

	if friendIdStrs == nil {
		return friends, nil
	}

	friendIds := make([]int64, 0, len(friendIdStrs))
	for _, f := range friendIdStrs {
		uid, err := strconv.ParseInt(f, 10, 64)
		if err == nil {
			continue
		}
		friendIds = append(friendIds, uid)
	}

	// 从缓存中读用户列表
	fs, err := MGetUserInfoByUserId(ctx, friendIds)
	if err != nil {
		klog.CtxErrorf(ctx, "redis get userinfo failed %v", err)
		return friends, err
	}

	// 封装
	for _, f := range fs {
		if f == nil {
			continue
		}
		friends = append(friends, &relation.FriendUser{
			Id:            f.Id,
			Name:          f.Name,
			FollowCount:   f.FollowCount,
			FollowerCount: f.FollowerCount,
			IsFollow:      true,
			Avatar:        constants.DefaultAvatarUrl, // 使用默认头像
		})
	}

	return friends, nil
}

func updateFollowList(ctx context.Context, userId int64) error {
	followListKey := fmt.Sprintf(constants.RedisFollowListKey, userId)
	res, err := RDB.Exists(ctx, followListKey).Result()
	if err != nil {
		return err
	}

	// 不存在关注列表 查库
	if res == 0 {
		followList, err := neo4j.FollowList(ctx, userId)
		if err != nil {
			klog.CtxErrorf(ctx, "neo4j get follow list failed %v", err)
			return err
		}

		// 数据库中无关注
		if len(followList) == 0 {
			return nil
		}

		// 读取关注用户id
		followIds := make([]int64, 0, len(followList))
		for _, f := range followList {
			followIds = append(followIds, f.Id)
		}

		// 把被关注的用户信息加入缓存
		err = MSetUserInfo(ctx, followList)
		if err != nil {
			klog.CtxErrorf(ctx, "redis set follow userinfo failed %v", err)
			return err
		}

		// 把被关注的用户id加入缓存
		err = RDB.SAdd(ctx, followListKey, followIds).Err()
		if err != nil {
			klog.CtxErrorf(ctx, "redis add follow id list failed %v", err)
			return err
		}

		// 设置list的过期时间
		err = RDB.Expire(ctx, followListKey, constants.FollowListExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second).Err()
		if err != nil {
			klog.CtxErrorf(ctx, "redis set follow list expiry failed %v", err)
		}
	}
	return nil
}

func updateFollowerList(ctx context.Context, userId int64) error {
	followerListKey := fmt.Sprintf(constants.RedisFollowerListKey, userId)
	res, err := RDB.Exists(ctx, followerListKey).Result()
	if err != nil {
		return err
	}

	// 不存在粉丝列表 查库
	if res == 0 {
		followerList, err := neo4j.FollowerList(ctx, userId)
		if err != nil {
			klog.CtxErrorf(ctx, "neo4j get follower list failed %v", err)
			return err
		}

		// 数据库中无粉丝
		if len(followerList) == 0 {
			return nil
		}

		// 读取粉丝用户id
		followerIds := make([]int64, 0, len(followerList))
		for _, f := range followerList {
			followerIds = append(followerIds, f.Id)
		}

		// 把粉丝的用户信息加入缓存
		err = MSetUserInfo(ctx, followerList)
		if err != nil {
			klog.CtxErrorf(ctx, "redis set follower userinfo failed %v", err)
			return err
		}

		// 把粉丝用户id加入缓存
		err = RDB.SAdd(ctx, followerListKey, followerIds).Err()
		if err != nil {
			klog.CtxErrorf(ctx, "redis add follower id list failed %v", err)
			return err
		}

		// 设置list的过期时间
		err = RDB.Expire(ctx, followerListKey, constants.FollowerListExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second).Err()
		if err != nil {
			klog.CtxErrorf(ctx, "redis set follower list expiry failed %v", err)
		}
	}
	return nil
}

func updateFriendList(ctx context.Context, userId int64) error {
	friendListKey := fmt.Sprintf(constants.RedisFriendListKey, userId)
	res, err := RDB.Exists(ctx, friendListKey).Result()
	if err != nil {
		return err
	}

	// 不存在好友列表 查库
	if res == 0 {
		friendList, err := neo4j.FriendList(ctx, userId)
		if err != nil {
			klog.CtxErrorf(ctx, "neo4j get friend list failed %v", err)
			return err
		}

		// 数据库中无好友
		if len(friendList) == 0 {
			return nil
		}

		// 读取好友用户id
		friendIds := make([]int64, 0, len(friendList))
		for _, f := range friendList {
			friendIds = append(friendIds, f.Id)
		}

		// 把好友的用户信息加入缓存
		friendUsers := make([]*user.User, 0, len(friendList))
		for _, f := range friendList {
			if f == nil {
				continue
			}
			friendUsers = append(friendUsers, &user.User{
				Id:            f.Id,
				Name:          f.Name,
				FollowCount:   f.FollowCount,
				FollowerCount: f.FollowerCount,
			})
		}
		err = MSetUserInfo(ctx, friendUsers)
		if err != nil {
			klog.CtxErrorf(ctx, "redis set friend userinfo failed %v", err)
			return err
		}

		// 把好友的用户id加入缓存
		err = RDB.SAdd(ctx, friendListKey, friendIds).Err()
		if err != nil {
			klog.CtxErrorf(ctx, "redis add friend id list failed %v", err)
			return err
		}

		// 设置list的过期时间
		err = RDB.Expire(ctx, friendListKey, constants.FriendListExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second).Err()
		if err != nil {
			klog.CtxErrorf(ctx, "redis set friend list expiry failed %v", err)
			return err
		}
	}
	return nil
}
