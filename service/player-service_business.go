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

	if cardCashFlow > 0 && card.AssetType != entity.BusinessTypes.Limited {
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
	} else if card.AssetType == entity.BusinessTypes.Limited {
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

		if card.AssetType == entity.BusinessTypes.Limited && part.Amount > 0 {
			card.Count = part.Amount
		} else if card.AssetType != entity.BusinessTypes.Limited {
			card.CashFlow = part.Passive
			card.Percent = part.Percent
		} else {
			return errors.New(storage.ErrorForbidden)
		}

		if owner.ID == currentPlayer.ID {
			card.WholeCost = cardCost
		} else {
			card.WholeCost = 0
		}

		card.IsOwner = card.AssetType == entity.BusinessTypes.Limited || owner.ID == currentPlayer.ID

		err := service.BuyBusiness(card, currentPlayer, part.Amount)

		if err != nil {
			logger.Error(err, nil)

			return err
		}
	}

	return nil
}

func (service *playerService) BuyBusiness(card entity.CardBusiness, player entity.Player, count int) error {
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

	if card.AssetType == entity.BusinessTypes.Limited && count > 0 {
		cost = count * cardCost

		if card.Limit < count {
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

	if card.IsOwner {
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

	err, player := service.UpdatePlayer(&sender)

	if err != nil {
		logger.Error(err, player)

		return err
	}

	err, _ = service.UpdatePlayer(&receiver)

	return err
}

func (service *playerService) SellBusiness(ID string, card entity.CardMarketBusiness, player entity.Player, count int) (error, int) {
	logger.Info("PlayerService.SellBusiness", map[string]interface{}{
		"ID":       ID,
		"card":     card,
		"playerId": player.ID,
	})

	var totalCash int

	if !player.HasBusiness() {
		return errors.New(storage.ErrorYouHaveNoProperties), 0
	}

	_, business := player.FindBusinessByID(ID)

	if !business.IsOwner {
		return errors.New(storage.ErrorForbidden), 0
	}

	if card.Cost < 10 {
		totalCash = business.Cost * card.Cost
	} else if card.Cost >= 10 && card.Cost <= 100 {
		totalCash = (business.Cost / 100) * card.Cost
	} else if card.Cost > 100 {
		totalCash = card.Cost
	}

	if business.AssetType == entity.BusinessTypes.Limited && count > 0 {
		totalCash *= count
		business.Count -= count
		player.ReduceLimitedShares(ID, count)
	}

	if business.Count <= 0 || business.AssetType != entity.BusinessTypes.Limited {
		player.RemoveBusiness(ID)
	}

	if totalCash > 0 {
		service.UpdateCash(&player, totalCash, card.Heading)
	}

	if business.AssetType == entity.BusinessTypes.Limited {
		return nil, totalCash
	}

	players := service.GetAllPlayersByRaceId(player.RaceID)

	for _, user := range players {
		_, asset := player.FindBusinessByID(ID)

		if ID == asset.ID && !asset.IsOwner && asset.AssetType != entity.BusinessTypes.Limited {
			user.RemoveBusiness(ID)

			err, play := service.UpdatePlayer(&user)

			if err != nil {
				logger.Error("SellBusiness.UpdatePlayer", play, ID, user.ID, user.RaceID)
			}
		}
	}

	return nil, totalCash
}

func (service *playerService) MarketBusiness(card entity.CardMarketBusiness, player entity.Player) error {
	logger.Info("PlayerService.MarketBusiness", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	//@toDo make percent for cards where happens cashflow for users / maybe take all amounts-cashflows and calculate it to percents and new cashflow give per cashflows previous values
	businesses := &player.Assets.Business

	if len(*businesses) == 0 {
		return errors.New(storage.ErrorNotFoundAssets)
	}

	if card.CashFlow > 0 {
		businessType := card.AssetType

		for _, asset := range *businesses {
			assetAssetType := asset.AssetType

			if (asset.Symbol == card.Symbol) ||
				(assetAssetType == businessType) {

				asset.CashFlow += card.CashFlow
			} else {
				continue
			}
		}
	}

	var err error

	if card.Cost > 0 {
		service.UpdateCash(&player, -card.Cost, card.Heading)
	} else {
		err, _ = service.UpdatePlayer(&player)
	}

	return err
}
