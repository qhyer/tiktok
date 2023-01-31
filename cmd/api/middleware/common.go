package middleware

import (
	"context"
	"time"

	"tiktok/pkg/jwt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Common() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()

		// 获取当前用户id
		userId := int64(0)
		isAuth := false
		tokenStr := c.PostForm("token")
		if tokenStr == "" {
			tokenStr = c.Query("token")
		}
		if tokenStr != "" {
			claims, err := jwt.ParseToken(tokenStr)
			if err == nil {
				userId = claims.UserID
				isAuth = true
			}
		}
		c.Set("UserID", userId)
		c.Set("Auth", isAuth)

		c.Next(ctx)

		end := time.Now()
		latency := end.Sub(start).Microseconds
		hlog.CtxTracef(ctx, "status=%d cost=%d method=%s full_path=%s client_ip=%s host=%s",
			c.Response.StatusCode(), latency,
			c.Request.Header.Method(), c.Request.URI().PathOriginal(), c.ClientIP(), c.Request.Host())
	}
}
