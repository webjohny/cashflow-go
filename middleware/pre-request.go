package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
)

func PreRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := ctx.Request
		wr := ctx.Writer
		sessionStore, err := session.Start(context.Background(), wr, req)

		if err != nil {
			fmt.Fprint(wr, err)
			return
		}

		ctx.Set("sessionStore", sessionStore)
		ctx.Next()
	}
}
