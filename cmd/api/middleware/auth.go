package middleware

import (
	"context"

	"tiktok/cmd/api/handler"
	"tiktok/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
)

func Auth() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		if c.GetBool("Auth") != true {
			handler.SendResponse(c, errno.AuthorizationFailedErr)
			c.Abort()
			return
		}
		c.Next(ctx)
	}
}
