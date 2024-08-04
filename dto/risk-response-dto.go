package dto

type RiskResponseDTO struct {
	RolledDice int    `json:"rolled_dice"`
	Message    string `json:"message,omitempty"`
	Error      string `json:"error,omitempty"`
}
