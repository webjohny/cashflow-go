package dto

type TakeLoanBodyDTO struct {
	Amount int `json:"amount" form:"amount"`
}
