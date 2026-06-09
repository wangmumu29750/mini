package service

import (
	"reflect"
	"testing"
)

func TestCalculateFareCentsUsesTrainTypeSeatCoefficient(t *testing.T) {
	tests := []struct {
		name          string
		mileage       int
		trainType     string
		seatClassCode string
		want          int64
	}{
		{name: "high speed second", mileage: 100, trainType: "G", seatClassCode: "SECOND", want: 4600},
		{name: "high speed business", mileage: 100, trainType: "G", seatClassCode: "BUSINESS", want: 13000},
		{name: "conventional hard seat", mileage: 100, trainType: "K", seatClassCode: "HARD_SEAT", want: 1300},
		{name: "emu second sleeper", mileage: 100, trainType: "D", seatClassCode: "SECOND_SLEEPER", want: 4900},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateFareCents(tt.mileage, tt.trainType, tt.seatClassCode)
			if err != nil {
				t.Fatalf("calculateFareCents returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("calculateFareCents() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestEnsureSeatClassAllowedRejectsMismatchedTrainType(t *testing.T) {
	if err := ensureSeatClassAllowed("G", "HARD_SEAT"); err == nil {
		t.Fatal("expected G train to reject hard seat")
	}
	if err := ensureSeatClassAllowed("Z", "SOFT_SLEEPER"); err != nil {
		t.Fatalf("expected Z train to accept soft sleeper: %v", err)
	}
}

func Test_conventionalSeatRules(t *testing.T) {
	tests := []struct {
		name string
		want []seatClassRule
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := conventionalSeatRules(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("conventionalSeatRules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizeTrainType(t *testing.T) {
	type args struct {
		value string
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
			if got := normalizeTrainType(tt.args.value); got != tt.want {
				t.Errorf("normalizeTrainType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_supportedTrainType(t *testing.T) {
	type args struct {
		value string
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
			if got := supportedTrainType(tt.args.value); got != tt.want {
				t.Errorf("supportedTrainType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_supportedSeatClasses(t *testing.T) {
	type args struct {
		trainType string
	}
	tests := []struct {
		name string
		args args
		want []seatClassRule
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := supportedSeatClasses(tt.args.trainType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("supportedSeatClasses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_trainTypeFromTrainNo(t *testing.T) {
	type args struct {
		trainNo string
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
			if got := trainTypeFromTrainNo(tt.args.trainNo); got != tt.want {
				t.Errorf("trainTypeFromTrainNo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ensureSeatClassAllowed(t *testing.T) {
	type args struct {
		trainType     string
		seatClassCode string
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
			if err := ensureSeatClassAllowed(tt.args.trainType, tt.args.seatClassCode); (err != nil) != tt.wantErr {
				t.Errorf("ensureSeatClassAllowed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_calculateFareCents(t *testing.T) {
	type args struct {
		mileage       int
		trainType     string
		seatClassCode string
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
			got, err := calculateFareCents(tt.args.mileage, tt.args.trainType, tt.args.seatClassCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("calculateFareCents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calculateFareCents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_seatClassName(t *testing.T) {
	type args struct {
		code string
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
			if got := seatClassName(tt.args.code); got != tt.want {
				t.Errorf("seatClassName() = %v, want %v", got, tt.want)
			}
		})
	}
}
