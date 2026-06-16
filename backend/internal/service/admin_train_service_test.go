package service

import (
	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/model"
	"reflect"
	"testing"
)

func TestAdminService_ListTrains(t *testing.T) {
	type args struct {
		query dto.AdminTrainListQuery
	}
	tests := []struct {
		name    string
		s       *AdminService
		args    args
		want    dto.PageResponse[dto.AdminTrainResponse]
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ListTrains(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminService.ListTrains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminService.ListTrains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminService_CreateTrain(t *testing.T) {
	type args struct {
		req dto.SaveTrainRequest
	}
	tests := []struct {
		name    string
		s       *AdminService
		args    args
		want    dto.AdminTrainResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.CreateTrain(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminService.CreateTrain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminService.CreateTrain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminService_UpdateTrain(t *testing.T) {
	type args struct {
		id  uint64
		req dto.SaveTrainRequest
	}
	tests := []struct {
		name    string
		s       *AdminService
		args    args
		want    dto.AdminTrainResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.UpdateTrain(tt.args.id, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminService.UpdateTrain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminService.UpdateTrain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminService_DeleteTrain(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		s       *AdminService
		args    args
		want    dto.AdminTrainResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.DeleteTrain(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminService.DeleteTrain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminService.DeleteTrain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminService_SellableStats(t *testing.T) {
	type args struct {
		query dto.SellableTrainStatsQuery
	}
	tests := []struct {
		name    string
		s       *AdminService
		args    args
		want    []dto.SellableTrainStatItem
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.SellableStats(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminService.SellableStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminService.SellableStats() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trainFromRequest(t *testing.T) {
	type args struct {
		req dto.SaveTrainRequest
	}
	tests := []struct {
		name    string
		args    args
		want    model.Train
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := trainFromRequest(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("trainFromRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("trainFromRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
