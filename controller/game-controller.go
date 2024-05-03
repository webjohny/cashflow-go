package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"github.com/webjohny/cashflow-go/storage"
	"strconv"
)

type GameController interface {
	Start(ctx *gin.Context)
	Cancel(ctx *gin.Context)
	Reset(ctx *gin.Context)
	MoveToBigRace(ctx *gin.Context)
	RollDice(ctx *gin.Context)
	ReRollDice(ctx *gin.Context)
	GetGame(ctx *gin.Context)
	ChangeTurn(ctx *gin.Context)
	GetTiles(ctx *gin.Context)
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
	userId := request.GetUserId(ctx)
	raceId := request.GetRaceId(ctx)
	bigRace := request.GetBigRace(ctx)

	var err error
	var response interface{}

	if userId != 0 {
		response = c.gameService.GetGame(raceId, userId, bigRace)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *gameController) Reset(ctx *gin.Context) {
	userId := request.GetUserId(ctx)
	raceId := request.GetRaceId(ctx)

	var err error
	var response interface{}

	if userId != 0 {
		response = c.gameService.Reset(raceId, userId)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *gameController) Start(ctx *gin.Context) {
	lobbyId := request.GetLobbyId(ctx)

	var response dto.StartGameResponseDto

	err, race := c.gameService.Start(lobbyId)

	if err == nil {
		response = dto.StartGameResponseDto{ID: race.ID, Redirect: storage.PathShowProfession}
	}

	request.FinalResponse(ctx, err, response)
}

func (c *gameController) RollDice(ctx *gin.Context) {
	dice, _ := strconv.Atoi(ctx.Param("dice"))
	raceId := request.GetRaceId(ctx)
	userId := request.GetUserId(ctx)
	bigRace := request.GetBigRace(ctx)

	var response dto.RollDiceResponseDto

	err, diceValues := c.gameService.RollDice(raceId, userId, dice, bigRace)

	if err == nil {
		response = dto.RollDiceResponseDto{
			DiceValues: diceValues,
		}
	}

	request.FinalResponse(ctx, err, response)
}

func (c *gameController) ReRollDice(ctx *gin.Context) {
	//dice, _ := strconv.Atoi(ctx.Param("dice"))
	//raceId := request.GetRaceId(ctx)
	//userId := request.GetUserId(ctx)
	//bigRace := request.GetBigRace(ctx)
	//
	//var response dto.RollDiceResponseDto
	//
	//err, diceValues := c.gameService.RollDice(raceId, userId, dice, bigRace)
	//
	//if err == nil {
	//	response = dto.RollDiceResponseDto{
	//		DiceValues: diceValues,
	//	}
	//}

	request.FinalResponse(ctx, nil, nil)
}

func (c *gameController) Cancel(ctx *gin.Context) {
	raceId := request.GetRaceId(ctx)
	userId := request.GetUserId(ctx)

	err := c.gameService.Cancel(raceId, userId)

	request.FinalResponse(ctx, err, nil)
}

func (c *gameController) ChangeTurn(ctx *gin.Context) {
	raceId := request.GetRaceId(ctx)
	bigRace := request.GetBigRace(ctx)

	err := c.gameService.ChangeTurn(raceId, bigRace)

	request.FinalResponse(ctx, err, nil)
}

func (c *gameController) GetTiles(ctx *gin.Context) {
	raceId := request.GetRaceId(ctx)
	bigRace := request.GetBigRace(ctx)

	tiles := c.gameService.GetTiles(raceId, bigRace)

	request.FinalResponse(ctx, nil, map[string][]string{
		"tiles": tiles,
	})
}

func (c *gameController) MoveToBigRace(ctx *gin.Context) {
	//lobbyId := uint64(ctx.GetInt("lobbyId"))

	var response request.Response

	//err, _ := c.gameService.MoveToBigRace(lobbyId)
	//
	//if err == nil {
	//	response = request.RedirectResponse(storage.PathShowProfession)
	//}

	request.FinalResponse(ctx, nil, response)
}
