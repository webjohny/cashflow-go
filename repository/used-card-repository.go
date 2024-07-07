package repository

import (
	"github.com/webjohny/cashflow-go/entity"
	"gorm.io/gorm"
)

type UsedCardRepository interface {
	SetCard(raceID uint64, cardID string, family string)
	HasCard(raceID uint64, cardID string, family string) bool
	CountFamilyCards(raceID uint64, action string) int64
}

type usedCardConnection struct {
	connection *gorm.DB
}

func NewUsedCardRepository(dbConn *gorm.DB) UsedCardRepository {
	return &usedCardConnection{
		connection: dbConn,
	}
}

func (db *usedCardConnection) HasCard(raceID uint64, cardID string, action string) bool {
	var usedCard = entity.UsedCard{
		RaceID: raceID,
		CardID: cardID,
		Action: action,
	}
	db.connection.Find(&usedCard)

	return usedCard.ID > 0
}

func (db *usedCardConnection) SetCard(raceID uint64, cardID string, action string) {
	db.connection.Save(&entity.UsedCard{
		RaceID: raceID,
		CardID: cardID,
		Action: action,
	})
}

func (db *usedCardConnection) CountFamilyCards(raceID uint64, action string) int64 {
	var count int64

	db.connection.Model(&entity.UsedCard{}).
		Where("race_id = ? AND action = ?", raceID, action).
		Count(&count)

	return count
}
