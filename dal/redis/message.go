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

func GetMessageIdsByUserId(ctx context.Context, userId int64, toUserId int64) ([]*mysql.Message, error) {
	msgListKey := fmt.Sprintf(constants.RedisMessageListKey, userId, toUserId)
	res := make([]*mysql.Message, 0)
	// 判断消息列表是否存在，不存在则创建列表
	err := updateMessageList(ctx, userId, toUserId)
	if err != nil {
		klog.CtxErrorf(ctx, "redis update message list failed %v", err)
		return res, err
	}

	// 查询id列表
	msgs, err := RDB.ZRevRangeByScoreWithScores(ctx, msgListKey, &redis.ZRangeBy{
		Min: "1", // 这里是时间戳 避免缓存穿透 因此置1
		Max: fmt.Sprintf("%d", time.Now().UnixMilli()),
	}).Result()
	if err != nil {
		klog.CtxErrorf(ctx, "redis get message id list failed %v", err)
		return res, err
	}

	// 把消息id加入到列表
	for _, m := range msgs {
		mid, err := strconv.ParseInt(m.Member.(string), 10, 64)
		if err != nil {
			continue
		}
		res = append(res, &mysql.Message{
			Id:        mid,
			CreatedAt: time.UnixMilli(int64(m.Score)),
		})
	}

	// 读完后要清空redis中的已读消息
	err = RDB.ZRem(ctx, msgListKey, msgs).Err()
	if err != nil {
		klog.CtxErrorf(ctx, "redis remove read messages failed %v", err)
		return res, nil
	}

	return res, nil
}

func MGetMessageByMessageId(ctx context.Context, redisMsgs []*mysql.Message) (messages []*mysql.Message, err error) {
	notInCacheMessageIds := make([]int64, 0)
	res, err := RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for _, m := range redisMsgs {
			messageKey := fmt.Sprintf(constants.RedisMessageKey, m.Id)
			pipeliner.HGetAll(ctx, messageKey)
			pipeliner.Expire(ctx, messageKey, constants.MessageExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second)
		}
		return nil
	})
	if err != nil {
		klog.CtxErrorf(ctx, "redis get message failed %v", err)
		return nil, err
	}

	// 处理getall结果
	for i, l := range res {
		if i%2 == 0 {
			var c mysql.Message
			err := MustScan(l.(*redis.MapStringStringCmd), &c)
			if err != nil {
				notInCacheMessageIds = append(notInCacheMessageIds, redisMsgs[i/2].Id)
				continue
			}
			msg := &mysql.Message{
				Id:        redisMsgs[i/2].Id,
				UserId:    c.UserId,
				ToUserId:  c.ToUserId,
				Content:   c.Content,
				CreatedAt: redisMsgs[i/2].CreatedAt,
			}
			messages = append(messages, msg)
		}
	}

	// 缓存没查到 查库
	if len(notInCacheMessageIds) > 0 {
		msgs, err := mysql.MGetMessageListByMessageId(ctx, notInCacheMessageIds)
		if err != nil {
			klog.CtxErrorf(ctx, "mysql get message failed %v", err)
			return nil, err
		}

		// 把消息加入缓存
		err = MSetMessage(ctx, msgs)
		if err != nil {
			klog.CtxErrorf(ctx, "redis set message failed %v", err)
			return nil, err
		}

		// 把消息加入结果
		messages = append(messages, msgs...)
	}
	return
}

func MSetMessage(ctx context.Context, messages []*mysql.Message) error {
	_, err := RDB.TxPipelined(ctx, func(pipeliner redis.Pipeliner) error {
		for _, c := range messages {
			if c == nil {
				continue
			}
			messageKey := fmt.Sprintf(constants.RedisMessageKey, c.Id)
			pipeliner.HMSet(ctx, messageKey,
				"user_id", c.UserId,
				"to_user_id", c.ToUserId,
				"content", c.Content,
			)
			pipeliner.Expire(ctx, messageKey, constants.MessageExpiry+time.Duration(rand.Intn(constants.MaxRandExpireSecond))*time.Second)
		}
		return nil
	})
	return err
}

func MAddMessageToMessageList(ctx context.Context, messages []*mysql.Message) error {
	for _, message := range messages {
		msgId := message.Id
		userId := message.UserId
		toUserId := message.ToUserId

		// 如果key存在 则加入聊天记录中
		myMsgListKey := fmt.Sprintf(constants.RedisMessageListKey, userId, toUserId)
		hisMsgListKey := fmt.Sprintf(constants.RedisMessageListKey, toUserId, userId)
		lua := redis.NewScript(`
					if redis.call("Exists", KEYS[1]) > 0 then
						redis.call("ZAdd", KEYS[1], "id", ARGV[1], "user_id", ARGV[2], "to_user_id", ARGV[3], "content", ARGV[4], "created_at", ARGV[5])
					end
					if redis.call("Exists", KEYS[2]) > 0 then
						redis.call("ZAdd", KEYS[1], "id", ARGV[1], "user_id", ARGV[3], "to_user_id", ARGV[2], "content", ARGV[4], "created_at", ARGV[5])
					end
					return true
					`)
		msgContent := message.Content
		curTime := message.CreatedAt.UnixMilli()
		keys := []string{myMsgListKey, hisMsgListKey}
		args := []interface{}{msgId, userId, toUserId, msgContent, curTime}
		if err := lua.Run(ctx, RDB, keys, args).Err(); err != nil {
			return err
		}
	}
	return nil
}

func updateMessageList(ctx context.Context, userId int64, toUserId int64) error {
	msgListKey := fmt.Sprintf(constants.RedisMessageListKey, userId, toUserId)
	exists, err := RDB.Exists(ctx, msgListKey).Result()
	if err != nil {
		return err
	}

	// 消息列表不存在 需要从数据库把历史消息全部取出放入缓存
	if exists == 0 {
		msgs, err := mysql.GetMessageListByUserId(ctx, userId, toUserId)
		if err != nil {
			klog.CtxErrorf(ctx, "mysql get message list failed %v", err)
			return err
		}

		// 从列表中获取消息id
		msgzs := make([]redis.Z, 0, len(msgs))
		for _, m := range msgs {
			if m == nil {
				continue
			}
			msgzs = append(msgzs, redis.Z{
				Score:  float64(m.CreatedAt.UnixMilli()),
				Member: m.Id,
			})
		}

		// 防止缓存穿透加入空消息id
		msgzs = append(msgzs, redis.Z{
			Score:  0,
			Member: 0,
		})

		// 把消息加入缓存
		err = MSetMessage(ctx, msgs)
		if err != nil {
			klog.CtxErrorf(ctx, "redis set message failed %v", err)
			return err
		}

		// 把消息id加入缓存
		err = RDB.ZAdd(ctx, msgListKey, msgzs...).Err()
		if err != nil {
			klog.CtxErrorf(ctx, "redis zadd messages failed %v", err)
			return err
		}

		// 设置list的过期时间
		err = RDB.Expire(ctx, msgListKey, constants.MessageListExpiry).Err()
		if err != nil {
			klog.CtxErrorf(ctx, "redis set comment list expiry failed %v", err)
			return err
		}
	}
	return nil
}
