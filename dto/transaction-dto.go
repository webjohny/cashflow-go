package dto

type TransactionCreatePlayerDTO struct {
	PlayerID    uint64 `json:"player_id" form:"player_id" binding:"required"`
	Details     string `json:"details" form:"details" binding:"required"`
	CurrentCash int    `json:"current_cash" form:"current_cash" binding:"required"`
	Cash        int    `json:"cash" form:"cash" binding:"required"`
	Amount      int    `json:"amount" form:"amount" binding:"required"`
}

type TransactionCreateRaceDTO struct {
	RaceID   uint64 `json:"race_id" form:"race_id" binding:"required"`
	CardID   string `json:"card_id" form:"card_id" binding:"required"`
	PlayerID uint64 `json:"player_id" form:"player_id" binding:"required"`
	Details  string `json:"details" form:"details" binding:"required"`
	CardType string `json:"card_type" form:"card_type" binding:"required"`
	Username string `json:"username" form:"username" binding:"required"`
	Color    string `json:"color" form:"color" binding:"required"`
}

type TransactionDTO struct {
	RaceID      uint64  `json:"race_id" form:"race_id" binding:"required"`
	CardID      *string `json:"card_id" form:"card_id" binding:"omitempty"`
	CardType    string  `json:"card_type" form:"card_type" binding:"required"`
	PlayerID    uint64  `json:"player_id" form:"player_id" binding:"required"`
	SenderID    *uint64 `json:"sender_id" form:"sender_id" binding:"omitempty"`
	Details     string  `json:"details" form:"details" binding:"required"`
	Username    string  `json:"username" form:"username" binding:"required"`
	Color       string  `json:"color" form:"color" binding:"required"`
	CurrentCash *int    `json:"current_cash,omitempty" form:"current_cash" binding:"omitempty"`
	UpdatedCash *int    `json:"updated_cash,omitempty" form:"updated_cash" binding:"omitempty"`
	Amount      *int    `json:"amount,omitempty" form:"amount" binding:"required"`
}

type TransactionCardDTO struct {
	CardID   string `json:"card_id" form:"card_id" binding:"omitempty"`
	CardType string `json:"card_type" form:"card_type" binding:"required"`
	Details  string `json:"details" form:"details" binding:"required"`
}

type TransactionCreateDTO struct {
	RaceID   uint64 `json:"race_id" form:"race_id" binding:"required"`
	Details  string `json:"details" form:"details" binding:"required"`
	TxType   string `json:"tx_type" form:"tx_type" binding:"required"`
	Username string `json:"username" form:"username" binding:"required"`
	Color    string `json:"color" form:"color" binding:"required"`
}

type TransactionUpdateDTO struct {
	ID               uint64 `json:"id" form:"id" binding:"required"`
	UserID           uint64 `json:"userid" form:"userid" binding:"required"`
	TransactionType  string `json:"trxtype" form:"trxtype" binding:"required"`
	Date             string `json:"date" form:"date" binding:"required"`
	Description      string `json:"description" form:"description" binding:"required"`
	TransactionValue int    `json:"trxvalue" form:"trxvalue" binding:"required"`
	TransactionGroup string `json:"trxgroup" form:"trxgroup" binding:"required"`
}
