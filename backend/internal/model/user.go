package model

import "time"

type User struct {
	BaseModel

	Username     string     `gorm:"size:64;not null;uniqueIndex" json:"username"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	Role         UserRole   `gorm:"size:20;not null;index" json:"role"`
	Status       UserStatus `gorm:"size:20;not null;default:ACTIVE;index" json:"status"`
	LastLoginAt  *time.Time `json:"lastLoginAt,omitempty"`

	Profile PassengerProfile `gorm:"foreignKey:UserID;references:ID" json:"profile,omitempty"`
}
