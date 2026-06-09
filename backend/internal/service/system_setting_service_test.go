package service

import (
	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/model"
	"mini-12306/backend/internal/repository"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestNewSystemSettingService(t *testing.T) {
	type args struct {
		settings *repository.SystemSettingRepository
	}
	tests := []struct {
		name string
		args args
		want *SystemSettingService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSystemSettingService(tt.args.settings); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSystemSettingService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemSettingService_List(t *testing.T) {
	tests := []struct {
		name    string
		s       *SystemSettingService
		want    []dto.SystemSettingResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("SystemSettingService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SystemSettingService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemSettingService_Update(t *testing.T) {
	type args struct {
		req dto.UpdateSystemSettingsRequest
	}
	tests := []struct {
		name    string
		s       *SystemSettingService
		args    args
		want    []dto.SystemSettingResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Update(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SystemSettingService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SystemSettingService.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystemSettingService_ensureDefaults(t *testing.T) {
	tests := []struct {
		name    string
		s       *SystemSettingService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.ensureDefaults(); (err != nil) != tt.wantErr {
				t.Errorf("SystemSettingService.ensureDefaults() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateSettingValue(t *testing.T) {
	type args struct {
		valueType string
		value     string
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
			if err := validateSettingValue(tt.args.valueType, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("validateSettingValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_settingResponse(t *testing.T) {
	type args struct {
		setting model.SystemSetting
	}
	tests := []struct {
		name string
		args args
		want dto.SystemSettingResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := settingResponse(tt.args.setting); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("settingResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_systemSettingInt(t *testing.T) {
	type args struct {
		tx       *gorm.DB
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
			if got := systemSettingInt(tt.args.tx, tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("systemSettingInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
