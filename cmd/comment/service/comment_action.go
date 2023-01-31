package service

import (
	"context"

	"tiktok/dal/mysql"
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
	content := req.GetCommentText()
	userId := req.GetUserId()
	videoId := req.GetVideoId()

	// 过滤敏感词
	content = censor.TextCensor.GetFilter().Replace(content, '*')

	// 插入数据
	c, err := mysql.CommentAction(s.ctx, &mysql.Comment{
		UserId:  userId,
		VideoId: videoId,
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
	userId := req.GetUserId()
	commentId := req.GetCommentId()

	err := mysql.DeleteCommentAction(s.ctx, &mysql.Comment{
		UserId: userId,
		Id:     commentId,
	})
	return err
}
