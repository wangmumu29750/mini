package config

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name   string
		envs   map[string]string
		assert func(t *testing.T, got Config)
	}{
		{
			name: "defaults when no env vars set",
			assert: func(t *testing.T, got Config) {
				if got.AppEnv != "dev" {
					t.Errorf("AppEnv = %v, want dev", got.AppEnv)
				}
				if got.AppPort != "8080" {
					t.Errorf("AppPort = %v, want 8080", got.AppPort)
				}
				if got.JWTSecret != "change-me-in-local-env" {
					t.Errorf("JWTSecret = %v, want default", got.JWTSecret)
				}
			},
		},
		{
			name: "production env",
			envs: map[string]string{"APP_ENV": "production", "APP_PORT": "3000"},
			assert: func(t *testing.T, got Config) {
				if got.AppEnv != "production" {
					t.Errorf("AppEnv = %v, want production", got.AppEnv)
				}
				if got.AppPort != "3000" {
					t.Errorf("AppPort = %v, want 3000", got.AppPort)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear relevant env vars before test
			for _, k := range []string{"APP_ENV", "APP_PORT", "MYSQL_DSN", "JWT_SECRET", "TOKEN_EXPIRE_MINUTES", "ORDER_PAY_EXPIRE_MINUTES"} {
				os.Unsetenv(k)
			}
			for k, v := range tt.envs {
				os.Setenv(k, v)
			}
			got := Load()
			tt.assert(t, got)
		})
	}
}

func TestConfig_IsProduction(t *testing.T) {
	tests := []struct {
		name string
		c    Config
		want bool
	}{
		{name: "dev env", c: Config{AppEnv: "dev"}, want: false},
		{name: "prod env", c: Config{AppEnv: "prod"}, want: true},
		{name: "production env", c: Config{AppEnv: "production"}, want: true},
		{name: "PROD uppercase", c: Config{AppEnv: "PROD"}, want: true},
		{name: "Production mixed case", c: Config{AppEnv: "Production"}, want: true},
		{name: "staging env", c: Config{AppEnv: "staging"}, want: false},
		{name: "empty env", c: Config{AppEnv: ""}, want: false},
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
		{name: "default 120 min", c: Config{TokenExpireMinutes: 120}, want: 120 * time.Minute},
		{name: "60 min", c: Config{TokenExpireMinutes: 60}, want: 60 * time.Minute},
		{name: "zero minutes", c: Config{TokenExpireMinutes: 0}, want: 0},
		{name: "1 minute", c: Config{TokenExpireMinutes: 1}, want: 1 * time.Minute},
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
		{name: "default 15 min", c: Config{OrderPayExpireMinutes: 15}, want: 15 * time.Minute},
		{name: "30 min", c: Config{OrderPayExpireMinutes: 30}, want: 30 * time.Minute},
		{name: "zero minutes", c: Config{OrderPayExpireMinutes: 0}, want: 0},
		{name: "5 minutes", c: Config{OrderPayExpireMinutes: 5}, want: 5 * time.Minute},
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
		envs map[string]string
		want string
	}{
		{name: "env var set", args: args{key: "TEST_KEY", fallback: "default"}, envs: map[string]string{"TEST_KEY": "hello"}, want: "hello"},
		{name: "env var not set returns fallback", args: args{key: "TEST_MISSING", fallback: "default"}, envs: nil, want: "default"},
		{name: "env var empty returns fallback", args: args{key: "TEST_EMPTY", fallback: "default"}, envs: map[string]string{"TEST_EMPTY": ""}, want: "default"},
		{name: "env var whitespace only returns fallback", args: args{key: "TEST_SPACE", fallback: "fallback"}, envs: map[string]string{"TEST_SPACE": "   "}, want: "fallback"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(tt.args.key)
			for k, v := range tt.envs {
				os.Setenv(k, v)
			}
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
		envs map[string]string
		want int
	}{
		{name: "env var set to valid int", args: args{key: "TEST_INT", fallback: 10}, envs: map[string]string{"TEST_INT": "42"}, want: 42},
		{name: "env var not set returns fallback", args: args{key: "TEST_INT_MISSING", fallback: 10}, envs: nil, want: 10},
		{name: "env var empty returns fallback", args: args{key: "TEST_INT_EMPTY", fallback: 10}, envs: map[string]string{"TEST_INT_EMPTY": ""}, want: 10},
		{name: "env var non-numeric returns fallback", args: args{key: "TEST_INT_NAN", fallback: 10}, envs: map[string]string{"TEST_INT_NAN": "abc"}, want: 10},
		{name: "env var zero returns fallback", args: args{key: "TEST_INT_ZERO", fallback: 10}, envs: map[string]string{"TEST_INT_ZERO": "0"}, want: 10},
		{name: "env var negative returns fallback", args: args{key: "TEST_INT_NEG", fallback: 10}, envs: map[string]string{"TEST_INT_NEG": "-5"}, want: 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(tt.args.key)
			for k, v := range tt.envs {
				os.Setenv(k, v)
			}
			if got := envInt(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("envInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
