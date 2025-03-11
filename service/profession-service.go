package service

import (
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/repository"
)

type ProfessionService interface {
	GetAll(language string) (error, []entity.Profession)
	GetByID(ID uint64, language string) (error, entity.Profession)
	SetProfessions(data dto.ProfessionsSetBodyDTO)
	GetRandomProfession(language string, excluded *[]int) (error, entity.Profession)
}

type professionService struct {
	professionRepository repository.ProfessionRepository
}

func NewProfessionService(professionRepository repository.ProfessionRepository) ProfessionService {
	return &professionService{
		professionRepository: professionRepository,
	}
}

func (service *professionService) GetAll(language string) (error, []entity.Profession) {
	return service.professionRepository.All(language)
}

func (service *professionService) GetByID(ID uint64, language string) (error, entity.Profession) {
	return service.professionRepository.FindProfessionById(ID, language)
}

func (service *professionService) GetRandomProfession(language string, excluded *[]int) (error, entity.Profession) {
	return service.professionRepository.PickProfession(language, excluded)
}

func (service *professionService) SetProfessions(professions dto.ProfessionsSetBodyDTO) {
	service.professionRepository.SetProfessions(professions)
}
