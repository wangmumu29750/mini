package model

import "time"

type Order struct {
	BaseModel

	OrderNo         string      `gorm:"size:32;not null;uniqueIndex" json:"orderNo"`
	UserID          uint64      `gorm:"not null;index;uniqueIndex:idx_order_idempotency,priority:1" json:"userId"`
	TrainID         uint64      `gorm:"not null;index" json:"trainId"`
	TrainNo         string      `gorm:"size:32;not null" json:"trainNo"`
	TravelDate      time.Time   `gorm:"type:date;not null;index" json:"travelDate"`
	FromStationID   uint64      `gorm:"not null" json:"fromStationId"`
	FromStationName string      `gorm:"size:64;not null" json:"fromStationName"`
	ToStationID     uint64      `gorm:"not null" json:"toStationId"`
	ToStationName   string      `gorm:"size:64;not null" json:"toStationName"`
	SeatClassCode   string      `gorm:"size:32;not null" json:"seatClassCode"`
	SeatClassName   string      `gorm:"size:64;not null" json:"seatClassName"`
	PassengerName   string      `gorm:"size:64;not null" json:"passengerName"`
	IDCardNo        string      `gorm:"size:32;not null" json:"idCardNo"`
	AmountCents     int64       `gorm:"not null" json:"amountCents"`
	Status          OrderStatus `gorm:"size:24;not null;index" json:"status"`
	PayExpiresAt    time.Time   `gorm:"not null" json:"payExpiresAt"`
	PaidAt          *time.Time  `json:"paidAt,omitempty"`
	IdempotencyKey  string      `gorm:"size:80;uniqueIndex:idx_order_idempotency,priority:2" json:"-"`

	Tickets  []Ticket  `gorm:"foreignKey:OrderID" json:"tickets,omitempty"`
	Payments []Payment `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
}
