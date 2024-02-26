package service

import (
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/repository"
)

type ProfessionService interface {
	GetAll() []entity.Profession
}

type professionService struct {
	professionRepository repository.ProfessionRepository
}

func NewProfessionService(professionRepository repository.ProfessionRepository) ProfessionService {
	return &professionService{
		professionRepository: professionRepository,
	}
}

func (service *professionService) GetAll() []entity.Profession {
	return service.professionRepository.All()
}
