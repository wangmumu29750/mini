package handler

import (
	"mini-12306/backend/internal/service"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewTrainHandler(t *testing.T) {
	type args struct {
		trains *service.TrainService
	}
	tests := []struct {
		name string
		args args
		want *TrainHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTrainHandler(tt.args.trains); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTrainHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrainHandler_Stations(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *TrainHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Stations(tt.args.c)
		})
	}
}

func TestTrainHandler_Search(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *TrainHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.Search(tt.args.c)
		})
	}
}
