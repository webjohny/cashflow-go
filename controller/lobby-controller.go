package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"github.com/webjohny/cashflow-go/session"
)

type LobbyController interface {
	CreateLobby(ctx *gin.Context)
	Join(ctx *gin.Context)
}

type lobbyController struct {
	lobbyService service.LobbyService
}

func NewLobbyController(lobbyService service.LobbyService) LobbyController {
	return &lobbyController{
		lobbyService: lobbyService,
	}
}

func (c *lobbyController) CreateLobby(ctx *gin.Context) {
	username := session.GetItem[string](ctx, "username")

	err, lobby := c.lobbyService.CreateLobby(username)

	session.SetItem(ctx, "gameId", lobby.ID)

	request.FinalResponse(ctx, err, lobby)
}

func (c *lobbyController) Join(ctx *gin.Context) {
	username := session.GetItem[string](ctx, "username")
	gameId := session.GetItem[uint64](ctx, "gameId")

	err := c.lobbyService.Join(gameId, username)

	request.FinalResponse(ctx, err, nil)
}
