package model

import "time"

type Ticket struct {
	BaseModel

	TicketNo        string       `gorm:"size:32;not null;uniqueIndex" json:"ticketNo"`
	OrderID         uint64       `gorm:"not null;index" json:"orderId"`
	UserID          uint64       `gorm:"not null;index" json:"userId"`
	TrainID         uint64       `gorm:"not null;index" json:"trainId"`
	TrainNo         string       `gorm:"size:32;not null" json:"trainNo"`
	TravelDate      time.Time    `gorm:"type:date;not null;index" json:"travelDate"`
	FromStationID   uint64       `gorm:"not null" json:"fromStationId"`
	FromStationName string       `gorm:"size:64;not null" json:"fromStationName"`
	ToStationID     uint64       `gorm:"not null" json:"toStationId"`
	ToStationName   string       `gorm:"size:64;not null" json:"toStationName"`
	SeatClassCode   string       `gorm:"size:32;not null" json:"seatClassCode"`
	SeatClassName   string       `gorm:"size:64;not null" json:"seatClassName"`
	PassengerName   string       `gorm:"size:64;not null" json:"passengerName"`
	IDCardNo        string       `gorm:"size:32;not null" json:"idCardNo"`
	Status          TicketStatus `gorm:"size:20;not null;index" json:"status"`
	IssuedAt        time.Time    `gorm:"not null" json:"issuedAt"`
}
