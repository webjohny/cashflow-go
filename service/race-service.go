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
	GetInjectedRace(ctx *gin.Context) *entity.Race
}

type raceService struct {
	raceRepository repository.RaceRepository
}

func NewRaceService(raceRepo repository.RaceRepository) RaceService {
	return &raceService{
		raceRepository: raceRepo,
	}
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
