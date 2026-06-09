package repository

import (
	"mini-12306/backend/internal/model"
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestNewUserRepository(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *UserRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_CreateWithProfile(t *testing.T) {
	type args struct {
		user    *model.User
		profile *model.PassengerProfile
	}
	tests := []struct {
		name    string
		r       *UserRepository
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.CreateWithProfile(tt.args.user, tt.args.profile); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.CreateWithProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_FindByUsername(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		r       *UserRepository
		args    args
		want    *model.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.FindByUsername(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.FindByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.FindByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		r       *UserRepository
		args    args
		want    *model.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.FindByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_UsernameExists(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		r       *UserRepository
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.UsernameExists(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.UsernameExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserRepository.UsernameExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_IDCardExists(t *testing.T) {
	type args struct {
		idCardNo string
	}
	tests := []struct {
		name    string
		r       *UserRepository
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.IDCardExists(tt.args.idCardNo)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.IDCardExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserRepository.IDCardExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_PhoneExists(t *testing.T) {
	type args struct {
		phone string
	}
	tests := []struct {
		name    string
		r       *UserRepository
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.PhoneExists(tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.PhoneExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserRepository.PhoneExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_ListPassengerProfilesByUser(t *testing.T) {
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		r       *UserRepository
		args    args
		want    []model.PassengerProfile
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.ListPassengerProfilesByUser(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.ListPassengerProfilesByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.ListPassengerProfilesByUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_FindPassengerProfileByID(t *testing.T) {
	type args struct {
		userID      uint64
		passengerID uint64
	}
	tests := []struct {
		name    string
		r       *UserRepository
		args    args
		want    *model.PassengerProfile
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.FindPassengerProfileByID(tt.args.userID, tt.args.passengerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.FindPassengerProfileByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepository.FindPassengerProfileByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserRepository_CreatePassengerProfile(t *testing.T) {
	type args struct {
		profile *model.PassengerProfile
	}
	tests := []struct {
		name    string
		r       *UserRepository
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.CreatePassengerProfile(tt.args.profile); (err != nil) != tt.wantErr {
				t.Errorf("UserRepository.CreatePassengerProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
