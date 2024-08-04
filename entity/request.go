package entity

import (
	"time"
)

type Request struct {
	ID            uint64      `gorm:"primary_key:auto_increment" json:"id"`
	RaceID        uint64      `gorm:"uniqueIndex:user_index;index" json:"race_id"`
	UserID        uint64      `gorm:"uniqueIndex:user_index;index" json:"user_id"`
	Type          string      `gorm:"type:varchar(20)" json:"type"`
	CurrentCard   string      `gorm:"type:varchar(150)" json:"current_card"`
	Amount        int         `gorm:"type:int(11)" json:"amount"`
	Message       string      `gorm:"type:text" json:"message"`
	RejectMessage string      `gorm:"type:text" json:"reject_message"`
	Approved      bool        `gorm:"default:true" json:"approved"`
	Data          interface{} `gorm:"type:json" json:"data"`
	CreatedAt     time.Time   `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP();not null" json:"created_at"`
}
