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

type GameService interface {
	GetRaceByCtx(ctx *gin.Context) entity.Race
}

type gameService struct {
	raceRepository repository.RaceRepository
}

func NewGameService(raceRepo repository.RaceRepository) GameService {
	return &gameService{
		raceRepository: raceRepo,
	}
}

func (service *gameService) GetRaceByCtx(ctx *gin.Context) entity.Race {
	raceId := ctx.MustGet("race_id").(string)
	var queryDTO dto.QueryBigRaceDTO
	errDTO := ctx.ShouldBind(&queryDTO)

	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	id, err := strconv.Atoi(raceId)

	if err != nil {
		panic(err)
	}

	var race entity.Race

	race = service.raceRepository.FindRaceById(uint64(id), queryDTO.IsBigRace)

	return race
}
