package database

import (
	"database/sql"
	"fmt"
	"time"

	"mini-12306/backend/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg config.Config) (*gorm.DB, error) {
	if cfg.MySQLDSN == "" {
		return nil, fmt.Errorf("MYSQL_DSN is required")
	}

	gormLogger := logger.Default.LogMode(logger.Warn)
	if !cfg.IsProduction() {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func Ping(db *gorm.DB) error {
	sqlDB, err := SQLDB(db)
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func SQLDB(db *gorm.DB) (*sql.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("database is not configured")
	}
	return db.DB()
}
