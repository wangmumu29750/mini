package model

import "time"

type ChangeRecord struct {
	BaseModel

	ChangeNo       string       `gorm:"size:32;not null;uniqueIndex" json:"changeNo"`
	OldTicketID    uint64       `gorm:"not null;index" json:"oldTicketId"`
	NewTicketID    uint64       `gorm:"not null;index" json:"newTicketId"`
	UserID         uint64       `gorm:"not null;index;uniqueIndex:idx_change_idempotency,priority:1" json:"userId"`
	PriceDiffCents int64        `gorm:"not null" json:"priceDiffCents"`
	Status         ChangeStatus `gorm:"size:20;not null;index" json:"status"`
	IdempotencyKey string       `gorm:"size:80;uniqueIndex:idx_change_idempotency,priority:2" json:"-"`
	ChangedAt      time.Time    `gorm:"not null" json:"changedAt"`
}
