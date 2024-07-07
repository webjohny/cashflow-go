package service

import (
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"github.com/webjohny/cashflow-go/repository"
	"github.com/webjohny/cashflow-go/storage"
	"math"
	"strconv"
)

type PlayerService interface {
	Payday(player entity.Player)
	CashFlowDay(player entity.Player)
	Doodad(card entity.CardDoodad, player entity.Player) error
	BuyBusiness(card entity.CardBusiness, player entity.Player, count int, updateCash bool) error
	BuyRealEstate(card entity.CardRealEstate, player entity.Player) error
	BuyRealEstateInPartnership(card entity.CardRealEstate, owner entity.Player, players []entity.Player, parts []dto.CardPurchasePlayerActionDTO) error
	BuyBusinessInPartnership(card entity.CardBusiness, owner entity.Player, players []entity.Player, parts []dto.CardPurchasePlayerActionDTO) error
	BuyLottery(card entity.CardLottery, player entity.Player, dice int) (error, bool)
	BuyRiskBusiness(card entity.CardRiskBusiness, player entity.Player, rolledDice int) (error, bool)
	BuyOtherAssets(card entity.CardOtherAssets, player entity.Player, count int) error
	BuyRiskStocks(card entity.CardRiskStocks, player entity.Player, rolledDice int) (error, bool)
	BuyDream(card entity.CardDream, player entity.Player) error
	BuyStocks(card entity.CardStocks, player entity.Player, updateCash bool) error
	SellOtherAssets(card entity.CardOtherAssets, player entity.Player) error
	SellStocks(card entity.CardStocks, player entity.Player, count int, updateCash bool) error
	SellRealEstate(ID string, card entity.CardMarketRealEstate, player entity.Player) error
	SellBusiness(ID string, card entity.CardMarketBusiness, player entity.Player) (error, int)
	TransferBusiness(ID string, sender entity.Player, receiver entity.Player, count int) error
	TransferStocks(ID string, sender entity.Player, receiver entity.Player, count int) error
	DecreaseStocks(card entity.CardStocks, player entity.Player) error
	IncreaseStocks(card entity.CardStocks, player entity.Player) error
	Charity(player entity.Player) error
	BigCharity(card entity.CardCharity, player entity.Player) error
	PayTax(card entity.CardPayTax, player entity.Player) error
	Downsized(player entity.Player) error
	MoveToBigRace(player entity.Player) error
	MarketPayDamages(card entity.CardMarket, player entity.Player) error
	MarketBusiness(card entity.CardMarketBusiness, player entity.Player) error
	SellAllProperties(player entity.Player) (error, int)
	TakeLoan(player entity.Player, amount int) error
	PayLoan(player entity.Player, actionType string, amount int) error
	UpdateCash(player *entity.Player, amount int, details string)
	SetTransaction(ID uint64, currentCash int, cash int, amount int, details string)
	GetPlayerByUsername(username string) entity.Player
	GetPlayerByUsernameAndRaceId(raceId uint64, username string) entity.Player
	GetPlayerByUserIdAndRaceId(raceId uint64, userId uint64) entity.Player
	GetAllPlayersByRaceId(raceId uint64) []entity.Player
	GetProfessionById(id uint8) (error, entity.Profession)
	GetRacePlayer(raceId uint64, userId uint64) (error, dto.GetRacePlayerResponseDTO)
	GetFormattedPlayerResponse(player entity.Player) dto.GetRacePlayerResponseDTO
	InsertPlayer(b *entity.Player) (error, entity.Player)
	UpdatePlayer(b *entity.Player) (error, entity.Player)
}

type playerService struct {
	playerRepository   repository.PlayerRepository
	professionService  ProfessionService
	transactionService TransactionService
}

func NewPlayerService(playerRepo repository.PlayerRepository, professionService ProfessionService, transactionService TransactionService) PlayerService {
	return &playerService{
		playerRepository:   playerRepo,
		professionService:  professionService,
		transactionService: transactionService,
	}
}

func (service *playerService) Payday(player entity.Player) {
	logger.Info("PlayerService.Payday", map[string]interface{}{
		"playerId": player.ID,
	})

	service.UpdateCash(&player, player.CalculateCashFlow(), "Зарплата")
}

func (service *playerService) CashFlowDay(player entity.Player) {
	logger.Info("PlayerService.CashFlowDay", map[string]interface{}{
		"playerId": player.ID,
	})

	service.UpdateCash(&player, player.CalculateCashFlow(), "Кэш-флоу день")
}

func (service *playerService) AreYouBankrupt(player entity.Player) {
	logger.Info("PlayerService.AreYouBankrupt", map[string]interface{}{
		"playerId": player.ID,
	})

	if player.IsBankrupt() {
		players := service.GetAllPlayersByRaceId(player.RaceID)

		businesses := player.Assets.Business
		realEstates := player.Assets.RealEstates
		//otherAssets := player.Assets.OtherAssets

		for _, business := range businesses {
			if business.IsOwner {
				for _, anotherPlayer := range players {
					_, item := anotherPlayer.FindBusinessBySymbol(business.Symbol)

					if item.ID != "" {
						anotherPlayer.RemoveBusiness(item.ID)
						_, _ = service.UpdatePlayer(&anotherPlayer)
					}
				}
			}
		}

		for _, realEs := range realEstates {
			if realEs.IsOwner {
				for _, anotherPlayer := range players {
					item := anotherPlayer.FindRealEstate(realEs.Symbol)

					if item.ID != "" {
						anotherPlayer.RemoveRealEstate(item.ID)
						_, _ = service.UpdatePlayer(&anotherPlayer)
					}
				}
			}
		}

		//player.Reset()
	}
}

func (service *playerService) Doodad(card entity.CardDoodad, player entity.Player) error {
	logger.Info("PlayerService.Doodad", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	cost := card.Cost

	if card.HasBabies && player.Babies <= 0 {
		return errors.New(storage.ErrorYouHaveNoBabies)
	}

	if player.Cash < cost {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	service.UpdateCash(&player, -cost, "Растраты")

	return nil
}

func (service *playerService) BuyDream(card entity.CardDream, player entity.Player) error {
	logger.Info("PlayerService.BuyDream", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	cost := card.Cost

	if player.Cash < cost {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	player.Assets.Dreams = append(player.Assets.Dreams, card)

	service.UpdateCash(&player, -cost, "Мечта")

	return nil
}

func (service *playerService) BuyLottery(card entity.CardLottery, player entity.Player, dice int) (error, bool) {
	logger.Info("PlayerService.BuyLottery", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	if card.Cost > player.Cash {
		return errors.New(storage.ErrorNotEnoughMoney), false
	}

	player.AllowReRoll()

	if helper.Contains[int](card.Success, dice) {
		amount := card.Outcome.Success - card.Cost
		service.UpdateCash(&player, amount, card.Symbol)

		return nil, true
	}

	service.UpdateCash(&player, -card.Cost, card.Symbol)

	return nil, false
}

func (service *playerService) SellOtherAssets(card entity.CardOtherAssets, player entity.Player) error {
	logger.Info("PlayerService.SellOtherAssets", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	_, asset := player.FindOtherAssets(card.Symbol)

	if asset.ID == "" {
		return errors.New(storage.ErrorNotFoundAssets)
	}

	if asset.Count < card.Count {
		return errors.New(storage.ErrorNotEnoughAsset)
	}

	asset.Count -= card.Count
	asset.Cost = asset.CostPerOne * asset.Count

	if asset.Count <= 0 {
		player.RemoveOtherAssets(asset.Symbol)
	} else {
		player.UpdateAsset(asset.Symbol, *asset)
	}

	service.UpdateCash(&player, card.CostPerOne*card.Count, card.Symbol)

	return nil
}

func (service *playerService) Charity(player entity.Player) error {
	logger.Info("PlayerService.Charity", map[string]interface{}{
		"playerId": player.ID,
	})

	amount := int(math.Floor(0.1 * float64(player.CalculateTotalIncome())))

	if player.Cash < amount {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	player.IncrementDualDiceCount()

	service.UpdateCash(&player, -amount, "Благотворительность")

	return nil
}

func (service *playerService) BigCharity(card entity.CardCharity, player entity.Player) error {
	logger.Info("PlayerService.BigCharity", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	if player.Cash < card.Cost {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	service.UpdateCash(&player, -card.Cost, "Акция милосердия")

	return nil
}

func (service *playerService) PayTax(card entity.CardPayTax, player entity.Player) error {
	logger.Info("PlayerService.PayTax", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	amount := (player.Cash / 100) * card.Percent

	if player.Cash < amount {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	service.UpdateCash(&player, -amount, "Налоги")

	return nil
}

func (service *playerService) Downsized(player entity.Player) error {
	logger.Info("PlayerService.Downsized", map[string]interface{}{
		"playerId": player.ID,
	})

	player.InitializeSkippedTurns()

	amount := player.CalculateTotalExpenses()

	if player.Cash < amount {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	service.UpdateCash(&player, -amount, "Уволен")

	return nil
}

func (service *playerService) MoveToBigRace(player entity.Player) error {
	logger.Info("PlayerService.MoveToBigRace", map[string]interface{}{
		"playerId": player.ID,
	})

	if !player.ConditionsForBigRace() {
		return errors.New(storage.ErrorMovingBigRaceDeclined)
	}

	cashFlow := player.CalculatePassiveIncome() * 100

	player.OnBigRace = 1
	player.CashFlow = cashFlow
	player.Cash = cashFlow + player.Cash
	player.TotalIncome = 0
	player.TotalExpenses = 0
	player.Expenses = make(map[string]int)
	player.Assets = entity.PlayerAssets{
		Savings:     0,
		Stocks:      make([]entity.CardStocks, 0),
		OtherAssets: make([]entity.CardOtherAssets, 0),
		RealEstates: make([]entity.CardRealEstate, 0),
		Business:    make([]entity.CardBusiness, 0),
		Dreams:      make([]entity.CardDream, 0),
	}

	err, _ := service.playerRepository.UpdatePlayer(&player)

	return err
}

func (service *playerService) MarketPayDamages(card entity.CardMarket, player entity.Player) error {
	logger.Info("PlayerService.MarketPayDamages", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	if player.Cash < card.Cost {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	if !player.HasRealEstates() {
		return errors.New(storage.ErrorYouHaveNoProperties)
	}

	service.UpdateCash(&player, -card.Cost, "Имущество поврежденно")

	return nil
}

func (service *playerService) BuyOtherAssets(card entity.CardOtherAssets, player entity.Player, count int) error {
	logger.Info("PlayerService.BuyOtherAssets", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	if player.Cash < card.Cost {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	if card.Count < count {
		return errors.New(storage.ErrorTooManyAssets)
	}

	_, asset := player.FindOtherAssets(card.Symbol)

	if asset.ID != "" {
		asset.Count += count
		asset.Cost += count * card.CostPerOne
		asset.CostPerOne = asset.Cost * asset.Count
	} else {
		card.Count = count
		card.Cost = count * card.CostPerOne

		player.Assets.OtherAssets = append(player.Assets.OtherAssets, card)
	}

	service.UpdateCash(&player, -card.Cost, "Другие активы: "+card.Heading)

	return nil
}

func (service *playerService) TakeLoan(player entity.Player, amount int) error {
	logger.Info("PlayerService.TakeLoan", map[string]interface{}{
		"playerId": player.ID,
		"amount":   amount,
	})

	player.Liabilities.BankLoan += amount
	player.Expenses["bankLoan"] = player.Liabilities.BankLoan / 10

	service.UpdateCash(&player, amount, fmt.Sprintf(
		"Взял(а) в кредит $%s",
		strconv.Itoa(amount),
	))

	return nil
}

func (service *playerService) PayLoan(player entity.Player, actionType string, amount int) error {
	logger.Info("PlayerService.PayLoan", map[string]interface{}{
		"playerId": player.ID,
		"amount":   amount,
	})

	if player.Cash < amount {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	loanMapper := map[string]string{
		"homeMortgage":   "homeMortgage",
		"schoolLoans":    "schoolLoans",
		"carLoans":       "carLoans",
		"creditCardDebt": "creditCardDebt",
	}

	var liabilityAmount int

	if actionType == "homeMortgage" {
		liabilityAmount = player.Liabilities.HomeMortgage
	} else if actionType == "schoolLoans" {
		liabilityAmount = player.Liabilities.SchoolLoans
	} else if actionType == "carLoans" {
		liabilityAmount = player.Liabilities.CarLoans
	} else if actionType == "creditCardDebt" {
		liabilityAmount = player.Liabilities.CreditCardDebt
	} else if actionType == "bankLoan" {
		liabilityAmount = player.Liabilities.BankLoan
	}

	if liabilityAmount > amount {
		liabilityAmount -= amount
	} else {
		liabilityAmount = 0
	}

	if actionType == "bankLoan" {
		player.Expenses["bankLoan"] = liabilityAmount / 10
	} else {
		player.Expenses[loanMapper[actionType]] = liabilityAmount
	}

	service.UpdateCash(&player, -amount, "Оплата по кредиту")

	return nil
}

func (service *playerService) UpdateCash(player *entity.Player, amount int, details string) {
	logger.Info("PlayerService.UpdateCash", map[string]interface{}{
		"playerId": player.ID,
		"cash":     player.Cash,
		"amount":   amount,
		"details":  details,
	})

	currentCash := player.Cash

	player.Cash += amount

	service.SetTransaction(player.ID, currentCash, player.Cash, amount, details)
	err, _ := service.playerRepository.UpdatePlayer(player)

	if err != nil {
		logger.Error(err, map[string]interface{}{
			"playerId": player.ID,
			"cash":     player.Cash,
		})
	}
}

func (service *playerService) SetTransaction(ID uint64, currentCash int, cash int, amount int, details string) {
	logger.Info("PlayerService.SetTransaction", map[string]interface{}{
		"playerId":    ID,
		"currentCash": currentCash,
		"cash":        cash,
		"details":     details,
	})

	service.transactionService.InsertPlayerTransaction(dto.TransactionCreatePlayerDTO{
		PlayerID:    ID,
		Details:     details,
		CurrentCash: currentCash,
		Cash:        cash,
		Amount:      amount,
	})
}

func (service *playerService) InsertPlayer(b *entity.Player) (error, entity.Player) {
	logger.Info("PlayerService.InsertPlayer", helper.JsonSerialize(b))

	return service.playerRepository.InsertPlayer(b)
}

func (service *playerService) UpdatePlayer(b *entity.Player) (error, entity.Player) {
	logger.Info("PlayerService.UpdatePlayer", helper.JsonSerialize(b))

	return service.playerRepository.UpdatePlayer(b)
}

func (service *playerService) GetPlayerByUsername(username string) entity.Player {
	logger.Info("PlayerService.GetPlayerByUsername", map[string]interface{}{
		"username": username,
	})

	return service.playerRepository.FindPlayerByUsername(username)
}

func (service *playerService) GetProfessionById(id uint8) (error, entity.Profession) {
	logger.Info("PlayerService.GetProfessionById", map[string]interface{}{
		"id": id,
	})

	profession := service.professionService.GetByID(uint64(id))

	if profession.ID == 0 {
		return errors.New(storage.ErrorUndefinedProfession), entity.Profession{}
	}

	return nil, profession
}

func (service *playerService) GetPlayerByUsernameAndRaceId(raceId uint64, username string) entity.Player {
	logger.Info("PlayerService.GetPlayerByUsernameAndRaceId", map[string]interface{}{
		"raceId":   raceId,
		"username": username,
	})

	return service.playerRepository.FindPlayerByUsernameAndRaceId(raceId, username)
}

func (service *playerService) GetAllPlayersByRaceId(raceId uint64) []entity.Player {
	//logger.Info("PlayerService.GetAllPlayersByRaceId", map[string]interface{}{
	//	"raceId": raceId,
	//})

	return service.playerRepository.AllByRaceId(raceId)
}

func (service *playerService) GetPlayerByUserIdAndRaceId(raceId uint64, userId uint64) entity.Player {
	logger.Info("PlayerService.GetPlayerByUserIdAndRaceId", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	return service.playerRepository.FindPlayerByUserIdAndRaceId(raceId, userId)
}

func (service *playerService) GetRacePlayer(raceId uint64, userId uint64) (error, dto.GetRacePlayerResponseDTO) {
	//logger.Info("PlayerService.GetRacePlayer", map[string]interface{}{
	//	"raceId": raceId,
	//	"userId": userId,
	//})

	player := service.playerRepository.FindPlayerByUserIdAndRaceId(raceId, userId)

	if player.ID != 0 {
		return nil, service.GetFormattedPlayerResponse(player)
	}

	return errors.New(storage.ErrorUndefinedPlayer), dto.GetRacePlayerResponseDTO{}
}

func (service *playerService) GetFormattedPlayerResponse(player entity.Player) dto.GetRacePlayerResponseDTO {
	//logger.Info("PlayerService.GetFormattedPlayerResponse", helper.JsonSerialize(player))

	profession := service.professionService.GetByID(uint64(player.ProfessionID))
	transactionsQuery := service.transactionService.GetPlayerTransactions(player.ID)

	transactions := make([]dto.RacePlayerTransactionsResponseDTO, 0)

	for i := 0; i < len(transactionsQuery); i++ {
		transactions = append(transactions, dto.RacePlayerTransactionsResponseDTO{
			CurrentCash: *transactionsQuery[i].Data.CurrentCash,
			Cash:        *transactionsQuery[i].Data.Cash,
			Amount:      *transactionsQuery[i].Data.Amount,
			Details:     transactionsQuery[i].Details,
		})
	}

	return dto.GetRacePlayerResponseDTO{
		ID:       player.ID,
		UserId:   player.UserID,
		Username: player.Username,
		Role:     player.Role,
		Color:    player.Color,
		Profile: dto.RacePlayerProfileResponseDTO{
			Income: dto.RacePlayerIncomeResponseDTO{
				RealEstates: player.Assets.RealEstates,
				Business:    player.Assets.Business,
				Salary:      player.Salary,
			},
			Babies:   player.Babies,
			Expenses: player.Expenses,
			Assets:   player.Assets,
			Liabilities: dto.RacePlayerLiabilitiesResponseDTO{
				RealEstates:    player.Assets.RealEstates,
				Business:       player.Assets.Business,
				BankLoan:       player.Liabilities.BankLoan,
				HomeMortgage:   player.Liabilities.HomeMortgage,
				SchoolLoans:    player.Liabilities.SchoolLoans,
				CarLoans:       player.Liabilities.CarLoans,
				CreditCardDebt: player.Liabilities.CreditCardDebt,
			},
			TotalIncome:   player.CalculateTotalIncome(),
			TotalExpenses: player.CalculateTotalExpenses(),
			CashFlow:      player.CalculateCashFlow(),
			PassiveIncome: player.CalculatePassiveIncome(),
			Cash:          player.Cash,
		},
		Profession:      profession,
		IsRolledDice:    player.IsRolledDice == 1,
		LastPosition:    player.LastPosition,
		Transactions:    transactions,
		CurrentPosition: player.CurrentPosition,
		DualDiceCount:   player.DualDiceCount,
		SkippedTurns:    player.SkippedTurns,
		CanReRoll:       player.CanReRoll == 1,
		OnBigRace:       player.OnBigRace == 1,
		HasBankrupt:     player.HasBankrupt == 1,
		AboutToBankrupt: player.AboutToBankrupt,
		HasMlm:          player.HasMlm == 1,
	}
}
