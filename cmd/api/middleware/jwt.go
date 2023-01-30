package middleware

import (
	"context"

	"tiktok/cmd/api/handler"
	"tiktok/pkg/errno"
	"tiktok/pkg/jwt"

	"github.com/cloudwego/hertz/pkg/app"
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

		claims, err := jwt.ParseToken(tokenStr)
		if err != nil {
			handler.SendResponse(c, errno.AuthorizationFailedErr)
			c.Abort()
			return
		}
		userID := claims.UserID

		// 请求中加入userID
		c.Set("UserID", userID)
		c.Next(ctx)
	}
}
