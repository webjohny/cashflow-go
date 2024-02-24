package service

import (
	"fmt"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/storage"
	"math"
)

type PlayerService interface {
	GetPlayerByUsername(username string) *entity.Player
	Payday(player entity.Player)
	CashFlowDay(player entity.Player)
	Doodad(card entity.CardDoodad, player entity.Player) error
	BuyBusiness(card entity.CardBusiness, player entity.Player) error
	BuyRealEstate(card entity.CardRealEstate, player entity.Player) error
	BuyRiskBusiness(card entity.CardRiskBusiness, player entity.Player, rolledDice int) (error, bool)
	BuyRiskStocks(card entity.CardRiskStocks, player entity.Player, rolledDice int) (error, bool)
	BuyDream(card entity.CardDream, player entity.Player) error
	BuyStocks(card entity.CardStocks, player entity.Player, count int, updateCash bool) error
	SellGold(card entity.CardPreciousMetals, player entity.Player, count int) error
	SellStocks(card entity.CardStocks, player entity.Player, count int, updateCash bool) error
	SellRealEstate(card entity.CardRealEstate, player entity.Player) error
	DecreaseStocks(card entity.CardStocks, player entity.Player) error
	IncreaseStocks(card entity.CardStocks, player entity.Player) error
	Charity(player entity.Player) error
	BigCharity(card entity.CardCharity, player entity.Player) error
	PayTax(card entity.CardPayTax, player entity.Player) error
	Downsized(player entity.Player) error
	MoveToBigRace(player entity.Player) error
	PayDamages(card entity.CardMarket, player entity.Player) error
	AddGoldCoins(card entity.CardPreciousMetals, player entity.Player) error
	SellRealEstates(player entity.Player) (error, int)
	SellBusiness(player entity.Player) (error, int)
	TakeLoan(player entity.Player, amount int) error
	PayLoan(player entity.Player, actionType string, amount int) error
	UpdateCash(player *entity.Player, amount int, details string)
}

type playerService struct {
	playerRepository   repository.PlayerRepository
	transactionService TransactionService
}

func NewPlayerService(playerRepo repository.PlayerRepository, transactionService TransactionService) PlayerService {
	return &playerService{
		playerRepository:   playerRepo,
		transactionService: transactionService,
	}
}

func (service *playerService) GetPlayerByUsername(username string) *entity.Player {
	return service.playerRepository.FindPlayerByUsername(username)
}

func (service *playerService) Payday(player entity.Player) {
	service.UpdateCash(&player, player.CalculateCashFlow(), "Зарплата")
}

func (service *playerService) CashFlowDay(player entity.Player) {
	service.UpdateCash(&player, player.CalculateCashFlow(), "Кэш-флоу день")
}

func (service *playerService) Doodad(card entity.CardDoodad, player entity.Player) error {
	cost := card.Cost

	if card.HasBabies && player.Babies <= 0 {
		return fmt.Errorf(storage.ErrorYouHaveNoBabies)
	}

	if player.Cash < cost {
		return fmt.Errorf(storage.ErrorNotEnoughMoney)
	}

	service.UpdateCash(&player, -cost, "Растраты")

	return nil
}

func (service *playerService) BuyDream(card entity.CardDream, player entity.Player) error {
	cost := card.Cost

	if player.Cash < cost {
		return fmt.Errorf(storage.ErrorNotEnoughMoney)
	}

	player.Assets.Dreams = append(player.Assets.Dreams, card)

	service.UpdateCash(&player, -cost, "Мечта")

	return nil
}

func (service *playerService) BuyStocks(card entity.CardStocks, player entity.Player, count int, updateCash bool) error {
	totalCost := int(float64(card.Price) * float64(count))

	if player.Cash < totalCost {
		return fmt.Errorf(storage.ErrorNotEnoughMoney)
	}

	key, stock := player.FindStocks(card.Symbol)

	if stock != nil {
		totalCount := count + *stock.Count
		*stock.Count = totalCount
		player.Assets.Stocks[key] = *stock
	} else {
		player.Assets.Stocks = append(player.Assets.Stocks, card)
	}

	if updateCash {
		service.UpdateCash(&player, -totalCost, card.Symbol)
	}

	return nil
}

func (service *playerService) SellGold(card entity.CardPreciousMetals, player entity.Player, count int) error {
	_, gold := player.FindPreciousMetals(card.Symbol)
	totalCost := card.Cost * count

	if gold.Count < count {
		return fmt.Errorf(storage.ErrorNotEnoughMoney)
	}

	gold.Count -= count

	service.UpdateCash(&player, totalCost, card.Symbol)

	if gold.Count <= 0 {
		player.RemovePreciousMetals(gold.Symbol)
	}

	return nil
}

func (service *playerService) SellStocks(card entity.CardStocks, player entity.Player, count int, updateCash bool) error {
	_, stock := player.FindStocks(card.Symbol)

	if stock != nil || *stock.Count < count {
		return fmt.Errorf(storage.ErrorNotFoundStocks)
	}

	totalCost := card.Price * count
	*stock.Count -= count

	if updateCash {
		service.UpdateCash(&player, totalCost, card.Symbol)
	}

	if *stock.Count <= 0 {
		player.RemoveStocks(stock.Symbol)
	}

	return nil
}

func (service *playerService) SellRealEstate(card entity.CardRealEstate, player entity.Player) error {
	realEstate := player.FindRealEstate(card.ID)

	if realEstate == nil {
		return fmt.Errorf(storage.ErrorNotFoundAssets)
	}

	value := (realEstate.Cost / 100) * card.Value
	totalCost := realEstate.Cost + value

	if card.Plus {
		totalCost = realEstate.Cost + card.Value
	}

	service.UpdateCash(&player, totalCost-*realEstate.Mortgage, card.Symbol)

	player.RemoveRealEstate(card.ID)

	return nil
}

func (service *playerService) DecreaseStocks(card entity.CardStocks, player entity.Player) error {
	key, stock := player.FindStocks(card.Symbol)

	if stock != nil {
		*stock.Count = int(math.Floor(float64(*stock.Count / *card.Decrease)))
		player.Assets.Stocks[key] = *stock
	} else {
		player.Assets.Stocks = append(player.Assets.Stocks, card)
	}

	return nil
}

func (service *playerService) IncreaseStocks(card entity.CardStocks, player entity.Player) error {
	key, stock := player.FindStocks(card.Symbol)

	if stock != nil {
		*stock.Count = int(math.Floor(float64(*stock.Count * *card.Increase)))
		player.Assets.Stocks[key] = *stock
	} else {
		player.Assets.Stocks = append(player.Assets.Stocks, card)
	}

	return nil
}

func (service *playerService) Charity(player entity.Player) error {
	amount := int(math.Floor(0.1 * float64(player.CalculateTotalIncome())))

	if player.Cash < amount {
		return fmt.Errorf(storage.ErrorNotEnoughMoney)
	}

	service.UpdateCash(&player, -amount, "Благотворительность")

	return nil
}

func (service *playerService) BigCharity(card entity.CardCharity, player entity.Player) error {
	if player.Cash < card.Cost {
		return fmt.Errorf(storage.ErrorNotEnoughMoney)
	}

	service.UpdateCash(&player, -card.Cost, "Акция милосердия")

	return nil
}

func (service *playerService) PayTax(card entity.CardPayTax, player entity.Player) error {
	amount := (player.Cash / 100) * card.Percent

	if player.Cash < amount {
		return fmt.Errorf(storage.ErrorNotEnoughMoney)
	}

	service.UpdateCash(&player, -amount, "Налоги")

	return nil
}

func (service *playerService) Downsized(player entity.Player) error {
	amount := player.CalculateTotalExpenses()

	if player.Cash < amount {
		return fmt.Errorf(storage.ErrorNotEnoughMoney)
	}

	service.UpdateCash(&player, -amount, "Уволен")

	return nil
}

func (service *playerService) MoveToBigRace(player entity.Player) error {
	if !player.ConditionsForBigRace() {
		return fmt.Errorf(storage.ErrorMovingBigRaceDeclined)
	}

	cashFlow := player.CalculatePassiveIncome() * 100

	player.OnBigRace = true
	player.CashFlow = cashFlow
	player.Cash = cashFlow + player.Cash
	player.TotalIncome = 0
	player.TotalExpenses = 0
	player.Expenses = make(map[string]int)
	player.Assets = entity.PlayerAssets{
		Savings:        0,
		Stocks:         make([]entity.CardStocks, 0),
		PreciousMetals: make([]entity.CardPreciousMetals, 0),
		RealEstates:    make([]entity.CardRealEstate, 0),
		Business:       make([]entity.CardBusiness, 0),
		Dreams:         make([]entity.CardDream, 0),
	}
	player.Income = entity.PlayerIncome{
		RealEstates: []entity.CardRealEstate{
			{
				CashFlow: &cashFlow,
			},
		},
	}

	return nil
}

func (service *playerService) PayDamages(card entity.CardMarket, player entity.Player) error {
	if player.Cash < *card.Cost {
		return fmt.Errorf(storage.ErrorNotEnoughMoney)
	}

	if !player.HasRealEstates() {
		return fmt.Errorf(storage.ErrorYouHaveNoProperties)
	}

	service.UpdateCash(&player, -*card.Cost, "Имущество поврежденно")

	return nil
}

func (service *playerService) AddGoldCoins(card entity.CardPreciousMetals, player entity.Player) error {
	if player.Cash < card.Cost {
		return fmt.Errorf(storage.ErrorNotEnoughMoney)
	}

	player.Assets.PreciousMetals = append(player.Assets.PreciousMetals, card)

	service.UpdateCash(&player, -card.Cost, "Золотые монеты")

	return nil
}

func (service *playerService) SellRealEstates(player entity.Player) (error, int) {
	var totalCash int

	if !player.HasRealEstates() {
		return fmt.Errorf(storage.ErrorYouHaveNoProperties), 0
	}

	for i := 0; i < len(player.Assets.RealEstates); i++ {
		property := player.Assets.RealEstates[i]
		totalCash += *property.DownPayment / 2
	}

	player.Assets.RealEstates = make([]entity.CardRealEstate, 0)
	player.Income.RealEstates = make([]entity.CardRealEstate, 0)
	player.Liabilities.RealEstates = make([]entity.CardRealEstate, 0)

	return nil, totalCash
}

func (service *playerService) SellBusiness(player entity.Player) (error, int) {
	var totalCash int

	if !player.HasBusiness() {
		return fmt.Errorf(storage.ErrorYouHaveNoProperties), 0
	}

	for i := 0; i < len(player.Assets.Business); i++ {
		property := player.Assets.Business[i]
		totalCash += property.Cost / 2
	}

	player.Assets.Business = make([]entity.CardBusiness, 0)
	player.Income.Business = make([]entity.CardBusiness, 0)
	player.Liabilities.Business = make([]entity.CardBusiness, 0)

	return nil, totalCash
}

func (service *playerService) TakeLoan(player entity.Player, amount int) error {
	service.UpdateCash(&player, amount, "Взял в кредит")

	player.Liabilities.BankLoan += amount
	player.Expenses["bankLoanPayment"] = player.Liabilities.BankLoan / 10

	return nil
}

func (service *playerService) PayLoan(player entity.Player, actionType string, amount int) error {
	if player.Cash < amount {
		return fmt.Errorf(storage.ErrorNotEnoughMoney)
	}

	loanMapper := map[string]string{
		"homeMortgage":   "homeMortgagePayment",
		"schoolLoans":    "schoolLoanPayment",
		"carLoans":       "carLoanPayment",
		"creditCardDebt": "creditCardPayment",
	}

	service.UpdateCash(&player, -amount, "Оплата по кредиту")

	var liabilityAmount int

	if actionType == "homeMortgage" {
		liabilityAmount = player.Liabilities.HomeMortgage
	} else if actionType == "homeMortgage" {
		liabilityAmount = player.Liabilities.SchoolLoans
	} else if actionType == "homeMortgage" {
		liabilityAmount = player.Liabilities.CarLoans
	} else if actionType == "homeMortgage" {
		liabilityAmount = player.Liabilities.CreditCardDebt
	}

	if liabilityAmount > 0 {
		liabilityAmount -= amount
	}

	if actionType == "bankLoan" {
		player.Expenses["bankLoanPayment"] = liabilityAmount / 10
	} else {
		player.Expenses[loanMapper[actionType]] = 0
	}

	return nil
}

func (service *playerService) BuyRealEstate(card entity.CardRealEstate, player entity.Player) error {
	if player.Cash < *card.DownPayment {
		return fmt.Errorf(storage.ErrorNotEnoughMoney)
	}

	player.Assets.RealEstates = append(player.Assets.RealEstates, card)
	player.Income.RealEstates = append(player.Income.RealEstates, card)
	player.Liabilities.RealEstates = append(player.Liabilities.RealEstates, card)

	service.UpdateCash(&player, -*card.DownPayment, card.Heading)

	return nil
}

func (service *playerService) BuyBusiness(card entity.CardBusiness, player entity.Player) error {
	if player.Cash < card.Cost {
		return fmt.Errorf(storage.ErrorNotEnoughMoney)
	}

	player.Assets.Business = append(player.Assets.Business, card)
	player.Income.Business = append(player.Income.Business, card)
	player.Liabilities.Business = append(player.Liabilities.Business, card)

	service.UpdateCash(&player, -card.Cost, card.Heading)

	return nil
}

func (service *playerService) BuyRiskBusiness(card entity.CardRiskBusiness, player entity.Player, rolledDice int) (error, bool) {
	cost := card.Cost

	if player.Cash < cost {
		return fmt.Errorf(storage.ErrorNotEnoughMoney), false
	}

	var cashFlow int
	for _, dice := range card.Dices {
		for _, value := range dice.Dices {
			if value == rolledDice {
				cashFlow = *dice.CashFlow
			}
		}
	}

	if cashFlow > 0 {
		service.UpdateCash(&player, -cost, card.Heading)

		business := entity.CardBusiness{
			ID:          card.ID,
			Type:        card.Type,
			Symbol:      card.Symbol,
			Heading:     card.Heading,
			Description: card.Description,
			Cost:        card.Cost,
			CashFlow:    &cashFlow,
		}

		player.Assets.Business = append(player.Assets.Business, business)
		player.Income.Business = append(player.Income.Business, business)
		player.Liabilities.Business = append(player.Liabilities.Business, business)

		return nil, true
	}

	return nil, false
}

func (service *playerService) BuyRiskStocks(card entity.CardRiskStocks, player entity.Player, rolledDice int) (error, bool) {
	cost := card.Cost

	if player.Cash < cost {
		return fmt.Errorf(storage.ErrorNotEnoughMoney), false
	}

	var costPerOne float32
	for _, dice := range card.Dices {
		for _, value := range dice.Dices {
			if value == rolledDice {
				costPerOne = *dice.CostPerOne
			}
		}
	}

	if costPerOne > 0 {
		service.UpdateCash(&player, -cost, card.Heading)
		service.UpdateCash(&player, int(float32(card.Count)*costPerOne), card.Heading)

		return nil, true
	}

	return nil, false
}

func (service *playerService) UpdateCash(player *entity.Player, amount int, details string) {
	currentCash := player.Cash

	player.Cash += amount

	go service.SetTransaction(player.ID, currentCash, player.Cash, amount, details)
}

func (service *playerService) SetTransaction(ID uint64, currentCash int, cash int, amount int, details string) {
	service.transactionService.InsertPlayerTransaction(dto.TransactionCreatePlayerDTO{
		PlayerID:    ID,
		Details:     details,
		CurrentCash: currentCash,
		Cash:        cash,
		Amount:      amount,
	})
}
