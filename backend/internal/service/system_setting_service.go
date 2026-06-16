package service

import (
	"net/http"
	"strconv"
	"strings"

	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/model"
	"mini-12306/backend/internal/repository"
	apperrors "mini-12306/backend/pkg/errors"
	"mini-12306/backend/pkg/response"

	"gorm.io/gorm"
)

type SystemSettingService struct {
	settings *repository.SystemSettingRepository
}

var defaultSystemSettings = []model.SystemSetting{
	{Key: "order_pay_expire_minutes", Value: "15", ValueType: "INT", Description: "Order payment timeout in minutes"},
	{Key: "refund_cutoff_minutes", Value: "30", ValueType: "INT", Description: "Minutes before departure after which refund is blocked"},
	{Key: "change_cutoff_minutes", Value: "30", ValueType: "INT", Description: "Minutes before departure after which change is blocked"},
	{Key: "refund_fee_percent", Value: "0", ValueType: "INT", Description: "Refund fee percentage for mock settlement"},
	{Key: "change_fee_percent", Value: "0", ValueType: "INT", Description: "Change fee percentage based on the original paid ticket amount"},
	{Key: "mock_payment_enabled", Value: "true", ValueType: "BOOL", Description: "Whether local mock payment is enabled"},
}

func NewSystemSettingService(settings *repository.SystemSettingRepository) *SystemSettingService {
	return &SystemSettingService{settings: settings}
}

func (s *SystemSettingService) List() ([]dto.SystemSettingResponse, error) {
	if err := s.ensureDefaults(); err != nil {
		return nil, err
	}
	settings, err := s.settings.List()
	if err != nil {
		return nil, err
	}
	result := make([]dto.SystemSettingResponse, 0, len(settings))
	for _, setting := range settings {
		result = append(result, settingResponse(setting))
	}
	return result, nil
}

func (s *SystemSettingService) Update(req dto.UpdateSystemSettingsRequest) ([]dto.SystemSettingResponse, error) {
	if err := s.ensureDefaults(); err != nil {
		return nil, err
	}

	allowed := map[string]model.SystemSetting{}
	for _, setting := range defaultSystemSettings {
		allowed[setting.Key] = setting
	}

	err := s.settings.Transaction(func(tx *gorm.DB) error {
		for _, item := range req.Settings {
			key := strings.TrimSpace(item.Key)
			value := strings.TrimSpace(item.Value)
			setting, ok := allowed[key]
			if !ok {
				return apperrors.New(http.StatusBadRequest, response.CodeValidationError, "unsupported system setting")
			}
			if err := validateSettingValue(setting.ValueType, value); err != nil {
				return err
			}
			if err := tx.Model(&model.SystemSetting{}).Where("`key` = ?", key).Update("value", value).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.List()
}

func (s *SystemSettingService) ensureDefaults() error {
	return s.settings.Transaction(func(tx *gorm.DB) error {
		for _, setting := range defaultSystemSettings {
			var count int64
			if err := tx.Model(&model.SystemSetting{}).Where("`key` = ?", setting.Key).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				if err := tx.Create(&setting).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func validateSettingValue(valueType, value string) error {
	switch valueType {
	case "INT":
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 0 {
			return apperrors.New(http.StatusBadRequest, response.CodeValidationError, "system setting must be a non-negative integer")
		}
	case "BOOL":
		if value != "true" && value != "false" {
			return apperrors.New(http.StatusBadRequest, response.CodeValidationError, "system setting must be true or false")
		}
	}
	return nil
}

func settingResponse(setting model.SystemSetting) dto.SystemSettingResponse {
	return dto.SystemSettingResponse{
		Key:         setting.Key,
		Value:       setting.Value,
		ValueType:   setting.ValueType,
		Description: setting.Description,
	}
}

func systemSettingInt(tx *gorm.DB, key string, fallback int) int {
	var setting model.SystemSetting
	if err := tx.Where("`key` = ?", key).First(&setting).Error; err != nil {
		return fallback
	}
	value, err := strconv.Atoi(strings.TrimSpace(setting.Value))
	if err != nil || value < 0 {
		return fallback
	}
	return value
}
