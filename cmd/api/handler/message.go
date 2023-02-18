package handler

import (
	"context"
	"net/http"

	"tiktok/cmd/rpc"
	"tiktok/kitex_gen/message"
	"tiktok/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type MessageActionParam struct {
	ToUserId   int64  `query:"to_user_id" vd:"$>0"`
	ActionType int32  `query:"action_type" vd:"$==1"`
	Content    string `query:"content" vd:"len($)>0"`
}

type MessageListParam struct {
	ToUserId       int64 `query:"to_user_id" vd:"$>0"`
	PreMessageTime int64 `query:"pre_msg_time" vd:"$>=0"`
}

type MessageListResponse struct {
	StatusCode  int32              `json:"status_code"`
	StatusMsg   string             `json:"status_msg"`
	MessageList []*message.Message `json:"message_list"`
}

// MessageAction 发送消息
func MessageAction(ctx context.Context, c *app.RequestContext) {
	var req MessageActionParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, errno.ParamErr)
		return
	}

	userId := c.GetInt64("UserID")
	toUserId := req.ToUserId
	actionType := req.ActionType
	content := req.Content

	// 两个用户不能相同
	if userId == toUserId {
		hlog.CtxWarnf(ctx, "param error userId == toUserId")
		SendResponse(c, errno.ParamErr)
		return
	}

	messageActionResp, err := rpc.MessageAction(context.Background(), &message.DouyinMessageActionRequest{
		UserId:     userId,
		ToUserId:   toUserId,
		ActionType: actionType,
		Content:    content,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc send message action error %v", err)
		SendResponse(c, err)
		return
	}

	c.JSON(consts.StatusOK, Response{
		StatusCode: messageActionResp.GetStatusCode(),
		StatusMsg:  messageActionResp.GetStatusMsg(),
	})
}

// MessageList 消息列表
func MessageList(ctx context.Context, c *app.RequestContext) {
	var req MessageListParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, errno.ParamErr)
		return
	}

	userId := c.GetInt64("UserID")
	toUserId := req.ToUserId
	preMsgTime := req.PreMessageTime

	messageResponse, err := rpc.MessageList(context.Background(), &message.DouyinMessageListRequest{
		UserId:     userId,
		ToUserId:   toUserId,
		PreMsgTime: preMsgTime,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc get message list error %v", err)
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, &MessageListResponse{
		StatusCode:  errno.Success.ErrCode,
		StatusMsg:   errno.Success.ErrMsg,
		MessageList: messageResponse.GetMessageList(),
	})
}
