package repository_mocks

import "github.com/webjohny/cashflow-go/entity"

type MockLobbyRepository struct {
	InsertLobbyFunc   func(lobby *entity.Lobby) entity.Lobby
	UpdateLobbyFunc   func(lobby *entity.Lobby) entity.Lobby
	DeleteLobbyFunc   func(lobby *entity.Lobby)
	FindLobbyByIdFunc func(ID uint64) entity.Lobby
	AllFunc           func() []entity.Lobby
}

func (m *MockLobbyRepository) InsertLobby(lobby *entity.Lobby) entity.Lobby {
	if m.InsertLobbyFunc != nil {
		return m.InsertLobbyFunc(lobby)
	}
	return entity.Lobby{}
}

func (m *MockLobbyRepository) UpdateLobby(lobby *entity.Lobby) entity.Lobby {
	if m.UpdateLobbyFunc != nil {
		return m.UpdateLobbyFunc(lobby)
	}
	return entity.Lobby{}
}

func (m *MockLobbyRepository) DeleteLobby(lobby *entity.Lobby) {
	if m.DeleteLobbyFunc != nil {
		m.DeleteLobbyFunc(lobby)
	}
}

func (m *MockLobbyRepository) FindLobbyById(ID uint64) entity.Lobby {
	if m.FindLobbyByIdFunc != nil {
		return m.FindLobbyByIdFunc(ID)
	}
	return entity.Lobby{}
}

func (m *MockLobbyRepository) All() []entity.Lobby {
	if m.AllFunc != nil {
		return m.AllFunc()
	}
	return []entity.Lobby{}
}
