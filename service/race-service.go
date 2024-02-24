package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/objects"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/request"
	"github.com/webjohny/cashflow-go/storage"
	"net/http"
	"strconv"
)

type RaceService interface {
	BusinessAction(raceId uint64, username string, actionType string) error
	RealEstateAction(raceId uint64, username string, actionType string) error
	DreamAction(raceId uint64, username string, actionType string) error
	RiskBusinessAction(raceId uint64, username string, actionType string) (error, dto.RiskResponseDTO)
	RiskStocksAction(raceId uint64, username string, actionType string) (error, dto.RiskResponseDTO)
	StocksAction(raceId uint64, username string, actionType string, count int) error
	LotteryAction(raceId uint64, username string, actionType string) error
	GoldCoinsAction(raceId uint64, username string, actionType string) error
	MlmAction(raceId uint64, username string, actionType string) error
	SkipAction(raceId uint64, username string, actionType string) error
	PaydayAction(raceId uint64, username string, actionType string) error
	MarketAction(raceId uint64, username string, actionType string) error
	GetRaceAndPlayer(raceId uint64, username string) (error, *entity.Race, *entity.Player)
	GetInjectedRace(ctx *gin.Context) *entity.Race
	GetRaceByRaceId(raceId uint64) *entity.Race
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

func (service *raceService) GetRaceAndPlayer(raceId uint64, username string) (error, *entity.Race, *entity.Player) {
	race := service.GetRaceByRaceId(raceId)
	player := service.playerService.GetPlayerByUsername(username)

	if player == nil {
		return fmt.Errorf(storage.ErrorUndefinedUser), nil, nil
	} else if race == nil {
		return fmt.Errorf(storage.ErrorUndefinedGame), nil, nil
	} else if race.CurrentCard == nil {
		return fmt.Errorf(storage.ErrorHaveNoDefinedCard), nil, nil
	}

	return nil, race, player
}

func (service *raceService) BusinessAction(raceId uint64, username string, actionType string) error {
	err, race, player := service.GetRaceAndPlayer(raceId, username)

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
	}, *player)

	//this.#log.addLog(player, `Купил бизнес ${this.#card.symbol} за $${this.#card.cost}`);
	race.Respond(player.ID, race.CurrentPlayer.ID)

	go service.SetTransaction(race.ID, *player, storage.MessageYouBoughtBusiness)

	return err
}

func (service *raceService) RealEstateAction(raceId uint64, username string, actionType string) error {
	err, race, player := service.GetRaceAndPlayer(raceId, username)

	if err != nil {
		return err
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
		Value:       *race.CurrentCard.Value,
		Plus:        *race.CurrentCard.Plus,
	}, *player)

	//this.#log.addLog(player, `Купил бизнес ${this.#card.symbol} за $${this.#card.cost}`);
	race.Respond(player.ID, race.CurrentPlayer.ID)

	go service.SetTransaction(race.ID, *player, storage.MessageYouBoughtRealEstate)

	return err
}

func (service *raceService) DreamAction(raceId uint64, username string, actionType string) error {
	err, race, player := service.GetRaceAndPlayer(raceId, username)

	if err != nil {
		return err
	}

	err = service.playerService.BuyDream(entity.CardDream{
		ID:          race.CurrentCard.ID,
		Type:        race.CurrentCard.Type,
		Heading:     race.CurrentCard.Heading,
		Description: race.CurrentCard.Description,
		Cost:        *race.CurrentCard.Cost,
	}, *player)

	//this.#log.addLog(player, `Купил бизнес ${this.#card.symbol} за $${this.#card.cost}`);
	race.Respond(player.ID, race.CurrentPlayer.ID)

	go service.SetTransaction(race.ID, *player, storage.MessageYouBoughtDream)

	return err
}

func (service *raceService) RiskBusinessAction(raceId uint64, username string, actionType string) (error, dto.RiskResponseDTO) {
	err, race, player := service.GetRaceAndPlayer(raceId, username)

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
	}, *player, rolledDice)

	if err == nil {
		if status {
			go service.SetTransaction(race.ID, *player, storage.MessageSuccessRiskDeal)
			//this.setTransactionState('risk', player.username, messages.SUCCESS_RISK_DEAL, { type: 'success', timeout: 1000 });
			//this.#log.addLog(player, `Рискованный бизнес - ${this.#card.symbol} за $${this.#card.cost}`);
		} else {
			go service.SetTransaction(race.ID, *player, storage.MessageFailRiskDeal)
			//this.setTransactionState('risk', player.username, messages.FAIL_RISK_DEAL, { type: 'warning', timeout: 1000 });
		}

		//this.#log.addLog(player, `Купил бизнес ${this.#card.symbol} за $${this.#card.cost}`);
		race.Respond(player.ID, race.CurrentPlayer.ID)
	}

	return err, dto.RiskResponseDTO{RolledDice: rolledDice}
}

func (service *raceService) RiskStocksAction(raceId uint64, username string, actionType string) (error, dto.RiskResponseDTO) {
	err, race, player := service.GetRaceAndPlayer(raceId, username)

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
	}, *player, rolledDice)

	if err == nil {
		if status {
			go service.SetTransaction(race.ID, *player, storage.MessageSuccessRiskStocksDeal)
			//this.setTransactionState('risk', player.username, messages.SUCCESS_RISK_DEAL, { type: 'success', timeout: 1000 });
			//this.#log.addLog(player, `Рискованный бизнес - ${this.#card.symbol} за $${this.#card.cost}`);
		} else {
			go service.SetTransaction(race.ID, *player, storage.MessageFailRiskStocksDeal)
			//this.setTransactionState('risk', player.username, messages.FAIL_RISK_DEAL, { type: 'warning', timeout: 1000 });
		}

		//this.#log.addLog(player, `Купил бизнес ${this.#card.symbol} за $${this.#card.cost}`);
		race.Respond(player.ID, race.CurrentPlayer.ID)
	}

	return err, dto.RiskResponseDTO{RolledDice: rolledDice}
}

func (service *raceService) StocksAction(raceId uint64, username string, actionType string, count int) error {
	err, race, player := service.GetRaceAndPlayer(raceId, username)

	if err != nil {
		return err
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
		Count:       race.CurrentCard.Count,
		OnlyYou:     race.CurrentCard.OnlyYou,
		Range:       race.CurrentCard.Range,
	}

	if race.CurrentCard.Increase != nil {
		err = service.playerService.IncreaseStocks(cardStocks, *player)
	} else if race.CurrentCard.Decrease != nil {
		err = service.playerService.DecreaseStocks(cardStocks, *player)
	} else {
		err = service.playerService.BuyStocks(cardStocks, *player, count, true)
	}

	//this.#log.addLog(player, `Купил бизнес ${this.#card.symbol} за $${this.#card.cost}`);
	race.Respond(player.ID, race.CurrentPlayer.ID)

	go service.SetTransaction(race.ID, *player, storage.MessageYouBoughtStocks)

	return err
}

func (service *raceService) LotteryAction(raceId uint64, username string, actionType string) error {
	//race := service.GetRaceByRaceId(raceId)

	return nil
}

func (service *raceService) GoldCoinsAction(raceId uint64, username string, actionType string) error {
	//race := service.GetRaceByRaceId(raceId)

	return nil
}

func (service *raceService) MlmAction(raceId uint64, username string, actionType string) error {
	//race := service.GetRaceByRaceId(raceId)

	return nil
}

func (service *raceService) SkipAction(raceId uint64, username string, actionType string) error {
	//race := service.GetRaceByRaceId(raceId)

	return nil
}

func (service *raceService) MarketAction(raceId uint64, username string, actionType string) error {
	err, race, player := service.GetRaceAndPlayer(raceId, username)

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
		err = service.playerService.PayDamages(cardMarket, *player)

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

	return err

}

func (service *raceService) PaydayAction(raceId uint64, username string, actionType string) error {
	err, _, player := service.GetRaceAndPlayer(raceId, username)

	if err != nil {
		return err
	}

	if actionType == "payday" {
		service.playerService.Payday(*player)
	} else if actionType == "cashFlowDay" {
		service.playerService.CashFlowDay(*player)
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

func (service *raceService) GetInjectedRace(ctx *gin.Context) *entity.Race {
	raceId := ctx.MustGet("race_id").(string)
	var queryDTO dto.QueryBigRaceDTO
	errDTO := ctx.ShouldBind(&queryDTO)

	if errDTO != nil {
		res := request.BuildErrorResponse("Failed to process request", errDTO.Error(), request.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return nil
	}

	id, err := strconv.Atoi(raceId)

	if err != nil {
		return nil
	}

	return service.raceRepository.FindRaceById(uint64(id), queryDTO.IsBigRace)
}

func (service *raceService) GetRaceByRaceId(raceId uint64) *entity.Race {
	return service.raceRepository.FindRaceById(raceId, false)
}

func (service *raceService) SetTransaction(ID uint64, player entity.Player, details string) {
	service.transactionService.InsertRaceTransaction(dto.TransactionCreateRaceDTO{
		RaceID:   ID,
		Details:  details,
		PlayerID: player.ID,
		Username: player.Username,
		Color:    player.Color,
	})
}
