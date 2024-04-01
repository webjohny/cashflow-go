package service_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	repository_mocks "github.com/webjohny/cashflow-go/repository/mocks"
	"github.com/webjohny/cashflow-go/service"
	"gorm.io/datatypes"
	"testing"
	"time"
)

func TestStartGame(t *testing.T) {

	// Set up mocks
	raceRepo := &repository_mocks.MockRaceRepository{}
	lobbyRepo := &repository_mocks.MockLobbyRepository{}
	playerRepo := &repository_mocks.MockPlayerRepository{}

	//raceDefault := entity.Lobby{
	//	ID:         1,
	//	Players:    make([]entity.LobbyPlayer, 0),
	//	MaxPlayers: service.LobbyMaxPlayers,
	//	Status:     entity.LobbyStatus.New,
	//	Options:    make(map[string]interface{}),
	//	CreatedAt:  datatypes.Date(time.Now()),
	//}

	lobbyDefault := entity.Lobby{
		ID:         1,
		Players:    make([]entity.LobbyPlayer, 0),
		MaxPlayers: service.LobbyMaxPlayers,
		Status:     entity.LobbyStatus.New,
		Options:    make(map[string]interface{}),
		CreatedAt:  datatypes.Date(time.Now()),
	}

	lobbyService := service.NewLobbyService(lobbyRepo)
	gameService := service.NewGameService(lobbyRepo, raceRepo, playerRepo)

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

	raceRepo.InsertRaceFunc = func(l *entity.Race) entity.Race {
		return entity.Race{
			ID:                0,
			Responses:         []entity.RaceResponse{},
			ParentID:          0,
			Status:            entity.RaceStatus.STARTED,
			CurrentPlayer:     &entity.RacePlayer{},
			CurrentCard:       &entity.Card{},
			Notifications:     []entity.RaceNotification{},
			BankruptedPlayers: []entity.RaceBankruptPlayer{},
			Logs:              []entity.RaceLog{},
			Dice:              []int{},
			Options:           entity.RaceOptions{},
		}
	}

	lobbyRepo.InsertLobbyFunc = func(l *entity.Lobby) entity.Lobby {
		return entity.Lobby{
			ID: lobbyDefault.ID,
			Players: []entity.LobbyPlayer{
				{
					Username: userOwner.Username,
					Role:     userOwner.Role,
					Color:    userOwner.Color,
				},
			},
			MaxPlayers: lobbyDefault.MaxPlayers,
			Status:     lobbyDefault.Status,
			Options:    lobbyDefault.Options,
			CreatedAt:  lobbyDefault.CreatedAt,
		}
	}

	t.Run("Starting a Game", func(t *testing.T) {
		// Creating a lobby
		err, lobby := lobbyService.CreateLobby(userOwner.Username)

		joinLobby := entity.Lobby{
			ID:         1,
			MaxPlayers: lobby.MaxPlayers,
			Players:    lobby.Players,
		}

		// Join a guest to the lobby
		if lobby.ID != 0 {
			lobbyRepo.FindLobbyByIdFunc = func(ID uint64) entity.Lobby {
				return joinLobby
			}

			err, lobby = lobbyService.Join(lobby.ID, userGuest1.Username)
		}

		// Start a game by the lobby
		err, ratRace := gameService.Start(lobby.ID)

		// Get game
		err, raceResponse := gameService.GetGame(ratRace.ID, 0, userOwner.Username, nil)

		playerRepo.FindPlayerByUsernameFunc = func(username string) entity.Player {
			return entity.Player{
				UserId:          0,
				RaceId:          0,
				Username:        userGuest1.Username,
				Role:            userGuest1.Role,
				Color:           userGuest1.Color,
				Income:          entity.PlayerIncome{},
				Babies:          0,
				Expenses:        make(map[string]int),
				Assets:          entity.PlayerAssets{},
				Liabilities:     entity.PlayerLiabilities{},
				Cash:            0,
				TotalIncome:     0,
				TotalExpenses:   0,
				CashFlow:        0,
				PassiveIncome:   0,
				ProfessionId:    0,
				LastPosition:    0,
				CurrentPosition: 0,
				DualDiceCount:   0,
				SkippedTurns:    0,
				IsRolledDice:    0,
				CanReRoll:       0,
				OnBigRace:       0,
				HasBankrupt:     0,
				AboutToBankrupt: "",
				HasMlm:          0,
			}
		}

		raceRepo.FindRaceByIdFunc = func(ID uint64, IsBigRace bool) entity.Race {
			return entity.Race{
				Responses: make([]entity.RaceResponse, 0),
				ParentID:  0,
				Status:    entity.RaceStatus.STARTED,
				CurrentPlayer: &entity.RacePlayer{
					ID:       1,
					Username: userGuest1.Username,
				},
				CurrentCard:       nil,
				Notifications:     make([]entity.RaceNotification, 0),
				BankruptedPlayers: make([]entity.RaceBankruptPlayer, 0),
				Logs:              make([]entity.RaceLog, 0),
				Dice:              make([]int, 0),
				Options: entity.RaceOptions{
					EnableWaitList: false,
				},
			}
		}

		assert.NoError(t, err)
		assert.Equal(t, raceResponse, dto.GetGameResponseDTO{
			Username: userOwner.Username,
			You: entity.Player{
				ID:       0,
				UserId:   0,
				RaceId:   0,
				Username: userGuest1.Username,
				Role:     userGuest1.Role,
				Color:    userGuest1.Color,
				Income: entity.PlayerIncome{
					RealEstates: []entity.CardRealEstate{},
					Business:    []entity.CardBusiness{},
					Salary:      0,
				},
				Babies:   0,
				Expenses: nil,
				Assets: entity.PlayerAssets{
					Dreams:         nil,
					PreciousMetals: nil,
					RealEstates:    nil,
					Business:       nil,
					Stocks:         nil,
					Savings:        0,
				},
				Liabilities: entity.PlayerLiabilities{
					RealEstates:    nil,
					Business:       nil,
					BankLoan:       0,
					HomeMortgage:   0,
					SchoolLoans:    0,
					CarLoans:       0,
					CreditCardDebt: 0,
				},
				Cash:            0,
				TotalIncome:     0,
				TotalExpenses:   0,
				CashFlow:        0,
				PassiveIncome:   0,
				ProfessionId:    0,
				LastPosition:    0,
				CurrentPosition: 0,
				DualDiceCount:   0,
				SkippedTurns:    0,
				IsRolledDice:    0,
				CanReRoll:       0,
				OnBigRace:       0,
				HasBankrupt:     0,
				AboutToBankrupt: "",
				HasMlm:          0,
				CreatedAt:       datatypes.Date{},
			},
			Hash:    "",
			Players: nil,
			Race: &entity.Race{
				ID:        1,
				Responses: nil,
				ParentID:  0,
				Status:    entity.RaceStatus.STARTED,
				CurrentPlayer: &entity.RacePlayer{
					ID:       1,
					Username: userGuest1.Username,
				},
				CurrentCard:       nil,
				Notifications:     make([]entity.RaceNotification, 0),
				BankruptedPlayers: make([]entity.RaceBankruptPlayer, 0),
				Logs:              make([]entity.RaceLog, 0),
				Dice:              make([]int, 0),
				Options: entity.RaceOptions{
					EnableWaitList: false,
				},
			},
			Lobby: &entity.Lobby{
				ID:         0,
				Players:    nil,
				MaxPlayers: 0,
				Status:     "",
				Options:    nil,
				CreatedAt:  datatypes.Date{},
			},
		})
	})
}
