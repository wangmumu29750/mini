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
		name      string
		minTrains int
	}{
		{name: "returns at least two explicit trains", minTrains: 2},
		{name: "includes G101", minTrains: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := demoTrainSeeds()
			if tt.name == "returns at least two explicit trains" && len(got) < tt.minTrains {
				t.Errorf("demoTrainSeeds() returned %d trains, want at least %d", len(got), tt.minTrains)
			}
			if tt.name == "includes G101" {
				found := false
				for _, s := range got {
					if s.TrainNo == "G101" {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("demoTrainSeeds() should include G101")
				}
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
		{
			name: "two station route",
			args: args{
				stations:    []string{"BJN", "SHH"},
				miles:       []int{0, 1318},
				baseMinutes: 480,
			},
			want: []stopSeed{
				{StationCode: "BJN", Order: 1, ArriveClock: "", DepartClock: "08:00:00", Mileage: 0},
				{StationCode: "SHH", Order: 2, ArriveClock: "08:32:00", DepartClock: "", Mileage: 1318},
			},
		},
		{
			name: "three station route",
			args: args{
				stations:    []string{"BJN", "TJN", "SHH"},
				miles:       []int{0, 122, 1318},
				baseMinutes: 480,
			},
			want: []stopSeed{
				{StationCode: "BJN", Order: 1, ArriveClock: "", DepartClock: "08:00:00", Mileage: 0},
				{StationCode: "TJN", Order: 2, ArriveClock: "08:32:00", DepartClock: "08:36:00", Mileage: 122},
				{StationCode: "SHH", Order: 3, ArriveClock: "09:08:00", DepartClock: "", Mileage: 1318},
			},
		},
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
		{name: "midnight", args: args{minutes: 0}, want: "00:00:00"},
		{name: "8:00 am", args: args{minutes: 480}, want: "08:00:00"},
		{name: "noon", args: args{minutes: 720}, want: "12:00:00"},
		{name: "23:59", args: args{minutes: 1439}, want: "23:59:00"},
		{name: "wraps around 24h", args: args{minutes: 1500}, want: "01:00:00"},
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
		{
			name: "G train type",
			args: args{trainType: "G", in1: 0, index: 0},
			want: []inventorySeed{
				{SeatClassCode: "SECOND", TotalCount: 120, AvailableCount: 72},
				{SeatClassCode: "FIRST", TotalCount: 48, AvailableCount: 24},
				{SeatClassCode: "BUSINESS", TotalCount: 12, AvailableCount: 6},
			},
		},
		{
			name: "D train type",
			args: args{trainType: "D", in1: 0, index: 0},
			want: []inventorySeed{
				{SeatClassCode: "SECOND", TotalCount: 120, AvailableCount: 72},
				{SeatClassCode: "SECOND_SLEEPER", TotalCount: 64, AvailableCount: 30},
				{SeatClassCode: "FIRST_SLEEPER", TotalCount: 32, AvailableCount: 16},
			},
		},
		{
			name: "Z train type conventional",
			args: args{trainType: "Z", in1: 0, index: 0},
			want: []inventorySeed{
				{SeatClassCode: "HARD_SEAT", TotalCount: 120, AvailableCount: 70},
				{SeatClassCode: "NO_SEAT", TotalCount: 80, AvailableCount: 55},
				{SeatClassCode: "HARD_SLEEPER", TotalCount: 72, AvailableCount: 36},
				{SeatClassCode: "SOFT_SLEEPER", TotalCount: 32, AvailableCount: 14},
				{SeatClassCode: "DELUXE_SOFT_SLEEPER", TotalCount: 12, AvailableCount: 5},
			},
		},
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
		{name: "always returns true day 0", args: args{trainNo: "G101", day: 0}, want: true},
		{name: "always returns true day 3", args: args{trainNo: "D501", day: 3}, want: true},
		{name: "always returns true day 6", args: args{trainNo: "K701", day: 6}, want: true},
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
		{name: "zero total count returns zero", args: args{trainNo: "G101", seatClassCode: "SECOND", baseAvailable: 10, totalCount: 0, day: 0}, want: 0},
		{name: "zero available day 0 first class (day%3=0 AND seatWeight=0)", args: args{trainNo: "G101", seatClassCode: "FIRST", baseAvailable: 10, totalCount: 32, day: 0}, want: 0},
		{name: "zero available day 1 business (day%4=1 AND BUSINESS)", args: args{trainNo: "G101", seatClassCode: "BUSINESS", baseAvailable: 3, totalCount: 8, day: 1}, want: 0},
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
		{name: "a less than b", args: args{a: 1, b: 5}, want: 1},
		{name: "b less than a", args: args{a: 10, b: 3}, want: 3},
		{name: "equal values", args: args{a: 7, b: 7}, want: 7},
		{name: "negative a", args: args{a: -5, b: 2}, want: -5},
		{name: "both negative", args: args{a: -3, b: -1}, want: -3},
		{name: "zero and positive", args: args{a: 0, b: 10}, want: 0},
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
		{name: "G second class 100km", args: args{mileage: 100, trainType: "G", seatClassCode: "SECOND"}, want: 4600},
		{name: "G business class 100km", args: args{mileage: 100, trainType: "G", seatClassCode: "BUSINESS"}, want: 13000},
		{name: "D second sleeper 100km", args: args{mileage: 100, trainType: "D", seatClassCode: "SECOND_SLEEPER"}, want: 4900},
		{name: "K hard seat 100km", args: args{mileage: 100, trainType: "K", seatClassCode: "HARD_SEAT"}, want: 1300},
		{name: "Z soft sleeper 100km", args: args{mileage: 100, trainType: "Z", seatClassCode: "SOFT_SLEEPER"}, want: 3900},
		{name: "unknown train type defaults coefficient to 1", args: args{mileage: 100, trainType: "X", seatClassCode: "HARD_SEAT"}, want: 1300},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := seedFareCents(tt.args.mileage, tt.args.trainType, tt.args.seatClassCode); got != tt.want {
				t.Errorf("seedFareCents() = %v, want %v", got, tt.want)
			}
		})
	}
}
