package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
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

	if player.Cash < totalCost {
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
		service.UpdateCash(&player, -totalCost, card.Symbol)
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

	if updateCash {
		service.UpdateCash(&player, totalCost, card.Symbol)
	}

	return nil
}

func (service *playerService) DecreaseStocks(card entity.CardStocks, player entity.Player) error {
	logger.Info("PlayerService.SellRealEstate", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	_, stock := player.FindStocksBySymbol(card.Symbol)

	if stock.ID == "" {
		return errors.New(storage.ErrorNotFoundStocks)
	}

	stock.Count = int(math.Floor(float64(stock.Count / card.Decrease)))

	for i := 0; i < len(stock.History); i++ {
		stock.History[i].Count = int(math.Floor(float64(stock.History[i].Count) / float64(card.Decrease)))
	}

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
	logger.Info("PlayerService.TransferStocks", map[string]interface{}{
		"ID":         ID,
		"senderId":   sender.ID,
		"receiverId": receiver.ID,
		"count":      count,
	})

	if count <= 0 {
		return errors.New(storage.ErrorIncorrectCount)
	}

	for index, item := range sender.Assets.Stocks {
		if item.ID == ID && item.Count > 0 {
			senderCount := item.Count

			if senderCount < count {
				return errors.New(storage.ErrorNotEnoughStocks)
			}

			item.Count = count
			receiver.AddStocks(item)

			item.Count = senderCount - count
			if item.Count > 0 {
				sender.Assets.Stocks[index] = item
			} else {
				sender.Assets.Stocks = append(sender.Assets.Stocks[:index], sender.Assets.Stocks[index+1:]...)
			}

			break
		}
	}

	err, _ := service.UpdatePlayer(&sender)

	if err != nil {
		logger.Error(err, map[string]interface{}{
			"ID":     sender.ID,
			"raceID": sender.RaceID,
		})

		return err
	}

	if err == nil {
		err, _ = service.UpdatePlayer(&receiver)

		if err != nil {
			logger.Error(err, map[string]interface{}{
				"ID":     receiver.ID,
				"raceID": receiver.RaceID,
			})

			return err
		}
	}

	return err
}

func (service *playerService) BuyRiskStocks(card entity.CardRiskStocks, player entity.Player, rolledDice int) (error, bool) {
	logger.Info("PlayerService.BuyRiskStocks", map[string]interface{}{
		"playerId":   player.ID,
		"card":       card,
		"rolledDice": rolledDice,
	})

	cost := card.Cost

	if player.Cash < cost {
		return errors.New(storage.ErrorNotEnoughMoney), false
	}

	var costPerOne int
	for _, dice := range card.Dices {
		for _, value := range dice.Dices {
			if value == rolledDice {
				costPerOne = *dice.CostPerOne
			}
		}
	}

	if costPerOne > 0 {
		service.UpdateCash(&player, -cost, card.Heading)
		service.UpdateCash(&player, card.Count*costPerOne, card.Heading)

		return nil, true
	}

	return nil, false
}
