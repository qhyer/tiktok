package redis

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"tiktok/dal/mysql"
	"tiktok/pkg/constants"

	"github.com/redis/go-redis/v9"
)

func MGetVideoInfoByVideoId(ctx context.Context, videoIds []int64) (videos []*mysql.Video, notInCacheVideoIds []int64) {
	res, err := RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for _, vid := range videoIds {
			videoKey := fmt.Sprintf(constants.RedisVideoKey, vid)
			pipeliner.HGetAll(ctx, videoKey)
			pipeliner.Expire(ctx, videoKey, constants.VideoInfoExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second)
		}
		return nil
	})
	if err != nil {
		return nil, videoIds
	}
	// 处理getall结果
	for i, l := range res {
		if i%2 == 0 {
			var v mysql.Video
			err := MustScan(l.(*redis.MapStringStringCmd), &v)
			if err != nil {
				notInCacheVideoIds = append(notInCacheVideoIds, videoIds[i/2])
				continue
			}
			videos = append(videos, &v)
		}
	}

	return
}

func MSetVideoInfo(ctx context.Context, videos []*mysql.Video) error {
	_, err := RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for _, video := range videos {
			if video == nil {
				continue
			}
			videoKey := fmt.Sprintf(constants.RedisVideoKey, video.Id)
			pipeliner.HMSet(ctx, videoKey,
				"id", video.Id,
				"author_user_id", video.AuthorUserId,
				"play_url", video.PlayUrl,
				"cover_url", video.CoverUrl,
				"favorite_count", video.FavoriteCount,
				"comment_count", video.CommentCount,
				"title", video.Title,
				"created_at", video.CreatedAt.UnixMilli(),
			)
			pipeliner.Expire(ctx, videoKey, constants.VideoInfoExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second)
		}
		return nil
	})
	return err
}
