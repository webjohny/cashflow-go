package repository_mocks

import "github.com/webjohny/cashflow-go/entity"

type MockRaceRepository struct {
	InsertRaceFunc func(race *entity.Race) *entity.Race
}

func (m *MockRaceRepository) UpdateRace(b *entity.Race) entity.Race {
	//TODO implement me
	panic("implement me")
}

func (m *MockRaceRepository) All(idUser string) []entity.Race {
	//TODO implement me
	panic("implement me")
}

func (m *MockRaceRepository) DeleteRace(b *entity.Race) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRaceRepository) FindRaceById(ID uint64, IsBigRace bool) *entity.Race {
	//TODO implement me
	panic("implement me")
}

func (m *MockRaceRepository) InsertRace(b *entity.Race) entity.Race {
	if m.InsertRaceFunc != nil {
		return *m.InsertRaceFunc(b)
	}
	return entity.Race{}
}
