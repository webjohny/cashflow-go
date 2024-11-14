package dto

import "github.com/webjohny/cashflow-go/entity"

type ModeratorUpdatePlayerDto struct {
	Cash            int                               `json:"cash" binding:"numeric"`
	Savings         int                               `json:"savings" binding:"numeric"`
	CashFlow        int                               `json:"cashFlow" binding:"numeric"`
	Babies          int                               `json:"babies" binding:"numeric"`
	LastPosition    int                               `json:"lastPosition" binding:"numeric"`
	CurrentPosition int                               `json:"currentPosition" binding:"numeric"`
	SkippedTurns    int                               `json:"skippedTurns" binding:"numeric"`
	OnBigRace       bool                              `json:"onBigRace" binding:"boolean"`
	RealEstate      map[string]entity.CardRealEstate  `json:"realEstate,omitempty" binding:"omitempty,dive"`
	Business        map[string]entity.CardBusiness    `json:"business,omitempty" binding:"omitempty,dive"`
	Stocks          map[string]entity.CardStocks      `json:"stocks,omitempty" binding:"omitempty,dive"`
	Other           map[string]entity.CardOtherAssets `json:"other,omitempty" binding:"omitempty,dive"`
	Expenses        ModeratorUpdatePlayerExpenseDto   `json:"expenses" binding:"required"`
	Liabilities     ModeratorUpdatePlayerLiabilityDto `json:"liabilities" binding:"required"`
}

type ModeratorUpdatePlayerAssetDto struct {
	ID          int    `json:"id" binding:"required"`
	Heading     string `json:"heading" binding:"required"`
	Description string `json:"description" binding:"required"`
	Symbol      string `json:"symbol" binding:"required"`
	Mortgage    string `json:"mortgage" binding:"required,numeric"`
	Cost        string `json:"cost" binding:"required,numeric"`
	CashFlow    string `json:"cash_flow" binding:"required,numeric"`
	IsOwner     string `json:"is_owner" binding:"required,boolean"`
}

type ModeratorUpdatePlayerExpenseDto struct {
	BankLoanPayment int `json:"bankLoanPayment" binding:"numeric"`
	PerChildExpense int `json:"perChildExpense" binding:"numeric"`
}

type ModeratorUpdatePlayerLiabilityDto struct {
	BankLoan int `json:"bankLoan" binding:"numeric"`
}
