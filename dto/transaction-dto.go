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
