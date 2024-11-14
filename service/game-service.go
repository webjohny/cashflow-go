package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/storage"
	"gorm.io/datatypes"
	"log"
	"time"
)

type GameService interface {
	Start(lobbyId uint64) (error, entity.Race)
	RollDice(raceId uint64, userId uint64, dto dto.RollDiceDto, isBigRace bool) (error, []int)
	GetGame(raceId uint64, userId uint64, isBigRace bool) (error, dto.GetGameResponseDTO)
	ChangeTurn(raceId uint64) error
	Cancel(raceId uint64, userId uint64) error
	Reset(raceId uint64, userId uint64) error
	GetTiles(raceId uint64, isBigRace bool) []string
}

type gameService struct {
	raceService       RaceService
	playerService     PlayerService
	lobbyService      LobbyService
	professionService ProfessionService
}

func NewGameService(
	raceService RaceService,
	playerService PlayerService,
	lobbyService LobbyService,
	professionService ProfessionService,
) GameService {
	return &gameService{
		raceService:       raceService,
		playerService:     playerService,
		lobbyService:      lobbyService,
		professionService: professionService,
	}
}

func (service *gameService) GetGame(raceId uint64, userId uint64, isBigRace bool) (error, dto.GetGameResponseDTO) {
	err, player := service.playerService.GetPlayerByUserIdAndRaceId(raceId, userId)

	if err != nil {
		return err, dto.GetGameResponseDTO{}
	}

	you := service.playerService.GetFormattedPlayerResponse(player, true)

	response := dto.GetGameResponseDTO{
		Username:          player.Username,
		You:               you,
		Notifications:     make([]entity.RaceNotification, 0),
		Logs:              make([]entity.RaceLog, 0),
		TurnResponses:     make([]entity.RaceResponse, 0),
		BankruptedPlayers: make([]dto.GetRacePlayerResponseDTO, 0),
	}

	race := service.raceService.GetFormattedRaceResponse(raceId)
	response.TurnResponses = race.TurnResponses
	response.Players = race.Players
	response.CurrentCard = &race.CurrentCard
	response.CurrentPlayer = &race.CurrentPlayer
	response.GameId = race.GameId
	response.IsTurnEnded = race.IsTurnEnded
	response.IsMultiFlow = race.IsMultiFlow
	response.Status = race.Status
	response.DiceValues = race.DiceValues
	response.Logs = race.Logs
	response.Hash = helper.CreateHashByJson(race)

	return nil, response
}

func (service *gameService) RollDice(raceId uint64, userId uint64, dto dto.RollDiceDto, isBigRace bool) (error, []int) {
	logger.Info("GameService.RollDice", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"dto":    dto,
	})

	err, race, player := service.raceService.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err, []int{}
	}

	getDice := race.GetDice()

	var dice = make([]int, 0)

	if dto.Dices > 0 {
		dice = getDice.Roll(dto.Dices)
	}

	if dto.IsFinished {
		dualDiceCount := player.DualDiceCount

		var totalCount int

		if len(player.Dices) > 0 {
			if dto.Dices > 0 {
				player.AddDices(dice)
			}
			race.Dice = player.Dices
		} else {
			race.Dice = dice
		}

		totalCount = race.CalculateDices()

		if dualDiceCount > 0 {
			player.DecrementDualDiceCount()
		}
		player.Dices = make([]int, 0)
		player.ChangeDiceStatus(true)
		player.Move(totalCount)
		race.ResetResponses()
	} else {
		player.AddDices(dice)
	}

	err, _ = service.playerService.UpdatePlayer(&player)
	err, _ = service.raceService.UpdateRace(&race)

	return err, dice
}

func (service *gameService) Cancel(raceId uint64, userId uint64) error {
	logger.Info("GameService.Cancel", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	players := service.playerService.GetAllPlayersByRaceId(raceId)
	race := service.raceService.GetRaceByRaceId(raceId)

	if race.ID == 0 {
		return errors.New(storage.ErrorUndefinedGame)
	}

	for _, player := range players {
		if player.UserID == userId && player.Role != entity.PlayerRoles.Owner {
			return errors.New(storage.ErrorPermissionDenied)
		}
	}

	race.Status = entity.LobbyStatus.Cancelled
	err, _ := service.raceService.UpdateRace(&race)

	return err
}

func (service *gameService) Reset(raceId uint64, userId uint64) error {
	logger.Info("GameService.Reset", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	players := service.playerService.GetAllPlayersByRaceId(raceId)
	race := service.raceService.GetRaceByRaceId(raceId)

	if race.ID == 0 {
		return errors.New(storage.ErrorUndefinedGame)
	}

	for _, player := range players {
		if player.UserID == userId && player.Role != entity.PlayerRoles.Owner {
			return errors.New(storage.ErrorPermissionDenied)
		}
	}

	for _, player := range players {
		profession := service.professionService.GetByID(uint64(player.ProfessionID))

		profession.Assets.Business = make([]entity.CardBusiness, 0)
		profession.Assets.RealEstates = make([]entity.CardRealEstate, 0)
		profession.Assets.OtherAssets = make([]entity.CardOtherAssets, 0)
		profession.Assets.Stocks = make([]entity.CardStocks, 0)
		profession.Assets.Dreams = make([]entity.CardDream, 0)

		_, _ = service.playerService.UpdatePlayer(&entity.Player{
			ID:           player.ID,
			UserID:       player.UserID,
			RaceID:       race.ID,
			Username:     player.Username,
			Role:         player.Role,
			Color:        player.Color,
			Salary:       profession.Income.Salary,
			Babies:       uint8(profession.Babies),
			Expenses:     profession.Expenses,
			Assets:       profession.Assets,
			Liabilities:  profession.Liabilities,
			ProfessionID: uint8(profession.ID),
		})
	}

	race.CurrentPlayer = entity.RacePlayer{
		ID:       players[0].ID,
		UserId:   players[0].UserID,
		Username: players[0].Username,
	}
	race.Responses = service.createResponses(race.ID, race.CurrentPlayer.ID)
	race.CurrentCard = entity.Card{}
	race.Status = entity.RaceStatus.STARTED
	race.Notifications = make([]entity.RaceNotification, 0)
	race.BankruptedPlayers = make([]entity.RaceBankruptPlayer, 0)
	race.Logs = make([]entity.RaceLog, 0)
	race.Dice = make([]int, 0)
	race.Options = entity.RaceOptions{}
	race.CreatedAt = datatypes.Date(time.Now())
	err, _ := service.raceService.UpdateRace(&race)

	return err
}

func (service *gameService) GetTiles(raceId uint64, isBigRace bool) []string {
	logger.Info("GameService.GetTiles", map[string]interface{}{
		"raceId": raceId,
	})

	if isBigRace {
		return []string{
			"dream", "business", "dream", "business", "dream", "business", "dream", "bigCharity", "business", "dream", "business",
			"cashFlowDay", "dream", "business", "dream", "tax50percent", "dream", "business", "dream", "business", "dream", "business", "dream", "business", "dream", "cashFlowDay", "dream", "business", "dream", "tax100percent",
			"dream", "business", "dream", "business", "dream", "business", "dream", "business", "dream", "business", "dream", "cashFlowDay", "dream", "business", "dream", "bankrupt",
		}
	}

	return []string{
		"deal", "doodad", "deal", "charity", "deal", "payday", "deal",
		"market", "deal", "doodad", "deal", "downsized", "deal", "payday", "deal",
		"market", "deal", "doodad", "deal", "baby", "deal", "payday", "deal", "market",
	}
}

func (service *gameService) ChangeTurn(raceId uint64) error {
	logger.Info("GameService.ChangeTurn", map[string]interface{}{
		"raceId": raceId,
	})

	race := service.raceService.GetRaceByRaceId(raceId)

	if race.ID == 0 {
		return errors.New(storage.ErrorUndefinedGame)
	}

	return service.raceService.ChangeTurn(race, 0)
}

func (service *gameService) Start(lobbyId uint64) (error, entity.Race) {
	logger.Info("GameService.Start", map[string]interface{}{
		"lobbyId": lobbyId,
	})

	lobby := service.lobbyService.GetByID(lobbyId)

	if lobby.ID == 0 {
		return errors.New(storage.ErrorUndefinedLobby), entity.Race{}
	}

	if !lobby.AvailableToStart() {
		return errors.New(storage.ErrorInsufficientPlayers), entity.Race{}
	}

	if lobby.IsStarted() {
		return errors.New(storage.ErrorGameIsStarted), entity.Race{}
	}

	err, race := service.raceService.InsertRace(&entity.Race{
		Responses:         make([]entity.RaceResponse, 0),
		Status:            entity.RaceStatus.STARTED,
		Notifications:     make([]entity.RaceNotification, 0),
		BankruptedPlayers: make([]entity.RaceBankruptPlayer, 0),
		Logs:              make([]entity.RaceLog, 0),
		Dice:              make([]int, 0),
		Options:           entity.RaceOptions{},
		CreatedAt:         datatypes.Date(time.Now()),
	})

	if err != nil {
		return errors.New(storage.ErrorCannotCreatedRace), entity.Race{}
	}

	var excluded []int
	var players []entity.RacePlayer
	var responses []entity.RaceResponse

	for i := 0; i < len(lobby.Players); i++ {
		lobbyPlayer := lobby.Players[i]
		profession := service.professionService.GetRandomProfession(&excluded)
		excluded = append(excluded, int(profession.ID))

		profession.Assets.Business = make([]entity.CardBusiness, 0)
		profession.Assets.Dreams = make([]entity.CardDream, 0)

		playerErr, player := service.playerService.InsertPlayer(&entity.Player{
			UserID:       lobbyPlayer.ID,
			RaceID:       race.ID,
			Username:     lobbyPlayer.Username,
			Role:         lobbyPlayer.Role,
			Color:        lobbyPlayer.Color,
			Salary:       profession.Income.Salary,
			Babies:       uint8(profession.Babies),
			Expenses:     profession.Expenses,
			Assets:       profession.Assets,
			Liabilities:  profession.Liabilities,
			ProfessionID: uint8(profession.ID),
			CreatedAt:    datatypes.Date(time.Now()),
		})

		if playerErr != nil {
			log.Panic(playerErr)
		}

		players = append(players, entity.RacePlayer{
			ID:       player.ID,
			UserId:   player.UserID,
			Username: player.Username,
		})

		responses = append(responses, player.CreateResponse())
	}

	race.CurrentPlayer = players[helper.Random(len(players)-1)]
	race.Responses = responses

	err, _ = service.raceService.UpdateRace(&race)

	lobby.Status = entity.LobbyStatus.Started
	lobby.GameId = race.ID
	err, _ = service.lobbyService.Update(&lobby)

	if err != nil {
		return errors.New(storage.ErrorProcessFailed), entity.Race{}
	}

	return nil, race
}

func (service *gameService) createResponses(raceId uint64, currentPlayerId uint64) []entity.RaceResponse {
	players := service.playerService.GetAllPlayersByRaceId(raceId)

	responses := make([]entity.RaceResponse, 0)

	if len(players) > 0 {
		for _, player := range players {
			response := player.CreateResponse()

			if player.OnBigRace && player.ID != currentPlayerId {
				response.Responded = true
			}
			responses = append(responses, response)
		}
	}

	return responses
}
