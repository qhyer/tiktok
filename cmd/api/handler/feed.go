package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"tiktok/cmd/api/rpc"
	"tiktok/kitex_gen/feed"
	"tiktok/pkg/errno"
	"time"
)

// TODO: 参数校验
type FeedParam struct {
	LatestTime int64 `query:"latest_time"`
}

type FeedResponse struct {
	StatusCode int32   `json:"status_code"`
	StatusMsg  string  `json:"status_msg"`
	VideoList  []Video `json:"video_list"`
	NextTime   int64   `json:"next_time"`
}

type Video struct {
	Id            int64  `json:"id"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}

func Feed(_ context.Context, c *app.RequestContext) {
	var req FeedParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		SendResponse(c, err)
		return
	}
	// 没有传参 把最后时间设为当前时间
	if req.LatestTime == 0 {
		req.LatestTime = time.Now().UnixMilli()
	}
	userId := c.GetInt64("UserID")

	// rpc通信
	feedResponse, err := rpc.Feed(context.Background(), &feed.DouyinFeedRequest{
		LatestTime: &req.LatestTime,
		UserId:     userId,
	})
	if err != nil {
		SendResponse(c, err)
		return
	}

	videoList := make([]Video, 0)
	for _, v := range feedResponse.VideoList {
		videoList = append(videoList, Video{
			Id: v.Id,
			Author: User{
				Id:            v.Author.Id,
				IsFollow:      v.Author.IsFollow,
				FollowerCount: *v.Author.FollowerCount,
				FollowCount:   *v.Author.FollowCount,
				Name:          v.Author.Name,
			},
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			Title:         v.Title,
		})
	}

	c.JSON(http.StatusOK, FeedResponse{
		StatusCode: errno.Success.ErrCode,
		StatusMsg:  errno.Success.ErrMsg,
		VideoList:  videoList,
		NextTime:   *feedResponse.NextTime,
	})
}
