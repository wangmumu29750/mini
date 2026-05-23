package model

type Train struct {
	BaseModel

	TrainNo   string      `gorm:"size:32;not null;uniqueIndex" json:"trainNo"`
	TrainType string      `gorm:"size:16;not null;index" json:"trainType"`
	Status    TrainStatus `gorm:"size:20;not null;default:ACTIVE;index" json:"status"`

	Stops       []TrainStop `gorm:"foreignKey:TrainID" json:"stops,omitempty"`
	Inventories []Inventory `gorm:"foreignKey:TrainID" json:"inventories,omitempty"`
}
