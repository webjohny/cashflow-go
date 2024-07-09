package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"github.com/webjohny/cashflow-go/storage"
)

type CardController interface {
	Prepare(ctx *gin.Context)
	Purchase(ctx *gin.Context)
	Selling(ctx *gin.Context)
	Accept(ctx *gin.Context)
	Skip(ctx *gin.Context)
	Type(ctx *gin.Context)
	TestCard(ctx *gin.Context)
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

	var err error
	var response interface{}

	if userId == 0 {
		err = errors.New(storage.ErrorUndefinedPlayer)
	} else if raceId == 0 {
		err = errors.New(storage.ErrorUndefinedGame)
	} else {
		err, response = c.cardService.GetCard("", raceId, userId, bigRace)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *cardController) TestCard(ctx *gin.Context) {
	raceId := request.GetRaceId(ctx)
	userId := request.GetUserId(ctx)
	bigRace := request.GetBigRace(ctx)

	var err error
	var response interface{}

	if userId == 0 {
		err = errors.New(storage.ErrorUndefinedPlayer)
	} else if raceId == 0 {
		err = errors.New(storage.ErrorUndefinedGame)
	} else {
		err, response = c.cardService.TestCard(ctx.Param("action"), raceId, userId, bigRace)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *cardController) ResetTransaction(ctx *gin.Context) {
}

func (c *cardController) Prepare(ctx *gin.Context) {
	family := ctx.Param("family")
	actionType := ctx.Param("type")

	raceId := request.GetRaceId(ctx)
	userId := request.GetUserId(ctx)
	bigRace := request.GetBigRace(ctx)

	var err error
	var response interface{}

	if raceId == 0 {
		err = errors.New(storage.ErrorUndefinedGame)
	} else if userId == 0 {
		err = errors.New(storage.ErrorUndefinedPlayer)
	} else {
		err, response = c.cardService.Prepare(actionType, raceId, family, userId, bigRace)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *cardController) Selling(ctx *gin.Context) {
	actionType := ctx.Param("type")

	raceId := request.GetRaceId(ctx)
	userId := request.GetUserId(ctx)
	bigRace := request.GetBigRace(ctx)

	var body dto.CardSellingActionDTO

	if err := ctx.BindJSON(&body); err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	var err error
	var response interface{}

	if raceId == 0 {
		err = errors.New(storage.ErrorUndefinedGame)
	} else if userId == 0 {
		err = errors.New(storage.ErrorUndefinedUser)
	} else {
		err, response = c.cardService.Selling(actionType, raceId, userId, bigRace, body)
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
		err = errors.New(storage.ErrorUndefinedGame)
	} else if userId == 0 {
		err = errors.New(storage.ErrorUndefinedUser)
	} else {
		err, response = c.cardService.Accept(actionType, raceId, family, userId, bigRace)
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
		err = errors.New(storage.ErrorUndefinedGame)
	} else if userId == 0 {
		err = errors.New(storage.ErrorUndefinedUser)
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

	var body dto.CardPurchaseActionDTO

	if err := ctx.BindJSON(&body); err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	var err error
	var response interface{}

	if raceId == 0 {
		err = errors.New(storage.ErrorUndefinedGame)
	} else if userId == 0 {
		err = errors.New(storage.ErrorUndefinedUser)
	} else {
		err, response = c.cardService.Purchase(actionType, raceId, userId, bigRace, body)
	}

	request.FinalResponse(ctx, err, response)
}
