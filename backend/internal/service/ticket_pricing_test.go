package service

import "testing"

func TestCalculateTicketPriceAppliesTicketTypeDiscount(t *testing.T) {
	tests := []struct {
		name        string
		basePrice   int64
		trainType   string
		seatType    string
		ticketType  string
		wantPrice   int64
		wantErr     bool
		description string
	}{
		{
			name:        "adult ticket keeps original price",
			basePrice:   10000,
			trainType:   "G",
			seatType:    "SECOND",
			ticketType:  "ADULT",
			wantPrice:   10000,
			description: "Adult tickets should not receive any discount.",
		},
		{
			name:        "student high speed ticket gets 75 percent price",
			basePrice:   10000,
			trainType:   "G",
			seatType:    "SECOND",
			ticketType:  "STUDENT",
			wantPrice:   7500,
			description: "Student tickets on G/C/D trains use the high-speed discount.",
		},
		{
			name:        "student conventional ticket gets 60 percent price",
			basePrice:   10000,
			trainType:   "K",
			seatType:    "HARD_SEAT",
			ticketType:  "STUDENT",
			wantPrice:   6000,
			description: "Student tickets on Z/T/K trains use the conventional-train discount.",
		},
		{
			name:        "child ticket gets half price",
			basePrice:   10001,
			trainType:   "D",
			seatType:    "SECOND",
			ticketType:  "CHILD",
			wantPrice:   5001,
			description: "Child tickets are rounded after applying the 50 percent discount.",
		},
		{
			name:        "unsupported ticket type returns validation error",
			basePrice:   10000,
			trainType:   "G",
			seatType:    "SECOND",
			ticketType:  "SENIOR",
			wantErr:     true,
			description: "Unknown ticket types must be rejected before an order is created.",
		},
		{
			name:        "missing seat type returns validation error",
			basePrice:   10000,
			trainType:   "G",
			seatType:    "",
			ticketType:  "ADULT",
			wantErr:     true,
			description: "Required pricing inputs must be present.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrice, err := CalculateTicketPrice(tt.basePrice, tt.trainType, tt.seatType, tt.ticketType)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("%s: expected error, got nil", tt.description)
				}
				return
			}
			if err != nil {
				t.Fatalf("%s: unexpected error: %v", tt.description, err)
			}
			if gotPrice != tt.wantPrice {
				t.Fatalf("%s: CalculateTicketPrice() = %d, want %d", tt.description, gotPrice, tt.wantPrice)
			}
		})
	}
}

func TestCalculateTicketPriceNormalizesInput(t *testing.T) {
	gotPrice, err := CalculateTicketPrice(8888, " g ", " second ", " student ")
	if err != nil {
		t.Fatalf("CalculateTicketPrice returned error for lowercase padded input: %v", err)
	}
	if gotPrice != 6666 {
		t.Fatalf("CalculateTicketPrice() = %d, want 6666", gotPrice)
	}
}
