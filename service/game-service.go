package service

type GameService interface {
}

type gameService struct {
	raceService RaceService
}

func NewGameService(raceService RaceService) GameService {
	return &gameService{
		raceService: raceService,
	}
}
