package dto

import "github.com/webjohny/cashflow-go/entity"

type RacePlayerProfileResponseDTO struct {
	Income        entity.PlayerIncome      `json:"income"`
	Babies        uint8                    `json:"babies"`
	Expenses      map[string]int           `json:"expenses"`
	Assets        entity.PlayerAssets      `json:"assets"`
	Liabilities   entity.PlayerLiabilities `json:"liabilities"`
	TotalIncome   int                      `json:"total_income"`
	TotalExpenses int                      `json:"total_expenses"`
	CashFlow      int                      `json:"cash_flow"`
	PassiveIncome int                      `json:"passive_income"`
	Cash          int                      `json:"cash"`
}

type RacePlayerTransactionsResponseDTO struct {
	CurrentCash int    `json:"current_cash"`
	Cash        int    `json:"cash"`
	Amount      int    `json:"amount"`
	TxType      string `json:"tx_type"`
	Details     string `json:"details"`
}

type RacePlayerResponseDTO struct {
	ID              uint64                              `json:"id"`
	UserId          uint64                              `json:"userId"`
	Username        string                              `json:"username"`
	Role            string                              `json:"role"`
	Color           string                              `json:"color"`
	Profile         RacePlayerProfileResponseDTO        `json:"profile"`
	Profession      entity.Profession                   `json:"profession"`
	IsRolledDice    bool                                `json:"is_rolled_dice"`
	LastPosition    uint8                               `json:"last_position"`
	Transactions    []RacePlayerTransactionsResponseDTO `json:"transactions"`
	CurrentPosition uint8                               `json:"current_position"`
	DualDiceCount   bool                                `json:"dual_dice_count"`
	SkippedTurns    bool                                `json:"skipped_turns"`
	CanReRoll       bool                                `json:"can_re_roll"`
	OnBigRace       bool                                `json:"on_big_race"`
	HasBankrupt     bool                                `json:"has_bankrupt"`
	AboutToBankrupt string                              `json:"about_to_bankrupt"`
	HasMlm          bool                                `json:"has_mlm"`
}
