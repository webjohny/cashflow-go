package entity

var RaceStatus = struct {
	STARTED   string
	LOBBY     string
	CANCELLED string
	FINISHED  string
}{
	STARTED:   "started",
	LOBBY:     "lobby",
	CANCELLED: "cancelled",
	FINISHED:  "finished",
}

type Race struct {
	ID                uint64 `gorm:"primary_key:auto_increment" json:"id"`
	Players           string `json:"players"`
	MaxPlayers        string `json:"max_players"`
	ParentID          uint64 `gorm:"index" json:"parent_id"`
	Status            string `json:"status"`
	CurrentPlayer     string `json:"current_player"`
	CurrentCard       string `json:"current_card"`
	Notifications     string `json:"notifications"`
	BankruptedPlayers string `json:"bankrupted_players"`
	Logs              string `json:"logs"`
	Dice              string `json:"dice"`
	Options           string `json:"options"`
	CreatedAt         string `json:"created_at"`
}

//numId: 22
//players: "[{"username":"webjohny","responded":true},{"username":"webtoolteam","responded":false}]"
//maxPlayers: 6
//status: "started"
//currentPlayer: "webjohny"
//notifications: "[]"
//logs: "[]"
//bankruptedPlayers: []
//createdAt: 1705974096391
//dice: [...]
//options: "{"enterAfterGameStarting":false}"
