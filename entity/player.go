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
	Dreams      []CardDream       `json:"dreams"`
	OtherAssets []CardOtherAssets `json:"other"`
	RealEstates []CardRealEstate  `json:"realEstates"`
	Business    []CardBusiness    `json:"business"`
	Stocks      []CardStocks      `json:"stocks"`
	Savings     int               `json:"savings"`
}

type PlayerLiabilities struct {
	BankLoan       int `json:"bankLoan"`
	HomeMortgage   int `json:"homeMortgage"`
	SchoolLoans    int `json:"schoolLoans"`
	CarLoans       int `json:"carLoans"`
	CreditCardDebt int `json:"creditCardDebt"`
}

type Player struct {
	ID              uint64            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint64            `gorm:"uniqueIndex:idx_username" json:"user_id"`
	RaceID          uint64            `gorm:"index" json:"race_id"`
	Username        string            `gorm:"uniqueIndex:idx_username;type:varchar(255)" json:"username"`
	Role            string            `gorm:"type:varchar(255)" json:"role"`
	Color           string            `gorm:"type:varchar(255)" json:"color"`
	Salary          int               `json:"salary" gorm:"allowzero"`
	Babies          uint8             `json:"babies" gorm:"allowzero"`
	Expenses        map[string]int    `gorm:"type:json;serializer:json" json:"expenses"`
	Assets          PlayerAssets      `gorm:"type:json;serializer:json" json:"assets"`
	Liabilities     PlayerLiabilities `gorm:"type:json;serializer:json" json:"liabilities"`
	Cash            int               `json:"cash" gorm:"allowzero"`
	TotalIncome     int               `json:"total_income" gorm:"allowzero"`
	TotalExpenses   int               `json:"total_expenses" gorm:"allowzero"`
	CashFlow        int               `json:"cash_flow" gorm:"allowzero"`
	PassiveIncome   int               `json:"passive_income" gorm:"allowzero"`
	ProfessionID    uint8             `json:"profession_id"`
	Profession      Profession        `gorm:"-" json:"profession"`
	LastPosition    uint8             `json:"last_position" gorm:"allowzero"`
	CurrentPosition uint8             `json:"current_position" gorm:"allowzero"`
	DualDiceCount   uint8             `json:"dual_dice_count" gorm:"allowzero"`
	SkippedTurns    uint8             `json:"skipped_turns" gorm:"allowzero"`
	IsRolledDice    uint8             `json:"is_rolled_dice"`
	CanReRoll       uint8             `json:"can_re_roll"`
	OnBigRace       uint8             `json:"on_big_race"`
	HasBankrupt     uint8             `json:"has_bankrupt"`
	AboutToBankrupt string            `gorm:"type:varchar(255)" json:"about_to_bankrupt"`
	HasMlm          uint8             `json:"has_mlm"`
	CreatedAt       datatypes.Date    `gorm:"column:created_at;type:datetime;default:current_timestamp;not null" json:"created_at"`
}

func (e *Player) FindStocksBySymbol(symbol string) (int, *CardStocks) {
	for i := 0; i < len(e.Assets.Stocks); i++ {
		if symbol == e.Assets.Stocks[i].Symbol {
			return i, &e.Assets.Stocks[i]
		}
	}

	return -1, &CardStocks{}
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

func (e *Player) Reset(profession Profession) {
	e.Assets.Business = make([]CardBusiness, 0)
	e.Assets.RealEstates = make([]CardRealEstate, 0)
	e.Assets.OtherAssets = make([]CardOtherAssets, 0)
	e.Assets.Stocks = make([]CardStocks, 0)
	e.Assets.Dreams = make([]CardDream, 0)

	if profession.ID != 0 {
		e.Salary = profession.Income.Salary
		e.Babies = uint8(profession.Babies)
		e.Expenses = profession.Expenses
		e.Assets = profession.Assets
		e.Liabilities = profession.Liabilities
		e.ProfessionID = uint8(profession.ID)
	} else {
		e.Assets.Savings = 0
	}

	e.CurrentPosition = 0
	e.SkippedTurns = 0
	e.SkippedTurns = 0
}

func (e *Player) CreateResponse() RaceResponse {
	return RaceResponse{
		ID:        e.ID,
		UserId:    e.UserID,
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

func (e *Player) FindBusinessBySymbol(symbol string) (int, *CardBusiness) {
	for i := 0; i < len(e.Assets.Business); i++ {
		if symbol == e.Assets.Business[i].Symbol {
			return i, &e.Assets.Business[i]
		}
	}

	return 0, &CardBusiness{}
}

func (e *Player) FindBusinessByID(ID string) (int, *CardBusiness) {
	for i := 0; i < len(e.Assets.Business); i++ {
		if ID == e.Assets.Business[i].ID {
			return i, &e.Assets.Business[i]
		}
	}

	return 0, &CardBusiness{}
}

func (e *Player) FindAllBusinessBySymbol(symbol string) []CardBusiness {
	items := make([]CardBusiness, 0)

	for i := 0; i < len(e.Assets.Business); i++ {
		if symbol == e.Assets.Business[i].Symbol {
			items = append(items, e.Assets.Business[i])
		}
	}

	return items
}

func (e *Player) ReduceLimitedShares(ID string, count int) {
	_, shares := e.FindBusinessByID(ID)

	for i := 0; i < len(shares.History); i++ {
		if shares.History[i].Count > count {
			shares.History[i].Count -= count
			shares.History[i].SumCost()

			break
		} else {
			count -= shares.History[i].Count
			shares.History = append(shares.History[:i], shares.History[i+1:]...)
			i--
		}
	}
}

func (e *Player) AddBusiness(card CardBusiness) {
	e.Assets.Business = append(e.Assets.Business, card)
}

func (e *Player) AddStocks(card CardStocks) {
	_, asset := e.FindStocksBySymbol(card.Symbol)

	if asset.ID != "" {
		asset.Count = card.Count
		asset.Price = card.Price
	} else {
		e.Assets.Stocks = append(e.Assets.Stocks, card)
	}
}

func (e *Player) FindRealEstateByID(id string) *CardRealEstate {
	for i := 0; i < len(e.Assets.RealEstates); i++ {
		if id == e.Assets.RealEstates[i].ID {
			return &e.Assets.RealEstates[i]
		}
	}

	return &CardRealEstate{}
}

func (e *Player) FindRealEstateBySymbol(symbol string) (int, *CardRealEstate) {
	for i := 0; i < len(e.Assets.RealEstates); i++ {
		if symbol == e.Assets.RealEstates[i].Symbol {
			return i, &e.Assets.RealEstates[i]
		}
	}

	return 0, &CardRealEstate{}
}

func (e *Player) FindAllRealEstateBySymbol(symbol string) []CardRealEstate {
	var realEstates []CardRealEstate

	for i := 0; i < len(e.Assets.RealEstates); i++ {
		if symbol == e.Assets.RealEstates[i].Symbol {
			realEstates = append(realEstates, e.Assets.RealEstates[i])
		}
	}

	return realEstates
}

func (e *Player) FindOtherAssetsBySymbol(symbol string) (int, *CardOtherAssets) {
	for i := 0; i < len(e.Assets.OtherAssets); i++ {
		if symbol == e.Assets.OtherAssets[i].Symbol {
			return i, &e.Assets.OtherAssets[i]
		}
	}

	return -1, &CardOtherAssets{}
}

func (e *Player) FindOtherAssetsByID(ID string) (int, *CardOtherAssets) {
	for i := 0; i < len(e.Assets.OtherAssets); i++ {
		if ID == e.Assets.OtherAssets[i].ID {
			return i, &e.Assets.OtherAssets[i]
		}
	}

	return -1, &CardOtherAssets{}
}

func (e *Player) UpdateAsset(symbol string, card CardOtherAssets) {
	index, _ := e.FindOtherAssetsBySymbol(symbol)

	e.Assets.OtherAssets[index] = card
}

func (e *Player) RemoveOtherAssets(symbol string) {
	index, _ := e.FindOtherAssetsBySymbol(symbol)
	if index >= 0 && index < len(e.Assets.OtherAssets) {
		e.Assets.OtherAssets = append(e.Assets.OtherAssets[:index], e.Assets.OtherAssets[index+1:]...)
	}
}

func (e *Player) RemoveOtherAssetsByID(ID string) {
	index, _ := e.FindOtherAssetsByID(ID)
	if index >= 0 && index < len(e.Assets.OtherAssets) {
		e.Assets.OtherAssets = append(e.Assets.OtherAssets[:index], e.Assets.OtherAssets[index+1:]...)
	}
}

func (h *CardOtherAssets) SumCost() {
	h.Cost = h.CostPerOne * h.Count
}

func (e *Player) RemoveStocks(symbol string) {
	index, _ := e.FindStocksBySymbol(symbol)
	if index >= 0 && index < len(e.Assets.Stocks) {
		e.Assets.Stocks = append(e.Assets.Stocks[:index], e.Assets.Stocks[index+1:]...)
	}
}

func (e *Player) ReduceStocks(symbol string, count int) {
	_, stocks := e.FindStocksBySymbol(symbol)

	for i := 0; i < len(stocks.History); i++ {
		if stocks.History[i].Count > count {
			stocks.History[i].Count -= count
			stocks.History[i].SumCost()

			break
		} else {
			count -= stocks.History[i].Count
			stocks.History = append(stocks.History[:i], stocks.History[i+1:]...)
			i--
		}
	}
}

func (e *Player) RemoveRealEstate(id string) CardRealEstate {
	for i := 0; i < len(e.Assets.RealEstates); i++ {
		if id == e.Assets.RealEstates[i].ID {
			e.Assets.RealEstates = append(e.Assets.RealEstates[:i], e.Assets.RealEstates[i+1:]...)
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

	return CardBusiness{}
}

func (e *Player) SplitStocks(card string) {
	_, stock := e.FindStocksBySymbol(card)
	stock.Count *= 2
	e.DeactivateReRoll()
}

func (e *Player) ReverseSplitStocks(card string) {
	_, stock := e.FindStocksBySymbol(card)
	stock.Count = int(math.Ceil(float64(stock.Count) / 2))
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

	for i := 0; i < len(e.Assets.RealEstates); i++ {
		passiveIncome += e.Assets.RealEstates[i].CashFlow
	}

	for i := 0; i < len(e.Assets.Business); i++ {
		if e.Assets.Business[i].Count > 0 {
			passiveIncome += e.Assets.Business[i].Count * e.Assets.Business[i].CashFlow
		} else {
			passiveIncome += e.Assets.Business[i].CashFlow
		}
	}

	return passiveIncome
}

func (e *Player) CalculateTotalIncome() int {
	return e.Salary + e.CalculatePassiveIncome()
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
