package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/logger"
	"github.com/webjohny/cashflow-go/objects"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/storage"
	"net/http"
	"strconv"
)

type RaceService interface {
	BusinessAction(raceId uint64, userId uint64, isBigRace bool) error
	RealEstateAction(raceId uint64, userId uint64, isBigRace bool) error
	DreamAction(raceId uint64, userId uint64) error
	RiskBusinessAction(raceId uint64, userId uint64) (error, dto.RiskResponseDTO)
	RiskStocksAction(raceId uint64, userId uint64) (error, dto.RiskResponseDTO)
	StocksAction(raceId uint64, userId uint64, count int) error
	LotteryAction(raceId uint64, userId uint64, isBigRace bool) error
	GoldCoinsAction(raceId uint64, userId uint64) error
	MlmAction(raceId uint64, userId uint64, isBigRace bool) error
	SellRealEstate(raceId uint64, userId uint64, realEstateId string) error
	SellStocks(raceId uint64, userId uint64, count int) error
	SellGoldCoins(raceId uint64, userId uint64, count int) error
	SkipAction(raceId uint64, userId uint64, isBigRace bool) error
	PaydayAction(raceId uint64, userId uint64, actionType string, isBigRace bool) error
	MarketAction(raceId uint64, userId uint64, actionType string) error
	CharityAction(raceId uint64, userId uint64, isBigRace bool) error
	BabyAction(raceId uint64, userId uint64) error
	DoodadAction(raceId uint64, userId uint64) error
	GetRaceAndPlayer(raceId uint64, userId uint64, isBigRace bool) (error, entity.Race, entity.Player)
	GetInjectedRace(ctx *gin.Context) entity.Race
	GetRaceByRaceId(raceId uint64, isBigRace bool) entity.Race
	GetRacePlayersByRaceId(raceId uint64) []dto.GetRacePlayerResponseDTO
	GetFormattedRaceResponse(raceId uint64, isBigRace bool) dto.GetRaceResponseDTO
	SetTransaction(ID uint64, player entity.Player, txType string, details string)
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

func (service *raceService) GetRaceAndPlayer(raceId uint64, userId uint64, isBigRace bool) (error, entity.Race, entity.Player) {
	race := service.GetRaceByRaceId(raceId, isBigRace)
	player := service.playerService.GetPlayerByUserIdAndRaceId(raceId, userId)

	if player.ID == 0 {
		return fmt.Errorf(storage.ErrorUndefinedUser), entity.Race{}, entity.Player{}
	} else if race.ID == 0 {
		return fmt.Errorf(storage.ErrorUndefinedGame), entity.Race{}, entity.Player{}
	}

	return nil, race, player
}

func (service *raceService) BusinessAction(raceId uint64, userId uint64, isBigRace bool) error {
	logger.Info("RaceService.BusinessAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId, isBigRace)

	if err != nil {
		return err
	}

	err = service.playerService.BuyBusiness(entity.CardBusiness{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		Symbol:      race.CurrentCard.Symbol,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
		Rule:        race.CurrentCard.Rule,
		Cost:        *race.CurrentCard.Cost,
		CashFlow:    race.CurrentCard.CashFlow,
	}, player)

	//this.#log.addLog(player, `Купил бизнес ${this.#card.symbol} за $${this.#card.cost}`);
	race.Respond(player.ID, race.CurrentPlayer.ID)
	err, _ = service.UpdateRace(&race)

	go service.SetTransaction(race.ID, player, entity.TxTypes.Business, storage.MessageYouBoughtBusiness)

	return err
}

func (service *raceService) RealEstateAction(raceId uint64, userId uint64, isBigRace bool) error {
	logger.Info("RaceService.RealEstateAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId, isBigRace)

	if err != nil {
		return err
	}

	if race.CurrentCard.Type != "realEstate" {
		return fmt.Errorf(storage.ErrorInvalidTypeOfCard)
	}

	err = service.playerService.BuyRealEstate(entity.CardRealEstate{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		Symbol:      race.CurrentCard.Symbol,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
		Rule:        race.CurrentCard.Rule,
		Cost:        *race.CurrentCard.Cost,
		CashFlow:    race.CurrentCard.CashFlow,
		Mortgage:    race.CurrentCard.Mortgage,
		DownPayment: race.CurrentCard.DownPayment,
	}, player)

	//this.#log.addLog(player, `Купил бизнес ${this.#card.symbol} за $${this.#card.cost}`);
	race.Respond(player.ID, race.CurrentPlayer.ID)
	err, _ = service.UpdateRace(&race)

	go service.SetTransaction(race.ID, player, entity.TxTypes.RealEstate, storage.MessageYouBoughtRealEstate)

	return err
}

func (service *raceService) DreamAction(raceId uint64, userId uint64) error {
	logger.Info("RaceService.DreamAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId, true)

	if err != nil {
		return err
	}

	err = service.playerService.BuyDream(entity.CardDream{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
		Cost:        *race.CurrentCard.Cost,
	}, player)

	//this.#log.addLog(player, `Купил бизнес ${this.#card.symbol} за $${this.#card.cost}`);

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
		go service.SetTransaction(race.ID, player, entity.TxTypes.Dream, storage.MessageYouBoughtDream)
	}

	return err
}

func (service *raceService) RiskBusinessAction(raceId uint64, userId uint64) (error, dto.RiskResponseDTO) {
	logger.Info("RaceService.RiskBusinessAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId, true)

	if err != nil {
		return err, dto.RiskResponseDTO{RolledDice: 0}
	}

	dice := objects.Dice{}
	roll := dice.Roll(0)
	rolledDice := roll[0] | 1

	var status bool

	err, status = service.playerService.BuyRiskBusiness(entity.CardRiskBusiness{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
		Cost:        *race.CurrentCard.Cost,
		Dices:       *race.CurrentCard.Dices,
		ExtraDices:  *race.CurrentCard.ExtraDices,
		Symbol:      race.CurrentCard.Symbol,
	}, player, rolledDice)

	if err == nil {
		if status {
			go service.SetTransaction(race.ID, player, entity.TxTypes.Business, storage.MessageSuccessRiskDeal)
			//this.setTransactionState('risk', player.username, messages.SUCCESS_RISK_DEAL, { type: 'success', timeout: 1000 });
			//this.#log.addLog(player, `Рискованный бизнес - ${this.#card.symbol} за $${this.#card.cost}`);
		} else {
			go service.SetTransaction(race.ID, player, entity.TxTypes.Business, storage.MessageFailRiskDeal)
			//this.setTransactionState('risk', player.username, messages.FAIL_RISK_DEAL, { type: 'warning', timeout: 1000 });
		}

		//this.#log.addLog(player, `Купил бизнес ${this.#card.symbol} за $${this.#card.cost}`);
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
	}

	return err, dto.RiskResponseDTO{RolledDice: rolledDice}
}

func (service *raceService) SellRealEstate(raceId uint64, userId uint64, realEstateId string) error {
	logger.Info("RaceService.SellRealEstate", map[string]interface{}{
		"raceId":       raceId,
		"userId":       userId,
		"realEstateId": realEstateId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId, false)

	if err != nil {
		return err
	}

	err = service.playerService.SellRealEstate(realEstateId, entity.CardMarketRealEstate{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		Symbol:      race.CurrentCard.Symbol,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
		Rule:        *race.CurrentCard.Rule,
		Value:       *race.CurrentCard.Value,
		Plus:        *race.CurrentCard.Plus,
	}, player)

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		go service.SetTransaction(race.ID, player, entity.TxTypes.RealEstate, storage.MessageRealEstateHasBeenSold)
	}

	return err
}

func (service *raceService) SellStocks(raceId uint64, userId uint64, count int) error {
	logger.Info("RaceService.SellStocks", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"count":  count,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId, false)

	if err != nil {
		return err
	}

	err = service.playerService.SellStocks(entity.CardStocks{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		Symbol:      race.CurrentCard.Symbol,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
		Rule:        *race.CurrentCard.Rule,
		Price:       *race.CurrentCard.Price,
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

func (service *raceService) SellGoldCoins(raceId uint64, userId uint64, count int) error {
	logger.Info("RaceService.SellGoldCoins", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"count":  count,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId, false)

	if err != nil {
		return err
	}

	err = service.playerService.SellGold(entity.CardPreciousMetals{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		Symbol:      race.CurrentCard.Symbol,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
		Count:       *race.CurrentCard.Count,
	}, player, count)

	if err == nil {
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
		go service.SetTransaction(race.ID, player, entity.TxTypes.Other, storage.MessageGoldHasBeenSold)
	}

	return err
}

func (service *raceService) RiskStocksAction(raceId uint64, userId uint64) (error, dto.RiskResponseDTO) {
	logger.Info("RaceService.RiskStocksAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId, true)

	if err != nil {
		return err, dto.RiskResponseDTO{RolledDice: 0}
	}

	dice := objects.Dice{}
	roll := dice.Roll(0)
	rolledDice := roll[0] | 1

	var status bool

	err, status = service.playerService.BuyRiskStocks(entity.CardRiskStocks{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
		Count:       *race.CurrentCard.Count,
		CostPerOne:  *race.CurrentCard.CostPerOne,
		Dices:       *race.CurrentCard.Dices,
		ExtraDices:  *race.CurrentCard.ExtraDices,
		Symbol:      race.CurrentCard.Symbol,
	}, player, rolledDice)

	if err == nil {
		if status {
			go service.SetTransaction(race.ID, player, entity.TxTypes.Stocks, storage.MessageSuccessRiskStocksDeal)
			//this.setTransactionState('risk', player.username, messages.SUCCESS_RISK_DEAL, { type: 'success', timeout: 1000 });
			//this.#log.addLog(player, `Рискованный бизнес - ${this.#card.symbol} за $${this.#card.cost}`);
		} else {
			go service.SetTransaction(race.ID, player, entity.TxTypes.Stocks, storage.MessageFailRiskStocksDeal)
			//this.setTransactionState('risk', player.username, messages.FAIL_RISK_DEAL, { type: 'warning', timeout: 1000 });
		}

		//this.#log.addLog(player, `Купил бизнес ${this.#card.symbol} за $${this.#card.cost}`);
		race.Respond(player.ID, race.CurrentPlayer.ID)
		err, _ = service.UpdateRace(&race)
	}

	return err, dto.RiskResponseDTO{RolledDice: rolledDice}
}

func (service *raceService) StocksAction(raceId uint64, userId uint64, count int) error {
	logger.Info("RaceService.StocksAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"count":  count,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId, false)

	if err != nil {
		return err
	}

	if race.CurrentCard.Type != "stock" {
		return fmt.Errorf(storage.ErrorInvalidTypeOfCard)
	}

	cardStocks := entity.CardStocks{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
		Symbol:      race.CurrentCard.Symbol,
		Rule:        *race.CurrentCard.Rule,
		Price:       *race.CurrentCard.Price,
		Increase:    race.CurrentCard.Increase,
		Decrease:    race.CurrentCard.Decrease,
		Count:       &count,
		OnlyYou:     race.CurrentCard.OnlyYou,
		Range:       race.CurrentCard.Range,
	}

	if race.CurrentCard.Increase != nil {
		err = service.playerService.IncreaseStocks(cardStocks, player)
	} else if race.CurrentCard.Decrease != nil {
		err = service.playerService.DecreaseStocks(cardStocks, player)
	} else {
		err = service.playerService.BuyStocks(cardStocks, player, true)
	}

	//this.#log.addLog(player, `Купил бизнес ${this.#card.symbol} за $${this.#card.cost}`);
	race.Respond(player.ID, race.CurrentPlayer.ID)
	err, _ = service.UpdateRace(&race)

	go service.SetTransaction(race.ID, player, entity.TxTypes.Stocks, storage.MessageYouBoughtStocks)

	return err
}

func (service *raceService) LotteryAction(raceId uint64, userId uint64, isBigRace bool) error {
	logger.Info("RaceService.LotteryAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId, isBigRace)

	if err != nil {
		return err
	}

	race.Respond(player.ID, race.CurrentPlayer.ID)
	err, _ = service.UpdateRace(&race)

	return nil
}

func (service *raceService) GoldCoinsAction(raceId uint64, userId uint64) error {
	logger.Info("RaceService.GoldCoinsAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId, false)

	if err != nil {
		return err
	}

	race.Respond(player.ID, race.CurrentPlayer.ID)
	err, _ = service.UpdateRace(&race)

	return nil
}

func (service *raceService) MlmAction(raceId uint64, userId uint64, isBigRace bool) error {
	logger.Info("RaceService.MlmAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId, isBigRace)

	if err != nil {
		return err
	}

	race.Respond(player.ID, race.CurrentPlayer.ID)
	err, _ = service.UpdateRace(&race)

	return nil
}

func (service *raceService) SkipAction(raceId uint64, userId uint64, isBigRace bool) error {
	logger.Info("RaceService.SkipAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId, isBigRace)

	if err != nil {
		return err
	}

	race.Respond(player.ID, race.CurrentPlayer.ID)
	err, _ = service.UpdateRace(&race)

	return nil
}

func (service *raceService) CharityAction(raceId uint64, userId uint64, isBigRace bool) error {
	logger.Info("RaceService.CharityAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId, isBigRace)

	if err != nil {
		return err
	}

	err = service.playerService.Charity(player)

	race.Respond(player.ID, race.CurrentPlayer.ID)

	err, race = service.UpdateRace(&race)

	helper.LogPrintJson(race)

	return err
}

func (service *raceService) BabyAction(raceId uint64, userId uint64) error {
	logger.Info("RaceService.BabyAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId, false)

	if err != nil {
		return err
	}

	if player.Babies > 2 {
		return fmt.Errorf(storage.ErrorYouHaveTooManyBabies)
	}

	player.BornBaby()

	err, _ = service.playerService.UpdatePlayer(&player)

	race.Respond(player.ID, race.CurrentPlayer.ID)
	err, race = service.UpdateRace(&race)

	return err
}

func (service *raceService) DoodadAction(raceId uint64, userId uint64) error {
	logger.Info("RaceService.DoodadAction", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId, false)

	if err != nil {
		return err
	}

	//service.playerService.Doodad()

	err, _ = service.playerService.UpdatePlayer(&player)

	race.Respond(player.ID, race.CurrentPlayer.ID)
	err, race = service.UpdateRace(&race)

	return err
}

func (service *raceService) MarketAction(raceId uint64, userId uint64, actionType string) error {
	logger.Info("RaceService.MarketAction", map[string]interface{}{
		"raceId":     raceId,
		"userId":     userId,
		"actionType": actionType,
	})

	err, race, player := service.GetRaceAndPlayer(raceId, userId, false)

	if err != nil {
		return err
	}

	cardMarket := entity.CardMarket{
		ID:                   race.CurrentCard.ID,
		Type:                 race.CurrentCard.Type,
		Heading:              race.CurrentCard.Heading,
		Symbol:               race.CurrentCard.Symbol,
		Description:          race.CurrentCard.Description,
		Rule:                 *race.CurrentCard.Rule,
		SubRule:              *race.CurrentCard.SubRule,
		Cost:                 race.CurrentCard.Cost,
		ApplicableToEveryOne: race.CurrentCard.ApplicableToEveryOne,
		Success:              race.CurrentCard.Success,
		Plus:                 race.CurrentCard.Plus,
	}

	switch actionType {
	case "damage":
		err = service.playerService.PayDamages(cardMarket, player)

		//player.payDamages(this.#card);
		//
		//this.setTransactionState('market', player.username, messages.YOU_HAVE_PAID_PROPERTY_DAMAGE, { type: 'warning' });
		//
		//this.#log.addLog(player, messages.YOU_HAVE_PAID_PROPERTY_DAMAGE);
		//this.respond(player.username);
		break

	case "realEstate":
		//err = service.playerService.MarketRealEstate(cardMarket, *player)
		break

	case "business":
		//err = service.playerService.MarketRealEstate(cardMarket, *player)
		break

	case "goldCoins":
		//err = service.playerService.MarketGoldCoins(cardMarket, *player)
		break

	case "lottery":
		//err = service.playerService.MarketLottery(cardMarket, *player)
		break
	}

	err, _ = service.UpdateRace(&race)

	return err

}

func (service *raceService) PaydayAction(raceId uint64, userId uint64, actionType string, isBigRace bool) error {
	logger.Info("RaceService.MarketAction", map[string]interface{}{
		"raceId":     raceId,
		"userId":     userId,
		"actionType": actionType,
	})

	err, _, player := service.GetRaceAndPlayer(raceId, userId, isBigRace)

	if err != nil {
		return err
	}

	if actionType == "payday" {
		service.playerService.Payday(player)
	} else if actionType == "cashFlowDay" {
		service.playerService.CashFlowDay(player)
	}

	//const status = this.#currentPlayer.payday();
	//this.#log.addLog(this.#currentPlayer, `received pay of $${this.#cashflow()}`);
	//
	//this.setTransactionState('payday', this.#currentPlayer.username, messages.PAYDAY_MESSAGE);
	//if (status) {
	//	this.#changeTurnIfNoCard(player.username);
	//}
	//return status;

	return nil
}

func (service *raceService) GetInjectedRace(ctx *gin.Context) entity.Race {
	raceId := ctx.MustGet("race_id").(string)
	var queryDTO dto.QueryBigRaceDTO
	errDTO := ctx.ShouldBind(&queryDTO)

	if errDTO != nil {
		res := request.BuildErrorResponse("Failed to process request", errDTO.Error(), request.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return entity.Race{}
	}

	id, err := strconv.Atoi(raceId)

	if err != nil {
		return entity.Race{}
	}

	return service.raceRepository.FindRaceById(uint64(id), queryDTO.IsBigRace)
}

func (service *raceService) GetRaceByRaceId(raceId uint64, isBigRace bool) entity.Race {
	logger.Info("RaceService.GetRaceByRaceId", map[string]interface{}{
		"raceId": raceId,
	})

	return service.raceRepository.FindRaceById(raceId, isBigRace)
}

func (service *raceService) GetRacePlayersByRaceId(raceId uint64) []dto.GetRacePlayerResponseDTO {
	logger.Info("RaceService.GetRacePlayersByRaceId", map[string]interface{}{
		"raceId": raceId,
	})

	players := service.playerService.GetAllPlayersByRaceId(raceId)

	racePlayers := make([]dto.GetRacePlayerResponseDTO, 0)

	for _, player := range players {
		racePlayer := service.playerService.GetFormattedPlayerResponse(player)
		racePlayers = append(racePlayers, racePlayer)
	}

	return racePlayers
}

func (service *raceService) GetFormattedRaceResponse(raceId uint64, isBigRace bool) dto.GetRaceResponseDTO {
	logger.Info("RaceService.GetFormattedRaceResponse", map[string]interface{}{
		"raceId": raceId,
	})

	race := service.GetRaceByRaceId(raceId, isBigRace)
	players := service.GetRacePlayersByRaceId(raceId)
	_, player := service.playerService.GetRacePlayer(raceId, race.CurrentPlayer.UserId)

	response := dto.GetRaceResponseDTO{
		Players:       players,
		TurnResponses: race.Responses,
		Status:        race.Status,
		DiceValues:    race.Dice,
		CurrentPlayer: player,
		CurrentCard:   race.CurrentCard,
		GameId:        race.ID,
		IsTurnEnded:   false,
		Logs:          race.Logs,
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
