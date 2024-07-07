package dto

type SendAssetsBodyDTO struct {
	Amount  int    `json:"amount" form:"amount"`
	AssetId string `json:"asset_id" form:"asset_id"`
	Asset   string `json:"asset" form:"asset"`
	Player  string `json:"player" form:"player"`
}
