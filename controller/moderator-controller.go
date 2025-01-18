package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"strconv"
)

type ModeratorController interface {
	GetRace(ctx *gin.Context)
	GetRacePlayer(ctx *gin.Context)
	GetRacePlayers(ctx *gin.Context)
	UpdatePlayer(ctx *gin.Context)
	UpdateRace(ctx *gin.Context)
	HandleUserRequest(ctx *gin.Context)
}

type moderatorController struct {
	playerService      service.PlayerService
	raceService        service.RaceService
	lobbyService       service.LobbyService
	userRequestService service.UserRequestService
}

func NewModeratorController(
	playerService service.PlayerService,
	raceService service.RaceService,
	lobbyService service.LobbyService,
	userRequestService service.UserRequestService,
) ModeratorController {
	return &moderatorController{
		playerService:      playerService,
		raceService:        raceService,
		lobbyService:       lobbyService,
		userRequestService: userRequestService,
	}
}

func (c *moderatorController) GetRace(ctx *gin.Context) {
	raceId := helper.GetRaceId(ctx)

	var err error
	var response dto.GetRaceResponseDTO

	if raceId != 0 {
		response = c.raceService.GetFormattedRaceResponse(raceId, true)

		if response.GameId > 0 {
			response.UserRequests = c.userRequestService.GetAllByRaceId(raceId)
		}
	}

	request.FinalResponse(ctx, err, response)
}

func (c *moderatorController) GetRacePlayer(ctx *gin.Context) {
	userIdParam, _ := strconv.Atoi(ctx.Query("playerId"))
	userId := uint64(userIdParam)

	raceId := helper.GetRaceId(ctx)

	var err error
	var response interface{}

	if userId != 0 {
		err, response = c.playerService.GetRacePlayer(raceId, userId)
	}

	request.FinalResponse(ctx, err, response)
}

func (c *moderatorController) GetRacePlayers(ctx *gin.Context) {
	raceId := helper.GetRaceId(ctx)

	var err error

	players := c.raceService.GetRacePlayersByRaceId(raceId)

	request.FinalResponse(ctx, err, map[string]interface{}{
		"players": players,
	})
}

func (c *moderatorController) UpdatePlayer(ctx *gin.Context) {
	raceId := helper.GetRaceId(ctx)
	playerId := helper.ConvertToUInt64(ctx.Param("playerId"))

	var err error
	var body dto.ModeratorUpdatePlayerDto

	if err = ctx.BindJSON(&body); err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	err, player := c.playerService.GetPlayerByPlayerIdAndRaceId(raceId, playerId)

	if err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	player.Cash = body.Cash
	player.CashFlow = body.CashFlow
	player.Babies = uint8(body.Babies)
	player.CurrentPosition = uint8(body.CurrentPosition)
	player.LastPosition = uint8(body.LastPosition)
	player.SkippedTurns = uint8(body.SkippedTurns)
	player.OnBigRace = body.OnBigRace
	player.Assets.Savings = body.Savings

	for _, realEstate := range player.Assets.RealEstates {
		if item, ok := body.RealEstate[realEstate.ID]; ok {
			player.CreateOrUpdateRealEstateByID(item)
		} else {
			player.RemoveRealEstate(realEstate.ID)
		}
	}

	for _, business := range player.Assets.Business {
		if item, ok := body.Business[business.ID]; ok {
			player.CreateOrUpdateBusinessByID(item)
		} else {
			player.RemoveBusiness(business.ID)
		}
	}

	for _, other := range player.Assets.OtherAssets {
		if item, ok := body.Other[other.ID]; ok {
			player.CreateOrUpdateOtherAssetByID(item)
		} else {
			player.RemoveOtherAssetsByID(other.ID)
		}
	}

	for _, stock := range player.Assets.Stocks {
		if item, ok := body.Stocks[stock.ID]; ok {
			player.CreateOrUpdateStocksByID(item)
		} else {
			player.RemoveStocks(stock.Symbol)
		}
	}

	err, player = c.playerService.UpdatePlayer(&player)

	request.FinalResponse(ctx, err, map[string]interface{}{
		"player": player,
	})
}

func (c *moderatorController) UpdateRace(ctx *gin.Context) {
	raceId := helper.GetRaceId(ctx)

	var err error

	var body dto.ModeratorUpdateRaceDto

	if err = ctx.BindJSON(&body); err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	race := c.raceService.GetRaceByRaceId(raceId)

	race.Options.EnableManager = body.EnableManager
	race.Options.HideCards = body.HideCards
	race.Status = body.Status

	if race.Status == entity.RaceStatus.FINISHED || race.Status == entity.RaceStatus.CANCELLED {
		_ = c.lobbyService.ChangeStatusByGameId(raceId, entity.LobbyStatus.Cancelled)
	}

	for k, raceResponse := range body.Responses {
		race.Responses[k].Responded = raceResponse
	}

	err, _ = c.raceService.UpdateRace(&race)

	if err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	if int(race.CurrentPlayer.ID) != body.CurrentPlayer || len(race.Responses) == 1 {
		err = c.raceService.ChangeTurn(race, false, body.CurrentPlayer)
	}

	request.FinalResponse(ctx, err, map[string]interface{}{
		"race": race,
	})
}

func (c *moderatorController) HandleUserRequest(ctx *gin.Context) {
	var err error

	var body dto.HandleUserRequestBodyDto

	if err = ctx.BindJSON(&body); err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	var userRequest entity.UserRequest
	var player entity.Player

	err, userRequest = c.userRequestService.HandleUserRequest(body)

	if userRequest.ID > 0 && err == nil {
		err, player = c.playerService.GetPlayerByUserIdAndRaceId(userRequest.RaceID, userRequest.UserID)

		cardType := entity.TransactionCardType.Payday

		if userRequest.Type == "baby" {
			cardType = entity.TransactionCardType.Baby
		}

		if player.ID > 0 && err == nil {
			err = c.playerService.UpdateCash(
				&player,
				userRequest.Amount,
				&dto.TransactionDTO{
					CardType: cardType,
					Details:  userRequest.Message,
				},
			)
		}
	}

	request.FinalResponse(ctx, err, nil)
}
