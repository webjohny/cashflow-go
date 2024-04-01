package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
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

	if raceId != 0 && userId != 0 {
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

	if raceId != 0 && userId != 0 {
		err, response = c.cardService.Prepare(raceId, family, actionType, userId, bigRace)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *cardController) Selling(ctx *gin.Context) {
	actionType := ctx.Param("type")

	raceId := request.GetRaceId(ctx)
	userId := request.GetUserId(ctx)
	bigRace := request.GetBigRace(ctx)
	value := ctx.Query("value")

	var err error
	var response interface{}

	if raceId != 0 && userId != 0 {
		err, response = c.cardService.Selling(raceId, actionType, userId, value, bigRace)
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

	if raceId != 0 && userId != 0 {
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

	if raceId != 0 && userId != 0 {
		err, response = c.cardService.Skip(raceId, userId, bigRace)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *cardController) Purchase(ctx *gin.Context) {
	actionType := ctx.Param("type")
	count, _ := strconv.Atoi(ctx.Query("count"))

	raceId := request.GetRaceId(ctx)
	userId := request.GetUserId(ctx)
	bigRace := request.GetBigRace(ctx)

	var err error
	var response interface{}

	if raceId != 0 && userId != 0 {
		err, response = c.cardService.Purchase(raceId, actionType, userId, count, bigRace)
	}

	request.FinalResponse(ctx, err, response)
}
