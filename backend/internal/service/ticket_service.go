package service

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/model"
	"mini-12306/backend/internal/repository"
	apperrors "mini-12306/backend/pkg/errors"
	"mini-12306/backend/pkg/response"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TicketService struct {
	tickets *repository.TicketRepository
}

func NewTicketService(tickets *repository.TicketRepository) *TicketService {
	return &TicketService{tickets: tickets}
}

func (s *TicketService) List(userID uint64) ([]dto.TicketResponse, error) {
	tickets, err := s.tickets.ListByUser(userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.TicketResponse, 0, len(tickets))
	for _, ticket := range tickets {
		result = append(result, ticketResponse(ticket))
	}
	return result, nil
}

func (s *TicketService) Detail(userID, ticketID uint64) (dto.TicketResponse, error) {
	ticket, err := s.tickets.FindByUserAndID(userID, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.TicketResponse{}, apperrors.New(http.StatusNotFound, response.CodeNotFound, "车票不存在")
		}
		return dto.TicketResponse{}, err
	}
	return ticketResponse(ticket), nil
}

func (s *TicketService) ChangeOptions(userID, ticketID uint64, query dto.ChangeOptionsQuery) (dto.ChangeOptionsResponse, error) {
	date, err := time.ParseInLocation("2006-01-02", query.Date, time.Local)
	if err != nil {
		return dto.ChangeOptionsResponse{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "改签日期格式不正确")
	}

	ticket, err := s.tickets.FindByUserAndID(userID, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ChangeOptionsResponse{}, apperrors.New(http.StatusNotFound, response.CodeNotFound, "车票不存在")
		}
		return dto.ChangeOptionsResponse{}, err
	}
	if ticket.Status != model.TicketStatusIssued {
		return dto.ChangeOptionsResponse{}, apperrors.New(http.StatusConflict, response.CodeTicketNotChangeable, "当前车票不可改签")
	}
	if ticket.DepartTime != nil && !ticket.DepartTime.After(time.Now()) {
		return dto.ChangeOptionsResponse{}, apperrors.New(http.StatusConflict, response.CodeTicketNotChangeable, "已发车车票不可改签")
	}

	var rows []repository.TrainSearchRow
	err = s.tickets.Transaction(func(tx *gorm.DB) error {
		var loadErr error
		rows, loadErr = repository.SearchAvailableTrains(tx, date, ticket.FromStationID, ticket.ToStationID)
		return loadErr
	})
	if err != nil {
		return dto.ChangeOptionsResponse{}, err
	}

	options := trainSearchRowsToResponses(date, rows)
	filtered := make([]dto.TrainSearchItemResponse, 0, len(options))
	now := time.Now()
	for _, option := range options {
		departTime, err := time.Parse(time.RFC3339, option.DepartTime)
		if err != nil {
			return dto.ChangeOptionsResponse{}, err
		}
		if !departTime.After(now) {
			continue
		}

		seatOptions := option.SeatOptions[:0]
		for _, seat := range option.SeatOptions {
			if option.TrainID == ticket.TrainID && option.TravelDate == ticket.TravelDate.Format("2006-01-02") && seat.SeatClassCode == ticket.SeatClassCode {
				continue
			}
			if seat.AvailableCount > 0 {
				seatOptions = append(seatOptions, seat)
			}
		}
		if len(seatOptions) == 0 {
			continue
		}
		option.SeatOptions = seatOptions
		filtered = append(filtered, option)
	}

	return dto.ChangeOptionsResponse{
		OriginalTicket: ticketResponse(ticket),
		Options:        filtered,
	}, nil
}

func (s *TicketService) Refund(userID, ticketID uint64, req dto.RefundTicketRequest) (dto.RefundTicketResponse, error) {
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		reason = "用户申请退票"
	}
	key := strings.TrimSpace(req.IdempotencyKey)
	if key == "" {
		key = fmt.Sprintf("refund-auto-%d-%d", userID, time.Now().UnixNano())
	}

	var refund model.Refund
	var refunded model.Ticket
	err := s.tickets.Transaction(func(tx *gorm.DB) error {
		var existing model.Refund
		err := tx.Where("user_id = ? AND idempotency_key = ?", userID, key).First(&existing).Error
		if err == nil {
			refund = existing
			return tx.Where("id = ? AND user_id = ?", existing.TicketID, userID).First(&refunded).Error
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		var ticket model.Ticket
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", ticketID, userID).
			First(&ticket).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.New(http.StatusNotFound, response.CodeNotFound, "车票不存在")
			}
			return err
		}
		if ticket.Status == model.TicketStatusRefunded {
			refunded = ticket
			_ = tx.Where("ticket_id = ? AND status = ?", ticket.ID, model.RefundStatusSuccess).Order("id DESC").First(&refund).Error
			return nil
		}
		if ticket.Status != model.TicketStatusIssued {
			return apperrors.New(http.StatusConflict, response.CodeTicketNotRefundable, "当前车票不可退")
		}

		departTime, err := departTimeForTicket(tx, ticket)
		if err != nil {
			return err
		}
		cutoffMinutes := systemSettingInt(tx, "refund_cutoff_minutes", 0)
		if !departTime.After(time.Now().Add(time.Duration(cutoffMinutes) * time.Minute)) {
			return apperrors.New(http.StatusConflict, response.CodeTicketNotRefundable, "已发车车票不可退")
		}

		row, err := repository.FindInventoryForOrder(tx.Clauses(clause.Locking{Strength: "UPDATE"}), ticket.TravelDate, ticket.TrainID, ticket.FromStationID, ticket.ToStationID, ticket.SeatClassCode)
		if err != nil {
			return err
		}
		if row.SoldCount <= 0 {
			return apperrors.New(http.StatusConflict, response.CodeTicketNotRefundable, "车票库存状态异常")
		}
		if err := tx.Model(&model.Inventory{}).
			Where("id = ? AND sold_count > 0", row.InventoryID).
			Updates(map[string]any{
				"available_count": gorm.Expr("available_count + ?", 1),
				"sold_count":      gorm.Expr("sold_count - ?", 1),
			}).Error; err != nil {
			return err
		}

		now := time.Now()
		if err := tx.Model(&ticket).Updates(map[string]any{
			"status":      model.TicketStatusRefunded,
			"refunded_at": now,
		}).Error; err != nil {
			return err
		}

		var payment model.Payment
		if err := tx.Where("order_id = ? AND status = ?", ticket.OrderID, model.PaymentStatusSuccess).Order("id DESC").First(&payment).Error; err != nil {
			return err
		}
		refund = model.Refund{
			RefundNo:       makeBizNo("R"),
			TicketID:       ticket.ID,
			PaymentID:      payment.ID,
			UserID:         userID,
			AmountCents:    payment.AmountCents,
			Status:         model.RefundStatusSuccess,
			Reason:         reason,
			IdempotencyKey: key,
			RefundedAt:     now,
		}
		if err := tx.Create(&refund).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Order{}).
			Where("id = ? AND user_id = ?", ticket.OrderID, userID).
			Update("status", model.OrderStatusClosed).Error; err != nil {
			return err
		}

		ticket.Status = model.TicketStatusRefunded
		ticket.RefundedAt = &now
		refunded = ticket
		return nil
	})
	if err != nil {
		return dto.RefundTicketResponse{}, err
	}

	return dto.RefundTicketResponse{
		RefundNo: refund.RefundNo,
		Ticket:   ticketResponse(refunded),
	}, nil
}

func (s *TicketService) Change(userID, ticketID uint64, req dto.ChangeTicketRequest) (dto.ChangeTicketResponse, error) {
	newDate, err := time.ParseInLocation("2006-01-02", req.NewTravelDate, time.Local)
	if err != nil {
		return dto.ChangeTicketResponse{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "改签日期格式不正确")
	}
	key := strings.TrimSpace(req.IdempotencyKey)
	if key == "" {
		key = fmt.Sprintf("change-auto-%d-%d", userID, time.Now().UnixNano())
	}

	var record model.ChangeRecord
	var oldTicket model.Ticket
	var newTicket model.Ticket
	err = s.tickets.Transaction(func(tx *gorm.DB) error {
		var existing model.ChangeRecord
		err := tx.Where("user_id = ? AND idempotency_key = ?", userID, key).First(&existing).Error
		if err == nil {
			record = existing
			if err := tx.Where("id = ? AND user_id = ?", existing.OldTicketID, userID).First(&oldTicket).Error; err != nil {
				return err
			}
			return tx.Where("id = ? AND user_id = ?", existing.NewTicketID, userID).First(&newTicket).Error
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		var ticket model.Ticket
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ?", ticketID, userID).
			First(&ticket).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.New(http.StatusNotFound, response.CodeNotFound, "车票不存在")
			}
			return err
		}
		if ticket.Status != model.TicketStatusIssued {
			return apperrors.New(http.StatusConflict, response.CodeTicketNotChangeable, "当前车票不可改签")
		}
		oldDepartTime, err := departTimeForTicket(tx, ticket)
		if err != nil {
			return err
		}
		cutoffMinutes := systemSettingInt(tx, "change_cutoff_minutes", 0)
		if !oldDepartTime.After(time.Now().Add(time.Duration(cutoffMinutes) * time.Minute)) {
			return apperrors.New(http.StatusConflict, response.CodeTicketNotChangeable, "已发车车票不可改签")
		}
		if ticket.TrainID == req.NewTrainID && ticket.TravelDate.Equal(newDate) && ticket.SeatClassCode == req.NewSeatClassCode {
			return apperrors.New(http.StatusBadRequest, response.CodeValidationError, "新车次和席别不能与原票相同")
		}

		oldRow, err := repository.FindInventoryForOrder(tx.Clauses(clause.Locking{Strength: "UPDATE"}), ticket.TravelDate, ticket.TrainID, ticket.FromStationID, ticket.ToStationID, ticket.SeatClassCode)
		if err != nil {
			return err
		}
		newRow, err := repository.FindInventoryForOrder(tx.Clauses(clause.Locking{Strength: "UPDATE"}), newDate, req.NewTrainID, ticket.FromStationID, ticket.ToStationID, req.NewSeatClassCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.New(http.StatusNotFound, response.CodeNotFound, "未找到可改签的车次席别")
			}
			return err
		}
		if newRow.AvailableCount <= 0 {
			return apperrors.New(http.StatusConflict, response.CodeInsufficientInventory, "改签车次余票不足")
		}

		if err := tx.Model(&model.Inventory{}).
			Where("id = ? AND sold_count > 0", oldRow.InventoryID).
			Updates(map[string]any{
				"available_count": gorm.Expr("available_count + ?", 1),
				"sold_count":      gorm.Expr("sold_count - ?", 1),
			}).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Inventory{}).
			Where("id = ? AND available_count > 0", newRow.InventoryID).
			Updates(map[string]any{
				"available_count": gorm.Expr("available_count - ?", 1),
				"sold_count":      gorm.Expr("sold_count + ?", 1),
			}).Error; err != nil {
			return err
		}

		now := time.Now()
		if err := tx.Model(&ticket).Update("status", model.TicketStatusChangedOut).Error; err != nil {
			return err
		}

		newTicket = model.Ticket{
			TicketNo:        makeBizNo("T"),
			OrderID:         ticket.OrderID,
			UserID:          userID,
			TrainID:         newRow.TrainID,
			TrainNo:         newRow.TrainNo,
			TravelDate:      newDate,
			FromStationID:   newRow.FromStationID,
			FromStationName: newRow.FromStationName,
			ToStationID:     newRow.ToStationID,
			ToStationName:   newRow.ToStationName,
			SeatClassCode:   newRow.SeatClassCode,
			SeatClassName:   seatClassName(newRow.SeatClassCode),
			CoachNo:         makeCoachNo(newRow.SoldCount + 1),
			SeatNo:          makeSeatNo(newRow.SoldCount + 1),
			PassengerName:   ticket.PassengerName,
			IDCardNo:        ticket.IDCardNo,
			Status:          model.TicketStatusIssued,
			IssuedAt:        now,
		}
		if err := tx.Create(&newTicket).Error; err != nil {
			return err
		}

		record = model.ChangeRecord{
			ChangeNo:       makeBizNo("C"),
			OldTicketID:    ticket.ID,
			NewTicketID:    newTicket.ID,
			UserID:         userID,
			PriceDiffCents: newRow.PriceCents - oldRow.PriceCents,
			Status:         model.ChangeStatusSuccess,
			IdempotencyKey: key,
			ChangedAt:      now,
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Order{}).
			Where("id = ? AND user_id = ?", ticket.OrderID, userID).
			Updates(map[string]any{
				"train_id":          newTicket.TrainID,
				"train_no":          newTicket.TrainNo,
				"travel_date":       newTicket.TravelDate,
				"from_station_id":   newTicket.FromStationID,
				"from_station_name": newTicket.FromStationName,
				"to_station_id":     newTicket.ToStationID,
				"to_station_name":   newTicket.ToStationName,
				"seat_class_code":   newTicket.SeatClassCode,
				"seat_class_name":   newTicket.SeatClassName,
				"amount_cents":      newRow.PriceCents,
				"status":            model.OrderStatusPaid,
			}).Error; err != nil {
			return err
		}

		ticket.Status = model.TicketStatusChangedOut
		oldTicket = ticket
		return nil
	})
	if err != nil {
		return dto.ChangeTicketResponse{}, err
	}

	return dto.ChangeTicketResponse{
		ChangeNo:       record.ChangeNo,
		PriceDiffCents: record.PriceDiffCents,
		OldTicket:      ticketResponse(oldTicket),
		NewTicket:      ticketResponse(newTicket),
	}, nil
}

func trainSearchRowsToResponses(date time.Time, rows []repository.TrainSearchRow) []dto.TrainSearchItemResponse {
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
	return result
}

func ticketResponse(ticket model.Ticket) dto.TicketResponse {
	refundedAt := (*string)(nil)
	if ticket.RefundedAt != nil {
		formatted := ticket.RefundedAt.Format(time.RFC3339)
		refundedAt = &formatted
	}

	result := dto.TicketResponse{
		ID:             ticket.ID,
		TicketNo:       ticket.TicketNo,
		OrderID:        ticket.OrderID,
		TrainID:        ticket.TrainID,
		TrainNo:        ticket.TrainNo,
		TravelDate:     ticket.TravelDate.Format("2006-01-02"),
		FromStation:    dto.StationResponse{ID: ticket.FromStationID, Name: ticket.FromStationName},
		ToStation:      dto.StationResponse{ID: ticket.ToStationID, Name: ticket.ToStationName},
		SeatClassCode:  ticket.SeatClassCode,
		SeatClassName:  ticket.SeatClassName,
		CoachNo:        ticket.CoachNo,
		SeatNo:         ticket.SeatNo,
		PassengerName:  ticket.PassengerName,
		IDCardNoMasked: maskIDCardNo(ticket.IDCardNo),
		Status:         string(ticket.Status),
		IssuedAt:       ticket.IssuedAt.Format(time.RFC3339),
		RefundedAt:     refundedAt,
	}
	if ticket.DepartTime != nil {
		result.DepartTime = ticket.DepartTime.Format(time.RFC3339)
	}
	if ticket.ArriveTime != nil {
		result.ArriveTime = ticket.ArriveTime.Format(time.RFC3339)
	}
	return result
}

func departTimeForTicket(tx *gorm.DB, ticket model.Ticket) (time.Time, error) {
	var stop model.TrainStop
	err := tx.Where("train_id = ? AND station_id = ?", ticket.TrainID, ticket.FromStationID).First(&stop).Error
	if err != nil {
		return time.Time{}, err
	}
	return combineDateClock(ticket.TravelDate, stop.DepartClock, stop.DayOffset), nil
}

func makeCoachNo(sequence int) string {
	if sequence <= 0 {
		sequence = 1
	}
	return fmt.Sprintf("%02d", sequence%8+1)
}

func makeSeatNo(sequence int) string {
	columns := []string{"A", "B", "C", "D", "F"}
	if sequence <= 0 {
		sequence = 1
	}
	return fmt.Sprintf("%02d%s", sequence%18+1, columns[sequence%len(columns)])
}

func maskIDCardNo(value string) string {
	text := strings.TrimSpace(value)
	if len(text) <= 8 {
		return text
	}
	return text[:4] + strings.Repeat("*", len(text)-8) + text[len(text)-4:]
}
