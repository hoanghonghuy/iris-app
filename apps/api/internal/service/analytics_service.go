package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type AnalyticsService struct {
	schoolRepo       analyticsServiceSchoolRepo
	classRepo        analyticsServiceClassRepo
	userRepo         analyticsServiceUserRepo
	studentRepo      analyticsServiceStudentRepo
	teacherScopeRepo teacherScopeStatsRepo
}

type analyticsServiceSchoolRepo interface {
	CountAll(ctx context.Context) (int, error)
}

type analyticsServiceClassRepo interface {
	CountBySchool(ctx context.Context, schoolID *uuid.UUID) (int, error)
}

type analyticsServiceUserRepo interface {
	CountUsersByRoleAndSchool(ctx context.Context, role string, schoolID *uuid.UUID) (int, error)
}

type analyticsServiceStudentRepo interface {
	CountStudentsBySchool(ctx context.Context, schoolID *uuid.UUID) (int, error)
}

func NewAnalyticsService(repos *repo.Repositories) *AnalyticsService {
	return &AnalyticsService{
		schoolRepo:       repos.SchoolRepo,
		classRepo:        repos.ClassRepo,
		userRepo:         repos.UserRepo,
		studentRepo:      repos.StudentRepo,
		teacherScopeRepo: repos.TeacherScopeRepo,
	}
}

// GetAdminAnalytics trả về các chỉ số thống kê cho trang Dashboard của Admin.
// Trả về nil schoolID nếu là SUPER_ADMIN để đếm toàn hệ thống, hoặc truyền schoolID nếu là SCHOOL_ADMIN.
func (s *AnalyticsService) GetAdminAnalytics(ctx context.Context, schoolID *uuid.UUID) (*model.AdminAnalytics, error) {
	var err error
	var totalSchools, totalClasses, totalTeachers, totalStudents, totalParents int

	// Thống kê trường học
	if schoolID == nil {
		totalSchools, err = s.schoolRepo.CountAll(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		// SCHOOL_ADMIN chỉ quản lý 1 trường của mình
		totalSchools = 1
	}

	// Thống kê lớp học
	totalClasses, err = s.classRepo.CountBySchool(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	// Thống kê Giáo viên
	totalTeachers, err = s.userRepo.CountUsersByRoleAndSchool(ctx, "TEACHER", schoolID)
	if err != nil {
		return nil, err
	}

	// Thống kê Phụ huynh
	totalParents, err = s.userRepo.CountUsersByRoleAndSchool(ctx, "PARENT", schoolID)
	if err != nil {
		return nil, err
	}

	// Thống kê Học sinh
	totalStudents, err = s.studentRepo.CountStudentsBySchool(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	return &model.AdminAnalytics{
		TotalSchools:  totalSchools,
		TotalClasses:  totalClasses,
		TotalTeachers: totalTeachers,
		TotalStudents: totalStudents,
		TotalParents:  totalParents,
	}, nil
}

// GetTeacherAnalytics trả về các chỉ số thống kê cho trang Dashboard của Giáo viên.
func (s *AnalyticsService) GetTeacherAnalytics(ctx context.Context, teacherUserID uuid.UUID) (*model.TeacherAnalytics, error) {
	var err error
	var totalClasses, totalStudents, totalPosts int

	classes, err := s.teacherScopeRepo.ListMyClasses(ctx, teacherUserID)
	if err != nil {
		return nil, err
	}
	totalClasses = len(classes)

	totalStudents, err = s.teacherScopeRepo.CountMyStudents(ctx, teacherUserID)
	if err != nil {
		return nil, err
	}

	totalPosts, err = s.teacherScopeRepo.CountMyPosts(ctx, teacherUserID)
	if err != nil {
		return nil, err
	}

	return &model.TeacherAnalytics{
		TotalClasses:  totalClasses,
		TotalStudents: totalStudents,
		TotalPosts:    totalPosts,
	}, nil
}
