package redis

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"tiktok/dal/mysql"
	"tiktok/pkg/constants"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
)

func MSetComment(ctx context.Context, comments []*mysql.Comment) error {
	_, err := RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for _, comment := range comments {
			if comment == nil {
				continue
			}
			commentKey := fmt.Sprintf(constants.RedisCommentKey, comment.Id)
			pipeliner.HMSet(ctx, commentKey,
				"id", comment.Id,
				"user_id", comment.UserId,
				"content", comment.Content,
			)
			pipeliner.Expire(ctx, commentKey, constants.CommentExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second)
		}
		return nil
	})
	return err
}

func GetCommentIdListByVideoId(ctx context.Context, videoId int64) ([]*mysql.Comment, error) {
	// 判断commentList是否存在 不存在则读库创建列表
	err := updateCommentList(ctx, videoId)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update comment list failed %v", err)
	}

	commentListKey := fmt.Sprintf(constants.RedisCommentListKey, videoId)

	// 查询评论id列表
	res, err := RDB.ZRevRangeByScoreWithScores(ctx, commentListKey, &redis.ZRangeBy{}).Result()
	if err != nil {
		return nil, err
	}

	// 把评论id和时间加入结果
	commentList := make([]*mysql.Comment, 0, len(res))
	for _, c := range res {
		com := &mysql.Comment{
			Id:        c.Member.(int64),
			CreatedAt: time.UnixMilli(int64(c.Score)),
		}
		commentList = append(commentList, com)
	}
	return commentList, nil
}

func updateCommentList(ctx context.Context, videoId int64) error {
	commentListKey := fmt.Sprintf(constants.RedisCommentListKey, videoId)
	res, err := RDB.Exists(ctx, commentListKey).Result()
	if err != nil {
		return err
	}

	// 不存在评论列表
	if res == 0 {
		commentList, err := mysql.GetCommentListByVideoId(ctx, videoId)
		if err != nil {
			klog.CtxErrorf(ctx, "mysql get comment list failed %v", err)
			return err
		}

		// 数据库也没评论
		if len(commentList) == 0 {
			return nil
		}

		// 从列表中读评论id
		commentIds := make([]redis.Z, 0, len(commentList))
		for _, c := range commentList {
			if c == nil {
				continue
			}
			commentIds = append(commentIds, redis.Z{
				Score:  float64(c.CreatedAt.UnixMilli()),
				Member: c.Id,
			})
		}

		// 把评论加入缓存
		err = MSetComment(ctx, commentList)
		if err != nil {
			klog.CtxErrorf(ctx, "redis set comment failed %v", err)
			return err
		}

		// 把评论id加入缓存
		err = RDB.ZAdd(ctx, commentListKey, commentIds...).Err()
		return err
	}
	return nil
}
