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
	transactionRepo := &repository_mocks.MockTransactionRepository{}
	playerRepo := &repository_mocks.MockPlayerRepository{}
	professionRepo := &repository_mocks.MockProfessionRepository{}

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
	professionService := service.NewProfessionService(professionRepo)
	transactionService := service.NewTransactionService(transactionRepo)
	playerService := service.NewPlayerService(playerRepo, professionService, transactionService)
	raceService := service.NewRaceService(raceRepo, playerService, transactionService)
	gameService := service.NewGameService(raceService, playerService, lobbyService, professionService)

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

	raceRepo.InsertRaceFunc = func(l *entity.Race) (error, entity.Race) {
		return nil, entity.Race{
			ID:                0,
			Responses:         []entity.RaceResponse{},
			Status:            entity.RaceStatus.STARTED,
			CurrentPlayer:     entity.RacePlayer{},
			CurrentCard:       entity.Card{},
			Notifications:     []entity.RaceNotification{},
			BankruptedPlayers: []entity.RaceBankruptPlayer{},
			Logs:              []entity.RaceLog{},
			Dice:              []int{},
			Options:           entity.RaceOptions{},
		}
	}

	lobbyRepo.InsertLobbyFunc = func(l *entity.Lobby) (error, entity.Lobby) {
		return nil, entity.Lobby{
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
		err, lobby := lobbyService.Create(userOwner.Username, userOwner.ID)

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

			err, _ = lobbyService.Join(lobby.ID, userGuest1.Username, userGuest1.ID)
		}

		// Start a game by the lobby
		err, ratRace := gameService.Start(lobby.ID)

		// Get game
		err, raceResponse := gameService.GetGame(ratRace.ID, 0, false)

		playerRepo.FindPlayerByUsernameFunc = func(username string) entity.Player {
			return entity.Player{
				UserID:          0,
				RaceID:          0,
				Username:        userGuest1.Username,
				Role:            userGuest1.Role,
				Color:           userGuest1.Color,
				Babies:          0,
				Expenses:        make(map[string]int),
				Assets:          entity.PlayerAssets{},
				Liabilities:     entity.PlayerLiabilities{},
				Cash:            0,
				TotalIncome:     0,
				TotalExpenses:   0,
				CashFlow:        0,
				PassiveIncome:   0,
				ProfessionID:    0,
				LastPosition:    0,
				CurrentPosition: 0,
				DualDiceCount:   0,
				SkippedTurns:    0,
				IsRolledDice:    0,
				CanReRoll:       0,
				OnBigRace:       false,
				HasBankrupt:     0,
				AboutToBankrupt: "",
				HasMlm:          0,
			}
		}

		raceRepo.FindRaceByIdFunc = func(ID uint64, IsBigRace bool) entity.Race {
			return entity.Race{
				Responses: make([]entity.RaceResponse, 0),
				Status:    entity.RaceStatus.STARTED,
				CurrentPlayer: entity.RacePlayer{
					ID:       1,
					Username: userGuest1.Username,
				},
				CurrentCard:       entity.Card{},
				Notifications:     make([]entity.RaceNotification, 0),
				BankruptedPlayers: make([]entity.RaceBankruptPlayer, 0),
				Logs:              make([]entity.RaceLog, 0),
				Dice:              make([]int, 0),
				Options: entity.RaceOptions{
					EnableWaitList: false,
				},
			}
		}

		player := entity.Player{
			ID:       0,
			UserID:   0,
			RaceID:   0,
			Username: userGuest1.Username,
			Role:     userGuest1.Role,
			Color:    userGuest1.Color,
			Babies:   0,
			Expenses: nil,
			Assets: entity.PlayerAssets{
				Dreams:      nil,
				OtherAssets: nil,
				RealEstates: nil,
				Business:    nil,
				Stocks:      nil,
				Savings:     0,
			},
			Liabilities: entity.PlayerLiabilities{
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
			ProfessionID:    0,
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
		}

		assert.NoError(t, err)
		assert.Equal(t, raceResponse, dto.GetGameResponseDTO{
			Username: userOwner.Username,
			You: dto.GetRacePlayerResponseDTO{
				Username: player.Username,
			},
			Hash:    "",
			Players: nil,
		})
	})
}
