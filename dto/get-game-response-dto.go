package dto

import "github.com/webjohny/cashflow-go/entity"

type GetGameResponseDTO struct {
	Username      string                    `json:"username"`
	You           entity.Player             `json:"you"`
	Hash          string                    `json:"hash"`
	Players       []entity.Player           `json:"players"`
	TurnResponses []entity.RaceResponse     `json:"turnResponses"`
	Status        string                    `json:"status"`
	DiceValues    []int                     `json:"diceValues"`
	CurrentPlayer *RacePlayerResponseDTO    `json:"currentPlayer"`
	GameId        uint64                    `json:"gameId"`
	IsTurnEnded   bool                      `json:"isTurnEnded"`
	Logs          []entity.RaceLog          `json:"logs"`
	Notifications []entity.RaceNotification `json:"notifications"`
	Transaction   entity.TransactionData    `json:"transaction"`
}

type GetRaceResponseDTO struct {
	Players       []entity.Player           `json:"players"`
	TurnResponses []entity.RaceResponse     `json:"turnResponses"`
	Status        string                    `json:"status"`
	DiceValues    []int                     `json:"diceValues"`
	CurrentPlayer RacePlayerResponseDTO     `json:"currentPlayer"`
	GameId        uint64                    `json:"gameId"`
	IsTurnEnded   bool                      `json:"isTurnEnded"`
	Logs          []entity.RaceLog          `json:"logs"`
	Notifications []entity.RaceNotification `json:"notifications"`
	Transaction   entity.TransactionData    `json:"transaction"`
}
