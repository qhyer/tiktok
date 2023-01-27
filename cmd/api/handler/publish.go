package handler

import (
	"bytes"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/h2non/filetype"
	"io"
	"mime/multipart"
	"net/http"
	"tiktok/cmd/api/rpc"
	"tiktok/kitex_gen/feed"
	"tiktok/kitex_gen/publish"
	"tiktok/pkg/errno"
)

// TODO 参数校验
type PublishActionParam struct {
	Data  *multipart.FileHeader `form:"data"`
	Title string                `form:"title"`
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

func PublishAction(_ context.Context, c *app.RequestContext) {
	var req PublishActionParam
	err := c.BindAndValidate(&req)
	if err != nil {
		SendResponse(c, err)
		return
	}

	userId := c.GetInt64("UserID")

	// 读取视频
	videoFile, err := req.Data.Open()
	if err != nil {
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
		SendResponse(c, errno.ParamErr)
		return
	}

	// rpc通信
	_, err = rpc.PublishAction(context.Background(), &publish.DouyinPublishActionRequest{
		Data:   videoData.Bytes(),
		Title:  req.Title,
		UserId: userId,
	})

	if err != nil {
		SendResponse(c, err)
		return
	}

	SendResponse(c, errno.Success)
}

func PublishList(_ context.Context, c *app.RequestContext) {
	var req PublishListParam
	// 参数校验
	err := c.BindAndValidate(&req)
	if err != nil {
		SendResponse(c, err)
		return
	}

	userId := c.GetInt64("UserID")
	toUserId := req.UserId

	// rpc通信
	publishListResponse, err := rpc.PublishList(context.Background(), &publish.DouyinPublishListRequest{
		UserId:   userId,
		ToUserId: toUserId,
	})

	c.JSON(http.StatusOK, PublishListResponse{
		StatusCode: errno.Success.ErrCode,
		StatusMsg:  errno.Success.ErrMsg,
		VideoList:  publishListResponse.VideoList,
	})
}
