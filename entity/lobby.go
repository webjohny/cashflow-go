package entity

var LobbyStatus = struct {
	NEW       string
	STARTED   string
	CANCELLED string
}{
	NEW:       "new",
	STARTED:   "started",
	CANCELLED: "cancelled",
}

type LobbyPlayer struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Color    string `json:"color"`
}

type Lobby struct {
	ID         uint64                 `gorm:"primary_key:auto_increment" json:"id"`
	Players    []LobbyPlayer          `gorm:"serializer:json" json:"players"`
	MaxPlayers int8                   `json:"max_players"`
	Status     string                 `json:"status"`
	Options    map[string]interface{} `gorm:"serializer:json" json:"options"`
	CreatedAt  string                 `json:"created_at"`
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
		RaceId:          raceId,
		Username:        username,
		Role:            player.Role,
		Color:           player.Color,
		Income:          profession.Income,
		Babies:          uint8(profession.Babies),
		Expenses:        profession.Expenses,
		Assets:          profession.Assets,
		Liabilities:     profession.Liabilities,
		Cash:            0,
		PassiveIncome:   0,
		ProfessionId:    uint8(profession.ID),
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

	instance.TotalExpenses = instance.CalculateTotalExpenses()
	instance.TotalIncome = instance.CalculateTotalIncome()
	instance.CashFlow = instance.CalculateCashFlow()

	return instance
}

func (l *Lobby) AddPlayer(username string, role string) {
	if l.Players == nil {
		l.Players = make([]LobbyPlayer, 0)
	}
	if !l.IsPlayerAlreadyJoined(username) {
		l.Players = append(l.Players, LobbyPlayer{Username: username, Role: role})
	}
}

func (l *Lobby) AddGuest(username string) {
	l.AddPlayer(username, PlayerRoles.GUEST)
}

func (l *Lobby) AddOwner(username string) {
	l.AddPlayer(username, PlayerRoles.OWNER)
}

func (l *Lobby) AddAdmin(username string) {
	l.AddPlayer(username, PlayerRoles.ADMIN)
}

func (l *Lobby) GetPlayer(username string) *LobbyPlayer {
	for _, player := range l.Players {
		if player.Username == username {
			return &player
		}
	}

	return nil
}

func (l *Lobby) IsFull() bool {
	return len(l.Players) == int(l.MaxPlayers)
}

func (l *Lobby) IsStarted() bool {
	return l.Status == LobbyStatus.STARTED
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
		if player.Role != PlayerRoles.ADMIN {
			count++
		}
	}

	return count < 2
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
