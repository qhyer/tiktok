package initialize

import (
	"context"
	"fmt"

	"tiktok/cmd/api/handler"
	"tiktok/cmd/api/middleware"
	"tiktok/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
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

	// 视频流接口
	apiRouter.GET("/feed/", handler.Feed)

	// 用户接口
	apiRouter.POST("/user/register/", handler.Register)
	apiRouter.POST("/user/login/", handler.Login)
	apiRouter.GET("/user/", handler.GetUserInfo)

	// 需要鉴权的接口路由
	authRouter := apiRouter.Group("/")
	// 中间件鉴权
	authRouter.Use(middleware.JWT())

	// 投稿路由
	authRouter.POST("/publish/action/", handler.PublishAction)

	// 发布列表路由
	apiRouter.GET("/publish/list/", handler.PublishList)

	// 点赞路由
	authRouter.POST("/favorite/action/", handler.FavoriteAction)

	// 喜欢列表路由
	apiRouter.GET("/favorite/list/", handler.FavoriteList)

	// 评论路由
	authRouter.POST("/comment/action/", handler.CommentAction)

	// 评论列表路由
	apiRouter.GET("/comment/list/", handler.CommentList)

	// TODO 其余接口路由
	h.Spin()
}
