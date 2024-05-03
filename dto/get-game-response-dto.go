package dto

import "github.com/webjohny/cashflow-go/entity"

type GetGameResponseDTO struct {
	Username          string                     `json:"username"`
	You               GetRacePlayerResponseDTO   `json:"you"`
	Hash              string                     `json:"hash"`
	Players           []GetRacePlayerResponseDTO `json:"players"`
	BankruptedPlayers []GetRacePlayerResponseDTO `json:"bankrupted_players"`
	TurnResponses     []entity.RaceResponse      `json:"turn_responses"`
	Status            string                     `json:"status"`
	DiceValues        []int                      `json:"dice_values"`
	CurrentPlayer     *GetRacePlayerResponseDTO  `json:"current_player"`
	CurrentCard       *entity.Card               `json:"current_card"`
	GameId            uint64                     `json:"game_id"`
	IsTurnEnded       bool                       `json:"is_turn_ended"`
	Logs              []entity.RaceLog           `json:"logs"`
	Notifications     []entity.RaceNotification  `json:"notifications"`
	Transaction       entity.TransactionData     `json:"transaction"`
}

type GetRaceResponseDTO struct {
	Players       []GetRacePlayerResponseDTO `json:"players"`
	TurnResponses []entity.RaceResponse      `json:"turn_responses"`
	Status        string                     `json:"status"`
	DiceValues    []int                      `json:"dice_values"`
	CurrentPlayer GetRacePlayerResponseDTO   `json:"current_player"`
	CurrentCard   entity.Card                `json:"current_card"`
	GameId        uint64                     `json:"game_id"`
	IsTurnEnded   bool                       `json:"is_turn_ended"`
	Logs          []entity.RaceLog           `json:"logs"`
	Notifications []entity.RaceNotification  `json:"notifications"`
	Transaction   entity.TransactionData     `json:"transaction"`
}

type GetLobbyResponseDTO struct {
	Username string               `json:"username"`
	You      entity.LobbyPlayer   `json:"you"`
	Hash     string               `json:"hash"`
	GameId   uint64               `json:"game_id"`
	LobbyId  uint64               `json:"lobby_id"`
	Players  []entity.LobbyPlayer `json:"players"`
	Status   string               `json:"status"`
}
