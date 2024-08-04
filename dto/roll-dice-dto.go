package dto

type RollDiceDto struct {
	IsFinished bool `json:"is_finished"`
	Dices      int  `json:"dices"`
}
