package service

type MoneyService interface{}

type moneyService struct{}

func NewMoneyService() MoneyService {
	return &moneyService{}
}
