package repository

import (
	"mini-12306/backend/internal/model"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestNewSystemSettingRepository(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *SystemSettingRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSystemSettingRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSystemSettingRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemSettingRepository_List(t *testing.T) {
	tests := []struct {
		name    string
		r       *SystemSettingRepository
		want    []model.SystemSetting
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("SystemSettingRepository.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SystemSettingRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemSettingRepository_Transaction(t *testing.T) {
	type args struct {
		fn func(tx *gorm.DB) error
	}
	tests := []struct {
		name    string
		r       *SystemSettingRepository
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Transaction(tt.args.fn); (err != nil) != tt.wantErr {
				t.Errorf("SystemSettingRepository.Transaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
