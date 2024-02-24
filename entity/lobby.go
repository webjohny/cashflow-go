package entity

var LobbyStatus = struct {
	STARTED   string
	CANCELLED string
}{
	STARTED:   "started",
	CANCELLED: "cancelled",
}

type LobbyPlayer struct {
	Username string `json:"username"`
}

type Lobby struct {
	ID         uint64        `gorm:"primary_key:auto_increment" json:"id"`
	Players    []LobbyPlayer `json:"players"`
	MaxPlayers int8          `json:"max_players"`
	Status     string        `json:"status"`
	Options    string        `json:"options"`
	CreatedAt  string        `json:"created_at"`
}

func (l *Lobby) AddPlayer(username string) {
	if len(l.Players) > 0 {
		l.Players = append(l.Players, LobbyPlayer{Username: username})
	}
}

func (l *Lobby) RemovePlayer(username string) {
	if len(l.Players) > 0 {
		l.Players = append(l.Players, LobbyPlayer{Username: username})
	}
}
