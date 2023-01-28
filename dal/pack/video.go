package pack

import (
	"tiktok/dal/db"
	"tiktok/kitex_gen/feed"
	"tiktok/kitex_gen/user"
)

// Video pack video
func Video(video *db.Video) *feed.Video {
	if video == nil {
		return nil
	}

	return &feed.Video{
		Id:            video.Id,
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		Title:         video.Title,
		Author: &user.User{
			Id: video.AuthorUserId,
		},
	}
}

// Videos pack list of video
func Videos(vs []*db.Video) ([]*feed.Video, int64) {
	var nextTime int64
	videos := make([]*feed.Video, 0, len(vs))

	if len(vs) == 0 {
		return videos, 0
	}

	// 找到最后一个视频的创建时间
	nextTime = vs[len(vs)-1].CreatedAt.UnixMilli()

	// pack video
	for _, v := range vs {
		if vp := Video(v); vp != nil {
			videos = append(videos, vp)
		}
	}

	return videos, nextTime
}
