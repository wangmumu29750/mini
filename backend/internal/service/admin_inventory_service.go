package service

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/model"
	apperrors "mini-12306/backend/pkg/errors"
	"mini-12306/backend/pkg/response"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *AdminService) ListInventories(query dto.InventoryListQuery) (dto.PageResponse[dto.InventoryResponse], error) {
	page, pageSize := normalizePage(query.Page, query.PageSize)
	date, err := optionalDate(query.Date)
	if err != nil {
		return dto.PageResponse[dto.InventoryResponse]{}, err
	}

	rows, total, err := s.admin.ListInventories(page, pageSize, query.TrainID, strings.TrimSpace(query.SeatClassCode), date)
	if err != nil {
		return dto.PageResponse[dto.InventoryResponse]{}, err
	}

	items := make([]dto.InventoryResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, inventoryRowResponse(row))
	}
	return dto.PageResponse[dto.InventoryResponse]{Items: items, Page: page, PageSize: pageSize, Total: total}, nil
}

func (s *AdminService) SaveInventory(req dto.SaveInventoryRequest) (dto.InventoryResponse, error) {
	inventory, err := inventoryFromRequest(req)
	if err != nil {
		return dto.InventoryResponse{}, err
	}

	err = s.admin.Transaction(func(tx *gorm.DB) error {
		train, err := ensureTrainAndStations(tx, inventory.TrainID, inventory.FromStationID, inventory.ToStationID)
		if err != nil {
			return err
		}
		if err := ensureSeatClassAllowed(train.TrainType, inventory.SeatClassCode); err != nil {
			return err
		}
		calculatedPrice, err := fareForInventory(tx, inventory, train.TrainType)
		if err != nil {
			return err
		}
		if inventory.PriceCents <= 0 {
			inventory.PriceCents = calculatedPrice
		}
		var existing model.Inventory
		err = tx.Where("train_id = ? AND DATE(travel_date) = ? AND from_station_id = ? AND to_station_id = ? AND seat_class_code = ?",
			inventory.TrainID, inventory.TravelDate.Format("2006-01-02"), inventory.FromStationID, inventory.ToStationID, inventory.SeatClassCode).
			First(&existing).Error
		if err == nil {
			existing.PriceCents = inventory.PriceCents
			existing.TotalCount = inventory.TotalCount
			existing.AvailableCount = inventory.AvailableCount
			existing.LockedCount = inventory.LockedCount
			existing.SoldCount = inventory.SoldCount
			existing.Status = inventory.Status
			return tx.Save(&existing).Error
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return tx.Create(&inventory).Error
	})
	if err != nil {
		return dto.InventoryResponse{}, err
	}

	rows, _, err := s.admin.ListInventories(1, 1, req.TrainID, req.SeatClassCode, &inventory.TravelDate)
	if err != nil {
		return dto.InventoryResponse{}, err
	}
	if len(rows) == 0 {
		return dto.InventoryResponse{}, apperrors.New(http.StatusNotFound, response.CodeNotFound, "票额不存在")
	}
	return inventoryRowResponse(rows[0]), nil
}

func (s *AdminService) QuoteStats(query dto.InventoryQuoteStatsQuery) (dto.InventoryQuoteStatsResponse, error) {
	db := s.admin.DB().Model(&model.Inventory{}).Where("train_id = ?", query.TrainID)
	if query.SeatClassCode != "" {
		db = db.Where("seat_class_code = ?", query.SeatClassCode)
	}

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return dto.InventoryQuoteStatsResponse{}, err
	}

	var lowest int64
	if count > 0 {
		if err := db.Select("COALESCE(MIN(price_cents), 0)").Scan(&lowest).Error; err != nil {
			return dto.InventoryQuoteStatsResponse{}, err
		}
	}
	return dto.InventoryQuoteStatsResponse{
		TrainID:       query.TrainID,
		SeatClassCode: query.SeatClassCode,
		QuoteCount:    count,
		LowestPrice:   lowest,
	}, nil
}

func (s *AdminService) FlowInventory(req dto.InventoryFlowRequest) (dto.InventoryFlowResponse, error) {
	if req.Quantity <= 0 {
		return dto.InventoryFlowResponse{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "流转数量必须大于 0")
	}
	action := strings.ToUpper(strings.TrimSpace(req.Action))

	var trainID uint64
	var updated model.Inventory
	err := s.admin.Transaction(func(tx *gorm.DB) error {
		var inventory model.Inventory
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&inventory, req.InventoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.New(http.StatusNotFound, response.CodeNotFound, "票额不存在")
			}
			return err
		}
		trainID = inventory.TrainID

		switch action {
		case "LOCK":
			if inventory.AvailableCount < req.Quantity {
				return apperrors.New(http.StatusConflict, response.CodeInsufficientInventory, "可售票额不足")
			}
			inventory.AvailableCount -= req.Quantity
			inventory.LockedCount += req.Quantity
		case "PAY":
			if inventory.LockedCount < req.Quantity {
				return apperrors.New(http.StatusConflict, response.CodeInvalidOrderState, "锁定票额不足")
			}
			inventory.LockedCount -= req.Quantity
			inventory.SoldCount += req.Quantity
		case "RELEASE":
			if inventory.LockedCount < req.Quantity {
				return apperrors.New(http.StatusConflict, response.CodeInvalidOrderState, "锁定票额不足")
			}
			inventory.LockedCount -= req.Quantity
			inventory.AvailableCount += req.Quantity
		case "REFUND", "CHANGE_OUT":
			if inventory.SoldCount < req.Quantity {
				return apperrors.New(http.StatusConflict, response.CodeInvalidOrderState, "已售票额不足")
			}
			inventory.SoldCount -= req.Quantity
			inventory.AvailableCount += req.Quantity
		case "CHANGE_IN":
			if inventory.AvailableCount < req.Quantity {
				return apperrors.New(http.StatusConflict, response.CodeInsufficientInventory, "可售票额不足")
			}
			inventory.AvailableCount -= req.Quantity
			inventory.SoldCount += req.Quantity
		default:
			return apperrors.New(http.StatusBadRequest, response.CodeValidationError, "不支持的票额流转动作")
		}

		if err := tx.Save(&inventory).Error; err != nil {
			return err
		}
		updated = inventory
		return nil
	})
	if err != nil {
		return dto.InventoryFlowResponse{}, err
	}

	rows, _, err := s.admin.ListInventories(1, 1, trainID, updated.SeatClassCode, &updated.TravelDate)
	if err != nil {
		return dto.InventoryFlowResponse{}, err
	}
	var item dto.InventoryResponse
	for _, row := range rows {
		if row.ID == updated.ID {
			item = inventoryRowResponse(row)
			break
		}
	}
	if item.ID == 0 {
		item = dto.InventoryResponse{
			ID:             updated.ID,
			TrainID:        updated.TrainID,
			TravelDate:     updated.TravelDate.Format("2006-01-02"),
			SeatClassCode:  updated.SeatClassCode,
			SeatClassName:  seatClassName(updated.SeatClassCode),
			PriceCents:     updated.PriceCents,
			TotalCount:     updated.TotalCount,
			AvailableCount: updated.AvailableCount,
			LockedCount:    updated.LockedCount,
			SoldCount:      updated.SoldCount,
			Status:         string(updated.Status),
			UpdatedAt:      updated.UpdatedAt.Format(time.RFC3339),
		}
	}

	lowest, err := lowestPriceForTrain(s.admin.DB(), trainID)
	if err != nil {
		return dto.InventoryFlowResponse{}, err
	}
	return dto.InventoryFlowResponse{Inventory: item, LowestPrice: lowest}, nil
}

func inventoryFromRequest(req dto.SaveInventoryRequest) (model.Inventory, error) {
	date, err := time.ParseInLocation("2006-01-02", req.TravelDate, time.Local)
	if err != nil {
		return model.Inventory{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "乘车日期格式不正确")
	}
	if req.FromStationID == req.ToStationID {
		return model.Inventory{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "出发站和到达站不能相同")
	}
	if req.PriceCents < 0 || req.TotalCount < 0 || req.AvailableCount < 0 || req.LockedCount < 0 || req.SoldCount < 0 {
		return model.Inventory{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "票价和票额数量不正确")
	}
	if req.AvailableCount+req.LockedCount+req.SoldCount > req.TotalCount {
		return model.Inventory{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "可售、锁定和已售数量不能超过总票额")
	}
	status := strings.ToUpper(strings.TrimSpace(req.Status))
	if status == "" {
		status = string(model.InventoryStatusActive)
	}
	if status != string(model.InventoryStatusActive) {
		return model.Inventory{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "票额状态只能是 ACTIVE")
	}
	return model.Inventory{
		TrainID:        req.TrainID,
		TravelDate:     startOfLocalDay(date),
		FromStationID:  req.FromStationID,
		ToStationID:    req.ToStationID,
		SeatClassCode:  strings.ToUpper(strings.TrimSpace(req.SeatClassCode)),
		PriceCents:     req.PriceCents,
		TotalCount:     req.TotalCount,
		AvailableCount: req.AvailableCount,
		LockedCount:    req.LockedCount,
		SoldCount:      req.SoldCount,
		Status:         model.InventoryStatus(status),
	}, nil
}

func fareForInventory(tx *gorm.DB, inventory model.Inventory, trainType string) (int64, error) {
	var fromStop model.TrainStop
	if err := tx.Where("train_id = ? AND station_id = ?", inventory.TrainID, inventory.FromStationID).First(&fromStop).Error; err != nil {
		return 0, err
	}
	var toStop model.TrainStop
	if err := tx.Where("train_id = ? AND station_id = ?", inventory.TrainID, inventory.ToStationID).First(&toStop).Error; err != nil {
		return 0, err
	}
	if fromStop.StopOrder >= toStop.StopOrder {
		return 0, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "出发站必须早于到达站")
	}
	return calculateFareCents(toStop.Mileage-fromStop.Mileage, trainType, inventory.SeatClassCode)
}
