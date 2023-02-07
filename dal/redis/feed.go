package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"tiktok/dal/mysql"
	"tiktok/pkg/constants"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
)

func MAddVideoIdToFeed(ctx context.Context, videos []*mysql.Video) error {
	// 判断feed是否存在，不存在则读库创建列表
	err := updateFeed(ctx)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update feed failed %v", err)
		return err
	}

	if len(videos) == 0 {
		return nil
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
	err = RDB.ZAdd(ctx, constants.RedisFeedKey, videoIds...).Err()

	return err
}

func GetVideoIdsByLatestTime(ctx context.Context, latestTime int64, limit int64) ([]int64, error) {
	feedKey := constants.RedisFeedKey

	videoIds := make([]int64, 0, limit)
	// 判断feed是否存在，不存在则读库创建列表
	err := updateFeed(ctx)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update feed failed %v", err)
		return videoIds, err
	}

	// 查询视频列表
	res, err := RDB.ZRevRangeByScore(ctx, feedKey, &redis.ZRangeBy{
		Min:   "1", // 这里是时间戳 为0可能有为了避免缓存穿透放的空视频 因此置1
		Max:   fmt.Sprintf("%d", latestTime-1),
		Count: limit,
	}).Result()
	if err != nil {
		return videoIds, err
	}

	// 把视频id加入结果
	for _, v := range res {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			continue
		}
		videoIds = append(videoIds, id)
	}
	return videoIds, nil
}

func updateFeed(ctx context.Context) error {
	res, err := RDB.Exists(ctx, constants.RedisFeedKey).Result()
	if err != nil {
		return err
	}

	// 不存在feed
	if res == 0 {
		videoList, err := mysql.GetVideosByLatestTime(ctx, 100, time.Now().UnixMilli())
		if err != nil {
			klog.CtxErrorf(ctx, "mysql get latest videos failed %v", err)
			return err
		}

		// 库里也没视频 避免缓存穿透 放入空数据
		if len(videoList) == 0 {
			err = RDB.ZAdd(ctx, constants.RedisFeedKey, redis.Z{
				Score:  0,
				Member: 0,
			}).Err()
			if err != nil {
				klog.CtxErrorf(ctx, "redis add video ids to feed failed %v", err)
				return err
			}
			return nil
		}

		// 从列表中读视频id
		videoIds := make([]redis.Z, 0, len(videoList))
		for _, v := range videoList {
			if v == nil {
				continue
			}
			videoIds = append(videoIds, redis.Z{
				Score:  float64(v.CreatedAt.UnixMilli()),
				Member: v.Id,
			})
		}

		// 把视频加入缓存
		err = MSetVideoInfo(ctx, videoList)
		if err != nil {
			klog.CtxErrorf(ctx, "redis set video info failed %v", err)
			return err
		}

		// 把视频id加入缓存
		err = RDB.ZAdd(ctx, constants.RedisFeedKey, videoIds...).Err()
		if err != nil {
			klog.CtxErrorf(ctx, "redis add video ids to feed failed %v", err)
			return err
		}
	}
	return nil
}
