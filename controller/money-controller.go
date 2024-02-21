package controller

import (
	"github.com/webjohny/cashflow-go/request"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/service"
)

type MoneyController interface {
	Send(ctx *gin.Context)
}

type moneyController struct {
	authService service.AuthService
}

func NewMoneyController(authService service.AuthService) MoneyController {
	return &moneyController{
		authService: authService,
	}
}

func (c *moneyController) Send(ctx *gin.Context) {
	var loginDTO dto.LoginDTO
	errDTO := ctx.ShouldBind(&loginDTO)
	if errDTO != nil {
		response := request.BuildErrorResponse("Failed to process request", errDTO.Error(), request.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	response := request.BuildErrorResponse("Please check again your credential", "Invalid Credential", request.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}
