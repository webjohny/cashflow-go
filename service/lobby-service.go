package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/storage"
	"log"
	"time"
)

type LobbyService interface {
	GetByID(lobbyId uint64) entity.Lobby
	Create(username string, userId uint64) (error, entity.Lobby)
	SetOptions(lobbyId uint64, body dto.SetOptionsLobbyRequestDTO) error
	Update(lobby *entity.Lobby) (error, entity.Lobby)
	Leave(ID uint64, userId uint64) (error, entity.Lobby)
	Cancel(ID uint64, userId uint64) (error, entity.Lobby)
	Join(ID uint64, username string, userId uint64) (error, entity.LobbyPlayer)
	GetLobby(lobbyId uint64, userId uint64) (error, dto.GetLobbyResponseDTO)
	ChangeStatusByGameId(gameId uint64, status string) error
	ChangeRoleByGameIdAndUserId(gameId uint64, userId uint64, role string) error
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

func (service *lobbyService) SetOptions(lobbyId uint64, body dto.SetOptionsLobbyRequestDTO) error {
	lobby := service.lobbyRepository.FindLobbyById(lobbyId)

	if lobby.ID == 0 {
		return errors.New(storage.ErrorUndefinedLobby)
	}

	lobby.Options.HandMode = body.HandMode
	lobby.Options.MeetLink = body.MeetLink
	lobby.Options.BannerLink = body.BannerLink
	lobby.Options.BannerImage = body.BannerImage
	lobby.Options.Language = body.Language

	err, _ := service.Update(&lobby)

	return err
}

func (service *lobbyService) ChangeRoleByGameIdAndUserId(gameId uint64, userId uint64, role string) error {
	logger.Info("LobbyService.ChangeRoleByGameIdAndUserId", map[string]interface{}{
		"gameId": gameId,
		"userId": userId,
		"role":   role,
	})

	lobby := service.lobbyRepository.FindLobbyByGameId(gameId)

	if lobby.ID == 0 {
		return errors.New(storage.ErrorUndefinedLobby)
	}

	lobby.ChangePlayerRole(userId, role)

	err, _ := service.lobbyRepository.UpdateLobby(&lobby)

	if err != nil {
		logger.Error(err)
	}

	return err
}

func (service *lobbyService) ChangeStatusByGameId(gameId uint64, status string) error {
	logger.Info("LobbyService.ChangeStatusByGameId", map[string]interface{}{
		"gameId": gameId,
		"status": status,
	})

	lobby := service.lobbyRepository.FindLobbyByGameId(gameId)

	if lobby.ID == 0 {
		return errors.New(storage.ErrorUndefinedLobby)
	}

	lobby.Status = status

	err, _ := service.lobbyRepository.UpdateLobby(&lobby)

	if err != nil {
		logger.Error(err)
	}

	return err
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
		Options:  lobby.Options,
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
		Options:    entity.RaceOptions{},
		CreatedAt:  time.Now(),
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

func (service *lobbyService) Leave(ID uint64, userId uint64) (error, entity.Lobby) {
	logger.Info("LobbyService.Leave", map[string]interface{}{
		"lobbyId": ID,
		"userID":  userId,
	})

	lobby := service.lobbyRepository.FindLobbyById(ID)

	log.Println("LobbyService.Leave:", lobby.ID, ID, userId)

	if lobby.ID != 0 {
		log.Println("LobbyService.Leave: exists lobby", lobby.CountPlayers(), ID, userId)

		lobby.RemovePlayer(userId)

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
