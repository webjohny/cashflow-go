package repository

import (
	logger "github.com/sirupsen/logrus"
	"github.com/webjohny/cashflow-go/entity"
	"github.com/webjohny/cashflow-go/helper"
	"gorm.io/gorm"
)

type PlayerRepository interface {
	InsertPlayer(b *entity.Player) (error, entity.Player)
	UpdatePlayer(b *entity.Player) (error, entity.Player)
	UpdateCash(b *entity.Player, cash int)
	AllByRaceId(raceId uint64) []entity.Player
	AllActiveByRaceId(raceId uint64) []entity.Player
	DeletePlayer(b *entity.Player) error
	IsCurrentPlayerOnTheRace(player entity.Player) bool
	FindPlayerByRaceIdAndInfoDreamId(raceId uint64, dreamId int) entity.Player
	FindPlayerById(ID uint64) entity.Player
	FindPlayerByUsername(username string) entity.Player
	FindPlayerByUsernameAndRaceId(raceId uint64, username string) entity.Player
	FindPlayerByUserIdAndRaceId(raceId uint64, userId uint64) entity.Player
	FindPlayerByPlayerIdAndRaceId(raceId uint64, playerId uint64) entity.Player
}

const PlayerTable = "players"

type playerConnection struct {
	connection *gorm.DB
}

func NewPlayerRepository(dbConn *gorm.DB) PlayerRepository {
	return &playerConnection{
		connection: dbConn,
	}
}

func (db *playerConnection) InsertPlayer(b *entity.Player) (error, entity.Player) {
	logger.Info("PlayerRepository.InsertPlayer", b.ID)

	result := db.connection.Save(&b)

	if result.Error != nil {
		logger.Error(result.Error, helper.JsonSerialize(b))

		return result.Error, entity.Player{}
	}

	db.connection.Find(&b)
	return nil, *b
}

func (db *playerConnection) AllActiveByRaceId(raceId uint64) []entity.Player {
	var players []entity.Player
	db.connection.Where("race_id = ? AND role != 'moderator' AND is_active = 1", raceId).Find(&players)
	return players
}

func (db *playerConnection) AllByRaceId(raceId uint64) []entity.Player {
	var players []entity.Player
	db.connection.Where("race_id = ? AND role != 'moderator'", raceId).Find(&players)
	return players
}

func (db *playerConnection) UpdatePlayer(b *entity.Player) (error, entity.Player) {
	logger.Info("PlayerRepository.UpdatePlayer", b.ID)

	result := db.connection.Select("*").Updates(&b)

	if result.Error != nil {
		logger.Error(result.Error, helper.JsonSerialize(b))

		return result.Error, entity.Player{}
	}

	db.connection.Find(&b)
	return nil, *b
}

func (db *playerConnection) UpdateCash(b *entity.Player, cash int) {
	logger.Info("PlayerRepository.UpdateCash", map[string]interface{}{
		"cash":     cash,
		"playerId": b.ID,
	})

	result := db.connection.Model(&b).Select("Cash").Update("cash", cash)

	if result.Error != nil {
		logger.Error(result.Error, helper.JsonSerialize(b))
	}
}

func (db *playerConnection) DeletePlayer(b *entity.Player) error {
	logger.Info("PlayerRepository.DeletePlayer", helper.JsonSerialize(b))

	result := db.connection.Delete(&b)

	if result.Error != nil {
		logger.Error(result.Error, helper.JsonSerialize(b))

		return result.Error
	}

	return nil
}

func (db *playerConnection) FindPlayerByRaceIdAndInfoDreamId(raceId uint64, dreamId int) entity.Player {
	logger.Info("PlayerRepository.FindPlayerByRaceIdAndInfoDreamId", map[string]interface{}{
		"raceId":  raceId,
		"dreamId": dreamId,
	})

	var player entity.Player

	db.connection.Where("JSON_EXTRACT(info, '$.dream.id') = ?", dreamId).Where("race_id = ? AND role != 'moderator'", raceId).Find(&player)

	return player
}

func (db *playerConnection) FindPlayerById(ID uint64) entity.Player {
	logger.Info("PlayerRepository.FindPlayerById", map[string]interface{}{
		"id": ID,
	})

	var player entity.Player

	db.connection.Find(&player, ID)

	return player
}

func (db *playerConnection) IsCurrentPlayerOnTheRace(player entity.Player) bool {
	logger.Info("PlayerRepository.IsCurrentPlayerOnRace", map[string]interface{}{
		"id": player.ID,
	})

	var count int64
	db.connection.Model(&entity.Race{}).
		Where("id = ? AND role != 'moderator' AND JSON_EXTRACT(current_player, '$.id') = ?", player.RaceID, player.ID).
		Count(&count)

	return count > 0
}

func (db *playerConnection) FindPlayerByUsername(username string) entity.Player {
	logger.Info("PlayerRepository.FindPlayerByUsername", map[string]interface{}{
		"username": username,
	})

	var player entity.Player

	db.connection.Find(&player, "`username` = ?", username).Find(&player)

	return player
}

func (db *playerConnection) FindPlayerByUsernameAndRaceId(raceId uint64, username string) entity.Player {
	logger.Info("PlayerRepository.FindPlayerByUsernameAndRaceId", map[string]interface{}{
		"raceId":   raceId,
		"username": username,
	})

	var player entity.Player

	db.connection.Where("`username` = ? AND `race_id` = ?", username, raceId).Find(&player)

	return player
}

func (db *playerConnection) FindPlayerByUserIdAndRaceId(raceId uint64, userId uint64) entity.Player {
	//logger.Info("PlayerRepository.FindPlayerByUserIdAndRaceId", map[string]interface{}{
	//	"raceId": raceId,
	//	"userId": userId,
	//})

	var player entity.Player

	db.connection.Where("`user_id` = ? AND `race_id` = ?", userId, raceId).Find(&player)

	return player
}

func (db *playerConnection) FindPlayerByPlayerIdAndRaceId(raceId uint64, playerId uint64) entity.Player {
	logger.Info("PlayerRepository.FindPlayerByPlayerIdAndRaceId", map[string]interface{}{
		"raceId":   raceId,
		"playerId": playerId,
	})

	var player entity.Player

	db.connection.Where("`id` = ? AND `race_id` = ?", playerId, raceId).Find(&player)

	return player
}
