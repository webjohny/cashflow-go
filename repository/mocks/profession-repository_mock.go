package repository_mocks

import "github.com/webjohny/cashflow-go/entity"

type MockProfessionRepository struct {
	AllFunc                func() []entity.Profession
	FindProfessionByIdFunc func(ID uint64) entity.Profession
	PickProfessionFunc     func(excluded *[]int) entity.Profession
}

func (m MockProfessionRepository) All() []entity.Profession {
	//TODO implement me
	panic("implement me")
}

func (m MockProfessionRepository) FindProfessionById(ID uint64) entity.Profession {
	//TODO implement me
	panic("implement me")
}

func (m MockProfessionRepository) PickProfession(excluded *[]int) entity.Profession {
	//TODO implement me
	panic("implement me")
}
