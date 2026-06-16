package repository

import (
	"errors"
	"strings"

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
	err := r.db.
		Model(&model.PassengerProfile{}).
		Distinct("passenger_profiles.*").
		Joins("LEFT JOIN passenger_associations pa ON pa.passenger_profile_id = passenger_profiles.id").
		Select("passenger_profiles.*, COALESCE(pa.passenger_type, passenger_profiles.passenger_type) AS passenger_type").
		Where("passenger_profiles.user_id = ? OR pa.owner_user_id = ?", userID, userID).
		Order("passenger_profiles.id DESC").
		Find(&profiles).Error
	return profiles, err
}

func (r *UserRepository) FindPassengerProfileByID(userID, passengerID uint64) (*model.PassengerProfile, error) {
	var profile model.PassengerProfile
	err := r.db.
		Model(&model.PassengerProfile{}).
		Distinct("passenger_profiles.*").
		Joins("LEFT JOIN passenger_associations pa ON pa.passenger_profile_id = passenger_profiles.id").
		Select("passenger_profiles.*, COALESCE(pa.passenger_type, passenger_profiles.passenger_type) AS passenger_type").
		Where("passenger_profiles.id = ? AND (passenger_profiles.user_id = ? OR pa.owner_user_id = ?)", passengerID, userID, userID).
		First(&profile).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &profile, err
}

func (r *UserRepository) FindVerifiedPassengerProfileByIdentity(idCardNo string) (*model.PassengerProfile, error) {
	var profile model.PassengerProfile
	err := r.db.
		Where("id_card_no = ? AND verified_status = ?", strings.TrimSpace(idCardNo), model.VerificationStatusVerified).
		First(&profile).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &profile, err
}

func (r *UserRepository) CreatePassengerProfile(profile *model.PassengerProfile) error {
	return r.db.Create(profile).Error
}

func (r *UserRepository) FindPassengerProfileByIdentity(idCardNo string) (*model.PassengerProfile, error) {
	var profile model.PassengerProfile
	err := r.db.Where("id_card_no = ?", strings.TrimSpace(idCardNo)).First(&profile).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &profile, err
}

func (r *UserRepository) LinkPassengerProfile(ownerUserID, passengerProfileID uint64) error {
	association := &model.PassengerAssociation{
		OwnerUserID:        ownerUserID,
		PassengerProfileID: passengerProfileID,
	}
	return r.db.Where("owner_user_id = ? AND passenger_profile_id = ?", ownerUserID, passengerProfileID).FirstOrCreate(association).Error
}

func (r *UserRepository) UpsertPassengerAssociation(ownerUserID, passengerProfileID uint64, passengerType model.PassengerType) error {
	association := &model.PassengerAssociation{
		OwnerUserID:        ownerUserID,
		PassengerProfileID: passengerProfileID,
		PassengerType:      passengerType,
	}
	return r.db.Where("owner_user_id = ? AND passenger_profile_id = ?", ownerUserID, passengerProfileID).
		Assign(association).
		FirstOrCreate(association).Error
}
