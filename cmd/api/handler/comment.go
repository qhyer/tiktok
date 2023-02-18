package handler

import (
	"context"
	"net/http"

	"tiktok/cmd/rpc"
	"tiktok/kitex_gen/comment"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type CommentActionParam struct {
	VideoId     int64  `query:"video_id" vd:"$>0"`
	ActionType  int32  `query:"action_type" vd:"$==1||$==2"`
	CommentText string `query:"comment_text"`
	CommentID   int64  `query:"comment_id"`
}

type CommentActionResponse struct {
	StatusCode int32            `json:"status_code"`
	StatusMsg  string           `json:"status_msg"`
	Comment    *comment.Comment `json:"comment"`
}

type CommentListParam struct {
	VideoId int64 `query:"video_id" vd:"$>0"`
}

type CommentListResponse struct {
	StatusCode  int32              `json:"status_code"`
	StatusMsg   string             `json:"status_msg"`
	CommentList []*comment.Comment `json:"comment_list"`
}

// CommentAction 评论和删除评论
func CommentAction(ctx context.Context, c *app.RequestContext) {
	var req CommentActionParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, errno.ParamErr)
		return
	}
	switch req.ActionType {
	case constants.CreateCommentAction:
		if len(req.CommentText) == 0 {
			hlog.CtxWarnf(ctx, "param comment_text error %v", req.CommentText)
			SendResponse(c, errno.ParamErr)
			return
		}
	case constants.DeleteCommentAction:
		if req.CommentID <= 0 {
			hlog.CtxWarnf(ctx, "param comment_id error %v", req.CommentID)
			SendResponse(c, errno.ParamErr)
			return
		}
	}

	userId := c.GetInt64("UserID")

	// rpc通信
	commentActionResponse, err := rpc.CommentAction(context.Background(), &comment.DouyinCommentActionRequest{
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

	c.JSON(consts.StatusOK, Response{
		StatusCode: commentActionResponse.GetStatusCode(),
		StatusMsg:  commentActionResponse.GetStatusMsg(),
	})
}

// CommentList 获取视频的评论列表
func CommentList(ctx context.Context, c *app.RequestContext) {
	var req CommentListParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, errno.ParamErr)
		return
	}
	userId := c.GetInt64("UserID")

	// rpc通信
	commentResponse, err := rpc.CommentList(context.Background(), &comment.DouyinCommentListRequest{
		UserId:  userId,
		VideoId: req.VideoId,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		StatusCode:  errno.Success.ErrCode,
		StatusMsg:   errno.Success.ErrMsg,
		CommentList: commentResponse.GetCommentList(),
	})
}
