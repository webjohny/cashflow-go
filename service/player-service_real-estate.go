package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
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

	return service.AreYouBankrupt(player)
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
	mortgage := card.Mortgage
	cost := card.Cost

	var cashFlow int

	for _, part := range parts {
		cashFlow += part.Passive
	}

	if cashFlow > card.CashFlow {
		return errors.New(storage.ErrorCommonPassiveIncomeGreaterThanCashFlowOfCard)
	}

	for _, pl := range parts {
		var currentPlayer entity.Player

		for _, person := range players {
			if int(person.ID) == pl.ID {
				currentPlayer = person
			}
		}

		card.CashFlow = pl.Passive
		card.DownPayment = pl.Amount

		if owner.ID == currentPlayer.ID {
			card.Mortgage = mortgage
			card.Cost = cost
			card.IsOwner = true
		} else {
			card.Cost = 0
			card.Mortgage = 0
			card.IsOwner = false
		}

		currentPlayer.Assets.RealEstates = append(currentPlayer.Assets.RealEstates, card)

		if owner.ID == currentPlayer.ID {
			service.UpdateCash(&currentPlayer, -cardCost, card.Heading)
		} else {
			err, player := service.UpdatePlayer(&currentPlayer)

			if err != nil {
				logger.Error(err, player)

				return err
			}
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

	return service.AreYouBankrupt(player), totalCash
}

func (service *playerService) SellRealEstate(ID string, card entity.CardMarketRealEstate, player entity.Player) error {
	logger.Info("PlayerService.SellRealEstate", map[string]interface{}{
		"playerId":     player.ID,
		"card":         card,
		"realEstateId": ID,
	})

	var totalCost int

	realEstate := player.FindRealEstateByID(ID)

	if realEstate.ID == "" {
		return errors.New(storage.ErrorNotFoundAssets)
	}

	if !realEstate.IsOwner {
		return errors.New(storage.ErrorForbiddenByOwner)
	}

	if card.AssetType == entity.RealEstateTypes.Building {
		if helper.Contains[int](card.Range, realEstate.Count) || len(card.Range) == 0 {
			totalCost = card.Cost * realEstate.Count
		} else {
			return errors.New(storage.ErrorNotSuitableBuilding)
		}
	} else if card.AssetType == entity.RealEstateTypes.Single {
		totalCost = card.Cost
	}

	player.RemoveRealEstate(ID)

	if totalCost > 0 && totalCost >= realEstate.Mortgage {
		service.UpdateCash(&player, totalCost-realEstate.Mortgage, card.Heading)
	}

	players := service.GetAllPlayersByRaceId(player.RaceID)

	for _, user := range players {
		asset := user.FindRealEstateByID(ID)

		if ID == asset.ID && !asset.IsOwner {
			user.RemoveRealEstate(ID)

			err, play := service.UpdatePlayer(&user)

			if err != nil {
				logger.Error("SellRealEstate.UpdatePlayer", play, ID, user.ID, user.RaceID)
			}
		}
	}

	return service.AreYouBankrupt(player)
}
