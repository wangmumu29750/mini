package handler

import (
	"mini-12306/backend/internal/service"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewTicketHandler(t *testing.T) {
	type args struct {
		tickets *service.TicketService
	}
	tests := []struct {
		name string
		args args
		want *TicketHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTicketHandler(tt.args.tickets); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTicketHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTicketHandler_List(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *TicketHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.List(tt.args.c)
		})
	}
}

func TestTicketHandler_Detail(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *TicketHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Detail(tt.args.c)
		})
	}
}

func TestTicketHandler_Refund(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *TicketHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Refund(tt.args.c)
		})
	}
}

func TestTicketHandler_ChangeOptions(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *TicketHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.ChangeOptions(tt.args.c)
		})
	}
}

func TestTicketHandler_Change(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *TicketHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Change(tt.args.c)
		})
	}
}
