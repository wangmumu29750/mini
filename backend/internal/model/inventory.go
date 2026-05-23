package model

import "time"

type Inventory struct {
	BaseModel

	TrainID        uint64          `gorm:"not null;uniqueIndex:idx_inventory_key,priority:1;index" json:"trainId"`
	TravelDate     time.Time       `gorm:"type:date;not null;uniqueIndex:idx_inventory_key,priority:2" json:"travelDate"`
	FromStationID  uint64          `gorm:"not null;uniqueIndex:idx_inventory_key,priority:3" json:"fromStationId"`
	ToStationID    uint64          `gorm:"not null;uniqueIndex:idx_inventory_key,priority:4" json:"toStationId"`
	SeatClassCode  string          `gorm:"size:32;not null;uniqueIndex:idx_inventory_key,priority:5" json:"seatClassCode"`
	PriceCents     int64           `gorm:"not null" json:"priceCents"`
	TotalCount     int             `gorm:"not null" json:"totalCount"`
	AvailableCount int             `gorm:"not null" json:"availableCount"`
	LockedCount    int             `gorm:"not null" json:"lockedCount"`
	SoldCount      int             `gorm:"not null" json:"soldCount"`
	Status         InventoryStatus `gorm:"size:20;not null;default:ACTIVE;index" json:"status"`
}
