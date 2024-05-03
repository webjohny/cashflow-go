package middleware

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetGameId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Query("raceId") != "" {
			raceId, _ := strconv.Atoi(ctx.Query("raceId"))
			ctx.Set("raceId", raceId)
		} else if ctx.Param("raceId") != "" {
			raceId, _ := strconv.Atoi(ctx.Param("raceId"))
			ctx.Set("raceId", raceId)
		} else if ctx.Query("lobbyId") != "" {
			lobbyId, _ := strconv.Atoi(ctx.Query("lobbyId"))
			ctx.Set("lobbyId", lobbyId)
		} else if ctx.Param("lobbyId") != "" {
			raceId, _ := strconv.Atoi(ctx.Param("lobbyId"))
			ctx.Set("lobbyId", raceId)
		}

		ctx.Next()
	}
}
