package repository_mocks

import "github.com/webjohny/cashflow-go/entity"

type MockRaceRepository struct {
	InsertRaceFunc   func(race *entity.Race) entity.Race
	UpdateRaceFunc   func(race *entity.Race) entity.Race
	DeleteRaceFunc   func(race *entity.Race)
	FindRaceByIdFunc func(ID uint64, IsBigRace bool) entity.Race
	AllFunc          func() []entity.Race
}

func (m *MockRaceRepository) UpdateRace(race *entity.Race) entity.Race {
	if m.UpdateRaceFunc != nil {
		return m.UpdateRaceFunc(race)
	}
	return entity.Race{}
}

func (m *MockRaceRepository) DeleteRace(race *entity.Race) {
	if m.DeleteRaceFunc != nil {
		m.DeleteRaceFunc(race)
	}
}

func (m *MockRaceRepository) FindRaceById(ID uint64, IsBigRace bool) entity.Race {
	if m.FindRaceByIdFunc != nil {
		return m.FindRaceByIdFunc(ID, IsBigRace)
	}
	return entity.Race{}
}

func (m *MockRaceRepository) InsertRace(race *entity.Race) entity.Race {
	if m.InsertRaceFunc != nil {
		return m.InsertRaceFunc(race)
	}
	return entity.Race{}
}

func (m *MockRaceRepository) All() []entity.Race {
	if m.AllFunc != nil {
		return m.AllFunc()
	}
	return []entity.Race{}
}
