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
	"time"
)

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

func (service *playerService) Payday(player entity.Player, card entity.Card) error {
	logger.Info("PlayerService.Payday", map[string]interface{}{
		"playerId": player.ID,
	})

	return service.UpdateCash(&player, player.CalculateCashFlow(), &dto.TransactionDTO{
		CardID:   &card.ID,
		CardType: entity.TransactionCardType.Payday,
		Details:  card.Heading,
	})
}

func (service *playerService) CashFlowDay(player entity.Player, card entity.Card) error {
	logger.Info("PlayerService.CashFlowDay", map[string]interface{}{
		"playerId": player.ID,
	})

	return service.UpdateCash(&player, player.CalculateCashFlow(), &dto.TransactionDTO{
		CardID:   &card.ID,
		CardType: entity.TransactionCardType.CashFlowDay,
		Details:  card.Heading,
	})
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

		profession := service.professionService.GetRandomProfession(&[]int{
			int(player.ProfessionID),
		})
		player.Reset(profession)
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
		return errors.New(storage.WarnYouHaveNoBabies)
	}

	if player.Cash < cost {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	transaction := dto.TransactionDTO{
		CardID:   &card.ID,
		CardType: entity.TransactionCardType.Doodad,
		Details:  card.Heading,
		PlayerID: player.ID,
		RaceID:   player.RaceID,
	}

	if trx := service.GetTransaction(transaction); trx.ID != 0 {
		return errors.New(storage.ErrorTransactionAlreadyExists)
	}

	return service.UpdateCash(&player, -cost, &transaction)
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

	return service.UpdateCash(&player, -cost, &dto.TransactionDTO{
		CardID:   &card.ID,
		CardType: entity.TransactionCardType.Dream,
		Details:  card.Heading,
	})
}

func (service *playerService) BuyLottery(card entity.CardLottery, player entity.Player, dice int) (error, bool) {
	logger.Info("PlayerService.BuyLottery", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	if card.Cost > player.Cash {
		return errors.New(storage.ErrorNotEnoughMoney), false
	}

	if helper.Contains[int](card.Success, dice) {
		var amount int

		if card.AssetType == entity.LotteryTypes.CashFlow {
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
		} else {
			amount = card.Outcome.Success - card.Cost
		}

		err := service.UpdateCash(&player, amount, &dto.TransactionDTO{
			CardID:   &card.ID,
			CardType: entity.TransactionCardType.Lottery,
			Details:  card.Heading,
		})

		return err, true
	}

	err := service.UpdateCash(&player, -card.Cost, &dto.TransactionDTO{
		CardID:   &card.ID,
		CardType: entity.TransactionCardType.Lottery,
		Details:  card.Heading,
	})

	return err, false
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
		err := service.UpdateCash(&player, totalCost, &dto.TransactionDTO{
			CardID:   &card.ID,
			CardType: entity.TransactionCardType.MarketOther,
			Details:  card.Heading,
		})

		if err != nil {
			return err
		}
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

	return service.UpdateCash(&player, -amount, &dto.TransactionDTO{
		CardID:   &card.ID,
		CardType: entity.TransactionCardType.Charity,
		Details:  card.Heading,
	})
}

func (service *playerService) PayTax(card entity.CardPayTax, player entity.Player) error {
	logger.Info("PlayerService.PayTax", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	var result int

	if card.Percent < 100 {
		result = int((float32(player.Cash) / 100) * float32(card.Percent))
	} else {
		result = player.Cash
	}

	err := service.UpdateCash(&player, -result, &dto.TransactionDTO{
		CardID:   &card.ID,
		CardType: entity.TransactionCardType.PayTax,
		Details:  card.Heading,
	})

	if err != nil {
		return err
	}

	return nil
}

func (service *playerService) Downsized(player entity.Player, card entity.Card) error {
	logger.Info("PlayerService.Downsized", map[string]interface{}{
		"playerId": player.ID,
	})

	player.InitializeSkippedTurns()

	amount := player.CalculateTotalExpenses()

	if player.Cash < amount {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	err := service.UpdateCash(&player, -amount, &dto.TransactionDTO{
		CardID:   &card.ID,
		CardType: entity.TransactionCardType.Downsized,
		Details:  card.Heading,
	})

	if err != nil {
		return err
	}

	return service.AreYouBankrupt(player)
}

func (service *playerService) BornBaby(player entity.Player, card entity.Card) (error, bool) {
	logger.Info("PlayerService.BornBaby", map[string]interface{}{
		"playerId": player.ID,
	})

	if player.Babies > 2 {
		return errors.New(storage.MessageYouHaveTooManyBabies), true
	}

	transaction := dto.TransactionDTO{
		CardID:   &card.ID,
		CardType: entity.TransactionCardType.Baby,
		Details:  card.Heading,
		PlayerID: player.ID,
		RaceID:   player.RaceID,
	}

	if trx := service.GetTransaction(transaction); trx.ID != 0 {
		return errors.New(storage.ErrorTransactionAlreadyExists), false
	}

	var err error

	if err = service.SetTransaction(player, transaction); err != nil {
		return err, false
	}

	player.BornBaby()

	err, _ = service.UpdatePlayer(&player)

	return err, false
}

func (service *playerService) BigBankrupt(player entity.Player) error {
	logger.Info("PlayerService.BigBankrupt", map[string]interface{}{
		"playerId": player.ID,
	})

	var profession entity.Profession

	if player.ProfessionID > 0 {
		profession = service.professionService.GetRandomProfession(&[]int{})
	} else {
		profession = player.Info.Profession
	}

	player.Reset(profession)

	err, _ := service.UpdatePlayer(&player)

	return err
}

func (service *playerService) MoveOnBigRace(player entity.Player) error {
	logger.Info("PlayerService.MoveOnBigRace", map[string]interface{}{
		"playerId": player.ID,
	})

	if !player.ConditionsForBigRace() {
		return errors.New(storage.ErrorMovingBigRaceDeclined)
	}

	cashFlow := player.CalculatePassiveIncome() * 100

	player.OnBigRace = true
	player.CashFlow = cashFlow
	player.TotalIncome = 0
	player.TotalExpenses = 0
	player.CurrentPosition = 0
	player.LastPosition = 0
	player.HasBankrupt = 0
	player.IsRolledDice = 0
	player.SkippedTurns = 0
	player.DualDiceCount = 0
	player.ExtraDices = 0
	player.Salary = 0
	player.Dices = make([]int, 0)
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

func (service *playerService) SetDream(raceId uint64, userId uint64, playerDream entity.PlayerDream) error {
	logger.Info("PlayerService.SetDream", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
		"dream":  playerDream,
	})

	anotherPlayer := service.playerRepository.FindPlayerByRaceIdAndInfoDreamId(raceId, playerDream.ID)

	if anotherPlayer.ID != 0 {
		return errors.New(storage.ErrorDreamPlaceHasAlreadyTaken)
	}

	err, player := service.GetPlayerByUserIdAndRaceId(raceId, userId)

	player.Info.Dream = playerDream

	err, _ = service.playerRepository.UpdatePlayer(&player)

	return err
}

func (service *playerService) MarketDamage(card entity.CardMarket, player entity.Player) error {
	logger.Info("PlayerService.MarketDamage", map[string]interface{}{
		"playerId": player.ID,
		"card":     card,
	})

	if player.HasOwnRealEstates() {
		if player.Cash < card.Cost {
			return errors.New(storage.ErrorNotEnoughMoney)
		}

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

		var cost int

		if card.AssetType == entity.MarketTypes.AnyRealEstate {
			cost = card.Cost
		} else if card.AssetType == entity.MarketTypes.EachRealEstate {
			cost = card.Cost * len(realEstates)
		}

		err := service.UpdateCash(&player, -cost, &dto.TransactionDTO{
			CardID:   &card.ID,
			CardType: entity.TransactionCardType.Damage,
			Details:  card.Heading,
		})

		if err != nil {
			return err
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
		err = service.UpdateCash(&player, -card.WholeCost, &dto.TransactionDTO{
			CardID:   &card.ID,
			CardType: entity.TransactionCardType.Other,
			Details:  card.Heading,
		})

		if err != nil {
			return err
		}
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

	err := service.UpdateCash(&player, amount, &dto.TransactionDTO{
		CardType: entity.TransactionCardType.TakeLoan,
		Details: fmt.Sprintf(
			"Взял(а) в кредит $%s",
			strconv.Itoa(amount),
		),
	})

	if err != nil {
		return err
	}

	return service.AreYouBankrupt(player)
}

func (service *playerService) PayLoan(player entity.Player, actionType string, amount int) error {
	logger.Info("PlayerService.PayLoan", map[string]interface{}{
		"playerId":   player.ID,
		"actionType": actionType,
		"amount":     amount,
	})

	if player.Cash < amount {
		return errors.New(storage.ErrorNotEnoughMoney)
	}

	loanMapper := map[string]string{
		"homeMortgage":   "homeMortgage",
		"schoolLoans":    "schoolLoans",
		"carLoans":       "carLoans",
		"bankLoan":       "bankLoan",
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

	if liabilityAmount >= amount {
		liabilityAmount -= amount
	} else {
		liabilityAmount = 0
	}

	logger.Info("PlayerService.PayLoan: dividing", map[string]interface{}{
		"playerId":        player.ID,
		"amount":          amount,
		"liabilityAmount": liabilityAmount,
		"actionType":      loanMapper[actionType],
	})

	if actionType == "homeMortgage" {
		player.Liabilities.HomeMortgage = liabilityAmount
	} else if actionType == "schoolLoans" {
		player.Liabilities.SchoolLoans = liabilityAmount
	} else if actionType == "carLoans" {
		player.Liabilities.CarLoans = liabilityAmount
	} else if actionType == "creditCardDebt" {
		player.Liabilities.CreditCardDebt = liabilityAmount
	} else if actionType == "bankLoan" {
		player.Liabilities.BankLoan = liabilityAmount
	}

	player.Expenses[loanMapper[actionType]] = liabilityAmount / 10

	return service.UpdateCash(&player, -amount, &dto.TransactionDTO{
		CardType: entity.TransactionCardType.PayLoan,
		Details:  "Оплата по кредиту",
	})
}

func (service *playerService) UpdateCash(player *entity.Player, amount int, data *dto.TransactionDTO) error {
	logger.Info("PlayerService.UpdateCash", map[string]interface{}{
		"playerId": player.ID,
		"cash":     player.Cash,
		"amount":   amount,
		"details":  &data.Details,
	})

	currentCash := player.Cash

	player.Cash += amount

	if &data != nil {
		data.Amount = &amount
		data.CurrentCash = &currentCash
		data.UpdatedCash = &player.Cash
		err := service.SetTransaction(*player, *data)

		if err != nil {
			return err
		}
	}

	err, _ := service.playerRepository.UpdatePlayer(player)

	return err
}

func (service *playerService) BecomeModerator(raceId uint64, userId uint64) error {
	logger.Info("PlayerService.BecomeModerator", map[string]interface{}{
		"raceId": raceId,
		"userId": userId,
	})

	player := service.playerRepository.FindPlayerByUserIdAndRaceId(raceId, userId)

	if player.ID == 0 {
		return errors.New(storage.ErrorUndefinedPlayer)
	}

	player.Role = entity.PlayerRoles.Moderator

	err, _ := service.playerRepository.UpdatePlayer(&player)

	return err
}

func (service *playerService) InsertPlayer(b *entity.Player) (error, entity.Player) {
	return service.playerRepository.InsertPlayer(b)
}

func (service *playerService) SetTransaction(player entity.Player, data dto.TransactionDTO) error {
	if data.CardID == nil {
		cardId := helper.CreateHash(data.Details + data.CardType + strconv.Itoa(int(time.Now().Unix())))
		data.CardID = &cardId
	}

	data.PlayerID = player.ID
	data.RaceID = player.RaceID
	data.Username = player.Username
	data.Color = player.Color
	data.CurrentCash = &player.Cash
	return service.transactionService.InsertTransaction(data)
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
	return service.playerRepository.AllActiveByRaceId(raceId)
}

func (service *playerService) GetAllStatePlayersByRaceId(raceId uint64) []entity.Player {
	return service.playerRepository.AllByRaceId(raceId)
}

func (service *playerService) GetTransaction(data dto.TransactionDTO) entity.Transaction {
	return service.transactionService.GetTransaction(data)
}

func (service *playerService) GetPlayerByUserIdAndRaceId(raceId uint64, userId uint64) (error, entity.Player) {
	player := service.playerRepository.FindPlayerByUserIdAndRaceId(raceId, userId)

	if player.ID == 0 {
		return errors.New(storage.ErrorUndefinedPlayer), entity.Player{}
	}

	return nil, player
}

func (service *playerService) GetPlayerByPlayerIdAndRaceId(raceId uint64, playerId uint64) (error, entity.Player) {
	player := service.playerRepository.FindPlayerByPlayerIdAndRaceId(raceId, playerId)

	if player.ID == 0 {
		return errors.New(storage.ErrorUndefinedPlayer), entity.Player{}
	}

	return nil, player
}

func (service *playerService) GetRacePlayer(raceId uint64, userId uint64, full bool) (error, dto.GetRacePlayerResponseDTO) {
	player := service.playerRepository.FindPlayerByUserIdAndRaceId(raceId, userId)

	if player.ID != 0 {
		return nil, service.GetFormattedPlayerResponse(player, full)
	}

	return errors.New(storage.ErrorUndefinedPlayer), dto.GetRacePlayerResponseDTO{}
}

func (service *playerService) GetFormattedPlayerResponse(player entity.Player, hasRestrictedFields bool) dto.GetRacePlayerResponseDTO {
	profession := service.professionService.GetByID(uint64(player.ProfessionID))

	response := dto.GetRacePlayerResponseDTO{
		ID:       player.ID,
		UserId:   player.UserID,
		Username: player.Username,
		Role:     player.Role,
		Color:    player.Color,
		Profile: dto.RacePlayerProfileResponseDTO{
			Babies:        player.Babies,
			TotalIncome:   player.CalculateTotalIncome(),
			TotalExpenses: player.CalculateTotalExpenses(),
			PassiveIncome: player.CalculatePassiveIncome(),
			CashFlow:      player.CalculateCashFlow(),
			ExtraCashFlow: player.CashFlow,
			Cash:          player.Cash,
		},
		Info:              player.Info,
		Profession:        profession,
		IsRolledDice:      player.IsRolledDice == 1,
		LastPosition:      player.LastPosition,
		Transactions:      make([]dto.RacePlayerTransactionsResponseDTO, 0),
		CurrentPosition:   player.CurrentPosition,
		ExtraDices:        player.ExtraDices,
		DualDiceCount:     player.DualDiceCount,
		SkippedTurns:      player.SkippedTurns,
		AllowOnBigRace:    player.ConditionsForBigRace(),
		GameIsCompleted:   player.ConditionsForCompletedBigRace(),
		GoalPassiveIncome: player.GoalPassiveIncomeOnBigRace(),
		GoalPersonalDream: player.GoalPersonalDream(),
		OnBigRace:         player.OnBigRace,
		IsActive:          player.IsActive,
		HasBankrupt:       player.HasBankrupt == 1,
		AboutToBankrupt:   player.AboutToBankrupt,
	}

	if hasRestrictedFields {
		response.Profile.Income = dto.RacePlayerIncomeResponseDTO{
			RealEstates: player.Assets.RealEstates,
			Business:    player.Assets.Business,
			Salary:      player.Salary,
		}
		response.Profile.Expenses = player.Expenses
		response.Profile.Assets = player.Assets
		response.Profile.Liabilities = dto.RacePlayerLiabilitiesResponseDTO{
			RealEstates:    player.Assets.RealEstates,
			Business:       player.Assets.Business,
			BankLoan:       player.Liabilities.BankLoan,
			HomeMortgage:   player.Liabilities.HomeMortgage,
			SchoolLoans:    player.Liabilities.SchoolLoans,
			CarLoans:       player.Liabilities.CarLoans,
			CreditCardDebt: player.Liabilities.CreditCardDebt,
		}
		response.Notifications = player.Notifications
		response.Dices = player.Dices
	}

	return response
}
