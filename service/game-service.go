package service

import (
	"github.com/webjohny/cashflow-go/repository"
)

type GameService interface{}

type gameService struct {
	raceRepository repository.RaceRepository
}

func NewGameService(raceRepo repository.RaceRepository) GameService {
	return &gameService{
		raceRepository: raceRepo,
	}
}
