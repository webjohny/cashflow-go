package service

type CardService interface {
	Prepare(raceId uint64, family string, actionType string, username string) string
	Accept(raceId uint64, family string, actionType string, username string) string
	Purchase(raceId uint64, family string, actionType string, username string) string
	Selling(raceId uint64, family string, actionType string, username string) string
}

type cardService struct {
	gameService GameService
	raceService RaceService
}

func NewCardService(gameService GameService, raceService RaceService) CardService {
	return &cardService{
		gameService: gameService,
		raceService: raceService,
	}
}

func (service *cardService) Prepare(raceId uint64, family string, actionType string, username string) string {
	if actionType == "risk" || actionType == "riskStock" {
		service.raceService.PreRiskAction(raceId)
	}
	return ""
}

func (service *cardService) Accept(raceId uint64, family string, actionType string, username string) string {
	if family == "payday" {
		service.raceService.PaydayAction(raceId)
	}
	return ""
}

func (service *cardService) Purchase(raceId uint64, family string, actionType string, username string) string {
	return ""
}

func (service *cardService) Selling(raceId uint64, family string, actionType string, username string) string {
	return ""
}
