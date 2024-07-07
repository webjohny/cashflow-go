package entity

type UsedCard struct {
	ID     uint64 `gorm:"primary_key:auto_increment" json:"id"`
	RaceID uint64 `gorm:"uniqueIndex:card_index;index" json:"race_id"`
	CardID string `gorm:"uniqueIndex:card_index;type:varchar(100)" json:"card_id"`
	Action string `gorm:"uniqueIndex:card_index;type:varchar(100)" json:"action"`
}
