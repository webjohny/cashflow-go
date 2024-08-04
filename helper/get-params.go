package helper

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetUserId(ctx *gin.Context) uint64 {
	userIdParam, _ := strconv.Atoi(ctx.GetString("userId"))
	return uint64(userIdParam)
}

func GetRaceId(ctx *gin.Context) uint64 {
	return uint64(ctx.GetInt("raceId"))
}

func GetLobbyId(ctx *gin.Context) uint64 {
	lobbyId, _ := strconv.Atoi(ctx.Param("lobbyId"))
	return uint64(lobbyId)
}

func GetBigRace(ctx *gin.Context) bool {
	bigRaceQuery := ctx.Query("bigRace")

	bigRace := false

	if bigRaceQuery != "" {
		bigRace = bigRaceQuery == "true"
	}
	return bigRace
}
