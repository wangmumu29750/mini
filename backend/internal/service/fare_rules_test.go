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
		{
			name: "returns 5 conventional seat rules",
			want: []seatClassRule{
				{Code: "DELUXE_SOFT_SLEEPER", Name: "高级软卧", Coefficient: 4.0},
				{Code: "SOFT_SLEEPER", Name: "软卧", Coefficient: 3.0},
				{Code: "HARD_SLEEPER", Name: "硬卧", Coefficient: 2.0},
				{Code: "HARD_SEAT", Name: "硬座", Coefficient: 1.0},
				{Code: "NO_SEAT", Name: "无座", Coefficient: 1.0},
			},
		},
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
		{name: "G train", args: args{value: "G"}, want: "G"},
		{name: "lowercase g", args: args{value: "g"}, want: "G"},
		{name: "full G train prefix", args: args{value: "G101"}, want: "G"},
		{name: "D train", args: args{value: "D"}, want: "D"},
		{name: "K train", args: args{value: "K"}, want: "K"},
		{name: "empty string", args: args{value: ""}, want: ""},
		{name: "whitespace only", args: args{value: "  "}, want: ""},
		{name: "unknown type X", args: args{value: "X"}, want: "X"},
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
		{name: "G supported", args: args{value: "G"}, want: true},
		{name: "D supported", args: args{value: "D"}, want: true},
		{name: "C supported", args: args{value: "C"}, want: true},
		{name: "Z supported", args: args{value: "Z"}, want: true},
		{name: "T supported", args: args{value: "T"}, want: true},
		{name: "K supported", args: args{value: "K"}, want: true},
		{name: "X not supported", args: args{value: "X"}, want: false},
		{name: "empty not supported", args: args{value: ""}, want: false},
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
		{
			name: "G train seat classes",
			args: args{trainType: "G"},
			want: []seatClassRule{
				{Code: "BUSINESS", Name: "商务座", Coefficient: 10.0},
				{Code: "FIRST", Name: "一等座", Coefficient: 5.8},
				{Code: "SECOND", Name: "二等座", Coefficient: 3.5},
			},
		},
		{
			name: "D train seat classes",
			args: args{trainType: "D"},
			want: []seatClassRule{
				{Code: "FIRST_SLEEPER", Name: "一等卧", Coefficient: 5.0},
				{Code: "SECOND_SLEEPER", Name: "二等卧", Coefficient: 3.8},
				{Code: "SECOND", Name: "二等座", Coefficient: 3.0},
			},
		},
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
		{name: "G101", args: args{trainNo: "G101"}, want: "G"},
		{name: "D501", args: args{trainNo: "D501"}, want: "D"},
		{name: "K701", args: args{trainNo: "K701"}, want: "K"},
		{name: "lowercase", args: args{trainNo: "g137"}, want: "G"},
		{name: "empty", args: args{trainNo: ""}, want: ""},
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
		{name: "G business allowed", args: args{trainType: "G", seatClassCode: "BUSINESS"}, wantErr: false},
		{name: "G second allowed", args: args{trainType: "G", seatClassCode: "SECOND"}, wantErr: false},
		{name: "G hard seat rejected", args: args{trainType: "G", seatClassCode: "HARD_SEAT"}, wantErr: true},
		{name: "D second sleeper allowed", args: args{trainType: "D", seatClassCode: "SECOND_SLEEPER"}, wantErr: false},
		{name: "D business rejected", args: args{trainType: "D", seatClassCode: "BUSINESS"}, wantErr: true},
		{name: "K hard seat allowed", args: args{trainType: "K", seatClassCode: "HARD_SEAT"}, wantErr: false},
		{name: "K soft sleeper allowed", args: args{trainType: "K", seatClassCode: "SOFT_SLEEPER"}, wantErr: false},
		{name: "unknown train type rejected", args: args{trainType: "X", seatClassCode: "SECOND"}, wantErr: true},
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
		{name: "G second 100km", args: args{mileage: 100, trainType: "G", seatClassCode: "SECOND"}, want: 4600, wantErr: false},
		{name: "G business 100km", args: args{mileage: 100, trainType: "G", seatClassCode: "BUSINESS"}, want: 13000, wantErr: false},
		{name: "D second sleeper 200km", args: args{mileage: 200, trainType: "D", seatClassCode: "SECOND_SLEEPER"}, want: 9900, wantErr: false},
		{name: "K hard seat 100km", args: args{mileage: 100, trainType: "K", seatClassCode: "HARD_SEAT"}, want: 1300, wantErr: false},
		{name: "zero mileage error", args: args{mileage: 0, trainType: "G", seatClassCode: "SECOND"}, want: 0, wantErr: true},
		{name: "negative mileage error", args: args{mileage: -10, trainType: "G", seatClassCode: "SECOND"}, want: 0, wantErr: true},
		{name: "invalid seat for train type", args: args{mileage: 100, trainType: "G", seatClassCode: "HARD_SEAT"}, want: 0, wantErr: true},
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
		{name: "BUSINESS", args: args{code: "BUSINESS"}, want: "商务座"},
		{name: "FIRST", args: args{code: "FIRST"}, want: "一等座"},
		{name: "SECOND", args: args{code: "SECOND"}, want: "二等座"},
		{name: "HARD_SEAT", args: args{code: "HARD_SEAT"}, want: "硬座"},
		{name: "HARD_SLEEPER", args: args{code: "HARD_SLEEPER"}, want: "硬卧"},
		{name: "SOFT_SLEEPER", args: args{code: "SOFT_SLEEPER"}, want: "软卧"},
		{name: "NO_SEAT", args: args{code: "NO_SEAT"}, want: "无座"},
		{name: "lowercase input", args: args{code: "business"}, want: "商务座"},
		{name: "unknown code returns itself", args: args{code: "UNKNOWN"}, want: "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := seatClassName(tt.args.code); got != tt.want {
				t.Errorf("seatClassName() = %v, want %v", got, tt.want)
			}
		})
	}
}
