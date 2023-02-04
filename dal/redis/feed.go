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

func AddVideoIdToFeed(ctx context.Context, video *mysql.Video) {
	// 判断feed是否为空，为空则读库创建列表
	err := updateFeed(ctx)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update feed failed %v", err)
	}

	// 加入当前视频
	err = RDB.ZAdd(ctx, constants.RedisFeedKey, redis.Z{
		Score:  float64(video.CreatedAt.UnixMilli()),
		Member: video.Id,
	}).Err()

	return
}

func MAddVideoIdToFeed(ctx context.Context, videos []*mysql.Video) error {
	// 判断feed是否为空，为空则读库创建列表
	err := updateFeed(ctx)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update feed failed %v", err)
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

	// 判断feed是否为空，为空则读库创建列表
	err := updateFeed(ctx)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update feed failed %v", err)
	}

	// 查询视频列表
	res, err := RDB.ZRevRangeByScore(ctx, feedKey, &redis.ZRangeBy{
		Min:    "0",
		Max:    fmt.Sprintf("%d", latestTime),
		Offset: 0,
		Count:  limit,
	}).Result()
	if err != nil {
		return nil, err
	}

	// 把视频id加入结果
	videoIds := make([]int64, 0, len(res))
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
	count, err := RDB.Exists(ctx, constants.RedisFeedKey).Result()
	if err != nil {
		return err
	}

	// 不存在feed
	if count == 0 {
		videoList, err := mysql.GetVideosByLatestTime(ctx, 100, time.Now().UnixMilli())
		if err != nil {
			klog.CtxErrorf(ctx, "mysql get latest videos failed %v", err)
			return err
		}

		// 库里也没视频
		if len(videoList) == 0 {
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
			klog.CtxErrorf(ctx, "redis add video info failed %v", err)
		}

		// 把视频id加入缓存
		err = RDB.ZAdd(ctx, constants.RedisFeedKey, videoIds...).Err()
		return err
	}
	return nil
}
