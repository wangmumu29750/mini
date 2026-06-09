package service

import (
	"mini-12306/backend/internal/repository"
	"reflect"
	"testing"
)

func TestNewAdminService(t *testing.T) {
	type args struct {
		admin *repository.AdminRepository
	}
	tests := []struct {
		name string
		args args
		want *AdminService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAdminService(tt.args.admin); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAdminService() = %v, want %v", got, tt.want)
			}
		})
	}
}
