package pack

import (
	"context"
	"time"

	"tiktok/dal/db"
	"tiktok/kitex_gen/comment"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/rpc"

	"github.com/cloudwego/kitex/pkg/klog"
)

// Comment pack comment
func Comment(com *db.Comment) *comment.Comment {
	if com == nil {
		return nil
	}
	return &comment.Comment{
		Id:         com.Id,
		Content:    com.Content,
		CreateDate: com.UpdatedAt.In(time.FixedZone("CST", 8*3600)).Format("01-02"),
	}
}

// Comments pack list of comment
func Comments(ctx context.Context, cms []*db.Comment, userId int64) ([]*comment.Comment, error) {
	comments := make([]*comment.Comment, 0, len(cms))
	if len(cms) == 0 {
		return comments, nil
	}

	commentUserIds := make([]int64, 0, len(cms))

	for _, c := range cms {
		commentUserIds = append(commentUserIds, c.UserId)
	}

	// 查询评论者的用户信息
	userInfoResponse, err := rpc.UserInfo(ctx, &user.DouyinUserInfoRequest{
		UserId:    userId,
		ToUserIds: commentUserIds,
	})
	if err != nil {
		klog.CtxErrorf(ctx, "query userinfo failed %v", userInfoResponse)
		return nil, err
	}

	userInfoMap := make(map[int64]*user.User, 0)
	for _, us := range userInfoResponse.User {
		userInfoMap[us.Id] = &user.User{
			Id:            us.Id,
			Name:          us.Name,
			FollowCount:   us.FollowCount,
			FollowerCount: us.FollowerCount,
			IsFollow:      us.IsFollow,
		}
	}

	// pack comment
	for _, c := range cms {
		if vp := Comment(c); vp != nil {
			if us := userInfoMap[c.UserId]; us != nil {
				vp.User = us
				comments = append(comments, vp)
			}
		}
	}

	return comments, err
}
