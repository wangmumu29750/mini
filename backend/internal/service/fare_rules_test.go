package service

import "testing"

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
