package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"github.com/webjohny/cashflow-go/session"
	"strconv"
)

type LobbyController interface {
	CreateLobby(ctx *gin.Context)
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

func (c *lobbyController) CreateLobby(ctx *gin.Context) {
	username := session.GetItem[string](ctx, "username")

	err, lobby := c.lobbyService.CreateLobby(username)

	session.SetItem(ctx, "gameId", lobby.ID)

	request.FinalResponse(ctx, err, lobby)
}

func (c *lobbyController) Join(ctx *gin.Context) {
	username := session.GetItem[string](ctx, "username")
	gameId, _ := strconv.Atoi(ctx.Param("gameId"))

	var err error
	var response request.Response

	err = c.lobbyService.Join(uint64(gameId), username)

	if err == nil {
		session.SetItem(ctx, "gameId", uint64(gameId))
		response = request.SuccessResponse()
	}

	request.FinalResponse(ctx, err, response)
}

func (c *lobbyController) Leave(ctx *gin.Context) {
	username := session.GetItem[string](ctx, "username")
	gameId, _ := strconv.Atoi(ctx.Param("gameId"))

	var err error
	var response request.Response

	err = c.lobbyService.Leave(uint64(gameId), username)

	if err == nil {
		session.DeleteItem(ctx, "gameId")
		response = request.SuccessResponse()
	}

	request.FinalResponse(ctx, err, response)
}
