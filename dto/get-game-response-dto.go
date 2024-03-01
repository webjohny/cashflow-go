package dto

import "github.com/webjohny/cashflow-go/entity"

type GetGameResponseDTO struct {
	Username string          `json:"username" form:"username"`
	You      entity.Player   `json:"you" form:"you"`
	Hash     string          `json:"hash" form:"hash"`
	Players  []entity.Player `json:"players" form:"players"`
	Race     *entity.Race    `json:"race" form:"race"`
	Lobby    *entity.Lobby   `json:"lobby" form:"lobby"`
}
