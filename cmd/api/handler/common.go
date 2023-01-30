package handler

import (
	"tiktok/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// SendResponse 封装返回结构体
func SendResponse(c *app.RequestContext, err error) {
	Err := errno.ConvertErr(err)

	// TODO 对外不应暴露服务器异常信息
	//if Err.ErrCode == errno.ServiceErrCode {
	//	Err.ErrMsg = "Server Internal error"
	//}

	c.JSON(consts.StatusOK, Response{
		StatusCode: Err.ErrCode,
		StatusMsg:  Err.ErrMsg,
	})
}
