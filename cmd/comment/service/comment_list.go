package service

import (
	"context"

	"tiktok/dal/mysql"
	"tiktok/dal/pack"
	"tiktok/kitex_gen/comment"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/rpc"

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

	if len(comments) == 0 {
		return nil, nil
	}

	// 查询用户信息
	userIds := make([]int64, 0, len(comments))
	for _, v := range comments {
		userIds = append(userIds, v.User.Id)
	}

	users, err := rpc.UserInfo(s.ctx, &user.DouyinUserInfoRequest{
		UserId:    req.UserId,
		ToUserIds: userIds,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "rpc get userinfo failed %v", err)
		return nil, err
	}

	// 加入用户信息
	for i := range comments {
		comments[i].User = users.User[i]
	}

	return comments, nil
}
