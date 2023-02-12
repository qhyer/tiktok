package service

import (
	"context"

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

	// 获取消息id
	msgIds, err := redis.GetMessageIdsByUserId(s.ctx, userId, toUserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "reids get message id list failed %v", err)
		return nil, err
	}

	// 获取消息记录
	msgs, err := redis.MGetMessageByMessageId(s.ctx, msgIds)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis get message list failed %v", err)
		return nil, err
	}

	messages := pack.Messages(msgs)

	return messages, nil
}
