package initialize

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"tiktok/cmd/api/handler"
	"tiktok/cmd/api/middleware"
)

func Router() {
	h := server.Default()

	// 接口路由
	apiRouter := h.Group("/douyin/")

	// 用户接口
	apiRouter.GET("/user/", handler.GetUserInfo)
	apiRouter.POST("/user/register/", handler.Register)
	apiRouter.POST("/user/login/", handler.Login)

	// 需要鉴权的接口路由
	authRouter := h.Group("/")
	// 中间件鉴权
	authRouter.Use(middleware.JWT())
	// TODO 其余接口路由

	h.Spin()
}
