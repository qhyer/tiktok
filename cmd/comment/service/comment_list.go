package service

import (
	"context"

	"tiktok/dal/db"
	"tiktok/dal/pack"
	"tiktok/kitex_gen/comment"

	"github.com/cloudwego/kitex/pkg/klog"
)

type CommentListService struct {
	ctx context.Context
}

// NewCommentListService new CommentListService
func NewCommentListService(ctx context.Context) *CommentListService {
	return &CommentListService{ctx: ctx}
}

func (s *CommentListService) CommentList(req *comment.DouyinCommentListRequest) ([]*comment.Comment, error) {
	cs, err := db.CommentList(s.ctx, req.VideoId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "db get comment list failed %v", err)
		return nil, err
	}

	comments, err := pack.Comments(s.ctx, cs, req.UserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "pack comment list failed %v", err)
		return nil, err
	}

	return comments, nil
}
