package repository

import (
	"github.com/webjohny/cashflow-go/entity"
	"gorm.io/gorm"
)

type RaceRepository interface {
	InsertRace(b *entity.Race) (error, entity.Race)
	UpdateRace(b *entity.Race) (error, entity.Race)
	All() []entity.Race
	DeleteRace(b *entity.Race) error
	FindRaceById(ID uint64, IsBigRace bool) entity.Race
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

func (db *raceConnection) InsertRace(b *entity.Race) (error, entity.Race) {
	result := db.connection.Save(&b)

	if result.Error != nil {
		return result.Error, entity.Race{}
	}

	db.connection.Preload(RaceTable).Find(&b)
	return nil, *b
}

func (db *raceConnection) All() []entity.Race {
	var races []entity.Race
	db.connection.Preload(RaceTable).Find(&races)
	return races
}

func (db *raceConnection) UpdateRace(b *entity.Race) (error, entity.Race) {
	result := db.connection.Save(&b)

	if result.Error != nil {
		return result.Error, entity.Race{}
	}

	db.connection.Preload(RaceTable).Find(&b)
	return nil, *b
}

func (db *raceConnection) DeleteRace(b *entity.Race) error {
	result := db.connection.Delete(&b)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *raceConnection) FindRaceById(ID uint64, isBigRace bool) entity.Race {
	var race entity.Race

	if !isBigRace {
		db.connection.Preload(RaceTable).Find(&race, ID)
	} else {
		db.connection.Preload(RaceTable).Where("parent_id = ? AND is_big_race = ?", ID, isBigRace).Find(&race)
	}

	return race
}
