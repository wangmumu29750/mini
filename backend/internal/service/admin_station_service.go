package service

import (
	"errors"
	"net/http"
	"strings"

	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/model"
	apperrors "mini-12306/backend/pkg/errors"
	"mini-12306/backend/pkg/response"

	"gorm.io/gorm"
)

func (s *AdminService) ListStations(query dto.StationListQuery) (dto.StationListResponse, error) {
	page, pageSize := normalizePage(query.Page, query.PageSize)
	status, err := stationStatus(query.Status, true)
	if err != nil {
		return dto.StationListResponse{}, err
	}

	stations, total, activeTotal, err := s.admin.ListStations(page, pageSize, status)
	if err != nil {
		return dto.StationListResponse{}, err
	}

	items := make([]dto.AdminStationResponse, 0, len(stations))
	for _, station := range stations {
		items = append(items, stationResponse(station))
	}
	return dto.StationListResponse{Items: items, Page: page, PageSize: pageSize, Total: total, ActiveTotal: activeTotal}, nil
}

func (s *AdminService) PublicStations(query dto.StationListQuery) (any, error) {
	page, pageSize := normalizePage(query.Page, query.PageSize)
	status, err := stationStatus(query.Status, true)
	if err != nil {
		return nil, err
	}
	if status == "" && query.Status == "" {
		status = string(model.StationStatusActive)
	}

	stations, total, _, err := s.admin.ListStations(page, pageSize, status)
	if err != nil {
		return nil, err
	}

	items := make([]dto.StationResponse, 0, len(stations))
	for _, station := range stations {
		items = append(items, dto.StationResponse{ID: station.ID, Name: station.Name})
	}
	if query.Page <= 0 && query.PageSize <= 0 && query.Status == "" {
		return items, nil
	}
	return dto.PageResponse[dto.StationResponse]{Items: items, Page: page, PageSize: pageSize, Total: total}, nil
}

func (s *AdminService) CreateStation(req dto.SaveStationRequest) (dto.AdminStationResponse, error) {
	station, err := stationFromRequest(req)
	if err != nil {
		return dto.AdminStationResponse{}, err
	}
	if err := s.admin.CreateStation(&station); err != nil {
		return dto.AdminStationResponse{}, err
	}
	return stationResponse(station), nil
}

func (s *AdminService) UpdateStation(id uint64, req dto.SaveStationRequest) (dto.AdminStationResponse, error) {
	station, err := s.admin.FindStation(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.AdminStationResponse{}, apperrors.New(http.StatusNotFound, response.CodeNotFound, "站点不存在")
		}
		return dto.AdminStationResponse{}, err
	}

	next, err := stationFromRequest(req)
	if err != nil {
		return dto.AdminStationResponse{}, err
	}
	station.Code = next.Code
	station.Name = next.Name
	station.City = next.City
	station.Status = next.Status
	if err := s.admin.SaveStation(&station); err != nil {
		return dto.AdminStationResponse{}, err
	}
	return stationResponse(station), nil
}

func (s *AdminService) DisableStation(id uint64) (dto.AdminStationResponse, error) {
	station, err := s.admin.FindStation(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.AdminStationResponse{}, apperrors.New(http.StatusNotFound, response.CodeNotFound, "站点不存在")
		}
		return dto.AdminStationResponse{}, err
	}
	station.Status = model.StationStatusDisabled
	if err := s.admin.SaveStation(&station); err != nil {
		return dto.AdminStationResponse{}, err
	}
	return stationResponse(station), nil
}

func stationFromRequest(req dto.SaveStationRequest) (model.Station, error) {
	status, err := stationStatus(req.Status, true)
	if err != nil {
		return model.Station{}, err
	}
	if status == "" {
		status = string(model.StationStatusActive)
	}
	code := strings.ToUpper(strings.TrimSpace(req.Code))
	name := strings.TrimSpace(req.Name)
	city := strings.TrimSpace(req.City)
	if code == "" || name == "" || city == "" {
		return model.Station{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "站点代码、名称和城市不能为空")
	}
	return model.Station{Code: code, Name: name, City: city, Status: model.StationStatus(status)}, nil
}
