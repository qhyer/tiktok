package pack

import (
	"errors"

	"tiktok/kitex_gen/relation"
	"tiktok/pkg/errno"
)

func BuildRelationActionResp(err error) *relation.DouyinRelationActionResponse {
	if err == nil {
		return relationActionResp(errno.Success)
	}
	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return relationActionResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return relationActionResp(s)
}

func relationActionResp(err errno.ErrNo) *relation.DouyinRelationActionResponse {
	return &relation.DouyinRelationActionResponse{StatusCode: err.ErrCode, StatusMsg: &err.ErrMsg}
}

func BuildFollowListResp(err error) *relation.DouyinRelationFollowListResponse {
	if err == nil {
		return followListResp(errno.Success)
	}
	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return followListResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return followListResp(s)
}

func followListResp(err errno.ErrNo) *relation.DouyinRelationFollowListResponse {
	return &relation.DouyinRelationFollowListResponse{StatusCode: err.ErrCode, StatusMsg: &err.ErrMsg}
}

func BuildFollowerListResp(err error) *relation.DouyinRelationFollowerListResponse {
	if err == nil {
		return followerListResp(errno.Success)
	}
	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return followerListResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return followerListResp(s)
}

func followerListResp(err errno.ErrNo) *relation.DouyinRelationFollowerListResponse {
	return &relation.DouyinRelationFollowerListResponse{StatusCode: err.ErrCode, StatusMsg: &err.ErrMsg}
}

func BuildFriendListResp(err error) *relation.DouyinRelationFriendListResponse {
	if err == nil {
		return friendListResp(errno.Success)
	}
	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return friendListResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return friendListResp(s)
}

func friendListResp(err errno.ErrNo) *relation.DouyinRelationFriendListResponse {
	return &relation.DouyinRelationFriendListResponse{StatusCode: err.ErrCode, StatusMsg: &err.ErrMsg}
}
