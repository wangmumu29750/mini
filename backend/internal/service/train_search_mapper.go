package service

import (
	"fmt"
	"sort"
	"time"

	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/repository"

	"gorm.io/gorm"
)

func buildTrainSearchResponses(db *gorm.DB, date time.Time, rows []repository.TrainSearchRow) ([]dto.TrainSearchItemResponse, error) {
	itemsByTrain := make(map[uint64]*dto.TrainSearchItemResponse)
	viaStationsBySegment := make(map[string][]dto.StationResponse)

	for _, row := range rows {
		key := viaStationKey(row.TrainID, row.FromOrder, row.ToOrder)
		if _, ok := viaStationsBySegment[key]; !ok {
			stations, err := repository.ListRouteStations(db, row.TrainID, row.FromOrder, row.ToOrder)
			if err != nil {
				return nil, err
			}
			viaStations := make([]dto.StationResponse, 0, len(stations))
			for _, station := range stations {
				viaStations = append(viaStations, dto.StationResponse{
					ID:   station.ID,
					Name: station.Name,
				})
			}
			viaStationsBySegment[key] = viaStations
		}

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
				TrainType:  row.TrainType,
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
				ViaStations:     viaStationsBySegment[key],
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

func viaStationKey(trainID uint64, fromOrder, toOrder int) string {
	return fmt.Sprintf("%d-%d-%d", trainID, fromOrder, toOrder)
}

func trainSearchRowsToResponses(date time.Time, rows []repository.TrainSearchRow) []dto.TrainSearchItemResponse {
	items := make([]dto.TrainSearchItemResponse, 0)
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
				TrainType:  row.TrainType,
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

	for _, item := range itemsByTrain {
		items = append(items, *item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].DepartTime < items[j].DepartTime
	})
	return items
}
