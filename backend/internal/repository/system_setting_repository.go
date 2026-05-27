package repository

import (
	"mini-12306/backend/internal/model"

	"gorm.io/gorm"
)

type SystemSettingRepository struct {
	db *gorm.DB
}

func NewSystemSettingRepository(db *gorm.DB) *SystemSettingRepository {
	return &SystemSettingRepository{db: db}
}

func (r *SystemSettingRepository) List() ([]model.SystemSetting, error) {
	var settings []model.SystemSetting
	err := r.db.Order("id ASC").Find(&settings).Error
	return settings, err
}

func (r *SystemSettingRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
