package repository_mocks

import (
	"errors"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/storage"
)

type MockRaceRepository struct {
	InsertRaceFunc   func(race *entity.Race) (error, entity.Race)
	UpdateRaceFunc   func(race *entity.Race) (error, entity.Race)
	DeleteRaceFunc   func(race *entity.Race) error
	FindRaceByIdFunc func(ID uint64) entity.Race
	AllFunc          func() []entity.Race
}

func (m *MockRaceRepository) UpdateRace(race *entity.Race) (error, entity.Race) {
	if m.UpdateRaceFunc != nil {
		return m.UpdateRaceFunc(race)
	}
	return errors.New(storage.ErrorUndefinedGame), entity.Race{}
}

func (m *MockRaceRepository) DeleteRace(race *entity.Race) error {
	if m.DeleteRaceFunc != nil {
		return m.DeleteRaceFunc(race)
	}
	return errors.New(storage.ErrorUndefinedGame)
}

func (m *MockRaceRepository) FindRaceById(ID uint64) entity.Race {
	if m.FindRaceByIdFunc != nil {
		return m.FindRaceByIdFunc(ID)
	}
	return entity.Race{}
}

func (m *MockRaceRepository) InsertRace(race *entity.Race) (error, entity.Race) {
	if m.InsertRaceFunc != nil {
		return m.InsertRaceFunc(race)
	}
	return errors.New(storage.ErrorUndefinedGame), entity.Race{}
}

func (m *MockRaceRepository) All() []entity.Race {
	if m.AllFunc != nil {
		return m.AllFunc()
	}
	return []entity.Race{}
}
