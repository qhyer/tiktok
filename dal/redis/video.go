package redis

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"tiktok/dal/mysql"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"

	"github.com/redis/go-redis/v9"
)

func GetVideoInfoByVideoId(ctx context.Context, videoId int64) (*mysql.Video, error) {
	videoKey := fmt.Sprintf(constants.RedisVideoKey, videoId)

	var v mysql.Video
	_, err := RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		err := pipeliner.HMGet(ctx, videoKey, "id", "author_user_id", "play_url", "cover_url",
			"favorite_count", "comment_count", "title", "created_at").Scan(&v)
		if err != nil {
			return err
		}
		pipeliner.Expire(ctx, videoKey, constants.UserInfoExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func SetVideoInfo(ctx context.Context, video *mysql.Video) error {
	videoId := video.Id

	if videoId == 0 {
		return errno.ParamErr
	}

	videoKey := fmt.Sprintf(constants.RedisVideoKey, videoId)

	_, err := RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		pipeliner.HMSet(ctx, videoKey,
			"id", video.Id,
			"author_user_id", video.AuthorUserId,
			"play_url", video.PlayUrl,
			"favorite_count", video.FavoriteCount,
			"commment_count", video.CommentCount,
			"title", video.Title,
			"created_at", video.CreatedAt.UnixMilli(),
		)
		pipeliner.Expire(ctx, videoKey, constants.VideoInfoExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second)
		return nil
	})
	return err
}

func MGetVideoInfoByVideoId(ctx context.Context, videoIds []int64) (videos []*mysql.Video, notInCacheVideoIds []int64) {
	for _, vid := range videoIds {
		res, err := GetVideoInfoByVideoId(ctx, vid)
		if err != nil || res == nil {
			notInCacheVideoIds = append(notInCacheVideoIds, vid)
		} else {
			videos = append(videos, res)
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
