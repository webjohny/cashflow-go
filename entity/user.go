package entity

type User struct {
	ID       uint64 `gorm:"primary_key:auto_increment" json:"id"`
	Name     string `gorm:"type:varchar(255)" json:"name"`
	Email    string `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password string `gorm:"->;<-;not null" json:"-"`
	Profile  string `gorm:"type:varchar(255)" json:"profile"`
	Jk       string `gorm:"type:varchar(255)" json:"jk"`
	Token    string `gorm:"-" json:"token,omitempty"`

	UserRequests []UserRequest `gorm:"foreignKey:UserID" json:"user_requests"`
}
