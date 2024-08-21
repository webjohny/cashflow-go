package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/storage"
	"gorm.io/datatypes"
	"log"
	"time"
)

type LobbyService interface {
	Create(username string, userId uint64) (error, entity.Lobby)
	Update(lobby *entity.Lobby) (error, entity.Lobby)
	Join(ID uint64, username string, userId uint64) (error, entity.LobbyPlayer)
	Leave(ID uint64, username string) (error, entity.Lobby)
	Cancel(ID uint64, userId uint64) (error, entity.Lobby)
	GetLobby(lobbyId uint64, userId uint64) (error, dto.GetLobbyResponseDTO)
	GetByID(lobbyId uint64) entity.Lobby
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

func (service *lobbyService) GetByID(lobbyId uint64) entity.Lobby {
	return service.lobbyRepository.FindLobbyById(lobbyId)
}

func (service *lobbyService) Update(lobby *entity.Lobby) (error, entity.Lobby) {
	return service.lobbyRepository.UpdateLobby(lobby)
}

func (service *lobbyService) GetLobby(lobbyId uint64, userId uint64) (error, dto.GetLobbyResponseDTO) {
	logger.Info("LobbyService.GetLobby", map[string]interface{}{
		"lobbyId": lobbyId,
		"userId":  userId,
	})

	lobby := service.lobbyRepository.FindLobbyById(lobbyId)
	player := lobby.GetPlayer(userId)

	if player.ID == 0 {
		return errors.New(storage.ErrorUndefinedPlayer), dto.GetLobbyResponseDTO{}
	}

	response := dto.GetLobbyResponseDTO{
		Username: player.Username,
		You:      player,
		Players:  lobby.Players,
		Status:   lobby.Status,
		LobbyId:  lobby.ID,
		GameId:   lobby.GameId,
		Hash:     helper.CreateHashByJson(lobby),
	}

	return nil, response
}

func (service *lobbyService) Create(username string, userId uint64) (error, entity.Lobby) {
	logger.Info("LobbyService.Create", map[string]interface{}{
		"username": username,
		"userId":   userId,
	})

	lobby := &entity.Lobby{
		Players:    make([]entity.LobbyPlayer, 0),
		MaxPlayers: LobbyMaxPlayers,
		Status:     entity.LobbyStatus.New,
		Options:    make(map[string]interface{}),
		CreatedAt:  datatypes.Date(time.Now()),
	}
	lobby.AddOwner(userId, username)
	err, instance := service.lobbyRepository.InsertLobby(lobby)

	if err != nil {
		return err, entity.Lobby{}
	}

	if instance.ID == 0 {
		return errors.New(storage.ErrorUndefinedLobby), entity.Lobby{}
	}

	return nil, instance
}

func (service *lobbyService) Join(ID uint64, username string, userId uint64) (error, entity.LobbyPlayer) {
	logger.Info("LobbyService.Join", map[string]interface{}{
		"lobbyId":  ID,
		"username": username,
		"userId":   userId,
	})

	var player entity.LobbyPlayer
	lobby := service.lobbyRepository.FindLobbyById(ID)

	log.Println("LobbyService.Join:", ID, username, userId)

	if lobby.ID != 0 {
		if lobby.IsFull() {
			return errors.New(storage.ErrorGameIsFull), entity.LobbyPlayer{}
		}

		if !lobby.IsGameStarted() && lobby.IsStarted() {
			return errors.New(storage.ErrorGameIsStarted), entity.LobbyPlayer{}
		}

		player = lobby.GetPlayer(userId)

		log.Println("LobbyService.Join: exists lobby", ID, userId, player)

		if player.ID == 0 {
			if lobby.IsGameStarted() {
				log.Println("LobbyService.Join.Waitlist:", ID, userId)
				//@toDo add waitlist in game
				//game.AddWaitList()
				lobby.AddWaitList(userId, username)
			} else if !lobby.IsStarted() {
				log.Println("LobbyService.Join.Guest:", ID, userId)

				lobby.AddGuest(userId, username)
			}

			err, _ := service.lobbyRepository.UpdateLobby(&lobby)

			if err != nil {
				return err, entity.LobbyPlayer{}
			}

			player = lobby.GetPlayer(userId)
		}
	} else {
		return errors.New(storage.ErrorUndefinedLobby), entity.LobbyPlayer{}
	}

	return nil, player
}

func (service *lobbyService) Leave(ID uint64, username string) (error, entity.Lobby) {
	logger.Info("LobbyService.Leave", map[string]interface{}{
		"lobbyId":  ID,
		"username": username,
	})

	lobby := service.lobbyRepository.FindLobbyById(ID)

	log.Println("LobbyService.Leave:", lobby.ID, ID, username)

	if lobby.ID != 0 {
		log.Println("LobbyService.Leave: exists lobby", lobby.CountPlayers(), ID, username)

		lobby.RemovePlayer(username)

		if lobby.CountPlayers() == 0 {
			service.lobbyRepository.DeleteLobby(&lobby)
		}

		err, _ := service.lobbyRepository.UpdateLobby(&lobby)

		if err != nil {
			return err, entity.Lobby{}
		}

		return nil, lobby
	}

	return errors.New(storage.ErrorUndefinedLobby), entity.Lobby{}
}

func (service *lobbyService) Cancel(ID uint64, userId uint64) (error, entity.Lobby) {
	logger.Info("LobbyService.Cancel", map[string]interface{}{
		"lobbyId": ID,
		"userId":  userId,
	})

	lobby := service.lobbyRepository.FindLobbyById(ID)

	if lobby.ID != 0 {
		logger.Info("LobbyService.Cancel exists lobby", map[string]interface{}{
			"countPlayers": lobby.CountPlayers(),
			"lobbyId":      ID,
			"userId":       userId,
		})

		service.lobbyRepository.CancelLobby(&lobby)

		return nil, lobby
	}

	return errors.New(storage.ErrorUndefinedLobby), entity.Lobby{}
}
