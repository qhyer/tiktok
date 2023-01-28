package handler

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"

	"tiktok/kitex_gen/feed"
	"tiktok/kitex_gen/publish"
	"tiktok/pkg/errno"
	"tiktok/pkg/rpc"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/h2non/filetype"
)

// TODO 参数校验
type PublishActionParam struct {
	Data  *multipart.FileHeader `form:"data"`
	Title string                `form:"title"`
}

type PublishActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// TODO 参数校验
type PublishListParam struct {
	UserId int64 `query:"user_id"`
}

type PublishListResponse struct {
	StatusCode int32         `json:"status_code"`
	StatusMsg  string        `json:"status_msg"`
	VideoList  []*feed.Video `json:"video_list"`
}

func PublishAction(ctx context.Context, c *app.RequestContext) {
	var req PublishActionParam
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, err)
		return
	}

	userId := c.GetInt64("UserID")

	// 读取视频
	videoFile, err := req.Data.Open()
	if err != nil {
		hlog.CtxWarnf(ctx, "read video error %v", err)
		SendResponse(c, errno.ParamErr)
		return
	}
	videoData := bytes.NewBuffer(nil)
	if _, err := io.Copy(videoData, videoFile); err != nil {
		SendResponse(c, err)
		return
	}

	// 校验视频文件合法性
	if !filetype.IsVideo(videoData.Bytes()) {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, errno.ParamErr)
		return
	}

	// rpc通信
	publishActionResponse, err := rpc.PublishAction(ctx, &publish.DouyinPublishActionRequest{
		Data:   videoData.Bytes(),
		Title:  req.Title,
		UserId: userId,
	})

	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, PublishActionResponse{
		StatusCode: publishActionResponse.StatusCode,
		StatusMsg:  *publishActionResponse.StatusMsg,
	})
}

func PublishList(ctx context.Context, c *app.RequestContext) {
	var req PublishListParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		hlog.CtxWarnf(ctx, "param error %v", err)
		SendResponse(c, err)
		return
	}

	userId := c.GetInt64("UserID")
	toUserId := req.UserId

	// rpc通信
	publishListResponse, err := rpc.PublishList(ctx, &publish.DouyinPublishListRequest{
		UserId:   userId,
		ToUserId: toUserId,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc response error %v", err)
		SendResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, PublishListResponse{
		StatusCode: publishListResponse.StatusCode,
		StatusMsg:  *publishListResponse.StatusMsg,
		VideoList:  publishListResponse.VideoList,
	})
}
