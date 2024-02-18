package repository

import (
	"github.com/webjohny/cashflow-go/entity"
	"gorm.io/gorm"
)

type RaceRepository interface {
	InsertRace(b *entity.Race) entity.Race
	UpdateRace(b *entity.Race) entity.Race
	All(idUser string) []entity.Race
	DeleteRace(b *entity.Race)
	FindRaceById(ID uint64, IsBigRace bool) *entity.Race
}

const RaceTable = "races"

type raceConnection struct {
	connection *gorm.DB
}

func NewRaceRepository(dbConn *gorm.DB) RaceRepository {
	return &raceConnection{
		connection: dbConn,
	}
}

func (db *raceConnection) InsertRace(b *entity.Race) entity.Race {
	db.connection.Save(&b)
	db.connection.Preload(RaceTable).Find(&b)
	return *b
}

func (db *raceConnection) All(idUser string) []entity.Race {
	var races []entity.Race
	db.connection.Preload(RaceTable).Where("user_id = ?", idUser).Find(&races)
	return races
}

func (db *raceConnection) UpdateRace(b *entity.Race) entity.Race {
	db.connection.Save(&b)
	db.connection.Preload(RaceTable).Find(&b)
	return *b
}

func (db *raceConnection) DeleteRace(b *entity.Race) {
	db.connection.Delete(&b)
}

func (db *raceConnection) FindRaceById(ID uint64, IsBigRace bool) *entity.Race {
	var race *entity.Race

	if !IsBigRace {
		db.connection.Preload(RaceTable).Find(&race, ID)
	} else {
		ParentID := ID
		db.connection.Preload(RaceTable).Find(&race, ParentID).Find(&race, IsBigRace)
	}
	return race
}
