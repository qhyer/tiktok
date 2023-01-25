package pack

import (
	"context"
	"tiktok/cmd/feed/dal/db"
	"tiktok/cmd/feed/rpc"
	"tiktok/kitex_gen/feed"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/constants"
)

// Video pack video
func Video(video *db.Video, author *feed.User, isFavorite bool) *feed.Video {
	if video == nil || author == nil {
		return nil
	}

	return &feed.Video{
		Id:            video.Id,
		PlayUrl:       constants.OSSBaseURL + video.PlayUrl,
		CoverUrl:      constants.OSSBaseURL + video.CoverUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		Title:         video.Title,
		Author:        author,
		IsFavorite:    isFavorite,
	}
}

// Videos pack list of video
func Videos(ctx context.Context, vs []*db.Video, userId int64) ([]*feed.Video, int64) {
	var nextTime int64
	videos := make([]*feed.Video, 0)

	if len(vs) == 0 {
		return videos, 0
	}

	authorUserIds := make([]int64, 0)

	for _, v := range vs {
		authorUserIds = append(authorUserIds, v.AuthorUserId)
	}

	// 查询作者的用户信息
	userInfoResponse, err := rpc.UserInfo(ctx, &user.DouyinUserInfoRequest{
		UserId:    userId,
		ToUserIds: authorUserIds,
	})
	if err != nil {
		return nil, nextTime
	}

	userInfoMap := make(map[int64]*feed.User, 0)
	for _, us := range userInfoResponse.User {
		userInfoMap[us.Id] = &feed.User{
			Id:            us.Id,
			Name:          us.Name,
			FollowCount:   us.FollowCount,
			FollowerCount: us.FollowerCount,
			IsFollow:      us.IsFollow,
		}
	}

	// TODO 查询用户是否给视频点赞
	userFavoriteMap := make(map[int64]bool, 0)

	// 找到最后一个视频的创建时间
	nextTime = vs[len(vs)-1].CreatedAt.UnixMilli()

	// pack video
	for _, v := range vs {
		if vp := Video(v, userInfoMap[v.AuthorUserId], userFavoriteMap[v.Id]); vp != nil {
			videos = append(videos, vp)
		}
	}

	return videos, nextTime
}
