package service

import (
	"context"
	"tiktok/cmd/comment/dal/db"
	"tiktok/kitex_gen/comment"
)

type CommentActionService struct {
	ctx context.Context
}

// NewCommentActionService new CommentActionService
func NewCommentActionService(ctx context.Context) *CommentActionService {
	return &CommentActionService{ctx: ctx}
}

// CommentAction user comment video action
func (s *CommentActionService) CommentAction(req *comment.DouyinCommentActionRequest) (*comment.Comment, error) {
	// TODO 文本内容审核 返回评论内容
	err := db.CommentAction(s.ctx, &db.Comment{
		UserId:  req.UserId,
		VideoId: req.VideoId,
		Content: *req.CommentText,
	})
	return nil, err
}

// DeleteCommentAction delete user comment action
func (s *CommentActionService) DeleteCommentAction(req *comment.DouyinCommentActionRequest) error {
	err := db.DeleteCommentAction(s.ctx, &db.Comment{
		UserId: req.UserId,
		Id:     *req.CommentId,
	})
	return err
}
