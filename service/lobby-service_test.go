package service_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	repository_mocks "github.com/webjohny/cashflow-go/repository/mocks"
	"github.com/webjohny/cashflow-go/service"
	"testing"
	"time"
)

func TestCreatingLobby(t *testing.T) {

	// Set up mocks
	lobbyRepo := &repository_mocks.MockLobbyRepository{}

	lobbyDefault := entity.Lobby{
		ID:         1,
		Players:    make([]entity.LobbyPlayer, 0),
		MaxPlayers: service.LobbyMaxPlayers,
		Status:     entity.LobbyStatus.New,
		Options:    make(map[string]interface{}),
		CreatedAt:  time.Now(),
	}

	lobbyService := service.NewLobbyService(lobbyRepo)

	userOwner := entity.LobbyPlayer{
		Username: "user1",
		Role:     entity.PlayerRoles.Owner,
		Color:    helper.PickColor(),
	}

	lobbyRepo.InsertLobbyFunc = func(l *entity.Lobby) entity.Lobby {
		return entity.Lobby{
			ID: lobbyDefault.ID,
			Players: []entity.LobbyPlayer{{
				Username: userOwner.Username,
				Role:     userOwner.Role,
				Color:    l.GetPlayer(userOwner.ID).Color,
			}},
			MaxPlayers: lobbyDefault.MaxPlayers,
			Status:     lobbyDefault.Status,
			Options:    lobbyDefault.Options,
			CreatedAt:  lobbyDefault.CreatedAt,
		}
	}

	t.Run("Creating a Lobby", func(t *testing.T) {
		err, lobby := lobbyService.Create(userOwner.Username, userOwner.ID)
		assert.NoError(t, err)
		assert.Equal(t, lobby, entity.Lobby{
			ID: lobbyDefault.ID,
			Players: []entity.LobbyPlayer{{
				Username: userOwner.Username,
				Role:     userOwner.Role,
				Color:    lobby.GetPlayer(userOwner.ID).Color,
			}},
			MaxPlayers: lobbyDefault.MaxPlayers,
			Status:     lobbyDefault.Status,
			Options:    lobbyDefault.Options,
			CreatedAt:  lobbyDefault.CreatedAt,
		})
	})
}

func TestJoinToLobby(t *testing.T) {

	// Set up mocks
	lobbyRepo := &repository_mocks.MockLobbyRepository{}

	lobbyDefault := entity.Lobby{
		ID:         1,
		Players:    make([]entity.LobbyPlayer, 0),
		MaxPlayers: service.LobbyMaxPlayers,
		Status:     entity.LobbyStatus.New,
		Options:    make(map[string]interface{}),
		CreatedAt:  time.Now(),
	}

	lobbyService := service.NewLobbyService(lobbyRepo)

	userOwner := entity.LobbyPlayer{
		Username: "user1",
		Role:     entity.PlayerRoles.Owner,
		Color:    helper.PickColor(),
	}
	userGuest1 := entity.LobbyPlayer{
		Username: "user2",
		Role:     entity.PlayerRoles.Guest,
		Color:    helper.PickColor(),
	}
	userWL1 := entity.LobbyPlayer{
		Username: "user1wl",
		Role:     entity.PlayerRoles.WaitList,
		Color:    helper.PickColor(),
	}

	lobbyRepo.InsertLobbyFunc = func(l *entity.Lobby) entity.Lobby {
		return entity.Lobby{
			ID: lobbyDefault.ID,
			Players: []entity.LobbyPlayer{{
				Username: userOwner.Username,
				Role:     userOwner.Role,
				Color:    l.GetPlayer(userOwner.ID).Color,
			}},
			MaxPlayers: lobbyDefault.MaxPlayers,
			Status:     lobbyDefault.Status,
			Options:    lobbyDefault.Options,
			CreatedAt:  lobbyDefault.CreatedAt,
		}
	}

	t.Run("Joining new guest to Lobby", func(t *testing.T) {
		err, lobby := lobbyService.Create(userOwner.Username, userOwner.ID)

		joinLobby := entity.Lobby{
			ID:         1,
			MaxPlayers: lobby.MaxPlayers,
			Players:    lobby.Players,
		}

		var lobbyPlayer entity.LobbyPlayer

		if lobby.ID != 0 {
			lobbyRepo.FindLobbyByIdFunc = func(ID uint64) entity.Lobby {
				return joinLobby
			}

			err, lobbyPlayer = lobbyService.Join(lobby.ID, userGuest1.Username, userGuest1.ID)
		}

		assert.NoError(t, err)
		assert.Contains(t, lobby.Players, lobbyPlayer)
	})

	t.Run("Joining wait list to Lobby", func(t *testing.T) {
		err, lobby := lobbyService.Create(userOwner.Username, userOwner.ID)

		joinLobby := entity.Lobby{
			ID:         1,
			MaxPlayers: lobby.MaxPlayers,
			Players:    lobby.Players,
			Status:     entity.LobbyStatus.Started,
			Options: map[string]interface{}{
				"enable_wait_list": true,
			},
		}

		var lobbyPlayer entity.LobbyPlayer

		if lobby.ID != 0 {
			lobbyRepo.FindLobbyByIdFunc = func(ID uint64) entity.Lobby {
				return joinLobby
			}

			err, lobbyPlayer = lobbyService.Join(lobby.ID, userWL1.Username, userWL1.ID)
		}

		assert.NoError(t, err)
		assert.Contains(t, lobby.Players, lobbyPlayer)
	})
}
