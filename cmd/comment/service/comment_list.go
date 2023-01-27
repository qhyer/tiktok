package service

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"tiktok/cmd/comment/dal/db"
	"tiktok/cmd/comment/pack"
	"tiktok/kitex_gen/comment"
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
		klog.CtxFatalf(s.ctx, "db get comment list failed %v", err)
		return nil, err
	}

	comments, err := pack.Comments(s.ctx, cs, req.UserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "pack comment list failed %v", err)
		return nil, err
	}

	return comments, nil
}
