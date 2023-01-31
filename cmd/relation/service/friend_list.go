package service

import (
	"context"

	"tiktok/dal/neo4j"
	"tiktok/kitex_gen/relation"
	"tiktok/kitex_gen/user"

	"github.com/cloudwego/kitex/pkg/klog"
)

type FriendListService struct {
	ctx context.Context
}

// NewFriendListService new FriendListService
func NewFriendListService(ctx context.Context) *FriendListService {
	return &FriendListService{ctx: ctx}
}

func (s *FriendListService) FriendList(req *relation.DouyinRelationFriendListRequest) ([]*user.User, error) {
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
		userFollowMap[u.Id] = true
	}

	// 交集就是朋友
	friends := make([]*user.User, 0)
	for _, u := range followerList {
		// 用户没有关注他
		if !userFollowMap[u.Id] {
			continue
		}
		friends = append(friends, &user.User{
			Id:            u.Id,
			Name:          u.Name,
			FollowCount:   u.FollowCount,
			FollowerCount: u.FollowerCount,
			IsFollow:      true,
		})
	}
	return friends, err

}
