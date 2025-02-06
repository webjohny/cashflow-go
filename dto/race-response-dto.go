package dto

import "github.com/webjohny/cashflow-go/entity"

type RacePlayerProfileResponseDTO struct {
	Income        RacePlayerIncomeResponseDTO      `json:"income,omitempty"`
	Babies        uint8                            `json:"babies,omitempty"`
	Expenses      map[string]int                   `json:"expenses,omitempty"`
	Assets        entity.PlayerAssets              `json:"assets,omitempty"`
	Liabilities   RacePlayerLiabilitiesResponseDTO `json:"liabilities,omitempty"`
	TotalIncome   int                              `json:"total_income"`
	TotalExpenses int                              `json:"total_expenses"`
	CashFlow      int                              `json:"cash_flow"`
	ExtraCashFlow int                              `json:"extra_cash_flow"`
	PassiveIncome int                              `json:"passive_income"`
	Cash          int                              `json:"cash"`
}

type RacePlayerTransactionsResponseDTO struct {
	CurrentCash int    `json:"current_cash"`
	UpdatedCash int    `json:"updated_cash"`
	Amount      int    `json:"amount"`
	Details     string `json:"details"`
}

type RacePlayerIncomeResponseDTO struct {
	RealEstates []entity.CardRealEstate `json:"realEstates"`
	Business    []entity.CardBusiness   `json:"business"`
	Salary      int                     `json:"salary"`
}

type RacePlayerLiabilitiesResponseDTO struct {
	RealEstates    []entity.CardRealEstate `json:"realEstates"`
	Business       []entity.CardBusiness   `json:"business"`
	BankLoan       int                     `json:"bankLoan"`
	HomeMortgage   int                     `json:"homeMortgage"`
	SchoolLoans    int                     `json:"schoolLoans"`
	CarLoans       int                     `json:"carLoans"`
	CreditCardDebt int                     `json:"creditCardDebt"`
}

type RacePlayerProfessionResponseDTO struct {
	ID         uint64 `json:"id"`
	Profession string `json:"profession"`
}

type RacePlayerInfoResponseDTO struct {
	ID                uint64                   `json:"id"`
	Dream             entity.PlayerDream       `json:"dream"`
	FullName          string                   `json:"fullName"`
	GoalPassiveIncome int                      `json:"goalPassiveIncome"`
	Conditions        entity.BigRaceConditions `json:"conditions"`
}

type GetRacePlayerResponseDTO struct {
	ID                uint64                              `json:"id"`
	UserId            uint64                              `json:"user_id"`
	Username          string                              `json:"username"`
	Role              string                              `json:"role"`
	Color             string                              `json:"color"`
	Profile           RacePlayerProfileResponseDTO        `json:"profile"`
	Info              RacePlayerInfoResponseDTO           `json:"info"`
	Profession        RacePlayerProfessionResponseDTO     `json:"profession"`
	IsRolledDice      bool                                `json:"is_rolled_dice"`
	LastPosition      uint8                               `json:"last_position"`
	Transactions      []RacePlayerTransactionsResponseDTO `json:"transactions"`
	CurrentPosition   uint8                               `json:"current_position"`
	ExtraDices        int                                 `json:"extra_dices"`
	Dices             []int                               `json:"dices,omitempty"`
	Notifications     []entity.PlayerNotification         `json:"notifications,omitempty"`
	DualDiceCount     int                                 `json:"dual_dice_count"`
	SkippedTurns      uint8                               `json:"skipped_turns"`
	AllowOnBigRace    bool                                `json:"allow_on_big_race"`
	GameIsCompleted   bool                                `json:"game_is_completed"`
	GoalPassiveIncome bool                                `json:"goal_passive_income"`
	GoalPersonalDream bool                                `json:"goal_personal_dream"`
	OnBigRace         bool                                `json:"on_big_race"`
	IsActive          bool                                `json:"is_active"`
	HasBankrupt       bool                                `json:"has_bankrupt"`
	AboutToBankrupt   string                              `json:"about_to_bankrupt"`
	HasMlm            bool                                `json:"has_mlm"`
}
