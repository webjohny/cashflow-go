package dto

type SendAssetsBodyDTO struct {
	Amount  int    `json:"amount" form:"amount"`
	AssetId string `json:"id" form:"id"`
	Asset   string `json:"asset" form:"asset"`
	Player  int    `json:"player" form:"player"`
}
