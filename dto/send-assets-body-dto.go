package dto

type SendAssetsBodyDTO struct {
	Amount int    `json:"amount" form:"amount"`
	Asset  string `json:"asset" form:"asset"`
	Player string `json:"player" form:"player"`
}
