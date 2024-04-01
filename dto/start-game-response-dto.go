package dto

type StartGameResponseDto struct {
	ID       uint64 `json:"id"`
	Redirect string `json:"redirect"`
}
