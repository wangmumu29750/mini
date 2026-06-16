package handler

import (
	"mini-12306/backend/internal/service"
	pkgauth "mini-12306/backend/pkg/auth"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewOrderHandler(t *testing.T) {
	type args struct {
		orders *service.OrderService
	}
	tests := []struct {
		name string
		args args
		want *OrderHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrderHandler(tt.args.orders); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrderHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderHandler_Create(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *OrderHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Create(tt.args.c)
		})
	}
}

func TestOrderHandler_ClerkCreate(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *OrderHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.ClerkCreate(tt.args.c)
		})
	}
}

func TestOrderHandler_List(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *OrderHandler
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

func TestOrderHandler_Detail(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *OrderHandler
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

func TestOrderHandler_Pay(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *OrderHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Pay(tt.args.c)
		})
	}
}

func TestOrderHandler_Cancel(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *OrderHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Cancel(tt.args.c)
		})
	}
}

func Test_currentPrincipal(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name  string
		args  args
		want  pkgauth.Principal
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := currentPrincipal(tt.args.c)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("currentPrincipal() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("currentPrincipal() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
