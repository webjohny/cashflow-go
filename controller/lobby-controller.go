package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"github.com/webjohny/cashflow-go/storage"
)

type LobbyController interface {
	Create(ctx *gin.Context)
	Join(ctx *gin.Context)
	Leave(ctx *gin.Context)
	Cancel(ctx *gin.Context)
	SetOptions(ctx *gin.Context)
	GetLobby(ctx *gin.Context)
}

type lobbyController struct {
	lobbyService service.LobbyService
}

func NewLobbyController(lobbyService service.LobbyService) LobbyController {
	return &lobbyController{
		lobbyService: lobbyService,
	}
}

func (c *lobbyController) GetLobby(ctx *gin.Context) {
	userId := helper.GetUserId(ctx)
	lobbyId := helper.GetLobbyId(ctx)

	var err error
	var response interface{}

	if userId != 0 {
		err, response = c.lobbyService.GetLobby(lobbyId, userId)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *lobbyController) Create(ctx *gin.Context) {
	userId := helper.GetUserId(ctx)
	username := ctx.GetString("name")

	var err error
	var lobby entity.Lobby

	if userId != 0 {
		err, lobby = c.lobbyService.Create(username, userId)
	}

	request.FinalResponse(ctx, err, lobby)
}

func (c *lobbyController) SetOptions(ctx *gin.Context) {
	userId := helper.GetUserId(ctx)
	lobbyId := helper.GetLobbyId(ctx)

	var err error
	var body dto.SetOptionsLobbyRequestDTO

	if err = ctx.BindJSON(&body); err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	if userId != 0 {
		err = c.lobbyService.SetOptions(lobbyId, body)
	}

	request.FinalResponse(ctx, err, nil)
}

func (c *lobbyController) Join(ctx *gin.Context) {
	userId := helper.GetUserId(ctx)
	username := ctx.GetString("name")
	lobbyId := helper.GetLobbyId(ctx)

	var err error
	var player entity.LobbyPlayer

	if userId != 0 {
		err, player = c.lobbyService.Join(lobbyId, username, userId)
	}

	request.FinalResponse(ctx, err, player)
}

func (c *lobbyController) Leave(ctx *gin.Context) {
	userId := helper.GetUserId(ctx)
	lobbyId := helper.GetLobbyId(ctx)

	var err error

	if userId > 0 {
		err, _ = c.lobbyService.Leave(lobbyId, userId)
	} else {
		err = errors.New(storage.ErrorUndefinedUser)
	}

	request.FinalResponse(ctx, err, nil)
}

func (c *lobbyController) Cancel(ctx *gin.Context) {
	userId := helper.GetUserId(ctx)
	lobbyId := helper.GetLobbyId(ctx)

	var err error

	if userId != 0 {
		err, _ = c.lobbyService.Cancel(lobbyId, userId)
	} else {
		err = errors.New(storage.ErrorUndefinedUser)
	}

	request.FinalResponse(ctx, err, nil)
}
