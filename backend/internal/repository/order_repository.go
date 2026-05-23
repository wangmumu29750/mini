package repository

import (
	"errors"
	"time"

	"mini-12306/backend/internal/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

type OrderInventoryRow struct {
	InventoryID     uint64
	TrainID         uint64
	TrainNo         string
	TravelDate      time.Time
	FromStationID   uint64
	FromStationName string
	ToStationID     uint64
	ToStationName   string
	SeatClassCode   string
	PriceCents      int64
	AvailableCount  int
	LockedCount     int
	SoldCount       int
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *OrderRepository) ListByUser(userID uint64) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("Tickets").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) FindByUserAndID(userID, orderID uint64) (model.Order, error) {
	var order model.Order
	err := r.db.Preload("Tickets").
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error
	return order, err
}

func FindInventoryForOrder(tx *gorm.DB, travelDate time.Time, trainID, fromStationID, toStationID uint64, seatClassCode string) (OrderInventoryRow, error) {
	var row OrderInventoryRow
	result := tx.Table("inventories AS i").
		Select(`
			i.id AS inventory_id,
			i.train_id,
			t.train_no,
			i.travel_date,
			i.from_station_id,
			fs.name AS from_station_name,
			i.to_station_id,
			ts.name AS to_station_name,
			i.seat_class_code,
			i.price_cents,
			i.available_count,
			i.locked_count,
			i.sold_count
		`).
		Joins("JOIN trains AS t ON t.id = i.train_id AND t.status = ?", model.TrainStatusActive).
		Joins("JOIN stations AS fs ON fs.id = i.from_station_id AND fs.status = ?", model.StationStatusActive).
		Joins("JOIN stations AS ts ON ts.id = i.to_station_id AND ts.status = ?", model.StationStatusActive).
		Where("DATE(i.travel_date) = ? AND i.train_id = ? AND i.from_station_id = ? AND i.to_station_id = ? AND i.seat_class_code = ? AND i.status = ?",
			travelDate.Format("2006-01-02"), trainID, fromStationID, toStationID, seatClassCode, model.InventoryStatusActive).
		Limit(1).
		Scan(&row)
	if result.Error != nil {
		return OrderInventoryRow{}, result.Error
	}
	if result.RowsAffected == 0 {
		return OrderInventoryRow{}, gorm.ErrRecordNotFound
	}
	if row.InventoryID == 0 {
		return OrderInventoryRow{}, errors.New("inventory row scan failed")
	}
	return row, nil
}
