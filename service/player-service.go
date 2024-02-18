package service

import (
	"fmt"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/repository"
)

type PlayerService interface {
	GetPlayerByUsername(username string) *entity.Player
	Payday(player entity.Player)
	Doodad(card entity.CardDoodad, player entity.Player) error
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

func (service *playerService) Doodad(card entity.CardDoodad, player entity.Player) error {
	cost := card.Cost

	if card.HasBabies && player.Babies <= 0 {
		return fmt.Errorf(helper.GetMessage("YOU_HAVE_NO_BABIES"))
	}

	if player.Cash < cost {
		return fmt.Errorf(helper.GetMessage("ERROR_NOT_ENOUGH_MONEY"))
	}

	service.UpdateCash(&player, -cost, "Растраты")

	return nil
}

func (service *playerService) Dream(card entity.CardDream, player entity.Player) error {
	cost := card.Cost

	if player.Cash < cost {
		return fmt.Errorf(helper.GetMessage("ERROR_NOT_ENOUGH_MONEY"))
	}

	player.Assets.Dreams = append(player.Assets.Dreams, card)

	service.UpdateCash(&player, -cost, "Мечта")

	return nil
}

func (service *playerService) BuyStocks(card entity.CardStocks, player entity.Player, count int, updateCash bool) error {
	totalCost := int(float64(card.Price) * float64(count))

	if player.Cash < totalCost {
		return fmt.Errorf(helper.GetMessage("ERROR_NOT_ENOUGH_MONEY"))
	}

	key, stock := player.FindStocks(card.Symbol)

	if stock != nil {
		totalCount := count + stock.Count
		stock.Count = totalCount
	} else {
		player.Assets.Stocks = append(player.Assets.Stocks, card)
	}

	if updateCash {
		service.UpdateCash(&player, -totalCost, card.Symbol)
	}

	return nil
}

func (service *playerService) DivideStocks(card entity.CardStocks, player entity.Player, count int) error {
	stock := player.FindStocks(card.Symbol)

	if stock != nil {
		stock.Count = stock.Count * count
	} else {
		player.Assets.Stocks = append(player.Assets.Stocks, card)
	}

	return nil
}

func (service *playerService) BuyRealEstate(card entity.CardRealEstate, player entity.Player) error {
	if player.Cash < *card.DownPayment {
		return fmt.Errorf(helper.GetMessage("ERROR_NOT_ENOUGH_MONEY"))
	}

	player.Assets.RealEstates = append(player.Assets.RealEstates, card)
	player.Income.RealEstates = append(player.Income.RealEstates, card)
	player.Liabilities.RealEstates = append(player.Liabilities.RealEstates, card)

	service.UpdateCash(&player, -*card.DownPayment, card.Heading)
}

func (service *playerService) RiskBusiness(card entity.CardRiskBusiness, player entity.Player, rolledDice int) error {
	cost := card.Cost

	if player.Cash < cost {
		return fmt.Errorf(helper.GetMessage("ERROR_NOT_ENOUGH_MONEY"))
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

		realEstate := entity.CardRealEstate{
			ID:          card.ID,
			Type:        card.Type,
			Symbol:      card.Symbol,
			Heading:     card.Heading,
			Description: card.Description,
			Cost:        card.Cost,
			CashFlow:    &cashFlow,
		}

		player.Assets.RealEstates = append(player.Assets.RealEstates, realEstate)
		player.Income.RealEstates = append(player.Income.RealEstates, realEstate)
		player.Liabilities.RealEstates = append(player.Liabilities.RealEstates, realEstate)

		return nil
	}

	return fmt.Errorf(helper.GetMessage("RISK_REQUEST_DECLINED"))
}

func (service *playerService) RiskStocks(card entity.CardRiskStocks, player entity.Player, rolledDice int) error {
	cost := card.Cost

	if player.Cash < cost {
		return fmt.Errorf(helper.GetMessage("ERROR_NOT_ENOUGH_MONEY"))
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

		return nil
	}

	return fmt.Errorf(helper.GetMessage("RISK_REQUEST_DECLINED"))
}

func (service *playerService) UpdateCash(player *entity.Player, amount int, details string) {
	currentCash := player.Cash

	player.Cash += amount

	go service.SetTransaction(player.ID, currentCash, player.Cash, amount, details)

	//const currentCash = this.#cash;
	//this.#cash += Number(amount);
	//const totalCash = this.#cash;
	//this.#recordToLedger(
	//{ currentCash, totalCash, amount, description: details }
	//);
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
