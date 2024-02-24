package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"github.com/webjohny/cashflow-go/session"
	"strconv"
)

type CardController interface {
	Prepare(ctx *gin.Context)
	Purchase(ctx *gin.Context)
	Selling(ctx *gin.Context)
	Accept(ctx *gin.Context)
	Skip(ctx *gin.Context)
}

type cardController struct {
	cardService service.CardService
}

func NewCardController(cardService service.CardService) CardController {
	return &cardController{
		cardService: cardService,
	}
}

func (c *cardController) Prepare(ctx *gin.Context) {
	family := ctx.Param("family")
	actionType := ctx.Param("type")

	raceId := session.GetItem[uint64](ctx, "raceId")
	username := session.GetItem[string](ctx, "username")

	err, response := c.cardService.Prepare(raceId, family, actionType, username)

	request.FinalResponse(ctx, err, response)
}

func (c *cardController) Selling(ctx *gin.Context) {
	actionType := ctx.Param("type")

	raceId := session.GetItem[uint64](ctx, "raceId")
	username := session.GetItem[string](ctx, "username")

	err, response := c.cardService.Selling(raceId, actionType, username)

	request.FinalResponse(ctx, err, response)
}

func (c *cardController) Accept(ctx *gin.Context) {
	family := ctx.Param("family")
	actionType := ctx.Param("type")

	raceId := session.GetItem[uint64](ctx, "raceId")
	username := session.GetItem[string](ctx, "username")

	err, response := c.cardService.Accept(raceId, family, actionType, username)

	request.FinalResponse(ctx, err, response)
}

func (c *cardController) Skip(ctx *gin.Context) {
	raceId := session.GetItem[uint64](ctx, "raceId")
	username := session.GetItem[string](ctx, "username")

	err, response := c.cardService.Skip(raceId, username)

	request.FinalResponse(ctx, err, response)
}

func (c *cardController) Purchase(ctx *gin.Context) {
	actionType := ctx.Param("type")
	count, _ := strconv.Atoi(ctx.Query("count"))

	raceId := session.GetItem[uint64](ctx, "raceId")
	username := session.GetItem[string](ctx, "username")

	err, response := c.cardService.Purchase(raceId, actionType, username, count)

	request.FinalResponse(ctx, err, response)
}
