package service

import (
	"context"

	"tiktok/dal/mysql"
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
	cs, err := mysql.CommentList(s.ctx, req.VideoId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "mysql get comment list failed %v", err)
		return nil, err
	}

	comments := pack.Comments(cs)

	// TODO add user
	return comments, nil
}
