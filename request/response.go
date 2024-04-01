package request

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

type TransactionGroupSum struct {
	TransactionGroup string `json:"transaction_group"`
	TotalTransaction int    `json:"total_transaction"`
}

type TransactionReport struct {
	TransactionOut int `json:"transaction_out"`
	TransactionIn  int `json:"total_in"`
}

type EmptyObj struct{}

func BuildResponse(status bool, message string, data interface{}) Response {
	res := Response{
		Status:  status,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
	return res
}

func BuildErrorResponse(message string, err string, data interface{}) Response {
	splittedError := strings.Split(err, "\n")
	res := Response{
		Status:  false,
		Message: message,
		Errors:  splittedError,
		Data:    data,
	}
	return res
}

func FinalResponse(ctx *gin.Context, err error, response interface{}) {
	if err != nil {
		response := BuildErrorResponse("Failed to process request", err.Error(), EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(response))
}
