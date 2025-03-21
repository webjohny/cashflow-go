package controller

import (
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
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
	SendMoney(ctx *gin.Context)
	UpdatePlayer(ctx *gin.Context)
	UpdateRace(ctx *gin.Context)
	UpdateStatusRace(ctx *gin.Context)
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

	logger.Info("Moderator.GetRacePlayer", map[string]interface{}{
		"playerId": userId,
		"raceId":   raceId,
	})

	var err error
	var player entity.Player
	var response dto.GetRacePlayerResponseDTO

	if userId != 0 {
		err, player = c.playerService.GetPlayerByUserIdAndRaceId(raceId, userId)

		if err == nil {
			response = c.playerService.GetFormattedPlayerResponse(player, true)
			response.Data = player.Info.Data
		}
	}

	request.FinalResponse(ctx, err, response)
}

func (c *moderatorController) GetRacePlayers(ctx *gin.Context) {
	raceId := helper.GetRaceId(ctx)

	var err error

	players := c.raceService.GetRacePlayersByRaceId(raceId, false)

	request.FinalResponse(ctx, err, map[string]interface{}{
		"players": players,
	})
}

func (c *moderatorController) UpdatePlayer(ctx *gin.Context) {
	raceId := helper.GetRaceId(ctx)
	playerId := helper.ConvertToUInt64(ctx.Param("playerId"))

	var err error
	var body dto.ModeratorUpdatePlayerDto
	var isUpdateRace bool

	if err = ctx.BindJSON(&body); err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	err, player := c.playerService.GetPlayerByPlayerIdAndRaceId(raceId, playerId)

	if err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}
	isUpdateRace = player.IsActive != body.IsActive

	player.Cash = body.Cash
	player.CashFlow = body.CashFlow
	player.Babies = uint8(body.Babies)
	player.SkippedTurns = uint8(body.SkippedTurns)
	player.DualDiceCount = body.DualDiceCount
	player.ExtraDices = body.ExtraDices
	player.OnBigRace = body.OnBigRace
	player.IsActive = body.IsActive
	player.Assets.Savings = body.Savings

	if body.Liabilities.BankLoan > 0 {
		player.Expenses["bankLoan"] = body.Liabilities.BankLoan / 10
	} else {
		player.Expenses["bankLoan"] = 0
	}

	player.Liabilities.BankLoan = body.Liabilities.BankLoan

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

	err, _ = c.playerService.UpdatePlayer(&player)

	if err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	if isUpdateRace {
		race := c.raceService.GetRaceByRaceId(raceId)
		user := race.GetNextPlayer()
		race.Responses = c.raceService.CreateResponses(raceId, player.ID)
		race.CurrentPlayer = entity.RacePlayer{
			ID:       user.ID,
			UserId:   user.UserId,
			Username: user.Username,
		}
		err = c.raceService.ChangeTurn(race, true, int(user.ID))
	}

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

	race.Options.EnableCardCategory = body.EnableCardCategory
	race.Options.EnableManager = body.EnableManager
	race.Options.HideCards = body.HideCards
	race.Options.HandMode = body.HandMode
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
		err = c.raceService.ChangeTurn(race, true, body.CurrentPlayer)
	}

	request.FinalResponse(ctx, err, map[string]interface{}{
		"race": race,
	})
}

func (c *moderatorController) UpdateStatusRace(ctx *gin.Context) {
	raceId := helper.GetRaceId(ctx)

	var err error

	var body dto.ModeratorUpdateStatusRaceDto

	if err = ctx.BindJSON(&body); err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	race := c.raceService.GetRaceByRaceId(raceId)

	race.Status = body.Status

	if race.Status == entity.RaceStatus.FINISHED || race.Status == entity.RaceStatus.CANCELLED {
		_ = c.lobbyService.ChangeStatusByGameId(raceId, entity.LobbyStatus.Cancelled)
	}

	err, _ = c.raceService.UpdateRace(&race)

	if err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	request.FinalResponse(ctx, err, nil)
}

func (c *moderatorController) HandleUserRequest(ctx *gin.Context) {
	var err error

	var body dto.HandleUserRequestBodyDto

	if err = ctx.BindJSON(&body); err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	var userRequest entity.UserRequest

	err, userRequest = c.userRequestService.HandleUserRequest(body)

	if err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	var player entity.Player

	if userRequest.ID > 0 && userRequest.Status == 1 {
		err, player = c.playerService.GetPlayerByUserIdAndRaceId(userRequest.RaceID, userRequest.UserID)

		if player.ID > 0 && err == nil {
			if userRequest.Type == entity.TransactionCardType.TakeLoan {
				err = c.playerService.TakeLoan(player, userRequest.Amount)
			} else {
				err = c.playerService.UpdateCash(
					&player,
					userRequest.Amount,
					&dto.TransactionDTO{
						CardType: userRequest.Type,
						Details:  userRequest.Message,
					},
				)
			}
		}
	} else if userRequest.Status == 2 {
		err, player = c.playerService.GetPlayerByUserIdAndRaceId(userRequest.RaceID, userRequest.UserID)

		player.SetNotification(userRequest.RejectMessage, entity.NotificationTypes.Error)

		err, _ = c.playerService.UpdatePlayer(&player)
	}

	request.FinalResponse(ctx, err, nil)
}

func (c *moderatorController) SendMoney(ctx *gin.Context) {
	var err error
	raceId := helper.GetRaceId(ctx)

	var body dto.ModeratorSendMoneyDTO

	if err = ctx.BindJSON(&body); err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(raceId, uint64(body.Player))

	if err != nil {
		request.FinalResponse(ctx, err, nil)
		return
	}

	player.SetNotification(body.Message, entity.NotificationTypes.Success)

	err = c.playerService.UpdateCash(
		&player,
		body.Amount,
		&dto.TransactionDTO{
			CardType: entity.TransactionCardType.SendMoneyFromBank,
			Details:  body.Message,
		},
	)

	request.FinalResponse(ctx, err, nil)
}
