package service

import (
	"mini-12306/backend/internal/dto"
	"reflect"
	"testing"
)

func TestAdminService_ListStops(t *testing.T) {
	type args struct {
		trainID uint64
	}
	tests := []struct {
		name    string
		s       *AdminService
		args    args
		want    []dto.TrainStopResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ListStops(tt.args.trainID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminService.ListStops() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminService.ListStops() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminService_SaveStops(t *testing.T) {
	type args struct {
		trainID uint64
		req     dto.SaveTrainStopsRequest
	}
	tests := []struct {
		name    string
		s       *AdminService
		args    args
		want    []dto.TrainStopResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.SaveStops(tt.args.trainID, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminService.SaveStops() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminService.SaveStops() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateStopItems(t *testing.T) {
	type args struct {
		stops []dto.SaveTrainStopItem
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
			if err := validateStopItems(tt.args.stops); (err != nil) != tt.wantErr {
				t.Errorf("validateStopItems() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
