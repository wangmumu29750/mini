package repository

import (
	"time"

	"mini-12306/backend/internal/model"

	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

type InventoryRow struct {
	ID              uint64
	TrainID         uint64
	TrainNo         string
	TrainType       string
	TravelDate      time.Time
	FromStationID   uint64
	FromStationName string
	ToStationID     uint64
	ToStationName   string
	SeatClassCode   string
	PriceCents      int64
	TotalCount      int
	AvailableCount  int
	LockedCount     int
	SoldCount       int
	Status          string
	UpdatedAt       time.Time
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *AdminRepository) DB() *gorm.DB {
	return r.db
}

func (r *AdminRepository) ListStations(page, pageSize int, status string) ([]model.Station, int64, int64, error) {
	query := r.db.Model(&model.Station{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	var activeTotal int64
	if err := r.db.Model(&model.Station{}).Where("status = ?", model.StationStatusActive).Count(&activeTotal).Error; err != nil {
		return nil, 0, 0, err
	}

	var stations []model.Station
	err := query.Order("id ASC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&stations).Error
	return stations, total, activeTotal, err
}

func (r *AdminRepository) FindStation(id uint64) (model.Station, error) {
	var station model.Station
	err := r.db.First(&station, id).Error
	return station, err
}

func (r *AdminRepository) CreateStation(station *model.Station) error {
	return r.db.Create(station).Error
}

func (r *AdminRepository) SaveStation(station *model.Station) error {
	return r.db.Save(station).Error
}

func (r *AdminRepository) ListTrains(page, pageSize int, status, trainNo string) ([]model.Train, int64, error) {
	query := r.db.Model(&model.Train{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if trainNo != "" {
		query = query.Where("train_no LIKE ?", "%"+trainNo+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var trains []model.Train
	err := query.Preload("Stops").Order("id DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&trains).Error
	return trains, total, err
}

func (r *AdminRepository) FindTrain(id uint64) (model.Train, error) {
	var train model.Train
	err := r.db.Preload("Stops").First(&train, id).Error
	return train, err
}

func (r *AdminRepository) CreateTrain(train *model.Train) error {
	return r.db.Create(train).Error
}

func (r *AdminRepository) SaveTrain(train *model.Train) error {
	return r.db.Save(train).Error
}

func (r *AdminRepository) ListStops(trainID uint64) ([]model.TrainStop, error) {
	var stops []model.TrainStop
	err := r.db.Preload("Station").Where("train_id = ?", trainID).Order("stop_order ASC").Find(&stops).Error
	return stops, err
}

func (r *AdminRepository) ListInventories(page, pageSize int, trainID uint64, seatClassCode string, travelDate *time.Time) ([]InventoryRow, int64, error) {
	query := r.db.Table("inventories AS i").Joins("JOIN trains AS t ON t.id = i.train_id").
		Joins("JOIN stations AS fs ON fs.id = i.from_station_id").
		Joins("JOIN stations AS ts ON ts.id = i.to_station_id")
	if trainID > 0 {
		query = query.Where("i.train_id = ?", trainID)
	}
	if seatClassCode != "" {
		query = query.Where("i.seat_class_code = ?", seatClassCode)
	}
	if travelDate != nil {
		query = query.Where("DATE(i.travel_date) = ?", travelDate.Format("2006-01-02"))
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []InventoryRow
	err := query.Select(`
			i.id,
			i.train_id,
			t.train_no,
			t.train_type,
			i.travel_date,
			i.from_station_id,
			fs.name AS from_station_name,
			i.to_station_id,
			ts.name AS to_station_name,
			i.seat_class_code,
			i.price_cents,
			i.total_count,
			i.available_count,
			i.locked_count,
			i.sold_count,
			i.status,
			i.updated_at
		`).
		Order("i.travel_date DESC, t.train_no ASC, i.seat_class_code ASC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&rows).Error
	return rows, total, err
}
