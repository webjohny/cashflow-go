package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/storage"
)

func (service *playerService) BuyRealEstate(card entity.CardRealEstate, player entity.Player) error {
	logger.Info("PlayerService.BuyRealEstate", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	if player.Cash < card.DownPayment {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	player.Assets.RealEstates = append(player.Assets.RealEstates, card)

	service.UpdateCash(&player, -card.DownPayment, card.Heading)

	return nil
}

func (service *playerService) BuyRealEstateInPartnership(card entity.CardRealEstate, owner entity.Player, players []entity.Player, parts []dto.CardPurchasePlayerActionDTO) error {
	logger.Info("PlayerService.BuyRealEstateInPartnership", map[string]interface{}{
		"ownerId": owner.ID,
		"card":    card,
		"parts":   parts,
	})

	if owner.Cash < card.DownPayment {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	cardCost := card.DownPayment

	var err error

	mortgage := card.Mortgage
	cost := card.Cost

	for _, pl := range parts {
		var currentPlayer entity.Player

		for _, person := range players {
			if int(person.ID) == pl.ID {
				currentPlayer = person
			}
		}

		card.CashFlow = pl.Passive
		card.DownPayment = pl.Amount

		if int(owner.ID) != pl.ID {
			card.Mortgage = 0
			card.Cost = 0
			card.IsOwner = false
		} else {
			card.Mortgage = mortgage
			card.Cost = cost
			card.IsOwner = true
		}

		currentPlayer.CashFlow += card.CashFlow
		currentPlayer.Assets.RealEstates = append(currentPlayer.Assets.RealEstates, card)

		if owner.ID == currentPlayer.ID {
			service.UpdateCash(&currentPlayer, -cardCost, card.Heading)
		} else {
			card.Cost = 0
			err, _ = service.UpdatePlayer(&currentPlayer)
		}

		if err != nil {
			logger.Error(err, nil)

			return err
		}
	}

	return nil
}

func (service *playerService) SellAllProperties(player entity.Player) (error, int) {
	logger.Info("PlayerService.SellAllProperties", map[string]interface{}{
		"playerId": player.ID,
	})

	var totalCash int

	if !player.HasRealEstates() {
		return errors.New(storage.ErrorYouHaveNoProperties), 0
	}

	for i := 0; i < len(player.Assets.RealEstates); i++ {
		property := player.Assets.RealEstates[i]
		totalCash += property.DownPayment / 2
	}

	player.Assets.RealEstates = make([]entity.CardRealEstate, 0)

	return nil, totalCash
}

func (service *playerService) SellRealEstate(ID string, card entity.CardMarketRealEstate, player entity.Player) error {
	logger.Info("PlayerService.SellRealEstate", map[string]interface{}{
		"playerId":     player.ID,
		"card":         card,
		"realEstateId": ID,
	})

	realEstate := player.FindRealEstate(ID)

	if realEstate.ID == "" {
		return errors.New(storage.ErrorNotFoundAssets)
	}

	value := (realEstate.Cost / 100) * card.Value
	totalCost := realEstate.Cost + value

	if card.Plus {
		totalCost = realEstate.Cost + card.Value
	}

	service.UpdateCash(&player, totalCost-realEstate.Mortgage, card.Symbol)

	player.RemoveRealEstate(card.ID)

	return nil
}
