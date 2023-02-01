package service

import (
	"context"

	"tiktok/dal/mysql"
	"tiktok/dal/pack"
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

func (s *MessageListService) MessageList(req *message.DouyinMessageListRequest) ([]*message.Message, error) {
	userId := req.GetUserId()
	toUserId := req.GetToUserId()
	limit := int(req.GetLimit())

	// 获取聊天消息记录
	msgs, err := mysql.MessageList(s.ctx, userId, toUserId, limit)
	if err != nil {
		klog.CtxErrorf(s.ctx, "db get message list failed %v", err)
		return nil, err
	}

	messages := pack.Messages(msgs)

	return messages, nil
}
