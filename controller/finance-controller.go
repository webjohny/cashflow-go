package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
)

type FinanceController interface {
	SendMoney(ctx *gin.Context)
	SendAssets(ctx *gin.Context)
	TakeLoan(ctx *gin.Context)
}

type financeController struct {
	financeService service.FinanceService
}

func NewFinanceController(financeService service.FinanceService) FinanceController {
	return &financeController{
		financeService: financeService,
	}
}

func (c *financeController) SendMoney(ctx *gin.Context) {
	raceId := request.GetRaceId(ctx)
	username := ctx.GetString("username")

	var sendMoneyBodyDTO dto.SendMoneyBodyDTO
	var err error
	var response interface{}

	errDTO := ctx.ShouldBind(&sendMoneyBodyDTO)

	if errDTO != nil {
		response = request.BuildErrorResponse("Failed to process request", errDTO.Error(), request.EmptyObj{})
	} else if raceId != 0 && username != "" {
		if sendMoneyBodyDTO.Player == "bankLoan" {
			err = c.financeService.PayLoan(raceId, username, sendMoneyBodyDTO.Amount)
		} else if sendMoneyBodyDTO.Player != "" {
			err = c.financeService.SendMoney(raceId, username, sendMoneyBodyDTO.Amount, sendMoneyBodyDTO.Player)
		}
	}

	request.FinalResponse(ctx, err, response)
}

func (c *financeController) SendAssets(ctx *gin.Context) {
	raceId := request.GetRaceId(ctx)
	username := ctx.GetString("username")

	var sendAssetsBodyDTO dto.SendAssetsBodyDTO
	var err error
	var response interface{}

	errDTO := ctx.ShouldBind(&sendAssetsBodyDTO)

	if errDTO != nil {
		response = request.BuildErrorResponse("Failed to process request", errDTO.Error(), request.EmptyObj{})
	} else if raceId != 0 && username != "" {
		err = c.financeService.SendAssets(raceId, username, sendAssetsBodyDTO)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *financeController) TakeLoan(ctx *gin.Context) {
	raceId := request.GetRaceId(ctx)
	username := ctx.GetString("username")

	var takeLoanBodyDTO dto.TakeLoanBodyDTO
	var err error
	var response interface{}

	errDTO := ctx.ShouldBind(&takeLoanBodyDTO)

	if errDTO != nil {
		response = request.BuildErrorResponse("Failed to process request", errDTO.Error(), request.EmptyObj{})
	} else if raceId != 0 && username != "" {
		err = c.financeService.TakeLoan(raceId, username, takeLoanBodyDTO.Amount)
	}

	request.FinalResponse(ctx, err, response)
}
