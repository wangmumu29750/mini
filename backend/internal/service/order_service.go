package service

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"mini-12306/backend/internal/config"
	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/model"
	"mini-12306/backend/internal/repository"
	apperrors "mini-12306/backend/pkg/errors"
	"mini-12306/backend/pkg/response"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderService struct {
	cfg    config.Config
	orders *repository.OrderRepository
}

func NewOrderService(cfg config.Config, orders *repository.OrderRepository) *OrderService {
	return &OrderService{cfg: cfg, orders: orders}
}

func (s *OrderService) Create(userID uint64, req dto.CreateOrderRequest) (dto.OrderResponse, error) {
	if req.FromStationID == req.ToStationID {
		return dto.OrderResponse{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "出发站和到达站不能相同")
	}

	travelDate, err := time.ParseInLocation("2006-01-02", req.TravelDate, time.Local)
	if err != nil {
		return dto.OrderResponse{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "乘车日期格式不正确")
	}

	key := strings.TrimSpace(req.IdempotencyKey)
	if key == "" {
		key = fmt.Sprintf("auto-%d-%d", userID, time.Now().UnixNano())
	}

	var created model.Order
	err = s.orders.Transaction(func(tx *gorm.DB) error {
		var existing model.Order
		err := tx.Preload("Tickets").
			Where("user_id = ? AND idempotency_key = ?", userID, key).
			First(&existing).Error
		if err == nil {
			created = existing
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		var profile model.PassengerProfile
		if err := tx.Where("user_id = ? AND verified_status = ?", userID, model.VerificationStatusVerified).First(&profile).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.New(http.StatusForbidden, response.CodeForbidden, "请先完成实名乘车人信息")
			}
			return err
		}

		row, err := repository.FindInventoryForOrder(tx.Clauses(clause.Locking{Strength: "UPDATE"}), travelDate, req.TrainID, req.FromStationID, req.ToStationID, req.SeatClassCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.New(http.StatusNotFound, response.CodeNotFound, "未找到可预订的车次席别")
			}
			return err
		}
		if row.AvailableCount <= 0 {
			return apperrors.New(http.StatusConflict, response.CodeInsufficientInventory, "余票不足")
		}

		if err := tx.Model(&model.Inventory{}).
			Where("id = ? AND available_count > 0", row.InventoryID).
			Updates(map[string]any{
				"available_count": gorm.Expr("available_count - ?", 1),
				"locked_count":    gorm.Expr("locked_count + ?", 1),
			}).Error; err != nil {
			return err
		}

		created = model.Order{
			OrderNo:         makeBizNo("O"),
			UserID:          userID,
			TrainID:         row.TrainID,
			TrainNo:         row.TrainNo,
			TravelDate:      travelDate,
			FromStationID:   row.FromStationID,
			FromStationName: row.FromStationName,
			ToStationID:     row.ToStationID,
			ToStationName:   row.ToStationName,
			SeatClassCode:   row.SeatClassCode,
			SeatClassName:   seatClassName(row.SeatClassCode),
			PassengerName:   profile.RealName,
			IDCardNo:        profile.IDCardNo,
			AmountCents:     row.PriceCents,
			Status:          model.OrderStatusPendingPayment,
			PayExpiresAt:    time.Now().Add(s.cfg.OrderPayExpireDuration()),
			IdempotencyKey:  key,
		}
		return tx.Create(&created).Error
	})
	if err != nil {
		return dto.OrderResponse{}, err
	}
	return orderResponse(created), nil
}

func (s *OrderService) List(userID uint64) ([]dto.OrderResponse, error) {
	orders, err := s.orders.ListByUser(userID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		result = append(result, orderResponse(order))
	}
	return result, nil
}

func (s *OrderService) Detail(userID, orderID uint64) (dto.OrderResponse, error) {
	order, err := s.orders.FindByUserAndID(userID, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.OrderResponse{}, apperrors.New(http.StatusNotFound, response.CodeNotFound, "订单不存在")
		}
		return dto.OrderResponse{}, err
	}
	return orderResponse(order), nil
}

func (s *OrderService) Pay(userID, orderID uint64, req dto.PayOrderRequest) (dto.PaymentResponse, error) {
	channel := strings.TrimSpace(req.Channel)
	if channel == "" {
		channel = "MOCK_BANK"
	}

	var payment model.Payment
	var paidOrder model.Order
	err := s.orders.Transaction(func(tx *gorm.DB) error {
		var order model.Order
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Tickets").
			Where("id = ? AND user_id = ?", orderID, userID).
			First(&order).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.New(http.StatusNotFound, response.CodeNotFound, "订单不存在")
			}
			return err
		}

		if order.Status == model.OrderStatusPaid {
			paidOrder = order
			var existing model.Payment
			_ = tx.Where("order_id = ? AND status = ?", order.ID, model.PaymentStatusSuccess).Order("id DESC").First(&existing).Error
			payment = existing
			return nil
		}
		if order.Status != model.OrderStatusPendingPayment {
			return apperrors.New(http.StatusConflict, response.CodeInvalidOrderState, "当前订单状态不可支付")
		}
		if time.Now().After(order.PayExpiresAt) {
			if err := releaseOrderLock(tx, order, model.OrderStatusClosed); err != nil {
				return err
			}
			return apperrors.New(http.StatusConflict, response.CodeInvalidOrderState, "订单已超过支付时间")
		}

		row, err := repository.FindInventoryForOrder(tx.Clauses(clause.Locking{Strength: "UPDATE"}), order.TravelDate, order.TrainID, order.FromStationID, order.ToStationID, order.SeatClassCode)
		if err != nil {
			return err
		}
		if row.LockedCount <= 0 {
			return apperrors.New(http.StatusConflict, response.CodeInvalidOrderState, "订单锁票状态异常")
		}
		if err := tx.Model(&model.Inventory{}).
			Where("id = ? AND locked_count > 0", row.InventoryID).
			Updates(map[string]any{
				"locked_count": gorm.Expr("locked_count - ?", 1),
				"sold_count":   gorm.Expr("sold_count + ?", 1),
			}).Error; err != nil {
			return err
		}

		now := time.Now()
		payment = model.Payment{
			PaymentNo:   makeBizNo("P"),
			OrderID:     order.ID,
			UserID:      userID,
			AmountCents: order.AmountCents,
			Channel:     channel,
			Status:      model.PaymentStatusSuccess,
			PaidAt:      now,
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		ticket := model.Ticket{
			TicketNo:        makeBizNo("T"),
			OrderID:         order.ID,
			UserID:          userID,
			TrainID:         order.TrainID,
			TrainNo:         order.TrainNo,
			TravelDate:      order.TravelDate,
			FromStationID:   order.FromStationID,
			FromStationName: order.FromStationName,
			ToStationID:     order.ToStationID,
			ToStationName:   order.ToStationName,
			SeatClassCode:   order.SeatClassCode,
			SeatClassName:   order.SeatClassName,
			CoachNo:         makeCoachNo(row.SoldCount + 1),
			SeatNo:          makeSeatNo(row.SoldCount + 1),
			PassengerName:   order.PassengerName,
			IDCardNo:        order.IDCardNo,
			Status:          model.TicketStatusIssued,
			IssuedAt:        now,
		}
		if err := tx.Create(&ticket).Error; err != nil {
			return err
		}

		paidAt := now
		if err := tx.Model(&order).Updates(map[string]any{
			"status":  model.OrderStatusPaid,
			"paid_at": paidAt,
		}).Error; err != nil {
			return err
		}
		order.Status = model.OrderStatusPaid
		order.PaidAt = &paidAt
		order.Tickets = []model.Ticket{ticket}
		paidOrder = order
		return nil
	})
	if err != nil {
		return dto.PaymentResponse{}, err
	}

	return dto.PaymentResponse{
		PaymentNo: payment.PaymentNo,
		Order:     orderResponse(paidOrder),
	}, nil
}

func (s *OrderService) Cancel(userID, orderID uint64) (dto.OrderResponse, error) {
	var cancelled model.Order
	err := s.orders.Transaction(func(tx *gorm.DB) error {
		var order model.Order
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Tickets").
			Where("id = ? AND user_id = ?", orderID, userID).
			First(&order).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperrors.New(http.StatusNotFound, response.CodeNotFound, "订单不存在")
			}
			return err
		}

		if order.Status == model.OrderStatusCancelled || order.Status == model.OrderStatusClosed {
			cancelled = order
			return nil
		}
		if order.Status != model.OrderStatusPendingPayment {
			return apperrors.New(http.StatusConflict, response.CodeInvalidOrderState, "只有待支付订单可以取消")
		}

		if err := releaseOrderLock(tx, order, model.OrderStatusCancelled); err != nil {
			return err
		}
		order.Status = model.OrderStatusCancelled
		cancelled = order
		return nil
	})
	if err != nil {
		return dto.OrderResponse{}, err
	}
	return orderResponse(cancelled), nil
}

func orderResponse(order model.Order) dto.OrderResponse {
	var paidAt *string
	if order.PaidAt != nil {
		formatted := order.PaidAt.Format(time.RFC3339)
		paidAt = &formatted
	}

	result := dto.OrderResponse{
		ID:            order.ID,
		OrderNo:       order.OrderNo,
		TrainID:       order.TrainID,
		TrainNo:       order.TrainNo,
		TravelDate:    order.TravelDate.Format("2006-01-02"),
		FromStation:   dto.StationResponse{ID: order.FromStationID, Name: order.FromStationName},
		ToStation:     dto.StationResponse{ID: order.ToStationID, Name: order.ToStationName},
		SeatClassCode: order.SeatClassCode,
		SeatClassName: order.SeatClassName,
		PassengerName: order.PassengerName,
		AmountCents:   order.AmountCents,
		Status:        string(order.Status),
		PayExpiresAt:  order.PayExpiresAt.Format(time.RFC3339),
		PaidAt:        paidAt,
	}
	if order.DepartTime != nil {
		result.DepartTime = order.DepartTime.Format(time.RFC3339)
	}
	if order.ArriveTime != nil {
		result.ArriveTime = order.ArriveTime.Format(time.RFC3339)
	}
	if len(order.Tickets) > 0 {
		result.TicketNo = order.Tickets[0].TicketNo
		result.TicketStatus = string(order.Tickets[0].Status)
	}
	return result
}

func makeBizNo(prefix string) string {
	return fmt.Sprintf("%s%s%06d", prefix, time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)
}

func releaseOrderLock(tx *gorm.DB, order model.Order, nextStatus model.OrderStatus) error {
	row, err := repository.FindInventoryForOrder(tx.Clauses(clause.Locking{Strength: "UPDATE"}), order.TravelDate, order.TrainID, order.FromStationID, order.ToStationID, order.SeatClassCode)
	if err != nil {
		return err
	}
	if row.LockedCount > 0 {
		if err := tx.Model(&model.Inventory{}).
			Where("id = ? AND locked_count > 0", row.InventoryID).
			Updates(map[string]any{
				"available_count": gorm.Expr("available_count + ?", 1),
				"locked_count":    gorm.Expr("locked_count - ?", 1),
			}).Error; err != nil {
			return err
		}
	}
	return tx.Model(&order).Update("status", nextStatus).Error
}
