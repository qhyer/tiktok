package handler

import (
	"context"
	"net/http"

	"tiktok/kitex_gen/favorite"
	"tiktok/kitex_gen/feed"
	"tiktok/pkg/rpc"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type FavoriteActionParam struct {
	VideoId    int64 `query:"video_id" vd:"$!=nil&&$>0"`
	ActionType int32 `query:"action_type" vd:"$!=nil&&$==1||$==2"`
}

type FavoriteActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func FavoriteAction(ctx context.Context, c *app.RequestContext) {
	var req FavoriteActionParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, err)
		return
	}
	userId := c.GetInt64("UserID")

	// rpc通信
	favoriteResponse, err := rpc.FavoriteAction(ctx, &favorite.DouyinFavoriteActionRequest{
		UserId:     userId,
		VideoId:    req.VideoId,
		ActionType: req.ActionType,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, FavoriteActionResponse{
		StatusCode: favoriteResponse.StatusCode,
		StatusMsg:  *favoriteResponse.StatusMsg,
	})
}

type FavoriteListParam struct {
	userId int64 `query:"user_id" vd:"$!=nil&&$>0"`
}

type FavoriteListResponse struct {
	StatusCode int32         `json:"status_code"`
	StatusMsg  string        `json:"status_msg"`
	VideoList  []*feed.Video `json:"video_list"`
}

func FavoriteList(ctx context.Context, c *app.RequestContext) {
	var req FavoriteListParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, err)
		return
	}
	userId := c.GetInt64("UserID")

	// rpc通信
	favoriteResponse, err := rpc.FavoriteList(ctx, &favorite.DouyinFavoriteListRequest{
		UserId: userId,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, FavoriteListResponse{
		StatusCode: favoriteResponse.StatusCode,
		StatusMsg:  *favoriteResponse.StatusMsg,
		VideoList:  favoriteResponse.VideoList,
	})
}
