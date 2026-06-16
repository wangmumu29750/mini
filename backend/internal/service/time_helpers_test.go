package service

import (
	"reflect"
	"testing"
	"time"
)

func Test_combineDateClock(t *testing.T) {
	baseDate := time.Date(2025, 1, 15, 0, 0, 0, 0, time.Local)

	type args struct {
		date      time.Time
		clock     string
		dayOffset int
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{name: "combines date with clock", args: args{date: baseDate, clock: "08:30:00", dayOffset: 0}, want: time.Date(2025, 1, 15, 8, 30, 0, 0, time.Local)},
		{name: "empty clock defaults to midnight", args: args{date: baseDate, clock: "", dayOffset: 0}, want: time.Date(2025, 1, 15, 0, 0, 0, 0, time.Local)},
		{name: "positive day offset", args: args{date: baseDate, clock: "12:00:00", dayOffset: 1}, want: time.Date(2025, 1, 16, 12, 0, 0, 0, time.Local)},
		{name: "negative day offset", args: args{date: baseDate, clock: "23:00:00", dayOffset: -1}, want: time.Date(2025, 1, 14, 23, 0, 0, 0, time.Local)},
		{name: "invalid clock returns original date", args: args{date: baseDate, clock: "not-a-time", dayOffset: 0}, want: baseDate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := combineDateClock(tt.args.date, tt.args.clock, tt.args.dayOffset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("combineDateClock() = %v, want %v", got, tt.want)
			}
		})
	}
}
