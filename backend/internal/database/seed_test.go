package database

import (
	"mini-12306/backend/internal/model"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestSeedDemoData(t *testing.T) {
	type args struct {
		db *gorm.DB
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
			if err := SeedDemoData(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("SeedDemoData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_seedUsers(t *testing.T) {
	type args struct {
		tx *gorm.DB
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
			if err := seedUsers(tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("seedUsers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_seedStations(t *testing.T) {
	type args struct {
		tx *gorm.DB
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]model.Station
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := seedStations(tt.args.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("seedStations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("seedStations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_seedTrainWithInventory(t *testing.T) {
	type args struct {
		tx       *gorm.DB
		stations map[string]model.Station
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
			if err := seedTrainWithInventory(tt.args.tx, tt.args.stations); (err != nil) != tt.wantErr {
				t.Errorf("seedTrainWithInventory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_demoTrainSeeds(t *testing.T) {
	tests := []struct {
		name string
		want []routeSeed
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := demoTrainSeeds(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("demoTrainSeeds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildStops(t *testing.T) {
	type args struct {
		stations    []string
		miles       []int
		baseMinutes int
	}
	tests := []struct {
		name string
		args args
		want []stopSeed
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildStops(tt.args.stations, tt.args.miles, tt.args.baseMinutes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildStops() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clockText(t *testing.T) {
	type args struct {
		minutes int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clockText(tt.args.minutes); got != tt.want {
				t.Errorf("clockText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildInventories(t *testing.T) {
	type args struct {
		trainType string
		in1       int
		index     int
	}
	tests := []struct {
		name string
		args args
		want []inventorySeed
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildInventories(tt.args.trainType, tt.args.in1, tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildInventories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findOrCreateTrain(t *testing.T) {
	type args struct {
		tx        *gorm.DB
		trainNo   string
		trainType string
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
			got, err := findOrCreateTrain(tt.args.tx, tt.args.trainNo, tt.args.trainType)
			if (err != nil) != tt.wantErr {
				t.Errorf("findOrCreateTrain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findOrCreateTrain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_seedStops(t *testing.T) {
	type args struct {
		tx       *gorm.DB
		trainID  uint64
		stations map[string]model.Station
		seeds    []stopSeed
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
			if err := seedStops(tt.args.tx, tt.args.trainID, tt.args.stations, tt.args.seeds); (err != nil) != tt.wantErr {
				t.Errorf("seedStops() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_seedInventories(t *testing.T) {
	type args struct {
		tx            *gorm.DB
		trainID       uint64
		fromStationID uint64
		toStationID   uint64
		seeds         []inventorySeed
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
			if err := seedInventories(tt.args.tx, tt.args.trainID, tt.args.fromStationID, tt.args.toStationID, tt.args.seeds); (err != nil) != tt.wantErr {
				t.Errorf("seedInventories() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_trainRunsOnDay(t *testing.T) {
	type args struct {
		trainNo string
		day     int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trainRunsOnDay(tt.args.trainNo, tt.args.day); got != tt.want {
				t.Errorf("trainRunsOnDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_demoAvailableCount(t *testing.T) {
	type args struct {
		trainNo       string
		seatClassCode string
		baseAvailable int
		totalCount    int
		day           int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := demoAvailableCount(tt.args.trainNo, tt.args.seatClassCode, tt.args.baseAvailable, tt.args.totalCount, tt.args.day); got != tt.want {
				t.Errorf("demoAvailableCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_minInt(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minInt(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("minInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_routeMileage(t *testing.T) {
	type args struct {
		tx            *gorm.DB
		trainID       uint64
		fromStationID uint64
		toStationID   uint64
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := routeMileage(tt.args.tx, tt.args.trainID, tt.args.fromStationID, tt.args.toStationID)
			if (err != nil) != tt.wantErr {
				t.Errorf("routeMileage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("routeMileage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_seedFareCents(t *testing.T) {
	type args struct {
		mileage       int
		trainType     string
		seatClassCode string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := seedFareCents(tt.args.mileage, tt.args.trainType, tt.args.seatClassCode); got != tt.want {
				t.Errorf("seedFareCents() = %v, want %v", got, tt.want)
			}
		})
	}
}
