package repository

import (
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"gorm.io/gorm"
)

type RequestRepository interface {
	Insert(b *entity.Request) (error, entity.Request)
	Update(b *entity.Request) (error, entity.Request)
	All() []entity.Request
}

type requestConnection struct {
	connection *gorm.DB
}

func NewRequestRepository(dbConn *gorm.DB) RequestRepository {
	return &requestConnection{
		connection: dbConn,
	}
}

func (db *requestConnection) Insert(b *entity.Request) (error, entity.Request) {
	result := db.connection.Save(&b)

	if result.Error != nil {
		logger.Error(result.Error, helper.JsonSerialize(b))

		return result.Error, entity.Request{}
	}

	db.connection.Find(&b)
	return nil, *b
}

func (db *requestConnection) All() []entity.Request {
	var requests []entity.Request
	db.connection.Find(&requests)
	return requests
}

func (db *requestConnection) Update(b *entity.Request) (error, entity.Request) {
	result := db.connection.Select("*").Updates(&b)

	if result.Error != nil {
		logger.Error(result.Error, helper.JsonSerialize(b))

		return result.Error, entity.Request{}
	}

	db.connection.Find(&b)
	return nil, *b
}
