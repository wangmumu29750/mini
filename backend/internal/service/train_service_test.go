package service

import (
	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/repository"
	"reflect"
	"testing"
)

func TestNewTrainService(t *testing.T) {
	type args struct {
		trains *repository.TrainRepository
	}
	tests := []struct {
		name string
		args args
		want *TrainService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTrainService(tt.args.trains); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTrainService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrainService_ListStations(t *testing.T) {
	tests := []struct {
		name    string
		s       *TrainService
		want    []dto.StationResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ListStations()
			if (err != nil) != tt.wantErr {
				t.Errorf("TrainService.ListStations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TrainService.ListStations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrainService_Search(t *testing.T) {
	type args struct {
		query dto.TrainSearchQuery
	}
	tests := []struct {
		name    string
		s       *TrainService
		args    args
		want    []dto.TrainSearchItemResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Search(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("TrainService.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TrainService.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}
