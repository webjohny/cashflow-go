package controller

import (
	"github.com/webjohny/cashflow-go/request"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/service"
)

type GameController interface {
	Start(ctx *gin.Context)
}

type gameController struct {
	authService service.AuthService
	jwtService  service.JWTService
}

func NewGameController(authService service.AuthService, jwtService service.JWTService) GameController {
	return &gameController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *gameController) Start(ctx *gin.Context) {
	var loginDTO dto.LoginDTO
	errDTO := ctx.ShouldBind(&loginDTO)
	if errDTO != nil {
		response := request.BuildErrorResponse("Failed to process request", errDTO.Error(), request.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	authResult := c.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)
	if v, ok := authResult.(entity.User); ok {
		generatedTokn := c.jwtService.GenerateToken(strconv.FormatUint(v.ID, 10), v.Email, v.Profile, v.Jk, v.Telephone, v.Pin, v.Name)
		v.Token = generatedTokn
		response := request.BuildResponse(true, "OK", v)
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := request.BuildErrorResponse("Please check again your credential", "Invalid Credential", request.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}
