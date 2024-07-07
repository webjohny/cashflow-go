package repository_mocks

import (
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/storage"
)

type MockPlayerRepository struct {
	InsertPlayerFunc                  func(b *entity.Player) (error, entity.Player)
	UpdatePlayerFunc                  func(b *entity.Player) (error, entity.Player)
	UpdateCashFunc                    func(b *entity.Player, cash int)
	AllByRaceIdFunc                   func(raceId uint64) []entity.Player
	DeletePlayerFunc                  func(b *entity.Player) error
	FindPlayerByIdFunc                func(ID uint64) entity.Player
	FindPlayerByUsernameFunc          func(username string) entity.Player
	FindPlayerByUsernameAndRaceIdFunc func(raceId uint64, username string) entity.Player
	FindPlayerByUserIdAndRaceIdFunc   func(raceId uint64, userId uint64) entity.Player
}

func (m *MockPlayerRepository) UpdatePlayer(player *entity.Player) (error, entity.Player) {
	if m.UpdatePlayerFunc != nil {
		return m.UpdatePlayerFunc(player)
	}
	return errors.New(storage.ErrorUndefinedPlayer), entity.Player{}
}

func (m *MockPlayerRepository) UpdateCash(player *entity.Player, cash int) {
	if m.UpdatePlayerFunc != nil {
		m.UpdateCashFunc(player, cash)
	}
}

func (m *MockPlayerRepository) AllByRaceId(raceId uint64) []entity.Player {
	if m.AllByRaceIdFunc != nil {
		return m.AllByRaceIdFunc(raceId)
	}

	return make([]entity.Player, 0)
}

func (m *MockPlayerRepository) DeletePlayer(player *entity.Player) error {
	if m.DeletePlayerFunc != nil {
		return m.DeletePlayerFunc(player)
	}

	return errors.New(storage.ErrorUndefinedPlayer)
}

func (m *MockPlayerRepository) FindPlayerById(ID uint64) entity.Player {
	if m.FindPlayerByIdFunc != nil {
		return m.FindPlayerByIdFunc(ID)
	}
	return entity.Player{}
}

func (m *MockPlayerRepository) FindPlayerByUsername(username string) entity.Player {
	if m.FindPlayerByUsernameFunc != nil {
		return m.FindPlayerByUsernameFunc(username)
	}
	return entity.Player{}
}

func (m *MockPlayerRepository) FindPlayerByUsernameAndRaceId(raceId uint64, username string) entity.Player {
	if m.FindPlayerByUsernameAndRaceIdFunc != nil {
		return m.FindPlayerByUsernameAndRaceIdFunc(raceId, username)
	}
	return entity.Player{}
}

func (m *MockPlayerRepository) FindPlayerByUserIdAndRaceId(raceId uint64, userId uint64) entity.Player {
	if m.FindPlayerByUsernameAndRaceIdFunc != nil {
		return m.FindPlayerByUserIdAndRaceIdFunc(raceId, userId)
	}
	return entity.Player{}
}

func (m *MockPlayerRepository) InsertPlayer(player *entity.Player) (error, entity.Player) {
	if m.InsertPlayerFunc != nil {
		return m.InsertPlayerFunc(player)
	}

	return errors.New(storage.ErrorUndefinedPlayer), entity.Player{}
}
