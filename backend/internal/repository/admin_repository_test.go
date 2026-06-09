package repository

import (
	"mini-12306/backend/internal/model"
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestNewAdminRepository(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *AdminRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAdminRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAdminRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminRepository_Transaction(t *testing.T) {
	type args struct {
		fn func(tx *gorm.DB) error
	}
	tests := []struct {
		name    string
		r       *AdminRepository
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Transaction(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("AdminRepository.Transaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdminRepository_DB(t *testing.T) {
	tests := []struct {
		name string
		r    *AdminRepository
		want *gorm.DB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.DB(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminRepository.DB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminRepository_ListStations(t *testing.T) {
	type args struct {
		page     int
		pageSize int
		status   string
	}
	tests := []struct {
		name    string
		r       *AdminRepository
		args    args
		want    []model.Station
		want1   int64
		want2   int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := tt.r.ListStations(tt.args.page, tt.args.pageSize, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminRepository.ListStations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminRepository.ListStations() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AdminRepository.ListStations() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("AdminRepository.ListStations() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestAdminRepository_FindStation(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		r       *AdminRepository
		args    args
		want    model.Station
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.FindStation(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminRepository.FindStation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminRepository.FindStation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminRepository_CreateStation(t *testing.T) {
	type args struct {
		station *model.Station
	}
	tests := []struct {
		name    string
		r       *AdminRepository
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.CreateStation(tt.args.station); (err != nil) != tt.wantErr {
				t.Errorf("AdminRepository.CreateStation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdminRepository_SaveStation(t *testing.T) {
	type args struct {
		station *model.Station
	}
	tests := []struct {
		name    string
		r       *AdminRepository
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.SaveStation(tt.args.station); (err != nil) != tt.wantErr {
				t.Errorf("AdminRepository.SaveStation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdminRepository_ListTrains(t *testing.T) {
	type args struct {
		page     int
		pageSize int
		status   string
		trainNo  string
	}
	tests := []struct {
		name    string
		r       *AdminRepository
		args    args
		want    []model.Train
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.r.ListTrains(tt.args.page, tt.args.pageSize, tt.args.status, tt.args.trainNo)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminRepository.ListTrains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminRepository.ListTrains() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AdminRepository.ListTrains() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAdminRepository_FindTrain(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		r       *AdminRepository
		args    args
		want    model.Train
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.FindTrain(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminRepository.FindTrain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminRepository.FindTrain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminRepository_CreateTrain(t *testing.T) {
	type args struct {
		train *model.Train
	}
	tests := []struct {
		name    string
		r       *AdminRepository
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.CreateTrain(tt.args.train); (err != nil) != tt.wantErr {
				t.Errorf("AdminRepository.CreateTrain() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdminRepository_SaveTrain(t *testing.T) {
	type args struct {
		train *model.Train
	}
	tests := []struct {
		name    string
		r       *AdminRepository
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.SaveTrain(tt.args.train); (err != nil) != tt.wantErr {
				t.Errorf("AdminRepository.SaveTrain() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAdminRepository_ListStops(t *testing.T) {
	type args struct {
		trainID uint64
	}
	tests := []struct {
		name    string
		r       *AdminRepository
		args    args
		want    []model.TrainStop
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.ListStops(tt.args.trainID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminRepository.ListStops() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminRepository.ListStops() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAdminRepository_ListInventories(t *testing.T) {
	type args struct {
		page          int
		pageSize      int
		trainID       uint64
		seatClassCode string
		travelDate    *time.Time
	}
	tests := []struct {
		name    string
		r       *AdminRepository
		args    args
		want    []InventoryRow
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.r.ListInventories(tt.args.page, tt.args.pageSize, tt.args.trainID, tt.args.seatClassCode, tt.args.travelDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("AdminRepository.ListInventories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AdminRepository.ListInventories() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AdminRepository.ListInventories() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
