package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/service"
	"net/http"
	"sort"
	"strconv"
)

type PlayerTestController interface {
	AddMoney(ctx *gin.Context)
	MarketManipulation(ctx *gin.Context)
	DamageRealEstate(ctx *gin.Context)
	SellStocks(ctx *gin.Context)
	SellBusiness(ctx *gin.Context)
	SellRealEstate(ctx *gin.Context)
	SellOtherAssets(ctx *gin.Context)
	DecreaseStocks(ctx *gin.Context)
	IncreaseStocks(ctx *gin.Context)
	BuyRealEstate(ctx *gin.Context)
	BuyStocks(ctx *gin.Context)
	BuyBusiness(ctx *gin.Context)
	BuyBigBusiness(ctx *gin.Context)
	BuyLottery(ctx *gin.Context)
	BuyOtherAssets(ctx *gin.Context)
	BuyOtherAssetsInPartnership(ctx *gin.Context)
	BuyDream(ctx *gin.Context)
	BuyRiskStocks(ctx *gin.Context)
	BuyRiskBusiness(ctx *gin.Context)
	BuyRealEstateInPartnership(ctx *gin.Context)
	BuyBusinessInPartnership(ctx *gin.Context)
	Index(ctx *gin.Context)
}

type PlayerResponse struct {
	ID               uint64
	UserID           uint64
	Card             interface{}
	Extra            interface{} `json:"Extra,omitempty"`
	OldCash          int
	NewCash          int
	NewPassiveIncome int         `json:"NewPassiveIncome,omitempty"`
	OldCashFlow      int         `json:"OldCashFlow,omitempty"`
	NewCashFlow      int         `json:"NewCashFlow,omitempty"`
	SingleAsset      interface{} `json:"SingleAsset,omitempty"`
	Assets           interface{} `json:"Assets,omitempty"`
}

type playerTestController struct {
	playerService service.PlayerService
}

func NewPlayerTestController(playerService service.PlayerService) PlayerTestController {
	return &playerTestController{
		playerService: playerService,
	}
}

var Links = map[string]map[string]string{
	"Покупка": {
		"Акции":        "/test/player/buy-stocks",
		"Недвижимость": "/test/player/buy-real-estate",
		"Недвижимость (сделка с кем-то)": "/test/player/buy-partner-real-estate",
		"Бизнес": "/test/player/buy-business",
		"Бизнес (лимитированное партнёрство)":                   "/test/player/buy-business?type=limited",
		"Бизнес (лимитированное партнёрство - сделка с кем-то)": "/test/player/buy-partner-business?type=limited",
		"Бизнес (сделка с кем-то)":                              "/test/player/buy-partner-business",
		"Лотерея":       "/test/player/buy-lottery",
		"Другие активы": "/test/player/buy-other-assets",
		"Другие активы (сделка с кем-то)":   "/test/player/buy-partner-other-assets",
		"Большой круг - Бизнес":             "/test/player/buy-big-business",
		"Большой круг - Рискованный бизнес": "/test/player/buy-risk-business",
		"Большой круг - Рискованные акции":  "/test/player/buy-risk-stocks",
		"Большой круг - Мечта":              "/test/player/buy-dream",
	},
	"Продажа": {
		"Акции":                  "/test/player/sell-stocks",
		"Недвижимость":           "/test/player/sell-real-estate?type=single",
		"Большая недвижимость":   "/test/player/sell-real-estate",
		"Бизнес":                 "/test/player/sell-business",
		"Бизнес огран. партнёр.": "/test/player/sell-business?type=limited",
		"Другие активы (монеты)": "/test/player/sell-other-assets",
		"Другие активы (земля)":  "/test/player/sell-other-assets?type=whole",
	},
	"Манипуляции": {
		"Удвоение кол-во акций":                       "/test/player/increase-stocks",
		"Разделение (на 3) кол-во акций":              "/test/player/decrease-stocks",
		"Недвижимость повреждена":                     "/test/player/damage-real-estate",
		"Бизнес идёт вверх (прибавка к пасс. доходу)": "/test/player/market-manipulation?type=success",
		//"Бизнес идёт вниз (убавление пасс. дохода)":   "/test/player/market-manipulation?type=failure",
		"Удары инфляции": "/test/player/market-manipulation?type=inflation",
	},
}

func (c *playerTestController) AddMoney(ctx *gin.Context) {
	b := c.getPlayer(ctx)

	cashQuery := ctx.Query("cash")
	cash, _ := strconv.Atoi(cashQuery)

	c.addCashForPlayer(&b, cash, false)
}

func (c *playerTestController) MarketManipulation(ctx *gin.Context) {
	c.cleaning(ctx)
	b := c.getPlayer(ctx)

	card := c.getCard(ctx, "inflation").(entity.CardMarket)

	err := c.playerService.MarketManipulation(card, b)

	if err != nil {
		logger.Error(err, nil)

		request.FinalResponse(ctx, err, nil)
		return
	}

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

	request.FinalResponse(ctx, err, PlayerResponse{
		ID:               player.ID,
		UserID:           player.UserID,
		Card:             card,
		OldCash:          b.Cash,
		NewCash:          player.Cash,
		OldCashFlow:      b.CalculateCashFlow(),
		NewCashFlow:      player.CalculateCashFlow(),
		NewPassiveIncome: player.CalculatePassiveIncome(),
		Assets:           player.Assets,
	})
}

func (c *playerTestController) DamageRealEstate(ctx *gin.Context) {
	c.cleaning(ctx)
	b := c.getPlayer(ctx)

	card := entity.CardMarket{
		ID:          helper.Uuid("damage"),
		Type:        "damage",
		Heading:     "Арендатор повредил вашу собственность",
		Symbol:      "ANY",
		Description: "Потеряв работу арендатор отказался платить и скрылся, заплатите $1000",
		AssetType:   entity.MarketTypes.AnyRealEstate,
		Cost:        1000,
		OnlyYou:     false,
	}

	c.addCashForPlayer(&b, card.Cost, false)

	err := c.playerService.MarketDamage(card, b)

	if err != nil {
		logger.Error(err, nil)

		request.FinalResponse(ctx, err, nil)
		return
	}

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

	request.FinalResponse(ctx, err, PlayerResponse{
		ID:               player.ID,
		UserID:           player.UserID,
		Card:             card,
		OldCash:          b.Cash,
		NewCash:          player.Cash,
		OldCashFlow:      b.CalculateCashFlow(),
		NewCashFlow:      player.CalculateCashFlow(),
		NewPassiveIncome: player.CalculatePassiveIncome(),
		Assets:           player.Assets,
	})
}

func (c *playerTestController) BuyBigBusiness(ctx *gin.Context) {
	c.cleaning(ctx)
	b := c.getPlayer(ctx)

	card := entity.CardBusiness{
		ID:          helper.Uuid("business"),
		Type:        "business",
		Symbol:      "pizzaFranchise",
		Heading:     "Франчайзинг пиццерий",
		Description: "Откройте пиццерию в своём родном городе",
		Cost:        125000,
		CashFlow:    6000,
		IsOwner:     true,
	}
	c.addCashForPlayer(&b, card.Cost, true)

	err := c.playerService.BuyBusiness(card, b, 1, true)

	if err != nil {
		logger.Error(err, nil)

		request.FinalResponse(ctx, err, nil)
		return
	}

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

	_, business := player.FindBusinessBySymbol(card.Symbol)

	request.FinalResponse(ctx, err, PlayerResponse{
		ID:               player.ID,
		UserID:           player.UserID,
		Card:             card,
		OldCash:          b.Cash,
		NewCash:          player.Cash,
		OldCashFlow:      b.CalculateCashFlow(),
		NewCashFlow:      player.CalculateCashFlow(),
		NewPassiveIncome: player.CalculatePassiveIncome(),
		SingleAsset:      business,
		Assets:           player.Assets.Business,
	})
}

func (c *playerTestController) BuyLottery(ctx *gin.Context) {
	c.cleaning(ctx)
	b := c.getPlayer(ctx)

	card := entity.CardLottery{
		ID:          helper.Uuid("lottery"),
		Type:        "lottery",
		Symbol:      "lottery",
		Heading:     "Sister-In-Law borrows Money",
		Description: "Sister-in-law is downsized.Needs $5,000 to make house payment.",
		Cost:        1000,
		AssetType:   entity.LotteryTypes.Money,
		Rule:        "If you choose to help. Pay $5,000 and roll 1 die:",
		SubRule: []string{
			"Die = 1-3, She never pays you back and you're out $5,000",
			"Die = 4-6, She pays you back $10,000, but family get-togethers are still awkward",
		},
		Failure: []int{1, 2, 3},
		Success: []int{4, 5, 6},
		Outcome: entity.CardLotteryOutcome{
			Failure: 0,
			Success: 5000,
		},
	}

	c.addCashForPlayer(&b, card.Cost, false)

	err, result := c.playerService.BuyLottery(card, b, helper.Random(6))

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

	request.FinalResponse(ctx, err, PlayerResponse{
		ID:      player.ID,
		UserID:  player.UserID,
		Extra:   result,
		Card:    card,
		OldCash: b.Cash,
		NewCash: player.Cash,
	})
}

func (c *playerTestController) BuyOtherAssets(ctx *gin.Context) {
	c.cleaning(ctx)
	b := c.getPlayer(ctx)

	card := c.getCard(ctx, entity.OtherAssetTypes.Piece).(entity.CardOtherAssets)

	c.addCashForPlayer(&b, card.WholeCost, false)

	err := c.playerService.BuyOtherAssets(card, b, 5)

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

	request.FinalResponse(ctx, err, PlayerResponse{
		ID:               b.ID,
		UserID:           b.UserID,
		Card:             card,
		OldCash:          b.Cash,
		NewCash:          player.Cash,
		OldCashFlow:      b.CalculateCashFlow(),
		NewCashFlow:      player.CalculateCashFlow(),
		NewPassiveIncome: player.CalculatePassiveIncome(),
		Assets:           player.Assets.OtherAssets,
	})
}

func (c *playerTestController) BuyOtherAssetsInPartnership(ctx *gin.Context) {
	c.cleaning(ctx)

	card := c.getCard(ctx, entity.OtherAssetTypes.Piece).(entity.CardOtherAssets)

	raceIDNum, _ := strconv.Atoi(ctx.DefaultQuery("raceID", "30"))
	raceID := uint64(raceIDNum)

	players := c.playerService.GetAllPlayersByRaceId(raceID)

	c.addCashForPlayer(&players[1], card.WholeCost, true)

	var parts []dto.CardPurchasePlayerActionDTO

	if card.AssetType == entity.OtherAssetTypes.Piece {
		parts = c.fillAmounts(players, card.Count, 1, 5, "amount")
	} else {
		parts = c.fillAmounts(players, card.Cost, 10, card.Cost, "amount")
	}

	err := c.playerService.BuyOtherAssetsInPartnership(card, players[1], players, parts)

	if err != nil {
		logger.Error(err, nil)

		request.FinalResponse(ctx, err, nil)
		return
	}

	newPlayers := c.playerService.GetAllPlayersByRaceId(raceID)

	responsePlayers := make([]PlayerResponse, 0)

	for i, player := range newPlayers {
		_, asset := player.FindOtherAssetsBySymbol(card.Symbol)

		responsePlayers = append(responsePlayers, PlayerResponse{
			ID:          player.ID,
			UserID:      player.UserID,
			Card:        card,
			OldCash:     players[i].Cash,
			NewCash:     player.Cash,
			OldCashFlow: players[i].CalculatePassiveIncome(),
			NewCashFlow: player.CalculatePassiveIncome(),
			SingleAsset: asset,
			Assets:      player.Assets.OtherAssets,
		})
	}

	request.FinalResponse(ctx, err, responsePlayers)
}

func (c *playerTestController) BuyRealEstateInPartnership(ctx *gin.Context) {
	c.cleaning(ctx)
	typeOf := ctx.DefaultQuery("type", entity.RealEstateTypes.Building)

	raceIDNum, _ := strconv.Atoi(ctx.DefaultQuery("raceID", "30"))
	raceID := uint64(raceIDNum)

	players := c.playerService.GetAllPlayersByRaceId(raceID)

	card := c.getCard(ctx, entity.RealEstateTypes.Building).(entity.CardRealEstate)

	parts := []dto.CardPurchasePlayerActionDTO{
		{
			ID:      int(players[0].ID),
			Amount:  1000,
			Passive: 150,
		},
		{
			ID:      int(players[1].ID),
			Amount:  1000,
			Passive: 150,
		},
	}

	if typeOf == entity.RealEstateTypes.Building {
		parts[0].Amount = 40000
		parts[0].Passive = 1000
		parts[1].Amount = 40000
		parts[1].Passive = 2000
	}

	for _, player := range players {
		c.addCashForPlayer(&player, card.Cost, false)
	}

	err := c.playerService.BuyRealEstateInPartnership(card, players[1], players, parts)

	if err != nil {
		logger.Error(err, nil)

		request.FinalResponse(ctx, err, nil)
		return
	}

	newPlayers := c.playerService.GetAllPlayersByRaceId(raceID)

	responsePlayers := make([]PlayerResponse, 0)

	for key, player := range newPlayers {
		var realEstate entity.CardRealEstate

		if len(player.Assets.RealEstates) > 0 {
			realEstate = player.Assets.RealEstates[0]
		}

		responsePlayers = append(responsePlayers, PlayerResponse{
			ID:          player.ID,
			Card:        card,
			OldCash:     players[key].Cash,
			NewCash:     player.Cash,
			OldCashFlow: players[key].CalculatePassiveIncome(),
			NewCashFlow: player.CalculatePassiveIncome(),
			SingleAsset: realEstate,
			Assets:      player.Assets.RealEstates,
		})
	}

	request.FinalResponse(ctx, err, responsePlayers)
}

func (c *playerTestController) BuyBusinessInPartnership(ctx *gin.Context) {
	c.cleaning(ctx)

	raceIDNum, _ := strconv.Atoi(ctx.DefaultQuery("raceID", "30"))
	raceID := uint64(raceIDNum)

	players := c.playerService.GetAllPlayersByRaceId(raceID)

	card := c.cardBusiness(ctx, &players[1])

	var parts []dto.CardPurchasePlayerActionDTO

	if card.AssetType == entity.BusinessTypes.Limited {
		parts = c.fillAmounts(players, card.Limit, 1, 5, "amount")
	} else {
		//parts = c.fillAmounts(players, card.CashFlow, 10, card.CashFlow/2, "passive")
		parts = append(parts, dto.CardPurchasePlayerActionDTO{
			ID:      int(players[0].ID),
			Passive: 112,
		})
		parts = append(parts, dto.CardPurchasePlayerActionDTO{
			ID:      int(players[1].ID),
			Passive: 300 - 112,
		})
	}

	logger.Warn(parts)

	err := c.playerService.BuyBusinessInPartnership(card, players[1], players, parts)

	if err != nil {
		logger.Error(err, nil)

		request.FinalResponse(ctx, err, nil)
		return
	}

	newPlayers := c.playerService.GetAllPlayersByRaceId(raceID)

	responsePlayers := make([]PlayerResponse, 0)

	for i, player := range newPlayers {
		_, business := player.FindBusinessBySymbol(card.Symbol)

		responsePlayers = append(responsePlayers, PlayerResponse{
			ID:          player.ID,
			UserID:      player.UserID,
			Card:        card,
			OldCash:     players[i].Cash,
			NewCash:     player.Cash,
			OldCashFlow: players[i].CalculatePassiveIncome(),
			NewCashFlow: player.CalculatePassiveIncome(),
			SingleAsset: business,
			Assets:      player.Assets.Business,
		})
	}

	request.FinalResponse(ctx, err, responsePlayers)
}

func (c *playerTestController) BuyStocks(ctx *gin.Context) {
	c.cleaning(ctx)
	b := c.getPlayer(ctx)

	card := entity.CardStocks{
		ID:          helper.Uuid("life10u40"),
		Type:        "stock",
		Symbol:      "Lifecell",
		Heading:     "Lifecell оператор мобильной связи",
		Description: "Research and development delays cause low share price for this longtime pharmaceutical company",
		Rule:        "Only you may buy as many shares as you want at this price. Everyone may sell at this price.",
		Price:       10,
		OnlyYou:     true,
		Count:       100,
		Range:       []int{5, 30},
	}

	c.addCashForPlayer(&b, card.Price*card.Count, false)

	err := c.playerService.BuyStocks(card, b, true)

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

	request.FinalResponse(ctx, err, PlayerResponse{
		ID:               b.ID,
		UserID:           b.UserID,
		Card:             card,
		OldCash:          b.Cash,
		NewCash:          player.Cash,
		OldCashFlow:      b.CalculateCashFlow(),
		NewCashFlow:      player.CalculateCashFlow(),
		NewPassiveIncome: player.CalculatePassiveIncome(),
		Assets:           player.Assets.Stocks,
	})
}

func (c *playerTestController) SellStocks(ctx *gin.Context) {
	c.cleaning(ctx)
	b := c.getPlayer(ctx)

	card := entity.CardStocks{
		ID:          helper.Uuid("life10u40"),
		Type:        "stock",
		Symbol:      "Lifecell",
		Heading:     "Lifecell оператор мобильной связи",
		Description: "Research and development delays cause low share price for this longtime pharmaceutical company",
		Rule:        "Only you may buy as many shares as you want at this price. Everyone may sell at this price.",
		Price:       30,
		OnlyYou:     true,
		Range:       []int{5, 30},
	}

	err := c.playerService.SellStocks(card, b, 2100, true)

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

	request.FinalResponse(ctx, err, PlayerResponse{
		ID:               b.ID,
		UserID:           b.UserID,
		Card:             card,
		OldCash:          b.Cash,
		NewCash:          player.Cash,
		OldCashFlow:      b.CalculateCashFlow(),
		NewCashFlow:      player.CalculateCashFlow(),
		NewPassiveIncome: player.CalculatePassiveIncome(),
		Assets:           player.Assets.Stocks,
	})
}

func (c *playerTestController) SellBusiness(ctx *gin.Context) {
	c.cleaning(ctx)
	b := c.getPlayer(ctx)

	typeOfBusiness := ctx.DefaultQuery("type", entity.BusinessTypes.Startup)

	_, business := b.FindBusinessBySymbol(typeOfBusiness)

	if business.ID == "" {
		card := c.cardBusiness(ctx, &b)
		card.IsOwner = true

		err := c.playerService.BuyBusiness(card, b, 0, true)

		if err != nil {
			request.FinalResponse(ctx, err, business)
			return
		}

		err, b = c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

		_, business = b.FindBusinessBySymbol(typeOfBusiness)
	}

	var card entity.CardMarketBusiness
	var count int

	if typeOfBusiness == entity.BusinessTypes.Startup {
		card = entity.CardMarketBusiness{
			ID:          helper.Uuid("itCompanyMarket"),
			Type:        "business",
			AssetType:   entity.BusinessTypes.Startup,
			Symbol:      "startup",
			Heading:     "Покупка IT компании",
			Description: "Крупная компания по производству ПО предлагает 100,000$ за не-большую компанию.",
			Rule:        "Only you may buy as many shares as you want at this price. Everyone may sell at this price.",
			Cost:        100000,
		}
	} else if typeOfBusiness == entity.BusinessTypes.Limited {
		card = entity.CardMarketBusiness{
			ID:          helper.Uuid("itCompanyMarket"),
			Type:        "business",
			AssetType:   entity.BusinessTypes.Limited,
			Symbol:      "limited",
			Heading:     "Покупка партнёрства",
			Description: "Крупная компания предлагает тройную стоимость за каждую долю.",
			Rule:        "Only you may buy as many shares as you want at this price. Everyone may sell at this price.",
			Cost:        3,
		}

		count = 2
	}

	err, result := c.playerService.SellBusiness(business.ID, card, b, count)

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

	request.FinalResponse(ctx, err, PlayerResponse{
		ID:               b.ID,
		Extra:            result,
		UserID:           b.UserID,
		Card:             card,
		OldCash:          b.Cash,
		NewCash:          player.Cash,
		OldCashFlow:      b.CalculateCashFlow(),
		NewCashFlow:      player.CalculateCashFlow(),
		NewPassiveIncome: player.CalculatePassiveIncome(),
		Assets:           player.Assets.Business,
	})
}

func (c *playerTestController) SellRealEstate(ctx *gin.Context) {
	c.cleaning(ctx)
	b := c.getPlayer(ctx)

	typeOfRealEstate := ctx.DefaultQuery("type", "building")

	_, realEstate := b.FindRealEstateBySymbol(typeOfRealEstate)

	if realEstate.ID == "" {
		card := c.getCard(ctx, entity.RealEstateTypes.Building).(entity.CardRealEstate)

		c.addCashForPlayer(&b, card.Cost, false)

		card.IsOwner = true

		err := c.playerService.BuyRealEstate(card, b)

		if err != nil {
			request.FinalResponse(ctx, err, realEstate)
			return
		}

		err, b = c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

		_, realEstate = b.FindRealEstateBySymbol(typeOfRealEstate)
	}

	var card entity.CardMarketRealEstate

	if typeOfRealEstate == "building" {
		card = entity.CardMarketRealEstate{
			ID:          helper.Uuid("marketRealEstate"),
			Type:        "realEstate",
			Heading:     "Покупатель на большой многоквартирный дом",
			Symbol:      "building",
			AssetType:   entity.RealEstateTypes.Building,
			Description: "Покупатель предлагает за каждую квартиру $40,000",
			Cost:        40000,
			Range:       []int{2, 4, 8, 12, 24},
		}
	} else if typeOfRealEstate == "single" {
		card = entity.CardMarketRealEstate{
			ID:          helper.Uuid("marketRealEstate"),
			Type:        "realEstate",
			Heading:     "Покупатель на 3КВ",
			Symbol:      "single",
			AssetType:   entity.RealEstateTypes.Single,
			Description: "Покупатель предлагает $100,000 за 3х комнатную квартиру",
			Cost:        100000,
		}
	}

	err := c.playerService.SellRealEstate(realEstate.ID, card, b)

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

	request.FinalResponse(ctx, err, PlayerResponse{
		ID:               b.ID,
		UserID:           b.UserID,
		Card:             card,
		OldCash:          b.Cash,
		NewCash:          player.Cash,
		OldCashFlow:      b.CalculateCashFlow(),
		NewCashFlow:      player.CalculateCashFlow(),
		NewPassiveIncome: player.CalculatePassiveIncome(),
		Assets:           player.Assets.Business,
	})
}

func (c *playerTestController) SellOtherAssets(ctx *gin.Context) {
	c.cleaning(ctx)
	b := c.getPlayer(ctx)

	typeOf := ctx.DefaultQuery("type", entity.OtherAssetTypes.Piece)

	_, otherAsset := b.FindOtherAssetsBySymbol(typeOf)

	if otherAsset.ID == "" {
		card := c.getCard(ctx, entity.OtherAssetTypes.Piece).(entity.CardOtherAssets)

		c.addCashForPlayer(&b, card.Cost, false)

		card.IsOwner = true

		var count int

		if typeOf == entity.OtherAssetTypes.Piece {
			count = 5
		}
		err := c.playerService.BuyOtherAssets(card, b, count)

		if err != nil {
			request.FinalResponse(ctx, err, otherAsset)
			return
		}

		err, b = c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

		_, otherAsset = b.FindOtherAssetsBySymbol(typeOf)
	}

	var card entity.CardMarketOtherAssets

	if typeOf == entity.OtherAssetTypes.Piece {
		card = entity.CardMarketOtherAssets{
			ID:          helper.Uuid("marketGoldCoins"),
			Type:        "other",
			Heading:     "Покупатель на золотые монеты",
			Symbol:      "piece",
			AssetType:   entity.OtherAssetTypes.Piece,
			Description: "Покупатель предлагает за каждую монету по $5,000",
			Cost:        5000,
		}
	} else if typeOf == entity.OtherAssetTypes.Whole {
		card = entity.CardMarketOtherAssets{
			ID:          helper.Uuid("marketLand"),
			Type:        "other",
			Heading:     "Покупатель на 20 акров земли",
			Symbol:      "whole",
			AssetType:   entity.OtherAssetTypes.Whole,
			Description: "Покупатель предлагает $200,000 за 20 акров земли",
			Cost:        200000,
		}
	}

	err := c.playerService.SellOtherAssets(otherAsset.ID, card, b, 5)

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

	request.FinalResponse(ctx, err, PlayerResponse{
		ID:               b.ID,
		UserID:           b.UserID,
		Card:             card,
		OldCash:          b.Cash,
		NewCash:          player.Cash,
		OldCashFlow:      b.CalculateCashFlow(),
		NewCashFlow:      player.CalculateCashFlow(),
		NewPassiveIncome: player.CalculatePassiveIncome(),
		Assets:           player.Assets.OtherAssets,
	})
}

func (c *playerTestController) IncreaseStocks(ctx *gin.Context) {
	c.cleaning(ctx)
	b := c.getPlayer(ctx)

	card := entity.CardStocks{
		ID:          helper.Uuid("life10u40"),
		Type:        "stock",
		Symbol:      "Lifecell",
		Heading:     "Lifecell оператор мобильной связи",
		Description: "Research and development delays cause low share price for this longtime pharmaceutical company",
		Rule:        "Only you may buy as many shares as you want at this price. Everyone may sell at this price.",
		Increase:    2,
	}

	err := c.playerService.IncreaseStocks(card, b)

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

	request.FinalResponse(ctx, err, PlayerResponse{
		ID:               b.ID,
		UserID:           b.UserID,
		Card:             card,
		OldCash:          b.Cash,
		NewCash:          player.Cash,
		OldCashFlow:      b.CalculateCashFlow(),
		NewCashFlow:      player.CalculateCashFlow(),
		NewPassiveIncome: player.CalculatePassiveIncome(),
		Assets:           player.Assets.Stocks,
	})
}

func (c *playerTestController) DecreaseStocks(ctx *gin.Context) {
	c.cleaning(ctx)
	b := c.getPlayer(ctx)

	card := entity.CardStocks{
		ID:          helper.Uuid("life10u40"),
		Type:        "stock",
		Symbol:      "Lifecell",
		Heading:     "Lifecell оператор мобильной связи",
		Description: "Research and development delays cause low share price for this longtime pharmaceutical company",
		Rule:        "Only you may buy as many shares as you want at this price. Everyone may sell at this price.",
		Decrease:    3,
	}

	err := c.playerService.DecreaseStocks(card, b)

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

	request.FinalResponse(ctx, err, PlayerResponse{
		ID:               b.ID,
		UserID:           b.UserID,
		Card:             card,
		OldCash:          b.Cash,
		NewCash:          player.Cash,
		OldCashFlow:      b.CalculateCashFlow(),
		NewCashFlow:      player.CalculateCashFlow(),
		NewPassiveIncome: player.CalculatePassiveIncome(),
		Assets:           player.Assets.Stocks,
	})
}

func (c *playerTestController) BuyRealEstate(ctx *gin.Context) {
	c.cleaning(ctx)
	b := c.getPlayer(ctx)

	card := c.getCard(ctx, entity.RealEstateTypes.Building).(entity.CardRealEstate)

	c.addCashForPlayer(&b, card.DownPayment*5, false)

	err := c.playerService.BuyRealEstate(card, b)

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

	request.FinalResponse(ctx, err, PlayerResponse{
		ID:               b.ID,
		UserID:           b.UserID,
		Card:             card,
		OldCash:          b.Cash,
		NewCash:          player.Cash,
		OldCashFlow:      b.CalculateCashFlow(),
		NewCashFlow:      player.CalculateCashFlow(),
		NewPassiveIncome: player.CalculatePassiveIncome(),
		Assets:           player.Assets.RealEstates,
	})
}

func (c *playerTestController) BuyDream(ctx *gin.Context) {
	c.cleaning(ctx)
	b := c.getPlayer(ctx)

	card := entity.CardDream{
		ID:          helper.Uuid("dreamParkAttraction1"),
		Heading:     "Парк аттракционов в Лондоне",
		Description: "Вы купили супер парк развлечений для своих детей",
		Type:        "dream",
		Cost:        100000,
	}

	c.addCashForPlayer(&b, card.Cost, false)

	err := c.playerService.BuyDream(card, b)

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

	request.FinalResponse(ctx, err, PlayerResponse{
		ID:               b.ID,
		UserID:           b.UserID,
		Card:             card,
		OldCash:          b.Cash,
		NewCash:          player.Cash,
		OldCashFlow:      b.CalculateCashFlow(),
		NewCashFlow:      player.CalculateCashFlow(),
		NewPassiveIncome: player.CalculatePassiveIncome(),
		Assets:           player.Assets.Dreams,
	})
}

func (c *playerTestController) BuyRiskStocks(ctx *gin.Context) {
	c.cleaning(ctx)
	b := c.getPlayer(ctx)

	card := entity.CardLottery{
		ID:          helper.Uuid("lottery"),
		Type:        "riskStock",
		Symbol:      "lottery",
		Heading:     "Рискованные акции",
		Description: "Крупнейшая компания распродаёт свои акции по 0.1$, купи 100,000 акций и сможешь забрать 100,000$",
		Cost:        10000,
		AssetType:   entity.LotteryTypes.Money,
		Rule:        "If you choose to help. Pay $10,000 and roll 1 die:",
		SubRule: []string{
			"Кубик = 1-3, Ты проиграл $10,000",
			"Кубик = 4-6, Ты увеличил в 10 раз прибыль",
		},
		Failure: []int{1, 2, 3},
		Success: []int{4, 5, 6},
		Outcome: entity.CardLotteryOutcome{
			Failure: 0,
			Success: 100000,
		},
	}

	c.addCashForPlayer(&b, card.Cost, false)

	err, result := c.playerService.BuyLottery(card, b, helper.Random(6))

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

	request.FinalResponse(ctx, err, PlayerResponse{
		ID:               b.ID,
		Extra:            result,
		UserID:           b.UserID,
		Card:             card,
		OldCash:          b.Cash,
		NewCash:          player.Cash,
		OldCashFlow:      b.CalculateCashFlow(),
		NewCashFlow:      player.CalculateCashFlow(),
		NewPassiveIncome: player.CalculatePassiveIncome(),
		Assets:           player.Assets.Stocks,
	})
}

func (c *playerTestController) BuyRiskBusiness(ctx *gin.Context) {
	c.cleaning(ctx)
	b := c.getPlayer(ctx)

	card := entity.CardLottery{
		ID:          helper.Uuid("lottery"),
		Type:        "riskBusiness",
		Symbol:      "lottery",
		Heading:     "Рискованный бизнес",
		Description: "Заплати $300,000 и получишь пассивный доход в $75,000",
		Cost:        300000,
		AssetType:   entity.LotteryTypes.CashFlow,
		Rule:        "If you choose to help. Pay $10,000 and roll 1 die:",
		SubRule: []string{
			"Кубик = 1-3, Ты проиграл $300,000",
			"Кубик = 4-6, Ты увеличил пассивный доход на $75,000",
		},
		Failure: []int{1, 2, 3},
		Success: []int{4, 5, 6},
		Outcome: entity.CardLotteryOutcome{
			Failure: 0,
			Success: 75000,
		},
	}

	c.addCashForPlayer(&b, card.Cost, false)

	err, result := c.playerService.BuyLottery(card, b, helper.Random(6))

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

	request.FinalResponse(ctx, err, PlayerResponse{
		ID:               b.ID,
		Extra:            result,
		UserID:           b.UserID,
		Card:             card,
		OldCash:          b.Cash,
		NewCash:          player.Cash,
		OldCashFlow:      b.CalculateCashFlow(),
		NewCashFlow:      player.CalculateCashFlow(),
		NewPassiveIncome: player.CalculatePassiveIncome(),
		Assets:           player.Assets.Business,
	})
}

func (c *playerTestController) BuyBusiness(ctx *gin.Context) {
	c.cleaning(ctx)

	card := c.cardBusiness(ctx, nil)

	b := c.getPlayer(ctx)

	count := 0

	if card.AssetType == entity.BusinessTypes.Limited {
		count = 3
	}

	card.IsOwner = true

	err := c.playerService.BuyBusiness(card, b, count, true)

	if err != nil {
		logger.Error(err, nil)

		request.FinalResponse(ctx, err, nil)
		return
	}

	err, player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

	_, business := player.FindBusinessBySymbol(card.Symbol)

	request.FinalResponse(ctx, err, PlayerResponse{
		ID:               player.ID,
		UserID:           player.UserID,
		Card:             card,
		OldCash:          b.Cash,
		NewCash:          player.Cash,
		OldCashFlow:      b.CalculateCashFlow(),
		NewCashFlow:      player.CalculateCashFlow(),
		NewPassiveIncome: player.CalculatePassiveIncome(),
		SingleAsset:      business,
		Assets:           player.Assets.Business,
	})
}

func (c *playerTestController) Index(ctx *gin.Context) {
	c.cleaning(ctx)

	var content string

	for typeOfDeal, links := range Links {
		content += "<h2>" + typeOfDeal + "</h2>"

		content += "<ul>"

		keys := make([]string, 0, len(links))
		for key := range links {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		for _, name := range keys {
			content += `<li><a target="_blank" href="` + links[name] + `">` + name + `</a></li>`
		}

		content += "</ul>"
	}

	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(content))
}

func (c *playerTestController) cleaning(ctx *gin.Context) {
	cleaning := ctx.DefaultQuery("cleaning", "0")
	raceIDNum, _ := strconv.Atoi(ctx.DefaultQuery("raceID", "30"))
	raceID := uint64(raceIDNum)

	if cleaning == "1" {
		players := c.playerService.GetAllPlayersByRaceId(raceID)

		for _, player := range players {
			player.Assets.OtherAssets = []entity.CardOtherAssets{}
			player.Assets.Stocks = []entity.CardStocks{}
			player.Assets.Business = []entity.CardBusiness{}
			player.Assets.RealEstates = []entity.CardRealEstate{}

			err, _ := c.playerService.UpdatePlayer(&player)

			if err != nil {
				logger.Error(err)
			}
		}
	}
}

func (c *playerTestController) getPlayer(ctx *gin.Context) entity.Player {
	raceIDNum, _ := strconv.Atoi(ctx.DefaultQuery("raceID", "30"))
	raceID := uint64(raceIDNum)
	userIDNum, _ := strconv.Atoi(ctx.DefaultQuery("userID", "3"))
	userID := uint64(userIDNum)

	_, player := c.playerService.GetPlayerByUserIdAndRaceId(raceID, userID)

	return player
}

func (c *playerTestController) cardBusiness(ctx *gin.Context, player *entity.Player) entity.CardBusiness {
	if player == nil {
		player = new(entity.Player)
		*player = c.getPlayer(ctx)
	}

	card := c.getCard(ctx, entity.BusinessTypes.Startup).(entity.CardBusiness)

	c.addCashForPlayer(player, card.WholeCost, false)
	card.WholeCost = 0

	return card
}

func (c *playerTestController) getCard(ctx *gin.Context, defaultValue string) interface{} {
	typeOf := ctx.DefaultQuery("type", defaultValue)

	switch typeOf {
	case entity.BusinessTypes.Limited:
		return entity.CardBusiness{
			ID:          helper.Uuid("limitedPartnership"),
			Type:        "business",
			Symbol:      "limited",
			Heading:     "Магазин бутербродов",
			Description: "Some description",
			AssetType:   entity.BusinessTypes.Limited,
			Cost:        5000,
			Limit:       10,
			CashFlow:    180,
			WholeCost:   50000,
		}

	case entity.BusinessTypes.Startup:
		return entity.CardBusiness{
			ID:          helper.Uuid("business"),
			Type:        "business",
			AssetType:   entity.BusinessTypes.Startup,
			Symbol:      "startup",
			Heading:     "Вы создали IT компанию",
			Description: "Successful 4-bay, coin-operated auto wash near busy intersection. ",
			Cost:        5000,
			CashFlow:    300,
			WholeCost:   5000,
		}

	case entity.RealEstateTypes.Single:
		return entity.CardRealEstate{
			ID:          helper.Uuid("realEstate"),
			Type:        "realEstate",
			Symbol:      "single",
			AssetType:   entity.RealEstateTypes.Single,
			Heading:     "3-х комнатная квартира в Софии",
			Description: "Research and development delays cause low share price for this longtime pharmaceutical company",
			Cost:        50000,
			Mortgage:    48000,
			DownPayment: 2000,
			CashFlow:    300,
		}

	case entity.RealEstateTypes.Building:
		return entity.CardRealEstate{
			ID:          helper.Uuid("realEstate"),
			Type:        "realEstate",
			Symbol:      "building",
			Heading:     "24-квартирный жилой дом",
			AssetType:   entity.RealEstateTypes.Building,
			Description: "Research and development delays cause low share price for this longtime pharmaceutical company",
			Count:       24,
			Cost:        500000,
			Mortgage:    420000,
			DownPayment: 80000,
			CashFlow:    3000,
		}

	case entity.OtherAssetTypes.Piece:
		return entity.CardOtherAssets{
			ID:          helper.Uuid("land20akr"),
			Type:        "other",
			IsOwner:     true,
			AssetType:   entity.OtherAssetTypes.Whole,
			Cost:        5000,
			Count:       20,
			Symbol:      "whole",
			Heading:     "20 акров земли",
			Description: "Супер возможность на покупку 20 акров земли",
			WholeCost:   5000,
		}

	case entity.OtherAssetTypes.Whole:
		return entity.CardOtherAssets{
			ID:          helper.Uuid("goldCoins"),
			Type:        "other",
			IsOwner:     true,
			AssetType:   entity.OtherAssetTypes.Piece,
			CostPerOne:  1000,
			Count:       10,
			Symbol:      "piece",
			Heading:     "Золотые монеты 12го века",
			Description: "Супер возможность на покупку уникальных золотых монет",
			WholeCost:   10000,
		}

	case "inflation":
		return entity.CardMarket{
			ID:          helper.Uuid("inflation"),
			Type:        "inflation",
			Heading:     "Удары инфляции!",
			Symbol:      "ANY",
			Description: "Все 3Вг/2Ва дома которыми вы владеете забираются банком без права выкупа.",
			AssetType:   entity.MarketTypes.EachRealEstate,
			OnlyYou:     true,
		}

	case "success":
		return entity.CardMarket{
			ID:          helper.Uuid("success"),
			Type:        "success",
			Heading:     "Экономический рост",
			Symbol:      "ANY",
			Description: "Ваша созданная компания заключила важный договор и ваш пасс. дох. увеличился на $400",
			AssetType:   entity.MarketTypes.EachStartup,
			CashFlow:    400,
			OnlyYou:     true,
		}

	case "failure":
		return entity.CardMarket{
			ID:          helper.Uuid("failure"),
			Type:        "failure",
			Heading:     "Арендатор повредил вашу собственность",
			Symbol:      "ANY",
			Description: "Потеряв работу арендатор отказался платить и скрылся, заплатите $1000",
			AssetType:   entity.MarketTypes.AnyRealEstate,
			Cost:        1000,
			OnlyYou:     false,
		}

	default:
		return entity.Card{}
	}
}

func (c *playerTestController) addCashForPlayer(player *entity.Player, cash int, cleaning bool) {
	if cleaning {
		player.Assets.OtherAssets = []entity.CardOtherAssets{}
		player.Assets.Stocks = []entity.CardStocks{}
		player.Assets.Business = []entity.CardBusiness{}
		player.Assets.RealEstates = []entity.CardRealEstate{}
	}
	player.Cash = cash

	c.playerService.UpdateCash(player, 0, "Прибавка к зп")
}

func (c *playerTestController) fillAmounts(players []entity.Player, count int, min int, max int, key string) []dto.CardPurchasePlayerActionDTO {
	minAmount := min
	maxAmount := max

	amounts := make([]int, len(players))

	// Заполняем amounts случайными значениями
	for i := range players {
		if count == 0 {
			break
		}

		// Максимальное количество, которое может получить текущий игрок
		maxPossible := maxAmount
		if count < maxAmount {
			maxPossible = count
		}

		amount := helper.RandomMinMax(minAmount, maxPossible)
		amounts[i] = amount
		count -= amount
	}

	// Распределяем оставшийся count среди игроков
	for count > 0 {
		for i := range players {
			if count == 0 {
				break
			}
			if amounts[i] < maxAmount {
				amounts[i]++
				count--
			}
		}
	}

	var parts []dto.CardPurchasePlayerActionDTO

	for i, amount := range amounts {
		fmt.Printf("Player %d: %d\n", players[i], amount)

		if key == "amount" {
			parts = append(parts, dto.CardPurchasePlayerActionDTO{
				ID:     int(players[i].ID),
				Amount: amount,
			})
		} else if key == "passive" {
			parts = append(parts, dto.CardPurchasePlayerActionDTO{
				ID:      int(players[i].ID),
				Passive: amount,
			})
		} else if key == "percent" {
			parts = append(parts, dto.CardPurchasePlayerActionDTO{
				ID:      int(players[i].ID),
				Percent: amount,
			})
		}
	}

	return parts
}
