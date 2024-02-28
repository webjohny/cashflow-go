package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"github.com/webjohny/cashflow-go/session"
	"strconv"
)

type GameController interface {
	Start(ctx *gin.Context)
}

type gameController struct {
	gameService service.GameService
}

func NewGameController(gameService service.GameService) GameController {
	return &gameController{
		gameService: gameService,
	}
}

func (c *gameController) Start(ctx *gin.Context) {
	username := session.GetItem[string](ctx, "username")
	gameId, _ := strconv.Atoi(ctx.Param("gameId"))

	var err error
	var response request.Response

	err = c.gameService.Start(uint64(gameId), username)

	if err == nil {
		response = request.SuccessResponse()
	}

	request.FinalResponse(ctx, err, response)
}
