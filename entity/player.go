package entity

var PlayerRoles = struct {
	GUEST string
	OWNER string
	ADMIN string
}{
	GUEST: "guest",
	OWNER: "owner",
	ADMIN: "admin",
}

type Player struct {
	ID              uint64 `gorm:"primary_key:auto_increment" json:"id"`
	GameId          uint64 `gorm:"index" json:"game_id"`
	Username        string `gorm:"uniqueIndex;type:varchar(255)" json:"username"`
	Role            string `json:"role"`
	Color           string `json:"color"`
	Income          string `json:"income"`
	Babies          uint8  `json:"babies"`
	Expenses        string `json:"expenses"`
	Assets          string `json:"assets"`
	Liabilities     string `json:"liabilities"`
	Cash            uint32 `json:"cash"`
	TotalIncome     uint32 `json:"total_income"`
	TotalExpenses   uint32 `json:"total_expenses"`
	CashFlow        uint32 `json:"cash_flow"`
	PassiveIncome   uint32 `json:"passive_income"`
	Profession      uint8  `json:"profession"`
	LastPosition    uint8  `json:"last_position"`
	Transactions    string `json:"transactions"`
	CurrentPosition uint8  `json:"current_position"`
	DualDiceCount   uint8  `json:"dual_dice_count"`
	SkippedTurns    uint8  `json:"skipped_turns"`
	CanReRoll       *bool  `json:"can_re_roll"`
	InBigRace       *bool  `json:"in_big_race"`
	HasBankrupt     *bool  `json:"has_bankrupt"`
	AboutToBankrupt *bool  `json:"about_to_bankrupt"`
	HasMlm          *bool  `json:"has_mlm"`
	CreatedAt       uint16 `json:"created_at"`
}

//gameId:"01HMT107G8J8PAECCDRRCRSNQA"
//isAdmin:false
//username:"webtoolteam"
//role:"guest"
//color:"violet"
//income:"{"salary":2500,"realEstates":[]}"
//babies:0
//expenses:"{"taxes":500,"homeMortgagePayment":400,"schoolLoanPayment":0,"carLoanPayment":100,"creditCardPayment":100,"otherExpenses":600,"bankLoanPayment":0,"perChildExpense":200}"
//assets:"{"savings":800,"preciousMetals":[],"stocks":[],"realEstates":[]}"
//liabilities:"{"homeMortgage":38000,"schoolLoans":0,"carLoans":4000,"creditCardDebt":3000,"bankLoan":0,"realEstates":[]}"
//cash:1600
//totalIncome:2500
//totalExpenses:1700
//cashFlow:800
//passiveIncome:0
//profession:6
//lastPosition:14
//transactions:"[]"
//currentPosition:19
//dualDiceCount:0
//skippedTurns:0
//canReRoll:false
//isInFastTrack:false
//hasBankrupt:false
//aboutToBankrupt:false
//hasMlm:false
//createdAt:1705975496233
