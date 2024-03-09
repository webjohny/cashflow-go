package service

type FinanceService interface {
	SendMoney(raceId uint64, username string, amount int, player string) error
	SendAssets(raceId uint64, username string, amount int, player string, asset string) error
	TakeLoan(raceId uint64, username string, amount int) error
}

type financeService struct {
	raceService RaceService
}

func NewFinanceService(raceService RaceService) FinanceService {
	return &financeService{
		raceService: raceService,
	}
}

func (service *financeService) SendMoney(raceId uint64, username string, amount int, receiver string) error {
	err, race, player := service.raceService.GetRaceAndPlayer(raceId, username)

	if err != nil {
		return err
	}
	
	return nil
}

func (service *financeService) SendAssets(raceId uint64, username string, amount int, receiver string, asset string) error {
	err, race, player := service.raceService.GetRaceAndPlayer(raceId, username)

	if err != nil {
		return err
	}

	return nil
}

func (service *financeService) TakeLoan(raceId uint64, username string, amount int) error {
	err, race, player := service.raceService.GetRaceAndPlayer(raceId, username)

	if err != nil {
		return err
	}

	return nil
}
