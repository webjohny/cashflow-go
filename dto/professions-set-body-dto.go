package dto

import "github.com/webjohny/cashflow-go/entity"

type ProfessionsSetBodyDTO struct {
	Type        string              `json:"type"`
	Language    string              `json:"language"`
	Professions []entity.Profession `json:"professions" form:"professions"`
}
