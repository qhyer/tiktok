package handler

import (
	"context"
	"net/http"

	"tiktok/cmd/rpc"
	"tiktok/kitex_gen/relation"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type RelationActionParam struct {
	ToUserId   int64 `query:"to_user_id" vd:"$>0"`
	ActionType int32 `query:"action_type" vd:"$==1||$==2"`
}

type RelationListParam struct {
	ToUserId int64 `query:"user_id" vd:"$>0"`
}

type RelationListResponse struct {
	StatusCode int32        `json:"status_code"`
	StatusMsg  string       `json:"status_msg"`
	UserList   []*user.User `json:"user_list"`
}

type FriendListResponse struct {
	StatusCode int32                  `json:"status_code"`
	StatusMsg  string                 `json:"status_msg"`
	UserList   []*relation.FriendUser `json:"user_list"`
}

// RelationAction 关注、取关操作
func RelationAction(ctx context.Context, c *app.RequestContext) {
	var req RelationActionParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, errno.ParamErr)
		return
	}

	userId := c.GetInt64("UserID")

	// 两个用户不能相同
	if req.ToUserId == userId {
		hlog.CtxWarnf(ctx, "param error userId == toUserId")
		SendResponse(c, errno.ParamErr)
		return
	}

	// rpc通信
	relationActionResponse, err := rpc.RelationAction(context.Background(), &relation.DouyinRelationActionRequest{
		UserId:     userId,
		ToUserId:   req.ToUserId,
		ActionType: req.ActionType,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err)
		return
	}

	c.JSON(consts.StatusOK, Response{
		StatusCode: relationActionResponse.GetStatusCode(),
		StatusMsg:  relationActionResponse.GetStatusMsg(),
	})
}

// FollowList 关注列表
func FollowList(ctx context.Context, c *app.RequestContext) {
	var req RelationListParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, errno.ParamErr)
		return
	}

	userId := c.GetInt64("UserID")

	// rpc通信
	relationResponse, err := rpc.FollowList(context.Background(), &relation.DouyinRelationFollowListRequest{
		UserId:   userId,
		ToUserId: req.ToUserId,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, RelationListResponse{
		StatusCode: relationResponse.GetStatusCode(),
		StatusMsg:  relationResponse.GetStatusMsg(),
		UserList:   relationResponse.GetUserList(),
	})
}

// FollowerList 粉丝列表
func FollowerList(ctx context.Context, c *app.RequestContext) {
	var req RelationListParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, errno.ParamErr)
		return
	}

	userId := c.GetInt64("UserID")

	// rpc通信
	relationResponse, err := rpc.FollowerList(context.Background(), &relation.DouyinRelationFollowerListRequest{
		UserId:   userId,
		ToUserId: req.ToUserId,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, RelationListResponse{
		StatusCode: relationResponse.GetStatusCode(),
		StatusMsg:  relationResponse.GetStatusMsg(),
		UserList:   relationResponse.GetUserList(),
	})
}

// FriendList 好友列表
func FriendList(ctx context.Context, c *app.RequestContext) {
	// 用户可以看到别人的关注和粉丝列表
	// 理论上就可以知道别人的好友列表
	// 但看了文档觉得这个接口是为消息功能设计的
	// 因此目前只支持查询当前登录用户的好友

	userId := c.GetInt64("UserID")

	// rpc通信
	relationResponse, err := rpc.FriendList(context.Background(), &relation.DouyinRelationFriendListRequest{
		UserId: userId,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, FriendListResponse{
		StatusCode: relationResponse.GetStatusCode(),
		StatusMsg:  relationResponse.GetStatusMsg(),
		UserList:   relationResponse.GetUserList(),
	})
}
