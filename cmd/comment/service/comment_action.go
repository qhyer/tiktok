package service

import (
	"context"

	"tiktok/dal/db"
	"tiktok/dal/pack"
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
	c, err := db.CommentAction(s.ctx, &db.Comment{
		UserId:  req.UserId,
		VideoId: req.VideoId,
		Content: *req.CommentText,
	})
	if err != nil {
		return nil, err
	}

	com := pack.Comment(c)
	return com, err
}

// DeleteCommentAction delete user comment action
func (s *CommentActionService) DeleteCommentAction(req *comment.DouyinCommentActionRequest) error {
	err := db.DeleteCommentAction(s.ctx, &db.Comment{
		UserId: req.UserId,
		Id:     *req.CommentId,
	})
	return err
}
