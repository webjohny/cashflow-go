package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"time"
)

type BackdoorController interface {
	ChangeCard(ctx *gin.Context)
	RollDice(ctx *gin.Context)
}

type backdoorController struct {
	cardService   service.CardService
	raceService   service.RaceService
	playerService service.PlayerService
	gameService   service.GameService
}

func NewBackdoorController(cardService service.CardService, raceService service.RaceService, playerService service.PlayerService, gameService service.GameService) BackdoorController {
	return &backdoorController{
		cardService:   cardService,
		raceService:   raceService,
		playerService: playerService,
		gameService:   gameService,
	}
}

func (c *backdoorController) ChangeCard(ctx *gin.Context) {
	raceId := helper.GetRaceId(ctx)

	var err error
	var response interface{}

	var body dto.BackdoorCardBodyDTO

	if err = ctx.BindJSON(&body); err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	race := c.raceService.GetRaceByRaceId(raceId)
	race.CurrentCard = body.Card

	err = c.cardService.ProcessCard(race)

	if race.CurrentCard.Family == "market" || race.CurrentCard.Type == "stock" {
		race.IsMultiFlow = race.CurrentCard.OnlyYou == false
	} else {
		race.IsMultiFlow = false
	}

	err, response = c.raceService.UpdateRace(&race)

	request.FinalResponse(ctx, err, response)
}

func (c *backdoorController) RollDice(ctx *gin.Context) {
	raceId := helper.GetRaceId(ctx)

	var err error
	var response interface{}

	var body dto.BackdoorRollDiceBodyDto

	if err = ctx.BindJSON(&body); err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	race := c.raceService.GetRaceByRaceId(raceId)

	diceValues := make([]int, 0)

	for k, v := range body.Dices {
		rollDiceDto := dto.RollDiceDto{
			IsFinished: false,
			Dices:      0,
			DiceValue:  v,
		}

		if k == len(body.Dices)-1 {
			rollDiceDto.IsFinished = true
		}
		var diceValue []int
		err, diceValue = c.gameService.RollDice(raceId, race.CurrentPlayer.UserId, rollDiceDto)

		if err == nil {
			break
		}

		diceValues = append(diceValues, diceValue...)

		time.Sleep(time.Second)
	}

	if err == nil {
		err, response = c.raceService.UpdateRace(&race)
	}

	request.FinalResponse(ctx, err, response)
}
