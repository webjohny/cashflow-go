package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"github.com/webjohny/cashflow-go/storage"
	"log"
	"strconv"
)

type CardController interface {
	Prepare(ctx *gin.Context)
	Purchase(ctx *gin.Context)
	Selling(ctx *gin.Context)
	Accept(ctx *gin.Context)
	Skip(ctx *gin.Context)
	Type(ctx *gin.Context)
	ResetTransaction(ctx *gin.Context)
}

type cardController struct {
	cardService service.CardService
}

func NewCardController(cardService service.CardService) CardController {
	return &cardController{
		cardService: cardService,
	}
}

func (c *cardController) Type(ctx *gin.Context) {
	raceId := request.GetRaceId(ctx)
	userId := request.GetUserId(ctx)
	bigRace := request.GetBigRace(ctx)

	log.Println(raceId)

	var err error
	var response interface{}

	if userId == 0 {
		err = fmt.Errorf(storage.ErrorUndefinedPlayer)
	} else if raceId == 0 {
		err = fmt.Errorf(storage.ErrorUndefinedGame)
	} else {
		err, response = c.cardService.GetCard("", raceId, userId, bigRace)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *cardController) ResetTransaction(ctx *gin.Context) {}

func (c *cardController) Prepare(ctx *gin.Context) {
	family := ctx.Param("family")
	actionType := ctx.Param("type")

	raceId := request.GetRaceId(ctx)
	userId := request.GetUserId(ctx)
	bigRace := request.GetBigRace(ctx)

	var err error
	var response interface{}

	if raceId == 0 {
		err = fmt.Errorf(storage.ErrorUndefinedGame)
	} else if userId == 0 {
		err = fmt.Errorf(storage.ErrorUndefinedPlayer)
	} else {
		err, response = c.cardService.Prepare(raceId, family, actionType, userId, bigRace)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *cardController) Selling(ctx *gin.Context) {
	actionType := ctx.Param("type")

	raceId := request.GetRaceId(ctx)
	userId := request.GetUserId(ctx)
	bigRace := request.GetBigRace(ctx)

	var body dto.CardActionDTO

	if err := ctx.BindJSON(&body); err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	var err error
	var response interface{}

	if raceId == 0 {
		err = fmt.Errorf(storage.ErrorUndefinedGame)
	} else if userId == 0 {
		err = fmt.Errorf(storage.ErrorUndefinedUser)
	} else {
		err, response = c.cardService.Selling(raceId, actionType, userId, body.Value, bigRace)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *cardController) Accept(ctx *gin.Context) {
	family := ctx.Param("family")
	actionType := ctx.Param("type")

	raceId := request.GetRaceId(ctx)
	userId := request.GetUserId(ctx)
	bigRace := request.GetBigRace(ctx)

	var err error
	var response interface{}

	if raceId == 0 {
		err = fmt.Errorf(storage.ErrorUndefinedGame)
	} else if userId == 0 {
		err = fmt.Errorf(storage.ErrorUndefinedUser)
	} else {
		err, response = c.cardService.Accept(raceId, family, actionType, userId, bigRace)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *cardController) Skip(ctx *gin.Context) {
	raceId := request.GetRaceId(ctx)
	userId := request.GetUserId(ctx)
	bigRace := request.GetBigRace(ctx)

	var err error
	var response interface{}

	if raceId == 0 {
		err = fmt.Errorf(storage.ErrorUndefinedGame)
	} else if userId == 0 {
		err = fmt.Errorf(storage.ErrorUndefinedUser)
	} else {
		err, response = c.cardService.Skip(raceId, userId, bigRace)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *cardController) Purchase(ctx *gin.Context) {
	actionType := ctx.Param("type")

	raceId := request.GetRaceId(ctx)
	userId := request.GetUserId(ctx)
	bigRace := request.GetBigRace(ctx)

	var body dto.CardActionDTO

	if err := ctx.BindJSON(&body); err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	value, _ := strconv.Atoi(body.Value)

	var err error
	var response interface{}

	if raceId == 0 {
		err = fmt.Errorf(storage.ErrorUndefinedGame)
	} else if userId == 0 {
		err = fmt.Errorf(storage.ErrorUndefinedUser)
	} else {
		err, response = c.cardService.Purchase(raceId, actionType, userId, value, bigRace)
	}

	request.FinalResponse(ctx, err, response)
}
