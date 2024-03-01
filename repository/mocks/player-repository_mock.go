package repository_mocks

import "github.com/webjohny/cashflow-go/entity"

type MockPlayerRepository struct {
	InsertPlayerFunc func(player *entity.Player)
}

func (m *MockPlayerRepository) UpdatePlayer(b *entity.Player) entity.Player {
	//TODO implement me
	panic("implement me")
}

func (m *MockPlayerRepository) All(idUser string) []entity.Player {
	//TODO implement me
	panic("implement me")
}

func (m *MockPlayerRepository) DeletePlayer(b *entity.Player) {
	//TODO implement me
	panic("implement me")
}

func (m *MockPlayerRepository) FindPlayerById(ID uint64) *entity.Player {
	//TODO implement me
	panic("implement me")
}

func (m *MockPlayerRepository) FindPlayerByUsername(username string) *entity.Player {
	//TODO implement me
	panic("implement me")
}

func (m *MockPlayerRepository) InsertPlayer(b *entity.Player) entity.Player {
	if m.InsertPlayerFunc != nil {
		m.InsertPlayerFunc(b)
	}

	return entity.Player{}
}
