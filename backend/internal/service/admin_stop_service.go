package service

import (
	"errors"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/model"
	apperrors "mini-12306/backend/pkg/errors"
	"mini-12306/backend/pkg/response"

	"gorm.io/gorm"
)

func (s *AdminService) ListStops(trainID uint64) ([]dto.TrainStopResponse, error) {
	if _, err := s.admin.FindTrain(trainID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New(http.StatusNotFound, response.CodeNotFound, "车次不存在")
		}
		return nil, err
	}
	stops, err := s.admin.ListStops(trainID)
	if err != nil {
		return nil, err
	}
	items := make([]dto.TrainStopResponse, 0, len(stops))
	for _, stop := range stops {
		items = append(items, stopResponse(stop))
	}
	return items, nil
}

func (s *AdminService) SaveStops(trainID uint64, req dto.SaveTrainStopsRequest) ([]dto.TrainStopResponse, error) {
	if len(req.Stops) < 2 {
		return nil, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "经停站至少需要两个站点")
	}
	if err := validateStopItems(req.Stops); err != nil {
		return nil, err
	}

	err := s.admin.Transaction(func(tx *gorm.DB) error {
		var train model.Train
		if err := tx.First(&train, trainID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.New(http.StatusNotFound, response.CodeNotFound, "车次不存在")
			}
			return err
		}
		if err := tx.Unscoped().Where("train_id = ?", trainID).Delete(&model.TrainStop{}).Error; err != nil {
			return err
		}
		for _, item := range req.Stops {
			stop := model.TrainStop{
				TrainID:     trainID,
				StationID:   item.StationID,
				StopOrder:   item.StopOrder,
				DayOffset:   item.DayOffset,
				ArriveClock: strings.TrimSpace(item.ArriveClock),
				DepartClock: strings.TrimSpace(item.DepartClock),
				Mileage:     item.Mileage,
			}
			if err := tx.Create(&stop).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.ListStops(trainID)
}

func validateStopItems(stops []dto.SaveTrainStopItem) error {
	seenOrder := map[int]struct{}{}
	seenStation := map[uint64]struct{}{}
	clockPattern := regexp.MustCompile(`^$|^([01]\d|2[0-3]):[0-5]\d:[0-5]\d$`)
	for _, stop := range stops {
		if stop.StopOrder <= 0 || stop.StationID == 0 || stop.Mileage < 0 || stop.DayOffset < 0 {
			return apperrors.New(http.StatusBadRequest, response.CodeValidationError, "经停站序、站点、里程或跨天偏移不正确")
		}
		if _, ok := seenOrder[stop.StopOrder]; ok {
			return apperrors.New(http.StatusBadRequest, response.CodeValidationError, "经停站序不能重复")
		}
		if _, ok := seenStation[stop.StationID]; ok {
			return apperrors.New(http.StatusBadRequest, response.CodeValidationError, "经停站点不能重复")
		}
		if !clockPattern.MatchString(strings.TrimSpace(stop.ArriveClock)) || !clockPattern.MatchString(strings.TrimSpace(stop.DepartClock)) {
			return apperrors.New(http.StatusBadRequest, response.CodeValidationError, "到达和发车时间必须为 HH:mm:ss")
		}
		seenOrder[stop.StopOrder] = struct{}{}
		seenStation[stop.StationID] = struct{}{}
	}

	sort.Slice(stops, func(i, j int) bool { return stops[i].StopOrder < stops[j].StopOrder })
	for i := 1; i < len(stops); i++ {
		if stops[i].Mileage < stops[i-1].Mileage {
			return apperrors.New(http.StatusBadRequest, response.CodeValidationError, "经停里程必须随站序递增")
		}
	}
	return nil
}
