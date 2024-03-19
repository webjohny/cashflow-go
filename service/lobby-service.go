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
	CreateLobby(username string) (error, entity.Lobby)
	Join(ID uint64, username string) (error, entity.Lobby)
	Leave(ID uint64, username string) (error, entity.Lobby)
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

func (service *lobbyService) CreateLobby(username string) (error, entity.Lobby) {
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
		return fmt.Errorf(storage.ErrorUndefinedLobby), entity.Lobby{}
	}

	return nil, instance
}

func (service *lobbyService) Join(ID uint64, username string) (error, entity.Lobby) {
	lobby := service.lobbyRepository.FindLobbyById(ID)

	if lobby.ID != 0 {
		if lobby.IsFull() {
			return fmt.Errorf(storage.ErrorGameIsFull), entity.Lobby{}
		}

		if !lobby.IsGameStarted() && lobby.IsStarted() {
			return fmt.Errorf(storage.ErrorGameIsStarted), entity.Lobby{}
		}

		player := lobby.GetPlayer(username)

		if player.Username == "" {
			if lobby.IsGameStarted() {
				lobby.AddWaitList(username)
			} else if !lobby.IsStarted() {
				lobby.AddGuest(username)
			}
		}
	} else {
		return fmt.Errorf(storage.ErrorUndefinedLobby), entity.Lobby{}
	}

	return nil, lobby
}

func (service *lobbyService) Leave(ID uint64, username string) (error, entity.Lobby) {
	lobby := service.lobbyRepository.FindLobbyById(ID)

	if lobby.ID != 0 {
		lobby.RemovePlayer(username)

		if lobby.CountPlayers() == 0 {
			service.lobbyRepository.DeleteLobby(&lobby)
		}

		return nil, lobby
	}

	return fmt.Errorf(storage.ErrorUndefinedLobby), entity.Lobby{}
}
