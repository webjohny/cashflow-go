package repository

import (
	"github.com/webjohny/cashflow-go/entity"
	"gorm.io/gorm"
)

type LobbyRepository interface {
	InsertLobby(b *entity.Lobby) entity.Lobby
	UpdateLobby(b *entity.Lobby) entity.Lobby
	All() []entity.Lobby
	DeleteLobby(b *entity.Lobby)
	FindLobbyById(ID uint64) *entity.Lobby
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

func (db *lobbyConnection) InsertLobby(b *entity.Lobby) entity.Lobby {
	db.connection.Save(&b)
	db.connection.Preload(LobbyTable).Find(&b)
	return *b
}

func (db *lobbyConnection) All() []entity.Lobby {
	var lobbies []entity.Lobby
	db.connection.Preload(LobbyTable).Find(&lobbies)
	return lobbies
}

func (db *lobbyConnection) UpdateLobby(b *entity.Lobby) entity.Lobby {
	db.connection.Save(&b)
	db.connection.Preload(LobbyTable).Find(&b)
	return *b
}

func (db *lobbyConnection) DeleteLobby(b *entity.Lobby) {
	db.connection.Delete(&b)
}

func (db *lobbyConnection) FindLobbyById(ID uint64) *entity.Lobby {
	var lobby *entity.Lobby

	db.connection.Preload(LobbyTable).Find(&lobby, ID)

	return lobby
}
