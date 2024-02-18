package service

import (
	"github.com/gin-gonic/gin"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/repository"
	"net/http"
	"strconv"
)

type RaceService interface {
	PreRiskAction(raceId uint64, username string)
	PaydayAction(raceId uint64, username string)
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

func (service *raceService) PreRiskAction(raceId uint64, username string) {
	//race := service.GetRaceByRaceId(raceId)
}

func (service *raceService) PaydayAction(raceId uint64, username string) {
	//race := service.GetRaceByRaceId(raceId)

	player := service.playerService.GetPlayerByUsername(username)

	if player != nil {
		service.playerService.Payday(*player)
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
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
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
