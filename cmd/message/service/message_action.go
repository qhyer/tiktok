package service

import (
	"context"

	"tiktok/dal/mysql"
	"tiktok/kitex_gen/message"
	"tiktok/pkg/errno"

	"github.com/cloudwego/kitex/pkg/klog"
)

type MessageActionService struct {
	ctx context.Context
}

// NewMessageActionService new MessageActionService
func NewMessageActionService(ctx context.Context) *MessageActionService {
	return &MessageActionService{ctx: ctx}
}

// SendMessage user send message to friend
func (s *MessageActionService) SendMessage(req *message.DouyinMessageActionRequest) error {
	userId := req.GetUserId()
	toUserId := req.GetToUserId()
	content := req.GetContent()

	// TODO 查询好友关系
	isFriend := false

	if isFriend {
		err := mysql.CreateMessage(s.ctx, []*mysql.Message{{
			UserId:   userId,
			ToUserId: toUserId,
			Content:  content,
		}})
		if err != nil {
			return errno.DBOperationFailedErr
		}
		return nil
	} else {
		klog.CtxWarnf(s.ctx, "not friend relation %v", req)
		return errno.ParamErr
	}
}
