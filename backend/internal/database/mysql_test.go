package database

import (
	"strings"
	"testing"

	"mini-12306/backend/internal/config"
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
