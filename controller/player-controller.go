package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"github.com/webjohny/cashflow-go/storage"
	"gopkg.in/errgo.v2/errors"
	"strconv"
)

type PlayerController interface {
	GetRacePlayer(ctx *gin.Context)
	GetPlayerData(ctx *gin.Context)
	SetPlayerData(ctx *gin.Context)
	MoveOnBigRace(ctx *gin.Context)
	SetDream(ctx *gin.Context)
	BecomeModerator(ctx *gin.Context)
	IsReadNotification(ctx *gin.Context)
}

type playerController struct {
	playerService service.PlayerService
	raceService   service.RaceService
	lobbyService  service.LobbyService
}

func NewPlayerController(
	playerService service.PlayerService,
	raceService service.RaceService,
	lobbyService service.LobbyService,
) PlayerController {
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
			err, response = c.playerService.GetRacePlayer(raceId, userId, true)
		}
	}

	request.FinalResponse(ctx, err, response)
}

func (c *playerController) GetPlayerData(ctx *gin.Context) {
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
			var player entity.Player
			err, player = c.playerService.GetPlayerByUserIdAndRaceId(raceId, userId)

			if player.ID > 0 {
				player.Info.Data.Assets.Savings = strconv.Itoa(player.Assets.Savings)
				response = player.Info.Data
			}
		}
	}

	request.FinalResponse(ctx, err, response)
}

func (c *playerController) SetPlayerData(ctx *gin.Context) {
	userId := helper.GetUserId(ctx)
	raceId := helper.GetRaceId(ctx)

	var setPlayerDataDTO entity.PlayerInfoData
	var err error
	var response interface{}

	errDTO := ctx.ShouldBind(&setPlayerDataDTO)

	if errDTO != nil {
		response = request.BuildErrorResponse("Failed to process request", errDTO.Error(), request.EmptyObj{})
		request.FinalResponse(ctx, err, response)
		return
	}

	if userId != 0 {
		race := c.raceService.GetRaceByRaceId(raceId)

		if race.ID == 0 {
			err = errors.New(storage.ErrorUndefinedGame)
		} else if race.Status == entity.RaceStatus.FINISHED {
			err = errors.New(storage.ErrorGameIsFinished)
		} else {
			err = c.playerService.SetPlayerData(raceId, userId, setPlayerDataDTO)

			response = request.SuccessResponse(nil)
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

func (c *playerController) IsReadNotification(ctx *gin.Context) {
	userId := helper.GetUserId(ctx)
	raceId := helper.GetRaceId(ctx)

	notificationId := ctx.Param("notificationId")

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(raceId, userId)

	if err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	player.RemoveNotification(notificationId)

	err, _ = c.playerService.UpdatePlayer(&player)

	request.FinalResponse(ctx, err, nil)
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

	if err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	err = c.raceService.SetOptions(raceId, entity.RaceOptions{
		EnableManager: true,
	})

	if err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	request.FinalResponse(ctx, err, response)
}
