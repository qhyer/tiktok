package service

import (
	"context"

	"tiktok/dal/mysql"
	"tiktok/dal/pack"
	"tiktok/dal/redis"
	"tiktok/kitex_gen/message"

	"github.com/cloudwego/kitex/pkg/klog"
)

type MessageListService struct {
	ctx context.Context
}

// NewMessageListService new MessageListService
func NewMessageListService(ctx context.Context) *MessageListService {
	return &MessageListService{ctx: ctx}
}

// MessageList get user unread message list
func (s *MessageListService) MessageList(req *message.DouyinMessageListRequest) ([]*message.Message, error) {
	userId := req.GetUserId()
	toUserId := req.GetToUserId()
	preMsgTime := req.GetPreMsgTime()

	// 获取消息id
	msgIds, err := redis.GetMessageIdsByUserIdAndPreMsgTime(s.ctx, userId, toUserId, preMsgTime)
	if err != nil {
		klog.CtxErrorf(s.ctx, "reids get message id list failed %v", err)
		return nil, err
	}

	// 获取消息记录
	rmsgs, err := redis.MGetMessageByMessageId(s.ctx, msgIds)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis get message list failed %v", err)
		return nil, err
	}

	// 合并消息
	msgMap := make(map[int64]*mysql.Message, 0)
	msgs := make([]*mysql.Message, 0)
	for _, m := range rmsgs {
		if m == nil {
			continue
		}
		msgMap[m.Id] = m
	}

	for _, m := range msgIds {
		if m == nil {
			continue
		}
		res := msgMap[m.Id]
		if res == nil {
			continue
		}
		msgs = append(msgs, res)
	}

	messages := pack.Messages(msgs)

	return messages, nil
}
