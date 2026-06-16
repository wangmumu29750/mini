package handler

import (
	"mini-12306/backend/internal/service"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewAdminHandler(t *testing.T) {
	type args struct {
		admin *service.AdminService
	}
	tests := []struct {
		name string
		args args
		want *AdminHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAdminHandler(tt.args.admin); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAdminHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminHandler_PublicStations(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.PublicStations(tt.args.c)
		})
	}
}

func TestAdminHandler_ListStations(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.ListStations(tt.args.c)
		})
	}
}

func TestAdminHandler_CreateStation(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.CreateStation(tt.args.c)
		})
	}
}

func TestAdminHandler_UpdateStation(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.UpdateStation(tt.args.c)
		})
	}
}

func TestAdminHandler_DisableStation(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.DisableStation(tt.args.c)
		})
	}
}

func TestAdminHandler_ListTrains(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.ListTrains(tt.args.c)
		})
	}
}

func TestAdminHandler_CreateTrain(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.CreateTrain(tt.args.c)
		})
	}
}

func TestAdminHandler_UpdateTrain(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.UpdateTrain(tt.args.c)
		})
	}
}

func TestAdminHandler_DeleteTrain(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.DeleteTrain(tt.args.c)
		})
	}
}

func TestAdminHandler_SellableStats(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.SellableStats(tt.args.c)
		})
	}
}

func TestAdminHandler_ListStops(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.ListStops(tt.args.c)
		})
	}
}

func TestAdminHandler_SaveStops(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.SaveStops(tt.args.c)
		})
	}
}

func TestAdminHandler_ListInventories(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.ListInventories(tt.args.c)
		})
	}
}

func TestAdminHandler_SaveInventory(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.SaveInventory(tt.args.c)
		})
	}
}

func TestAdminHandler_QuoteStats(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.QuoteStats(tt.args.c)
		})
	}
}

func TestAdminHandler_FlowInventory(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		h    *AdminHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.h.FlowInventory(tt.args.c)
		})
	}
}

func Test_parseIDParam(t *testing.T) {
	type args struct {
		c       *gin.Context
		name    string
		message string
	}
	tests := []struct {
		name  string
		args  args
		want  uint64
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parseIDParam(tt.args.c, tt.args.name, tt.args.message)
			if got != tt.want {
				t.Errorf("parseIDParam() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseIDParam() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
