package service

import (
	"fmt"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/storage"
	"gorm.io/datatypes"
	"log"
	"time"
)

type GameService interface {
	Start(lobbyId uint64) (error, entity.Race)
	RollDice(raceId uint64, userId uint64, dice int, isBigRace bool) (error, []int)
	GetGame(raceId uint64, lobbyId uint64, userId uint64, isBigRace bool) (error, dto.GetGameResponseDTO)
	ChangeTurn(raceId uint64, isBigRace bool) error
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

func (service *gameService) GetGame(raceId uint64, lobbyId uint64, userId uint64, isBigRace bool) (error, dto.GetGameResponseDTO) {
	player := service.playerService.GetPlayerByUserIdAndRaceId(raceId, userId)

	response := dto.GetGameResponseDTO{
		Username: player.Username,
		You:      player,
	}

	if lobbyId > 0 {
		lobby := service.lobbyRepository.FindLobbyById(lobbyId)
		//response.Lobby = &lobby
		response.Hash = helper.CreateHashByJson(lobby)
	} else if raceId > 0 {
		race := service.raceService.GetFormattedRaceResponse(raceId, isBigRace)
		response.CurrentPlayer = &race.CurrentPlayer
		//response.Race = &race
		response.Hash = helper.CreateHashByJson(race)
	}

	return nil, response
}

func (service *gameService) RollDice(raceId uint64, userId uint64, dice int, isBigRace bool) (error, []int) {
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

func (service *gameService) ChangeTurn(raceId uint64, isBigRace bool) error {
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

func (service *gameService) createResponses(raceId uint64) []entity.RaceResponse {
	players := service.playerService.GetAllPlayersByRaceId(raceId)

	responses := make([]entity.RaceResponse, 0)

	if len(players) > 0 {
		for _, player := range players {
			responses = append(responses, entity.RaceResponse{
				ID:        player.UserId,
				Username:  player.Username,
				Responded: false,
			})
		}
	}

	return responses
}

func (service *gameService) Start(lobbyId uint64) (error, entity.Race) {
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

	lobby.Status = entity.LobbyStatus.Started

	service.lobbyRepository.UpdateLobby(&lobby)

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

	for i := 0; i < len(lobby.Players); i++ {
		lobbyPlayer := lobby.Players[i]
		profession := service.professionRepository.PickProfession(&excluded)
		excluded = append(excluded, int(profession.ID))

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
	}

	race.CurrentPlayer = players[helper.Random(len(players)-1)]
	race.Responses = service.createResponses(race.ID)

	err, _ = service.raceService.UpdateRace(&race)

	if err != nil {
		log.Panic(err)

		return fmt.Errorf(storage.ErrorProcessFailed), entity.Race{}
	}

	return nil, race
}
