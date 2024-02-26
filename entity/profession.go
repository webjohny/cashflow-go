package entity

type Profession struct {
	ID          uint64            `json:"id"`
	Profession  string            `json:"profession"`
	Income      PlayerIncome      `json:"income"`
	Babies      int               `json:"babies"`
	Expenses    map[string]int    `json:"expenses"`
	Assets      PlayerAssets      `json:"assets"`
	Liabilities PlayerLiabilities `json:"liabilities"`
}
