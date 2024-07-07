package repository_mocks

import "github.com/webjohny/cashflow-go/entity"

type MockTransactionRepository struct {
	InsertTransactionFunc func(b *entity.Transaction) entity.Transaction
	UpdateTransactionFunc func(b *entity.Transaction) entity.Transaction
	GetPlayerTransactionsFunc func(playerId uint64) []entity.Transaction
	GetRaceTransactionsFunc func(raceId uint64) []entity.Transaction
	DeleteTransactionFunc func(b *entity.Transaction)
	FindTransactionByPlayerIdFunc func(ID uint64) entity.Transaction
}

func (m MockTransactionRepository) InsertTransaction(b *entity.Transaction) entity.Transaction {
	//TODO implement me
	panic("implement me")
}

func (m MockTransactionRepository) UpdateTransaction(b *entity.Transaction) entity.Transaction {
	//TODO implement me
	panic("implement me")
}

func (m MockTransactionRepository) GetPlayerTransactions(playerId uint64) []entity.Transaction {
	//TODO implement me
	panic("implement me")
}

func (m MockTransactionRepository) GetRaceTransactions(raceId uint64) []entity.Transaction {
	//TODO implement me
	panic("implement me")
}

func (m MockTransactionRepository) DeleteTransaction(b *entity.Transaction) {
	//TODO implement me
	panic("implement me")
}

func (m MockTransactionRepository) FindTransactionByPlayerId(ID uint64) entity.Transaction {
	//TODO implement me
	panic("implement me")
}

