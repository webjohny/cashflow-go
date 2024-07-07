package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/storage"
)

func (service *playerService) BuyBusinessInPartnership(card entity.CardBusiness, owner entity.Player, players []entity.Player, parts []dto.CardPurchasePlayerActionDTO) error {
	logger.Info("PlayerService.BuyBusinessInPartnership", map[string]interface{}{
		"ownerId": owner.ID,
		"card":    card,
		"parts":   parts,
	})

	cardCost := card.Cost
	cardCashFlow := card.CashFlow

	if cardCashFlow > 0 && card.Limit == 0 {
		fullPassiveIncome := 0

		for _, part := range parts {
			fullPassiveIncome += part.Passive
		}

		if fullPassiveIncome > cardCashFlow {
			return errors.New(storage.ErrorCommonPassiveIncomeGreaterThanCashFlowOfCard)
		}

		for _, part := range parts {
			if part.Percent == 0 {
				part.Percent = (part.Passive * fullPassiveIncome) * 100
			}
		}
	} else if card.Limit > 0 {
		cardCost = 0

		for _, part := range parts {
			cardCost += part.Amount * card.Cost
		}
	}

	if owner.Cash < cardCost {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	for _, part := range parts {
		var currentPlayer entity.Player

		for _, person := range players {
			if int(person.ID) == part.ID {
				currentPlayer = person
			}
		}

		if card.Limit > 0 && part.Amount > 0 {
			card.Count = part.Amount
		} else if part.Passive > 0 {
			card.CashFlow = part.Passive
			card.Percent = part.Percent
		} else {
			return errors.New(storage.ErrorForbidden)
		}

		if owner.ID == currentPlayer.ID {
			card.WholeCost = cardCost
			card.IsOwner = true
		} else {
			card.WholeCost = 0
			card.IsOwner = false
		}

		err := service.BuyBusiness(card, currentPlayer, part.Amount, owner.ID == currentPlayer.ID)

		if err != nil {
			return err
		}
	}

	return nil
}

func (service *playerService) BuyBusiness(card entity.CardBusiness, player entity.Player, count int, updateCash bool) error {
	logger.Info("PlayerService.BuyBusiness", map[string]interface{}{
		"playerId":         player.ID,
		"player.Cash":      player.Cash,
		"card.DownPayment": card.DownPayment,
		"card.Cost":        card.Cost,
		"card.CashFlow":    card.CashFlow,
	})

	_, asset := player.FindBusinessBySymbol(card.Symbol)
	var cost int

	cardCost := card.Cost

	if card.DownPayment > 0 {
		cardCost = card.DownPayment
	}

	if count == 0 {
		count = 1
	}

	if card.Limit > 0 && count > 0 {
		cost = count * cardCost

		if card.Limit > 0 && card.Limit < count {
			return errors.New(storage.ErrorLimitedPartnership)
		}

		if asset.ID != "" {
			asset.Count += count

			asset.SetCardHistory(entity.CardHistory{
				Cost:  cost,
				Price: cardCost,
				Count: count,
			})
		} else {
			card.Count = count

			card.SetCardHistory(entity.CardHistory{
				Cost:  cost,
				Price: cardCost,
				Count: count,
			})

			player.Assets.Business = append(player.Assets.Business, card)
		}
	} else {
		cost = cardCost

		player.Assets.Business = append(player.Assets.Business, card)
	}

	var err error

	if updateCash {
		if card.WholeCost > 0 {
			cost = card.WholeCost
		}

		if player.Cash < cost {
			return errors.New(storage.ErrorNotEnoughMoney)
		}

		service.UpdateCash(&player, -cost, card.Heading)
	} else {
		err, _ = service.UpdatePlayer(&player)
	}

	return err
}

func (service *playerService) BuyRiskBusiness(card entity.CardRiskBusiness, player entity.Player, rolledDice int) (error, bool) {
	logger.Info("PlayerService.BuyRiskBusiness", map[string]interface{}{
		"playerId":   player.ID,
		"card":       card,
		"rolledDice": rolledDice,
	})

	cost := card.Cost

	if player.Cash < cost {
		return errors.New(storage.ErrorNotEnoughMoney), false
	}

	var cashFlow int
	for _, dice := range card.Dices {
		for _, value := range dice.Dices {
			if value == rolledDice {
				cashFlow = *dice.CashFlow
			}
		}
	}

	if cashFlow > 0 {
		service.UpdateCash(&player, -cost, card.Heading)

		business := entity.CardBusiness{
			ID:          card.ID,
			Type:        card.Type,
			Symbol:      card.Symbol,
			Heading:     card.Heading,
			Description: card.Description,
			Cost:        card.Cost,
			CashFlow:    cashFlow,
		}

		player.Assets.Business = append(player.Assets.Business, business)

		return nil, true
	}

	return nil, false
}

func (service *playerService) TransferBusiness(ID string, sender entity.Player, receiver entity.Player, count int) error {
	logger.Info("PlayerService.TransferBusiness", map[string]interface{}{
		"ID":         ID,
		"senderId":   sender.ID,
		"receiverId": receiver.ID,
		"count":      count,
	})

	for index, item := range sender.Assets.Business {
		//@toDo Way for identity item, if in business assets have 2 items with 2 same IDs
		if item.ID == ID {
			if count > 0 {
				if item.Count > 0 {
					senderCount := item.Count

					if senderCount >= count {
						item.Count = count
						receiver.AddBusiness(item)

						item.Count = senderCount - count

						if item.Count > 0 {
							sender.Assets.Business[index] = item
						} else {
							sender.Assets.Business = append(sender.Assets.Business[:index], sender.Assets.Business[index+1:]...)
						}
					} else {
						return errors.New(storage.ErrorNotEnoughAsset)
					}
				}
			} else {
				return errors.New(storage.ErrorPermissionDenied)
			}

			break
		}
	}

	err, _ := service.UpdatePlayer(&sender)

	if err == nil {
		err, _ = service.UpdatePlayer(&receiver)
	}

	return err
}

func (service *playerService) SellBusiness(ID string, card entity.CardMarketBusiness, player entity.Player) (error, int) {
	logger.Info("PlayerService.SellBusiness", map[string]interface{}{
		"ID":       ID,
		"card":     card,
		"playerId": player.ID,
	})

	var totalCash int

	if !player.HasBusiness() {
		return errors.New(storage.ErrorYouHaveNoProperties), 0
	}

	for i := 0; i < len(player.Assets.Business); i++ {
		property := player.Assets.Business[i]
		totalCash += property.Cost / 2
	}

	player.Assets.Business = make([]entity.CardBusiness, 0)

	return nil, totalCash
}

func (service *playerService) MarketBusiness(card entity.CardMarketBusiness, player entity.Player) error {
	logger.Info("PlayerService.MarketBusiness", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	//@toDo make percent for cards where happens cashflow for users / maybe take all amounts-cashflows and calculate it to percents and new cashflow give per cashflows previous values
	assets := player.Assets.Business

	if len(assets) == 0 {
		return errors.New(storage.ErrorNotFoundAssets)
	}

	if card.CashFlow > 0 {
		businessType := card.BusinessType

		for index, asset := range assets {
			assetBusinessType := asset.BusinessType

			if (asset.Symbol == card.Symbol) ||
				(assetBusinessType == businessType) {

				assets[index].CashFlow += card.CashFlow
			} else {
				continue
			}
		}
	}

	player.Assets.Business = assets

	var err error

	if card.Cost > 0 {
		service.UpdateCash(&player, -card.Cost, card.Heading)
	} else {
		err, _ = service.UpdatePlayer(&player)
	}

	return err
}
