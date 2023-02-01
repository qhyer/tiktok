package main

import (
	"context"

	"tiktok/cmd/relation/service"
	"tiktok/dal/pack"
	"tiktok/kitex_gen/relation"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
)

// RelationSrvImpl implements the last service interface defined in the IDL.
type RelationSrvImpl struct{}

// RelationAction implements the RelationSrvImpl interface.
func (s *RelationSrvImpl) RelationAction(ctx context.Context, req *relation.DouyinRelationActionRequest) (resp *relation.DouyinRelationActionResponse, err error) {
	resp = new(relation.DouyinRelationActionResponse)
	switch req.GetActionType() {
	case constants.FollowAction:
		err = service.NewRelationActionService(ctx).Follow(req)
	case constants.UnfollowAction:
		err = service.NewRelationActionService(ctx).Unfollow(req)
	default:
		err = errno.ParamErr
	}
	if err != nil {
		resp = pack.BuildRelationActionResp(err)
		return resp, err
	}
	resp = pack.BuildRelationActionResp(errno.Success)
	return resp, nil
}

// RelationFollowList implements the RelationSrvImpl interface.
func (s *RelationSrvImpl) RelationFollowList(ctx context.Context, req *relation.DouyinRelationFollowListRequest) (resp *relation.DouyinRelationFollowListResponse, err error) {
	resp = new(relation.DouyinRelationFollowListResponse)

	users, err := service.NewFollowListService(ctx).FollowList(req)
	if err != nil {
		resp = pack.BuildFollowListResp(err)
		return resp, err
	}

	resp = pack.BuildFollowListResp(errno.Success)
	resp.UserList = users
	return resp, nil
}

// RelationFollowerList implements the RelationSrvImpl interface.
func (s *RelationSrvImpl) RelationFollowerList(ctx context.Context, req *relation.DouyinRelationFollowerListRequest) (resp *relation.DouyinRelationFollowerListResponse, err error) {
	resp = new(relation.DouyinRelationFollowerListResponse)

	users, err := service.NewFollowerListService(ctx).FollowerList(req)
	if err != nil {
		resp = pack.BuildFollowerListResp(err)
		return resp, err
	}

	resp = pack.BuildFollowerListResp(errno.Success)
	resp.UserList = users
	return resp, nil
}

// RelationFriendList implements the RelationSrvImpl interface.
func (s *RelationSrvImpl) RelationFriendList(ctx context.Context, req *relation.DouyinRelationFriendListRequest) (resp *relation.DouyinRelationFriendListResponse, err error) {
	resp = new(relation.DouyinRelationFriendListResponse)

	users, err := service.NewFriendListService(ctx).FriendList(req)
	if err != nil {
		resp = pack.BuildFriendListResp(err)
		return resp, err
	}

	resp = pack.BuildFriendListResp(errno.Success)
	resp.UserList = users
	return resp, nil
}

// RelationIsFriend implements the RelationSrvImpl interface.
func (s *RelationSrvImpl) RelationIsFriend(ctx context.Context, req *relation.DouyinRelationIsFriendRequest) (resp *relation.DouyinRelationIsFriendResponse, err error) {
	resp = new(relation.DouyinRelationIsFriendResponse)

	isFriend, err := service.NewQueryRelationService(ctx).IsFriend(req)
	if err != nil {
		resp.IsFriend = false
		return resp, err
	}

	resp.IsFriend = isFriend
	return resp, err
}
