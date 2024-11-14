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
	ChangeTurn(race entity.Race, definedPlayerId int) error
	GetRaceAndPlayer(raceId uint64, userId uint64) (error, entity.Race, entity.Player)
	GetRaceByRaceId(raceId uint64) entity.Race
	GetRacePlayersByRaceId(raceId uint64) []dto.GetRacePlayerResponseDTO
	GetFormattedRaceResponse(raceId uint64) dto.GetRaceResponseDTO
	SetTransaction(ID uint64, player entity.Player, txType string, details string)
	RemovePlayer(raceId uint64, userId uint64) error
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
	logger.Info("RaceService.InsertRace", b)

	return service.raceRepository.InsertRace(b)
}

func (service *raceService) UpdateRace(b *entity.Race) (error, entity.Race) {
	logger.Info("RaceService.UpdateRace", b)

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

func (service *raceService) BusinessAction(raceId uint64, userId uint64, isBigRace bool, dto dto.CardPurchaseActionDTO) error {
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

	if len(dto.Players) > 0 {
		players := service.playerService.GetAllPlayersByRaceId(raceId)
		err = service.playerService.BuyBusinessInPartnership(card, player, players, dto.Players)
	} else {
		card.IsOwner = true
		err = service.playerService.BuyBusiness(card, player, dto.Count, true)
	}

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
		go service.SetTransaction(race.ID, player, entity.TxTypes.Business, storage.MessageYouBoughtBusiness)
	}

	return err
}

func (service *raceService) RealEstateAction(raceId uint64, userId uint64, isBigRace bool, dto dto.CardPurchaseActionDTO) error {
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

	if len(dto.Players) > 0 {
		players := service.playerService.GetAllPlayersByRaceId(raceId)
		err = service.playerService.BuyRealEstateInPartnership(card, player, players, dto.Players)
	} else {
		card.IsOwner = true
		err = service.playerService.BuyRealEstate(card, player)
	}

	//this.#log.addLog(player, `Купил бизнес ${this.#card.symbol} за $${this.#card.cost}`);
	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
		go service.SetTransaction(race.ID, player, entity.TxTypes.RealEstate, storage.MessageYouBoughtRealEstate)
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
	}, player)

	//this.#log.addLog(player, `Купил бизнес ${this.#card.symbol} за $${this.#card.cost}`);

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
		go service.SetTransaction(race.ID, player, entity.TxTypes.Dream, storage.MessageYouBoughtDream)
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
	rolledDice := roll[0] | 1

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
		go service.SetTransaction(race.ID, player, race.CurrentCard.Type, storage.MessageSuccessRiskDeal)
		response.Message = storage.MessageSuccessRiskDeal
	} else if status && race.CurrentCard.Type == entity.SmallDealTypes.Lottery {
		go service.SetTransaction(race.ID, player, entity.SmallDealTypes.Lottery, storage.MessageSuccessRiskDeal)
		response.Message = storage.MessageSuccessLottery
	} else if !status && checkRiskDeal {
		go service.SetTransaction(race.ID, player, race.CurrentCard.Type, storage.MessageFailRiskDeal)
		response.Error = storage.MessageFailRiskDeal
	} else if !status && race.CurrentCard.Type == entity.SmallDealTypes.Lottery {
		go service.SetTransaction(race.ID, player, entity.SmallDealTypes.Lottery, storage.MessageFailLottery)
		response.Error = storage.MessageFailLottery
	} else {
		return errors.New(storage.ErrorPermissionDenied), response
	}

	race.Respond(player.ID, race.CurrentPlayer.ID)
	err, _ = service.UpdateRace(&race)

	return err, response
}

func (service *raceService) ChangeTurn(race entity.Race, definedPlayerId int) error {
	logger.Info("RaceService.ChangeTurn", map[string]interface{}{
		"raceId": race.ID,
	})

	if definedPlayerId > 0 {
		race.PickCurrentPlayer(definedPlayerId)
	} else if race.CurrentCard.ID != "" && !race.IsReceived(race.CurrentPlayer.Username) {
		return nil
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
	race.Notifications = make([]entity.RaceNotification, 0)

	err, race = service.UpdateRace(&race)

	if err != nil {
		return err
	}

	if player.SkippedTurns > 0 {
		player.DecrementSkippedTurns()

		if player.DualDiceCount > 0 {
			player.DecrementDualDiceCount()
		}

		err, _ = service.playerService.UpdatePlayer(&player)

		if err != nil {
			return err
		}

		err = service.ChangeTurn(race, 0)

		if err != nil {
			return err
		}
	}

	return nil
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
		go service.SetTransaction(race.ID, player, entity.TxTypes.RealEstate, storage.MessageRealEstateHasBeenSold)
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
		go service.SetTransaction(race.ID, player, entity.TxTypes.Stocks, storage.MessageStocksHaveBeenSold)
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
		go service.SetTransaction(race.ID, player, entity.TxTypes.Other, storage.MessageGoldHasBeenSold)
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
		go service.SetTransaction(race.ID, player, entity.TxTypes.Other, storage.MessageGoldHasBeenSold)
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

	//this.#log.addLog(player, `Купил бизнес ${this.#card.symbol} за $${this.#card.cost}`);
	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
		go service.SetTransaction(race.ID, player, entity.TxTypes.Stocks, storage.MessageYouBoughtStocks)
	}

	return err
}

func (service *raceService) OtherAssetsAction(raceId uint64, userId uint64, dto dto.CardPurchaseActionDTO) error {
	logger.Info("RaceService.OtherAssetsAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"dto":    dto,
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

	if len(dto.Players) > 0 {
		players := service.playerService.GetAllPlayersByRaceId(raceId)
		err = service.playerService.BuyOtherAssetsInPartnership(card, player, players, dto.Players)
	} else {
		card.IsOwner = true
		card.WholeCost = card.Cost

		if card.AssetType == entity.OtherAssetTypes.Piece {
			card.WholeCost = dto.Count * card.CostPerOne
		}

		err = service.playerService.BuyOtherAssets(card, player, dto.Count)
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
		err = service.ChangeTurn(race, 0)
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

	if player.Babies <= 2 {
		player.BornBaby()
		err, _ = service.playerService.UpdatePlayer(&player)

		if err != nil {
			return err, dto.MessageResponseDto{}
		}
	}

	race.Respond(player.ID, race.CurrentPlayer.ID)
	err, race = service.UpdateRace(&race)

	if err != nil {
		return err, dto.MessageResponseDto{}
	}

	if player.Babies > 2 {
		return nil, dto.MessageResponseDto{Message: storage.MessageYouHaveTooManyBabies, Result: false}
	}

	return nil, dto.MessageResponseDto{Message: storage.MessageYouHadBaby, Result: true}
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

	err = service.playerService.Doodad(entity.CardDoodad{
		ID:            race.CurrentCard.ID,
		Type:          race.CurrentCard.Type,
		Heading:       race.CurrentCard.Heading,
		Symbol:        race.CurrentCard.Symbol,
		Description:   race.CurrentCard.Description,
		Cost:          race.CurrentCard.Cost,
		Rule:          race.CurrentCard.Rule,
		IsConditional: race.CurrentCard.IsConditional,
		HasBabies:     race.CurrentCard.HasBabies,
	}, player)

	if err == nil || err.Error() == storage.ErrorYouHaveNoBabies {
		go service.SetTransaction(race.ID, player, entity.TxTypes.Other, race.CurrentCard.Heading)

		race.Respond(player.ID, race.CurrentPlayer.ID)
		err = service.ChangeTurn(race, 0)
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

	err = service.playerService.Downsized(player)

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
	}

	return err
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

	err = service.playerService.BigBankrupt(player)

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
	}

	return err
}

func (service *raceService) MarketAction(raceId uint64, userId uint64, actionType string) error {
	logger.Info("RaceService.MarketAction", map[string]interface{}{
		"raceId":     raceId,
		"userId":     userId,
		"actionType": actionType,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	cardMarket := entity.CardMarket{}
	cardMarket.Fill(race.CurrentCard)

	if actionType == "damage" {
		err = service.playerService.MarketDamage(cardMarket, player)
	} else if actionType == "business" {
		cardMarketBusiness := entity.CardMarketBusiness{}
		cardMarketBusiness.Fill(race.CurrentCard)

		err = service.playerService.MarketBusiness(cardMarketBusiness, player)
	} else {
		return errors.New(storage.ErrorUndefinedTypeOfDeal)
	}

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, race = service.UpdateRace(&race)

		err, _ = service.UpdateRace(&race)
	}

	return err
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
		service.playerService.Payday(player)
	} else if actionType == "cashFlowDay" {
		service.playerService.CashFlowDay(player)
	}

	return service.ChangeTurn(race, 0)
}

func (service *raceService) RemovePlayer(raceId uint64, userId uint64) error {
	logger.Info("RaceService.RemovePlayer", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId)

	if err != nil {
		return err
	}

	race.RemoveResponsePlayer(player.ID)

	if race.CurrentPlayer.ID == player.ID {
		race.PickCurrentPlayer(int(race.Responses[0].ID))
	}

	err, _ = service.UpdateRace(&race)

	return err
}

func (service *raceService) GetRaceByRaceId(raceId uint64) entity.Race {
	return service.raceRepository.FindRaceById(raceId)
}

func (service *raceService) GetRacePlayersByRaceId(raceId uint64) []dto.GetRacePlayerResponseDTO {
	players := service.playerService.GetAllPlayersByRaceId(raceId)

	racePlayers := make([]dto.GetRacePlayerResponseDTO, 0)

	for _, player := range players {
		racePlayer := service.playerService.GetFormattedPlayerResponse(player, false)
		racePlayers = append(racePlayers, racePlayer)
	}

	return racePlayers
}

func (service *raceService) GetFormattedRaceResponse(raceId uint64) dto.GetRaceResponseDTO {
	race := service.GetRaceByRaceId(raceId)
	logs := service.transactionService.GetRaceLogs(raceId)
	players := service.GetRacePlayersByRaceId(raceId)
	err, player := service.playerService.GetRacePlayer(raceId, race.CurrentPlayer.UserId)

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
		GameId:        race.ID,
		IsTurnEnded:   race.IsReceived(player.Username),
		IsMultiFlow:   race.IsMultiFlow,
		Logs:          logs,
		Notifications: race.Notifications,
		Transaction:   entity.TransactionData{},
	}

	return response
}

func (service *raceService) SetTransaction(ID uint64, player entity.Player, txType string, details string) {
	service.transactionService.InsertRaceTransaction(dto.TransactionCreateRaceDTO{
		RaceID:   ID,
		Details:  details,
		PlayerID: player.ID,
		Username: player.Username,
		Color:    player.Color,
		TxType:   txType,
	})
}
