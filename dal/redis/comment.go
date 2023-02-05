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

func AddNewCommentToCommentList(ctx context.Context, comment *mysql.Comment, videoId int64) error {
	// 判断评论列表是否存在，不存在则创建列表
	err := updateCommentList(ctx, videoId)
	if err != nil {
		return err
	}

	if comment == nil {
		return nil
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
						return true
					end
					return false
					`)
		keys := []string{videoKey}
		if ret, err := incrBy.Run(ctx, RDB, keys).Result(); err != nil || ret != true {
			return err
		}
		return nil
	})
	return err
}

func DeleteCommentFromCommentList(ctx context.Context, comment *mysql.Comment, videoId int64) error {
	// 判断评论列表是否存在，不存在则创建列表
	err := updateCommentList(ctx, videoId)
	if err != nil {
		return err
	}

	if comment == nil {
		return nil
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
						return true
					end
					return false
					`)
		keys := []string{videoKey}
		if ret, err := incrBy.Run(ctx, RDB, keys).Result(); err != nil || ret != true {
			return err
		}
		return nil
	})
	return err
}

func MAddCommentIdToCommentList(ctx context.Context, comments []*mysql.Comment, videoId int64) error {
	// 判断评论列表是否存在，不存在则创建列表
	err := updateCommentList(ctx, videoId)
	if err != nil {
		return err
	}

	if len(comments) == 0 {
		return nil
	}

	commentListKey := fmt.Sprintf(constants.RedisCommentListKey, videoId)
	commentIds := make([]redis.Z, 0, len(comments))
	for _, c := range comments {
		if c == nil {
			continue
		}
		commentIds = append(commentIds, redis.Z{
			Score:  float64(c.CreatedAt.UnixMilli()),
			Member: c.Id,
		})
	}
	err = RDB.ZAdd(ctx, commentListKey, commentIds...).Err()

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

func MGetCommentByCommentId(ctx context.Context, redisComments []*mysql.Comment) (comments []*mysql.Comment, notInCacheCommentIds []int64) {
	res, err := RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for _, c := range redisComments {
			commentKey := fmt.Sprintf(constants.RedisCommentKey, c.Id)
			pipeliner.HGetAll(ctx, commentKey)
			pipeliner.Expire(ctx, commentKey, constants.CommentExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second)
		}
		return nil
	})
	if err != nil {
		// 出错 把所有id都加入未找到
		for _, c := range redisComments {
			notInCacheCommentIds = append(notInCacheCommentIds, c.Id)
		}
		return nil, notInCacheCommentIds
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
		Min: "0",
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
		if err != nil {
			klog.CtxErrorf(ctx, "redis add comment id to list failed %v", err)
			return err
		}

		// 设置list的过期时间
		err = RDB.Expire(ctx, commentListKey, constants.CommentListExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second).Err()
		if err != nil {
			klog.CtxErrorf(ctx, "redis set comment list expire failed %v", err)
			return err
		}
		return nil
	}
	return nil
}
