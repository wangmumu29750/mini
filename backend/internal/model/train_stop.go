package model

type TrainStop struct {
	BaseModel

	TrainID     uint64 `gorm:"not null;uniqueIndex:idx_train_stop_order,priority:1;uniqueIndex:idx_train_station,priority:1;index" json:"trainId"`
	StationID   uint64 `gorm:"not null;uniqueIndex:idx_train_station,priority:2;index" json:"stationId"`
	StopOrder   int    `gorm:"not null;uniqueIndex:idx_train_stop_order,priority:2" json:"stopOrder"`
	DayOffset   int    `gorm:"not null;default:0" json:"dayOffset"`
	ArriveClock string `gorm:"size:8;not null;default:''" json:"arriveClock"`
	DepartClock string `gorm:"size:8;not null;default:''" json:"departClock"`
	Mileage     int    `gorm:"not null;default:0" json:"mileage"`

	Station Station `gorm:"foreignKey:StationID;references:ID" json:"station,omitempty"`
}
