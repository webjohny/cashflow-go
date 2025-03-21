package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/storage"
	"log"
	"strconv"
)

type CardService interface {
	SetCards(body dto.CreateCardsDTO)
	Prepare(actionType string, raceId uint64, family string, userId uint64, data dto.PrepareCardBodyDTO) (error, interface{})
	Accept(actionType string, raceId uint64, family string, userId uint64, isBigRace bool) (error, interface{})
	Purchase(actionType string, raceId uint64, userId uint64, isBigRace bool, dto dto.CardPurchaseActionDTO) (error, interface{})
	Selling(actionType string, raceId uint64, userId uint64, isBigRace bool, dto dto.CardSellingActionDTO) (error, interface{})
	Skip(raceId uint64, userId uint64, isBigRace bool) (error, interface{})
	GetCard(action string, raceId uint64, userId uint64, cardType string) (error, entity.Card)
	TestCard(action string, raceId uint64, userId uint64, isBigRace bool) (error, entity.Card)
	CheckPayDay(player entity.Player) int
	ProcessCard(race entity.Race) error
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
	gameService   GameService
	raceService   RaceService
	playerService PlayerService
	ratRace       CardRatRace
	bigRace       CardBigRace

	cards map[string]map[string]map[string][]entity.Card
}

func NewCardService(gameService GameService, raceService RaceService, playerService PlayerService) CardService {
	return &cardService{
		gameService:   gameService,
		raceService:   raceService,
		playerService: playerService,
		cards:         make(map[string]map[string]map[string][]entity.Card),
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
	err, race, player := service.raceService.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err, entity.Card{}
	}

	var tile string

	tileName := action

	if !player.OnBigRace {
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

	err, cardList := service.getCards(race)

	if err != nil {
		return err, entity.Card{}
	}

	if !race.CardMap.HasMapping() {
		race.CardMap.SetMap(cardList)
	}

	race.CardMap.Next(tile)

	race.CurrentCard = service.getCardByTile(tile, race.CardMap.Active[tile], cardList)

	if race.CurrentCard.Family == "market" || race.CurrentCard.Type == "stock" {
		race.IsMultiFlow = race.CurrentCard.OnlyYou == false
	} else {
		race.IsMultiFlow = false
	}

	err = service.ProcessCard(race)

	if err == nil {
		err, _ = service.raceService.UpdateRace(&race)
	}

	return err, race.CurrentCard
}

func (service *cardService) CheckPayDay(player entity.Player) int {
	var count int
	var tiles []int

	current := int(player.CurrentPosition)
	last := int(player.LastPosition)

	if current == 0 && last == 0 {
		return 1
	}

	countTiles := 24

	if player.OnBigRace {
		tiles = service.bigRace.Tiles["cashFlowDay"]
		countTiles = 46
	} else {
		tiles = service.ratRace.Tiles["payday"]
	}

	for i := last + 1; ; i++ {
		key := i % countTiles

		if key == 0 {
			key = countTiles
		}

		if helper.Contains[int](tiles, key) {
			count++
		}

		if key == current {
			break
		}
	}

	return count
}

func (service *cardService) Prepare(actionType string, raceId uint64, family string, userId uint64, data dto.PrepareCardBodyDTO) (error, interface{}) {
	logger.Info("CardService.Prepare", map[string]interface{}{
		"raceId":     raceId,
		"family":     family,
		"actionType": actionType,
		"userId":     userId,
	})

	if family == "deal" || family == "market" {
		err, card := service.GetCard(actionType, raceId, userId, data.Type)

		return err, card
	}

	//if actionType == "risk" || actionType == "riskStock" {
	//	err = service.raceService.PreRiskAction(raceId, username, actionType)
	//}
	return errors.New(storage.ErrorForbidden), nil
}

func (service *cardService) Accept(actionType string, raceId uint64, family string, userId uint64, isBigRace bool) (error, interface{}) {
	logger.Info("Accept", map[string]interface{}{
		"raceId":     raceId,
		"family":     family,
		"actionType": actionType,
		"userId":     userId,
	})

	var err error
	var response dto.MessageResponseDto

	if family == "market" && actionType == "damage" {
		err = service.raceService.MarketAction(raceId, userId, actionType)
	} else if family == "charity" || family == "bigCharity" {
		err = service.raceService.CharityAction(raceId, userId, actionType, isBigRace)
	} else if family == "doodad" {
		err = service.raceService.DoodadAction(raceId, userId)
	} else if family == "baby" {
		err, response = service.raceService.BabyAction(raceId, userId)
	} else if family == "downsized" {
		err = service.raceService.DownsizedAction(raceId, userId)
	} else if family == "bankrupt" {
		err = service.raceService.BigBankruptAction(raceId, userId)
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
		err = service.raceService.OtherAssetsAction(raceId, userId, dto)
		break

	case "dream":
		err = service.raceService.DreamAction(raceId, userId)
		break

	case "stock":
		err = service.raceService.StocksAction(raceId, userId, dto.Count)
		break

	case "bigCharity":
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
		if dto.ID == "" {
			return errors.New(storage.ErrorIsNotValidRealEstate), nil
		}

		err = service.raceService.SellRealEstateAction(raceId, userId, dto.ID)
		break
	case "business":
		if dto.ID == "" {
			return errors.New(storage.ErrorIsNotValidBusiness), nil
		}

		err = service.raceService.SellBusinessAction(raceId, userId, dto.ID, dto.Count)
		break
	case "stock":
		if dto.Count < 1 {
			return errors.New(storage.ErrorIsNotValidCountValue), nil
		}

		err = service.raceService.SellStocksAction(raceId, userId, dto.Count)
		break
	case "other":
		if dto.ID == "" {
			return errors.New(storage.ErrorIsNotValidOtherAssets), nil
		}

		err = service.raceService.SellOtherAssetsAction(raceId, userId, dto.ID, dto.Count)
		break
	default:
		err = service.raceService.SkipAction(raceId, userId, isBigRace)
		break
	}

	return err, response
}

func (service *cardService) GetCard(action string, raceId uint64, userId uint64, cardType string) (error, entity.Card) {
	logger.Info("CardService.GetCard", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"action": action,
	})

	err, race, player := service.raceService.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err, entity.Card{}
	}

	var tile string

	if player.OnBigRace {
		tile = service.getBigCardType(int(player.CurrentPosition))
	} else {
		tile = service.getRatCardType(int(player.CurrentPosition))

		if action == "big" && player.Cash < 10000 {
			return errors.New(storage.ErrorCannotTakeBigDeals), entity.Card{}
		}

		if action != "" && (action == "small" || action == "big") {
			tile = action + "Deal"
		}
	}

	var card entity.Card
	err, cardList := service.getCards(race)

	if err != nil {
		return err, entity.Card{}
	}

	if tile == "deals" {
		card = cardList["defaultDeal"][0]
		card.ID = helper.Uuid(card.ID)
	} else if race.Options.EnableCardCategory && tile == "market" && action != "buy" {
		card = cardList["defaultMarket"][0]
		card.ID = helper.Uuid(card.ID)
	} else {
		if !race.CardMap.HasMapping() {
			race.CardMap.SetMap(cardList)
		}

		if tile == "dream" {
			players := service.raceService.GetRacePlayersByRaceId(raceId, false)

			for _, currentPlayer := range players {
				if currentPlayer.Info.Dream.ID == int(player.CurrentPosition) {
					card = entity.Card{
						ID:        "own-dream-" + strconv.Itoa(int(currentPlayer.ID)),
						Heading:   currentPlayer.Info.Dream.Name,
						Type:      "dream",
						AssetType: "personal",
						Cost:      currentPlayer.Info.Dream.Price,
						PlayerId:  int(currentPlayer.ID),
					}
				}
			}
		}

		log.Println("CARD_TYPE", cardType)

		if cardType == "" {
			race.CardMap.Next(tile)
		} else {
			for i := 0; i < len(race.CardMap.Map[tile]); i++ {
				race.CardMap.Next(tile)

				card = service.getCardByTile(tile, race.CardMap.Active[tile], cardList)

				if card.Type == cardType {
					break
				}
			}
		}

		if card.ID == "" {
			card = service.getCardByTile(tile, race.CardMap.Active[tile], cardList)
		}
	}

	if card.Family == "market" || card.Type == "stock" {
		race.IsMultiFlow = card.OnlyYou == false
	} else {
		race.IsMultiFlow = false
	}

	race.CurrentCard = card

	err = service.ProcessCard(race)

	logger.Info("CardService.GetCard", map[string]interface{}{
		"raceId":   raceId,
		"playerId": player.ID,
		"card":     helper.JsonSerialize(card),
	})

	if err == nil {
		err, _ = service.raceService.UpdateRace(&race)
	}

	return err, race.CurrentCard
}

func (service *cardService) SetCards(body dto.CreateCardsDTO) {
	if service.cards[body.Type] == nil {
		service.cards[body.Type] = make(map[string]map[string][]entity.Card)
	}
	service.cards[body.Type][body.Language] = body.Cards
}

func (service *cardService) ProcessCard(race entity.Race) error {
	card := race.CurrentCard

	if card.Name == "smallDeal" {
		if card.AssetType == entity.StockTypes.Manipulation {
			players := service.playerService.GetAllPlayersByRaceId(race.ID)
			cardStocks := entity.CardStocks{}
			cardStocks.Fill(card)

			for _, player := range players {
				var err error
				if card.Increase > 0 {
					err = service.playerService.IncreaseStocks(cardStocks, player)
				} else if card.Decrease > 0 {
					err = service.playerService.DecreaseStocks(cardStocks, player)
				}
				if err != nil {
					logger.Error(err)
				}
			}
		}
	} else if card.Type == "success" || card.Type == "inflation" {
		players := service.playerService.GetAllPlayersByRaceId(race.ID)

		for _, pl := range players {
			if (race.CurrentCard.OnlyYou && race.CurrentPlayer.ID == pl.ID) || !race.CurrentCard.OnlyYou {
				cardMarket := entity.CardMarket{}
				cardMarket.Fill(race.CurrentCard)
				err := service.playerService.MarketManipulation(cardMarket, pl, players)

				if err != nil {
					logger.Error(err)
				}
			}
		}
	}

	return nil
}

func (service *cardService) getCards(race entity.Race) (error, map[string][]entity.Card) {
	logger.Info("CardService.getCards", race.Options.CardCollection, race.Options.Language)

	if race.Options.CardCollection == "" {
		race.Options.CardCollection = "default"
	}

	cards := service.cards[race.Options.CardCollection][race.Options.Language]

	if cards == nil {
		return errors.New(storage.ErrorCardsNotFound), make(map[string][]entity.Card)
	}

	return nil, cards
}

func (service *cardService) getCardByTile(cardType string, currentPosition int, cardList map[string][]entity.Card) entity.Card {
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
		"cashFlowDay",
		"bigCharity",
		"business",
		"dream",
		"tax50percent",
		"tax100percent",
		"bankrupt",
	)

	if helper.Contains[string](validTypes, cardType) {
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

	for tile, positions := range service.bigRace.Tiles {
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
