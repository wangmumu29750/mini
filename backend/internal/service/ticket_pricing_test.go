package service

import (
	"testing"
)

func Test_newTicketPriceCalculator(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "creates calculator with adult, student, child rules"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newTicketPriceCalculator()
			if got == nil {
				t.Fatal("newTicketPriceCalculator() returned nil")
			}
			if got.rules == nil {
				t.Fatal("newTicketPriceCalculator() rules map is nil")
			}
			if len(got.rules) != 3 {
				t.Errorf("newTicketPriceCalculator() has %d rules, want 3", len(got.rules))
			}
			for _, key := range []string{"ADULT", "STUDENT", "CHILD"} {
				if _, ok := got.rules[key]; !ok {
					t.Errorf("newTicketPriceCalculator() missing rule for %s", key)
				}
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
		{name: "adult ticket full price", args: args{basePriceCents: 10000, seatType: "SECOND", ticketType: "ADULT"}, want: 10000, wantErr: false},
		{name: "student ticket 75%", args: args{basePriceCents: 10000, seatType: "SECOND", ticketType: "STUDENT"}, want: 7500, wantErr: false},
		{name: "child ticket 50%", args: args{basePriceCents: 10000, seatType: "SECOND", ticketType: "CHILD"}, want: 5000, wantErr: false},
		{name: "student not second class", args: args{basePriceCents: 10000, seatType: "FIRST", ticketType: "STUDENT"}, want: 0, wantErr: true},
		{name: "unknown ticket type", args: args{basePriceCents: 10000, seatType: "SECOND", ticketType: "UNKNOWN"}, want: 0, wantErr: true},
		{name: "empty ticket type", args: args{basePriceCents: 10000, seatType: "SECOND", ticketType: ""}, want: 0, wantErr: true},
		{name: "empty seat type", args: args{basePriceCents: 10000, seatType: "", ticketType: "ADULT"}, want: 0, wantErr: true},
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
		{name: "adult full price", c: newTicketPriceCalculator(), args: args{basePriceCents: 20000, seatType: "SECOND", ticketType: "ADULT"}, want: 20000, wantErr: false},
		{name: "student lowercased", c: newTicketPriceCalculator(), args: args{basePriceCents: 10000, seatType: "second", ticketType: "student"}, want: 7500, wantErr: false},
		{name: "empty seat type errors", c: newTicketPriceCalculator(), args: args{basePriceCents: 10000, seatType: "", ticketType: "ADULT"}, want: 0, wantErr: true},
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
		{name: "returns base price unchanged", a: adultTicketRule{}, args: args{ctx: TicketPriceContext{BasePriceCents: 10000, SeatType: "SECOND", TicketType: "ADULT"}}, want: 10000, wantErr: false},
		{name: "zero base price", a: adultTicketRule{}, args: args{ctx: TicketPriceContext{BasePriceCents: 0, SeatType: "FIRST", TicketType: "ADULT"}}, want: 0, wantErr: false},
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
		{name: "second class 75% price", s: studentTicketRule{}, args: args{ctx: TicketPriceContext{BasePriceCents: 10000, SeatType: "SECOND", TicketType: "STUDENT"}}, want: 7500, wantErr: false},
		{name: "first class rejected", s: studentTicketRule{}, args: args{ctx: TicketPriceContext{BasePriceCents: 10000, SeatType: "FIRST", TicketType: "STUDENT"}}, want: 0, wantErr: true},
		{name: "business class rejected", s: studentTicketRule{}, args: args{ctx: TicketPriceContext{BasePriceCents: 10000, SeatType: "BUSINESS", TicketType: "STUDENT"}}, want: 0, wantErr: true},
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
		{name: "child 50% price", c: childTicketRule{}, args: args{ctx: TicketPriceContext{BasePriceCents: 10000, SeatType: "SECOND", TicketType: "CHILD"}}, want: 5000, wantErr: false},
		{name: "child rounds up for odd", c: childTicketRule{}, args: args{ctx: TicketPriceContext{BasePriceCents: 10003, SeatType: "SECOND", TicketType: "CHILD"}}, want: 5002, wantErr: false},
		{name: "zero base price", c: childTicketRule{}, args: args{ctx: TicketPriceContext{BasePriceCents: 0, SeatType: "FIRST", TicketType: "CHILD"}}, want: 0, wantErr: false},
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
		{name: "75 percent of 10000", args: args{base: 10000, factor: 0.75}, want: 7500},
		{name: "50 percent of 10000", args: args{base: 10000, factor: 0.5}, want: 5000},
		{name: "rounds down at .5 ties to even", args: args{base: 10001, factor: 0.5}, want: 5000},
		{name: "rounds up at .75", args: args{base: 10001, factor: 0.75}, want: 7501},
		{name: "zero base returns zero", args: args{base: 0, factor: 0.75}, want: 0},
		{name: "negative base returns zero", args: args{base: -100, factor: 0.5}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := roundedPrice(tt.args.base, tt.args.factor); got != tt.want {
				t.Errorf("roundedPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
