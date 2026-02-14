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
}

func NewTeacherService(teacherRepo *repo.TeacherRepo, teacherClassRepo *repo.TeacherClassRepo) *TeacherService {
	return &TeacherService{
		teacherRepo:      teacherRepo,
		teacherClassRepo: teacherClassRepo,
	}
}

func (s *TeacherService) List(ctx context.Context, limit, offset int) ([]model.Teacher, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.teacherRepo.List(ctx, limit, offset)
}

func (s *TeacherService) Assign(ctx context.Context, teacherID, classID uuid.UUID) error {
	// kiểm tra teacher có tồn tại không
	_, err := s.teacherRepo.GetByTeacherID(ctx, teacherID)
	if err != nil {
		return err
	}

	// validate teacher (status active, etc.)
	// if teacher.Status != "active" {
	//    return errors.New("teacher không active")
	// }

	return s.teacherClassRepo.Assign(ctx, teacherID, classID)
}

func (s *TeacherService) ListTeachersOfClass(ctx context.Context, classID uuid.UUID) ([]model.Teacher, error) {
	return s.teacherClassRepo.ListTeacherDetailsOfClass(ctx, classID)
}

func (s *TeacherService) Unassign(ctx context.Context, teacherID, classID uuid.UUID) error {
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

func (s *TeacherService) GetByTeacherID(ctx context.Context, teacherID uuid.UUID) (*model.Teacher, error) {
	return s.teacherRepo.GetByTeacherID(ctx, teacherID)
}

// Update updates a teacher's information (admin only - can update all fields)
func (s *TeacherService) Update(ctx context.Context, teacherID uuid.UUID, fullName, phone string, schoolID uuid.UUID) error {
	// Check if teacher exists
	_, err := s.teacherRepo.GetByTeacherID(ctx, teacherID)
	if err != nil {
		return err
	}

	// Admin can update all fields
	return s.teacherRepo.Update(ctx, teacherID, fullName, phone, schoolID)
}
