package service

import (
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/storage"
	"strconv"
)

type FinanceService interface {
	SendMoney(raceId uint64, userId uint64, amount int, player string) error
	SendAssets(raceId uint64, userId uint64, dto dto.SendAssetsBodyDTO) error
	PayLoan(raceId uint64, userId uint64, amount int) error
	PayTax(raceId uint64, userId uint64, amount int) error
	TakeLoan(raceId uint64, userId uint64, amount int) error
	AskMoney(raceId uint64, userId uint64, dto dto.AskMoneyBodyDto) (error, bool)
}

type financeService struct {
	userRequestRepository repository.UserRequestRepository
	cardService           CardService
	raceService           RaceService
	playerService         PlayerService
}

func NewFinanceService(userRequestRepository repository.UserRequestRepository, cardService CardService, raceService RaceService, playerService PlayerService) FinanceService {
	return &financeService{
		userRequestRepository: userRequestRepository,
		cardService:           cardService,
		raceService:           raceService,
		playerService:         playerService,
	}
}

func (service *financeService) AskMoney(raceId uint64, userId uint64, data dto.AskMoneyBodyDto) (error, bool) {
	logger.Info("FinanceService.AskMoney", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"dto":    data,
	})

	err, race, player := service.raceService.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err, false
	}

	if race.ID == 0 {
		return errors.New(storage.ErrorUndefinedGame), false
	}
	if player.ID == 0 {
		return errors.New(storage.ErrorUndefinedPlayer), false
	}

	if data.Type == "" {
		data.Type = entity.UserRequestTypes.Salary
	}

	//@toDo need to make checking if its not repeated

	cardType := entity.TransactionCardType.ReceiveMoney

	if !race.Options.EnableManager {
		countPayDay := service.cardService.CheckPayDay(player)

		logger.Info("FinanceService.AskMoney: have no manager", map[string]interface{}{
			"playerId":        player.ID,
			"lastPosition":    player.LastPosition,
			"currentPosition": player.CurrentPosition,
			"playerCashFlow":  player.CalculateCashFlow(),
			"countPayDay":     countPayDay,
			"dto.Amount":      data.Amount,
			"dto.Type":        data.Type,
		})

		if player.LastPosition == 0 && player.CurrentPosition == 0 {
			if data.Amount != (player.CalculateCashFlow() + player.Assets.Savings) {
				return errors.New(storage.ErrorWrongAmount), false
			}

			player.Assets.Savings = 0
			cardType = entity.TransactionCardType.StartMoney
			race.CurrentCard.ID = cardType
		} else if race.CurrentCard.Type == entity.TransactionCardType.Payday ||
			race.CurrentCard.Type == entity.TransactionCardType.CashFlowDay {
			cardType = race.CurrentCard.Type

			if data.Amount != player.CalculateCashFlow() && data.Amount != (player.CalculateCashFlow()*countPayDay) {
				return errors.New(storage.ErrorWrongAmount), false
			}
		} else if data.Type == entity.UserRequestTypes.Salary {

			cardType = entity.TransactionCardType.Payday
			if countPayDay == 0 {
				return errors.New(storage.ErrorTransactionDeclined), false
			} else if data.Amount != player.CalculateCashFlow() && data.Amount != (player.CalculateCashFlow()*countPayDay) {
				return errors.New(storage.ErrorWrongAmount), false
			}
		} else if data.Type == entity.UserRequestTypes.Baby && data.Amount != 1000 {
			return errors.New(storage.ErrorWrongAmount), false
		}
	}

	updatedCash := player.Cash - data.Amount

	if race.CurrentCard.ID == "" {
		race.CurrentCard.ID = helper.Uuid(cardType)
	}

	transaction := dto.TransactionDTO{
		CardID:          race.CurrentCard.ID,
		CardType:        cardType,
		Details:         "",
		PlayerID:        player.ID,
		RaceID:          player.RaceID,
		CurrentCashFlow: player.CalculateCashFlow(),
		CurrentCash:     player.Cash,
		UpdatedCash:     updatedCash,
	}

	if trx := service.playerService.GetTransaction(transaction); trx.ID != 0 {
		return errors.New(storage.ErrorTransactionAlreadyExists), false
	}

	if race.Options.EnableManager {
		var request entity.UserRequest

		request.Type = data.Type
		request.CurrentCard = race.CurrentCard.ID
		request.RaceID = raceId
		request.UserID = userId
		request.Message = data.Message
		request.Amount = data.Amount
		request.Data = map[string]interface{}{
			"last":    player.LastPosition,
			"current": player.CurrentPosition,
		}

		err, _ = service.userRequestRepository.Insert(&request)

		if err != nil {
			return err, false
		}
	} else {
		err = service.playerService.UpdateCash(&player, data.Amount, &transaction)

		return err, true
	}

	return nil, false
}

func (service *financeService) SendMoney(raceId uint64, userId uint64, amount int, receiverUsername string) error {
	logger.Info("FinanceService.SendMoney", map[string]interface{}{
		"raceId":           raceId,
		"userId":           userId,
		"amount":           amount,
		"receiverUsername": receiverUsername,
	})

	err, race, sender := service.raceService.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}
	if sender.ID == 0 {
		return errors.New(storage.ErrorUndefinedPlayer)
	}
	if sender.Cash < amount {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	if receiverUsername == "" {
		err = service.playerService.UpdateCash(
			&sender,
			-amount,
			&dto.TransactionDTO{
				CardType: entity.TransactionCardType.SendMoneyToBank,
				Details: fmt.Sprintf(
					"Вы перевели $%s банк",
					strconv.Itoa(amount),
				),
			},
		)

		return err
	}

	receiver := service.playerService.GetPlayerByUsernameAndRaceId(raceId, receiverUsername)

	if receiver.ID == 0 {
		return errors.New(storage.ErrorUndefinedReceiverPlayer)
	}

	var transactionData = dto.TransactionDTO{
		CardType: entity.TransactionCardType.SendMoney,
		Details: fmt.Sprintf(
			"Перевёл $%s игроку %s (#%s)",
			strconv.Itoa(amount),
			helper.CamelToCapitalize(receiverUsername),
			receiver.GetStringID(),
		),
	}

	if race.CurrentCard.ID != "" {
		cardID := fmt.Sprintf("%s-%s", race.CurrentCard.ID, strconv.Itoa(helper.Random(1000)))
		transactionData.CardID = cardID
	}

	err = service.playerService.UpdateCash(
		&sender,
		-amount,
		&transactionData,
	)

	if err != nil {
		return err
	}

	transactionData.CardType = entity.TransactionCardType.ReceiveMoney
	transactionData.Details = fmt.Sprintf(
		"%s (#%s) перевёл Вам $%s",
		helper.CamelToCapitalize(sender.Username),
		sender.GetStringID(),
		strconv.Itoa(amount),
	)

	receiver.SetNotification(transactionData.Details, entity.NotificationTypes.Success)

	return service.playerService.UpdateCash(
		&receiver,
		amount,
		&transactionData,
	)
}

func (service *financeService) SendAssets(raceId uint64, userId uint64, data dto.SendAssetsBodyDTO) error {
	logger.Info("FinanceService.SendAssets", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"dto":    data,
	})

	err, race, sender := service.raceService.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	err, receiver := service.playerService.GetPlayerByPlayerIdAndRaceId(raceId, uint64(data.Player))

	if err != nil {
		return errors.New(storage.ErrorUndefinedReceiverPlayer)
	}

	if sender.ID == 0 {
		return errors.New(storage.ErrorUndefinedPlayer)
	}

	if receiver.ID == 0 {
		return errors.New(storage.ErrorUndefinedReceiverPlayer)
	}

	var transactionSenderMessage string
	var transactionReceiverMessage string

	if data.Asset == "stock" {
		transactionSenderMessage = fmt.Sprintf(
			"Перевёл акции (%s шт.) игроку %s",
			strconv.Itoa(data.Amount),
			helper.CamelToCapitalize(receiver.Username),
		)
		transactionReceiverMessage = fmt.Sprintf(
			"%s (#%s) перевёл Вам (%s шт.) акций",
			helper.CamelToCapitalize(sender.Username),
			sender.GetStringID(),
			strconv.Itoa(data.Amount),
		)
		receiver.SetNotification(transactionReceiverMessage, entity.NotificationTypes.Success)

		err = service.playerService.TransferStocks(data.AssetId, sender, receiver, data.Amount)
	} else if data.Asset == "business" {
		transactionSenderMessage = fmt.Sprintf(
			"Перевёл акции (%s шт.) игроку %s",
			strconv.Itoa(data.Amount),
			helper.CamelToCapitalize(receiver.Username),
		)
		transactionReceiverMessage = fmt.Sprintf(
			"%s (#%s) перевёл Вам (%s шт.) акций",
			helper.CamelToCapitalize(sender.Username),
			sender.GetStringID(),
			strconv.Itoa(data.Amount),
		)
		receiver.SetNotification(transactionReceiverMessage, entity.NotificationTypes.Success)

		err = service.playerService.TransferBusiness(data.AssetId, sender, receiver, data.Amount)
	}

	if err == nil {
		transactionData := dto.TransactionDTO{
			CardType:    entity.TransactionCardType.SendAssets,
			Details:     transactionSenderMessage,
			CurrentCash: sender.Cash,
			Amount:      data.Amount,
		}

		if race.CurrentCard.ID != "" {
			cardID := fmt.Sprintf("%s-%s", race.CurrentCard.ID, strconv.Itoa(helper.Random(1000)))
			transactionData.CardID = cardID
		}

		err = service.playerService.SetTransaction(sender, transactionData)

		if err != nil {
			return err
		}

		transactionData.CardType = entity.TransactionCardType.ReceiveAssets
		transactionData.Details = transactionReceiverMessage
		transactionData.CurrentCash = receiver.Cash

		err = service.playerService.SetTransaction(receiver, transactionData)
	}

	return err
}

func (service *financeService) PayLoan(raceId uint64, userId uint64, amount int) error {
	logger.Info("FinanceService.PayLoan", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"amount": amount,
	})

	err, player := service.playerService.GetPlayerByUserIdAndRaceId(raceId, userId)

	if err != nil {
		return err
	}

	if player.ID == 0 {
		return errors.New(storage.ErrorUndefinedPlayer)
	}

	if amount%1000 != 0 {
		return errors.New(storage.ErrorWrongAmountForPayingLoan)
	}

	return service.playerService.PayLoan(player, "bankLoan", amount)
}

func (service *financeService) PayTax(raceId uint64, userId uint64, amount int) error {
	logger.Info("FinanceService.PayTax", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"amount": amount,
	})

	err, race, player := service.raceService.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	card := entity.CardPayTax{}
	card.Fill(race.CurrentCard)

	var result int

	if card.Percent < 100 {
		result = int((float32(player.Cash) / 100) * float32(card.Percent))
	} else {
		result = player.Cash
	}

	if result > amount {
		return errors.New(storage.ErrorWrongAmount)
	}

	err = service.playerService.UpdateCash(&player, -amount, &dto.TransactionDTO{
		CardID:   card.ID,
		CardType: entity.TransactionCardType.PayTax,
		Details:  card.Heading,
	})

	race.Respond(player.ID, race.CurrentPlayer.ID)
	race.CurrentCard = entity.Card{}
	err, _ = service.raceService.UpdateRace(&race)

	return err
}

func (service *financeService) TakeLoan(raceId uint64, userId uint64, amount int) error {
	logger.Info("FinanceService.TakeLoan", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"amount": amount,
	})

	err, race, player := service.raceService.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	if player.ID == 0 {
		return errors.New(storage.ErrorUndefinedPlayer)
	}

	transaction := dto.TransactionDTO{
		CardID:      race.CurrentCard.ID,
		CardType:    entity.TransactionCardType.TakeLoan,
		Details:     "",
		PlayerID:    player.ID,
		RaceID:      player.RaceID,
		CurrentCash: player.Cash,
	}

	if trx := service.playerService.GetTransaction(transaction); trx.ID != 0 {
		return errors.New(storage.ErrorTransactionAlreadyExists)
	}

	if race.Options.EnableManager {
		var request entity.UserRequest

		request.Type = entity.TransactionCardType.TakeLoan
		request.CurrentCard = race.CurrentCard.ID
		request.RaceID = raceId
		request.UserID = userId
		request.Message = entity.TransactionCardType.TakeLoan
		request.Amount = amount
		request.Data = map[string]interface{}{
			"last":    player.LastPosition,
			"current": player.CurrentPosition,
		}

		err, _ = service.userRequestRepository.Insert(&request)

		if err != nil {
			return err
		}
	} else {
		if amount%1000 != 0 {
			return errors.New(storage.ErrorWrongAmountForTakingLoan)
		}

		return service.playerService.TakeLoan(player, amount)
	}

	return err
}
