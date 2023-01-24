package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"tiktok/cmd/api/handler"
	"tiktok/cmd/api/util"
	"tiktok/pkg/errno"
)

func JWT() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		tokenStr := c.PostForm("token")
		if tokenStr == "" {
			tokenStr = c.Query("token")
		}

		// token不能为空
		if tokenStr == "" {
			handler.SendResponse(c, errno.ParamErr)
			c.Abort()
			return
		}

		claims, err := util.ParseToken(tokenStr)
		if err != nil {
			handler.SendResponse(c, errno.AuthorizationFailedErr)
			c.Abort()
			return
		}
		userID := claims.UserID

		// 请求中加入userID
		c.Set("UserID", userID)

	}
}
