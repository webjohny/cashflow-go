package service

import (
	"fmt"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/logger"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/storage"
	"gorm.io/datatypes"
	"log"
	"time"
)

type GameService interface {
	Start(lobbyId uint64) (error, entity.Race)
	RollDice(raceId uint64, userId uint64, dice int, isBigRace bool) (error, []int)
	GetGame(raceId uint64, userId uint64, isBigRace bool) dto.GetGameResponseDTO
	ChangeTurn(raceId uint64, isBigRace bool) error
	Cancel(raceId uint64, userId uint64) error
	Reset(raceId uint64, userId uint64) error
	GetTiles(raceId uint64, isBigRace bool) []string
}

type gameService struct {
	raceService          RaceService
	playerService        PlayerService
	lobbyRepository      repository.LobbyRepository
	professionRepository repository.ProfessionRepository
}

func NewGameService(
	raceService RaceService,
	playerService PlayerService,
	lobbyRepository repository.LobbyRepository,
	professionRepository repository.ProfessionRepository,
) GameService {
	return &gameService{
		raceService:          raceService,
		playerService:        playerService,
		lobbyRepository:      lobbyRepository,
		professionRepository: professionRepository,
	}
}

func (service *gameService) GetGame(raceId uint64, userId uint64, isBigRace bool) dto.GetGameResponseDTO {
	_, player := service.playerService.GetRacePlayer(raceId, userId)

	response := dto.GetGameResponseDTO{
		Username:          player.Username,
		You:               player,
		Notifications:     make([]entity.RaceNotification, 0),
		Logs:              make([]entity.RaceLog, 0),
		TurnResponses:     make([]entity.RaceResponse, 0),
		BankruptedPlayers: make([]dto.GetRacePlayerResponseDTO, 0),
	}

	race := service.raceService.GetFormattedRaceResponse(raceId, isBigRace)
	response.TurnResponses = race.TurnResponses
	response.Players = race.Players
	response.CurrentCard = &race.CurrentCard
	response.CurrentPlayer = &race.CurrentPlayer
	response.GameId = race.GameId
	//response.Race = &race
	response.Hash = helper.CreateHashByJson(race)

	return response
}

func (service *gameService) RollDice(raceId uint64, userId uint64, dice int, isBigRace bool) (error, []int) {
	logger.Info("GameService.RollDice", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"dice":   dice,
	})

	err, race, player := service.raceService.GetRaceAndPlayer(raceId, userId, isBigRace)

	if err != nil {
		return err, []int{}
	}

	if dice == 0 {
		return fmt.Errorf(storage.ErrorUndefinedDiceValue), []int{}
	}

	getDice := race.GetDice()

	diceValues := getDice.Roll(dice)

	dualDiceCount := player.DualDiceCount

	totalCount := race.CalculateTotalSteps(diceValues, dice)

	if dualDiceCount > 0 {
		player.DecrementDualDiceCount()
	}

	player.ChangeDiceStatus(true)
	player.Move(totalCount)

	err, _ = service.playerService.UpdatePlayer(&player)

	//this.addLog(currentPlayer.username, `rolled ${totalCount}`);

	return err, diceValues
}

func (service *gameService) ReRollDice(raceId uint64, userId uint64, dice int, isBigRace bool) (error, []int) {
	logger.Info("GameService.ReRollDice", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"dice":   dice,
	})

	err, race, player := service.raceService.GetRaceAndPlayer(raceId, userId, isBigRace)

	if err != nil {
		return err, []int{}
	}

	if dice == 0 {
		return fmt.Errorf(storage.ErrorUndefinedDiceValue), []int{}
	}

	log.Println(player, race)

	//console.log('Game.reRollDice');
	//
	//this.#diceValues = this.#dice.roll(1);
	//const diceValue = this.#diceValues[0];
	//const currentPlayer = this.currentPlayer;
	//
	//this.addLog(username, `rolled ${diceValue} again`);
	//this.#currentTurn.lottery(this.#players, currentPlayer, diceValue);
	//currentPlayer.changeDiceStatus(true);
	//
	//await this.insertData();
	//
	//return this.#diceValues;

	return err, []int{}
}

func (service *gameService) Cancel(raceId uint64, userId uint64) error {
	logger.Info("GameService.Cancel", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	players := service.playerService.GetAllPlayersByRaceId(raceId)
	race := service.raceService.GetRaceByRaceId(raceId, false)

	if race.ID == 0 {
		return fmt.Errorf(storage.ErrorUndefinedGame)
	}

	for _, player := range players {
		if player.UserId == userId && player.Role != entity.PlayerRoles.Owner {
			return fmt.Errorf(storage.ErrorPermissionDenied)
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
	race := service.raceService.GetRaceByRaceId(raceId, false)

	if race.ID == 0 {
		return fmt.Errorf(storage.ErrorUndefinedGame)
	}

	for _, player := range players {
		if player.UserId == userId && player.Role != entity.PlayerRoles.Owner {
			return fmt.Errorf(storage.ErrorPermissionDenied)
		}
	}

	for _, player := range players {
		profession := service.professionRepository.FindProfessionById(uint64(player.ProfessionId))

		profession.Income.Business = make([]entity.CardBusiness, 0)
		profession.Liabilities.Business = make([]entity.CardBusiness, 0)
		profession.Assets.Business = make([]entity.CardBusiness, 0)
		profession.Assets.Dreams = make([]entity.CardDream, 0)

		_, _ = service.playerService.UpdatePlayer(&entity.Player{
			UserId:       player.UserId,
			RaceId:       race.ID,
			Username:     player.Username,
			Role:         player.Role,
			Color:        player.Color,
			Income:       profession.Income,
			Babies:       uint8(profession.Babies),
			Expenses:     profession.Expenses,
			Assets:       profession.Assets,
			Liabilities:  profession.Liabilities,
			ProfessionId: uint8(profession.ID),
			CreatedAt:    datatypes.Date(time.Now()),
		})
	}

	race.Responses = service.createResponses(race.ID)
	race.CurrentPlayer = entity.RacePlayer{
		ID:       players[0].ID,
		UserId:   players[0].UserId,
		Username: players[0].Username,
	}
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

func (service *gameService) ChangeTurn(raceId uint64, isBigRace bool) error {
	logger.Info("GameService.ChangeTurn", map[string]interface{}{
		"raceId": raceId,
	})

	race := service.raceService.GetRaceByRaceId(raceId, isBigRace)

	if race.ID == 0 {
		return fmt.Errorf(storage.ErrorUndefinedGame)
	}

	race.NextPlayer()

	currentPlayer := race.CurrentPlayer

	player := service.playerService.GetPlayerByUserIdAndRaceId(raceId, currentPlayer.ID)

	if player.ID == 0 {
		return fmt.Errorf(storage.ErrorUndefinedPlayer)
	}

	var err error

	if player.HasBankrupt == 1 {
		return service.ChangeTurn(raceId, isBigRace)
	} else {
		player.ChangeDiceStatus(false)
		race.CurrentCard = entity.Card{}
		race.Notifications = make([]entity.RaceNotification, 0)
		race.Responses = service.createResponses(raceId)

		if player.SkippedTurns > 0 {
			player.DecrementSkippedTurns()

			if player.DualDiceCount > 0 {
				player.DecrementDualDiceCount()
			}

			err = service.ChangeTurn(raceId, isBigRace)

			if err != nil {
				return err
			}
		}
	}

	err, _ = service.playerService.UpdatePlayer(&player)

	if err != nil {
		return err
	}

	err, _ = service.raceService.UpdateRace(&race)

	return err
}

func (service *gameService) Start(lobbyId uint64) (error, entity.Race) {
	logger.Info("GameService.Start", map[string]interface{}{
		"lobbyId": lobbyId,
	})

	lobby := service.lobbyRepository.FindLobbyById(lobbyId)

	if lobby.ID == 0 {
		return fmt.Errorf(storage.ErrorUndefinedLobby), entity.Race{}
	}

	if !lobby.AvailableToStart() {
		return fmt.Errorf(storage.ErrorInsufficientPlayers), entity.Race{}
	}

	if lobby.IsStarted() {
		return fmt.Errorf(storage.ErrorGameIsStarted), entity.Race{}
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
		log.Panic(err)

		return fmt.Errorf(storage.ErrorCannotCreatedRace), entity.Race{}
	}

	var excluded []int
	var players []entity.RacePlayer
	var responses []entity.RaceResponse

	for i := 0; i < len(lobby.Players); i++ {
		lobbyPlayer := lobby.Players[i]
		profession := service.professionRepository.PickProfession(&excluded)
		excluded = append(excluded, int(profession.ID))

		profession.Income.Business = make([]entity.CardBusiness, 0)
		profession.Liabilities.Business = make([]entity.CardBusiness, 0)
		profession.Assets.Business = make([]entity.CardBusiness, 0)
		profession.Assets.Dreams = make([]entity.CardDream, 0)

		playerErr, player := service.playerService.InsertPlayer(&entity.Player{
			UserId:       lobbyPlayer.ID,
			RaceId:       race.ID,
			Username:     lobbyPlayer.Username,
			Role:         lobbyPlayer.Role,
			Color:        lobbyPlayer.Color,
			Income:       profession.Income,
			Babies:       uint8(profession.Babies),
			Expenses:     profession.Expenses,
			Assets:       profession.Assets,
			Liabilities:  profession.Liabilities,
			ProfessionId: uint8(profession.ID),
			CreatedAt:    datatypes.Date(time.Now()),
		})

		if playerErr != nil {
			log.Panic(playerErr)
		}

		players = append(players, entity.RacePlayer{
			ID:       player.ID,
			UserId:   player.UserId,
			Username: player.Username,
		})

		responses = append(responses, player.CreateResponse())
	}

	race.CurrentPlayer = players[helper.Random(len(players)-1)]
	race.Responses = responses

	err, _ = service.raceService.UpdateRace(&race)

	lobby.Status = entity.LobbyStatus.Started
	lobby.GameId = race.ID
	_ = service.lobbyRepository.UpdateLobby(&lobby)

	if err != nil {
		log.Panic(err)

		return fmt.Errorf(storage.ErrorProcessFailed), entity.Race{}
	}

	return nil, race
}

func (service *gameService) createResponses(raceId uint64) []entity.RaceResponse {
	players := service.playerService.GetAllPlayersByRaceId(raceId)

	responses := make([]entity.RaceResponse, 0)

	if len(players) > 0 {
		for _, player := range players {
			responses = append(responses, player.CreateResponse())
		}
	}

	return responses
}
