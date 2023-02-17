package service

import (
	"context"

	"tiktok/cmd/rpc"
	"tiktok/dal/mysql"
	"tiktok/dal/pack"
	"tiktok/dal/redis"
	"tiktok/kitex_gen/comment"
	"tiktok/kitex_gen/user"
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

	// 数据库中插入数据
	c, err := mysql.CreateComment(s.ctx, &mysql.Comment{
		UserId:  userId,
		VideoId: videoId,
		Content: content,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "mysql create comment failed %v", err)
		return nil, err
	}

	com := pack.Comment(c)

	// 获取发布者信息
	userResp, err := rpc.UserInfo(s.ctx, &user.DouyinUserInfoRequest{
		UserId:    userId,
		ToUserIds: []int64{userId},
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "rpc get userinfo failed %v", err)
		return nil, err
	}
	com.User = userResp.User[0]

	// 在缓存中加入评论数据
	err = redis.SetComment(s.ctx, c)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis set comment failed %v", err)
		return nil, err
	}

	// 更新缓存评论列表
	err = redis.AddNewCommentToCommentList(s.ctx, c)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis add new comment failed %v", err)
		return nil, err
	}

	return com, nil
}

// DeleteComment delete user comment action
func (s *CommentActionService) DeleteComment(req *comment.DouyinCommentActionRequest) error {
	userId := req.GetUserId()
	commentId := req.GetCommentId()

	// 在数据库中删除评论
	com, err := mysql.DeleteComment(s.ctx, &mysql.Comment{
		UserId: userId,
		Id:     commentId,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "mysql delete comment failed %v", err)
		return err
	}

	// 从缓存中删除评论
	err = redis.DeleteCommentFromCommentList(s.ctx, com)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis delete comment from comment list failed %v", err)
		return err
	}

	return nil
}
