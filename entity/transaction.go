package entity

import (
	"gorm.io/datatypes"
)

var TransactionType = struct {
	PLAYER string
	RACE   string
}{
	PLAYER: "player",
	RACE:   "race",
}

var TxTypes = struct {
	Stocks     string
	RealEstate string
	Dream      string
	Other      string
	Business   string
}{
	Stocks:     "stocks",
	RealEstate: "realEstate",
	Dream:      "dream",
	Other:      "other",
	Business:   "business",
}

type TransactionData struct {
	CurrentCash *int    `json:"current_cash"`
	Cash        *int    `json:"cash"`
	Amount      *int    `json:"amount"`
	TxType      *string `json:"tx_type"`
	Username    *string `json:"username"`
	Color       *string `json:"color"`
}

type Transaction struct {
	ID              uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	PlayerID        *uint64          `gorm:"index" json:"player_id,omitempty"`
	RaceID          *uint64          `gorm:"index" json:"race_id,omitempty"`
	TransactionType string           `gorm:"type:varchar(20)" json:"transaction_type"` // Handle enum in application logic
	Details         string           `gorm:"type:varchar(255)" json:"description"`
	Data            *TransactionData `gorm:"type:json;serializer:json" json:"data,omitempty"`
	CreatedAt       datatypes.Date   `gorm:"column:created_at;type:datetime;default:current_timestamp;not null" json:"created_at"`
}
