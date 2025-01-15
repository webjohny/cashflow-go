package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"github.com/webjohny/cashflow-go/storage"
	"gopkg.in/errgo.v2/errors"
)

type PlayerController interface {
	GetRacePlayer(ctx *gin.Context)
	MoveOnBigRace(ctx *gin.Context)
	SetDream(ctx *gin.Context)
	BecomeModerator(ctx *gin.Context)
}

type playerController struct {
	playerService service.PlayerService
	raceService   service.RaceService
	lobbyService  service.LobbyService
}

func NewPlayerController(playerService service.PlayerService, raceService service.RaceService, lobbyService service.LobbyService) PlayerController {
	return &playerController{
		playerService: playerService,
		raceService:   raceService,
		lobbyService:  lobbyService,
	}
}

func (c *playerController) GetRacePlayer(ctx *gin.Context) {
	userId := helper.GetUserId(ctx)
	raceId := helper.GetRaceId(ctx)

	var err error
	var response interface{}

	if userId != 0 {
		race := c.raceService.GetRaceByRaceId(raceId)

		if race.ID == 0 {
			err = errors.New(storage.ErrorUndefinedGame)
		} else if race.Status == entity.RaceStatus.FINISHED {
			err = errors.New(storage.ErrorGameIsFinished)
		} else {
			err, response = c.playerService.GetRacePlayer(raceId, userId)
		}
	}

	request.FinalResponse(ctx, err, response)
}

func (c *playerController) MoveOnBigRace(ctx *gin.Context) {
	userId := helper.GetUserId(ctx)
	raceId := helper.GetRaceId(ctx)

	var err error
	var response interface{}

	if userId != 0 {
		var player entity.Player

		err, player = c.playerService.GetPlayerByUserIdAndRaceId(raceId, userId)

		if err == nil {
			err = c.playerService.MoveOnBigRace(player)
		}
	}

	request.FinalResponse(ctx, err, response)
}

func (c *playerController) SetDream(ctx *gin.Context) {
	userId := helper.GetUserId(ctx)
	raceId := helper.GetRaceId(ctx)

	var err error
	var response interface{}
	var playerDream entity.PlayerDream

	errDTO := ctx.ShouldBind(&playerDream)

	if errDTO != nil {
		response = request.BuildErrorResponse("Failed to process request", errDTO.Error(), request.EmptyObj{})
	} else if userId != 0 {
		err = c.playerService.SetDream(raceId, userId, playerDream)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *playerController) BecomeModerator(ctx *gin.Context) {
	userId := helper.GetUserId(ctx)
	raceId := helper.GetRaceId(ctx)

	var err error
	var response interface{}

	err = c.lobbyService.ChangeRoleByGameIdAndUserId(raceId, userId, entity.PlayerRoles.Moderator)

	if err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	err = c.playerService.BecomeModerator(raceId, userId)

	if err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	err = c.raceService.RemovePlayer(raceId, userId)

	request.FinalResponse(ctx, err, response)
}
