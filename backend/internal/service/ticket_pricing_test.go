package service

import (
	"reflect"
	"testing"
)

func Test_newTicketPriceCalculator(t *testing.T) {
	tests := []struct {
		name string
		want *ticketPriceCalculator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newTicketPriceCalculator(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTicketPriceCalculator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateTicketPrice(t *testing.T) {
	type args struct {
		basePriceCents int64
		seatType       string
		ticketType     string
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
			got, err := CalculateTicketPrice(tt.args.basePriceCents, tt.args.seatType, tt.args.ticketType)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateTicketPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalculateTicketPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ticketPriceCalculator_Calculate(t *testing.T) {
	type args struct {
		basePriceCents int64
		seatType       string
		ticketType     string
	}
	tests := []struct {
		name    string
		c       *ticketPriceCalculator
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Calculate(tt.args.basePriceCents, tt.args.seatType, tt.args.ticketType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ticketPriceCalculator.Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ticketPriceCalculator.Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_adultTicketRule_Apply(t *testing.T) {
	type args struct {
		ctx TicketPriceContext
	}
	tests := []struct {
		name    string
		a       adultTicketRule
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Apply(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("adultTicketRule.Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("adultTicketRule.Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentTicketRule_Apply(t *testing.T) {
	type args struct {
		ctx TicketPriceContext
	}
	tests := []struct {
		name    string
		s       studentTicketRule
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Apply(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("studentTicketRule.Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("studentTicketRule.Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_childTicketRule_Apply(t *testing.T) {
	type args struct {
		ctx TicketPriceContext
	}
	tests := []struct {
		name    string
		c       childTicketRule
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Apply(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("childTicketRule.Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("childTicketRule.Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_roundedPrice(t *testing.T) {
	type args struct {
		base   int64
		factor float64
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
			if got := roundedPrice(tt.args.base, tt.args.factor); got != tt.want {
				t.Errorf("roundedPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
