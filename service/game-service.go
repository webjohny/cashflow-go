package service

import (
	"fmt"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/storage"
)

type GameService interface {
	Start(ID uint64, username string) error
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

func (service *gameService) Start(ID uint64, username string) error {
	lobby := service.lobbyRepository.FindLobbyById(ID)

	if lobby != nil {
		lobby.RemovePlayer(username)
		return nil
	}

	if !lobby.AvailableToStart() {
		return fmt.Errorf(storage.ErrorInsufficientPlayers)
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

	return fmt.Errorf(storage.ErrorUndefinedLobby)
}
