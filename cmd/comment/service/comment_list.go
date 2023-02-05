package service

import (
	"context"

	"tiktok/cmd/rpc"
	"tiktok/dal/mysql"
	"tiktok/dal/pack"
	"tiktok/dal/redis"
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

	// 从缓存中读评论id列表，其中没读到缓存会查库
	redisComments, err := redis.GetCommentIdListByVideoId(s.ctx, videoId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "redis get comment id list failed %v", err)
		return nil, err
	}

	commentMap := make(map[int64]*comment.Comment, 0)
	comments := make([]*comment.Comment, 0)

	// 缓存中查评论详情
	rcs, notInCacheCommentIds := redis.MGetCommentByCommentId(s.ctx, redisComments)
	packRedisComments := pack.Comments(rcs)
	for _, c := range packRedisComments {
		if c == nil {
			continue
		}
		commentMap[c.Id] = c
	}

	// 缓存没查到 查库
	if len(notInCacheCommentIds) > 0 {
		cs, err := mysql.GetCommentListByVideoId(s.ctx, videoId)
		if err != nil {
			klog.CtxErrorf(s.ctx, "mysql get comment list failed %v", err)
			return nil, err
		}

		// 把评论加入缓存
		err = redis.MSetComment(s.ctx, cs)
		if err != nil {
			klog.CtxErrorf(s.ctx, "redis set comment failed %v", err)
		}

		// 把评论id加入缓存
		err = redis.MAddCommentIdToCommentList(s.ctx, cs, videoId)
		if err != nil {
			klog.CtxErrorf(s.ctx, "redis add comment id to comment list failed %v", err)
		}

		// 把评论放入map中
		coms := pack.Comments(cs)
		for _, c := range coms {
			if c == nil {
				continue
			}
			commentMap[c.Id] = c
		}
	}

	// 合并评论
	for _, c := range redisComments {
		if c == nil {
			continue
		}
		res := commentMap[c.Id]
		if res == nil {
			continue
		}
		comments = append(comments, res)
	}

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
