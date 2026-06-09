package config

import (
	"reflect"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Load(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_IsProduction(t *testing.T) {
	tests := []struct {
		name string
		c    Config
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsProduction(); got != tt.want {
				t.Errorf("Config.IsProduction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_TokenExpireDuration(t *testing.T) {
	tests := []struct {
		name string
		c    Config
		want time.Duration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.TokenExpireDuration(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.TokenExpireDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_OrderPayExpireDuration(t *testing.T) {
	tests := []struct {
		name string
		c    Config
		want time.Duration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.OrderPayExpireDuration(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.OrderPayExpireDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_envString(t *testing.T) {
	type args struct {
		key      string
		fallback string
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
			if got := envString(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("envString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_envInt(t *testing.T) {
	type args struct {
		key      string
		fallback int
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
			if got := envInt(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("envInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
