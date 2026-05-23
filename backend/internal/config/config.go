package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultAppEnv                = "dev"
	defaultAppPort               = "8080"
	defaultJWTSecret             = "change-me-in-local-env"
	defaultTokenExpireMinutes    = 120
	defaultOrderPayExpireMinutes = 15
)

type Config struct {
	AppEnv                string
	AppPort               string
	MySQLDSN              string
	JWTSecret             string
	TokenExpireMinutes    int
	OrderPayExpireMinutes int
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppEnv:                envString("APP_ENV", defaultAppEnv),
		AppPort:               envString("APP_PORT", defaultAppPort),
		MySQLDSN:              strings.TrimSpace(os.Getenv("MYSQL_DSN")),
		JWTSecret:             envString("JWT_SECRET", defaultJWTSecret),
		TokenExpireMinutes:    envInt("TOKEN_EXPIRE_MINUTES", defaultTokenExpireMinutes),
		OrderPayExpireMinutes: envInt("ORDER_PAY_EXPIRE_MINUTES", defaultOrderPayExpireMinutes),
	}
}

func (c Config) IsProduction() bool {
	return strings.EqualFold(c.AppEnv, "prod") || strings.EqualFold(c.AppEnv, "production")
}

func (c Config) TokenExpireDuration() time.Duration {
	return time.Duration(c.TokenExpireMinutes) * time.Minute
}

func (c Config) OrderPayExpireDuration() time.Duration {
	return time.Duration(c.OrderPayExpireMinutes) * time.Minute
}

func envString(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func envInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
