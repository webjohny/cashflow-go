package repository_mocks

import "github.com/webjohny/cashflow-go/entity"

type MockPlayerRepository struct {
	AllByRaceIdFunc                   func(raceId uint64) []entity.Player
	UpdatePlayerFunc                  func(player *entity.Player) entity.Player
	InsertPlayerFunc                  func(player *entity.Player) entity.Player
	DeletePlayerFunc                  func(player *entity.Player)
	FindPlayerByIdFunc                func(ID uint64) entity.Player
	FindPlayerByUsernameFunc          func(username string) entity.Player
	FindPlayerByUsernameAndRaceIdFunc func(raceId uint64, username string) entity.Player
}

func (m *MockPlayerRepository) UpdatePlayer(player *entity.Player) entity.Player {
	if m.UpdatePlayerFunc != nil {
		return m.UpdatePlayerFunc(player)
	}
	return entity.Player{}
}

func (m *MockPlayerRepository) AllByRaceId(raceId uint64) []entity.Player {
	if m.AllByRaceIdFunc != nil {
		return m.AllByRaceIdFunc(raceId)
	}

	return make([]entity.Player, 0)
}

func (m *MockPlayerRepository) DeletePlayer(player *entity.Player) {
	if m.DeletePlayerFunc != nil {
		m.DeletePlayerFunc(player)
	}
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

func (m *MockPlayerRepository) InsertPlayer(player *entity.Player) entity.Player {
	if m.InsertPlayerFunc != nil {
		m.InsertPlayerFunc(player)
	}

	return entity.Player{}
}
