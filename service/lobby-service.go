package service

import (
	"fmt"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/storage"
)

type LobbyService interface {
	CreateLobby(username string) (error, *entity.Lobby)
	Join(ID uint64, username string) error
	Leave(ID uint64, username string) error
}

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
		MaxPlayers: 0,
		Status:     entity.LobbyStatus.STARTED,
		Options:    nil,
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

		if lobby.IsStarted() {
			return fmt.Errorf(storage.ErrorGameIsFull)
		}

		player := lobby.GetPlayer(username)

		if player != nil {
			lobby.AddGuest(username)
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
		return nil
	}

	return fmt.Errorf(storage.ErrorUndefinedLobby)
}
