package entity

import "gorm.io/datatypes"

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
	ID              uint64           `gorm:"primary_key:auto_increment" json:"id"`
	PlayerID        *uint64          `gorm:"index" json:"user_id,omitempty"`
	RaceID          *uint64          `gorm:"index" json:"race_id,omitempty"`
	TransactionType string           `json:"transaction_type"`
	Details         string           `json:"description"`
	Data            *TransactionData `gorm:"serializer:json" json:"data,omitempty"`
	CreatedAt       datatypes.Date   `json:"created_at"`
}
