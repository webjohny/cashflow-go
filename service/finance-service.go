package service

import (
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/storage"
	"strconv"
)

type FinanceService interface {
	SendMoney(raceId uint64, username string, amount int, player string) error
	SendAssets(raceId uint64, username string, dto dto.SendAssetsBodyDTO) error
	PayLoan(raceId uint64, username string, amount int) error
	TakeLoan(raceId uint64, username string, amount int) error
}

type financeService struct {
	raceService   RaceService
	playerService PlayerService
}

func NewFinanceService(raceService RaceService, playerService PlayerService) FinanceService {
	return &financeService{
		raceService:   raceService,
		playerService: playerService,
	}
}

func (service *financeService) SendMoney(raceId uint64, username string, amount int, receiverUsername string) error {
	logger.Info("FinanceService.SendMoney", map[string]interface{}{
		"raceId":           raceId,
		"username":         username,
		"amount":           amount,
		"receiverUsername": receiverUsername,
	})

	sender := service.playerService.GetPlayerByUsernameAndRaceId(raceId, username)
	receiver := service.playerService.GetPlayerByUsernameAndRaceId(raceId, receiverUsername)

	if sender.ID == 0 {
		return errors.New(storage.ErrorUndefinedPlayer)
	}

	if receiver.ID == 0 {
		return errors.New(storage.ErrorUndefinedReceiverPlayer)
	}

	if sender.Cash < amount {
		return errors.New(storage.ErrorNotEnoughMoney)
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
			helper.CamelToCapitalize(username),
			strconv.Itoa(amount),
		),
	)

	go service.raceService.SetTransaction(raceId, sender, "", fmt.Sprintf(
		"%s перевёл $%s игроку %s",
		helper.CamelToCapitalize(username),
		strconv.Itoa(amount),
		helper.CamelToCapitalize(receiverUsername),
	))

	return nil
}

func (service *financeService) SendAssets(raceId uint64, username string, dto dto.SendAssetsBodyDTO) error {
	logger.Info("FinanceService.SendAssets", map[string]interface{}{
		"raceId": raceId,
		"dto":    dto,
	})

	sender := service.playerService.GetPlayerByUsernameAndRaceId(raceId, username)
	receiver := service.playerService.GetPlayerByUsernameAndRaceId(raceId, dto.Player)

	if sender.ID == 0 {
		return errors.New(storage.ErrorUndefinedPlayer)
	}

	if receiver.ID == 0 {
		return errors.New(storage.ErrorUndefinedReceiverPlayer)
	}

	var transactionMessage string
	var transactionSenderMessage string
	var transactionReceiverMessage string
	var err error

	if dto.Asset == "stock" {
		err = service.playerService.TransferStocks(dto.AssetId, sender, receiver, dto.Amount)

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
				helper.CamelToCapitalize(username),
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
				helper.CamelToCapitalize(username),
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

func (service *financeService) PayLoan(raceId uint64, username string, amount int) error {
	logger.Info("FinanceService.PayLoan", map[string]interface{}{
		"raceId":   raceId,
		"username": username,
		"amount":   amount,
	})

	player := service.playerService.GetPlayerByUsernameAndRaceId(raceId, username)

	if player.ID == 0 {
		return errors.New(storage.ErrorUndefinedPlayer)
	}

	if amount%1000 != 0 {
		return errors.New(storage.ErrorWrongAmountForPayingLoan)
	}

	err := service.playerService.PayLoan(player, "bankLoan", amount)

	if err == nil {
		go service.raceService.SetTransaction(raceId, player, "", fmt.Sprintf(
			"%s отдал кредит $%s",
			helper.CamelToCapitalize(username),
			strconv.Itoa(amount),
		))
	}

	return err
}

func (service *financeService) TakeLoan(raceId uint64, username string, amount int) error {
	logger.Info("FinanceService.TakeLoan", map[string]interface{}{
		"raceId":   raceId,
		"username": username,
		"amount":   amount,
	})

	player := service.playerService.GetPlayerByUsernameAndRaceId(raceId, username)

	if player.ID == 0 {
		return errors.New(storage.ErrorUndefinedPlayer)
	}

	if amount%1000 != 0 {
		return errors.New(storage.ErrorWrongAmountForTakingLoan)
	}

	err := service.playerService.TakeLoan(player, amount)

	if err == nil {
		go service.raceService.SetTransaction(raceId, player, "", fmt.Sprintf(
			"%s взял в кредит $%s",
			helper.CamelToCapitalize(username),
			strconv.Itoa(amount),
		))
	}

	return err
}
