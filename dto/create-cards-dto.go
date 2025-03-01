package dto

import "github.com/webjohny/cashflow-go/entity"

type CreateCardsDTO struct {
	Type     string                   `json:"type"`
	Language string                   `json:"language"`
	Cards    map[string][]entity.Card `json:"cards" form:"cards"`
}
