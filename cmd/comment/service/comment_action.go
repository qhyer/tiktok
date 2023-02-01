package service

import (
	"context"

	"tiktok/dal/mysql"
	"tiktok/dal/pack"
	"tiktok/kitex_gen/comment"
	"tiktok/pkg/censor"

	"github.com/cloudwego/kitex/pkg/klog"
)

type CommentActionService struct {
	ctx context.Context
}

// NewCommentActionService new CommentActionService
func NewCommentActionService(ctx context.Context) *CommentActionService {
	return &CommentActionService{ctx: ctx}
}

// CreateComment user comment video action
func (s *CommentActionService) CreateComment(req *comment.DouyinCommentActionRequest) (*comment.Comment, error) {
	content := req.GetCommentText()
	userId := req.GetUserId()
	videoId := req.GetVideoId()

	// 过滤敏感词
	content = censor.TextCensor.GetFilter().Replace(content, '*')

	// 插入数据
	c, err := mysql.CreateComment(s.ctx, &mysql.Comment{
		UserId:  userId,
		VideoId: videoId,
		Content: content,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "db create comment failed %v", err)
		return nil, err
	}

	com := pack.Comment(c)
	return com, nil
}

// DeleteComment delete user comment action
func (s *CommentActionService) DeleteComment(req *comment.DouyinCommentActionRequest) error {
	userId := req.GetUserId()
	commentId := req.GetCommentId()

	err := mysql.DeleteComment(s.ctx, &mysql.Comment{
		UserId: userId,
		Id:     commentId,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "db delete comment failed %v", err)
		return err
	}

	return nil
}
