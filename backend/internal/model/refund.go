package model

import "time"

type Refund struct {
	BaseModel

	RefundNo       string       `gorm:"size:32;not null;uniqueIndex" json:"refundNo"`
	TicketID       uint64       `gorm:"not null;index" json:"ticketId"`
	PaymentID      uint64       `gorm:"not null;index" json:"paymentId"`
	UserID         uint64       `gorm:"not null;index;uniqueIndex:idx_refund_idempotency,priority:1" json:"userId"`
	AmountCents    int64        `gorm:"not null" json:"amountCents"`
	Status         RefundStatus `gorm:"size:20;not null;index" json:"status"`
	Reason         string       `gorm:"size:200" json:"reason"`
	IdempotencyKey string       `gorm:"size:80;uniqueIndex:idx_refund_idempotency,priority:2" json:"-"`
	RefundedAt     time.Time    `gorm:"not null" json:"refundedAt"`
}
