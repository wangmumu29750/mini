package repository

import (
	"time"

	"mini-12306/backend/internal/model"

	"gorm.io/gorm"
)

type TrainRepository struct {
	db *gorm.DB
}

type TrainSearchRow struct {
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
	FromOrder       int
	ToOrder         int
	SeatClassCode   string
	PriceCents      int64
	AvailableCount  int
}

func NewTrainRepository(db *gorm.DB) *TrainRepository {
	return &TrainRepository{db: db}
}

func (r *TrainRepository) ListActiveStations() ([]model.Station, error) {
	var stations []model.Station
	err := r.db.Where("status = ?", model.StationStatusActive).Order("id ASC").Find(&stations).Error
	return stations, err
}

func (r *TrainRepository) SearchAvailableTrains(date time.Time, fromStationID, toStationID uint64) ([]TrainSearchRow, error) {
	return SearchAvailableTrains(r.db, date, fromStationID, toStationID)
}

func SearchAvailableTrains(db *gorm.DB, date time.Time, fromStationID, toStationID uint64) ([]TrainSearchRow, error) {
	var rows []TrainSearchRow
	err := db.Table("inventories AS i").
		Select(`
			t.id AS train_id,
			t.train_no,
			t.train_type,
			i.travel_date,
			fs.id AS from_station_id,
			fs.name AS from_station_name,
			ts.id AS to_station_id,
			ts.name AS to_station_name,
			fstop.depart_clock,
			fstop.day_offset AS depart_day_offset,
			tstop.arrive_clock,
			tstop.day_offset AS arrive_day_offset,
			fstop.stop_order AS from_order,
			tstop.stop_order AS to_order,
			i.seat_class_code,
			i.price_cents,
			i.available_count
		`).
		Joins("JOIN trains AS t ON t.id = i.train_id AND t.status = ?", model.TrainStatusActive).
		Joins("JOIN stations AS fs ON fs.id = i.from_station_id AND fs.status = ?", model.StationStatusActive).
		Joins("JOIN stations AS ts ON ts.id = i.to_station_id AND ts.status = ?", model.StationStatusActive).
		Joins("JOIN train_stops AS fstop ON fstop.train_id = i.train_id AND fstop.station_id = i.from_station_id").
		Joins("JOIN train_stops AS tstop ON tstop.train_id = i.train_id AND tstop.station_id = i.to_station_id").
		Where("DATE(i.travel_date) = ? AND i.from_station_id = ? AND i.to_station_id = ? AND i.status = ?", date.Format("2006-01-02"), fromStationID, toStationID, model.InventoryStatusActive).
		Where("fstop.stop_order < tstop.stop_order").
		Order("fstop.depart_clock ASC, t.train_no ASC, i.seat_class_code ASC").
		Scan(&rows).Error
	return rows, err
}
