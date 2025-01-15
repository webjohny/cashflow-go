package repository

import (
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"gorm.io/gorm"
)

type LobbyRepository interface {
	InsertLobby(b *entity.Lobby) (error, entity.Lobby)
	UpdateLobby(b *entity.Lobby) (error, entity.Lobby)
	All() []entity.Lobby
	DeleteLobby(b *entity.Lobby)
	CancelLobby(b *entity.Lobby)
	FindLobbyById(ID uint64) entity.Lobby
	FindLobbyByGameId(gameId uint64) entity.Lobby
}

const LobbyTable = "lobbies"

type lobbyConnection struct {
	connection *gorm.DB
}

func NewLobbyRepository(dbConn *gorm.DB) LobbyRepository {
	return &lobbyConnection{
		connection: dbConn,
	}
}

func (db *lobbyConnection) InsertLobby(b *entity.Lobby) (error, entity.Lobby) {
	result := db.connection.Save(&b)

	if result.Error != nil {
		logger.Error(result.Error, helper.JsonSerialize(b))

		return result.Error, entity.Lobby{}
	}

	db.connection.Find(&b)
	return nil, *b
}

func (db *lobbyConnection) All() []entity.Lobby {
	var lobbies []entity.Lobby
	db.connection.Find(&lobbies)
	return lobbies
}

func (db *lobbyConnection) UpdateLobby(b *entity.Lobby) (error, entity.Lobby) {
	result := db.connection.Select("*").Updates(&b)

	if result.Error != nil {
		logger.Error(result.Error, helper.JsonSerialize(b))

		return result.Error, entity.Lobby{}
	}

	db.connection.Find(&b)
	return nil, *b
}

func (db *lobbyConnection) DeleteLobby(b *entity.Lobby) {
	result := db.connection.Delete(&b)

	if result.Error != nil {
		logger.Error(result.Error, helper.JsonSerialize(b))
	}
}

func (db *lobbyConnection) CancelLobby(b *entity.Lobby) {
	b.Status = entity.LobbyStatus.Cancelled
	result := db.connection.Save(&b)

	if result.Error != nil {
		logger.Error(result.Error, helper.JsonSerialize(b))
	}
}

func (db *lobbyConnection) FindLobbyById(ID uint64) entity.Lobby {
	var lobby entity.Lobby

	db.connection.Where("id = ?", ID).Find(&lobby)

	return lobby
}

func (db *lobbyConnection) FindLobbyByGameId(gameId uint64) entity.Lobby {
	var lobby entity.Lobby

	db.connection.Where("game_id = ?", gameId).Find(&lobby)

	return lobby
}
