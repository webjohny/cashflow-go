package service

import (
	"github.com/mashingan/smapping"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/repository"
	"log"
)

type TransactionService interface {
	InsertPlayerTransaction(b dto.TransactionCreatePlayerDTO) entity.Transaction
	InsertRaceTransaction(b dto.TransactionCreateRaceDTO) entity.Transaction
	UpdateTransaction(b dto.TransactionUpdateDTO) entity.Transaction
	Delete(b entity.Transaction)
	GetPlayerTransactions(playerId uint64) []entity.Transaction
	GetRaceTransactions(raceId uint64) []entity.Transaction
	GetRaceLogs(raceId uint64) []entity.RaceLog
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
}

func NewTransactionService(transactionRepo repository.TransactionRepository) TransactionService {
	return &transactionService{
		transactionRepository: transactionRepo,
	}
}

func (service *transactionService) InsertPlayerTransaction(b dto.TransactionCreatePlayerDTO) entity.Transaction {
	trx := entity.Transaction{}
	trx.PlayerID = &b.PlayerID
	trx.TransactionType = entity.TransactionType.PLAYER
	trx.Details = b.Details
	trx.Data = &entity.TransactionData{
		CurrentCash: &b.CurrentCash,
		Cash:        &b.Cash,
		Amount:      &b.Amount,
	}
	res := service.transactionRepository.InsertTransaction(&trx)
	return res
}

func (service *transactionService) InsertRaceTransaction(b dto.TransactionCreateRaceDTO) entity.Transaction {
	trx := entity.Transaction{}
	trx.RaceID = &b.RaceID
	trx.PlayerID = &b.PlayerID
	trx.TransactionType = entity.TransactionType.RACE
	trx.Details = b.Details
	trx.Data = &entity.TransactionData{
		TxType:   &b.TxType,
		Color:    &b.Color,
		Username: &b.Username,
	}
	res := service.transactionRepository.InsertTransaction(&trx)
	return res
}

func (service *transactionService) GetPlayerTransactions(playerId uint64) []entity.Transaction {
	return service.transactionRepository.GetPlayerTransactions(playerId)
}

func (service *transactionService) GetRaceTransactions(raceId uint64) []entity.Transaction {
	return service.transactionRepository.GetRaceTransactions(raceId)
}

func (service *transactionService) GetRaceLogs(raceId uint64) []entity.RaceLog {
	logs := make([]entity.RaceLog, 0)
	transactions := service.GetRaceTransactions(raceId)

	for _, transaction := range transactions {
		logs = append(logs, entity.RaceLog{
			Username: *transaction.Data.Username,
			Color:    *transaction.Data.Color,
			Message:  transaction.Details,
		})
	}

	return logs
}

func (service *transactionService) UpdateTransaction(b dto.TransactionUpdateDTO) entity.Transaction {
	transaction := entity.Transaction{}
	err := smapping.FillStruct(&transaction, smapping.MapFields(&b))
	if err != nil {
		log.Fatalf("Failed map %v : ", err)
	}
	res := service.transactionRepository.UpdateTransaction(&transaction)
	return res
}

func (service *transactionService) Delete(b entity.Transaction) {
	service.transactionRepository.DeleteTransaction(&b)
}

//func (service *transactionService) IsAllowedToEdit(userID string, transactionID uint64) bool {
//	b := service.transactionRepository.FindTransactionById(transactionID)
//	id := fmt.Sprintf("%v", b.UserID)
//	return userID == id
//}
