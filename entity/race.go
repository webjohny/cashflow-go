package entity

import "gorm.io/datatypes"

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

type RaceNotification struct {
	AlertType   string                 `json:"alert_type"`
	Description string                 `json:"description"`
	Options     map[string]interface{} `json:"options"`
}

type RaceLog struct {
	Username string `json:"username"`
	Color    string `json:"color"`
	Message  string `json:"message"`
}

type RaceResponse struct {
	ID        uint64 `json:"id,omitempty"`
	Username  string `json:"username"`
	Responded bool   `json:"responded"`
}

type RacePlayer struct {
	ID       uint64 `json:"id,omitempty"`
	Username string `json:"username"`
}

type RaceBankruptPlayer struct {
	ID         uint64 `json:"id,omitempty"`
	Username   string `json:"username"`
	CountDices int    `json:"count_dices"`
}

type RaceOptions struct {
	EnterAfterGameStarting bool `json:"enter_after_game_starting"`
}

type Race struct {
	ID                uint64               `gorm:"primary_key:auto_increment" json:"id"`
	Responses         []RaceResponse       `gorm:"serializer:json" json:"responses"`
	ParentID          uint64               `gorm:"index" json:"parent_id"`
	Status            string               `json:"status"`
	CurrentPlayer     *RacePlayer          `gorm:"serializer:json" json:"current_player,omitempty"`
	CurrentCard       *Card                `gorm:"serializer:json" json:"current_card,omitempty"`
	Notifications     []RaceNotification   `gorm:"serializer:json" json:"notifications"`
	BankruptedPlayers []RaceBankruptPlayer `gorm:"serializer:json" json:"bankrupted_players"`
	Logs              []RaceLog            `gorm:"serializer:json" json:"logs"`
	Dice              []int                `gorm:"serializer:json" json:"dice"`
	Options           RaceOptions          `gorm:"serializer:json" json:"options"`
	CreatedAt         datatypes.Date       `json:"created_at"`
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
