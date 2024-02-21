package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/request"
	"net/http"
	"strconv"
)

type RaceService interface {
	PreRiskAction(raceId uint64, username string, actionType string)
	BusinessAction(raceId uint64, username string, actionType string) error
	RealEstateAction(raceId uint64, username string, actionType string) error
	DreamAction(raceId uint64, username string, actionType string)
	RiskBusinessAction(raceId uint64, username string, actionType string) error
	RiskStocksAction(raceId uint64, username string, actionType string) error
	StocksAction(raceId uint64, username string, actionType string, count int) error
	LotteryAction(raceId uint64, username string, actionType string)
	GoldCoinsAction(raceId uint64, username string, actionType string)
	MlmAction(raceId uint64, username string, actionType string)
	SkipAction(raceId uint64, username string, actionType string)
	PaydayAction(raceId uint64, username string, actionType string)
	MarketAction(raceId uint64, username string, actionType string) error
	GetInjectedRace(ctx *gin.Context) *entity.Race
	GetRaceByRaceId(raceId uint64) *entity.Race
}

type raceService struct {
	raceRepository repository.RaceRepository
	playerService  PlayerService
}

func NewRaceService(raceRepo repository.RaceRepository, playerService PlayerService) RaceService {
	return &raceService{
		raceRepository: raceRepo,
		playerService:  playerService,
	}
}

func (service *raceService) PreRiskAction(raceId uint64, username string, actionType string) {
	//race := service.GetRaceByRaceId(raceId)
}
func (service *raceService) BusinessAction(raceId uint64, username string, actionType string) error {
	return nil
}
func (service *raceService) RealEstateAction(raceId uint64, username string, actionType string) error {
	return nil
}
func (service *raceService) DreamAction(raceId uint64, username string, actionType string) {}
func (service *raceService) RiskBusinessAction(raceId uint64, username string, actionType string) error {
	return nil
}
func (service *raceService) RiskStocksAction(raceId uint64, username string, actionType string) error {
	return nil
}
func (service *raceService) StocksAction(raceId uint64, username string, actionType string, count int) error {
	return nil
}
func (service *raceService) LotteryAction(raceId uint64, username string, actionType string)   {}
func (service *raceService) GoldCoinsAction(raceId uint64, username string, actionType string) {}
func (service *raceService) MlmAction(raceId uint64, username string, actionType string)       {}
func (service *raceService) SkipAction(raceId uint64, username string, actionType string)      {}

func (service *raceService) MarketAction(raceId uint64, username string, actionType string) error {
	race := service.GetRaceByRaceId(raceId)

	if race.CurrentCard == nil {
		return fmt.Errorf(helper.GetMessage("HAVE_NO_DEFINED_CARD"))
	}

	player := service.playerService.GetPlayerByUsername(username)

	var err error

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

	if player != nil {
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

		case "goldCoins":
			//err = service.playerService.MarketGoldCoins(cardMarket, *player)
			break

		case "lottery":
			//err = service.playerService.MarketLottery(cardMarket, *player)
			break
		}
	}

	return err

}

func (service *raceService) PaydayAction(raceId uint64, username string, actionType string) {
	//race := service.GetRaceByRaceId(raceId)

	player := service.playerService.GetPlayerByUsername(username)

	if player != nil {
		if actionType == "payday" {
			service.playerService.Payday(*player)
		} else if actionType == "cashFlowDay" {
			service.playerService.CashFlowDay(*player)
		}
	}

	//const status = this.#currentPlayer.payday();
	//this.#log.addLog(this.#currentPlayer, `received pay of $${this.#cashflow()}`);
	//
	//this.setTransactionState('payday', this.#currentPlayer.username, messages.PAYDAY_MESSAGE);
	//if (status) {
	//	this.#changeTurnIfNoCard(player.username);
	//}
	//return status;
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
