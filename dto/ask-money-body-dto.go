package dto

type AskMoneyBodyDto struct {
	Amount  int    `json:"amount" form:"amount"`
	Message string `json:"message" form:"message"`
	Type    string `json:"type" form:"type"`
}
