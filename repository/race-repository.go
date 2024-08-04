package repository

import (
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"gorm.io/gorm"
)

type RaceRepository interface {
	InsertRace(b *entity.Race) (error, entity.Race)
	UpdateRace(b *entity.Race) (error, entity.Race)
	All() []entity.Race
	DeleteRace(b *entity.Race) error
	FindRaceById(ID uint64) entity.Race
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
		logger.Error(result.Error, helper.JsonSerialize(b))

		return result.Error, entity.Race{}
	}

	db.connection.Find(&b)
	return nil, *b
}

func (db *raceConnection) All() []entity.Race {
	var races []entity.Race
	db.connection.Find(&races)
	return races
}

func (db *raceConnection) UpdateRace(b *entity.Race) (error, entity.Race) {
	result := db.connection.Select("*").Updates(&b)

	if result.Error != nil {
		logger.Error(result.Error, helper.JsonSerialize(b))

		return result.Error, entity.Race{}
	}

	db.connection.Find(&b)
	return nil, *b
}

func (db *raceConnection) DeleteRace(b *entity.Race) error {
	result := db.connection.Delete(&b)

	if result.Error != nil {
		logger.Error(result.Error, helper.JsonSerialize(b))

		return result.Error
	}

	return nil
}

func (db *raceConnection) FindRaceById(ID uint64) entity.Race {
	var race entity.Race

	db.connection.Where("id = ?", ID).Find(&race)

	return race
}
