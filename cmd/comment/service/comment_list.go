package service

import (
	"context"

	"tiktok/cmd/rpc"
	"tiktok/dal/mysql"
	"tiktok/dal/pack"
	"tiktok/kitex_gen/comment"
	"tiktok/kitex_gen/user"

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
	videoId := req.GetVideoId()
	userId := req.GetUserId()

	cs, err := mysql.GetCommentListByVideoId(s.ctx, videoId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "mysql get comment list failed %v", err)
		return nil, err
	}

	comments := pack.Comments(cs)

	if len(comments) == 0 {
		return comments, nil
	}

	// 查询用户信息
	userIds := make([]int64, 0, len(comments))
	for _, v := range comments {
		if v == nil || v.User == nil {
			continue
		}
		userIds = append(userIds, v.User.Id)
	}

	users, err := rpc.UserInfo(s.ctx, &user.DouyinUserInfoRequest{
		UserId:    userId,
		ToUserIds: userIds,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "rpc get userinfo failed %v", err)
		return nil, err
	}

	// 加入用户信息
	us := users.GetUser()
	for i := range comments {
		if us[i] == nil {
			continue
		}
		comments[i].User = us[i]
	}

	return comments, nil
}
