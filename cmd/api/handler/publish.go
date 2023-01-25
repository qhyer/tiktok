package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/h2non/filetype"
	"mime/multipart"
	"tiktok/pkg/errno"
)

type PublishParam struct {
	Data *multipart.FileHeader `form:"data"`
}

func PublishAction(ctx context.Context, c *app.RequestContext) {
	var req PublishParam
	err := c.BindAndValidate(&req)
	if err != nil {
		SendResponse(c, err)
		return
	}

	// 读取视频
	videoFile, err := req.Data.Open()
	videoData := make([]byte, req.Data.Size)
	read, err := videoFile.Read(videoData)
	if err != nil {
		SendResponse(c, errno.ParamErr)
		return
	}

	// 校验视频文件合法性
	if !filetype.IsVideo(videoData) {
		SendResponse(c, errno.ParamErr)
		return
	}

}
