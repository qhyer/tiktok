package handler

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"tiktok/pkg/errno"
)

// Response 接口基础相应结构体
type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// SendResponse pack response
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
