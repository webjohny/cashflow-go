package dto

type CardSellingActionDTO struct {
	ID    string `json:"id" form:"id"`
	Count int    `json:"count" form:"count"`
}
