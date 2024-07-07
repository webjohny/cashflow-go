package service

import (
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/repository"
)

type ProfessionService interface {
	GetAll() []entity.Profession
	GetByID(ID uint64) entity.Profession
	GetRandomProfession(excluded *[]int) entity.Profession
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

func (service *professionService) GetByID(ID uint64) entity.Profession {
	return service.professionRepository.FindProfessionById(ID)
}

func (service *professionService) GetRandomProfession(excluded *[]int) entity.Profession {
	return service.professionRepository.PickProfession(excluded)
}
