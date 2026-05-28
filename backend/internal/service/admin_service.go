package service

import (
	"mini-12306/backend/internal/repository"
)

type AdminService struct {
	admin *repository.AdminRepository
}

func NewAdminService(admin *repository.AdminRepository) *AdminService {
	return &AdminService{admin: admin}
}
