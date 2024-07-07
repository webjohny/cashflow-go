package dto

type CardPurchaseActionDTO struct {
	Count   int                           `json:"count" form:"count"`
	Players []CardPurchasePlayerActionDTO `json:"players,omitempty" form:"players"`
}

type CardPurchasePlayerActionDTO struct {
	ID      int `json:"id"`
	Passive int `json:"passive,omitempty"`
	Amount  int `json:"amount"`
	Percent int `json:"percent,omitempty"`
}
