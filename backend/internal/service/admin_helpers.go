package service

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/model"
	"mini-12306/backend/internal/repository"
	apperrors "mini-12306/backend/pkg/errors"
	"mini-12306/backend/pkg/response"

	"gorm.io/gorm"
)

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func stationStatus(value string, allowEmpty bool) (string, error) {
	status := strings.ToUpper(strings.TrimSpace(value))
	if status == "" && allowEmpty {
		return "", nil
	}
	if status != string(model.StationStatusActive) && status != string(model.StationStatusDisabled) {
		return "", apperrors.New(http.StatusBadRequest, response.CodeValidationError, "站点状态只能是 ACTIVE 或 DISABLED")
	}
	return status, nil
}

func trainStatus(value string, allowEmpty bool) (string, error) {
	status := strings.ToUpper(strings.TrimSpace(value))
	if status == "" && allowEmpty {
		return "", nil
	}
	if status != string(model.TrainStatusActive) && status != string(model.TrainStatusDisabled) {
		return "", apperrors.New(http.StatusBadRequest, response.CodeValidationError, "车次状态只能是 ACTIVE 或 DISABLED")
	}
	return status, nil
}

func optionalDate(value string) (*time.Time, error) {
	text := strings.TrimSpace(value)
	if text == "" {
		return nil, nil
	}
	date, err := time.ParseInLocation("2006-01-02", text, time.Local)
	if err != nil {
		return nil, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "日期格式不正确")
	}
	date = startOfLocalDay(date)
	return &date, nil
}

func ensureTrainAndStations(tx *gorm.DB, trainID, fromStationID, toStationID uint64) error {
	var train model.Train
	if err := tx.First(&train, trainID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.New(http.StatusNotFound, response.CodeNotFound, "车次不存在")
		}
		return err
	}

	for _, stationID := range []uint64{fromStationID, toStationID} {
		var station model.Station
		if err := tx.First(&station, stationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.New(http.StatusNotFound, response.CodeNotFound, "站点不存在")
			}
			return err
		}
	}
	return nil
}

func lowestPriceForTrain(db *gorm.DB, trainID uint64) (int64, error) {
	var lowest int64
	err := db.Model(&model.Inventory{}).
		Where("train_id = ? AND status = ?", trainID, model.InventoryStatusActive).
		Select("COALESCE(MIN(price_cents), 0)").
		Scan(&lowest).Error
	return lowest, err
}

func startOfLocalDay(value time.Time) time.Time {
	date, _ := time.ParseInLocation("2006-01-02", value.Format("2006-01-02"), time.Local)
	return date
}

func stationResponse(station model.Station) dto.AdminStationResponse {
	return dto.AdminStationResponse{
		ID:        station.ID,
		Code:      station.Code,
		Name:      station.Name,
		City:      station.City,
		Status:    string(station.Status),
		CreatedAt: station.CreatedAt.Format(time.RFC3339),
		UpdatedAt: station.UpdatedAt.Format(time.RFC3339),
	}
}

func trainResponse(train model.Train) dto.AdminTrainResponse {
	return dto.AdminTrainResponse{
		ID:        train.ID,
		TrainNo:   train.TrainNo,
		TrainType: train.TrainType,
		Status:    string(train.Status),
		StopCount: int64(len(train.Stops)),
		CreatedAt: train.CreatedAt.Format(time.RFC3339),
		UpdatedAt: train.UpdatedAt.Format(time.RFC3339),
	}
}

func stopResponse(stop model.TrainStop) dto.TrainStopResponse {
	return dto.TrainStopResponse{
		ID:          stop.ID,
		TrainID:     stop.TrainID,
		Station:     dto.StationResponse{ID: stop.StationID, Name: stop.Station.Name},
		StopOrder:   stop.StopOrder,
		DayOffset:   stop.DayOffset,
		ArriveClock: stop.ArriveClock,
		DepartClock: stop.DepartClock,
		Mileage:     stop.Mileage,
	}
}

func inventoryRowResponse(row repository.InventoryRow) dto.InventoryResponse {
	return dto.InventoryResponse{
		ID:             row.ID,
		TrainID:        row.TrainID,
		TrainNo:        row.TrainNo,
		TravelDate:     row.TravelDate.Format("2006-01-02"),
		FromStation:    dto.StationResponse{ID: row.FromStationID, Name: row.FromStationName},
		ToStation:      dto.StationResponse{ID: row.ToStationID, Name: row.ToStationName},
		SeatClassCode:  row.SeatClassCode,
		SeatClassName:  seatClassName(row.SeatClassCode),
		PriceCents:     row.PriceCents,
		TotalCount:     row.TotalCount,
		AvailableCount: row.AvailableCount,
		LockedCount:    row.LockedCount,
		SoldCount:      row.SoldCount,
		Status:         row.Status,
		UpdatedAt:      row.UpdatedAt.Format(time.RFC3339),
	}
}
