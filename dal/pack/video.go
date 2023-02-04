package pack

import (
	"time"

	"tiktok/dal/mysql"
	"tiktok/kitex_gen/feed"
	"tiktok/kitex_gen/user"
)

// Video pack video
func Video(video *mysql.Video) *feed.Video {
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
func Videos(vs []*mysql.Video) ([]*feed.Video, int64) {
	videos := make([]*feed.Video, 0, len(vs))

	if len(vs) == 0 {
		return videos, 0
	}

	earliestTime := time.Now().UnixMilli()

	// pack video
	for _, v := range vs {
		if vp := Video(v); vp != nil {
			ts := v.CreatedAt.UnixMilli()
			if ts != 0 && ts < earliestTime {
				earliestTime = ts
			}

			ts = v.CreatedTimestamp
			if ts != 0 && ts < earliestTime {
				earliestTime = ts
			}

			videos = append(videos, vp)
		}
	}

	return videos, earliestTime
}
