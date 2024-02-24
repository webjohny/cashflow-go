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

type RaceResponse struct {
	ID        uint64
	Username  string
	Responded bool
}

type RacePlayer struct {
	ID       uint64
	Username string
}

type Race struct {
	ID                uint64         `gorm:"primary_key:auto_increment" json:"id"`
	Responses         []RaceResponse `json:"responses"`
	ParentID          uint64         `gorm:"index" json:"parent_id"`
	Status            string         `json:"status"`
	CurrentPlayer     RacePlayer     `json:"current_player"`
	CurrentCard       *CardDefault   `json:"current_card"`
	Notifications     string         `json:"notifications"`
	BankruptedPlayers string         `json:"bankrupted_players"`
	Logs              string         `json:"logs"`
	Dice              string         `json:"dice"`
	Options           string         `json:"options"`
	CreatedAt         string         `json:"created_at"`
}

func (r *Race) Respond(ID uint64, currentPlayerID uint64) {
	if len(r.Responses) > 0 {
		playerId := ID | currentPlayerID
		for i := 0; i < len(r.Responses); i++ {
			if playerId == r.Responses[i].ID {
				r.Responses[i].Responded = true
			}
		}
	}
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
