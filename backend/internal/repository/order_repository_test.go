package repository

import (
	"mini-12306/backend/internal/model"
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestNewOrderRepository(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *OrderRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrderRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrderRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderRepository_Transaction(t *testing.T) {
	type args struct {
		fn func(tx *gorm.DB) error
	}
	tests := []struct {
		name    string
		r       *OrderRepository
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Transaction(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.Transaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrderRepository_ListByUser(t *testing.T) {
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		r       *OrderRepository
		args    args
		want    []model.Order
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.ListByUser(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.ListByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderRepository.ListByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderRepository_FindByUserAndID(t *testing.T) {
	type args struct {
		userID  uint64
		orderID uint64
	}
	tests := []struct {
		name    string
		r       *OrderRepository
		args    args
		want    model.Order
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.FindByUserAndID(tt.args.userID, tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.FindByUserAndID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderRepository.FindByUserAndID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_attachOrderTimes(t *testing.T) {
	type args struct {
		db     *gorm.DB
		orders []model.Order
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
			if err := attachOrderTimes(tt.args.db, tt.args.orders); (err != nil) != tt.wantErr {
				t.Errorf("attachOrderTimes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_stationTimes(t *testing.T) {
	type args struct {
		db            *gorm.DB
		travelDate    time.Time
		trainID       uint64
		fromStationID uint64
		toStationID   uint64
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		want1   time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := stationTimes(tt.args.db, tt.args.travelDate, tt.args.trainID, tt.args.fromStationID, tt.args.toStationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("stationTimes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stationTimes() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("stationTimes() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_combineDateClock(t *testing.T) {
	type args struct {
		travelDate time.Time
		clock      string
		dayOffset  int
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := combineDateClock(tt.args.travelDate, tt.args.clock, tt.args.dayOffset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("combineDateClock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindInventoryForOrder(t *testing.T) {
	type args struct {
		tx            *gorm.DB
		travelDate    time.Time
		trainID       uint64
		fromStationID uint64
		toStationID   uint64
		seatClassCode string
	}
	tests := []struct {
		name    string
		args    args
		want    OrderInventoryRow
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindInventoryForOrder(tt.args.tx, tt.args.travelDate, tt.args.trainID, tt.args.fromStationID, tt.args.toStationID, tt.args.seatClassCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindInventoryForOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindInventoryForOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
