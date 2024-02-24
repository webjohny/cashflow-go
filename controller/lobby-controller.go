package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/service"
)

type LobbyController interface {
	CreateLobby(ctx *gin.Context)
}

type lobbyController struct {
	gameService service.GameService
	raceService service.RaceService
}

func NewLobbyController(gameService service.GameService, raceService service.RaceService) LobbyController {
	return &lobbyController{
		gameService: gameService,
		raceService: raceService,
	}
}

func (c *lobbyController) CreateLobby(ctx *gin.Context) {
	//const dice = new Dice(diceValues?.[0] || 1, 2, 6);
	//
	//const lobby = await req.games.newLobby(username, dice);
	//req.session.gameID = lobby.state.gameID;
	//req.game = lobby;
}
