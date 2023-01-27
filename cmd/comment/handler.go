package main

import (
	"context"
	"tiktok/cmd/comment/pack"
	"tiktok/cmd/comment/service"
	"tiktok/kitex_gen/comment"
	"tiktok/pkg/errno"
)

// CommentSrvImpl implements the last service interface defined in the IDL.
type CommentSrvImpl struct{}

const (
	DoCommentAction     = 1
	DeleteCommentAction = 2
)

// CommentAction implements the CommentSrvImpl interface.
func (s *CommentSrvImpl) CommentAction(ctx context.Context, req *comment.DouyinCommentActionRequest) (resp *comment.DouyinCommentActionResponse, err error) {
	resp = new(comment.DouyinCommentActionResponse)
	var newComment *comment.Comment

	switch req.ActionType {
	case DoCommentAction:
		newComment, err = service.NewCommentActionService(ctx).CommentAction(req)
	case DeleteCommentAction:
		err = service.NewCommentActionService(ctx).DeleteCommentAction(req)
	default:
		err = errno.ParamErr
	}
	if err != nil {
		resp = pack.BuildCommentActionResp(err)
		return resp, err
	}

	resp = pack.BuildCommentActionResp(errno.Success)
	resp.Comment = newComment
	return
}

// CommentList implements the CommentSrvImpl interface.
func (s *CommentSrvImpl) CommentList(ctx context.Context, req *comment.DouyinCommentListRequest) (resp *comment.DouyinCommentListResponse, err error) {
	resp = new(comment.DouyinCommentListResponse)

	comments, err := service.NewCommentListService(ctx).CommentList(req)
	if err != nil {
		resp = pack.BuildCommentListResp(err)
		return resp, err
	}

	resp = pack.BuildCommentListResp(errno.Success)
	resp.CommentList = comments
	return resp, nil
}
