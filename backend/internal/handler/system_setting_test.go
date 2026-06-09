package handler

import (
	"mini-12306/backend/internal/service"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewSystemSettingHandler(t *testing.T) {
	type args struct {
		settings *service.SystemSettingService
	}
	tests := []struct {
		name string
		args args
		want *SystemSettingHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSystemSettingHandler(tt.args.settings); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSystemSettingHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemSettingHandler_List(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *SystemSettingHandler
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

func TestSystemSettingHandler_Update(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *SystemSettingHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Update(tt.args.c)
		})
	}
}
