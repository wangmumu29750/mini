package service

import (
	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/model"
	"mini-12306/backend/internal/repository"
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestNewTicketService(t *testing.T) {
	type args struct {
		tickets *repository.TicketRepository
	}
	tests := []struct {
		name string
		args args
		want *TicketService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTicketService(tt.args.tickets); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTicketService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTicketService_List(t *testing.T) {
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		s       *TicketService
		args    args
		want    []dto.TicketResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.List(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("TicketService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TicketService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTicketService_Detail(t *testing.T) {
	type args struct {
		userID   uint64
		ticketID uint64
	}
	tests := []struct {
		name    string
		s       *TicketService
		args    args
		want    dto.TicketResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Detail(tt.args.userID, tt.args.ticketID)
			if (err != nil) != tt.wantErr {
				t.Errorf("TicketService.Detail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TicketService.Detail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTicketService_ChangeOptions(t *testing.T) {
	type args struct {
		userID   uint64
		ticketID uint64
		query    dto.ChangeOptionsQuery
	}
	tests := []struct {
		name    string
		s       *TicketService
		args    args
		want    dto.ChangeOptionsResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ChangeOptions(tt.args.userID, tt.args.ticketID, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("TicketService.ChangeOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TicketService.ChangeOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTicketService_Refund(t *testing.T) {
	type args struct {
		userID   uint64
		ticketID uint64
		req      dto.RefundTicketRequest
	}
	tests := []struct {
		name    string
		s       *TicketService
		args    args
		want    dto.RefundTicketResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Refund(tt.args.userID, tt.args.ticketID, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("TicketService.Refund() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TicketService.Refund() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTicketService_Change(t *testing.T) {
	type args struct {
		userID   uint64
		ticketID uint64
		req      dto.ChangeTicketRequest
	}
	tests := []struct {
		name    string
		s       *TicketService
		args    args
		want    dto.ChangeTicketResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Change(tt.args.userID, tt.args.ticketID, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("TicketService.Change() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TicketService.Change() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trainSearchRowsToResponses(t *testing.T) {
	type args struct {
		date time.Time
		rows []repository.TrainSearchRow
	}
	tests := []struct {
		name string
		args args
		want []dto.TrainSearchItemResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trainSearchRowsToResponses(tt.args.date, tt.args.rows); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("trainSearchRowsToResponses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ticketResponse(t *testing.T) {
	type args struct {
		ticket model.Ticket
	}
	tests := []struct {
		name string
		args args
		want dto.TicketResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ticketResponse(tt.args.ticket); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ticketResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_departTimeForTicket(t *testing.T) {
	type args struct {
		tx     *gorm.DB
		ticket model.Ticket
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := departTimeForTicket(tt.args.tx, tt.args.ticket)
			if (err != nil) != tt.wantErr {
				t.Errorf("departTimeForTicket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("departTimeForTicket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeCoachNo(t *testing.T) {
	type args struct {
		sequence int
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
			if got := makeCoachNo(tt.args.sequence); got != tt.want {
				t.Errorf("makeCoachNo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeSeatNo(t *testing.T) {
	type args struct {
		sequence int
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
			if got := makeSeatNo(tt.args.sequence); got != tt.want {
				t.Errorf("makeSeatNo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_maskIDCardNo(t *testing.T) {
	type args struct {
		value string
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
			if got := maskIDCardNo(tt.args.value); got != tt.want {
				t.Errorf("maskIDCardNo() = %v, want %v", got, tt.want)
			}
		})
	}
}
