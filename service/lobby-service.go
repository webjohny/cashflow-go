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
	Create(username string, userId uint64) (error, entity.Lobby)
	Join(ID uint64, username string, userId uint64) (error, entity.LobbyPlayer)
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

func (service *lobbyService) Create(username string, userId uint64) (error, entity.Lobby) {
	lobby := &entity.Lobby{
		Players:    make([]entity.LobbyPlayer, 0),
		MaxPlayers: LobbyMaxPlayers,
		Status:     entity.LobbyStatus.New,
		Options:    make(map[string]interface{}),
		CreatedAt:  datatypes.Date(time.Now()),
	}
	lobby.AddOwner(userId, username)
	instance := service.lobbyRepository.InsertLobby(lobby)

	if instance.ID == 0 {
		return fmt.Errorf(storage.ErrorUndefinedLobby), entity.Lobby{}
	}

	return nil, instance
}

func (service *lobbyService) Join(ID uint64, username string, userId uint64) (error, entity.LobbyPlayer) {
	var player entity.LobbyPlayer
	lobby := service.lobbyRepository.FindLobbyById(ID)

	if lobby.ID != 0 {
		if lobby.IsFull() {
			return fmt.Errorf(storage.ErrorGameIsFull), entity.LobbyPlayer{}
		}

		if !lobby.IsGameStarted() && lobby.IsStarted() {
			return fmt.Errorf(storage.ErrorGameIsStarted), entity.LobbyPlayer{}
		}

		player = lobby.GetPlayer(userId)

		if player.ID == 0 {
			if lobby.IsGameStarted() {
				//@toDo add waitlist in game
				//game.AddWaitList()
				lobby.AddWaitList(userId, username)
			} else if !lobby.IsStarted() {
				lobby.AddGuest(userId, username)
			}

			player = lobby.GetPlayer(userId)
		}
	} else {
		return fmt.Errorf(storage.ErrorUndefinedLobby), entity.LobbyPlayer{}
	}

	return nil, player
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
