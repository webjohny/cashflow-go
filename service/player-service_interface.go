package service

import (
	"github.com/webjohny/cashflow-go/dto"
	"github.com/webjohny/cashflow-go/entity"
)

type PlayerService interface {
	Payday(player entity.Player, card entity.Card) error
	BecomeModerator(raceId uint64, userId uint64) error
	CashFlowDay(player entity.Player, card entity.Card) error
	Doodad(card entity.CardDoodad, player entity.Player) error
	BigBankrupt(player entity.Player) error
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
	TransferStocks(ID string, sender entity.Player, receiver entity.Player, count int) error
	DecreaseStocks(card entity.CardStocks, player entity.Player) error
	IncreaseStocks(card entity.CardStocks, player entity.Player) error
	Charity(card entity.CardCharity, player entity.Player) error
	PayTax(card entity.CardPayTax, player entity.Player) error
	Downsized(player entity.Player, card entity.Card) error
	BornBaby(player entity.Player, card entity.Card) (error, bool)
	MoveOnBigRace(player entity.Player) error
	SetDream(raceId uint64, userId uint64, playerDream entity.PlayerDream) error
	MarketDamage(card entity.CardMarket, player entity.Player) error
	MarketManipulation(card entity.CardMarket, player entity.Player) error
	MarketBusiness(card entity.CardMarketBusiness, player entity.Player) error
	SellAllProperties(player entity.Player) (error, int)
	SetTransaction(player entity.Player, data dto.TransactionDTO) error
	TakeLoan(player entity.Player, amount int) error
	PayLoan(player entity.Player, actionType string, amount int) error
	UpdateCash(player *entity.Player, amount int, data *dto.TransactionDTO) error
	SetPlayerData(raceId uint64, userId uint64, dto entity.PlayerInfoData) error
	GetTransaction(data dto.TransactionDTO) entity.Transaction
	GetPlayerByUsername(username string) entity.Player
	GetPlayerByUsernameAndRaceId(raceId uint64, username string) entity.Player
	GetPlayerByUserIdAndRaceId(raceId uint64, userId uint64) (error, entity.Player)
	GetPlayerByPlayerIdAndRaceId(raceId uint64, playerId uint64) (error, entity.Player)
	GetAllPlayersByRaceId(raceId uint64) []entity.Player
	GetAllStatePlayersByRaceId(raceId uint64) []entity.Player
	GetProfessionById(id uint8) (error, entity.Profession)
	GetRacePlayer(raceId uint64, userId uint64, full bool) (error, dto.GetRacePlayerResponseDTO)
	GetFormattedPlayerResponse(player entity.Player, hasRestrictedFields bool) dto.GetRacePlayerResponseDTO
	InsertPlayer(b *entity.Player) (error, entity.Player)
	UpdatePlayer(b *entity.Player) (error, entity.Player)
}
