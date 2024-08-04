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
	BuyOtherAssets(card entity.CardOtherAssets, player entity.Player, count int) error
	BuyOtherAssetsInPartnership(card entity.CardOtherAssets, owner entity.Player, players []entity.Player, parts []dto.CardPurchasePlayerActionDTO) error
	BuyDream(card entity.CardDream, player entity.Player) error
	BuyStocks(card entity.CardStocks, player entity.Player, updateCash bool) error
	SellOtherAssets(ID string, card entity.CardMarketOtherAssets, player entity.Player, count int) error
	SellStocks(card entity.CardStocks, player entity.Player, count int, updateCash bool) error
	SellRealEstate(ID string, card entity.CardMarketRealEstate, player entity.Player) error
	SellBusiness(ID string, card entity.CardMarketBusiness, player entity.Player, count int) (error, int)
	TransferBusiness(ID string, sender entity.Player, receiver entity.Player, count int) error
	TransferStocks(card entity.CardStocks, ID string, sender entity.Player, receiver entity.Player, count int) error
	DecreaseStocks(card entity.CardStocks, player entity.Player) error
	IncreaseStocks(card entity.CardStocks, player entity.Player) error
	Charity(card entity.CardCharity, player entity.Player) error
	PayTax(card entity.CardPayTax, player entity.Player) error
	Downsized(player entity.Player) error
	MoveToBigRace(player entity.Player) error
	MarketDamage(card entity.CardMarket, player entity.Player) error
	MarketManipulation(card entity.CardMarket, player entity.Player) error
	MarketBusiness(card entity.CardMarketBusiness, player entity.Player) error
	SellAllProperties(player entity.Player) (error, int)
	TakeLoan(player entity.Player, amount int) error
	PayLoan(player entity.Player, actionType string, amount int) error
	UpdateCash(player *entity.Player, amount int, details string)
	SetTransaction(ID uint64, currentCash int, cash int, amount int, details string)
	GetPlayerByUsername(username string) entity.Player
	GetPlayerByUsernameAndRaceId(raceId uint64, username string) entity.Player
	GetPlayerByUserIdAndRaceId(raceId uint64, userId uint64) (error, entity.Player)
	GetAllPlayersByRaceId(raceId uint64) []entity.Player
	GetProfessionById(id uint8) (error, entity.Profession)
	GetRacePlayer(raceId uint64, userId uint64) (error, dto.GetRacePlayerResponseDTO)
	GetFormattedPlayerResponse(player entity.Player, hasRestrictedFields bool) dto.GetRacePlayerResponseDTO
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

func (service *playerService) AreYouBankrupt(player entity.Player) error {
	if player.IsBankrupt() {
		if !service.playerRepository.IsCurrentPlayerOnTheRace(player) {
			return nil
		}
		logger.Info("PlayerService.AreYouBankrupt", map[string]interface{}{
			"playerId": player.ID,
		})

		players := service.GetAllPlayersByRaceId(player.RaceID)

		businesses := player.Assets.Business
		realEstates := player.Assets.RealEstates
		otherAssets := player.Assets.OtherAssets

		for _, business := range businesses {
			if business.IsOwner {
				for _, anotherPlayer := range players {
					anotherPlayer.RemoveBusiness(business.ID)
					_, _ = service.UpdatePlayer(&anotherPlayer)
				}
			}
		}

		for _, realEs := range realEstates {
			if realEs.IsOwner {
				for _, anotherPlayer := range players {
					anotherPlayer.RemoveRealEstate(realEs.ID)
					_, _ = service.UpdatePlayer(&anotherPlayer)
				}
			}
		}

		for _, asset := range otherAssets {
			if asset.IsOwner {
				for _, anotherPlayer := range players {
					anotherPlayer.RemoveOtherAssetsByID(asset.ID)
					_, _ = service.UpdatePlayer(&anotherPlayer)
				}
			}
		}

		player.Reset(entity.Profession{})
		player.HasBankrupt = 1

		err, _ := service.playerRepository.UpdatePlayer(&player)

		if err != nil {
			logger.Error(err)

			return err
		}

		return errors.New(storage.ErrorYouAreBankrupt)
	}

	return nil
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
		var amount int

		if card.AssetType == entity.LotteryTypes.Money {
			amount = card.Outcome.Success - card.Cost
		} else {
			amount = -card.Cost
			business := entity.CardBusiness{
				ID:          card.ID,
				Type:        card.Type,
				Symbol:      card.Symbol,
				Heading:     card.Heading,
				Description: card.Description,
				Cost:        card.Cost,
				CashFlow:    card.Outcome.Success,
			}

			player.Assets.Business = append(player.Assets.Business, business)
		}
		service.UpdateCash(&player, amount, card.Symbol)

		return nil, true
	}

	service.UpdateCash(&player, -card.Cost, card.Symbol)

	return nil, false
}

func (service *playerService) SellOtherAssets(ID string, card entity.CardMarketOtherAssets, player entity.Player, count int) error {
	logger.Info("PlayerService.SellOtherAssets", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	_, asset := player.FindOtherAssetsByID(ID)

	if asset.ID == "" {
		return errors.New(storage.ErrorNotFoundAssets)
	}

	if !asset.IsOwner {
		return errors.New(storage.ErrorForbiddenByOwner)
	}

	if asset.Count < count {
		return errors.New(storage.ErrorNotEnoughAsset)
	}

	var totalCost = card.Cost

	if asset.AssetType == entity.OtherAssetTypes.Piece && count > 0 {
		totalCost *= count
		asset.Count -= count
	}

	if asset.Count <= 0 || asset.AssetType != entity.OtherAssetTypes.Piece {
		player.RemoveOtherAssetsByID(ID)
	}

	if totalCost > 0 {
		service.UpdateCash(&player, totalCost, card.Heading)
	}

	if asset.AssetType == entity.OtherAssetTypes.Piece {
		return nil
	}

	players := service.GetAllPlayersByRaceId(player.RaceID)

	for _, user := range players {
		_, item := user.FindOtherAssetsByID(ID)

		if item.ID != "" {
			user.RemoveOtherAssetsByID(ID)

			err, play := service.UpdatePlayer(&user)

			if err != nil {
				logger.Error("SellOtherAssets.UpdatePlayer", play, ID, user.ID, user.RaceID)
			}
		}
	}

	return nil
}

func (service *playerService) Charity(card entity.CardCharity, player entity.Player) error {
	logger.Info("PlayerService.Charity", map[string]interface{}{
		"playerId": player.ID,
	})

	amount := card.Cost

	if card.Percent > 0 {
		amount = (player.CalculateTotalIncome() / 100) * card.Percent
	}

	if player.Cash < amount {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	player.DualDiceCount += card.Limit
	player.ExtraDices = card.ExtraDices + 1

	service.UpdateCash(&player, -amount, card.Heading)

	return nil
}

func (service *playerService) PayTax(card entity.CardPayTax, player entity.Player) error {
	logger.Info("PlayerService.PayTax", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	amount := (player.Cash / 100) * card.Percent

	service.UpdateCash(&player, -amount, card.Heading)

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

	player.OnBigRace = true
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

func (service *playerService) MarketDamage(card entity.CardMarket, player entity.Player) error {
	logger.Info("PlayerService.MarketDamage", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	if player.Cash < card.Cost {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	if player.HasOwnRealEstates() {
		realEstates := player.Assets.RealEstates

		if card.Symbol != "ANY" {
			_, asset := player.FindRealEstateBySymbol(card.Symbol)

			if asset.ID == "" || !asset.IsOwner {
				return nil
			}
			realEstates = []entity.CardRealEstate{
				*asset,
			}
		}

		if card.AssetType == entity.MarketTypes.AnyRealEstate {
			service.UpdateCash(&player, -card.Cost, "Имущество поврежденно")
		} else if card.AssetType == entity.MarketTypes.EachRealEstate {
			cost := card.Cost * len(realEstates)
			service.UpdateCash(&player, -cost, "Каждое имущество поврежденно")
		}
	}

	return nil
}

func (service *playerService) MarketManipulation(card entity.CardMarket, player entity.Player) error {
	logger.Info("PlayerService.MarketManipulation", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	if card.Type == "inflation" {
		realEstates := player.Assets.RealEstates

		//@toDo make removing by ID for any realEstate
		if card.Symbol != "ANY" {
			assets := player.FindAllRealEstateBySymbol(card.Symbol)

			if len(assets) == 0 {
				return nil
			}

			realEstates = assets
		}

		for _, asset := range realEstates {
			if !asset.IsOwner {
				continue
			}

			player.RemoveRealEstate(asset.ID)

			if card.AssetType == entity.MarketTypes.AnyRealEstate {
				break
			}
		}
	}

	if card.Type == "success" {
		businesses := player.Assets.Business

		//@toDo make by ID for any business
		if card.Symbol != "ANY" {
			assets := player.FindAllBusinessBySymbol(card.Symbol)

			if len(assets) == 0 {
				return nil
			}

			businesses = assets
		}

		for i, asset := range businesses {
			if !asset.IsOwner || asset.AssetType == entity.BusinessTypes.Limited {
				continue
			}

			percent := asset.Percent

			if percent == 0 {
				percent = 100
			}

			businesses[i].CashFlow += (card.CashFlow / 100) * percent

			if card.AssetType == entity.MarketTypes.AnyBusiness ||
				card.AssetType == entity.MarketTypes.AnyStartup {
				break
			}
		}
		player.Assets.Business = businesses
	}

	err, _ := service.UpdatePlayer(&player)

	if err != nil {
		return err
	}

	return service.AreYouBankrupt(player)
}

func (service *playerService) BuyOtherAssets(card entity.CardOtherAssets, player entity.Player, count int) error {
	logger.Info("PlayerService.BuyOtherAssets", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	if card.IsOwner {
		if player.Cash < card.WholeCost {
			return errors.New(storage.ErrorNotEnoughMoney)
		}
	}

	if card.AssetType == entity.OtherAssetTypes.Piece && card.Count < count {
		return errors.New(storage.ErrorTooManyAssets)
	}

	_, asset := player.FindOtherAssetsBySymbol(card.Symbol)

	var err error

	if card.AssetType == entity.OtherAssetTypes.Piece {
		if asset.ID != "" {
			asset.Count += count
			asset.SumCost()
		} else if asset.ID == "" {
			card.Count = count
			card.SumCost()
		}
	}

	if asset.ID == "" || card.AssetType == entity.OtherAssetTypes.Whole {
		player.Assets.OtherAssets = append(player.Assets.OtherAssets, card)
	}

	if card.IsOwner && card.WholeCost > 0 {
		service.UpdateCash(&player, -card.WholeCost, "Другие активы: "+card.Heading)
	} else {
		err, _ = service.UpdatePlayer(&player)
	}

	return err
}

func (service *playerService) BuyOtherAssetsInPartnership(card entity.CardOtherAssets, owner entity.Player, players []entity.Player, parts []dto.CardPurchasePlayerActionDTO) error {
	logger.Info("PlayerService.BuyOtherAssetsInPartnership", map[string]interface{}{
		"ownerId": owner.ID,
		"card":    helper.JsonSerialize(card),
		"parts":   parts,
	})

	var cardCost = card.Cost

	if card.AssetType == entity.OtherAssetTypes.Piece {
		cardCost = 0
		for _, pl := range parts {
			cardCost += pl.Amount * card.CostPerOne
		}
	}

	if owner.Cash < cardCost {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	for _, pl := range parts {
		var currentPlayer entity.Player

		for _, person := range players {
			if int(person.ID) == pl.ID {
				currentPlayer = person
			}
		}

		if card.AssetType == entity.OtherAssetTypes.Piece && pl.Amount > 0 {
			card.Cost = pl.Amount * card.CostPerOne
		} else if card.AssetType == entity.OtherAssetTypes.Whole {
			card.Cost = pl.Amount
		} else {
			return errors.New(storage.ErrorForbidden)
		}

		if owner.ID == currentPlayer.ID {
			card.WholeCost = cardCost
		} else {
			card.WholeCost = 0
		}

		card.IsOwner = owner.ID == currentPlayer.ID || card.AssetType == entity.OtherAssetTypes.Piece

		err := service.BuyOtherAssets(card, currentPlayer, pl.Amount)

		if err != nil {
			logger.Error(err, nil)

			return err
		}
	}

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

	return service.AreYouBankrupt(player)
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
	return service.playerRepository.InsertPlayer(b)
}

func (service *playerService) UpdatePlayer(b *entity.Player) (error, entity.Player) {
	return service.playerRepository.UpdatePlayer(b)
}

func (service *playerService) GetPlayerByUsername(username string) entity.Player {
	return service.playerRepository.FindPlayerByUsername(username)
}

func (service *playerService) GetProfessionById(id uint8) (error, entity.Profession) {
	profession := service.professionService.GetByID(uint64(id))

	if profession.ID == 0 {
		return errors.New(storage.ErrorUndefinedProfession), entity.Profession{}
	}

	return nil, profession
}

func (service *playerService) GetPlayerByUsernameAndRaceId(raceId uint64, username string) entity.Player {
	return service.playerRepository.FindPlayerByUsernameAndRaceId(raceId, username)
}

func (service *playerService) GetAllPlayersByRaceId(raceId uint64) []entity.Player {
	return service.playerRepository.AllByRaceId(raceId)
}

func (service *playerService) GetPlayerByUserIdAndRaceId(raceId uint64, userId uint64) (error, entity.Player) {
	player := service.playerRepository.FindPlayerByUserIdAndRaceId(raceId, userId)

	if player.ID == 0 {
		return errors.New(storage.ErrorUndefinedPlayer), entity.Player{}
	}

	return nil, player
}

func (service *playerService) GetRacePlayer(raceId uint64, userId uint64) (error, dto.GetRacePlayerResponseDTO) {
	player := service.playerRepository.FindPlayerByUserIdAndRaceId(raceId, userId)

	if player.ID != 0 {
		return nil, service.GetFormattedPlayerResponse(player, false)
	}

	return errors.New(storage.ErrorUndefinedPlayer), dto.GetRacePlayerResponseDTO{}
}

func (service *playerService) GetFormattedPlayerResponse(player entity.Player, hasRestrictedFields bool) dto.GetRacePlayerResponseDTO {
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

	response := dto.GetRacePlayerResponseDTO{
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
		ExtraDices:      player.ExtraDices,
		DualDiceCount:   player.DualDiceCount,
		SkippedTurns:    player.SkippedTurns,
		CanReRoll:       player.CanReRoll == 1,
		OnBigRace:       player.OnBigRace,
		HasBankrupt:     player.HasBankrupt == 1,
		AboutToBankrupt: player.AboutToBankrupt,
		HasMlm:          player.HasMlm == 1,
	}

	if hasRestrictedFields {
		response.Dices = player.Dices
	}

	return response
}
