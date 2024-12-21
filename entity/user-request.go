package entity

import (
	"time"
)

type UserRequest struct {
	ID            uint64                 `gorm:"primary_key:auto_increment" json:"id"`
	RaceID        uint64                 `gorm:"type:int(11)" json:"race_id"`
	UserID        uint64                 `gorm:"type:int(11)" json:"user_id"`
	Type          string                 `gorm:"type:varchar(20)" json:"type"`
	CurrentCard   string                 `gorm:"type:varchar(150)" json:"current_card"`
	Amount        int                    `gorm:"type:int(11)" json:"amount"`
	Message       string                 `gorm:"type:text" json:"message"`
	RejectMessage string                 `gorm:"type:text" json:"reject_message"`
	Status        int                    `gorm:"type:int(1)" json:"status"`
	Data          map[string]interface{} `gorm:"type:json;serializer:json" json:"data"`
	CreatedAt     time.Time              `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP();not null" json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"user"`
}
