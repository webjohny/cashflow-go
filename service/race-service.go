package service

import (
	"github.com/webjohny/cashflow-go/repository"
)

type RaceService interface{}

type raceService struct {
	raceRepository repository.RaceRepository
}

func NewRaceService(raceRepo repository.RaceRepository) RaceService {
	return &raceService{
		raceRepository: raceRepo,
	}
}
