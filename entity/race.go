package entity

import (
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/objects"
	"math/rand"
	"time"
)

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
	UserId    uint64 `json:"user_id,omitempty"`
	Username  string `json:"username,omitempty"`
	Responded bool   `json:"responded"`
}

type RacePlayer struct {
	ID       uint64 `json:"id,omitempty"`
	UserId   uint64 `json:"user_id,omitempty"`
	Username string `json:"username"`
}

type RaceBankruptPlayer struct {
	ID         uint64 `json:"id,omitempty"`
	Username   string `json:"username"`
	CountDices int    `json:"count_dices"`
}

type BigRaceConditions struct {
	CashFlow      bool            `json:"cashFlow,omitempty"`
	RequiredDeals map[string]bool `json:"requiredDeals,omitempty"`
	CashPerMonths int             `json:"cashPerMonths,omitempty"`
}

type RaceOptions struct {
	BigRaceGoalForPassiveIncome int               `json:"bigRaceGoalForPassiveIncome,omitempty"`
	BigRaceConditions           BigRaceConditions `json:"bigRaceConditions,omitempty"`
	EnableManager               bool              `json:"enableManager,omitempty"`
	HideCards                   bool              `json:"hideCards,omitempty"`
	HandMode                    bool              `json:"handMode,omitempty"`
	BannerLink                  string            `json:"bannerLink,omitempty"`
	BannerImage                 string            `json:"bannerImage,omitempty"`
	MeetLink                    string            `json:"meetLink,omitempty"`
	EnableWaitList              bool              `json:"enableWaitList,omitempty"`
	CardCollection              string            `json:"cardCollection,omitempty"`
}

func (c *RaceOptions) Merge(override RaceOptions) {
	if override.BigRaceGoalForPassiveIncome != c.BigRaceGoalForPassiveIncome {
		c.BigRaceGoalForPassiveIncome = override.BigRaceGoalForPassiveIncome
	}

	if override.EnableManager != c.EnableManager {
		c.EnableManager = override.EnableManager
	}
	if override.HideCards != c.HideCards {
		c.HideCards = override.HideCards
	}
	if override.HandMode != c.HandMode {
		c.HandMode = override.HandMode
	}
	if override.EnableWaitList != c.EnableWaitList {
		c.EnableWaitList = override.EnableWaitList
	}
	if override.CardCollection != "" {
		c.CardCollection = override.CardCollection
	}
}

type RaceCardMap struct {
	Active map[string]int   `json:"active,omitempty"`
	Map    map[string][]int `json:"map,omitempty"`
}

func (rcm *RaceCardMap) HasMapping() bool {
	return len(rcm.Map) > 0
}

func (rcm *RaceCardMap) SetMap(cards map[string][]Card) {
	rcm.Map = make(map[string][]int)

	for action, items := range cards {
		var slice []int
		for key := range items {
			slice = append(slice, key)
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(slice), func(i, j int) {
			slice[i], slice[j] = slice[j], slice[i]
		})

		if _, ok := rcm.Map[action]; !ok {
			rcm.Map[action] = []int{}
		}

		rcm.Map[action] = slice
	}

	helper.LogPrintJson(rcm.Map)
}

func (rcm *RaceCardMap) Next(action string) {
	if rcm.Active == nil {
		rcm.Active = make(map[string]int)
	}

	current := rcm.Active[action]

	var index int
	for i, item := range rcm.Map[action] {
		if current == item {
			index = i + 1
		}
	}

	var nextItem = index

	if index > len(rcm.Map[action])-1 {
		nextItem = 0
	}

	rcm.Active[action] = rcm.Map[action][nextItem]
}

type Race struct {
	ID                uint64               `gorm:"primary_key:auto_increment" json:"id"`
	Responses         []RaceResponse       `gorm:"type:json;serializer:json" json:"responses"`
	IsMultiFlow       bool                 `gorm:"is_multi_flow" json:"is_multi_flow"`
	Status            string               `gorm:"status;type:enum('lobby','started','cancelled','finished')" json:"status"`
	Hash              string               `gorm:"hash" json:"hash"`
	CurrentPlayer     RacePlayer           `gorm:"type:json;serializer:json" json:"current_player,omitempty"`
	CurrentCard       Card                 `gorm:"type:json;serializer:json" json:"current_card,omitempty"`
	Notifications     []RaceNotification   `gorm:"type:json;serializer:json" json:"notifications"`
	BankruptedPlayers []RaceBankruptPlayer `gorm:"type:json;serializer:json" json:"bankrupted_players"`
	Logs              []RaceLog            `gorm:"type:json;serializer:json" json:"logs"`
	Dice              []int                `gorm:"type:json;serializer:json" json:"dice"`
	Options           RaceOptions          `gorm:"type:json;serializer:json" json:"options"`
	CardMap           RaceCardMap          `gorm:"type:json;serializer:json" json:"card_map"`
	CreatedAt         time.Time            `gorm:"column:created_at;type:datetime;default:current_timestamp;not null" json:"created_at"`
}

func (r *Race) Respond(ID uint64, currentPlayerID uint64) {
	if len(r.Responses) > 0 {
		playerId := ID

		if ID == 0 {
			playerId = currentPlayerID
		}
		for i := 0; i < len(r.Responses); i++ {
			if playerId == r.Responses[i].ID {
				r.Responses[i].Responded = true
			}
		}
	}
}

func (r *Race) ResetResponses() {
	if len(r.Responses) > 0 {
		for i := 0; i < len(r.Responses); i++ {
			r.Responses[i].Responded = false
		}
	}
}

func (r *Race) RemoveResponsePlayer(playerId uint64) {
	if len(r.Responses) > 0 {
		for i := 0; i < len(r.Responses); i++ {
			if r.Responses[i].ID == playerId {
				r.Responses = append(r.Responses[:i], r.Responses[i+1:]...)
				return
			}
		}
	}
}

func (r *Race) IsReceived(username string) bool {
	if r.IsMultiFlow {
		return r.AreReceived()
	}

	if len(r.Responses) > 0 {
		for i := 0; i < len(r.Responses); i++ {
			if username == r.Responses[i].Username {
				return r.Responses[i].Responded
			}
		}
	}

	return false
}

func (r *Race) AreReceived() bool {
	if len(r.Responses) > 0 {
		for i := 0; i < len(r.Responses); i++ {
			if !r.Responses[i].Responded {
				return false
			}
		}
	}

	return true
}

func (r *Race) CalculateDices() int {
	var dices int
	for i := 0; i < len(r.Dice); i++ {
		dices += r.Dice[i]
	}
	return dices
}

func (r *Race) GetDice() objects.Dice {
	dice := 1

	if len(r.Dice) > 0 {
		dice = r.Dice[0]
	}

	return objects.NewDice(dice, 1, 6)
}

func (r *Race) GetNextPlayer() RaceResponse {
	players := r.Responses
	username := r.CurrentPlayer.Username

	for i, player := range players {
		if player.Username == username {
			nextIndex := (i + 1) % len(players)
			return players[nextIndex]
		}
	}

	return RaceResponse{}
}

func (r *Race) NextPlayer() {
	next := r.GetNextPlayer()

	r.CurrentPlayer.ID = next.ID
	r.CurrentPlayer.UserId = next.UserId
	r.CurrentPlayer.Username = next.Username
}

func (r *Race) PickCurrentPlayer(playerId int) {
	players := r.Responses

	var next RaceResponse

	for i, player := range players {
		if int(player.ID) == playerId {
			next = players[i]
		}
	}

	r.CurrentPlayer.ID = next.ID
	r.CurrentPlayer.UserId = next.UserId
	r.CurrentPlayer.Username = next.Username
}

func (r *Race) CalculateTotalSteps(diceValues []int, diceCount int) int {
	totalCount := diceValues[0]

	if diceCount == 2 {
		totalCount += diceValues[1]
	}

	return totalCount
}
