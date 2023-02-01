package handler

import (
	"context"
	"net/http"

	"tiktok/cmd/rpc"
	"tiktok/kitex_gen/favorite"
	"tiktok/kitex_gen/feed"
	"tiktok/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type FavoriteActionParam struct {
	VideoId    int64 `query:"video_id" vd:"$>0"`
	ActionType int32 `query:"action_type" vd:"$==1||$==2"`
}

type FavoriteListParam struct {
	userId int64 `query:"user_id" vd:"$>0"`
}

type FavoriteListResponse struct {
	StatusCode int32         `json:"status_code"`
	StatusMsg  string        `json:"status_msg"`
	VideoList  []*feed.Video `json:"video_list"`
}

// FavoriteAction 关注、取关
func FavoriteAction(ctx context.Context, c *app.RequestContext) {
	var req FavoriteActionParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, errno.ParamErr)
		return
	}
	userId := c.GetInt64("UserID")

	// rpc通信
	_, err = rpc.FavoriteAction(ctx, &favorite.DouyinFavoriteActionRequest{
		UserId:     userId,
		VideoId:    req.VideoId,
		ActionType: req.ActionType,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err)
		return
	}

	SendResponse(c, errno.Success)
}

// FavoriteList 获取个人点赞列表
func FavoriteList(ctx context.Context, c *app.RequestContext) {
	var req FavoriteListParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, errno.ParamErr)
		return
	}

	userId := c.GetInt64("UserID")

	// rpc通信
	favoriteResponse, err := rpc.FavoriteList(ctx, &favorite.DouyinFavoriteListRequest{
		UserId:   userId,
		ToUserId: req.userId,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, FavoriteListResponse{
		StatusCode: errno.Success.ErrCode,
		StatusMsg:  errno.Success.ErrMsg,
		VideoList:  favoriteResponse.GetVideoList(),
	})
}
