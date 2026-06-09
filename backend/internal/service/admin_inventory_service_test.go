package service

import (
	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/model"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestAdminService_ListInventories(t *testing.T) {
	type args struct {
		query dto.InventoryListQuery
	}
	tests := []struct {
		name    string
		s       *AdminService
		args    args
		want    dto.PageResponse[dto.InventoryResponse]
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ListInventories(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminService.ListInventories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminService.ListInventories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminService_SaveInventory(t *testing.T) {
	type args struct {
		req dto.SaveInventoryRequest
	}
	tests := []struct {
		name    string
		s       *AdminService
		args    args
		want    dto.InventoryResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.SaveInventory(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminService.SaveInventory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminService.SaveInventory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminService_QuoteStats(t *testing.T) {
	type args struct {
		query dto.InventoryQuoteStatsQuery
	}
	tests := []struct {
		name    string
		s       *AdminService
		args    args
		want    dto.InventoryQuoteStatsResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.QuoteStats(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminService.QuoteStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminService.QuoteStats() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminService_FlowInventory(t *testing.T) {
	type args struct {
		req dto.InventoryFlowRequest
	}
	tests := []struct {
		name    string
		s       *AdminService
		args    args
		want    dto.InventoryFlowResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.FlowInventory(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminService.FlowInventory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminService.FlowInventory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inventoryFromRequest(t *testing.T) {
	type args struct {
		req dto.SaveInventoryRequest
	}
	tests := []struct {
		name    string
		args    args
		want    model.Inventory
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := inventoryFromRequest(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("inventoryFromRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inventoryFromRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fareForInventory(t *testing.T) {
	type args struct {
		tx        *gorm.DB
		inventory model.Inventory
		trainType string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fareForInventory(tt.args.tx, tt.args.inventory, tt.args.trainType)
			if (err != nil) != tt.wantErr {
				t.Errorf("fareForInventory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("fareForInventory() = %v, want %v", got, tt.want)
			}
		})
	}
}
