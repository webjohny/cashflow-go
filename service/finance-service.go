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

func (service *financeService) AskMoney(raceId uint64, userId uint64, dto dto.AskMoneyBodyDto) (error, bool) {
	logger.Info("FinanceService.AskMoney", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"dto":    dto,
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

	if dto.Type == "" {
		dto.Type = entity.UserRequestTypes.Salary
	}

	//@toDo need to make checking if its not repeated

	if !race.Options.EnableManager {
		countPayDay := service.cardService.CheckPayDay(player)

		logger.Info("FinanceService.AskMoney: have no manager", map[string]interface{}{
			"playerId":        player.ID,
			"lastPosition":    player.LastPosition,
			"currentPosition": player.CurrentPosition,
			"playerCashFlow":  player.CalculateCashFlow(),
			"countPayDay":     countPayDay,
			"dto.Amount":      dto.Amount,
			"dto.Type":        dto.Type,
		})

		if player.LastPosition == 0 && player.CurrentPosition == 0 {
			if dto.Amount != (player.CalculateCashFlow() + player.Assets.Savings) {
				return errors.New(storage.ErrorWrongAmount), false
			}

			player.Assets.Savings = 0
		} else if race.CurrentCard.Type == "payday" || race.CurrentCard.Type == "cashFlowDay" {

			if dto.Amount != player.CalculateCashFlow() && dto.Amount != (player.CalculateCashFlow()*countPayDay) {
				return errors.New(storage.ErrorWrongAmount), false
			}
		} else if dto.Type == entity.UserRequestTypes.Salary {

			if countPayDay == 0 {
				return errors.New(storage.ErrorTransactionDeclined), false
			} else if dto.Amount != player.CalculateCashFlow() && dto.Amount != (player.CalculateCashFlow()*countPayDay) {
				return errors.New(storage.ErrorWrongAmount), false
			}
		} else if dto.Type == entity.UserRequestTypes.Baby && dto.Amount != 1000 {
			return errors.New(storage.ErrorWrongAmount), false
		}
	}

	var request entity.UserRequest

	request.Type = dto.Type
	request.CurrentCard = race.CurrentCard.ID
	request.RaceID = raceId
	request.UserID = userId
	request.Message = dto.Message
	request.Amount = dto.Amount
	request.Data = map[string]interface{}{
		"last":    player.LastPosition,
		"current": player.CurrentPosition,
	}

	if !race.Options.EnableManager {
		request.Approved = true
	}

	err, _ = service.userRequestRepository.Insert(&request)

	if err != nil {
		return err, false
	}

	if !race.Options.EnableManager {
		service.playerService.UpdateCash(&player, dto.Amount, "Запрос: "+dto.Message)

		return nil, true
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

	err, sender := service.playerService.GetPlayerByUserIdAndRaceId(raceId, userId)

	if err != nil {
		return err
	}
	if sender.ID == 0 {
		return errors.New(storage.ErrorUndefinedPlayer)
	}
	if sender.Cash < amount {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	if receiverUsername != "" {
		receiver := service.playerService.GetPlayerByUsernameAndRaceId(raceId, receiverUsername)

		if receiver.ID == 0 {
			return errors.New(storage.ErrorUndefinedReceiverPlayer)
		}

		service.playerService.UpdateCash(
			&sender,
			-amount,
			fmt.Sprintf(
				"Перевёл $%s игроку %s",
				strconv.Itoa(amount),
				helper.CamelToCapitalize(receiverUsername),
			),
		)

		service.playerService.UpdateCash(
			&receiver,
			amount,
			fmt.Sprintf(
				"%s перевёл Вам $%s",
				helper.CamelToCapitalize(sender.Username),
				strconv.Itoa(amount),
			),
		)

		go service.raceService.SetTransaction(raceId, sender, "", fmt.Sprintf(
			"%s перевёл $%s игроку %s",
			helper.CamelToCapitalize(sender.Username),
			strconv.Itoa(amount),
			helper.CamelToCapitalize(receiverUsername),
		))
	} else {
		service.playerService.UpdateCash(
			&sender,
			-amount,
			fmt.Sprintf(
				"Перевёл $%s банк",
				strconv.Itoa(amount),
			),
		)

		go service.raceService.SetTransaction(raceId, sender, "", fmt.Sprintf(
			"%s перевёл $%s в банк",
			helper.CamelToCapitalize(sender.Username),
			strconv.Itoa(amount),
		))
	}

	return nil
}

func (service *financeService) SendAssets(raceId uint64, userId uint64, dto dto.SendAssetsBodyDTO) error {
	logger.Info("FinanceService.SendAssets", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"dto":    dto,
	})

	err, sender := service.playerService.GetPlayerByUserIdAndRaceId(raceId, userId)

	if err != nil {
		return err
	}

	receiver := service.playerService.GetPlayerByUsernameAndRaceId(raceId, dto.Player)
	race := service.raceService.GetRaceByRaceId(raceId)

	if sender.ID == 0 {
		return errors.New(storage.ErrorUndefinedPlayer)
	}

	if receiver.ID == 0 {
		return errors.New(storage.ErrorUndefinedReceiverPlayer)
	}

	var transactionMessage string
	var transactionSenderMessage string
	var transactionReceiverMessage string

	if dto.Asset == "stock" {
		cardStocks := entity.CardStocks{}
		cardStocks.Fill(race.CurrentCard)
		err = service.playerService.TransferStocks(cardStocks, dto.AssetId, sender, receiver, dto.Amount)

		if err == nil {
			transactionSenderMessage = fmt.Sprintf(
				"Перевёл акции (%s шт.) игроку %s",
				strconv.Itoa(dto.Amount),
				helper.CamelToCapitalize(dto.Player),
			)
			transactionReceiverMessage = fmt.Sprintf(
				"%s перевёл Вам (%s шт.) акций",
				helper.CamelToCapitalize(dto.Player),
				strconv.Itoa(dto.Amount),
			)
			transactionMessage = fmt.Sprintf(
				"%s перевёл акции (%s шт.) игроку %s",
				helper.CamelToCapitalize(sender.Username),
				strconv.Itoa(dto.Amount),
				helper.CamelToCapitalize(dto.Player),
			)
		}
	} else if dto.Asset == "business" {
		err = service.playerService.TransferBusiness(dto.AssetId, sender, receiver, dto.Amount)

		if err == nil {
			transactionSenderMessage = fmt.Sprintf(
				"Перевёл акции (%s шт.) игроку %s",
				strconv.Itoa(dto.Amount),
				helper.CamelToCapitalize(dto.Player),
			)
			transactionReceiverMessage = fmt.Sprintf(
				"%s перевёл Вам (%s шт.) акций",
				helper.CamelToCapitalize(dto.Player),
				strconv.Itoa(dto.Amount),
			)
			transactionMessage = fmt.Sprintf(
				"%s перевёл акции (%s шт.) игроку %s",
				helper.CamelToCapitalize(sender.Username),
				strconv.Itoa(dto.Amount),
				helper.CamelToCapitalize(dto.Player),
			)
		}
	}

	if err == nil {
		go service.playerService.SetTransaction(sender.ID, sender.Cash, sender.Cash, dto.Amount, transactionSenderMessage)
		go service.playerService.SetTransaction(receiver.ID, receiver.Cash, receiver.Cash, dto.Amount, transactionReceiverMessage)
		go service.raceService.SetTransaction(raceId, sender, "", transactionMessage)
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

	err = service.playerService.PayLoan(player, "bankLoan", amount)

	if err == nil {
		go service.raceService.SetTransaction(raceId, player, "", fmt.Sprintf(
			"%s отдал кредит $%s",
			helper.CamelToCapitalize(player.Username),
			strconv.Itoa(amount),
		))
	}

	return err
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

	if result != amount {
		return errors.New(storage.ErrorWrongAmount)
	}

	service.playerService.UpdateCash(&player, -result, card.Heading)

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		race.CurrentCard = entity.Card{}
		err, _ = service.raceService.UpdateRace(&race)
	}

	return err
}

func (service *financeService) TakeLoan(raceId uint64, userId uint64, amount int) error {
	logger.Info("FinanceService.TakeLoan", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"amount": amount,
	})

	err, player := service.playerService.GetPlayerByUserIdAndRaceId(raceId, userId)

	if player.ID == 0 {
		return errors.New(storage.ErrorUndefinedPlayer)
	}

	if amount%1000 != 0 {
		return errors.New(storage.ErrorWrongAmountForTakingLoan)
	}

	err = service.playerService.TakeLoan(player, amount)

	if err == nil {
		go service.raceService.SetTransaction(raceId, player, "", fmt.Sprintf(
			"%s взял в кредит $%s",
			helper.CamelToCapitalize(player.Username),
			strconv.Itoa(amount),
		))
	}

	return err
}
