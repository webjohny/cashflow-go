package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/objects"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/storage"
)

type RaceService interface {
	BusinessAction(raceId uint64, userId uint64, isBigRace bool, dto dto.CardPurchaseActionDTO) error
	RealEstateAction(raceId uint64, userId uint64, isBigRace bool, dto dto.CardPurchaseActionDTO) error
	DreamAction(raceId uint64, userId uint64) error
	OtherAssetsAction(raceId uint64, userId uint64, dto dto.CardPurchaseActionDTO) error
	StocksAction(raceId uint64, userId uint64, count int) error
	LotteryAction(raceId uint64, userId uint64, isBigRace bool) (error, dto.RiskResponseDTO)
	MlmAction(raceId uint64, userId uint64, isBigRace bool) error
	SellBusinessAction(raceId uint64, userId uint64, assetId string, count int) error
	SellRealEstateAction(raceId uint64, userId uint64, realEstateId string) error
	SellStocksAction(raceId uint64, userId uint64, count int) error
	SellOtherAssetsAction(raceId uint64, userId uint64, assetId string, count int) error
	SkipAction(raceId uint64, userId uint64, isBigRace bool) error
	PaydayAction(raceId uint64, userId uint64, actionType string, isBigRace bool) error
	MarketAction(raceId uint64, userId uint64, actionType string) error
	CharityAction(raceId uint64, userId uint64, actionType string, isBigRace bool) error
	BabyAction(raceId uint64, userId uint64) (error, dto.MessageResponseDto)
	DoodadAction(raceId uint64, userId uint64) error
	DownsizedAction(raceId uint64, userId uint64) error
	BigBankruptAction(raceId uint64, userId uint64) error
	ChangeTurn(race entity.Race, forced bool, definedPlayerId int) error
	CreateResponses(raceId uint64, currentPlayerId uint64) []entity.RaceResponse
	GetRaceAndPlayer(raceId uint64, userId uint64) (error, entity.Race, entity.Player)
	GetRaceByRaceId(raceId uint64) entity.Race
	GetRacePlayersByRaceId(raceId uint64, all bool) []dto.GetRacePlayerResponseDTO
	GetFormattedRaceResponse(raceId uint64, hasExtraInfo bool) dto.GetRaceResponseDTO
	SetTransaction(player entity.Player, card dto.TransactionCardDTO) error
	InsertRace(b *entity.Race) (error, entity.Race)
	UpdateRace(b *entity.Race) (error, entity.Race)
}

type raceService struct {
	raceRepository     repository.RaceRepository
	playerService      PlayerService
	transactionService TransactionService
}

func NewRaceService(raceRepo repository.RaceRepository, playerService PlayerService, transactionService TransactionService) RaceService {
	return &raceService{
		raceRepository:     raceRepo,
		playerService:      playerService,
		transactionService: transactionService,
	}
}

func (service *raceService) InsertRace(b *entity.Race) (error, entity.Race) {
	logger.Info("RaceService.InsertRace")

	return service.raceRepository.InsertRace(b)
}

func (service *raceService) UpdateRace(b *entity.Race) (error, entity.Race) {
	logger.Info("RaceService.UpdateRace")

	return service.raceRepository.UpdateRace(b)
}

func (service *raceService) GetRaceAndPlayer(raceId uint64, userId uint64) (error, entity.Race, entity.Player) {
	race := service.GetRaceByRaceId(raceId)
	err, player := service.playerService.GetPlayerByUserIdAndRaceId(raceId, userId)

	if err != nil {
		return err, entity.Race{}, entity.Player{}
	} else if race.ID == 0 {
		return errors.New(storage.ErrorUndefinedGame), entity.Race{}, entity.Player{}
	}

	return nil, race, player
}

func (service *raceService) BusinessAction(raceId uint64, userId uint64, isBigRace bool, data dto.CardPurchaseActionDTO) error {
	logger.Info("RaceService.BusinessAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	if race.CurrentCard.Type != "business" {
		return errors.New(storage.ErrorInvalidTypeOfCard)
	}

	card := entity.CardBusiness{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		AssetType:   race.CurrentCard.AssetType,
		Symbol:      race.CurrentCard.Symbol,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
		Rule:        race.CurrentCard.Rule,
		History:     race.CurrentCard.History,
		Percent:     race.CurrentCard.Percent,
		Cost:        race.CurrentCard.Cost,
		Limit:       race.CurrentCard.Limit,
		ExtraDices:  race.CurrentCard.ExtraDices,
		Mortgage:    race.CurrentCard.Mortgage,
		DownPayment: race.CurrentCard.DownPayment,
		CashFlow:    race.CurrentCard.CashFlow,
	}

	if len(data.Players) > 0 {
		players := service.playerService.GetAllPlayersByRaceId(raceId)
		err = service.playerService.BuyBusinessInPartnership(card, player, players, data.Players)
	} else {
		card.IsOwner = true
		err = service.playerService.BuyBusiness(card, player, data.Count, true)
	}

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
	}

	return err
}

func (service *raceService) RealEstateAction(raceId uint64, userId uint64, isBigRace bool, data dto.CardPurchaseActionDTO) error {
	logger.Info("RaceService.RealEstateAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	if race.CurrentCard.Type != "realEstate" {
		return errors.New(storage.ErrorInvalidTypeOfCard)
	}

	card := entity.CardRealEstate{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		AssetType:   race.CurrentCard.AssetType,
		Symbol:      race.CurrentCard.Symbol,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
		Rule:        race.CurrentCard.Rule,
		Percent:     race.CurrentCard.Percent,
		Count:       race.CurrentCard.Count,
		Cost:        race.CurrentCard.Cost,
		CashFlow:    race.CurrentCard.CashFlow,
		Mortgage:    race.CurrentCard.Mortgage,
		DownPayment: race.CurrentCard.DownPayment,
	}

	if len(data.Players) > 0 {
		players := service.playerService.GetAllPlayersByRaceId(raceId)
		err = service.playerService.BuyRealEstateInPartnership(card, player, players, data.Players)
	} else {
		card.IsOwner = true
		err = service.playerService.BuyRealEstate(card, player)
	}

	//this.#log.addLog(player, `Купил бизнес ${this.#card.symbol} за $${this.#card.cost}`);
	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)

		go service.SetTransaction(player, dto.TransactionCardDTO{
			CardID:   race.CurrentCard.ID,
			CardType: entity.TransactionCardType.RealEstate,
			Details:  race.CurrentCard.Heading,
		})
	}

	return err
}

func (service *raceService) DreamAction(raceId uint64, userId uint64) error {
	logger.Info("RaceService.DreamAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	err = service.playerService.BuyDream(entity.CardDream{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
		Cost:        race.CurrentCard.Cost,
		AssetType:   race.CurrentCard.AssetType,
		PlayerId:    race.CurrentCard.PlayerId,
	}, player)

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
	}

	return err
}

func (service *raceService) LotteryAction(raceId uint64, userId uint64, isBigRace bool) (error, dto.RiskResponseDTO) {
	logger.Info("RaceService.LotteryAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err, dto.RiskResponseDTO{RolledDice: 0}
	}
	if race.CurrentCard.ID == "" {
		return errors.New(storage.ErrorInvalidCard), dto.RiskResponseDTO{}
	}

	dice := objects.NewDice(1, 2, 6)
	roll := dice.Roll(1)
	rolledDice := roll[0]

	var status bool

	outcome := race.CurrentCard.Outcome.(map[string]interface{})

	err, status = service.playerService.BuyLottery(entity.CardLottery{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		Symbol:      race.CurrentCard.Symbol,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
		Cost:        race.CurrentCard.Cost,
		AssetType:   race.CurrentCard.AssetType,
		Rule:        race.CurrentCard.Rule,
		SubRule:     race.CurrentCard.SubRule,
		Failure:     race.CurrentCard.Failure,
		Success:     race.CurrentCard.Success,
		Outcome: entity.CardLotteryOutcome{
			Failure: int(outcome["failure"].(float64)),
			Success: int(outcome["success"].(float64)),
		},
	}, player, rolledDice)

	var response = dto.RiskResponseDTO{RolledDice: rolledDice}

	if err != nil {
		return err, response
	}

	checkRiskDeal := helper.Contains([]string{entity.BigBusinessTypes.RiskBusiness, entity.BigBusinessTypes.RiskStocks}, race.CurrentCard.Type)

	if status && checkRiskDeal {
		response.Message = storage.MessageSuccessRiskDeal
	} else if status && race.CurrentCard.Type == entity.SmallDealTypes.Lottery {
		response.Message = storage.MessageSuccessLottery
	} else if !status && checkRiskDeal {
		response.Error = storage.MessageFailRiskDeal
	} else if !status && race.CurrentCard.Type == entity.SmallDealTypes.Lottery {
		response.Error = storage.MessageFailLottery
	} else {
		return errors.New(storage.ErrorPermissionDenied), response
	}

	race.Respond(player.ID, race.CurrentPlayer.ID)
	err, _ = service.UpdateRace(&race)

	return err, response
}

func (service *raceService) ChangeTurn(race entity.Race, forced bool, definedPlayerId int) error {
	logger.Info("RaceService.ChangeTurn", map[string]interface{}{
		"raceId": race.ID,
	})

	if !forced && race.CurrentCard.ID != "" && !race.IsReceived(race.CurrentPlayer.Username) {
		return nil
	}

	if definedPlayerId > 0 {
		race.PickCurrentPlayer(definedPlayerId)
	} else {
		race.NextPlayer()
	}

	race.ResetResponses()
	race.IsMultiFlow = false

	currentPlayer := race.CurrentPlayer

	err, player := service.playerService.GetPlayerByUserIdAndRaceId(race.ID, currentPlayer.UserId)

	if err != nil {
		return err
	}

	player.ChangeDiceStatus(false)

	err, _ = service.playerService.UpdatePlayer(&player)

	if err != nil {
		return err
	}

	race.CurrentCard = entity.Card{}

	err, race = service.UpdateRace(&race)

	if err != nil {
		return err
	}

	check := false

	if definedPlayerId == 0 && player.SkippedTurns > 0 {
		player.DecrementSkippedTurns()
		check = true
	} else if !forced && player.DualDiceCount > 0 {
		player.DecrementDualDiceCount()
		check = true
	}

	if check {
		err, _ = service.playerService.UpdatePlayer(&player)

		if err != nil {
			return err
		}
	}

	if player.SkippedTurns == 0 {
		return nil
	}

	err = service.ChangeTurn(race, true, 0)

	if err != nil {
		return err
	}

	return nil
}

func (service *raceService) CreateResponses(raceId uint64, currentPlayerId uint64) []entity.RaceResponse {
	logger.Info("RaceService.createResponses", map[string]interface{}{
		"raceId":          raceId,
		"currentPlayerId": currentPlayerId,
	})

	players := service.playerService.GetAllPlayersByRaceId(raceId)

	responses := make([]entity.RaceResponse, 0)

	if len(players) > 0 {
		for _, player := range players {
			response := player.CreateResponse()

			if player.OnBigRace && player.ID != currentPlayerId {
				response.Responded = true
			}
			responses = append(responses, response)
		}
	}

	return responses
}

func (service *raceService) SellRealEstateAction(raceId uint64, userId uint64, realEstateId string) error {
	logger.Info("RaceService.SellRealEstate", map[string]interface{}{
		"raceId":       raceId,
		"userId":       userId,
		"realEstateId": realEstateId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	err = service.playerService.SellRealEstate(realEstateId, entity.CardMarketRealEstate{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		Heading:     race.CurrentCard.Heading,
		Symbol:      race.CurrentCard.Symbol,
		Description: race.CurrentCard.Description,
		Rule:        race.CurrentCard.Rule,
		AssetType:   race.CurrentCard.AssetType,
		SubRule:     race.CurrentCard.SubRule,
		Cost:        race.CurrentCard.Cost,
		Range:       race.CurrentCard.Range,
	}, player)

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
	}

	return err
}

func (service *raceService) SellStocksAction(raceId uint64, userId uint64, count int) error {
	logger.Info("RaceService.SellStocksAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"count":  count,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	err = service.playerService.SellStocks(entity.CardStocks{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		Symbol:      race.CurrentCard.Symbol,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
		Rule:        race.CurrentCard.Rule,
		Price:       race.CurrentCard.Price,
		Count:       race.CurrentCard.Count,
		Increase:    race.CurrentCard.Increase,
		Decrease:    race.CurrentCard.Decrease,
		OnlyYou:     race.CurrentCard.OnlyYou,
		Range:       race.CurrentCard.Range,
	}, player, count, true)

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
	}

	return err
}

func (service *raceService) SellOtherAssetsAction(raceId uint64, userId uint64, assetId string, count int) error {
	logger.Info("RaceService.SellOtherAssetsAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"count":  count,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	err = service.playerService.SellOtherAssets(assetId, entity.CardMarketOtherAssets{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		Symbol:      race.CurrentCard.Symbol,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
		Cost:        race.CurrentCard.Cost,
		Rule:        race.CurrentCard.Rule,
		SubRule:     race.CurrentCard.SubRule,
		AssetType:   race.CurrentCard.AssetType,
	}, player, count)

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
	}

	return err
}

func (service *raceService) SellBusinessAction(raceId uint64, userId uint64, assetId string, count int) error {
	logger.Info("RaceService.SellBusinessAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"count":  count,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	err, _ = service.playerService.SellBusiness(assetId, entity.CardMarketBusiness{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		Symbol:      race.CurrentCard.Symbol,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
		Cost:        race.CurrentCard.Cost,
		CashFlow:    race.CurrentCard.CashFlow,
		Rule:        race.CurrentCard.Rule,
		SubRule:     race.CurrentCard.SubRule,
		AssetType:   race.CurrentCard.AssetType,
	}, player, count)

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
	}

	return err
}

func (service *raceService) StocksAction(raceId uint64, userId uint64, count int) error {
	logger.Info("RaceService.StocksAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"count":  count,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	if race.CurrentCard.Type != "stock" {
		return errors.New(storage.ErrorInvalidTypeOfCard)
	}

	cardStocks := entity.CardStocks{}
	race.CurrentCard.Count = count
	cardStocks.Fill(race.CurrentCard)

	if count > 0 {
		err = service.playerService.BuyStocks(cardStocks, player, true)
	}

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
	}

	return err
}

func (service *raceService) OtherAssetsAction(raceId uint64, userId uint64, data dto.CardPurchaseActionDTO) error {
	logger.Info("RaceService.OtherAssetsAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"dto":    data,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	card := entity.CardOtherAssets{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		Cost:        race.CurrentCard.Cost,
		CostPerOne:  race.CurrentCard.CostPerOne,
		Count:       race.CurrentCard.Count,
		AssetType:   race.CurrentCard.AssetType,
		Symbol:      race.CurrentCard.Symbol,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
	}

	if len(data.Players) > 0 {
		players := service.playerService.GetAllPlayersByRaceId(raceId)
		err = service.playerService.BuyOtherAssetsInPartnership(card, player, players, data.Players)
	} else {
		card.IsOwner = true
		card.WholeCost = card.Cost

		if card.AssetType == entity.OtherAssetTypes.Piece {
			card.WholeCost = data.Count * card.CostPerOne
		}

		err = service.playerService.BuyOtherAssets(card, player, data.Count)
	}

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
	}

	return err
}

func (service *raceService) MlmAction(raceId uint64, userId uint64, isBigRace bool) error {
	logger.Info("RaceService.MlmAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
	}

	return err
}

func (service *raceService) SkipAction(raceId uint64, userId uint64, isBigRace bool) error {
	logger.Info("RaceService.SkipAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
	}

	return err
}

func (service *raceService) CharityAction(raceId uint64, userId uint64, actionType string, isBigRace bool) error {
	logger.Info("RaceService.CharityAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	card := entity.CardCharity{}
	card.Fill(race.CurrentCard)

	err = service.playerService.Charity(card, player)

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err = service.ChangeTurn(race, false, 0)
	}

	return err
}

func (service *raceService) BabyAction(raceId uint64, userId uint64) (error, dto.MessageResponseDto) {
	logger.Info("RaceService.BabyAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err, dto.MessageResponseDto{}
	}

	err, response := service.playerService.BornBaby(player, race.CurrentCard)

	if err != nil && !response {
		return err, dto.MessageResponseDto{}
	}

	race.Respond(player.ID, race.CurrentPlayer.ID)
	err = service.ChangeTurn(race, false, 0)

	if response && err != nil {
		return nil, dto.MessageResponseDto{Message: err.Error(), Result: false}
	}

	return err, dto.MessageResponseDto{Message: storage.MessageYouHadBaby, Result: true}
}

func (service *raceService) DoodadAction(raceId uint64, userId uint64) error {
	logger.Info("RaceService.DoodadAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	race.Respond(player.ID, race.CurrentPlayer.ID)

	if (player.Babies == 0 && race.CurrentCard.HasBabies) || player.HasHealthyInsurance() {
		err = service.ChangeTurn(race, false, 0)

		if err != nil {
			return err
		}

		if player.HasHealthyInsurance() {
			return errors.New(storage.WarnYouHaveHealthyInsurance)
		}

		return errors.New(storage.WarnYouHaveNoBabies)
	}

	err = service.playerService.Doodad(entity.CardDoodad{
		ID:            race.CurrentCard.ID,
		Type:          race.CurrentCard.Type,
		Heading:       race.CurrentCard.Heading,
		Symbol:        race.CurrentCard.Symbol,
		Description:   race.CurrentCard.Description,
		Cost:          race.CurrentCard.Cost,
		Rule:          race.CurrentCard.Rule,
		IsConditional: race.CurrentCard.IsConditional,
		AssetType:     race.CurrentCard.AssetType,
		HasBabies:     race.CurrentCard.HasBabies,
	}, player)

	if err == nil {
		return service.ChangeTurn(race, false, 0)
	}

	return err
}

func (service *raceService) DownsizedAction(raceId uint64, userId uint64) error {
	logger.Info("RaceService.DownsizedAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	if err = service.playerService.Downsized(player, race.CurrentCard); err != nil {
		return err
	}

	race.Respond(player.ID, race.CurrentPlayer.ID)

	if err, _ = service.UpdateRace(&race); err != nil {
		return err
	}

	return nil
}

func (service *raceService) BigBankruptAction(raceId uint64, userId uint64) error {
	logger.Info("RaceService.BigBankruptAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	if err = service.playerService.BigBankrupt(player); err != nil {
		return err
	}

	race.Respond(player.ID, race.CurrentPlayer.ID)

	if err, _ = service.UpdateRace(&race); err != nil {
		return err
	}

	return nil
}

func (service *raceService) MarketAction(raceId uint64, userId uint64, actionType string) error {
	logger.Info("RaceService.MarketAction", map[string]interface{}{
		"raceId":     raceId,
		"userId":     userId,
		"actionType": actionType,
	})

	// Retrieve race and player
	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	// Fill the card market
	cardMarket := entity.CardMarket{}
	cardMarket.Fill(race.CurrentCard)

	transactionCard := dto.TransactionCardDTO{
		CardID:   race.CurrentCard.ID,
		CardType: "",
		Details:  race.CurrentCard.Heading,
	}

	if cardMarket.Type == entity.TransactionCardType.Damage {
		if !player.HasOwnRealEstates() {
			return errors.New(storage.ErrorNotFoundTheRealEstate)
		}

		transactionCard.CardType = entity.TransactionCardType.Damage
	} else if cardMarket.Type == entity.TransactionCardType.Business {
		if !player.HasOwnRealEstates() {
			return errors.New(storage.ErrorNotFoundTheBusiness)
		}

		transactionCard.CardType = entity.TransactionCardType.MarketBusiness
	}

	trx := service.transactionService.GetRaceTransaction(player, transactionCard)

	if trx.ID > 0 {
		return errors.New(storage.ErrorTransactionAlreadyExists)
	}

	var actionErr error

	// Determine the action type and execute logic
	switch cardMarket.Type {
	case entity.TransactionCardType.Damage:
		actionErr = service.playerService.MarketDamage(cardMarket, player)
	case entity.TransactionCardType.Business:
		cardMarketBusiness := entity.CardMarketBusiness{}
		cardMarketBusiness.Fill(race.CurrentCard)

		actionErr = service.playerService.MarketBusiness(cardMarketBusiness, player)
	}

	if actionErr != nil {
		return actionErr
	}

	race.Respond(player.ID, race.CurrentPlayer.ID)
	if err, _ = service.UpdateRace(&race); err != nil {
		return err
	}

	return nil
}

func (service *raceService) PaydayAction(raceId uint64, userId uint64, actionType string, isBigRace bool) error {
	logger.Info("RaceService.MarketAction", map[string]interface{}{
		"raceId":     raceId,
		"userId":     userId,
		"actionType": actionType,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	if actionType == "payday" {
		return service.playerService.Payday(player, race.CurrentCard)
	} else if actionType == "cashFlowDay" {
		return service.playerService.CashFlowDay(player, race.CurrentCard)
	}

	return service.ChangeTurn(race, false, 0)
}

func (service *raceService) GetRaceByRaceId(raceId uint64) entity.Race {
	return service.raceRepository.FindRaceById(raceId)
}

func (service *raceService) GetRacePlayersByRaceId(raceId uint64, all bool) []dto.GetRacePlayerResponseDTO {
	players := make([]entity.Player, 0)

	if all {
		players = service.playerService.GetAllStatePlayersByRaceId(raceId)
	} else {
		players = service.playerService.GetAllPlayersByRaceId(raceId)
	}

	racePlayers := make([]dto.GetRacePlayerResponseDTO, 0)

	for _, player := range players {
		racePlayer := service.playerService.GetFormattedPlayerResponse(player, all)
		racePlayers = append(racePlayers, racePlayer)
	}

	return racePlayers
}

func (service *raceService) GetFormattedRaceResponse(raceId uint64, hasExtraInfo bool) dto.GetRaceResponseDTO {
	race := service.GetRaceByRaceId(raceId)

	var logs []entity.RaceLog

	if hasExtraInfo {
		logs = service.transactionService.GetRaceLogs(raceId)
	}

	players := service.GetRacePlayersByRaceId(raceId, hasExtraInfo)
	err, player := service.playerService.GetRacePlayer(raceId, race.CurrentPlayer.UserId, hasExtraInfo)

	if err != nil {
		logger.Error(err, map[string]interface{}{
			"raceId": raceId,
			"userId": race.CurrentPlayer.UserId,
		})

		return dto.GetRaceResponseDTO{}
	}

	response := dto.GetRaceResponseDTO{
		Players:       players,
		TurnResponses: race.Responses,
		Status:        race.Status,
		DiceValues:    race.Dice,
		CurrentPlayer: player,
		CurrentCard:   race.CurrentCard,
		Options:       race.Options,
		GameId:        race.ID,
		IsTurnEnded:   race.IsReceived(player.Username),
		IsMultiFlow:   race.IsMultiFlow,
		Logs:          logs,
	}

	return response
}

func (service *raceService) SetTransaction(player entity.Player, card dto.TransactionCardDTO) error {
	return service.transactionService.InsertRaceTransaction(dto.TransactionCreateRaceDTO{
		RaceID:   player.RaceID,
		CardID:   card.CardID,
		Details:  card.Details,
		PlayerID: player.ID,
		Username: player.Username,
		Color:    player.Color,
		CardType: card.CardType,
	})
}
