package service

import (
	"reflect"
	"testing"
	"time"
)

func Test_combineDateClock(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := combineDateClock(tt.args.date, tt.args.clock, tt.args.dayOffset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("combineDateClock() = %v, want %v", got, tt.want)
			}
		})
	}
}
