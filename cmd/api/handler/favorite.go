package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"tiktok/cmd/api/rpc"
	"tiktok/kitex_gen/favorite"
	"tiktok/kitex_gen/feed"
)

// TODO 参数校验
type FavoriteActionParam struct {
	VideoId    int64 `query:"video_id"`
	ActionType int32 `query:"action_type"`
}

type FavoriteActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func FavoriteAction(_ context.Context, c *app.RequestContext) {
	var req FavoriteActionParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		SendResponse(c, err)
		return
	}
	userId := c.GetInt64("UserID")

	// rpc通信
	favoriteResponse, err := rpc.FavoriteAction(context.Background(), &favorite.DouyinFavoriteActionRequest{
		UserId:     userId,
		VideoId:    req.VideoId,
		ActionType: req.ActionType,
	})
	if err != nil {
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, FavoriteActionResponse{
		StatusCode: favoriteResponse.StatusCode,
		StatusMsg:  *favoriteResponse.StatusMsg,
	})
}

// TODO 参数校验
type FavoriteListParam struct {
	userId int64 `query:"user_id"`
}

type FavoriteListResponse struct {
	StatusCode int32         `json:"status_code"`
	StatusMsg  string        `json:"status_msg"`
	VideoList  []*feed.Video `json:"video_list"`
}

func FavoriteList(_ context.Context, c *app.RequestContext) {
	var req FavoriteListParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		SendResponse(c, err)
		return
	}
	userId := c.GetInt64("UserID")

	// rpc通信
	favoriteResponse, err := rpc.FavoriteList(context.Background(), &favorite.DouyinFavoriteListRequest{
		UserId: userId,
	})
	if err != nil {
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, FavoriteListResponse{
		StatusCode: favoriteResponse.StatusCode,
		StatusMsg:  *favoriteResponse.StatusMsg,
		VideoList:  favoriteResponse.VideoList,
	})
}
