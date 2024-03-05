package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"log"
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
	username := ctx.GetString("username")
	raceId := uint64(ctx.GetInt("raceId"))
	lobbyId := uint64(ctx.GetInt("lobbyId"))
	bigRaceQuery := ctx.Query("bigRace")

	var bigRace *bool

	if bigRaceQuery != "" {
		*bigRace = bigRaceQuery == "true"
	}

	var err error
	var response interface{}

	if username != "" {
		err, response = c.gameService.GetGame(raceId, lobbyId, username, bigRace)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *gameController) Start(ctx *gin.Context) {
	lobbyId := uint64(ctx.GetInt("lobbyId"))

	var response request.Response

	err, race := c.gameService.Start(lobbyId)

	if err == nil {
		log.Println(race)
		response = request.SuccessResponse()
	}

	request.FinalResponse(ctx, err, response)
}
