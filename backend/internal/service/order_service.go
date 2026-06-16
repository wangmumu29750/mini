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
	"mini-12306/backend/pkg/mock"
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
	return s.create(userID, req, nil)
}

func (s *OrderService) ClerkCreate(clerkUserID uint64, req dto.ClerkCreateOrderRequest) (dto.OrderResponse, error) {
	req.PassengerName = strings.TrimSpace(req.PassengerName)
	req.IDCardNo = strings.TrimSpace(req.IDCardNo)
	req.Phone = strings.TrimSpace(req.Phone)
	req.BankCardNo = strings.TrimSpace(req.BankCardNo)

	if err := mock.VerifyIdentity(req.PassengerName, req.IDCardNo, req.Phone, req.BankCardNo); err != nil {
		return dto.OrderResponse{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, err.Error())
	}

	passenger := &model.PassengerProfile{
		RealName:       req.PassengerName,
		IDCardNo:       req.IDCardNo,
		Phone:          req.Phone,
		BankCardNo:     req.BankCardNo,
		PassengerType:  model.PassengerTypeAdult,
		VerifiedStatus: model.VerificationStatusVerified,
	}
	if len(req.Passengers) == 0 {
		req.Passengers = []dto.OrderPassengerDTO{{
			PassengerID: 1,
			SeatType:    "SECOND",
			TicketType:  string(model.TicketTypeAdult),
		}}
	}
	return s.create(clerkUserID, req.CreateOrderRequest, passenger)
}

func (s *OrderService) create(userID uint64, req dto.CreateOrderRequest, walkUpPassenger *model.PassengerProfile) (dto.OrderResponse, error) {
	if req.FromStationID == req.ToStationID {
		return dto.OrderResponse{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "出发站和到达站不能相同")
	}
	if len(req.Passengers) == 0 {
		return dto.OrderResponse{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "至少选择一位乘车人")
	}

	travelDate, err := time.ParseInLocation("2006-01-02", req.TravelDate, time.Local)
	if err != nil {
		return dto.OrderResponse{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "乘车日期格式不正确")
	}
	if err := validateTicketTravelDate(travelDate); err != nil {
		return dto.OrderResponse{}, err
	}

	key := strings.TrimSpace(req.IdempotencyKey)
	if key == "" {
		key = fmt.Sprintf("auto-%d-%d", userID, time.Now().UnixNano())
	}

	var created model.Order
	err = s.orders.Transaction(func(tx *gorm.DB) error {
		var existing model.Order
		err := tx.Preload("Items").Preload("Tickets").
			Where("user_id = ? AND idempotency_key = ?", userID, key).
			First(&existing).Error
		if err == nil {
			created = existing
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		calculator := newTicketPriceCalculator()
		orderItems := make([]model.OrderItem, 0, len(req.Passengers))
		totalAmount := int64(0)
		var snapshot repository.OrderInventoryRow

		for i, passengerReq := range req.Passengers {
			passengerReq.SeatType = strings.ToUpper(strings.TrimSpace(passengerReq.SeatType))
			passengerReq.TicketType = strings.ToUpper(strings.TrimSpace(passengerReq.TicketType))
			if passengerReq.SeatType == "" || passengerReq.TicketType == "" {
				return apperrors.New(http.StatusBadRequest, response.CodeValidationError, "席别和票种不能为空")
			}

			profile, err := passengerProfileForOrder(tx, userID, passengerReq.PassengerID, walkUpPassenger)
			if err != nil {
				return err
			}

			row, err := repository.FindInventoryForOrder(tx.Clauses(clause.Locking{Strength: "UPDATE"}), travelDate, req.TrainID, req.FromStationID, req.ToStationID, passengerReq.SeatType)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return apperrors.New(http.StatusNotFound, response.CodeNotFound, "未找到可预订的车次席别")
				}
				return err
			}
			if row.AvailableCount <= 0 {
				return apperrors.New(http.StatusConflict, response.CodeInsufficientInventory, "余票不足")
			}
			if err := validateFutureDeparture(row.TravelDate, row.DepartClock, row.DepartDayOffset); err != nil {
				return err
			}

			result := tx.Model(&model.Inventory{}).
				Where("id = ? AND available_count > 0", row.InventoryID).
				Updates(map[string]any{
					"available_count": gorm.Expr("available_count - ?", 1),
					"locked_count":    gorm.Expr("locked_count + ?", 1),
				})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return apperrors.New(http.StatusConflict, response.CodeInsufficientInventory, "余票不足")
			}

			realPrice, err := calculator.Calculate(row.PriceCents, row.TrainType, passengerReq.SeatType, passengerReq.TicketType)
			if err != nil {
				return err
			}

			if i == 0 {
				snapshot = row
			}
			totalAmount += realPrice
			orderItems = append(orderItems, model.OrderItem{
				PassengerID:    profile.ID,
				PassengerName:  profile.RealName,
				IDCardNo:       profile.IDCardNo,
				PassengerType:  profile.PassengerType,
				SeatType:       model.SeatType(passengerReq.SeatType),
				TicketType:     model.TicketType(passengerReq.TicketType),
				BasePriceCents: row.PriceCents,
				RealPriceCents: realPrice,
			})
		}

		created = model.Order{
			OrderNo:         makeBizNo("O"),
			UserID:          userID,
			TrainID:         snapshot.TrainID,
			TrainNo:         snapshot.TrainNo,
			TravelDate:      travelDate,
			FromStationID:   snapshot.FromStationID,
			FromStationName: snapshot.FromStationName,
			ToStationID:     snapshot.ToStationID,
			ToStationName:   snapshot.ToStationName,
			SeatClassCode:   string(orderItems[0].SeatType),
			SeatClassName:   seatClassName(string(orderItems[0].SeatType)),
			PassengerName:   orderItems[0].PassengerName,
			IDCardNo:        orderItems[0].IDCardNo,
			AmountCents:     totalAmount,
			Status:          model.OrderStatusPendingPayment,
			PayExpiresAt:    time.Now().Add(time.Duration(systemSettingInt(tx, "order_pay_expire_minutes", int(s.cfg.OrderPayExpireDuration().Minutes()))) * time.Minute),
			IdempotencyKey:  key,
		}
		if err := tx.Create(&created).Error; err != nil {
			return err
		}

		for i := range orderItems {
			orderItems[i].OrderID = created.ID
		}
		if err := tx.Create(&orderItems).Error; err != nil {
			return err
		}
		created.Items = orderItems
		return nil
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
			Preload("Items", func(db *gorm.DB) *gorm.DB { return db.Order("id ASC") }).
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
			if err := releaseOrderLocks(tx, order); err != nil {
				return err
			}
			if err := tx.Model(&order).Update("status", model.OrderStatusClosed).Error; err != nil {
				return err
			}
			return apperrors.New(http.StatusConflict, response.CodeInvalidOrderState, "订单已超过支付时间")
		}
		if len(order.Items) == 0 {
			return apperrors.New(http.StatusConflict, response.CodeInvalidOrderState, "订单明细为空")
		}

		now := time.Now()
		createdTickets := make([]model.Ticket, 0, len(order.Items))
		for i, item := range order.Items {
			row, err := repository.FindInventoryForOrder(tx.Clauses(clause.Locking{Strength: "UPDATE"}), order.TravelDate, order.TrainID, order.FromStationID, order.ToStationID, string(item.SeatType))
			if err != nil {
				return err
			}
			if row.LockedCount <= 0 {
				return apperrors.New(http.StatusConflict, response.CodeInvalidOrderState, "订单锁票状态异常")
			}
			result := tx.Model(&model.Inventory{}).
				Where("id = ? AND locked_count > 0", row.InventoryID).
				Updates(map[string]any{
					"locked_count": gorm.Expr("locked_count - ?", 1),
					"sold_count":   gorm.Expr("sold_count + ?", 1),
				})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return apperrors.New(http.StatusConflict, response.CodeInvalidOrderState, "订单锁票状态异常")
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
				SeatClassCode:   string(item.SeatType),
				SeatClassName:   seatClassName(string(item.SeatType)),
				TicketType:      item.TicketType,
				RealPriceCents:  item.RealPriceCents,
				CoachNo:         makeCoachNo(row.SoldCount + i + 1),
				SeatNo:          makeSeatNo(row.SoldCount + i + 1),
				PassengerName:   item.PassengerName,
				IDCardNo:        item.IDCardNo,
				Status:          model.TicketStatusIssued,
				IssuedAt:        now,
			}
			if err := tx.Create(&ticket).Error; err != nil {
				return err
			}
			item.TicketID = &ticket.ID
			item.TicketNo = ticket.TicketNo
			if err := tx.Save(&item).Error; err != nil {
				return err
			}
			createdTickets = append(createdTickets, ticket)
		}

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

		paidAt := now
		if err := tx.Model(&order).Updates(map[string]any{
			"status":  model.OrderStatusPaid,
			"paid_at": paidAt,
		}).Error; err != nil {
			return err
		}
		order.Status = model.OrderStatusPaid
		order.PaidAt = &paidAt
		order.Tickets = createdTickets
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
			Preload("Items").
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

		if err := releaseOrderLocks(tx, order); err != nil {
			return err
		}
		if err := tx.Model(&order).Update("status", model.OrderStatusCancelled).Error; err != nil {
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

	ticketItems := make([]dto.OrderTicketItem, 0, len(order.Items))
	for _, item := range order.Items {
		ticketItems = append(ticketItems, dto.OrderTicketItem{
			PassengerName:  item.PassengerName,
			SeatType:       string(item.SeatType),
			TicketType:     string(item.TicketType),
			RealPriceCents: item.RealPriceCents,
			TicketNo:       item.TicketNo,
			Status:         "",
		})
	}

	itemCount := len(order.Items)
	if itemCount == 0 && len(order.Tickets) > 0 {
		itemCount = len(order.Tickets)
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
		ItemCount:     itemCount,
		Status:        string(order.Status),
		PayExpiresAt:  order.PayExpiresAt.Format(time.RFC3339),
		PaidAt:        paidAt,
		Tickets:       ticketItems,
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

func passengerProfileForOrder(tx *gorm.DB, userID, passengerID uint64, walkUpPassenger *model.PassengerProfile) (model.PassengerProfile, error) {
	if walkUpPassenger != nil {
		copy := *walkUpPassenger
		copy.ID = passengerID
		if copy.ID == 0 {
			copy.ID = 1
		}
		return copy, nil
	}

	var profile model.PassengerProfile
	err := tx.Where("id = ? AND user_id = ? AND verified_status = ?", passengerID, userID, model.VerificationStatusVerified).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.PassengerProfile{}, apperrors.New(http.StatusNotFound, response.CodeNotFound, "乘车人不存在或未实名")
		}
		return model.PassengerProfile{}, err
	}
	return profile, nil
}

func releaseOrderLocks(tx *gorm.DB, order model.Order) error {
	if len(order.Items) == 0 {
		return releaseSingleOrderLock(tx, order)
	}
	for _, item := range order.Items {
		row, err := repository.FindInventoryForOrder(tx.Clauses(clause.Locking{Strength: "UPDATE"}), order.TravelDate, order.TrainID, order.FromStationID, order.ToStationID, string(item.SeatType))
		if err != nil {
			return err
		}
		if row.LockedCount <= 0 {
			continue
		}
		result := tx.Model(&model.Inventory{}).
			Where("id = ? AND locked_count > 0", row.InventoryID).
			Updates(map[string]any{
				"available_count": gorm.Expr("available_count + ?", 1),
				"locked_count":    gorm.Expr("locked_count - ?", 1),
			})
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func releaseSingleOrderLock(tx *gorm.DB, order model.Order) error {
	row, err := repository.FindInventoryForOrder(tx.Clauses(clause.Locking{Strength: "UPDATE"}), order.TravelDate, order.TrainID, order.FromStationID, order.ToStationID, order.SeatClassCode)
	if err != nil {
		return err
	}
	if row.LockedCount <= 0 {
		return nil
	}
	result := tx.Model(&model.Inventory{}).
		Where("id = ? AND locked_count > 0", row.InventoryID).
		Updates(map[string]any{
			"available_count": gorm.Expr("available_count + ?", 1),
			"locked_count":    gorm.Expr("locked_count - ?", 1),
		})
	return result.Error
}
