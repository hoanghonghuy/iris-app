package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type SchoolAdminService struct {
	schoolAdminRepo *repo.SchoolAdminRepo
	userRepo        *repo.UserRepo
}

func NewSchoolAdminService(schoolAdminRepo *repo.SchoolAdminRepo, userRepo *repo.UserRepo) *SchoolAdminService {
	return &SchoolAdminService{
		schoolAdminRepo: schoolAdminRepo,
		userRepo:        userRepo,
	}
}

// Create tạo mới school admin (SUPER_ADMIN only).
// Gán user vào trường với role SCHOOL_ADMIN, đồng thời gán role trong bảng user_roles.
func (s *SchoolAdminService) Create(ctx context.Context, userID, schoolID uuid.UUID, fullName, phone string) (*model.SchoolAdmin, error) {
	// Gán role SCHOOL_ADMIN cho user (nếu chưa có)
	if err := s.userRepo.AssignRole(ctx, userID, "SCHOOL_ADMIN"); err != nil {
		return nil, ErrFailedToAssignRole
	}

	// Tạo record trong bảng school_admins
	adminID, err := s.schoolAdminRepo.Create(ctx, userID, schoolID, fullName, phone)
	if err != nil {
		return nil, err
	}

	return s.schoolAdminRepo.GetByAdminID(ctx, adminID)
}

// GetByAdminID lấy thông tin school admin theo admin_id
func (s *SchoolAdminService) GetByAdminID(ctx context.Context, adminID uuid.UUID) (*model.SchoolAdmin, error) {
	admin, err := s.schoolAdminRepo.GetByAdminID(ctx, adminID)
	if err != nil {
		return nil, ErrSchoolAdminNotFound
	}
	return admin, nil
}

// List lấy danh sách tất cả school admins (SUPER_ADMIN only)
func (s *SchoolAdminService) List(ctx context.Context, limit, offset int) ([]model.SchoolAdmin, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.schoolAdminRepo.List(ctx, limit, offset)
}

// Delete xóa school admin theo admin_id (SUPER_ADMIN only).
// Chỉ xóa record trong school_admins, không xóa user.
func (s *SchoolAdminService) Delete(ctx context.Context, adminID uuid.UUID) error {
	return s.schoolAdminRepo.Delete(ctx, adminID)
}
