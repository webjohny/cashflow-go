package service

import (
	"fmt"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/storage"
	"gorm.io/datatypes"
	"time"
)

type LobbyService interface {
	CreateLobby(username string) (error, *entity.Lobby)
	Join(ID uint64, username string) error
	Leave(ID uint64, username string) error
}

const LobbyMaxPlayers = 6

type lobbyService struct {
	lobbyRepository repository.LobbyRepository
}

func NewLobbyService(lobbyRepository repository.LobbyRepository) LobbyService {
	return &lobbyService{
		lobbyRepository: lobbyRepository,
	}
}

func (service *lobbyService) CreateLobby(username string) (error, *entity.Lobby) {
	lobby := &entity.Lobby{
		Players:    make([]entity.LobbyPlayer, 0),
		MaxPlayers: LobbyMaxPlayers,
		Status:     entity.LobbyStatus.New,
		Options:    make(map[string]interface{}),
		CreatedAt:  datatypes.Date(time.Now()),
	}
	lobby.AddOwner(username)
	instance := service.lobbyRepository.InsertLobby(lobby)

	if instance.ID == 0 {
		return fmt.Errorf(storage.ErrorUndefinedLobby), nil
	}

	return nil, lobby
}

func (service *lobbyService) Join(ID uint64, username string) error {
	lobby := service.lobbyRepository.FindLobbyById(ID)

	if lobby != nil {
		if lobby.IsFull() {
			return fmt.Errorf(storage.ErrorGameIsFull)
		}

		if lobby.IsGameStarted() {
			return fmt.Errorf(storage.ErrorGameIsStarted)
		}

		player := lobby.GetPlayer(username)

		if player != nil {
			if lobby.IsGameStarted() {
				lobby.AddWaitList(username)
			} else {
				lobby.AddGuest(username)
			}
		}
	} else {
		return fmt.Errorf(storage.ErrorUndefinedLobby)
	}

	return nil
}

func (service *lobbyService) Leave(ID uint64, username string) error {
	lobby := service.lobbyRepository.FindLobbyById(ID)

	if lobby != nil {
		lobby.RemovePlayer(username)

		if lobby.CountPlayers() == 0 {
			service.lobbyRepository.DeleteLobby(lobby)
		}

		return nil
	}

	return fmt.Errorf(storage.ErrorUndefinedLobby)
}
