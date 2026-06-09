package database

import (
	"database/sql"
	"mini-12306/backend/internal/config"
	"reflect"
	"strings"
	"testing"

	"gorm.io/gorm"
)

func TestConnectRequiresMySQLDSN(t *testing.T) {
	t.Parallel()

	db, err := Connect(config.Config{})
	if err == nil {
		t.Fatal("expected missing MYSQL_DSN to fail")
	}
	if db != nil {
		t.Fatal("expected missing MYSQL_DSN to return nil db")
	}
	if !strings.Contains(err.Error(), "MYSQL_DSN") {
		t.Fatalf("expected error to mention MYSQL_DSN, got %q", err.Error())
	}
}

func TestConnect(t *testing.T) {
	type args struct {
		cfg config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Connect(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Connect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPing(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Ping(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLDB(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		args    args
		want    *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SQLDB(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLDB() = %v, want %v", got, tt.want)
			}
		})
	}
}
