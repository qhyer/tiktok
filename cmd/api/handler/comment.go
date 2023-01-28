package handler

import (
	"context"
	"net/http"

	"tiktok/kitex_gen/comment"
	"tiktok/pkg/errno"
	"tiktok/pkg/rpc"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type CommentActionParam struct {
	VideoId     int64  `query:"video_id" vd:"$!=nil&&$>0"`
	ActionType  int32  `query:"action_type" vd:"$!=nil&&$==1||$==2"`
	CommentText string `query:"comment_text"`
	CommentID   int64  `query:"comment_id"`
}

type CommentActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

const (
	DoCommentAction     = 1
	DeleteCommentAction = 2
)

func CommentAction(ctx context.Context, c *app.RequestContext) {
	var req CommentActionParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, err)
		return
	}
	switch req.ActionType {
	case DoCommentAction:
		if req.CommentText == "" {
			hlog.CtxWarnf(ctx, "param comment_text error %v", req.CommentText)
			SendResponse(c, errno.ParamErr)
			return
		}
	case DeleteCommentAction:
		if req.CommentID <= 0 {
			hlog.CtxWarnf(ctx, "param comment_id error %v", req.CommentID)
			SendResponse(c, errno.ParamErr)
			return
		}
	}

	userId := c.GetInt64("UserID")

	// rpc通信
	CommentResponse, err := rpc.CommentAction(ctx, &comment.DouyinCommentActionRequest{
		UserId:      userId,
		VideoId:     req.VideoId,
		ActionType:  req.ActionType,
		CommentText: &req.CommentText,
		CommentId:   &req.CommentID,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, CommentActionResponse{
		StatusCode: CommentResponse.StatusCode,
		StatusMsg:  *CommentResponse.StatusMsg,
	})
}

type CommentListParam struct {
	VideoId int64 `query:"video_id" vd:"$!=nil&&$>0"`
}

type CommentListResponse struct {
	StatusCode  int32              `json:"status_code"`
	StatusMsg   string             `json:"status_msg"`
	CommentList []*comment.Comment `json:"comment_list"`
}

func CommentList(ctx context.Context, c *app.RequestContext) {
	var req CommentListParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, err)
		return
	}
	userId := c.GetInt64("UserID")

	// rpc通信
	CommentResponse, err := rpc.CommentList(ctx, &comment.DouyinCommentListRequest{
		UserId:  userId,
		VideoId: req.VideoId,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		StatusCode:  CommentResponse.StatusCode,
		StatusMsg:   *CommentResponse.StatusMsg,
		CommentList: CommentResponse.CommentList,
	})
}
