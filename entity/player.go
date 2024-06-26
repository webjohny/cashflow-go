package entity

import (
	"gorm.io/datatypes"
	"math"
)

var PlayerRoles = struct {
	Guest    string
	WaitList string
	Owner    string
	Admin    string
}{
	Guest:    "guest",
	WaitList: "wait_list",
	Owner:    "owner",
	Admin:    "admin",
}

type PlayerIncome struct {
	RealEstates []CardRealEstate `json:"realEstates"`
	Business    []CardBusiness   `json:"business"`
	Salary      int              `json:"salary"`
}

type PlayerAssets struct {
	Dreams         []CardDream          `json:"dreams"`
	PreciousMetals []CardPreciousMetals `json:"preciousMetals"`
	RealEstates    []CardRealEstate     `json:"realEstates"`
	Business       []CardBusiness       `json:"business"`
	Stocks         []CardStocks         `json:"stocks"`
	Savings        int                  `json:"savings"`
}

type PlayerLiabilities struct {
	RealEstates    []CardRealEstate `json:"realEstates"`
	Business       []CardBusiness   `json:"business"`
	BankLoan       int              `json:"bankLoan"`
	HomeMortgage   int              `json:"homeMortgage"`
	SchoolLoans    int              `json:"schoolLoans"`
	CarLoans       int              `json:"carLoans"`
	CreditCardDebt int              `json:"creditCardDebt"`
}

type Player struct {
	ID              uint64            `gorm:"primary_key:auto_increment" json:"id"`
	UserId          uint64            `gorm:"index" json:"user_id"`
	RaceId          uint64            `gorm:"index" gorm:"index" json:"race_id"`
	Username        string            `gorm:"uniqueIndex;type:varchar(255)" json:"username"`
	Role            string            `json:"role"`
	Color           string            `json:"color"`
	Income          PlayerIncome      `gorm:"serializer:json" json:"income"`
	Babies          uint8             `json:"babies"`
	Expenses        map[string]int    `gorm:"serializer:json" json:"expenses"`
	Assets          PlayerAssets      `gorm:"serializer:json" json:"assets"`
	Liabilities     PlayerLiabilities `gorm:"serializer:json" json:"liabilities"`
	Cash            int               `json:"cash"`
	TotalIncome     int               `json:"total_income"`
	TotalExpenses   int               `json:"total_expenses"`
	CashFlow        int               `json:"cash_flow"`
	PassiveIncome   int               `json:"passive_income"`
	ProfessionId    uint8             `json:"profession_id"`
	Profession      Profession        `gorm:"-" sql:"-" json:"profession"`
	LastPosition    uint8             `json:"last_position"`
	CurrentPosition uint8             `json:"current_position"`
	DualDiceCount   uint8             `json:"dual_dice_count"`
	SkippedTurns    uint8             `json:"skipped_turns"`
	IsRolledDice    uint8             `json:"is_rolled_dice"`
	CanReRoll       uint8             `json:"can_re_roll"`
	OnBigRace       uint8             `json:"on_big_race"`
	HasBankrupt     uint8             `json:"has_bankrupt"`
	AboutToBankrupt string            `json:"about_to_bankrupt"`
	HasMlm          uint8             `json:"has_mlm"`
	CreatedAt       datatypes.Date    `json:"created_at"`
}

func (e *Player) FindStocks(symbol string) (int, CardStocks) {
	for i := 0; i < len(e.Assets.Stocks); i++ {
		if symbol == e.Assets.Stocks[i].Symbol {
			return i, e.Assets.Stocks[i]
		}
	}

	return -1, CardStocks{}
}

func (e *Player) BornBaby() {
	e.Babies++
}

func (e *Player) ChangeDiceStatus(status bool) {
	if status {
		e.IsRolledDice = 1
	} else {
		e.IsRolledDice = 0
	}
}

func (e *Player) Move(steps int) {
	e.LastPosition = e.CurrentPosition

	if e.OnBigRace == 1 {
		e.CurrentPosition = uint8((int(e.CurrentPosition) + steps) % 46)

		if e.CurrentPosition == 0 {
			e.CurrentPosition = 46
		}
	} else {
		e.CurrentPosition = uint8((int(e.CurrentPosition) + steps) % 24)

		if e.CurrentPosition == 0 {
			e.CurrentPosition = 24
		}
	}
}

func (e *Player) IncrementDualDiceCount() {
	e.DualDiceCount += 3
}

func (e *Player) CreateResponse() RaceResponse {
	return RaceResponse{
		ID:        e.ID,
		UserId:    e.UserId,
		Username:  e.Username,
		Responded: false,
	}
}

func (e *Player) DecrementDualDiceCount() {
	e.DualDiceCount--
}

func (e *Player) AllowReRoll() {
	e.ChangeDiceStatus(false)
	e.CanReRoll = 1
}

func (e *Player) DeactivateReRoll() {
	e.ChangeDiceStatus(true)
	e.CanReRoll = 0
}

func (e *Player) InitializeSkippedTurns() {
	e.SkippedTurns = 2
}

func (e *Player) DecrementSkippedTurns() {
	e.SkippedTurns--
}

func (e *Player) HasRealEstates() bool {
	return len(e.Assets.RealEstates) > 0
}

func (e *Player) HasBusiness() bool {
	return len(e.Assets.Business) > 0
}

func (e *Player) FindBusiness(id string) CardBusiness {
	for i := 0; i < len(e.Assets.Business); i++ {
		if id == e.Assets.Business[i].ID {
			return e.Assets.Business[i]
		}
	}

	return CardBusiness{}
}

func (e *Player) FindRealEstate(id string) CardRealEstate {
	for i := 0; i < len(e.Assets.RealEstates); i++ {
		if id == e.Assets.RealEstates[i].ID {
			return e.Assets.RealEstates[i]
		}
	}

	return CardRealEstate{}
}

func (e *Player) FindPreciousMetals(symbol string) (int, CardPreciousMetals) {
	for i := 0; i < len(e.Assets.PreciousMetals); i++ {
		if symbol == e.Assets.PreciousMetals[i].Symbol {
			return i, e.Assets.PreciousMetals[i]
		}
	}

	return -1, CardPreciousMetals{}
}

func (e *Player) RemovePreciousMetals(symbol string) CardPreciousMetals {
	index, _ := e.FindPreciousMetals(symbol)
	if index >= 0 && index < len(e.Assets.PreciousMetals) {
		e.Assets.PreciousMetals = append(e.Assets.PreciousMetals[:index], e.Assets.PreciousMetals[index+1:]...)
	}

	return CardPreciousMetals{}
}

func (e *Player) RemoveStocks(symbol string) CardPreciousMetals {
	index, _ := e.FindStocks(symbol)
	if index >= 0 && index < len(e.Assets.Stocks) {
		e.Assets.Stocks = append(e.Assets.Stocks[:index], e.Assets.Stocks[index+1:]...)
	}

	return CardPreciousMetals{}
}

func (e *Player) RemoveRealEstate(id string) CardRealEstate {
	for i := 0; i < len(e.Assets.RealEstates); i++ {
		if id == e.Assets.RealEstates[i].ID {
			e.Assets.RealEstates = append(e.Assets.RealEstates[:i], e.Assets.RealEstates[i+1:]...)
		}
	}
	for i := 0; i < len(e.Liabilities.RealEstates); i++ {
		if id == e.Liabilities.RealEstates[i].ID {
			e.Liabilities.RealEstates = append(e.Liabilities.RealEstates[:i], e.Liabilities.RealEstates[i+1:]...)
		}
	}
	for i := 0; i < len(e.Income.RealEstates); i++ {
		if id == e.Income.RealEstates[i].ID {
			e.Income.RealEstates = append(e.Income.RealEstates[:i], e.Income.RealEstates[i+1:]...)
		}
	}

	return CardRealEstate{}
}

func (e *Player) RemoveBusiness(id string) CardBusiness {
	for i := 0; i < len(e.Assets.Business); i++ {
		if id == e.Assets.Business[i].ID {
			e.Assets.Business = append(e.Assets.Business[:i], e.Assets.Business[i+1:]...)
		}
	}
	for i := 0; i < len(e.Liabilities.Business); i++ {
		if id == e.Liabilities.Business[i].ID {
			e.Liabilities.Business = append(e.Liabilities.Business[:i], e.Liabilities.Business[i+1:]...)
		}
	}
	for i := 0; i < len(e.Income.Business); i++ {
		if id == e.Income.Business[i].ID {
			e.Income.Business = append(e.Income.Business[:i], e.Income.Business[i+1:]...)
		}
	}

	return CardBusiness{}
}

func (e *Player) SplitStocks(card string) {
	_, stock := e.FindStocks(card)
	*stock.Count *= 2
	e.DeactivateReRoll()
}

func (e *Player) ReverseSplitStocks(card string) {
	_, stock := e.FindStocks(card)
	*stock.Count = int(math.Ceil(float64(*stock.Count) / 2))
	e.DeactivateReRoll()
}

func (e *Player) CanContinue() bool {
	return e.Cash > 0
}

func (e *Player) ConditionsForBigRace() bool {
	return e.IsIncomeStable()
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
