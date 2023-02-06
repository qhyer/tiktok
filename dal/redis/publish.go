package redis

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"tiktok/dal/mysql"
	"tiktok/pkg/constants"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
)

func GetPublishedVideoIdsByUserId(ctx context.Context, userId int64) ([]int64, error) {
	publishListKey := fmt.Sprintf(constants.RedisPublishListKey, userId)
	videoIds := make([]int64, 0)

	// 判断列表是否存在，不存在则读库创建列表
	err := updatePublishList(ctx, userId)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update publish list failed %v", err)
		return videoIds, err
	}

	// 查询视频列表
	res, err := RDB.ZRevRange(ctx, publishListKey, 0, -1).Result()
	if err != nil {
		klog.CtxErrorf(ctx, "redis get publish list failed %v", err)
	}
	for _, v := range res {
		vid, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}
		videoIds = append(videoIds, vid)
	}
	return videoIds, nil
}

func MAddVideoIdToPublishList(ctx context.Context, videos []*mysql.Video, userId int64) error {
	publishListKey := fmt.Sprintf(constants.RedisPublishListKey, userId)
	// 判断列表是否存在，不存在则读库创建列表
	err := updatePublishList(ctx, userId)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update publish list failed %v", err)
		return err
	}
	// 加入视频
	videoIds := make([]redis.Z, 0, len(videos))
	for _, v := range videos {
		if v == nil {
			continue
		}
		videoIds = append(videoIds, redis.Z{
			Score:  float64(v.CreatedAt.UnixMilli()),
			Member: v.Id,
		})
	}
	err = RDB.ZAdd(ctx, publishListKey, videoIds...).Err()
	return err
}

func updatePublishList(ctx context.Context, userId int64) error {
	publishListKey := fmt.Sprintf(constants.RedisPublishListKey, userId)
	res, err := RDB.Exists(ctx, publishListKey).Result()
	if err != nil {
		return err
	}

	// 不存在发布视频列表 查库
	if res == 0 {
		publishList, err := mysql.GetPublishedVideosByUserId(ctx, userId)
		if err != nil {
			klog.CtxErrorf(ctx, "mysql get user published video ids failed %v", err)
			return err
		}

		// 数据库也没有视频列表
		if len(publishList) == 0 {
			return nil
		}

		// 把视频加入缓存
		err = MSetVideoInfo(ctx, publishList)
		if err != nil {
			klog.CtxErrorf(ctx, "redis set video info failed %v", err)
			return err
		}

		// 把视频id加入缓存
		videoIds := make([]redis.Z, 0, len(publishList))
		for _, v := range publishList {
			if v == nil {
				continue
			}
			videoIds = append(videoIds, redis.Z{
				Score:  float64(v.CreatedAt.UnixMilli()),
				Member: v.Id,
			})
		}
		err = RDB.ZAdd(ctx, publishListKey, videoIds...).Err()
		if err != nil {
			klog.CtxErrorf(ctx, "redis add publish list failed %v", err)
			return err
		}

		// 设置list的过期时间
		err = RDB.Expire(ctx, publishListKey, constants.PublishListExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second).Err()
		if err != nil {
			klog.CtxErrorf(ctx, "redis set publish list expire failed %v", err)
			return err
		}
		return nil
	}
	return nil
}
