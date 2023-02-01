package initialize

import (
	"context"

	"tiktok/cmd/api/handler"
	"tiktok/cmd/api/middleware"
	"tiktok/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Router() {
	h := server.New(
		server.WithHostPorts("127.0.0.1:8080"),
		server.WithHandleMethodNotAllowed(true),
		server.WithMaxRequestBodySize(104857600), // 100MB
	)
	h.Use(recovery.Recovery(recovery.WithRecoveryHandler(
		func(ctx context.Context, c *app.RequestContext, err interface{}, stack []byte) {
			hlog.SystemLogger().CtxErrorf(ctx, "[Recovery] err=%v\nstack=%s", err, stack)
			handler.SendResponse(c, errno.ServiceErr)
		})))

	// 接口路由
	apiRouter := h.Group("/douyin/")
	// 中间件
	apiRouter.Use(middleware.Common())

	// 需要鉴权的接口路由
	authRouter := apiRouter.Group("/")
	// 中间件鉴权
	authRouter.Use(middleware.Auth())

	{
		// 视频流
		apiRouter.GET("/feed/", handler.Feed)
	}
	{
		// 用户注册
		apiRouter.POST("/user/register/", handler.Register)
		// 用户登录
		apiRouter.POST("/user/login/", handler.Login)
		// 用户信息
		apiRouter.GET("/user/", handler.GetUserInfo)
	}
	{
		// 投稿
		authRouter.POST("/publish/action/", handler.PublishAction)
		// 发布列表
		apiRouter.GET("/publish/list/", handler.PublishList)
	}
	{
		// 点赞操作
		authRouter.POST("/favorite/action/", handler.FavoriteAction)
		// 喜欢列表
		apiRouter.GET("/favorite/list/", handler.FavoriteList)
	}
	{
		// 评论操作
		authRouter.POST("/comment/action/", handler.CommentAction)
		// 评论列表
		apiRouter.GET("/comment/list/", handler.CommentList)
	}
	{
		// 关注路由
		authRouter.POST("/relation/action/", handler.RelationAction)
		// 关注列表
		apiRouter.GET("/relation/follow/list/", handler.FollowList)
		// 粉丝列表
		apiRouter.GET("/relation/follower/list", handler.FollowerList)
		// 好友列表
		authRouter.GET("/relation/friend/list/", handler.FriendList)
	}
	{
		// 发送消息
		authRouter.POST("/message/action/", handler.MessageAction)
		// 聊天记录
		authRouter.GET("/message/chat/", handler.MessageList)
	}

	h.Spin()
}
