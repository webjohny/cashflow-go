package dto

import "github.com/webjohny/cashflow-go/entity"

type StartGameResponseDto struct {
	ID       uint64             `json:"id"`
	Options  entity.RaceOptions `json:"options"`
	Redirect string             `json:"redirect"`
}
