package handler

import (
	"context"

	"tiktok/kitex_gen/comment"
	"tiktok/pkg/constants"
	"tiktok/pkg/errno"
	"tiktok/pkg/rpc"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type CommentActionParam struct {
	VideoId     int64  `query:"video_id" vd:"$>0"`
	ActionType  int32  `query:"action_type" vd:"$==1||$==2"`
	CommentText string `query:"comment_text"`
	CommentID   int64  `query:"comment_id"`
}

// CommentAction 评论和删除评论
func CommentAction(ctx context.Context, c *app.RequestContext) {
	var req CommentActionParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, err, nil)
		return
	}
	switch req.ActionType {
	case constants.DoCommentAction:
		if len(req.CommentText) == 0 {
			hlog.CtxWarnf(ctx, "param comment_text error %v", req.CommentText)
			SendResponse(c, errno.ParamErr, nil)
			return
		}
	case constants.DeleteCommentAction:
		if req.CommentID <= 0 {
			hlog.CtxWarnf(ctx, "param comment_id error %v", req.CommentID)
			SendResponse(c, errno.ParamErr, nil)
			return
		}
	}

	userId := c.GetInt64("UserID")

	// rpc通信
	_, err = rpc.CommentAction(ctx, &comment.DouyinCommentActionRequest{
		UserId:      userId,
		VideoId:     req.VideoId,
		ActionType:  req.ActionType,
		CommentText: &req.CommentText,
		CommentId:   &req.CommentID,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err, nil)
		return
	}

	SendResponse(c, errno.Success, nil)
}

type CommentListParam struct {
	VideoId int64 `query:"video_id" vd:"$>0"`
}

type CommentListResponse struct {
	CommentList []*comment.Comment `json:"comment_list"`
}

// CommentList 获取视频的评论列表
func CommentList(ctx context.Context, c *app.RequestContext) {
	var req CommentListParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, err, nil)
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
		SendResponse(c, err, nil)
		return
	}

	SendResponse(c, errno.Success, CommentListResponse{
		CommentList: CommentResponse.CommentList,
	})
}