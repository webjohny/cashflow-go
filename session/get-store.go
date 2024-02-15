package session

import (
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
)

func GetStore(ctx *gin.Context) session.Store {
	return ctx.MustGet("sessionStore").(session.Store)
}
