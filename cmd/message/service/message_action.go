package service

import (
	"context"

	"tiktok/cmd/rpc"
	"tiktok/dal/mysql"
	"tiktok/kitex_gen/message"
	"tiktok/kitex_gen/relation"
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

	// 获取两个用户是不是朋友关系
	relationResp, err := rpc.IsFriendRelation(s.ctx, &relation.DouyinRelationIsFriendRequest{
		UserId:   userId,
		ToUserId: toUserId,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "rpc get friend relation failed %v", err)
		return err
	}
	isFriend := relationResp.GetIsFriend()

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
