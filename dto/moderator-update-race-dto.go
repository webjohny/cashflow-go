package dto

type ModeratorUpdateRaceDto struct {
	Status        string `json:"status" binding:"required"`
	CurrentPlayer int    `json:"current_player" binding:"numeric"`
	Responses     []bool `json:"responses" binding:"required"`
}
