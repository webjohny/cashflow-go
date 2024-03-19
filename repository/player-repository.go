package repository

import (
	"github.com/webjohny/cashflow-go/entity"
	"gorm.io/gorm"
)

type PlayerRepository interface {
	InsertPlayer(b *entity.Player) entity.Player
	UpdatePlayer(b *entity.Player) entity.Player
	AllByRaceId(raceId uint64) []entity.Player
	DeletePlayer(b *entity.Player)
	FindPlayerById(ID uint64) entity.Player
	FindPlayerByUsername(username string) entity.Player
	FindPlayerByUsernameAndRaceId(raceId uint64, username string) entity.Player
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

func (db *playerConnection) AllByRaceId(raceId uint64) []entity.Player {
	var players []entity.Player
	db.connection.Preload(PlayerTable).Where("race_id = ?", raceId).Find(&players)
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

func (db *playerConnection) FindPlayerById(ID uint64) entity.Player {
	var player entity.Player

	db.connection.Preload(PlayerTable).Find(&player, ID)

	return player
}

func (db *playerConnection) FindPlayerByUsername(username string) entity.Player {
	var player entity.Player

	db.connection.Preload(PlayerTable).Find(&player, "`username` = ?", username).Find(&player)

	return player
}

func (db *playerConnection) FindPlayerByUsernameAndRaceId(raceId uint64, username string) entity.Player {
	var player entity.Player

	db.connection.Preload(PlayerTable).Where("`username` = ? AND `race_id` = ?", username, raceId).Find(&player)

	return player
}
