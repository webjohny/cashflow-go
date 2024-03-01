package repository_mocks

import "github.com/webjohny/cashflow-go/entity"

type MockLobbyRepository struct {
	FindLobbyByIdFunc    func(ID uint64) *entity.Lobby
	AvailableToStartFunc func() bool
	RemovePlayerFunc     func(username string)
}

func (m *MockLobbyRepository) InsertLobby(b *entity.Lobby) entity.Lobby {
	//TODO implement me
	panic("implement me")
}

func (m *MockLobbyRepository) UpdateLobby(b *entity.Lobby) entity.Lobby {
	//TODO implement me
	panic("implement me")
}

func (m *MockLobbyRepository) All() []entity.Lobby {
	//TODO implement me
	panic("implement me")
}

func (m *MockLobbyRepository) DeleteLobby(b *entity.Lobby) {
	//TODO implement me
	panic("implement me")
}

func (m *MockLobbyRepository) FindLobbyById(ID uint64) *entity.Lobby {
	if m.FindLobbyByIdFunc != nil {
		return m.FindLobbyByIdFunc(ID)
	}
	return nil
}

func (m *MockLobbyRepository) AvailableToStart() bool {
	if m.AvailableToStartFunc != nil {
		return m.AvailableToStartFunc()
	}
	return false
}

func (m *MockLobbyRepository) RemovePlayer(username string) {
	if m.RemovePlayerFunc != nil {
		m.RemovePlayerFunc(username)
	}
}
