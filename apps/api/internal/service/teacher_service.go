package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type TeacherService struct {
	TeacherRepo      *repo.TeacherRepo
	TeacherClassRepo *repo.TeacherClassRepo
}

// UpdateTeacherRequest represents the request to update a teacher (admin only)
type UpdateTeacherRequest struct {
	FullName string    `json:"full_name"`
	Phone    string    `json:"phone"`
	SchoolID uuid.UUID `json:"school_id"`
}

func (s *TeacherService) List(ctx context.Context) ([]model.Teacher, error) {
	return s.TeacherRepo.List(ctx)
}

func (s *TeacherService) Assign(ctx context.Context, teacherID, classID uuid.UUID) error {
	// kiểm tra teacher có tồn tại không
	_, err := s.TeacherRepo.GetByTeacherID(ctx, teacherID)
	if err != nil {
		return err
	}

	// validate teacher (status active, etc.)
	// if teacher.Status != "active" {
	//    return errors.New("teacher không active")
	// }

	return s.TeacherClassRepo.Assign(ctx, teacherID, classID)
}

func (s *TeacherService) ListTeachersOfClass(ctx context.Context, classID uuid.UUID) ([]model.Teacher, error) {
	teacherIDs, err := s.TeacherClassRepo.ListTeachersOfClass(ctx, classID)
	if err != nil {
		return nil, err
	}

	var teachers []model.Teacher

	// duyệt qua danh sách các ID giáo viên (teacherIDs) của một lớp học cụ thể.
	// với mỗi teacherID, gọi hàm GetByTeacherID để lấy thông tin chi tiết của giáo viên từ repository.
	// lỗi khi lấy thông tin giáo viên => trả về lỗi.
	// không lỗi => giáo viên được thêm vào danh sách teachers.
	for _, teacherID := range teacherIDs {
		teacher, err := s.TeacherRepo.GetByTeacherID(ctx, teacherID)
		if err != nil {
			// continue
			return nil, err
		}
		teachers = append(teachers, *teacher)
	}

	return teachers, nil
}

func (s *TeacherService) Unassign(ctx context.Context, teacherID, classID uuid.UUID) error {
	exists, err := s.TeacherClassRepo.IsTeacherAssignedToClass(ctx, teacherID, classID)
	if err != nil {
		return err
	}
	// nếu chưa assign (exists == false)
	if !exists {
		return errors.New("teacher is not assigned to this class")
	}

	return s.TeacherClassRepo.Unassign(ctx, teacherID, classID)
}

func (s *TeacherService) GetByTeacherID(ctx context.Context, teacherID uuid.UUID) (*model.Teacher, error) {
	return s.TeacherRepo.GetByTeacherID(ctx, teacherID)
}

// Update updates a teacher's information (admin only - can update all fields)
func (s *TeacherService) Update(ctx context.Context, teacherID uuid.UUID, req UpdateTeacherRequest) error {
	// Check if teacher exists
	_, err := s.TeacherRepo.GetByTeacherID(ctx, teacherID)
	if err != nil {
		return err
	}

	// Admin can update all fields
	return s.TeacherRepo.Update(ctx, teacherID, req.FullName, req.Phone, req.SchoolID)
}
