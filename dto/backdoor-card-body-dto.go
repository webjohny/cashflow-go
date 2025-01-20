package dto

import "github.com/webjohny/cashflow-go/entity"

type BackdoorCardBodyDTO struct {
	Card entity.Card `json:"card" form:"card"`
}
