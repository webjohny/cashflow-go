package service

import (
	"fmt"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/storage"
	"strconv"
)

type FinanceService interface {
	SendMoney(raceId uint64, username string, amount int, player string) error
	SendAssets(raceId uint64, username string, amount int, player string, asset string) error
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
	sender := service.playerService.GetPlayerByUsernameAndRaceId(raceId, username)
	receiver := service.playerService.GetPlayerByUsernameAndRaceId(raceId, receiverUsername)

	if sender.ID == 0 {
		return fmt.Errorf(storage.ErrorUndefinedPlayer)
	}

	if receiver.ID == 0 {
		return fmt.Errorf(storage.ErrorUndefinedReceiverPlayer)
	}

	if sender.Cash < amount {
		return fmt.Errorf(storage.ErrorNotEnoughMoney)
	}

	service.playerService.UpdateCash(
		sender,
		amount,
		fmt.Sprintf(
			"Перевёл $%s игроку %s",
			strconv.Itoa(amount),
			helper.CamelToCapitalize(receiverUsername),
		),
	)

	service.playerService.UpdateCash(
		receiver,
		amount,
		fmt.Sprintf(
			"%s перевёл Вам $%s",
			helper.CamelToCapitalize(username),
			strconv.Itoa(amount),
		),
	)

	go service.raceService.SetTransaction(raceId, sender, fmt.Sprintf(
		"%s перевёл $%s игроку %s",
		helper.CamelToCapitalize(username),
		strconv.Itoa(amount),
		helper.CamelToCapitalize(receiverUsername),
	))

	return nil
}

func (service *financeService) SendAssets(raceId uint64, username string, amount int, receiverUsername string, asset string) error {
	sender := service.playerService.GetPlayerByUsernameAndRaceId(raceId, username)
	receiver := service.playerService.GetPlayerByUsernameAndRaceId(raceId, receiverUsername)

	if sender.ID == 0 {
		return fmt.Errorf(storage.ErrorUndefinedPlayer)
	}

	if receiver.ID == 0 {
		return fmt.Errorf(storage.ErrorUndefinedReceiverPlayer)
	}

	var transactionMessage string
	var transactionSenderMessage string
	var transactionReceiverMessage string
	var err error

	if asset == "stocks" {
		for i := 0; i < len(sender.Assets.Stocks); i++ {
			stocks := sender.Assets.Stocks[i]

			if amount > *stocks.Count {
				return fmt.Errorf(storage.ErrorNotEnoughStocks)
			}

			err = service.playerService.SellStocks(stocks, sender, amount, false)

			if err == nil {
				transactionSenderMessage = fmt.Sprintf(
					"Перевёл акции (%s шт.) игроку %s",
					strconv.Itoa(amount),
					helper.CamelToCapitalize(receiverUsername),
				)
				transactionReceiverMessage = fmt.Sprintf(
					"%s перевёл Вам (%s шт.) акций",
					helper.CamelToCapitalize(receiverUsername),
					strconv.Itoa(amount),
				)
				transactionMessage = fmt.Sprintf(
					"%s перевёл акции (%s шт.) игроку %s",
					helper.CamelToCapitalize(username),
					strconv.Itoa(amount),
					helper.CamelToCapitalize(receiverUsername),
				)

				err = service.playerService.BuyStocks(stocks, sender, amount, false)
			}
		}
	}

	if err == nil {
		go service.playerService.SetTransaction(sender.ID, sender.Cash, sender.Cash, amount, transactionSenderMessage)
		go service.playerService.SetTransaction(receiver.ID, receiver.Cash, receiver.Cash, amount, transactionReceiverMessage)
		go service.raceService.SetTransaction(raceId, sender, transactionMessage)
	}

	return err
}

func (service *financeService) TakeLoan(raceId uint64, username string, amount int) error {
	player := service.playerService.GetPlayerByUsernameAndRaceId(raceId, username)

	if player.ID == 0 {
		return fmt.Errorf(storage.ErrorUndefinedPlayer)
	}

	err := service.playerService.TakeLoan(player, amount)

	if err == nil {
		go service.raceService.SetTransaction(raceId, player, fmt.Sprintf(
			"%s взял в кредит $%s",
			helper.CamelToCapitalize(username),
			strconv.Itoa(amount),
		))
	}

	return err
}
