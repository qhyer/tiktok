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

func AddNewFavoriteToFavoriteList(ctx context.Context, favorite *mysql.Favorite) error {
	if favorite == nil {
		return nil
	}

	userId := favorite.UserId
	videoId := favorite.VideoId

	// 判断喜欢列表是否存在，不存在则创建列表
	err := updateFavoriteList(ctx, userId)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update favorite list failed %v", err)
		return err
	}

	favoriteListKey := fmt.Sprintf(constants.RedisFavoriteListKey, userId)
	videoKey := fmt.Sprintf(constants.RedisVideoKey, videoId)
	_, err = RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		// 把视频加入喜欢列表
		if err := RDB.ZAdd(ctx, favoriteListKey, redis.Z{
			Score:  float64(favorite.CreatedAt.UnixMilli()),
			Member: videoId,
		}).Err(); err != nil {
			return err
		}

		// 喜欢数+1
		incrBy := redis.NewScript(`
					if redis.call("Exists", KEYS[1]) > 0 then
						redis.call("HIncrBy", KEYS[1], "favorite_count", 1)
					end
					return true 
					`)
		keys := []string{videoKey}
		if err := incrBy.Run(ctx, RDB, keys).Err(); err != nil {
			return err
		}
		return nil
	})
	return nil
}

func DeleteFavoriteFromFavoriteList(ctx context.Context, favorite *mysql.Favorite) error {
	if favorite == nil {
		return nil
	}
	videoId := favorite.VideoId
	userId := favorite.UserId

	// 判断喜欢列表是否存在，不存在则创建列表
	err := updateFavoriteList(ctx, videoId)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update favorite list failed %v", err)
		return err
	}

	favoriteListKey := fmt.Sprintf(constants.RedisFavoriteListKey, userId)
	videoKey := fmt.Sprintf(constants.RedisVideoKey, videoId)
	_, err = RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		// 把视频id从list中删除
		if err := RDB.ZRem(ctx, favoriteListKey, videoId).Err(); err != nil {
			return err
		}

		// 喜欢数-1
		incrBy := redis.NewScript(`
					if redis.call("Exists", KEYS[1]) > 0 then
						redis.call("HIncrBy", KEYS[1], "favorite_count", -1)
					end
					return true
					`)
		keys := []string{videoKey}
		if err := incrBy.Run(ctx, RDB, keys).Err(); err != nil {
			return err
		}
		return nil
	})
	return err
}

func GetFavoriteVideoIdsByUserId(ctx context.Context, userId int64) ([]int64, error) {
	favoriteListKey := fmt.Sprintf(constants.RedisFavoriteListKey, userId)
	videoIds := make([]int64, 0)

	err := updateFavoriteList(ctx, userId)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update favorite list failed %v", err)
		return videoIds, err
	}

	// 把视频id加入结果
	res, err := RDB.ZRevRangeByScore(ctx, favoriteListKey, &redis.ZRangeBy{
		Min: "1", // 这里是时间戳 0的位置可能有为了避免缓存穿透放的空视频 因此置1
		Max: fmt.Sprintf("%d", time.Now().UnixMilli()),
	}).Result()
	if err != nil {
		return videoIds, err
	}
	for _, f := range res {
		vid, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			continue
		}
		videoIds = append(videoIds, vid)
	}

	return videoIds, nil
}

func updateFavoriteList(ctx context.Context, userId int64) error {
	favoriteListKey := fmt.Sprintf(constants.RedisFavoriteListKey, userId)
	exist, err := RDB.Exists(ctx, favoriteListKey).Result()
	if err != nil {
		return err
	}

	// 不存在列表 更新
	if exist == 0 {
		videos, err := mysql.GetFavoriteVideoIdsByUserId(ctx, userId)
		if err != nil {
			klog.CtxErrorf(ctx, "mysql get favorite list failed %v", err)
			return err
		}

		// 数据库没数据 避免缓存穿透 放入空数据
		if len(videos) == 0 {
			err = RDB.ZAdd(ctx, favoriteListKey, redis.Z{
				Score:  0,
				Member: 0,
			}).Err()
			if err != nil {
				klog.CtxErrorf(ctx, "redis add favorite list failed %v", err)
				return err
			}

			// 设置list的过期时间
			err = RDB.Expire(ctx, favoriteListKey, constants.FavoriteListExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second).Err()
			if err != nil {
				klog.CtxErrorf(ctx, "redis set favorite list expiry failed %v", err)
				return err
			}
			return nil
		}

		// 从列表中读视频id
		videoZs := make([]redis.Z, 0, len(videos))
		for _, f := range videos {
			if f == nil {
				continue
			}
			videoZs = append(videoZs, redis.Z{
				Score:  float64(f.CreatedAt.UnixMilli()),
				Member: f.VideoId,
			})
		}

		// 把视频id加入缓存
		err = RDB.ZAdd(ctx, favoriteListKey, videoZs...).Err()
		if err != nil {
			klog.CtxErrorf(ctx, "redis add favorite list failed %v", err)
			return err
		}

		// 设置list的过期时间
		err = RDB.Expire(ctx, favoriteListKey, constants.FavoriteListExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second).Err()
		if err != nil {
			klog.CtxErrorf(ctx, "redis set favorite list expiry failed %v", err)
			return err
		}

		return nil
	}
	return nil
}
