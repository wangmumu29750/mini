package service

import (
	"mini-12306/backend/internal/config"
	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/model"
	"mini-12306/backend/internal/repository"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestNewOrderService(t *testing.T) {
	type args struct {
		cfg    config.Config
		orders *repository.OrderRepository
	}
	tests := []struct {
		name string
		args args
		want *OrderService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrderService(tt.args.cfg, tt.args.orders); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrderService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_Create(t *testing.T) {
	type args struct {
		userID uint64
		req    dto.CreateOrderRequest
	}
	tests := []struct {
		name    string
		s       *OrderService
		args    args
		want    dto.OrderResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Create(tt.args.userID, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderService.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_ClerkCreate(t *testing.T) {
	type args struct {
		clerkUserID uint64
		req         dto.ClerkCreateOrderRequest
	}
	tests := []struct {
		name    string
		s       *OrderService
		args    args
		want    dto.OrderResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ClerkCreate(tt.args.clerkUserID, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.ClerkCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderService.ClerkCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_create(t *testing.T) {
	type args struct {
		userID          uint64
		req             dto.CreateOrderRequest
		walkUpPassenger *model.PassengerProfile
	}
	tests := []struct {
		name    string
		s       *OrderService
		args    args
		want    dto.OrderResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.create(tt.args.userID, tt.args.req, tt.args.walkUpPassenger)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderService.create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_List(t *testing.T) {
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		s       *OrderService
		args    args
		want    []dto.OrderResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.List(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_Detail(t *testing.T) {
	type args struct {
		userID  uint64
		orderID uint64
	}
	tests := []struct {
		name    string
		s       *OrderService
		args    args
		want    dto.OrderResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Detail(tt.args.userID, tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.Detail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderService.Detail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_Pay(t *testing.T) {
	type args struct {
		userID  uint64
		orderID uint64
		req     dto.PayOrderRequest
	}
	tests := []struct {
		name    string
		s       *OrderService
		args    args
		want    dto.PaymentResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Pay(tt.args.userID, tt.args.orderID, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.Pay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderService.Pay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_Cancel(t *testing.T) {
	type args struct {
		userID  uint64
		orderID uint64
	}
	tests := []struct {
		name    string
		s       *OrderService
		args    args
		want    dto.OrderResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Cancel(tt.args.userID, tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.Cancel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderService.Cancel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_orderResponse(t *testing.T) {
	type args struct {
		order model.Order
	}
	tests := []struct {
		name string
		args args
		want dto.OrderResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := orderResponse(tt.args.order); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("orderResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeBizNo(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeBizNo(tt.args.prefix); got != tt.want {
				t.Errorf("makeBizNo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_passengerProfileForOrder(t *testing.T) {
	type args struct {
		tx              *gorm.DB
		userID          uint64
		passengerID     uint64
		walkUpPassenger *model.PassengerProfile
	}
	tests := []struct {
		name    string
		args    args
		want    model.PassengerProfile
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := passengerProfileForOrder(tt.args.tx, tt.args.userID, tt.args.passengerID, tt.args.walkUpPassenger)
			if (err != nil) != tt.wantErr {
				t.Errorf("passengerProfileForOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("passengerProfileForOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_releaseOrderLocks(t *testing.T) {
	type args struct {
		tx    *gorm.DB
		order model.Order
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := releaseOrderLocks(tt.args.tx, tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("releaseOrderLocks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_releaseSingleOrderLock(t *testing.T) {
	type args struct {
		tx    *gorm.DB
		order model.Order
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := releaseSingleOrderLock(tt.args.tx, tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("releaseSingleOrderLock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
