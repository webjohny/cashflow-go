package controller

import (
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
	BuyStocks(ctx *gin.Context)
	SellStocks(ctx *gin.Context)
	DecreaseStocks(ctx *gin.Context)
	IncreaseStocks(ctx *gin.Context)
	BuyRealEstate(ctx *gin.Context)
	BuyBusiness(ctx *gin.Context)
	BuyBigBusiness(ctx *gin.Context)
	BuyLottery(ctx *gin.Context)
	BuyOtherAssets(ctx *gin.Context)
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
		"Недвижимость (в партнёрстве)": "/test/player/buy-partner-real-estate",
		"Бизнес": "/test/player/buy-business",
		"Бизнес (в партнёрстве)": "/test/player/buy-partner-business",
		"Лотерея":                           "/test/player/buy-lottery",
		"Другие активы":                     "/test/player/buy-other-assets",
		"Большой круг - Бизнес":             "/test/player/buy-big-business",
		"Большой круг - Рискованный бизнес": "/test/player/buy-risk-business",
		"Большой круг - Рискованные акции":  "/test/player/buy-risk-stocks",
		"Большой круг - Мечта":              "/test/player/buy-dream",
	},
	"Продажа": {
		"Акции": "/test/player/sell-stocks",
	},
	"Манипуляции": {
		"Удвоение кол-во акций":          "/test/player/increase-stocks",
		"Разделение (на 3) кол-во акций": "/test/player/decrease-stocks",
	},
}

func (c *playerTestController) BuyBusiness(ctx *gin.Context) {
	b := c.getPlayer(ctx)

	card := entity.CardBusiness{
		ID:          helper.Uuid("business"),
		Type:        "business",
		Symbol:      "4-BAY,AUTO-WASH",
		Heading:     "Automated Business",
		Description: "Successful 4-bay, coin-operated auto wash near busy intersection. ",
		Cost:        125000,
		Mortgage:    100000,
		DownPayment: 25000,
		CashFlow:    1800,
	}
	c.addCashForPlayer(&b, card.DownPayment, false)

	//card := entity.CardBusiness{
	//	ID:          "limitedPartnershipS3",
	//	Type:        "business",
	//	Symbol:      "SANDWICHES_SHOP",
	//	Heading:     "Магазин бутербродов",
	//	Description: "Some description",
	//	Cost:        5000,
	//	Limit:       10,
	//	CashFlow:    180,
	//}
	//c.addCashForPlayer(&b, card.Cost*card.Limit, true)

	err := c.playerService.BuyBusiness(card, b, 5, true)

	if err != nil {
		logger.Error(err, nil)

		request.FinalResponse(ctx, err, nil)
		return
	}

	player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

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
	b := c.getPlayer(ctx)

	card := entity.CardLottery{
		ID:          helper.Uuid("lottery"),
		Type:        "lottery",
		Symbol:      "lottery",
		Heading:     "Sister-In-Law borrows Money",
		Description: "Sister-in-law is downsized.Needs $5,000 to make house payment.",
		Cost:        1000,
		Lottery:     "money",
		Rule:        "If you choose to help. Pay $5,000 and roll 1 die:",
		SubRule: []string{
			"Die = 1-3, She never pays you back and you're out $5,000",
			"Die = 4-6, She pays you back $10,000, but family get-togethers are still awkward",
		},
		Failure: []int{1, 2, 3},
		Success: []int{4, 5, 6},
		Outcome: struct {
			Failure int `json:"failure"`
			Success int `json:"success"`
		}(struct {
			Failure int
			Success int
		}{
			Failure: 0,
			Success: 5000,
		}),
	}
	c.addCashForPlayer(&b, card.Cost, false)

	err, result := c.playerService.BuyLottery(card, b, helper.Random(6))

	player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

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
	b := c.getPlayer(ctx)

	card := entity.CardOtherAssets{
		ID:          helper.Uuid("goldCoins"),
		Type:        "other",
		CostPerOne:  1000,
		Count:       10,
		Symbol:      "goldCoins",
		Heading:     "Золотые монеты 12го века",
		Description: "Супер возможность на покупку уникальных золотых монет",
	}
	c.addCashForPlayer(&b, card.CostPerOne*card.Count, false)

	err := c.playerService.BuyOtherAssets(card, b, 10)

	player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

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

func (c *playerTestController) BuyRealEstateInPartnership(ctx *gin.Context) {
	raceIDNum, _ := strconv.Atoi(ctx.DefaultQuery("raceID", "30"))
	raceID := uint64(raceIDNum)

	players := c.playerService.GetAllPlayersByRaceId(raceID)

	card := entity.CardRealEstate{
		ID:          helper.Uuid("realEstateS1"),
		Type:        "realEstate",
		Symbol:      "2Ком/Кв",
		Heading:     "Супер-сделка",
		Description: "Отличная 2-х комнатная квартира в центре Киева.",
		Cost:        45000,
		Mortgage:    43000,
		DownPayment: 2000,
		CashFlow:    300,
	}

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

	i := 0
	for key := range players {
		if i > 1 {
			break
		}

		c.addCashForPlayer(&players[key], card.Cost, false)

		i++
	}

	err := c.playerService.BuyRealEstateInPartnership(card, players[0], players, parts)

	if err != nil {
		logger.Error(err, nil)

		request.FinalResponse(ctx, err, nil)
		return
	}

	newPlayers := c.playerService.GetAllPlayersByRaceId(raceID)

	responsePlayers := make([]PlayerResponse, 0)

	i = 0
	for key, player := range newPlayers {
		if i > 1 {
			break
		}

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

		i++
	}

	request.FinalResponse(ctx, err, responsePlayers)
}

func (c *playerTestController) BuyBusinessInPartnership(ctx *gin.Context) {
	raceIDNum, _ := strconv.Atoi(ctx.DefaultQuery("raceID", "30"))
	raceID := uint64(raceIDNum)

	players := c.playerService.GetAllPlayersByRaceId(raceID)

	card := entity.CardBusiness{
		ID:          helper.Uuid("business"),
		Type:        "business",
		Symbol:      "IT_COMPANY",
		Heading:     "You're create IT company",
		Description: "Some description",
		Cost:        5000,
		CashFlow:    300,
	}
	c.addCashForPlayer(&players[0], card.Cost, false)

	//card := entity.CardBusiness{
	//	ID:          "limitedPartnershipS3",
	//	Type:        "business",
	//	Symbol:      "SANDWICHES_SHOP",
	//	Heading:     "Магазин бутербродов",
	//	Description: "Some description",
	//	Cost:        5000,
	//	Limit:       10,
	//	CashFlow:    180,
	//}

	parts := []dto.CardPurchasePlayerActionDTO{
		{
			ID:      int(players[0].ID),
			Amount:  5,
			Passive: 120,
		},
		{
			ID:      int(players[1].ID),
			Amount:  3,
			Passive: 180,
		},
	}

	i := 0
	for key := range players {
		if i > 1 {
			break
		}

		limit := card.Limit

		if limit == 0 {
			limit = 1
		}

		c.addCashForPlayer(&players[key], card.Cost*limit, false)

		i++
	}

	err := c.playerService.BuyBusinessInPartnership(card, players[0], players, parts)

	if err != nil {
		logger.Error(err, nil)

		request.FinalResponse(ctx, err, nil)
		return
	}

	newPlayers := c.playerService.GetAllPlayersByRaceId(raceID)

	responsePlayers := make([]PlayerResponse, 0)

	i = 0
	for key, player := range newPlayers {
		if i > 1 {
			break
		}

		_, business := player.FindBusinessBySymbol(card.Symbol)

		responsePlayers = append(responsePlayers, PlayerResponse{
			ID:          player.ID,
			UserID:      player.UserID,
			Card:        card,
			OldCash:     players[key].Cash,
			NewCash:     player.Cash,
			OldCashFlow: players[key].CalculatePassiveIncome(),
			NewCashFlow: player.CalculatePassiveIncome(),
			SingleAsset: business,
			Assets:      player.Assets.Business,
		})

		i++
	}

	request.FinalResponse(ctx, err, responsePlayers)
}

func (c *playerTestController) BuyStocks(ctx *gin.Context) {
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

	player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

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

	player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

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

func (c *playerTestController) IncreaseStocks(ctx *gin.Context) {
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

	player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

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

	player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

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
	b := c.getPlayer(ctx)

	card := entity.CardRealEstate{
		ID:          helper.Uuid("realEstate"),
		Type:        "realEstate",
		Symbol:      "3КВ",
		Heading:     "3-х комнатная квартира в Софии",
		Description: "Research and development delays cause low share price for this longtime pharmaceutical company",
		Cost:        50000,
		Mortgage:    47000,
		DownPayment: 3000,
		CashFlow:    200,
	}

	c.addCashForPlayer(&b, card.DownPayment*5, false)

	err := c.playerService.BuyRealEstate(card, b)

	player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

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
	b := c.getPlayer(ctx)

	card := entity.CardDream{
		ID:          helper.Uuid("dreamParkAttraction1"),
		Heading:     "Парк аттракционов в Лондоне",
		Description: "Вы купили супер парк развлечений для своих детей",
		Type:        "dream",
		Cost:        100000,
	}

	c.addCashForPlayer(&b, card.Cost, true)

	err := c.playerService.BuyDream(card, b)

	player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

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
	b := c.getPlayer(ctx)

	card := entity.CardRiskStocks{
		ID:          helper.Uuid("riskStocks"),
		Type:        "",
		Count:       0,
		Cost:        0,
		Dices:       nil,
		ExtraDices:  0,
		Symbol:      "",
		Heading:     "",
		Description: "",
		CostPerOne:  0,
	}

	err, _ := c.playerService.BuyRiskStocks(card, b, 1)

	player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

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

func (c *playerTestController) BuyRiskBusiness(ctx *gin.Context) {
	b := c.getPlayer(ctx)

	card := entity.CardRiskBusiness{
		ID:          helper.Uuid("riskBusiness"),
		Type:        "",
		Dices:       nil,
		ExtraDices:  0,
		Symbol:      "",
		Heading:     "",
		Description: "",
		Cost:        0,
	}

	err, _ := c.playerService.BuyRiskBusiness(card, b, 1)

	player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

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

func (c *playerTestController) BuyBigBusiness(ctx *gin.Context) {
	b := c.getPlayer(ctx)

	card := entity.CardBusiness{
		ID:          helper.Uuid("business"),
		Type:        "business",
		Symbol:      "4-BAY,AUTO-WASH",
		Heading:     "Automated Business",
		Description: "Successful 4-bay, coin-operated auto wash near busy intersection. ",
		Cost:        125000,
		Mortgage:    100000,
		DownPayment: 25000,
		CashFlow:    1800,
	}
	c.addCashForPlayer(&b, card.DownPayment, false)

	//card := entity.CardBusiness{
	//	ID:          "limitedPartnershipS3",
	//	Type:        "business",
	//	Symbol:      "SANDWICHES_SHOP",
	//	Heading:     "Магазин бутербродов",
	//	Description: "Some description",
	//	Cost:        5000,
	//	Limit:       10,
	//	CashFlow:    180,
	//}
	//c.addCashForPlayer(&b, card.Cost*card.Limit, true)

	err := c.playerService.BuyBusiness(card, b, 5, true)

	if err != nil {
		logger.Error(err, nil)

		request.FinalResponse(ctx, err, nil)
		return
	}

	player := c.playerService.GetPlayerByUserIdAndRaceId(b.RaceID, b.UserID)

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

func (c *playerTestController) getPlayer(ctx *gin.Context) entity.Player {
	raceIDNum, _ := strconv.Atoi(ctx.DefaultQuery("raceID", "30"))
	raceID := uint64(raceIDNum)
	userIDNum, _ := strconv.Atoi(ctx.DefaultQuery("userID", "3"))
	userID := uint64(userIDNum)

	return c.playerService.GetPlayerByUserIdAndRaceId(raceID, userID)
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
