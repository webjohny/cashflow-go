package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/repository/mocks"
	"github.com/webjohny/cashflow-go/service"
	"github.com/webjohny/cashflow-go/storage"
)

func TestGameService_Start(t *testing.T) {

	// Set up mocks
	lobbyRepo := &repository_mocks.MockLobbyRepository{}
	raceRepo := &repository_mocks.MockRaceRepository{}
	playerRepo := &repository_mocks.MockPlayerRepository{}

	gameService := service.NewGameService(lobbyRepo, raceRepo, playerRepo)

	t.Run("Remove Player from Lobby", func(t *testing.T) {
		lobbyRepo.FindLobbyByIdFunc = func(ID uint64) *entity.Lobby {
			return &entity.Lobby{
				ID:      1,
				Players: []entity.LobbyPlayer{{Username: "user1"}, {Username: "user2"}},
			}
		}

		err := gameService.Start(1, "user1")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(lobbyRepo.FindLobbyByIdFunc(1).Players))
	})

	t.Run("Insufficient Players in Lobby", func(t *testing.T) {
		lobbyRepo.FindLobbyByIdFunc = func(ID uint64) *entity.Lobby {
			return &entity.Lobby{
				ID:      2,
				Players: []entity.LobbyPlayer{{Username: "user1"}},
			}
		}
		lobbyRepo.AvailableToStartFunc = func() bool {
			return false
		}

		err := gameService.Start(2, "user1")
		assert.Error(t, err)
		assert.Equal(t, storage.ErrorInsufficientPlayers, err.Error())
	})

	t.Run("Start New Race", func(t *testing.T) {
		lobbyRepo.FindLobbyByIdFunc = func(ID uint64) *entity.Lobby {
			return &entity.Lobby{
				ID:      3,
				Players: []entity.LobbyPlayer{{Username: "user1"}, {Username: "user2"}},
			}
		}
		lobbyRepo.AvailableToStartFunc = func() bool {
			return true
		}
		raceRepo.InsertRaceFunc = func(race *entity.Race) *entity.Race {
			return &entity.Race{ID: 123}
		}

		err := gameService.Start(3, "user1")
		assert.NoError(t, err)
		// Add assertions based on the expected behavior after starting a new race
	})

	t.Run("Undefined Lobby", func(t *testing.T) {
		lobbyRepo.FindLobbyByIdFunc = nil

		err := gameService.Start(4, "user1")
		assert.Error(t, err)
		assert.Equal(t, storage.ErrorUndefinedLobby, err.Error())
	})
}
