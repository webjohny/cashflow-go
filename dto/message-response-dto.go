package dto

type MessageResponseDto struct {
	Message string `json:"message"`
	Result  bool   `json:"result,omitempty"`
}
