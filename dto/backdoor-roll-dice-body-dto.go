package dto

type BackdoorRollDiceBodyDto struct {
	Dices []int `json:"dices" form:"dices"`
}
