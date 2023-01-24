package initialize

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"tiktok/cmd/api/handler"
	"tiktok/cmd/api/middleware"
	"tiktok/pkg/errno"
)

func Router() {
	h := server.New(
		server.WithHostPorts("127.0.0.1:8080"),
		server.WithHandleMethodNotAllowed(true),
	)
	h.Use(recovery.Recovery(recovery.WithRecoveryHandler(
		func(ctx context.Context, c *app.RequestContext, err interface{}, stack []byte) {
			hlog.SystemLogger().CtxErrorf(ctx, "[Recovery] err=%v\nstack=%s", err, stack)
			c.JSON(consts.StatusInternalServerError, handler.Response{
				StatusCode: errno.ServiceErrCode,
				StatusMsg:  fmt.Sprintf("[Recovery] err=%v\nstack=%s", err, stack),
			})
		})))

	// 接口路由
	apiRouter := h.Group("/douyin/")

	// 用户接口
	apiRouter.POST("/user/register/", handler.Register)
	apiRouter.POST("/user/login/", handler.Login)

	// 需要鉴权的接口路由
	authRouter := apiRouter.Group("/")
	// 中间件鉴权
	authRouter.Use(middleware.JWT())

	authRouter.GET("/user/", handler.GetUserInfo)
	// TODO 其余接口路由

	h.Spin()
}
