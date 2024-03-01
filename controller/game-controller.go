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
	GetGame(ctx *gin.Context)
}

type gameController struct {
	gameService service.GameService
}

func NewGameController(gameService service.GameService) GameController {
	return &gameController{
		gameService: gameService,
	}
}

func (c *gameController) GetGame(ctx *gin.Context) {
	username := session.GetItem[string](ctx, "username")
	raceId, _ := strconv.Atoi(ctx.Param("raceId"))
	lobbyId, _ := strconv.Atoi(ctx.Param("lobbyId"))
	bigRaceQuery := ctx.Query("bigRace")

	var bigRace *bool

	if bigRaceQuery != "" {
		*bigRace = bigRaceQuery == "true"
	}

	err, game := c.gameService.GetGame(uint64(raceId), uint64(lobbyId), username, bigRace)

	request.FinalResponse(ctx, err, game)
}

func (c *gameController) Start(ctx *gin.Context) {
	lobbyId, _ := strconv.Atoi(ctx.Param("lobbyId"))

	var response request.Response

	err, race := c.gameService.Start(uint64(lobbyId))

	if err == nil {
		session.DeleteItem(ctx, "lobbyId")
		session.SetItem(ctx, "raceId", race.ID)
		response = request.SuccessResponse()
	}

	request.FinalResponse(ctx, err, response)
}
