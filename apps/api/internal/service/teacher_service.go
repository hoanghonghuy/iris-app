package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type TeacherService struct {
	teacherRepo      *repo.TeacherRepo
	teacherClassRepo *repo.TeacherClassRepo
	classRepo        *repo.ClassRepo
}

func NewTeacherService(teacherRepo *repo.TeacherRepo, teacherClassRepo *repo.TeacherClassRepo, classRepo *repo.ClassRepo) *TeacherService {
	return &TeacherService{
		teacherRepo:      teacherRepo,
		teacherClassRepo: teacherClassRepo,
		classRepo:        classRepo,
	}
}

// List lấy danh sách giáo viên.
func (s *TeacherService) List(ctx context.Context, adminSchoolID *uuid.UUID, limit, offset int) ([]model.Teacher, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.teacherRepo.List(ctx, adminSchoolID, limit, offset)
}

// Assign gán giáo viên vào lớp.
func (s *TeacherService) Assign(ctx context.Context, adminSchoolID *uuid.UUID, teacherID, classID uuid.UUID) error {
	// kiểm tra teacher có tồn tại không
	teacher, err := s.teacherRepo.GetByTeacherID(ctx, teacherID)
	if err != nil {
		return err
	}

	// SCHOOL_ADMIN: validate teacher thuộc cùng school với admin
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil && teacher.SchoolID != *adminSchoolID {
		return ErrSchoolAccessDenied
	}

	// SCHOOL_ADMIN: validate class thuộc cùng school với admin
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil {
		class, err := s.classRepo.GetByClassID(ctx, classID)
		if err != nil {
			return ErrInvalidClassID
		}
		if class.SchoolID != *adminSchoolID {
			return ErrSchoolAccessDenied
		}
	}

	return s.teacherClassRepo.Assign(ctx, teacherID, classID)
}

// ListTeachersOfClass lấy danh sách giáo viên của một lớp.
func (s *TeacherService) ListTeachersOfClass(ctx context.Context, adminSchoolID *uuid.UUID, classID uuid.UUID) ([]model.Teacher, error) {
	// SCHOOL_ADMIN: validate class thuộc cùng school với admin
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil {
		class, err := s.classRepo.GetByClassID(ctx, classID)
		if err != nil {
			return nil, ErrInvalidClassID
		}
		if class.SchoolID != *adminSchoolID {
			return nil, ErrSchoolAccessDenied
		}
	}

	return s.teacherClassRepo.ListTeacherDetailsOfClass(ctx, classID)
}

// Unassign hủy gán giáo viên khỏi lớp.
func (s *TeacherService) Unassign(ctx context.Context, adminSchoolID *uuid.UUID, teacherID, classID uuid.UUID) error {
	// SCHOOL_ADMIN: validate teacher thuộc cùng school với admin
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil {
		teacher, err := s.teacherRepo.GetByTeacherID(ctx, teacherID)
		if err != nil {
			return err
		}
		if teacher.SchoolID != *adminSchoolID {
			return ErrSchoolAccessDenied
		}

		// validate class thuộc cùng school với admin
		class, err := s.classRepo.GetByClassID(ctx, classID)
		if err != nil {
			return ErrInvalidClassID
		}
		if class.SchoolID != *adminSchoolID {
			return ErrSchoolAccessDenied
		}
	}

	exists, err := s.teacherClassRepo.IsTeacherAssignedToClass(ctx, teacherID, classID)
	if err != nil {
		return err
	}
	// nếu chưa assign (exists == false)
	if !exists {
		return ErrTeacherNotAssigned
	}

	return s.teacherClassRepo.Unassign(ctx, teacherID, classID)
}

// GetByTeacherID lấy thông tin giáo viên theo ID.
func (s *TeacherService) GetByTeacherID(ctx context.Context, adminSchoolID *uuid.UUID, teacherID uuid.UUID) (*model.Teacher, error) {
	teacher, err := s.teacherRepo.GetByTeacherID(ctx, teacherID)
	if err != nil {
		return nil, err
	}

	// SCHOOL_ADMIN: validate teacher thuộc cùng school với admin
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil && teacher.SchoolID != *adminSchoolID {
		return nil, ErrSchoolAccessDenied
	}

	return teacher, nil
}

// Update cập nhật thông tin giáo viên (admin only).
func (s *TeacherService) Update(ctx context.Context, adminSchoolID *uuid.UUID, teacherID uuid.UUID, fullName, phone string, schoolID uuid.UUID) error {
	teacher, err := s.teacherRepo.GetByTeacherID(ctx, teacherID)
	if err != nil {
		return err
	}

	// SCHOOL_ADMIN: validate teacher thuộc cùng school với admin
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil && teacher.SchoolID != *adminSchoolID {
		return ErrSchoolAccessDenied
	}

	// SCHOOL_ADMIN: không được đổi school_id sang school khác
	// adminSchoolID == nil => SUPER_ADMIN: không cần validate
	if adminSchoolID != nil && schoolID != *adminSchoolID {
		return ErrSchoolAccessDenied
	}

	return s.teacherRepo.Update(ctx, teacherID, fullName, phone, schoolID)
}
