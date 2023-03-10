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

func AddNewCommentToCommentList(ctx context.Context, comment *mysql.Comment) error {
	if comment == nil {
		return nil
	}
	videoId := comment.VideoId

	// 判断评论列表是否存在，不存在则创建列表
	err := updateCommentList(ctx, videoId)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update comment list failed %v", err)
		return err
	}

	commentListKey := fmt.Sprintf(constants.RedisCommentListKey, videoId)
	videoKey := fmt.Sprintf(constants.RedisVideoKey, videoId)
	_, err = RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		// 把评论加入list
		if err := RDB.ZAdd(ctx, commentListKey, redis.Z{
			Score:  float64(comment.CreatedAt.UnixMilli()),
			Member: comment.Id,
		}).Err(); err != nil {
			return err
		}
		// 评论数+1
		incrBy := redis.NewScript(`
					if redis.call("Exists", KEYS[1]) > 0 then
						redis.call("HIncrBy", KEYS[1], "comment_count", 1)
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

func DeleteCommentFromCommentList(ctx context.Context, comment *mysql.Comment) error {
	if comment == nil {
		return nil
	}
	videoId := comment.VideoId

	// 判断评论列表是否存在，不存在则创建列表
	err := updateCommentList(ctx, videoId)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update comment list failed %v", err)
		return err
	}

	commentListKey := fmt.Sprintf(constants.RedisCommentListKey, videoId)
	videoKey := fmt.Sprintf(constants.RedisVideoKey, videoId)
	_, err = RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		// 把评论从list中删除
		if err := RDB.ZRem(ctx, commentListKey, comment.Id).Err(); err != nil {
			return err
		}

		// 评论数-1
		incrBy := redis.NewScript(`
					if redis.call("Exists", KEYS[1]) > 0 then
						redis.call("HIncrBy", KEYS[1], "comment_count", -1)
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

func SetComment(ctx context.Context, comment *mysql.Comment) error {
	if comment == nil {
		return nil
	}
	_, err := RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		commentKey := fmt.Sprintf(constants.RedisCommentKey, comment.Id)
		pipeliner.HMSet(ctx, commentKey,
			"user_id", comment.UserId,
			"content", comment.Content,
		)
		pipeliner.Expire(ctx, commentKey, constants.CommentExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second)
		return nil
	})
	return err
}

func MSetComment(ctx context.Context, comments []*mysql.Comment) error {
	_, err := RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for _, c := range comments {
			if c == nil {
				continue
			}
			commentKey := fmt.Sprintf(constants.RedisCommentKey, c.Id)
			pipeliner.HMSet(ctx, commentKey,
				"user_id", c.UserId,
				"content", c.Content,
			)
			pipeliner.Expire(ctx, commentKey, constants.CommentExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second)
		}
		return nil
	})
	return err
}

func MGetCommentByCommentId(ctx context.Context, redisComments []*mysql.Comment) (comments []*mysql.Comment, err error) {
	notInCacheCommentIds := make([]int64, 0)
	res, err := RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for _, c := range redisComments {
			commentKey := fmt.Sprintf(constants.RedisCommentKey, c.Id)
			pipeliner.HGetAll(ctx, commentKey)
			pipeliner.Expire(ctx, commentKey, constants.CommentExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second)
		}
		return nil
	})
	if err != nil {
		klog.CtxErrorf(ctx, "redis get comment failed %v", err)
		return nil, err
	}

	// 处理getall结果
	for i, l := range res {
		if i%2 == 0 {
			var c mysql.Comment
			err := MustScan(l.(*redis.MapStringStringCmd), &c)
			if err != nil {
				notInCacheCommentIds = append(notInCacheCommentIds, redisComments[i/2].Id)
				continue
			}
			com := &mysql.Comment{
				Id:        redisComments[i/2].Id,
				UserId:    c.UserId,
				Content:   c.Content,
				CreatedAt: redisComments[i/2].CreatedAt,
			}
			comments = append(comments, com)
		}
	}

	// 缓存没查到 查库
	if len(notInCacheCommentIds) > 0 {
		cs, err := mysql.MGetCommentListByCommentId(ctx, notInCacheCommentIds)
		if err != nil {
			klog.CtxErrorf(ctx, "mysql get comment list failed %v", err)
			return nil, err
		}

		// 把评论加入缓存
		err = MSetComment(ctx, cs)
		if err != nil {
			klog.CtxErrorf(ctx, "redis set comment failed %v", err)
			return nil, err
		}

		// 把评论加入结果
		comments = append(comments, cs...)
	}

	return
}

func GetCommentIdListByVideoId(ctx context.Context, videoId int64) ([]*mysql.Comment, error) {
	// 判断commentList是否存在 不存在则读库创建列表
	err := updateCommentList(ctx, videoId)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update comment list failed %v", err)
		return nil, err
	}

	commentListKey := fmt.Sprintf(constants.RedisCommentListKey, videoId)

	// 查询评论id列表
	res, err := RDB.ZRevRangeByScoreWithScores(ctx, commentListKey, &redis.ZRangeBy{
		Min: "1", // 这是时间戳 0的时候可能会碰到避免缓存穿透放的空评论 因此置1
		Max: fmt.Sprintf("%d", time.Now().UnixMilli()),
	}).Result()
	if err != nil {
		return nil, err
	}

	// 把评论id和时间加入结果
	commentList := make([]*mysql.Comment, 0, len(res))
	for _, c := range res {
		cid, err := strconv.ParseInt(c.Member.(string), 10, 64)
		if err != nil {
			continue
		}
		com := &mysql.Comment{
			Id:        cid,
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

	// 不存在评论列表 查库
	if res == 0 {
		commentList, err := mysql.GetCommentListByVideoId(ctx, videoId)
		if err != nil {
			klog.CtxErrorf(ctx, "mysql get comment list failed %v", err)
			return err
		}

		// 数据库也没评论 避免缓存穿透 放入空数据
		if len(commentList) == 0 {
			err = RDB.ZAdd(ctx, commentListKey, redis.Z{
				Score:  0,
				Member: 0,
			}).Err()
			if err != nil {
				klog.CtxErrorf(ctx, "redis add comment id to list failed %v", err)
				return err
			}

			// 设置list的过期时间
			err = RDB.Expire(ctx, commentListKey, constants.CommentListExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second).Err()
			if err != nil {
				klog.CtxErrorf(ctx, "redis set comment list expiry failed %v", err)
				return err
			}
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
		if err != nil {
			klog.CtxErrorf(ctx, "redis add comment id to list failed %v", err)
			return err
		}

		// 设置list的过期时间
		err = RDB.Expire(ctx, commentListKey, constants.CommentListExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second).Err()
		if err != nil {
			klog.CtxErrorf(ctx, "redis set comment list expiry failed %v", err)
			return err
		}

		return nil
	}
	return nil
}
