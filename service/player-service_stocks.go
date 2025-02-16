package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/storage"
	"math"
)

func (service *playerService) BuyStocks(card entity.CardStocks, player entity.Player, updateCash bool) error {
	logger.Info("PlayerService.BuyStocks", map[string]interface{}{
		"playerId":   player.ID,
		"card":       card,
		"updateCash": updateCash,
	})

	totalCost := card.Price * card.Count

	if updateCash && player.Cash < totalCost {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	key, stock := player.FindStocksBySymbol(card.Symbol)

	if stock.ID != "" {
		totalCount := stock.Count + card.Count
		stock.Count = totalCount

		stock.SetCardHistory(entity.CardHistory{
			Price: card.Price,
			Count: card.Count,
			Cost:  card.Count * card.Price,
		})

		player.Assets.Stocks[key] = *stock
	} else {
		card.SetCardHistory(entity.CardHistory{
			Price: card.Price,
			Count: card.Count,
			Cost:  card.Count * card.Price,
		})

		player.Assets.Stocks = append(player.Assets.Stocks, card)
	}

	var err error

	if updateCash {
		err = service.UpdateCash(&player, -totalCost, &dto.TransactionDTO{
			CardID:   card.ID,
			CardType: entity.TransactionCardType.RealEstate,
			Details:  card.Heading,
		})
	} else {
		err, _ = service.UpdatePlayer(&player)
	}

	return err
}

func (service *playerService) SellStocks(card entity.CardStocks, player entity.Player, count int, updateCash bool) error {
	logger.Info("PlayerService.SellStocks", map[string]interface{}{
		"playerId":   player.ID,
		"card":       card,
		"count":      count,
		"updateCash": updateCash,
	})

	if count <= 0 {
		return errors.New(storage.ErrorIncorrectCount)
	}

	_, stock := player.FindStocksBySymbol(card.Symbol)

	if stock.ID == "" || stock.Count < count {
		return errors.New(storage.ErrorNotFoundStocks)
	}

	totalCost := card.Price * count
	stock.Count -= count

	if stock.Count <= 0 {
		player.RemoveStocks(stock.Symbol)
	} else {
		player.ReduceStocks(stock.Symbol, count)
	}

	var err error

	if updateCash {
		err = service.UpdateCash(&player, totalCost, &dto.TransactionDTO{
			CardID:   card.ID,
			CardType: entity.TransactionCardType.RealEstate,
			Details:  card.Heading,
		})
	} else {
		err, _ = service.UpdatePlayer(&player)
	}

	return err
}

func (service *playerService) DecreaseStocks(card entity.CardStocks, player entity.Player) error {
	logger.Info("PlayerService.DecreaseStocks", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	_, stock := player.FindStocksBySymbol(card.Symbol)

	if stock.ID == "" {
		return errors.New(storage.ErrorNotFoundStocks)
	}

	var count int
	for i := 0; i < len(stock.History); i++ {
		stock.History[i].Count = int(math.Floor(float64(stock.History[i].Count) / float64(card.Decrease)))
		count += stock.History[i].Count
	}

	stock.Count = count

	err, _ := service.UpdatePlayer(&player)

	return err
}

func (service *playerService) IncreaseStocks(card entity.CardStocks, player entity.Player) error {
	logger.Info("PlayerService.IncreaseStocks", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	_, stock := player.FindStocksBySymbol(card.Symbol)

	if stock.ID == "" {
		return errors.New(storage.ErrorNotFoundStocks)
	}

	stock.Count = stock.Count * card.Increase

	for i := 0; i < len(stock.History); i++ {
		stock.History[i].Count = stock.History[i].Count * card.Increase
	}

	err, _ := service.UpdatePlayer(&player)

	return err
}

func (service *playerService) TransferStocks(ID string, sender entity.Player, receiver entity.Player, count int) error {
	logger.Info("TransferStocks: init", map[string]interface{}{
		"ID":         ID,
		"senderId":   sender.ID,
		"receiverId": receiver.ID,
		"count":      count,
	})

	if count <= 0 {
		return errors.New(storage.ErrorIncorrectCount)
	}

	_, item := sender.FindStocksByID(ID)

	logger.Info("TransferStocks: getting info about a sender", map[string]interface{}{
		"senderId":    sender.ID,
		"playerCount": item.Count,
		"count":       count,
	})

	if item.Count < count {
		return errors.New(storage.ErrorNotEnoughStocks)
	}

	err := service.SellStocks(*item, sender, count, false)

	if err != nil {
		logger.Error(err, sender, map[string]interface{}{
			"ID":     sender.ID,
			"raceID": sender.RaceID,
		})

		return err
	}

	item.Count = count
	item.History = make([]entity.CardHistory, 0)
	err = service.BuyStocks(*item, receiver, false)

	if err != nil {
		logger.Error(err, receiver, map[string]interface{}{
			"ID":     receiver.ID,
			"raceID": receiver.RaceID,
		})

		return err
	}

	return nil
}
