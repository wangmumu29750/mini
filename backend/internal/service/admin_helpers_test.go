package service

import (
	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/model"
	"mini-12306/backend/internal/repository"
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm"
)

func Test_normalizePage(t *testing.T) {
	type args struct {
		page     int
		pageSize int
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 int
	}{
		{name: "valid page and size", args: args{page: 1, pageSize: 20}, want: 1, want1: 20},
		{name: "zero page defaults to 1", args: args{page: 0, pageSize: 20}, want: 1, want1: 20},
		{name: "negative page defaults to 1", args: args{page: -1, pageSize: 20}, want: 1, want1: 20},
		{name: "zero page size defaults to 20", args: args{page: 1, pageSize: 0}, want: 1, want1: 20},
		{name: "negative page size defaults to 20", args: args{page: 1, pageSize: -5}, want: 1, want1: 20},
		{name: "page size capped at 100", args: args{page: 1, pageSize: 200}, want: 1, want1: 100},
		{name: "page size exactly 100", args: args{page: 2, pageSize: 100}, want: 2, want1: 100},
		{name: "both zero use defaults", args: args{page: 0, pageSize: 0}, want: 1, want1: 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := normalizePage(tt.args.page, tt.args.pageSize)
			if got != tt.want {
				t.Errorf("normalizePage() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("normalizePage() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_stationStatus(t *testing.T) {
	type args struct {
		value      string
		allowEmpty bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stationStatus(tt.args.value, tt.args.allowEmpty)
			if (err != nil) != tt.wantErr {
				t.Errorf("stationStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("stationStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trainStatus(t *testing.T) {
	type args struct {
		value      string
		allowEmpty bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := trainStatus(tt.args.value, tt.args.allowEmpty)
			if (err != nil) != tt.wantErr {
				t.Errorf("trainStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("trainStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_optionalDate(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    *time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := optionalDate(tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("optionalDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("optionalDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ensureTrainAndStations(t *testing.T) {
	type args struct {
		tx            *gorm.DB
		trainID       uint64
		fromStationID uint64
		toStationID   uint64
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
			got, err := ensureTrainAndStations(tt.args.tx, tt.args.trainID, tt.args.fromStationID, tt.args.toStationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ensureTrainAndStations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ensureTrainAndStations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lowestPriceForTrain(t *testing.T) {
	type args struct {
		db      *gorm.DB
		trainID uint64
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
			got, err := lowestPriceForTrain(tt.args.db, tt.args.trainID)
			if (err != nil) != tt.wantErr {
				t.Errorf("lowestPriceForTrain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("lowestPriceForTrain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_startOfLocalDay(t *testing.T) {
	type args struct {
		value time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := startOfLocalDay(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("startOfLocalDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stationResponse(t *testing.T) {
	type args struct {
		station model.Station
	}
	tests := []struct {
		name string
		args args
		want dto.AdminStationResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stationResponse(tt.args.station); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stationResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trainResponse(t *testing.T) {
	type args struct {
		train model.Train
	}
	tests := []struct {
		name string
		args args
		want dto.AdminTrainResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trainResponse(tt.args.train); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("trainResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stopResponse(t *testing.T) {
	type args struct {
		stop model.TrainStop
	}
	tests := []struct {
		name string
		args args
		want dto.TrainStopResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stopResponse(tt.args.stop); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stopResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inventoryRowResponse(t *testing.T) {
	type args struct {
		row repository.InventoryRow
	}
	tests := []struct {
		name string
		args args
		want dto.InventoryResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inventoryRowResponse(tt.args.row); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inventoryRowResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
