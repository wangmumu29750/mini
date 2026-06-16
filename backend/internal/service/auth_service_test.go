package service

import (
	"mini-12306/backend/internal/config"
	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/model"
	"mini-12306/backend/internal/repository"
	"reflect"
	"testing"
)

func TestNewAuthService(t *testing.T) {
	type args struct {
		cfg   config.Config
		users *repository.UserRepository
	}
	tests := []struct {
		name string
		args args
		want *AuthService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthService(tt.args.cfg, tt.args.users); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	type args struct {
		req dto.RegisterRequest
	}
	tests := []struct {
		name    string
		s       *AuthService
		args    args
		want    dto.AuthResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Register(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthService.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	type args struct {
		req dto.LoginRequest
	}
	tests := []struct {
		name    string
		s       *AuthService
		args    args
		want    dto.AuthResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Login(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthService.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_CurrentUser(t *testing.T) {
	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		s       *AuthService
		args    args
		want    dto.CurrentUserResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.CurrentUser(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.CurrentUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthService.CurrentUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_ListPassengerProfiles(t *testing.T) {
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		s       *AuthService
		args    args
		want    []dto.PassengerSummaryResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ListPassengerProfiles(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.ListPassengerProfiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthService.ListPassengerProfiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_CreatePassengerProfile(t *testing.T) {
	type args struct {
		userID uint64
		req    dto.PassengerProfileRequest
	}
	tests := []struct {
		name    string
		s       *AuthService
		args    args
		want    dto.PassengerSummaryResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.CreatePassengerProfile(tt.args.userID, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.CreatePassengerProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthService.CreatePassengerProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_authResponse(t *testing.T) {
	type args struct {
		user model.User
	}
	tests := []struct {
		name    string
		s       *AuthService
		args    args
		want    dto.AuthResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.authResponse(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.authResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthService.authResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_currentUserResponse(t *testing.T) {
	type args struct {
		user model.User
	}
	tests := []struct {
		name string
		args args
		want dto.CurrentUserResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := currentUserResponse(tt.args.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("currentUserResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
