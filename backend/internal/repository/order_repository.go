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
	TrainType       string
	TravelDate      time.Time
	FromStationID   uint64
	FromStationName string
	ToStationID     uint64
	ToStationName   string
	DepartClock     string
	DepartDayOffset int
	ArriveClock     string
	ArriveDayOffset int
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
	err := r.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("id ASC")
	}).Preload("Tickets", func(db *gorm.DB) *gorm.DB {
		return db.Order("id DESC")
	}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	if err := attachOrderTimes(r.db, orders); err != nil {
		return nil, err
	}
	return orders, err
}

func (r *OrderRepository) FindByUserAndID(userID, orderID uint64) (model.Order, error) {
	var order model.Order
	err := r.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("id ASC")
	}).Preload("Tickets", func(db *gorm.DB) *gorm.DB {
		return db.Order("id DESC")
	}).
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error
	if err != nil {
		return order, err
	}
	orders := []model.Order{order}
	if err := attachOrderTimes(r.db, orders); err != nil {
		return order, err
	}
	order = orders[0]
	return order, err
}

func (r *OrderRepository) ExistsActiveTripForPassenger(
	db *gorm.DB,
	passengerID uint64,
	travelDate time.Time,
	trainID uint64,
) (bool, error) {
	var count int64
	query := db
	if query == nil {
		query = r.db
	}
	err := query.Table("order_items AS oi").
		Joins("JOIN orders AS o ON o.id = oi.order_id").
		Where(
			`oi.passenger_id = ? AND DATE(o.travel_date) = ? AND o.train_id = ? AND o.status IN ?`,
			passengerID,
			travelDate.Format("2006-01-02"),
			trainID,
			[]model.OrderStatus{model.OrderStatusPendingPayment, model.OrderStatusPaid},
		).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func attachOrderTimes(db *gorm.DB, orders []model.Order) error {
	for i := range orders {
		depart, arrive, err := stationTimes(db, orders[i].TravelDate, orders[i].TrainID, orders[i].FromStationID, orders[i].ToStationID)
		if err != nil {
			return err
		}
		orders[i].DepartTime = &depart
		orders[i].ArriveTime = &arrive
		for j := range orders[i].Tickets {
			orders[i].Tickets[j].DepartTime = &depart
			orders[i].Tickets[j].ArriveTime = &arrive
		}
	}
	return nil
}

func stationTimes(db *gorm.DB, travelDate time.Time, trainID, fromStationID, toStationID uint64) (time.Time, time.Time, error) {
	var fromStop model.TrainStop
	if err := db.Where("train_id = ? AND station_id = ?", trainID, fromStationID).First(&fromStop).Error; err != nil {
		return time.Time{}, time.Time{}, err
	}
	var toStop model.TrainStop
	if err := db.Where("train_id = ? AND station_id = ?", trainID, toStationID).First(&toStop).Error; err != nil {
		return time.Time{}, time.Time{}, err
	}
	return combineDateClock(travelDate, fromStop.DepartClock, fromStop.DayOffset),
		combineDateClock(travelDate, toStop.ArriveClock, toStop.DayOffset),
		nil
}

func combineDateClock(travelDate time.Time, clock string, dayOffset int) time.Time {
	if clock == "" {
		clock = "00:00:00"
	}
	parsed, err := time.ParseInLocation("2006-01-02 15:04:05", travelDate.Format("2006-01-02")+" "+clock, time.Local)
	if err != nil {
		return travelDate
	}
	return parsed.AddDate(0, 0, dayOffset)
}

func FindInventoryForOrder(tx *gorm.DB, travelDate time.Time, trainID, fromStationID, toStationID uint64, seatClassCode string) (OrderInventoryRow, error) {
	var row OrderInventoryRow
	result := tx.Table("inventories AS i").
		Select(`
			i.id AS inventory_id,
			i.train_id,
			t.train_no,
			t.train_type,
			i.travel_date,
			i.from_station_id,
			fs.name AS from_station_name,
			i.to_station_id,
			ts.name AS to_station_name,
			fstop.depart_clock,
			fstop.day_offset AS depart_day_offset,
			tstop.arrive_clock,
			tstop.day_offset AS arrive_day_offset,
			i.seat_class_code,
			i.price_cents,
			i.available_count,
			i.locked_count,
			i.sold_count
		`).
		Joins("JOIN trains AS t ON t.id = i.train_id AND t.status = ?", model.TrainStatusActive).
		Joins("JOIN stations AS fs ON fs.id = i.from_station_id AND fs.status = ?", model.StationStatusActive).
		Joins("JOIN stations AS ts ON ts.id = i.to_station_id AND ts.status = ?", model.StationStatusActive).
		Joins("JOIN train_stops AS fstop ON fstop.train_id = i.train_id AND fstop.station_id = i.from_station_id").
		Joins("JOIN train_stops AS tstop ON tstop.train_id = i.train_id AND tstop.station_id = i.to_station_id").
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
