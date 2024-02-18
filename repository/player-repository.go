package repository

import (
	"github.com/webjohny/cashflow-go/entity"
	"gorm.io/gorm"
)

type PlayerRepository interface {
	InsertPlayer(b *entity.Player) entity.Player
	UpdatePlayer(b *entity.Player) entity.Player
	All(idUser string) []entity.Player
	DeletePlayer(b *entity.Player)
	FindPlayerById(ID uint64) *entity.Player
	FindPlayerByUsername(username string) *entity.Player
}

const PlayerTable = "players"

type playerConnection struct {
	connection *gorm.DB
}

func NewPlayerRepository(dbConn *gorm.DB) PlayerRepository {
	return &playerConnection{
		connection: dbConn,
	}
}

func (db *playerConnection) InsertPlayer(b *entity.Player) entity.Player {
	db.connection.Save(&b)
	db.connection.Preload(PlayerTable).Find(&b)
	return *b
}

func (db *playerConnection) All(idUser string) []entity.Player {
	var players []entity.Player
	db.connection.Preload(PlayerTable).Where("user_id = ?", idUser).Find(&players)
	return players
}

func (db *playerConnection) UpdatePlayer(b *entity.Player) entity.Player {
	db.connection.Save(&b)
	db.connection.Preload(PlayerTable).Find(&b)
	return *b
}

func (db *playerConnection) DeletePlayer(b *entity.Player) {
	db.connection.Delete(&b)
}

func (db *playerConnection) FindPlayerById(ID uint64) *entity.Player {
	var player *entity.Player

	db.connection.Preload(PlayerTable).Find(&player, ID)

	return player
}

func (db *playerConnection) FindPlayerByUsername(Username string) *entity.Player {
	var player *entity.Player

	db.connection.Preload(PlayerTable).Find(&player, Username)

	return player
}
