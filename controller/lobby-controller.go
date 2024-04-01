package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
)

type LobbyController interface {
	Create(ctx *gin.Context)
	Join(ctx *gin.Context)
	Leave(ctx *gin.Context)
}

type lobbyController struct {
	lobbyService service.LobbyService
}

func NewLobbyController(lobbyService service.LobbyService) LobbyController {
	return &lobbyController{
		lobbyService: lobbyService,
	}
}

func (c *lobbyController) Create(ctx *gin.Context) {
	userId := request.GetUserId(ctx)
	username := ctx.GetString("username")

	var err error
	var lobby entity.Lobby

	if userId != 0 {
		err, lobby = c.lobbyService.Create(username, userId)
	}

	request.FinalResponse(ctx, err, lobby)
}

func (c *lobbyController) Join(ctx *gin.Context) {
	userId := request.GetUserId(ctx)
	username := ctx.GetString("username")
	lobbyId := request.GetLobbyId(ctx)

	var err error
	var player entity.LobbyPlayer

	if userId != 0 {
		err, player = c.lobbyService.Join(lobbyId, username, userId)
	}

	request.FinalResponse(ctx, err, player)
}

func (c *lobbyController) Leave(ctx *gin.Context) {
	username := ctx.GetString("username")
	lobbyId := request.GetLobbyId(ctx)

	var err error

	if username != "" {
		err, _ = c.lobbyService.Leave(lobbyId, username)
	}

	request.FinalResponse(ctx, err, nil)
}
