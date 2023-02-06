package service

import (
	"context"

	"tiktok/dal/neo4j"
	"tiktok/dal/redis"
	"tiktok/kitex_gen/message"
	"tiktok/kitex_gen/relation"
	"tiktok/pkg/constants"

	"github.com/cloudwego/kitex/pkg/klog"
)

type FriendListService struct {
	ctx context.Context
}

// NewFriendListService new FriendListService
func NewFriendListService(ctx context.Context) *FriendListService {
	return &FriendListService{ctx: ctx}
}

func (s *FriendListService) FriendList(req *relation.DouyinRelationFriendListRequest) ([]*relation.FriendUser, error) {
	userId := req.GetUserId()

	// 获取当前用户的朋友
	friends, err := redis.GetFriendListByUserId(s.ctx, userId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis get friend list failed %v", err)
		return nil, err
	}

	if len(friends) == 0 {
		return friends, nil
	}

	userIds := make([]int64, 0, len(friends))
	for _, f := range friends {
		if f == nil {
			continue
		}
		userIds = append(userIds, f.Id)
	}

	// 获取最新消息
	chats, err := neo4j.MQueryLastMessage(s.ctx, userId, userIds)
	if err != nil {
		klog.CtxErrorf(s.ctx, "neo4j query last message error %v", err)
		return nil, err
	}
	if chats != nil {
		chatMap := make(map[int64]*message.Message, 0)
		for _, c := range chats {
			if c == nil {
				continue
			}
			// 是发送者
			if c.FromUserId == userId {
				chatMap[c.ToUserId] = c
			} else {
				chatMap[c.FromUserId] = c
			}
		}

		// 添加最新消息
		for i, u := range friends {
			if u == nil {
				continue
			}
			chat := chatMap[u.Id]
			if chat == nil {
				continue
			}
			friends[i].Message = &chat.Content
			if chat.FromUserId == userId {
				friends[i].MsgType = constants.Sender
			} else {
				friends[i].MsgType = constants.Receiver
			}
		}
	}
	return friends, nil
}
