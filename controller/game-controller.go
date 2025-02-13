package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/objects"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"github.com/webjohny/cashflow-go/storage"
	"strconv"
	"time"
)

type GameController interface {
	Start(ctx *gin.Context)
	Cancel(ctx *gin.Context)
	Reset(ctx *gin.Context)
	RollDice(ctx *gin.Context)
	GetGame(ctx *gin.Context)
	ChangeTurn(ctx *gin.Context)
	GetTiles(ctx *gin.Context)
}

type gameController struct {
	gameService service.GameService
	mutex       *objects.MutexMap
}

func NewGameController(gameService service.GameService) GameController {
	return &gameController{
		gameService: gameService,
		mutex:       &objects.MutexMap{},
	}
}

func (c *gameController) GetGame(ctx *gin.Context) {
	userId := helper.GetUserId(ctx)
	raceId := helper.GetRaceId(ctx)

	var err error
	var response interface{}

	if userId != 0 {
		err, response = c.gameService.GetGame(raceId, userId)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *gameController) Reset(ctx *gin.Context) {
	userId := helper.GetUserId(ctx)
	raceId := helper.GetRaceId(ctx)

	var err error
	var response interface{}

	if userId != 0 {
		response = c.gameService.Reset(raceId, userId)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *gameController) Start(ctx *gin.Context) {
	lobbyId := helper.GetLobbyId(ctx)

	var response dto.StartGameResponseDto

	err, race := c.gameService.Start(lobbyId)

	if err == nil {
		response = dto.StartGameResponseDto{ID: race.ID, Redirect: storage.PathShowProfession}
	}

	request.FinalResponse(ctx, err, response)
}

func (c *gameController) RollDice(ctx *gin.Context) {
	raceId := helper.GetRaceId(ctx)

	if !c.mutex.LockMethodRace("RollDice", raceId, time.Second) {
		request.TooManyRequests(ctx)
		return
	}

	userId := helper.GetUserId(ctx)

	var response dto.RollDiceResponseDto
	var body dto.RollDiceDto

	if err := ctx.BindJSON(&body); err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	var err error

	err, response.DiceValues = c.gameService.RollDice(raceId, userId, body)

	request.FinalResponse(ctx, err, response)
}

func (c *gameController) Cancel(ctx *gin.Context) {
	raceId := helper.GetRaceId(ctx)
	userId := helper.GetUserId(ctx)

	err := c.gameService.Cancel(raceId, userId)

	request.FinalResponse(ctx, err, nil)
}

func (c *gameController) ChangeTurn(ctx *gin.Context) {
	raceId := helper.GetRaceId(ctx)

	if !c.mutex.LockMethodRace("ChangeTurn", raceId, time.Second) {
		request.TooManyRequests(ctx)
		return
	}

	forced, _ := strconv.ParseBool(ctx.Query("forced"))

	var err error

	err = c.gameService.ChangeTurn(raceId, forced)

	request.FinalResponse(ctx, err, nil)
}

func (c *gameController) GetTiles(ctx *gin.Context) {
	raceId := helper.GetRaceId(ctx)
	bigRace := helper.GetBigRace(ctx)

	tiles := c.gameService.GetTiles(raceId, bigRace)

	request.FinalResponse(ctx, nil, map[string][]string{
		"tiles": tiles,
	})
}
