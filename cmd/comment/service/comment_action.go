package service

import (
	"context"

	"tiktok/dal/db"
	"tiktok/dal/pack"
	"tiktok/kitex_gen/comment"
	"tiktok/pkg/censor"
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
	// 过滤敏感词
	content := *req.CommentText
	content = censor.TextCensor.GetFilter().Replace(content, '*')

	// 插入数据
	c, err := db.CommentAction(s.ctx, &db.Comment{
		UserId:  req.UserId,
		VideoId: req.VideoId,
		Content: content,
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
