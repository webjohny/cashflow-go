package service

import (
	"fmt"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/storage"
)

type GameService interface {
	Start(lobbyId uint64) (error, *entity.Race)
	GetGame(raceId uint64, lobbyId uint64, username string, isBigRace *bool) (error, dto.GetGameResponseDTO)
}

type gameService struct {
	lobbyRepository  repository.LobbyRepository
	raceRepository   repository.RaceRepository
	playerRepository repository.PlayerRepository
}

func NewGameService(lobbyRepository repository.LobbyRepository, raceRepository repository.RaceRepository, playerRepository repository.PlayerRepository) GameService {
	return &gameService{
		lobbyRepository:  lobbyRepository,
		raceRepository:   raceRepository,
		playerRepository: playerRepository,
	}
}

func (service *gameService) GetGame(raceId uint64, lobbyId uint64, username string, isBigRace *bool) (error, dto.GetGameResponseDTO) {
	player := service.playerRepository.FindPlayerByUsername(username)

	response := dto.GetGameResponseDTO{
		Username: username,
		You:      *player,
	}
	if lobbyId > 0 {
		lobby := service.lobbyRepository.FindLobbyById(lobbyId)
		response.Lobby = lobby
		response.Hash = helper.CreateHashByJson(lobby)
	} else if raceId > 0 {
		bigRace := player.OnBigRace

		if isBigRace != nil {
			bigRace = *isBigRace
		}

		race := service.raceRepository.FindRaceById(raceId, bigRace)
		response.Race = race
		response.Hash = helper.CreateHashByJson(race)
	}

	return nil, response
}

func (service *gameService) Start(lobbyId uint64) (error, *entity.Race) {
	lobby := service.lobbyRepository.FindLobbyById(lobbyId)

	if lobby == nil {
		return fmt.Errorf(storage.ErrorUndefinedLobby), nil
	}

	if !lobby.AvailableToStart() {
		return fmt.Errorf(storage.ErrorInsufficientPlayers), nil
	}

	race := service.raceRepository.InsertRace(&entity.Race{
		Responses:         make([]entity.RaceResponse, 0),
		Status:            entity.RaceStatus.STARTED,
		Notifications:     make([]entity.RaceNotification, 0),
		BankruptedPlayers: make([]entity.RaceBankruptPlayer, 0),
		Logs:              make([]entity.RaceLog, 0),
		Dice:              make([]int, 0),
		Options:           entity.RaceOptions{},
	})

	for i := 0; i < len(lobby.Players); i++ {
		player := lobby.Players[i]
		service.playerRepository.InsertPlayer(&entity.Player{
			RaceId:      race.ID,
			Username:    player.Username,
			Role:        player.Role,
			Color:       player.Color,
			Income:      entity.PlayerIncome{},
			Expenses:    make(map[string]int),
			Assets:      entity.PlayerAssets{},
			Liabilities: entity.PlayerLiabilities{},
		})
	}

	return nil, &race
}
