package dto

type RollDiceDto struct {
	IsFinished bool `json:"is_finished"`
	Dices      int  `json:"dices"`
	DiceValue  int  `json:"dice_value,omitempty"`
}
