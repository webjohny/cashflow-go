package dto

type ModeratorSendMoneyDTO struct {
	Message string `json:"message" form:"message"`
	Amount  int    `json:"amount" form:"amount"`
	Player  int    `json:"player" form:"player"`
}
