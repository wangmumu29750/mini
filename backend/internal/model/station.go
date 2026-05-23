package model

type Station struct {
	BaseModel

	Code   string        `gorm:"size:16;not null;uniqueIndex" json:"code"`
	Name   string        `gorm:"size:64;not null;uniqueIndex" json:"name"`
	City   string        `gorm:"size:64;not null;index" json:"city"`
	Status StationStatus `gorm:"size:20;not null;default:ACTIVE;index" json:"status"`
}
