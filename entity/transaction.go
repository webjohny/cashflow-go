package entity

import (
	"gorm.io/datatypes"
)

var TransactionCardType = struct {
	Skip             string
	Stock            string
	Payday           string
	CashFlowDay      string
	RealEstate       string
	Damage           string
	Baby             string
	Charity          string
	BigCharity       string
	Bankrupt         string
	BigBankrupt      string
	Downsized        string
	Doodad           string
	Other            string
	Dream            string
	PayTax           string
	Lottery          string
	Business         string
	RiskBusiness     string
	RiskStock        string
	MarketBusiness   string
	MarketRealEstate string
	SellStock        string
	MarketOther      string
	SendMoney        string
	SendMoneyToBank  string
	SendAssets       string
	ReceiveMoney     string
	StartMoney       string
	ReceiveAssets    string
	PayLoan          string
	TakeLoan         string
}{
	Skip:             "skip",
	Stock:            "stock",
	RealEstate:       "realEstate",
	Payday:           "payday",
	CashFlowDay:      "cashFlowDay",
	Damage:           "damage",
	Baby:             "baby",
	Other:            "other",
	PayTax:           "payTax",
	Bankrupt:         "bankrupt",
	BigBankrupt:      "bigBankrupt",
	Charity:          "charity",
	BigCharity:       "bigCharity",
	Downsized:        "downsized",
	Doodad:           "doodad",
	Dream:            "dream",
	Lottery:          "lottery",
	Business:         "business",
	RiskBusiness:     "riskBusiness",
	RiskStock:        "riskStock",
	MarketBusiness:   "marketBusiness",
	MarketRealEstate: "marketRealEstate",
	SellStock:        "sellStock",
	MarketOther:      "marketOther",
	SendMoney:        "sendMoney",
	ReceiveMoney:     "receiveMoney",
	SendMoneyToBank:  "sendMoneyToBank",
	SendAssets:       "sendAssets",
	ReceiveAssets:    "receiveAssets",
	StartMoney:       "startMoney",
	PayLoan:          "payLoan",
	TakeLoan:         "takeLoan",
}

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
	CurrentCash *int    `json:"current_cash,omitempty"`
	UpdatedCash *int    `json:"updated_cash,omitempty"`
	Amount      *int    `json:"amount,omitempty"`
	Username    *string `json:"username,omitempty"`
	Color       *string `json:"color,omitempty"`
}

type Transaction struct {
	ID              uint             `gorm:"primaryKey;autoIncrement" json:"id"`
	PlayerID        *uint64          `gorm:"uniqueIndex:trx;index:idx_player" json:"player_id,omitempty"`
	RaceID          *uint64          `gorm:"uniqueIndex:trx;index:idx_player" json:"race_id,omitempty"`
	CardID          string           `gorm:"uniqueIndex:trx;type:varchar(150)" json:"card_id"`
	CardType        string           `gorm:"uniqueIndex:trx;type:varchar(20)" json:"card_type"`
	TransactionType string           `gorm:"type:varchar(20)" json:"transaction_type"` // Handle enum in application logic
	Details         string           `gorm:"type:varchar(255)" json:"description"`
	Data            *TransactionData `gorm:"type:json;serializer:json" json:"data,omitempty"`
	CreatedAt       datatypes.Date   `gorm:"column:created_at;type:datetime;default:current_timestamp;not null" json:"created_at"`
}
