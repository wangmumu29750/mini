package model

type OrderItem struct {
	BaseModel

	OrderID        uint64        `gorm:"not null;index;uniqueIndex:idx_order_passenger,priority:1" json:"orderId"`
	PassengerID    uint64        `gorm:"not null;index;uniqueIndex:idx_order_passenger,priority:2" json:"passengerId"`
	PassengerName  string        `gorm:"size:64;not null" json:"passengerName"`
	IDCardNo       string        `gorm:"size:32;not null" json:"idCardNo"`
	PassengerType  PassengerType `gorm:"size:20;not null" json:"passengerType"`
	SeatType       SeatType      `gorm:"size:20;not null" json:"seatType"`
	TicketType     TicketType    `gorm:"size:20;not null" json:"ticketType"`
	BasePriceCents int64         `gorm:"not null" json:"basePriceCents"`
	RealPriceCents int64         `gorm:"not null" json:"realPriceCents"`
	TicketID       *uint64       `gorm:"index" json:"ticketId,omitempty"`
	TicketNo       string        `gorm:"size:32;index" json:"ticketNo,omitempty"`
}
