package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type ParentService struct {
	parentRepo        *repo.ParentRepo
	studentParentRepo *repo.StudentParentRepo
	studentRepo       *repo.StudentRepo
}

func NewParentService(parentRepo *repo.ParentRepo,
	studentParentRepo *repo.StudentParentRepo,
	studentRepo *repo.StudentRepo) *ParentService {
	return &ParentService{
		parentRepo:        parentRepo,
		studentParentRepo: studentParentRepo,
		studentRepo:       studentRepo,
	}
}

// List lấy danh sách phụ huynh.
func (s *ParentService) List(ctx context.Context, adminSchoolID *uuid.UUID, limit, offset int) ([]model.Parent, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.parentRepo.List(ctx, adminSchoolID, limit, offset)
}

// Create tạo mới phụ huynh
func (s *ParentService) Create(ctx context.Context, userID, schoolID uuid.UUID, fullName, phone string) (*model.Parent, error) {
	id, err := s.parentRepo.Create(ctx, userID, schoolID, fullName, phone)
	if err != nil {
		return nil, err
	}

	return &model.Parent{
		ParentID: id,
		UserID:   userID,
		SchoolID: schoolID,
		FullName: fullName,
		Phone:    phone,
	}, nil
}

// GetByParentID lấy thông tin phụ huynh theo parent_id.
func (s *ParentService) GetByParentID(ctx context.Context, adminSchoolID *uuid.UUID, parentID uuid.UUID) (*model.Parent, error) {
	parent, err := s.parentRepo.GetByParentID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	// SCHOOL_ADMIN: validate parent thuộc cùng school với admin
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil && parent.SchoolID != *adminSchoolID {
		return nil, ErrSchoolAccessDenied
	}

	return parent, nil
}

// AssignStudent gán phụ huynh cho học sinh.
func (s *ParentService) AssignStudent(ctx context.Context, adminSchoolID *uuid.UUID, parentID, studentID uuid.UUID, relationship string) error {
	if parentID == uuid.Nil {
		return ErrInvalidUserID
	}

	if studentID == uuid.Nil {
		return ErrInvalidUserID
	}

	// SCHOOL_ADMIN: validate parent thuộc cùng school với admin
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil {
		parent, err := s.parentRepo.GetByParentID(ctx, parentID)
		if err != nil {
			return err
		}
		if parent.SchoolID != *adminSchoolID {
			return ErrSchoolAccessDenied
		}

		// validate student thuộc cùng school với admin
		student, err := s.studentRepo.GetByStudentID(ctx, studentID)
		if err != nil {
			return err
		}
		if student.SchoolID != *adminSchoolID {
			return ErrSchoolAccessDenied
		}
	}

	return s.studentParentRepo.Assign(ctx, studentID, parentID, relationship)
}

// UnassignStudent hủy gán phụ huynh khỏi học sinh.
//
// adminSchoolID == nil → SUPER_ADMIN: không giới hạn
//
// adminSchoolID != nil → SCHOOL_ADMIN: validate cả parent AND student thuộc cùng school với admin
func (s *ParentService) UnassignStudent(ctx context.Context, adminSchoolID *uuid.UUID, parentID, studentID uuid.UUID) error {
	if parentID == uuid.Nil {
		return ErrInvalidUserID
	}

	if studentID == uuid.Nil {
		return ErrInvalidUserID
	}

	// SCHOOL_ADMIN: validate parent thuộc school với admin
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil {
		parent, err := s.parentRepo.GetByParentID(ctx, parentID)
		if err != nil {
			return err
		}
		if parent.SchoolID != *adminSchoolID {
			return ErrSchoolAccessDenied
		}

		// validate student thuộc cùng school với admin
		student, err := s.studentRepo.GetByStudentID(ctx, studentID)
		if err != nil {
			return err
		}
		if student.SchoolID != *adminSchoolID {
			return ErrSchoolAccessDenied
		}
	}

	return s.studentParentRepo.Unassign(ctx, studentID, parentID)
}
