package dto

type TakeLoanBodyDTO struct {
	Amount int    `json:"amount" form:"amount"`
	Player string `json:"player" form:"player"`
}
