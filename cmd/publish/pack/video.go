package pack

import (
	"context"
	"tiktok/cmd/publish/dal/db"
	"tiktok/cmd/publish/rpc"
	"tiktok/kitex_gen/publish"
	"tiktok/kitex_gen/user"
	"tiktok/pkg/constants"
)

// Video pack video
func Video(video *db.Video, author *publish.User, isFavorite bool) *publish.Video {
	if video == nil || author == nil {
		return nil
	}

	return &publish.Video{
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
func Videos(ctx context.Context, vs []*db.Video, userId int64) []*publish.Video {
	videos := make([]*publish.Video, 0, len(vs))

	if len(vs) == 0 {
		return videos
	}

	authorUserIds := make([]int64, 0, len(vs))

	for _, v := range vs {
		authorUserIds = append(authorUserIds, v.AuthorUserId)
	}

	// 查询作者的用户信息
	userInfoResponse, err := rpc.UserInfo(ctx, &user.DouyinUserInfoRequest{
		UserId:    userId,
		ToUserIds: authorUserIds,
	})
	if err != nil {
		return nil
	}

	userInfoMap := make(map[int64]*publish.User, 0)
	for _, us := range userInfoResponse.User {
		userInfoMap[us.Id] = &publish.User{
			Id:            us.Id,
			Name:          us.Name,
			FollowCount:   us.FollowCount,
			FollowerCount: us.FollowerCount,
			IsFollow:      us.IsFollow,
		}
	}

	// TODO 查询用户是否给视频点赞
	userFavoriteMap := make(map[int64]bool, 0)

	// pack video
	for _, v := range vs {
		if vp := Video(v, userInfoMap[v.AuthorUserId], userFavoriteMap[v.Id]); vp != nil {
			videos = append(videos, vp)
		}
	}

	return videos
}
