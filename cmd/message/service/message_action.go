package service

import (
	"context"

	"tiktok/dal/mysql"
	"tiktok/dal/neo4j"
	"tiktok/dal/redis"
	"tiktok/kitex_gen/message"
	"tiktok/pkg/censor"
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

	// 过滤敏感词
	content = censor.TextCensor.GetFilter().Replace(content, '*')

	// 关系中插入最后一条消息 这里会判断是否为好友关系
	ok, err := neo4j.UpsertLastMessage(s.ctx, userId, toUserId, content)
	if err != nil {
		klog.CtxErrorf(s.ctx, "neo4j upsert last message failed %v", err)
		return err
	}

	if ok {
		// 把消息加入数据库
		msgs, err := mysql.CreateMessage(s.ctx, []*mysql.Message{{
			UserId:   userId,
			ToUserId: toUserId,
			Content:  content,
		}})
		if err != nil {
			klog.CtxErrorf(s.ctx, "mysql create message failed %v", err)
			return errno.DatabaseOperationFailedErr
		}

		// 把消息加入缓存
		err = redis.AddNewMessageToMessageList(s.ctx, msgs[0])
		if err != nil {
			klog.CtxErrorf(s.ctx, "redis add message to list failed %v", err)
			return err
		}
	} else {
		klog.CtxWarnf(s.ctx, "not friend relation %v", req)
		return errno.ParamErr
	}
	return nil
}
