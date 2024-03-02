package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"github.com/webjohny/cashflow-go/session"
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
	username := session.GetItem[string](ctx, "username")
	raceId, _ := strconv.Atoi(ctx.Param("raceId"))

	var err error
	var response interface{}

	if username != nil {
		err, response = c.playerService.GetRacePlayer(uint64(raceId), *username)
	}

	request.FinalResponse(ctx, err, response)
}
