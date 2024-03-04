package controller

import (
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/request"
	"net/http"
	"os"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/service"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
	authService service.AuthService
	jwtService  service.JWTService
}

// NewAuthController is for blabla
func NewAuthController(authService service.AuthService, jwtService service.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) Login(ctx *gin.Context) {
	var loginDTO dto.LoginDTO
	errDTO := ctx.ShouldBind(&loginDTO)
	if errDTO != nil {
		response := request.BuildErrorResponse("Failed to process request", errDTO.Error(), request.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	authResult := c.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)
	if v, ok := authResult.(entity.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(v.ID, 10), v.Email, v.Profile, v.Jk, v.Name)
		v.Token = generatedToken
		response := request.BuildResponse(true, "OK", v)
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := request.BuildErrorResponse("Please check again your credential", "Invalid Credential", request.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

func (c *authController) Register(ctx *gin.Context) {
	var registerDTO dto.RegisterDTO
	errDTO := ctx.ShouldBind(&registerDTO)
	if errDTO != nil {
		response := request.BuildErrorResponse("Failed to process request", errDTO.Error(), request.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !c.authService.IsDuplicateEmail(registerDTO.Email) {
		response := request.BuildErrorResponse("Failed to process request", "Duplicate email", request.EmptyObj{})
		ctx.JSON(http.StatusConflict, response)
	} else {
		registerDTO.Jk = helper.CreateHash(registerDTO.Email + strconv.Itoa(int(registerDTO.ID)) + os.Getenv("SECRET"))
		createdUser := c.authService.CreateUser(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10), createdUser.Email, createdUser.Profile, createdUser.Jk, createdUser.Name)
		createdUser.Token = token
		response := request.BuildResponse(true, "OK", createdUser)
		ctx.JSON(http.StatusCreated, response)
	}
}
