package repository

import (
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"gorm.io/gorm"
)

type UserRequestRepository interface {
	Insert(b *entity.UserRequest) (error, entity.UserRequest)
	Update(b *entity.UserRequest) (error, entity.UserRequest)
	All(where string, params []interface{}) []entity.UserRequest
	FindOneById(userRequestId uint64) entity.UserRequest
}

type userRequestConnection struct {
	connection *gorm.DB
}

func NewUserRequestRepository(dbConn *gorm.DB) UserRequestRepository {
	return &userRequestConnection{
		connection: dbConn,
	}
}

func (db *userRequestConnection) Insert(b *entity.UserRequest) (error, entity.UserRequest) {
	result := db.connection.Save(&b)

	if result.Error != nil {
		logger.Error(result.Error, helper.JsonSerialize(b))

		return result.Error, entity.UserRequest{}
	}

	db.connection.Find(&b)
	return nil, *b
}

func (db *userRequestConnection) FindOneById(userRequestId uint64) entity.UserRequest {
	var request entity.UserRequest

	db.connection.Find(&request, userRequestId)

	return request
}

func (db *userRequestConnection) All(where string, params []interface{}) []entity.UserRequest {
	var requests []entity.UserRequest

	if where != "" {
		db.connection.Preload("User").Where(where, params...).Find(&requests)
	} else {
		db.connection.Find(&requests)
	}

	return requests
}

func (db *userRequestConnection) Update(b *entity.UserRequest) (error, entity.UserRequest) {
	result := db.connection.Select("*").Updates(&b)

	if result.Error != nil {
		logger.Error(result.Error, helper.JsonSerialize(b))

		return result.Error, entity.UserRequest{}
	}

	db.connection.Find(&b)

	return nil, *b
}
