package service

type CardService interface{}

type cardService struct{}

func NewCardService() CardService {
	return &cardService{}
}

func (service *cardService) Prepare() {
	//
}

func (service *cardService) Accept() {
	//
}

func (service *cardService) Purchase() {
	//
}

func (service *cardService) Selling() {
	//
}
