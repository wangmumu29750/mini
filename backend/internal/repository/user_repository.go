package repository

import (
	"errors"

	"mini-12306/backend/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateWithProfile(user *model.User, profile *model.PassengerProfile) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		profile.UserID = user.ID
		return tx.Create(profile).Error
	})
}

func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Profile").Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) FindByID(id uint64) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Profile").First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) UsernameExists(username string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) IDCardExists(idCardNo string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.PassengerProfile{}).Where("id_card_no = ?", idCardNo).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) PhoneExists(phone string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.PassengerProfile{}).Where("phone = ?", phone).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) ListPassengerProfilesByUser(userID uint64) ([]model.PassengerProfile, error) {
	var profiles []model.PassengerProfile
	err := r.db.Where("user_id = ?", userID).Order("id DESC").Find(&profiles).Error
	return profiles, err
}

func (r *UserRepository) FindPassengerProfileByID(userID, passengerID uint64) (*model.PassengerProfile, error) {
	var profile model.PassengerProfile
	err := r.db.Where("user_id = ? AND id = ?", userID, passengerID).First(&profile).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &profile, err
}

func (r *UserRepository) CreatePassengerProfile(profile *model.PassengerProfile) error {
	return r.db.Create(profile).Error
}
