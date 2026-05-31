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

func (s *AdminService) ListTrains(query dto.AdminTrainListQuery) (dto.PageResponse[dto.AdminTrainResponse], error) {
	page, pageSize := normalizePage(query.Page, query.PageSize)
	status, err := trainStatus(query.Status, true)
	if err != nil {
		return dto.PageResponse[dto.AdminTrainResponse]{}, err
	}

	trains, total, err := s.admin.ListTrains(page, pageSize, status, strings.TrimSpace(query.TrainNo))
	if err != nil {
		return dto.PageResponse[dto.AdminTrainResponse]{}, err
	}
	items := make([]dto.AdminTrainResponse, 0, len(trains))
	for _, train := range trains {
		items = append(items, trainResponse(train))
	}
	return dto.PageResponse[dto.AdminTrainResponse]{Items: items, Page: page, PageSize: pageSize, Total: total}, nil
}

func (s *AdminService) CreateTrain(req dto.SaveTrainRequest) (dto.AdminTrainResponse, error) {
	train, err := trainFromRequest(req)
	if err != nil {
		return dto.AdminTrainResponse{}, err
	}
	if err := s.admin.CreateTrain(&train); err != nil {
		return dto.AdminTrainResponse{}, err
	}
	return trainResponse(train), nil
}

func (s *AdminService) UpdateTrain(id uint64, req dto.SaveTrainRequest) (dto.AdminTrainResponse, error) {
	train, err := s.admin.FindTrain(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.AdminTrainResponse{}, apperrors.New(http.StatusNotFound, response.CodeNotFound, "车次不存在")
		}
		return dto.AdminTrainResponse{}, err
	}

	next, err := trainFromRequest(req)
	if err != nil {
		return dto.AdminTrainResponse{}, err
	}
	train.TrainNo = next.TrainNo
	train.TrainType = next.TrainType
	train.Status = next.Status
	if err := s.admin.SaveTrain(&train); err != nil {
		return dto.AdminTrainResponse{}, err
	}
	return trainResponse(train), nil
}

func (s *AdminService) DeleteTrain(id uint64) (dto.AdminTrainResponse, error) {
	train, err := s.admin.FindTrain(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.AdminTrainResponse{}, apperrors.New(http.StatusNotFound, response.CodeNotFound, "车次不存在")
		}
		return dto.AdminTrainResponse{}, err
	}
	train.Status = model.TrainStatusDisabled
	if err := s.admin.SaveTrain(&train); err != nil {
		return dto.AdminTrainResponse{}, err
	}
	return trainResponse(train), nil
}

func (s *AdminService) SellableStats(query dto.SellableTrainStatsQuery) ([]dto.SellableTrainStatItem, error) {
	if query.FromStationID == query.ToStationID {
		return nil, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "出发站和到达站不能相同")
	}

	today := startOfLocalDay(time.Now())
	result := make([]dto.SellableTrainStatItem, 0, 2)
	for i := 0; i < 2; i++ {
		date := today.AddDate(0, 0, i)
		rows, err := repository.SearchAvailableTrains(s.admin.DB(), date, query.FromStationID, query.ToStationID)
		if err != nil {
			return nil, err
		}
		seen := map[uint64]struct{}{}
		for _, row := range rows {
			if row.AvailableCount > 0 {
				seen[row.TrainID] = struct{}{}
			}
		}
		result = append(result, dto.SellableTrainStatItem{Date: date.Format("2006-01-02"), TrainCount: int64(len(seen))})
	}
	return result, nil
}

func trainFromRequest(req dto.SaveTrainRequest) (model.Train, error) {
	status, err := trainStatus(req.Status, true)
	if err != nil {
		return model.Train{}, err
	}
	if status == "" {
		status = string(model.TrainStatusActive)
	}
	trainNo := strings.ToUpper(strings.TrimSpace(req.TrainNo))
	trainType := normalizeTrainType(req.TrainType)
	if trainNo == "" || trainType == "" {
		return model.Train{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "车次号和车次类型不能为空")
	}
	if inferred := trainTypeFromTrainNo(trainNo); inferred != "" && trainType != inferred {
		return model.Train{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "车次类型必须与车次号前缀一致")
	}
	if !supportedTrainType(trainType) {
		return model.Train{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "车次类型只支持 G/C/D/Z/T/K")
	}
	return model.Train{TrainNo: trainNo, TrainType: trainType, Status: model.TrainStatus(status)}, nil
}
