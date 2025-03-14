package service

import (
	"github.com/mashingan/smapping"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/repository"
	"log"
)

type TransactionService interface {
	InsertPlayerTransaction(b dto.TransactionCreatePlayerDTO) error
	InsertTransaction(b dto.TransactionDTO) error
	InsertRaceTransaction(b dto.TransactionCreateRaceDTO) error
	UpdateTransaction(b dto.TransactionUpdateDTO) entity.Transaction
	Delete(b entity.Transaction)
	GetTransaction(data dto.TransactionDTO) entity.Transaction
	GetRaceTransaction(player entity.Player, data dto.TransactionCardDTO) entity.Transaction
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

func (service *transactionService) InsertPlayerTransaction(b dto.TransactionCreatePlayerDTO) error {
	trx := entity.Transaction{}
	trx.PlayerID = &b.PlayerID
	trx.TransactionType = entity.TransactionType.PLAYER
	trx.Details = b.Details
	trx.Data = &entity.TransactionData{
		CurrentCash: b.CurrentCash,
		UpdatedCash: b.Cash,
		Amount:      b.Amount,
	}
	return service.transactionRepository.InsertTransaction(&trx)
}

func (service *transactionService) InsertRaceTransaction(b dto.TransactionCreateRaceDTO) error {
	trx := entity.Transaction{}
	trx.RaceID = &b.RaceID
	trx.CardID = b.CardID
	trx.CardType = b.CardType
	trx.PlayerID = &b.PlayerID
	trx.TransactionType = entity.TransactionType.RACE
	trx.Details = b.Details
	trx.Data = &entity.TransactionData{
		Color:    b.Color,
		Username: b.Username,
	}
	return service.transactionRepository.InsertTransaction(&trx)
}

func (service *transactionService) InsertTransaction(b dto.TransactionDTO) error {
	trx := entity.Transaction{}
	trx.RaceID = &b.RaceID
	trx.CardID = b.CardID
	trx.CardType = b.CardType
	trx.PlayerID = &b.PlayerID
	trx.TransactionType = entity.TransactionType.PLAYER
	trx.Details = b.Details
	trx.Data = &entity.TransactionData{
		Color:           b.Color,
		Username:        b.Username,
		CurrentCash:     b.CurrentCash,
		UpdatedCash:     b.UpdatedCash,
		CurrentCashFlow: b.CurrentCashFlow,
		UpdatedCashFlow: b.UpdatedCashFlow,
		Amount:          b.Amount,
	}
	return service.transactionRepository.InsertTransaction(&trx)
}

func (service *transactionService) GetRaceTransaction(player entity.Player, data dto.TransactionCardDTO) entity.Transaction {
	return service.transactionRepository.FindRaceTransaction(player, data)
}

func (service *transactionService) GetTransaction(data dto.TransactionDTO) entity.Transaction {
	return service.transactionRepository.FindTransaction(data)
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
			Username:        transaction.Data.Username,
			PlayerId:        int(*transaction.PlayerID),
			Color:           transaction.Data.Color,
			CurrentCash:     transaction.Data.CurrentCash,
			UpdatedCash:     transaction.Data.UpdatedCash,
			CurrentCashFlow: transaction.Data.CurrentCashFlow,
			UpdatedCashFlow: transaction.Data.UpdatedCashFlow,
			Message:         transaction.Details,
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
