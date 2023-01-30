package handler

import (
	"context"

	"tiktok/kitex_gen/relation"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/errno"
	"tiktok/pkg/rpc"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type RelationActionParam struct {
	ToUserId   int64 `query:"to_user_id" vd:"$>0"`
	ActionType int32 `query:"action_type" vd:"$==1||$==2"`
}

// RelationAction 关注、取关操作
func RelationAction(ctx context.Context, c *app.RequestContext) {
	var req RelationActionParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, err, nil)
		return
	}

	userId := c.GetInt64("UserID")

	// rpc通信
	_, err = rpc.RelationAction(ctx, &relation.DouyinRelationActionRequest{
		UserId:     userId,
		ToUserId:   req.ToUserId,
		ActionType: req.ActionType,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err, nil)
		return
	}

	SendResponse(c, errno.Success, nil)
}

type RelationListParam struct {
	ToUserId int64 `query:"user_id" vd:"$>0"`
}

type RelationListResponse struct {
	UserList []*user.User `json:"user_list"`
}

// FollowList 关注列表
func FollowList(ctx context.Context, c *app.RequestContext) {
	var req RelationListParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, err, nil)
		return
	}

	userId := c.GetInt64("UserID")

	// rpc通信
	relationResponse, err := rpc.FollowList(ctx, &relation.DouyinRelationFollowListRequest{
		UserId:   userId,
		ToUserId: req.ToUserId,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err, nil)
		return
	}

	SendResponse(c, errno.Success, RelationListResponse{UserList: relationResponse.UserList})
}

// FollowerList 粉丝列表
func FollowerList(ctx context.Context, c *app.RequestContext) {
	var req RelationListParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, err, nil)
		return
	}

	userId := c.GetInt64("UserID")

	// rpc通信
	relationResponse, err := rpc.FollowerList(ctx, &relation.DouyinRelationFollowerListRequest{
		UserId:   userId,
		ToUserId: req.ToUserId,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err, nil)
		return
	}

	SendResponse(c, errno.Success, RelationListResponse{UserList: relationResponse.UserList})
}

// FriendList 好友列表
func FriendList(ctx context.Context, c *app.RequestContext) {
	// 用户可以看到别人的关注和粉丝列表
	// 理论上就可以知道别人的好友列表
	// 但看了文档觉得这个接口是为消息功能设计的
	// 因此目前只支持查询当前登录用户的好友

	userId := c.GetInt64("UserID")

	// rpc通信
	relationResponse, err := rpc.FriendList(ctx, &relation.DouyinRelationFriendListRequest{
		UserId: userId,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err, nil)
		return
	}

	SendResponse(c, errno.Success, RelationListResponse{UserList: relationResponse.UserList})
}
