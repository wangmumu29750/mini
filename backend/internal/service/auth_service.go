package service

import (
	"net/http"
	"strings"
	"time"

	"mini-12306/backend/internal/config"
	"mini-12306/backend/internal/dto"
	"mini-12306/backend/internal/model"
	"mini-12306/backend/internal/repository"
	"mini-12306/backend/pkg/auth"
	apperrors "mini-12306/backend/pkg/errors"
	"mini-12306/backend/pkg/mock"
	"mini-12306/backend/pkg/response"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	cfg   config.Config
	users *repository.UserRepository
}

func NewAuthService(cfg config.Config, users *repository.UserRepository) *AuthService {
	return &AuthService{cfg: cfg, users: users}
}

func (s *AuthService) Register(req dto.RegisterRequest) (dto.AuthResponse, error) {
	req.Username = strings.TrimSpace(req.Username)
	req.RealName = strings.TrimSpace(req.RealName)
	req.IDCardNo = strings.TrimSpace(req.IDCardNo)
	req.Phone = strings.TrimSpace(req.Phone)
	req.BankCardNo = strings.TrimSpace(req.BankCardNo)

	if err := mock.VerifyIdentity(req.RealName, req.IDCardNo, req.Phone, req.BankCardNo); err != nil {
		return dto.AuthResponse{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, err.Error())
	}

	exists, err := s.users.UsernameExists(req.Username)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if exists {
		return dto.AuthResponse{}, apperrors.New(http.StatusConflict, response.CodeConflict, "用户名已存在")
	}

	exists, err = s.users.IDCardExists(req.IDCardNo)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if exists {
		return dto.AuthResponse{}, apperrors.New(http.StatusConflict, response.CodeConflict, "身份证号已被注册")
	}

	exists, err = s.users.PhoneExists(req.Phone)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if exists {
		return dto.AuthResponse{}, apperrors.New(http.StatusConflict, response.CodeConflict, "手机号已被注册")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	user := &model.User{
		Username:     req.Username,
		PasswordHash: string(passwordHash),
		Role:         model.UserRolePassenger,
		Status:       model.UserStatusActive,
	}
	profile := &model.PassengerProfile{
		RealName:       req.RealName,
		IDCardNo:       req.IDCardNo,
		Phone:          req.Phone,
		BankCardNo:     req.BankCardNo,
		VerifiedStatus: model.VerificationStatusVerified,
	}

	if err := s.users.CreateWithProfile(user, profile); err != nil {
		return dto.AuthResponse{}, err
	}

	return s.authResponse(*user)
}

func (s *AuthService) Login(req dto.LoginRequest) (dto.AuthResponse, error) {
	user, err := s.users.FindByUsername(strings.TrimSpace(req.Username))
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return dto.AuthResponse{}, apperrors.New(http.StatusUnauthorized, response.CodeUnauthorized, "用户名或密码错误")
	}
	if user.Status != model.UserStatusActive {
		return dto.AuthResponse{}, apperrors.New(http.StatusForbidden, response.CodeForbidden, "账号已被禁用")
	}

	now := time.Now()
	user.LastLoginAt = &now

	return s.authResponse(*user)
}

func (s *AuthService) CurrentUser(id uint64) (dto.CurrentUserResponse, error) {
	user, err := s.users.FindByID(id)
	if err != nil {
		return dto.CurrentUserResponse{}, err
	}
	if user == nil {
		return dto.CurrentUserResponse{}, apperrors.New(http.StatusUnauthorized, response.CodeUnauthorized, "登录状态已失效")
	}
	return currentUserResponse(*user), nil
}

func (s *AuthService) ListPassengerProfiles(userID uint64) ([]dto.PassengerSummaryResponse, error) {
	profiles, err := s.users.ListPassengerProfilesByUser(userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.PassengerSummaryResponse, 0, len(profiles))
	for _, profile := range profiles {
		result = append(result, dto.PassengerSummaryResponse{
			ID:             profile.ID,
			RealName:       profile.RealName,
			IDCardNoMasked: maskIDCardNo(profile.IDCardNo),
			PassengerType:  string(profile.PassengerType),
		})
	}
	return result, nil
}

func (s *AuthService) CreatePassengerProfile(userID uint64, req dto.PassengerProfileRequest) (dto.PassengerSummaryResponse, error) {
	req.RealName = strings.TrimSpace(req.RealName)
	req.IDCardNo = strings.TrimSpace(req.IDCardNo)
	req.Phone = strings.TrimSpace(req.Phone)
	req.BankCardNo = strings.TrimSpace(req.BankCardNo)
	req.PassengerType = strings.ToUpper(strings.TrimSpace(req.PassengerType))

	if err := mock.VerifyIdentity(req.RealName, req.IDCardNo, req.Phone, req.BankCardNo); err != nil {
		return dto.PassengerSummaryResponse{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, err.Error())
	}
	if !isSupportedPassengerType(req.PassengerType) {
		return dto.PassengerSummaryResponse{}, apperrors.New(http.StatusBadRequest, response.CodeValidationError, "不支持的乘车人类型")
	}

	idCardExists, err := s.users.IDCardExists(req.IDCardNo)
	if err != nil {
		return dto.PassengerSummaryResponse{}, err
	}
	if idCardExists {
		return dto.PassengerSummaryResponse{}, apperrors.New(http.StatusConflict, response.CodeConflict, "身份证号已绑定其他乘车人")
	}

	phoneExists, err := s.users.PhoneExists(req.Phone)
	if err != nil {
		return dto.PassengerSummaryResponse{}, err
	}
	if phoneExists {
		return dto.PassengerSummaryResponse{}, apperrors.New(http.StatusConflict, response.CodeConflict, "手机号已绑定其他乘车人")
	}

	profile := &model.PassengerProfile{
		UserID:         userID,
		RealName:       req.RealName,
		IDCardNo:       req.IDCardNo,
		Phone:          req.Phone,
		BankCardNo:     req.BankCardNo,
		PassengerType:  model.PassengerType(req.PassengerType),
		VerifiedStatus: model.VerificationStatusVerified,
	}
	if err := s.users.CreatePassengerProfile(profile); err != nil {
		return dto.PassengerSummaryResponse{}, err
	}
	return dto.PassengerSummaryResponse{
		ID:             profile.ID,
		RealName:       profile.RealName,
		IDCardNoMasked: maskIDCardNo(profile.IDCardNo),
		PassengerType:  string(profile.PassengerType),
	}, nil
}

func isSupportedPassengerType(value string) bool {
	switch model.PassengerType(value) {
	case model.PassengerTypeAdult, model.PassengerTypeStudent, model.PassengerTypeChild:
		return true
	default:
		return false
	}
}

func (s *AuthService) authResponse(user model.User) (dto.AuthResponse, error) {
	token, err := auth.GenerateToken(s.cfg.JWTSecret, s.cfg.TokenExpireDuration(), auth.Principal{
		UserID:   user.ID,
		Username: user.Username,
		Role:     string(user.Role),
	})
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return dto.AuthResponse{
		AccessToken: token,
		User:        currentUserResponse(user),
	}, nil
}

func currentUserResponse(user model.User) dto.CurrentUserResponse {
	return dto.CurrentUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     string(user.Role),
	}
}
