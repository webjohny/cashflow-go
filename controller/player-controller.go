package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
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
	userId := request.GetUserId(ctx)
	raceId := request.GetRaceId(ctx)

	var err error
	var response interface{}

	if userId != 0 {
		err, response = c.playerService.GetRacePlayer(raceId, userId)
	}

	request.FinalResponse(ctx, err, response)
}
