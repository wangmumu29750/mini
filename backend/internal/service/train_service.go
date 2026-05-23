package service

import (
	"net/http"
	"sort"
	"time"

	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/repository"
	apperrors "mini-12306/backend/pkg/errors"
	"mini-12306/backend/pkg/response"
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

	rows, err := s.trains.SearchAvailableTrains(date, query.FromStationID, query.ToStationID)
	if err != nil {
		return nil, err
	}

	itemsByTrain := make(map[uint64]*dto.TrainSearchItemResponse)
	for _, row := range rows {
		item, ok := itemsByTrain[row.TrainID]
		if !ok {
			departTime := combineDateClock(date, row.DepartClock, row.DepartDayOffset)
			arriveTime := combineDateClock(date, row.ArriveClock, row.ArriveDayOffset)
			duration := int(arriveTime.Sub(departTime).Minutes())
			if duration < 0 {
				duration = 0
			}

			item = &dto.TrainSearchItemResponse{
				TrainID:    row.TrainID,
				TrainNo:    row.TrainNo,
				TravelDate: date.Format("2006-01-02"),
				FromStation: dto.StationResponse{
					ID:   row.FromStationID,
					Name: row.FromStationName,
				},
				ToStation: dto.StationResponse{
					ID:   row.ToStationID,
					Name: row.ToStationName,
				},
				DepartTime:      departTime.Format(time.RFC3339),
				ArriveTime:      arriveTime.Format(time.RFC3339),
				DurationMinutes: duration,
				SeatOptions:     []dto.SeatOptionResponse{},
			}
			itemsByTrain[row.TrainID] = item
		}

		item.SeatOptions = append(item.SeatOptions, dto.SeatOptionResponse{
			SeatClassCode:  row.SeatClassCode,
			SeatClassName:  seatClassName(row.SeatClassCode),
			PriceCents:     row.PriceCents,
			AvailableCount: row.AvailableCount,
		})
	}

	result := make([]dto.TrainSearchItemResponse, 0, len(itemsByTrain))
	for _, item := range itemsByTrain {
		result = append(result, *item)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].DepartTime < result[j].DepartTime
	})
	return result, nil
}

func combineDateClock(date time.Time, clock string, dayOffset int) time.Time {
	if clock == "" {
		clock = "00:00:00"
	}
	parsed, err := time.ParseInLocation("2006-01-02 15:04:05", date.Format("2006-01-02")+" "+clock, time.Local)
	if err != nil {
		return date
	}
	return parsed.AddDate(0, 0, dayOffset)
}

func seatClassName(code string) string {
	switch code {
	case "BUSINESS":
		return "商务座"
	case "FIRST":
		return "一等座"
	case "SECOND":
		return "二等座"
	case "HARD_SEAT":
		return "硬座"
	case "HARD_SLEEPER":
		return "硬卧"
	default:
		return code
	}
}
