package service

import (
	"net/http"
	"time"

	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/repository"
	apperrors "mini-12306/backend/pkg/errors"
	"mini-12306/backend/pkg/response"

	"gorm.io/gorm"
)

type TrainService struct {
	trains *repository.TrainRepository
}

func NewTrainService(trains *repository.TrainRepository) *TrainService {
	return &TrainService{trains: trains}
}

func (s *TrainService) ListStations() ([]dto.StationResponse, error) {
	stations, err := s.trains.ListActiveStations()
	if err != nil {
		return nil, err
	}

	result := make([]dto.StationResponse, 0, len(stations))
	for _, station := range stations {
		result = append(result, dto.StationResponse{
			ID:   station.ID,
			Name: station.Name,
		})
	}
	return result, nil
}

func (s *TrainService) Search(query dto.TrainSearchQuery) ([]dto.TrainSearchItemResponse, error) {
	if query.FromStationID == query.ToStationID {
		return nil, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "出发站和到达站不能相同")
	}

	date, err := time.ParseInLocation("2006-01-02", query.Date, time.Local)
	if err != nil {
		return nil, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "乘车日期格式不正确")
	}
	if err := validateTicketTravelDate(date); err != nil {
		return nil, err
	}

	result := make([]dto.TrainSearchItemResponse, 0)
	err = s.trains.Transaction(func(tx *gorm.DB) error {
		rows, err := repository.SearchAvailableTrains(tx, date, query.FromStationID, query.ToStationID)
		if err != nil {
			return err
		}

		items, err := buildTrainSearchResponses(tx, date, rows)
		if err != nil {
			return err
		}

		now := time.Now()
		for _, item := range items {
			departTime, err := time.Parse(time.RFC3339, item.DepartTime)
			if err != nil {
				return err
			}
			if departTime.After(now) {
				result = append(result, item)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
