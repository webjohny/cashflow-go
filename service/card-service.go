package service

import (
	"encoding/json"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"os"
)

type CardService interface {
	Prepare(raceId uint64, family string, actionType string, username string) (error, interface{})
	Accept(raceId uint64, family string, actionType string, username string) (error, interface{})
	Purchase(raceId uint64, actionType string, username string, count int) (error, interface{})
	Selling(raceId uint64, actionType string, username string) (error, interface{})
	Skip(raceId uint64, username string) (error, interface{})
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

func (service *cardService) Prepare(raceId uint64, family string, actionType string, username string) (error, interface{}) {
	var err error

	//if actionType == "risk" || actionType == "riskStock" {
	//	err = service.raceService.PreRiskAction(raceId, username, actionType)
	//}
	return err, nil
}

func (service *cardService) Accept(raceId uint64, family string, actionType string, username string) (error, interface{}) {
	var err error
	var response interface{}

	if family == "payday" {
		err = service.raceService.PaydayAction(raceId, username, actionType)
	} else if family == "market" && actionType == "damage" {
		err = service.raceService.MarketAction(raceId, username, actionType)
	}
	return err, response
}

func (service *cardService) Skip(raceId uint64, username string) (error, interface{}) {
	var err error
	var response interface{}

	return err, response
}

func (service *cardService) Purchase(raceId uint64, actionType string, username string, count int) (error, interface{}) {
	var err error
	var response interface{}

	switch actionType {
	case "business":
		err = service.raceService.BusinessAction(raceId, username, actionType)
		break

	case "realEstate":
		err = service.raceService.RealEstateAction(raceId, username, actionType)
		break

	case "dream":
		err = service.raceService.DreamAction(raceId, username, actionType)
		break

	case "riskBusiness":
		err, response = service.raceService.RiskBusinessAction(raceId, username, actionType)
		break

	case "riskStocks":
		err, response = service.raceService.RiskStocksAction(raceId, username, actionType)
		break

	case "stocks":
		err = service.raceService.StocksAction(raceId, username, actionType, count)
		break

	case "lottery":
		err = service.raceService.LotteryAction(raceId, username, actionType)
		break

	case "goldCoins":
		err = service.raceService.GoldCoinsAction(raceId, username, actionType)
		break

	case "mlm":
		err = service.raceService.MlmAction(raceId, username, actionType)
		break

	default:
		err = service.raceService.SkipAction(raceId, username, actionType)
		break
	}

	return err, response
}

func (service *cardService) Selling(raceId uint64, actionType string, username string) (error, interface{}) {
	return nil, nil
}

func (service *cardService) GetCards() []entity.Card {
	data, err := os.ReadFile(os.Getenv("CARDS_PATH"))
	if err != nil {
		panic(err)
	}

	var cards []entity.Card

	err = json.Unmarshal(data, &cards)
	if err != nil {
		panic(err)
	}

	return cards
}

func (service *cardService) GetCard(cardType string) entity.Card {
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
		card := service.PickCard(cardType)
		card.Family = service.GetFamily(deals, cardType)
		card.Name = cardType
		return card
	}

	return entity.Card{}
}

func (service *cardService) GetFamily(deals []string, dealType string) string {
	if helper.Contains(deals, dealType) {
		dealType = "deal"
	}

	return dealType
}

func (service *cardService) GetBigCardType(tilePosition int) string {
	for tile, positions := range service.ratRace.Tiles {
		for _, position := range positions {
			if position == tilePosition {
				return tile
			}
		}
	}
	return ""
}

func (service *cardService) GetRatCardType(tilePosition int) string {
	for tile, positions := range service.bigRace.Tiles {
		for _, position := range positions {
			if position == tilePosition {
				return tile
			}
		}
	}
	return ""
}

func (service *cardService) PickCard(cardType string) entity.Card {
	if helper.Contains([]string{}, cardType) {
		return entity.Card{}
	}

	cardList := service.GetCards()
	return cardList[helper.Random(len(cardList)-1)]
}
