package initialize

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"tiktok/cmd/api/middleware"
)

func Router() {
	h := server.Default()

	// 接口路由
	apiRouter := h.Group("/douyin/")

	// TODO: 用户注册和登录接口
	apiRouter.POST("/user/login/")

	// 需要鉴权的接口路由
	authRouter := h.Group("/")
	// 中间件鉴权
	authRouter.Use(middleware.JWT())
	// TODO: 其余接口路由

	h.Spin()
}
