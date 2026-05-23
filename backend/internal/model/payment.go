package model

import "time"

type Payment struct {
	BaseModel

	PaymentNo   string        `gorm:"size:32;not null;uniqueIndex" json:"paymentNo"`
	OrderID     uint64        `gorm:"not null;index" json:"orderId"`
	UserID      uint64        `gorm:"not null;index" json:"userId"`
	AmountCents int64         `gorm:"not null" json:"amountCents"`
	Channel     string        `gorm:"size:32;not null" json:"channel"`
	Status      PaymentStatus `gorm:"size:20;not null;index" json:"status"`
	PaidAt      time.Time     `gorm:"not null" json:"paidAt"`
}
