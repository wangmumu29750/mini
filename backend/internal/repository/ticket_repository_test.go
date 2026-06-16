package repository

import (
	"mini-12306/backend/internal/model"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestNewTicketRepository(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *TicketRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTicketRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTicketRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTicketRepository_Transaction(t *testing.T) {
	type args struct {
		fn func(tx *gorm.DB) error
	}
	tests := []struct {
		name    string
		r       *TicketRepository
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Transaction(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("TicketRepository.Transaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTicketRepository_ListByUser(t *testing.T) {
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		r       *TicketRepository
		args    args
		want    []model.Ticket
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.ListByUser(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("TicketRepository.ListByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TicketRepository.ListByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTicketRepository_FindByUserAndID(t *testing.T) {
	type args struct {
		userID   uint64
		ticketID uint64
	}
	tests := []struct {
		name    string
		r       *TicketRepository
		args    args
		want    model.Ticket
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.FindByUserAndID(tt.args.userID, tt.args.ticketID)
			if (err != nil) != tt.wantErr {
				t.Errorf("TicketRepository.FindByUserAndID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TicketRepository.FindByUserAndID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_attachTicketTimes(t *testing.T) {
	type args struct {
		db      *gorm.DB
		tickets []model.Ticket
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
			if err := attachTicketTimes(tt.args.db, tt.args.tickets); (err != nil) != tt.wantErr {
				t.Errorf("attachTicketTimes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
