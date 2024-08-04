package entity

import (
	"github.com/webjohny/cashflow-go/helper"
	"gorm.io/datatypes"
)

var LobbyStatus = struct {
	New       string
	Started   string
	Cancelled string
}{
	New:       "new",
	Started:   "started",
	Cancelled: "cancelled",
}

type LobbyPlayer struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Color    string `json:"color"`
}

type Lobby struct {
	ID         uint64                 `gorm:"primary_key:auto_increment" json:"id"`
	GameId     uint64                 `gorm:"index;type:int(11)" json:"game_id"`
	Players    []LobbyPlayer          `gorm:"type:json;serializer:json" json:"players"`
	MaxPlayers int8                   `gorm:"max_players:int(3)" json:"max_players"`
	Status     string                 `gorm:"status;type:enum('new','started','cancelled')" json:"status"`
	Options    map[string]interface{} `gorm:"type:json;serializer:json" json:"options"`
	CreatedAt  datatypes.Date         `gorm:"column:created_at;type:datetime;default:current_timestamp;not null" json:"created_at"`
}

func (l *Lobby) PreparePlayer(raceId uint64, username string, profession Profession) Player {
	var player LobbyPlayer

	for _, p := range l.Players {
		if p.Username == username {
			player = p
			break
		}
	}

	instance := Player{
		RaceID:          raceId,
		Username:        username,
		Role:            player.Role,
		Color:           player.Color,
		Salary:          profession.Income.Salary,
		Babies:          uint8(profession.Babies),
		Expenses:        profession.Expenses,
		Assets:          profession.Assets,
		Liabilities:     profession.Liabilities,
		Cash:            0,
		PassiveIncome:   0,
		ProfessionID:    uint8(profession.ID),
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

	instance.TotalExpenses = instance.CalculateTotalExpenses()
	instance.TotalIncome = instance.CalculateTotalIncome()
	instance.CashFlow = instance.CalculateCashFlow()

	return instance
}

func (l *Lobby) AddPlayer(userId uint64, username string, role string) {
	if !l.IsPlayerAlreadyJoined(username) {
		l.Players = append(l.Players, LobbyPlayer{ID: userId, Username: username, Role: role, Color: helper.PickColor()})
	}
}

func (l *Lobby) CountPlayers() int {
	return len(l.Players)
}

func (l *Lobby) AddWaitList(userId uint64, username string) {
	l.AddPlayer(userId, username, PlayerRoles.WaitList)
}

func (l *Lobby) AddGuest(userId uint64, username string) {
	l.AddPlayer(userId, username, PlayerRoles.Guest)
}

func (l *Lobby) AddOwner(userId uint64, username string) {
	l.AddPlayer(userId, username, PlayerRoles.Owner)
}

func (l *Lobby) AddAdmin(userId uint64, username string) {
	l.AddPlayer(userId, username, PlayerRoles.Admin)
}

func (l *Lobby) GetPlayer(userId uint64) LobbyPlayer {
	for _, player := range l.Players {
		if player.ID == userId {
			return player
		}
	}

	return LobbyPlayer{}
}

func (l *Lobby) IsFull() bool {
	return len(l.Players) == int(l.MaxPlayers)
}

func (l *Lobby) IsStarted() bool {
	return l.Status == LobbyStatus.Started
}

func (l *Lobby) IsGameStarted() bool {
	return l.IsStarted() && l.Options["enable_wait_list"] != nil && l.Options["enable_wait_list"] == true
}

func (l *Lobby) IsPlayerAlreadyJoined(username string) bool {
	for _, player := range l.Players {
		if player.Username == username {
			return true
		}
	}

	return false
}

func (l *Lobby) AddOption(key string, value interface{}) {
	l.Options[key] = value
}

func (l *Lobby) AvailableToStart() bool {
	var count int

	for i := 0; i < len(l.Players); i++ {
		player := l.Players[i]
		if player.Role != PlayerRoles.Admin {
			count++
		}
	}

	return count >= 2
}

func (l *Lobby) RemovePlayer(username string) {
	index := -1

	for i, player := range l.Players {
		if player.Username == username {
			index = i
			break
		}
	}

	if index != -1 {
		l.Players = append(l.Players[:index], l.Players[index+1:]...)
	}
}
