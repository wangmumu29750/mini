package repository

import (
	"mini-12306/backend/internal/model"
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestNewTrainRepository(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *TrainRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTrainRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTrainRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrainRepository_ListActiveStations(t *testing.T) {
	tests := []struct {
		name    string
		r       *TrainRepository
		want    []model.Station
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.ListActiveStations()
			if (err != nil) != tt.wantErr {
				t.Errorf("TrainRepository.ListActiveStations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TrainRepository.ListActiveStations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrainRepository_SearchAvailableTrains(t *testing.T) {
	type args struct {
		date          time.Time
		fromStationID uint64
		toStationID   uint64
	}
	tests := []struct {
		name    string
		r       *TrainRepository
		args    args
		want    []TrainSearchRow
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.SearchAvailableTrains(tt.args.date, tt.args.fromStationID, tt.args.toStationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("TrainRepository.SearchAvailableTrains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TrainRepository.SearchAvailableTrains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchAvailableTrains(t *testing.T) {
	type args struct {
		db            *gorm.DB
		date          time.Time
		fromStationID uint64
		toStationID   uint64
	}
	tests := []struct {
		name    string
		args    args
		want    []TrainSearchRow
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SearchAvailableTrains(tt.args.db, tt.args.date, tt.args.fromStationID, tt.args.toStationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchAvailableTrains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchAvailableTrains() = %v, want %v", got, tt.want)
			}
		})
	}
}
