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
		if ctx.Query("raceId") != "" {
			raceId, _ := strconv.Atoi(ctx.Query("raceId"))
			ctx.Set("raceId", raceId)
		} else if request.GetLobbyId(ctx) == 0 {
			request.FinalResponse(ctx, errors.New(storage.ErrorUndefinedGame), nil)
			return
		}

		ctx.Next()
	}
}
