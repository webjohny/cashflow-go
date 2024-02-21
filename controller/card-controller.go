package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"github.com/webjohny/cashflow-go/session"
	"strconv"
)

type CardController interface {
	Action(ctx *gin.Context)
}

type cardController struct {
	cardService service.CardService
}

func NewCardController(cardService service.CardService) CardController {
	return &cardController{
		cardService: cardService,
	}
}

func (c *cardController) Action(ctx *gin.Context) {
	action := ctx.Param("action")
	family := ctx.Param("family")
	actionType := ctx.Param("type")
	count, _ := strconv.Atoi(ctx.Query("count"))

	raceId := session.GetItem[uint64](ctx, "raceId")
	username := session.GetItem[string](ctx, "username")

	var err error
	var response interface{}

	switch action {
	case "pre":
		err, response = c.cardService.Prepare(raceId, family, actionType, username)
		break

	case "buy":
		err, response = c.cardService.Purchase(raceId, family, actionType, username, count)
		break

	case "sell":
		err, response = c.cardService.Selling(raceId, family, actionType, username)
		break

	case "ok":
		err, response = c.cardService.Accept(raceId, family, actionType, username)
		break
	}

	request.FinalResponse(ctx, err, response)
}
