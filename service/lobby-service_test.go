package service_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	repository_mocks "github.com/webjohny/cashflow-go/repository/mocks"
	"github.com/webjohny/cashflow-go/service"
	"gorm.io/datatypes"
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
		CreatedAt:  datatypes.Date(time.Now()),
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
				Color:    l.GetPlayer(userOwner.Username).Color,
			}},
			MaxPlayers: lobbyDefault.MaxPlayers,
			Status:     lobbyDefault.Status,
			Options:    lobbyDefault.Options,
			CreatedAt:  lobbyDefault.CreatedAt,
		}
	}

	t.Run("Creating a Lobby", func(t *testing.T) {
		err, lobby := lobbyService.CreateLobby(userOwner.Username)
		assert.NoError(t, err)
		assert.Equal(t, lobby, entity.Lobby{
			ID: lobbyDefault.ID,
			Players: []entity.LobbyPlayer{{
				Username: userOwner.Username,
				Role:     userOwner.Role,
				Color:    lobby.GetPlayer(userOwner.Username).Color,
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
		CreatedAt:  datatypes.Date(time.Now()),
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
				Color:    l.GetPlayer(userOwner.Username).Color,
			}},
			MaxPlayers: lobbyDefault.MaxPlayers,
			Status:     lobbyDefault.Status,
			Options:    lobbyDefault.Options,
			CreatedAt:  lobbyDefault.CreatedAt,
		}
	}

	t.Run("Joining new guest to Lobby", func(t *testing.T) {
		err, lobby := lobbyService.CreateLobby(userOwner.Username)

		joinLobby := entity.Lobby{
			ID:         1,
			MaxPlayers: lobby.MaxPlayers,
			Players:    lobby.Players,
		}

		if lobby.ID != 0 {
			lobbyRepo.FindLobbyByIdFunc = func(ID uint64) entity.Lobby {
				return joinLobby
			}

			err, lobby = lobbyService.Join(lobby.ID, userGuest1.Username)
		}

		assert.NoError(t, err)
		assert.Contains(t, lobby.Players, entity.LobbyPlayer{
			Username: userGuest1.Username,
			Role:     userGuest1.Role,
			Color:    lobby.GetPlayer(userGuest1.Username).Color,
		})
	})

	t.Run("Joining wait list to Lobby", func(t *testing.T) {
		err, lobby := lobbyService.CreateLobby(userOwner.Username)

		joinLobby := entity.Lobby{
			ID:         1,
			MaxPlayers: lobby.MaxPlayers,
			Players:    lobby.Players,
			Status:     entity.LobbyStatus.Started,
			Options: map[string]interface{}{
				"enabled_wait_list": true,
			},
		}

		if lobby.ID != 0 {
			lobbyRepo.FindLobbyByIdFunc = func(ID uint64) entity.Lobby {
				return joinLobby
			}

			err, lobby = lobbyService.Join(lobby.ID, userWL1.Username)
		}

		assert.NoError(t, err)
		assert.Contains(t, lobby.Players, entity.LobbyPlayer{
			Username: userWL1.Username,
			Role:     userWL1.Role,
			Color:    lobby.GetPlayer(userWL1.Username).Color,
		})
	})
}
