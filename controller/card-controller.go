package controller

import (
	"github.com/webjohny/cashflow-go/session"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/service"
)

type CardController interface {
	Action(ctx *gin.Context) string
}

type cardController struct {
	cardService service.CardService
}

func NewCardController(cardService service.CardService) CardController {
	return &cardController{
		cardService: cardService,
	}
}

func (c *cardController) Action(ctx *gin.Context) string {
	action := ctx.Param("action")
	family := ctx.Param("family")
	actionType := ctx.Param("type")

	raceId := session.GetItem[uint64](ctx, "raceId")
	username := session.GetItem[string](ctx, "username")

	var resp string

	switch action {
	case "pre":
		resp = c.cardService.Prepare(raceId, family, actionType, username)
		break

	case "buy":
		resp = c.cardService.Purchase(raceId, family, actionType, username)
		break

	case "sell":
		resp = c.cardService.Selling(raceId, family, actionType, username)
		break

	case "ok":
		resp = c.cardService.Accept(raceId, family, actionType, username)
		break
	}

	log.Println(resp, action, family, actionType)

	var actionDTO dto.CardActionDTO
	errDTO := ctx.ShouldBind(&actionDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return ""
	}
	response := helper.BuildErrorResponse("Please check again your credential", "Invalid Credential", helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)

	return resp
}
