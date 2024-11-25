package repository

import (
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/request"
	"gorm.io/datatypes"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	InsertTransaction(b *entity.Transaction) error
	UpdateTransaction(b *entity.Transaction) entity.Transaction
	GetPlayerTransactions(playerId uint64) []entity.Transaction
	GetRaceTransactions(raceId uint64) []entity.Transaction
	DeleteTransaction(b *entity.Transaction)
	FindTransactionByPlayerId(ID uint64) entity.Transaction
	FindRaceTransaction(player entity.Player, data dto.TransactionCardDTO) entity.Transaction
}

const TransactionsTable = "transactions"

type transactionConnection struct {
	connection *gorm.DB
}

func NewTransactionRepository(dbConn *gorm.DB) TransactionRepository {
	return &transactionConnection{
		connection: dbConn,
	}
}

func (db *transactionConnection) InsertTransaction(b *entity.Transaction) error {
	b.CreatedAt = datatypes.Date(time.Now())
	result := db.connection.Save(&b)

	if result.Error != nil {
		logger.Error(result.Error, helper.JsonSerialize(b))

		return result.Error
	}

	db.connection.Find(&b)
	return nil
}

func (db *transactionConnection) GetPlayerTransactions(playerId uint64) []entity.Transaction {
	var transactions []entity.Transaction
	db.connection.Model(&entity.Transaction{}).Where("player_id = ?", playerId).Where("transaction_type = ?", entity.TransactionType.PLAYER).Scan(&transactions)
	return transactions
}

func (db *transactionConnection) GetRaceTransactions(raceId uint64) []entity.Transaction {
	var transactions []entity.Transaction
	db.connection.Model(&entity.Transaction{}).Where("race_id = ?", raceId).Where("transaction_type = ?", entity.TransactionType.RACE).Scan(&transactions)
	return transactions
}

func (db *transactionConnection) TransactionReport(idUser string) request.TransactionReport {
	var result request.TransactionReport

	db.connection.Model(&entity.Transaction{}).
		Select("SUM(CASE WHEN transaction_type = '1' THEN transaction_value ELSE 0 END) AS transaction_in, "+
			"SUM(CASE WHEN transaction_type = '2' THEN transaction_value ELSE 0 END) AS transaction_out").
		Where("user_id = ?", idUser).
		Scan(&result)

	return result
}

func (db *transactionConnection) UpdateTransaction(b *entity.Transaction) entity.Transaction {
	result := db.connection.Select("*").Updates(&b)

	if result.Error != nil {
		logger.Error(result.Error, helper.JsonSerialize(b))

		return entity.Transaction{}
	}

	db.connection.Find(&b)
	return *b
}

func (db *transactionConnection) DeleteTransaction(b *entity.Transaction) {
	result := db.connection.Delete(&b)

	if result.Error != nil {
		logger.Error(result.Error, helper.JsonSerialize(b))
	}
}

func (db *transactionConnection) FindRaceTransaction(player entity.Player, data dto.TransactionCardDTO) entity.Transaction {
	var transaction entity.Transaction

	db.connection.Model(TransactionsTable).
		Where("race_id", player.RaceID).
		Where("player_id", player.ID).
		Where("details", data.Details).
		Where("card_type", data.CardType).
		Where("card_id", data.CardID).
		Scan(&transaction)

	return transaction
}

func (db *transactionConnection) FindTransactionByPlayerId(ID uint64) entity.Transaction {
	var transaction entity.Transaction
	db.connection.Model(TransactionsTable).Find(&transaction, ID)
	return transaction
}

func (db *transactionConnection) FindTransactionByRaceId(ID uint64) entity.Transaction {
	var transaction entity.Transaction
	db.connection.Model(TransactionsTable).Find(&transaction, ID)
	return transaction
}
