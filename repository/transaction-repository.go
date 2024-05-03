package repository

import (
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/request"
	"gorm.io/datatypes"
	"time"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	InsertTransaction(b *entity.Transaction) entity.Transaction
	UpdateTransaction(b *entity.Transaction) entity.Transaction
	GetPlayerTransactions(playerId uint64) []entity.Transaction
	GetRaceTransactions(raceId uint64) []entity.Transaction
	DeleteTransaction(b *entity.Transaction)
	FindTransactionByPlayerId(ID uint64) entity.Transaction
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

func (db *transactionConnection) InsertTransaction(b *entity.Transaction) entity.Transaction {
	b.CreatedAt = datatypes.Date(time.Now())
	db.connection.Save(&b)
	db.connection.Preload(TransactionsTable).Find(&b)
	return *b
}

func (db *transactionConnection) GetPlayerTransactions(playerId uint64) []entity.Transaction {
	var transactions []entity.Transaction
	db.connection.Preload(TransactionsTable).Where("player_id = ?", playerId).Find(&transactions)
	return transactions
}

func (db *transactionConnection) GetRaceTransactions(raceId uint64) []entity.Transaction {
	var transactions []entity.Transaction
	db.connection.Preload(TransactionsTable).Where("race_id = ?", raceId).Find(&transactions)
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
	db.connection.Save(&b)
	db.connection.Preload(TransactionsTable).Find(&b)
	return *b
}

func (db *transactionConnection) DeleteTransaction(b *entity.Transaction) {
	db.connection.Delete(&b)
}

func (db *transactionConnection) FindTransactionByPlayerId(ID uint64) entity.Transaction {
	var transaction entity.Transaction
	db.connection.Preload(TransactionsTable).Find(&transaction, ID)
	return transaction
}

func (db *transactionConnection) FindTransactionByRaceId(ID uint64) entity.Transaction {
	var transaction entity.Transaction
	db.connection.Preload(TransactionsTable).Find(&transaction, ID)
	return transaction
}
