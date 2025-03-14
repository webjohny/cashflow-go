package dto

type ModeratorUpdateStatusRaceDto struct {
	Status string `json:"status" binding:"required"`
}
