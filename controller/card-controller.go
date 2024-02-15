package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/service"
)

type CardController interface {
	Action(ctx *gin.Context)
}

type cardController struct {
	cardService service.CardService
	gameService service.GameService
}

func NewCardController(cardService service.CardService, gameService service.GameService) CardController {
	return &cardController{
		cardService: cardService,
		gameService: gameService,
	}
}

func (c *cardController) Action(ctx *gin.Context) {
	action := ctx.Param("action")
	family := ctx.Param("family")
	actionType := ctx.Param("type")

	race := c.gameService.GetRaceByCtx(ctx)

	//username := session.GetItem[string](ctx, "username")
	//game := ctx.MustGet("game")

	log.Println(action, family, actionType)

	var actionDTO dto.CardActionDTO
	errDTO := ctx.ShouldBind(&actionDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	response := helper.BuildErrorResponse("Please check again your credential", "Invalid Credential", helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}
