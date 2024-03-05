package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/storage"
	"gopkg.in/errgo.v2/errors"
	"strconv"
)

func GetGameId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := ctx.Request

		if req.Header.Get("X-Race-ID") != "" {
			raceId, _ := strconv.Atoi(req.Header.Get("X-Race-ID"))
			ctx.Set("raceId", raceId)
		} else if req.Header.Get("X-Lobby-ID") != "" {
			lobbyId, _ := strconv.Atoi(req.Header.Get("X-Lobby-ID"))
			ctx.Set("lobbyId", lobbyId)
		} else {
			request.FinalResponse(ctx, errors.New(storage.ErrorUndefinedGame), nil)
			return
		}

		ctx.Next()
	}
}
