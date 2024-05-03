package service

import (
	"encoding/json"
	"fmt"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/logger"
	"github.com/webjohny/cashflow-go/storage"
	"log"
	"os"
	"strconv"
)

type CardService interface {
	Prepare(raceId uint64, family string, actionType string, userId uint64, isBigRace bool) (error, interface{})
	Accept(raceId uint64, family string, actionType string, userId uint64, isBigRace bool) (error, interface{})
	Purchase(raceId uint64, actionType string, userId uint64, count int, isBigRace bool) (error, interface{})
	Selling(raceId uint64, actionType string, userId uint64, value string, isBigRace bool) (error, interface{})
	Skip(raceId uint64, userId uint64, isBigRace bool) (error, interface{})
	GetCard(action string, raceId uint64, userId uint64, isBigRace bool) (error, entity.Card)
}

type CardRatRace struct {
	Tiles         map[string][]int
	Notifications []string
}

type CardBigRace struct {
	Tiles         map[string][]int
	Notifications []string
}

type cardService struct {
	gameService GameService
	raceService RaceService
	ratRace     CardRatRace
	bigRace     CardBigRace
}

func NewCardService(gameService GameService, raceService RaceService) CardService {
	return &cardService{
		gameService: gameService,
		raceService: raceService,
		ratRace: CardRatRace{
			Tiles: map[string][]int{
				"deals":     {1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23},
				"payday":    {6, 14, 22},
				"market":    {8, 16, 24},
				"doodad":    {2, 10, 18},
				"charity":   {4},
				"downsized": {12},
				"baby":      {20},
			},
			Notifications: []string{"payday", "bigCharity", "baby"},
		},
		bigRace: CardBigRace{
			Tiles: map[string][]int{
				"business":      {2, 4, 6, 9, 11, 14, 18, 20, 22, 24, 28, 32, 34, 36, 38, 40, 44},
				"cashFlowDay":   {12, 26, 42},
				"bigCharity":    {8},
				"dream":         {1, 3, 5, 7, 10, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31, 33, 35, 37, 39, 41, 43, 45},
				"bankrupt":      {46},
				"tax50percent":  {16},
				"tax100percent": {30},
			},
			Notifications: []string{"cashFlowDay"},
		},
	}
}

func (service *cardService) Prepare(raceId uint64, family string, actionType string, userId uint64, isBigRace bool) (error, interface{}) {
	logger.Info("CardService.Prepare", map[string]interface{}{
		"raceId":     raceId,
		"family":     family,
		"actionType": actionType,
		"userId":     userId,
	})

	if family == "deal" {
		err, card := service.GetCard(actionType, raceId, userId, isBigRace)

		return err, card
	}

	//if actionType == "risk" || actionType == "riskStock" {
	//	err = service.raceService.PreRiskAction(raceId, username, actionType)
	//}
	return fmt.Errorf(storage.ErrorForbidden), nil
}

func (service *cardService) Accept(raceId uint64, family string, actionType string, userId uint64, isBigRace bool) (error, interface{}) {
	logger.Info("CardService.Accept", map[string]interface{}{
		"raceId":     raceId,
		"family":     family,
		"actionType": actionType,
		"userId":     userId,
	})

	var err error
	var response interface{}

	if family == "payday" {
		err = service.raceService.PaydayAction(raceId, userId, actionType, isBigRace)
	} else if family == "market" && actionType == "damage" {
		err = service.raceService.MarketAction(raceId, userId, actionType)
	} else if family == "charity" {
		err = service.raceService.CharityAction(raceId, userId, isBigRace)
	} else if family == "doodad" {
		err = service.raceService.DoodadAction(raceId, userId)
	} else if family == "baby" {
		err = service.raceService.BabyAction(raceId, userId)
	} else if family == "downsized" {
		err = service.raceService.DownsizedAction(raceId, userId)
	} else {
		err = service.raceService.SkipAction(raceId, userId, isBigRace)
	}

	return err, response
}

func (service *cardService) Skip(raceId uint64, userId uint64, isBigRace bool) (error, interface{}) {
	logger.Info("CardService.Skip", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	var err error
	var response interface{}

	err = service.raceService.SkipAction(raceId, userId, isBigRace)

	return err, response
}

func (service *cardService) Purchase(raceId uint64, actionType string, userId uint64, count int, isBigRace bool) (error, interface{}) {
	logger.Info("CardService.Purchase", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	var err error
	var response interface{}

	switch actionType {
	case "business":
		err = service.raceService.BusinessAction(raceId, userId, isBigRace)
		break

	case "realEstate":
		err = service.raceService.RealEstateAction(raceId, userId, isBigRace)
		break

	case "dream":
		err = service.raceService.DreamAction(raceId, userId)
		break

	case "riskBusiness":
		err, response = service.raceService.RiskBusinessAction(raceId, userId)
		break

	case "riskStocks":
		err, response = service.raceService.RiskStocksAction(raceId, userId)
		break

	case "stocks":
		err = service.raceService.StocksAction(raceId, userId, count)
		break

	case "lottery":
		err = service.raceService.LotteryAction(raceId, userId, isBigRace)
		break

	case "goldCoins":
		err = service.raceService.GoldCoinsAction(raceId, userId)
		break

	case "mlm":
		err = service.raceService.MlmAction(raceId, userId, isBigRace)
		break

	default:
		err = service.raceService.SkipAction(raceId, userId, isBigRace)
		break
	}

	return err, response
}

func (service *cardService) Selling(raceId uint64, actionType string, userId uint64, value string, isBigRace bool) (error, interface{}) {
	logger.Info("CardService.Selling", map[string]interface{}{
		"raceId":     raceId,
		"userId":     userId,
		"actionType": actionType,
		"value":      value,
	})

	var err error
	var response interface{}

	switch actionType {
	case "realEstate":
		if value != "" {
			return fmt.Errorf(storage.ErrorIsNotValidRealEstate), nil
		}

		err = service.raceService.SellRealEstate(raceId, userId, value)
		break
	case "stock":
		count, _ := strconv.Atoi(value)

		if count < 1 {
			return fmt.Errorf(storage.ErrorIsNotValidCountValue), nil
		}

		err = service.raceService.SellStocks(raceId, userId, count)
		break
	case "goldCoins":
		count, _ := strconv.Atoi(value)

		if count < 1 {
			return fmt.Errorf(storage.ErrorIsNotValidCountValue), nil
		}

		err = service.raceService.SellGoldCoins(raceId, userId, count)
		break
	case "skip":
		err = service.raceService.SkipAction(raceId, userId, isBigRace)
		break
	}

	return err, response
}

func (service *cardService) GetCard(action string, raceId uint64, userId uint64, isBigRace bool) (error, entity.Card) {
	logger.Info("CardService.GetCard", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"action": action,
	})

	err, race, player := service.raceService.GetRaceAndPlayer(raceId, userId, isBigRace)

	if err != nil {
		return err, entity.Card{}
	}

	var tile string

	if isBigRace {
		tile = service.getBigCardType(int(player.CurrentPosition))
	} else {
		tile = service.getRatCardType(int(player.CurrentPosition))
	}

	if action == "big" && player.Cash < 10000 {
		return fmt.Errorf(storage.ErrorCannotTakeBigDeals), entity.Card{}
	}

	if action != "" {
		tile = action + "Deal"
	}

	race.CurrentCard = service.getCardByTile(tile)

	err, _ = service.raceService.UpdateRace(&race)

	return err, race.CurrentCard
}

func (service *cardService) GetCards() map[string][]entity.Card {
	logger.Info("CardService.GetCards", nil)

	data, err := os.ReadFile(os.Getenv("CARDS_PATH"))
	if err != nil {
		panic(err)
	}

	var cards map[string][]entity.Card

	err = json.Unmarshal(data, &cards)
	if err != nil {
		panic(err)
	}

	return cards
}

func (service *cardService) getCardByTile(cardType string) entity.Card {
	logger.Info("CardService.getCardByTile", map[string]interface{}{
		"cardType": cardType,
	})

	deals := []string{"smallDeal", "bigDeal"}
	validTypes := append(deals,
		"market",
		"doodad",
		"charity",
		"baby",
		"downsized",
		"payday",
		"business",
		"dream",
		"tax50percent",
		"tax100percent",
		"bankrupt",
	)

	if cardType == "deals" {
		return entity.Card{
			ID:      "deal",
			Heading: "Выберите маленькую или большую сделку",
			Family:  "deal",
			Type:    "deal",
		}
	}

	if helper.Contains(validTypes, cardType) {
		card := service.getPickCard(cardType)
		card.Family = service.getFamily(deals, cardType)
		card.Name = cardType
		log.Println(card.DownPayment, card.CashFlow)
		return card
	}

	return entity.Card{}
}

func (service *cardService) getFamily(deals []string, dealType string) string {
	logger.Info("CardService.getFamily", map[string]interface{}{
		"dealType": dealType,
	})

	if helper.Contains(deals, dealType) {
		dealType = "deal"
	}

	return dealType
}

func (service *cardService) getBigCardType(tilePosition int) string {
	for tile, positions := range service.ratRace.Tiles {
		for _, position := range positions {
			if position == tilePosition {
				return tile
			}
		}
	}
	return ""
}

func (service *cardService) getRatCardType(tilePosition int) string {
	for tile, positions := range service.ratRace.Tiles {
		for _, position := range positions {
			if position == tilePosition {
				return tile
			}
		}
	}
	return ""
}

func (service *cardService) getPickCard(cardType string) entity.Card {
	logger.Info("CardService.getPickCard", map[string]interface{}{
		"cardType": cardType,
	})

	if helper.Contains([]string{}, cardType) {
		return entity.Card{}
	}

	cardList := service.GetCards()
	log.Println(len(cardList[cardType]))
	return cardList[cardType][helper.Random(len(cardList[cardType])-1)]
}
