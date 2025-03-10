package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/objects"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"github.com/webjohny/cashflow-go/storage"
	"time"
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
	SetCards(ctx *gin.Context)
}

type cardController struct {
	cardService service.CardService
	mutex       *objects.MutexMap
}

func NewCardController(cardService service.CardService) CardController {
	return &cardController{
		cardService: cardService,
		mutex:       &objects.MutexMap{},
	}
}

func (c *cardController) Type(ctx *gin.Context) {
	raceId := helper.GetRaceId(ctx)

	if !c.mutex.LockMethodRace("CardType", raceId, time.Second) {
		request.TooManyRequests(ctx)
		return
	}

	userId := helper.GetUserId(ctx)
	cardType := ctx.Query("type")

	var err error
	var response interface{}

	if userId == 0 {
		err = errors.New(storage.ErrorUndefinedPlayer)
	} else if raceId == 0 {
		err = errors.New(storage.ErrorUndefinedGame)
	} else {
		err, response = c.cardService.GetCard("", raceId, userId, cardType)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *cardController) TestCard(ctx *gin.Context) {
	raceId := helper.GetRaceId(ctx)
	userId := helper.GetUserId(ctx)
	bigRace := helper.GetBigRace(ctx)

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

func (c *cardController) SetCards(ctx *gin.Context) {
	raceId := helper.GetRaceId(ctx)

	var err error
	var response interface{}

	c.mutex.UnlockMethodRace("CardType", raceId)

	var body dto.CreateCardsDTO

	if err = ctx.BindJSON(&body); err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	c.cardService.SetCards(body)

	request.FinalResponse(ctx, err, response)
}

func (c *cardController) Prepare(ctx *gin.Context) {
	family := ctx.Param("family")
	actionType := ctx.Param("type")

	raceId := helper.GetRaceId(ctx)
	userId := helper.GetUserId(ctx)

	var err error
	var response interface{}
	var body dto.PrepareCardBodyDTO

	if err = ctx.BindJSON(&body); err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	if raceId == 0 {
		err = errors.New(storage.ErrorUndefinedGame)
	} else if userId == 0 {
		err = errors.New(storage.ErrorUndefinedPlayer)
	} else {
		err, response = c.cardService.Prepare(actionType, raceId, family, userId, body)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *cardController) Selling(ctx *gin.Context) {
	actionType := ctx.Param("type")

	raceId := helper.GetRaceId(ctx)
	userId := helper.GetUserId(ctx)
	bigRace := helper.GetBigRace(ctx)

	if !c.mutex.LockMethodRace("Selling", raceId, time.Second*3) {
		request.TooManyRequests(ctx)
		return
	}

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

	raceId := helper.GetRaceId(ctx)
	userId := helper.GetUserId(ctx)
	bigRace := helper.GetBigRace(ctx)

	if !c.mutex.LockMethodRace("Accept", raceId, time.Second*3) {
		request.TooManyRequests(ctx)
		return
	}

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
	raceId := helper.GetRaceId(ctx)
	userId := helper.GetUserId(ctx)
	bigRace := helper.GetBigRace(ctx)

	if !c.mutex.LockMethodRace("Skip", raceId, time.Second*3) {
		request.TooManyRequests(ctx)
		return
	}

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

	raceId := helper.GetRaceId(ctx)
	userId := helper.GetUserId(ctx)
	bigRace := helper.GetBigRace(ctx)

	if !c.mutex.LockMethodRace("Purchase", raceId, time.Second*3) {
		request.TooManyRequests(ctx)
		return
	}

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
