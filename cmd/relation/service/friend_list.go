package service

import (
	"context"

	"tiktok/dal/neo4j"
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

	// 获取当前用户的粉丝
	followerList, err := neo4j.FollowerList(s.ctx, userId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "neo4j get follower list failed %v", err)
		return nil, err
	}

	// 获取当前用户的关注
	followList, err := neo4j.FollowList(s.ctx, userId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "neo4j get follow list failed %v", err)
		return nil, err
	}
	userFollowMap := make(map[int64]bool, 0)
	for _, u := range followList {
		if u == nil {
			continue
		}
		userFollowMap[u.Id] = true
	}

	friends := make([]*relation.FriendUser, 0)
	// 交集就是朋友
	for _, u := range followerList {
		if u == nil {
			continue
		}
		// 用户没有关注他
		if !userFollowMap[u.Id] {
			continue
		}
		friends = append(friends, &relation.FriendUser{
			Id:            u.Id,
			Name:          u.Name,
			FollowCount:   u.FollowCount,
			FollowerCount: u.FollowerCount,
			IsFollow:      true,
			Avatar:        constants.DefaultAvatarUrl, // 没有找到上传头像的地方 先返回一个固定头像
		})
	}

	// 查询和朋友的最新消息

	return friends, nil
}
