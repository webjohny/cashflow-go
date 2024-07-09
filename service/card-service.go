package service

import (
	"encoding/json"
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/storage"
	"os"
)

type CardService interface {
	Prepare(actionType string, raceId uint64, family string, userId uint64, isBigRace bool) (error, interface{})
	Accept(actionType string, raceId uint64, family string, userId uint64, isBigRace bool) (error, interface{})
	Purchase(actionType string, raceId uint64, userId uint64, isBigRace bool, dto dto.CardPurchaseActionDTO) (error, interface{})
	Selling(actionType string, raceId uint64, userId uint64, isBigRace bool, dto dto.CardSellingActionDTO) (error, interface{})
	Skip(raceId uint64, userId uint64, isBigRace bool) (error, interface{})
	GetCard(action string, raceId uint64, userId uint64, isBigRace bool) (error, entity.Card)
	TestCard(action string, raceId uint64, userId uint64, isBigRace bool) (error, entity.Card)
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
	usedCardRepository repository.UsedCardRepository
	gameService        GameService
	raceService        RaceService
	playerService      PlayerService
	ratRace            CardRatRace
	bigRace            CardBigRace
}

func NewCardService(usedCardRepository repository.UsedCardRepository, gameService GameService, raceService RaceService, playerService PlayerService) CardService {
	return &cardService{
		usedCardRepository: usedCardRepository,
		gameService:        gameService,
		raceService:        raceService,
		playerService:      playerService,
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

func (service *cardService) TestCard(action string, raceId uint64, userId uint64, isBigRace bool) (error, entity.Card) {
	logger.Info("CardService.TestCard", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"action": action,
	})
	err, race, player := service.raceService.GetRaceAndPlayer(raceId, userId, isBigRace)

	if err != nil {
		return err, entity.Card{}
	}

	var tile string

	tileName := action

	if !isBigRace {
		if tileName == "small" || tileName == "big" {
			tileName = "deals"
		}

		player.CurrentPosition = uint8(service.ratRace.Tiles[tileName][0])

		tile = service.getRatCardType(int(player.CurrentPosition))

		if action == "small" || action == "big" {
			tile = action + "Deal"
		}
	} else {
		player.CurrentPosition = uint8(service.bigRace.Tiles[tileName][0])

		tile = service.getBigCardType(int(player.CurrentPosition))
	}

	race.CurrentCard = service.getCardByTile(tile, race.CardMap.Active[action])

	err = service.processCard(action, race, player)

	if err == nil {
		err, _ = service.raceService.UpdateRace(&race)
	}

	return err, race.CurrentCard
}

func (service *cardService) Prepare(actionType string, raceId uint64, family string, userId uint64, isBigRace bool) (error, interface{}) {
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
	return errors.New(storage.ErrorForbidden), nil
}

func (service *cardService) Accept(actionType string, raceId uint64, family string, userId uint64, isBigRace bool) (error, interface{}) {
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

func (service *cardService) Purchase(actionType string, raceId uint64, userId uint64, isBigRace bool, dto dto.CardPurchaseActionDTO) (error, interface{}) {
	logger.Info("CardService.Purchase", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	var err error
	var response interface{}

	switch actionType {
	case "business":
		err = service.raceService.BusinessAction(raceId, userId, isBigRace, dto)
		break

	case "realEstate":
		err = service.raceService.RealEstateAction(raceId, userId, isBigRace, dto)
		break

	case "other":
		err = service.raceService.OtherAssetsAction(raceId, userId, dto.Count)
		break

	case "dream":
		err = service.raceService.DreamAction(raceId, userId)
		break

	case "stock":
		err = service.raceService.StocksAction(raceId, userId, dto.Count)
		break

	case "lottery", "riskBusiness", "riskStock":
		err, response = service.raceService.LotteryAction(raceId, userId, isBigRace)
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

func (service *cardService) Selling(actionType string, raceId uint64, userId uint64, isBigRace bool, dto dto.CardSellingActionDTO) (error, interface{}) {
	logger.Info("CardService.Selling", map[string]interface{}{
		"raceId":     raceId,
		"userId":     userId,
		"actionType": actionType,
		"dto":        dto,
	})

	var err error
	var response interface{}

	switch actionType {
	case "realEstate":
		if dto.ID != "" {
			return errors.New(storage.ErrorIsNotValidRealEstate), nil
		}

		err = service.raceService.SellRealEstate(raceId, userId, dto.ID)
		break
	case "business":
		if dto.ID != "" {
			return errors.New(storage.ErrorIsNotValidBusiness), nil
		}

		err = service.raceService.SellBusiness(raceId, userId, dto.ID, dto.Count)
		break
	case "stock":
		if dto.Count < 1 {
			return errors.New(storage.ErrorIsNotValidCountValue), nil
		}

		err = service.raceService.SellStocks(raceId, userId, dto.Count)
		break
	case "other":
		if dto.ID != "" {
			return errors.New(storage.ErrorIsNotValidOtherAssets), nil
		}

		err = service.raceService.SellOtherAssets(raceId, userId, dto.ID, dto.Count)
		break
	default:
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
		return errors.New(storage.ErrorCannotTakeBigDeals), entity.Card{}
	}

	if action != "" && (action == "small" || action == "big") {
		tile = action + "Deal"
	}

	if tile == "deals" {
		race.CurrentCard = entity.Card{
			ID:      "deal",
			Heading: "Выберите маленькую или большую сделку",
			Family:  "deal",
			Type:    "deal",
		}
	} else {
		if !race.CardMap.HasMapping() {
			cardList := service.GetCards()

			race.CardMap.SetMap(cardList)
		}

		race.CardMap.Next(tile)
		race.CurrentCard = service.getCardByTile(tile, race.CardMap.Active[tile])
	}

	err = service.processCard(action, race, player)

	if err == nil {
		err, _ = service.raceService.UpdateRace(&race)
	}

	return err, race.CurrentCard
}

func (service *cardService) processCard(action string, race entity.Race, player entity.Player) error {
	if action == "market" {
		cardBusinessMarket := entity.CardMarketBusiness{}
		cardBusinessMarket.Fill(race.CurrentCard)

		if race.CurrentCard.ApplicableToEveryOne {
			players := service.playerService.GetAllPlayersByRaceId(race.ID)

			for _, pl := range players {
				if race.CurrentCard.Type == "business" {
					return service.playerService.MarketBusiness(cardBusinessMarket, pl)
				}
			}
		} else {
			if race.CurrentCard.Type == "business" {
				return service.playerService.MarketBusiness(cardBusinessMarket, player)
			}
		}
	}

	return nil
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

func (service *cardService) getCardByTile(cardType string, currentPosition int) entity.Card {
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

	if helper.Contains[string](validTypes, cardType) {
		cardList := service.GetCards()

		if currentPosition < 0 {
			currentPosition = helper.Random(len(cardList[cardType]) - 1)
		}

		card := cardList[cardType][currentPosition]
		card.ID = helper.Uuid(cardType)
		card.Family = service.getFamily(deals, cardType)
		card.Name = cardType
		return card
	}

	return entity.Card{}
}

func (service *cardService) getFamily(deals []string, dealType string) string {
	logger.Info("CardService.getFamily", map[string]interface{}{
		"dealType": dealType,
	})

	if helper.Contains[string](deals, dealType) {
		dealType = "deal"
	}

	return dealType
}

func (service *cardService) getBigCardType(tilePosition int) string {
	logger.Info("CardService.getBigCardType", map[string]interface{}{
		"cardType": tilePosition,
	})

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
	logger.Info("CardService.getRatCardType", map[string]interface{}{
		"cardType": tilePosition,
	})

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

	if helper.Contains[string]([]string{}, cardType) {
		return entity.Card{}
	}

	cardList := service.GetCards()

	return cardList[cardType][helper.Random(len(cardList[cardType])-1)]
}
