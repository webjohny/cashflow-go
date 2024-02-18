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

type CardDoodad struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	Heading       string `json:"heading"`
	Description   string `json:"description"`
	Cost          int    `json:"cost"`
	Rule          string `json:"rule"`
	IsConditional bool   `json:"is_conditional"`
	HasBabies     bool   `json:"has_babies"`
}

type CardRealEstate struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Symbol      string  `json:"symbol"`
	Heading     string  `json:"heading"`
	Description string  `json:"description"`
	Rule        *string `json:"rule"`
	Cost        int     `json:"cost"`
	Mortgage    *int    `json:"mortgage"`
	DownPayment *int    `json:"down_payment"`
	CashFlow    *int    `json:"cash_flow"`
}

type CardStocks struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Symbol      string `json:"symbol"`
	Heading     string `json:"heading"`
	Description string `json:"description"`
	Rule        string `json:"rule"`
	Price       int    `json:"price"`
	Count       int    `json:"count"`
	OnlyYou     bool   `json:"only_you"`
	Range       []int  `json:"range"`
}

type CardPreciousMetals struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Symbol      string `json:"symbol"`
	Heading     string `json:"heading"`
	Description string `json:"description"`
}

type CardDream struct {
	ID          string `json:"id"`
	Heading     string `json:"heading"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Cost        int    `json:"cost"`
}

type CardDice struct {
	Dices      []int    `json:"dices"`
	CashFlow   *int     `json:"cashFlow"`
	CostPerOne *float32 `json:"costPerOne"`
}

type CardRiskBusiness struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Dices       []CardDice `json:"dices"`
	ExtraDices  int        `json:"extraDices"`
	Symbol      string     `json:"symbol"`
	Heading     string     `json:"heading"`
	Description string     `json:"description"`
	Cost        int        `json:"cost"`
}

type CardRiskStocks struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Count       int        `json:"count"`
	Cost        int        `json:"cost"`
	Dices       []CardDice `json:"dices"`
	ExtraDices  int        `json:"extraDices"`
	Symbol      string     `json:"symbol"`
	Heading     string     `json:"heading"`
	Description string     `json:"description"`
	CostPerOne  float64    `json:"costPerOne"`
}

type PlayerIncome struct {
	RealEstates []CardRealEstate `json:"real_estates"`
	Salary      int              `json:"salary"`
}

type PlayerAssets struct {
	Dreams         []CardDream          `json:"dreams"`
	PreciousMetals []CardPreciousMetals `json:"precious_metals"`
	RealEstates    []CardRealEstate     `json:"real_estates"`
	Stocks         []CardStocks         `json:"stocks"`
	Savings        int                  `json:"savings"`
}

type PlayerLiabilities struct {
	RealEstates []CardRealEstate `json:"real_estates"`
}

type Player struct {
	ID              uint64            `gorm:"primary_key:auto_increment" json:"id"`
	RaceId          uint64            `gorm:"index" gorm:"index" json:"race_id"`
	Username        string            `gorm:"uniqueIndex;type:varchar(255)" json:"username"`
	Role            string            `json:"role"`
	Color           string            `json:"color"`
	Income          PlayerIncome      `json:"income"`
	Babies          uint8             `json:"babies"`
	Expenses        map[string]int    `json:"expenses"`
	Assets          PlayerAssets      `json:"assets"`
	Liabilities     PlayerLiabilities `json:"liabilities"`
	Cash            int               `json:"cash"`
	TotalIncome     int               `json:"total_income"`
	TotalExpenses   int               `json:"total_expenses"`
	CashFlow        int               `json:"cash_flow"`
	PassiveIncome   int               `json:"passive_income"`
	Profession      uint8             `json:"profession"`
	LastPosition    uint8             `json:"last_position"`
	CurrentPosition uint8             `json:"current_position"`
	DualDiceCount   uint8             `json:"dual_dice_count"`
	SkippedTurns    uint8             `json:"skipped_turns"`
	CanReRoll       *bool             `json:"can_re_roll"`
	InBigRace       *bool             `json:"in_big_race"`
	HasBankrupt     *bool             `json:"has_bankrupt"`
	AboutToBankrupt *bool             `json:"about_to_bankrupt"`
	HasMlm          *bool             `json:"has_mlm"`
	CreatedAt       string            `json:"created_at"`
}

func (e *Player) FindStocks(symbol string) (int, *CardStocks) {
	for i := 0; i < len(e.Assets.Stocks); i++ {
		if symbol == e.Assets.Stocks[i].Symbol {
			return i, &e.Assets.Stocks[i]
		}
	}

	return -1, nil
}

func (e *Player) FindRealEstate(id string) *CardRealEstate {
	for i := 0; i < len(e.Assets.RealEstates); i++ {
		if id == e.Assets.RealEstates[i].ID {
			return &e.Assets.RealEstates[i]
		}
	}

	return nil
}

func (e *Player) FindPreciousMetals(symbol string) *CardPreciousMetals {
	for i := 0; i < len(e.Assets.PreciousMetals); i++ {
		if symbol == e.Assets.PreciousMetals[i].Symbol {
			return &e.Assets.PreciousMetals[i]
		}
	}

	return nil
}

func (e *Player) CanContinue() bool {
	return e.Cash > 0
}

func (e *Player) IsIncomeStable() bool {
	return e.CalculatePassiveIncome() < e.CalculateTotalExpenses()
}

func (e *Player) IsBankrupt() bool {
	return e.CalculateCashFlow() < 300
}

func (e *Player) CalculateCashFlow() int {
	return e.CalculateTotalIncome() - e.CalculateTotalExpenses()
}

func (e *Player) CalculatePassiveIncome() int {
	passiveIncome := 0

	for i := 0; i < len(e.Income.RealEstates); i++ {
		passiveIncome += *e.Income.RealEstates[i].CashFlow
	}

	return passiveIncome
}

func (e *Player) CalculateTotalIncome() int {
	return e.Income.Salary + e.CalculatePassiveIncome()
}

func (e *Player) CalculateTotalExpenses() int {
	totalExpenses := 0

	for key, expense := range e.Expenses {
		if key == "perChildExpense" {
			expense = int(e.Babies) * expense
		}
		totalExpenses += expense
	}

	return totalExpenses
}
