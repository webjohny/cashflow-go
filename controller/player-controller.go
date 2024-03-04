package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"strconv"
)

type PlayerController interface {
	GetRacePlayer(ctx *gin.Context)
}

type playerController struct {
	playerService service.PlayerService
}

func NewPlayerController(playerService service.PlayerService) PlayerController {
	return &playerController{
		playerService: playerService,
	}
}

func (c *playerController) GetRacePlayer(ctx *gin.Context) {
	username := ctx.GetString("username")
	raceId, _ := strconv.Atoi(ctx.Param("raceId"))

	var err error
	var response interface{}

	if username != "" {
		err, response = c.playerService.GetRacePlayer(uint64(raceId), username)
	}

	request.FinalResponse(ctx, err, response)
}
