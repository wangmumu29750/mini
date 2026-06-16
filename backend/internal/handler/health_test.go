package handler

import (
	"mini-12306/backend/internal/config"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestNewHealthHandler(t *testing.T) {
	type args struct {
		cfg config.Config
		db  *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *HealthHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHealthHandler(tt.args.cfg, tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHealthHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealthHandler_Check(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *HealthHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Check(tt.args.c)
		})
	}
}
