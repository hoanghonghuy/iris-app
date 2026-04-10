package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type AnalyticsService struct {
	repos *repo.Repositories
}

func NewAnalyticsService(repos *repo.Repositories) *AnalyticsService {
	return &AnalyticsService{
		repos: repos,
	}
}

// GetAdminAnalytics trả về các chỉ số thống kê cho trang Dashboard của Admin.
// Trả về nil schoolID nếu là SUPER_ADMIN để đếm toàn hệ thống, hoặc truyền schoolID nếu là SCHOOL_ADMIN.
func (s *AnalyticsService) GetAdminAnalytics(ctx context.Context, schoolID *uuid.UUID) (*model.AdminAnalytics, error) {
	var err error
	var totalSchools, totalClasses, totalTeachers, totalStudents, totalParents int
	var schoolName string
	isSuperAdmin := schoolID == nil

	// Thống kê trường học
	if schoolID == nil {
		totalSchools, err = s.repos.SchoolRepo.CountAll(ctx)
		if err != nil {
			return nil, err
		}
		schoolName = "Toan he thong"
	} else {
		// SCHOOL_ADMIN chỉ quản lý 1 trường của mình
		totalSchools = 1
		school, schoolErr := s.repos.SchoolRepo.GetByID(ctx, *schoolID)
		if schoolErr != nil {
			return nil, schoolErr
		}
		schoolName = school.Name
	}

	// Thống kê lớp học
	totalClasses, err = s.repos.ClassRepo.CountBySchool(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	// Thống kê Giáo viên
	totalTeachers, err = s.repos.UserRepo.CountUsersByRoleAndSchool(ctx, "TEACHER", schoolID)
	if err != nil {
		return nil, err
	}

	// Thống kê Phụ huynh
	totalParents, err = s.repos.UserRepo.CountUsersByRoleAndSchool(ctx, "PARENT", schoolID)
	if err != nil {
		return nil, err
	}

	// Thống kê Học sinh
	totalStudents, err = s.repos.StudentRepo.CountStudentsBySchool(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	todayPresent, err := s.repos.StudentRepo.CountTodayAttendancePresentBySchool(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	recentHealthAlerts24h, err := s.repos.HealthLogRepo.CountRecentAlertsBySchool(ctx, schoolID, time.Now().Add(-24*time.Hour))
	if err != nil {
		return nil, err
	}

	todayAttendanceRate := 0.0
	if totalStudents > 0 {
		todayAttendanceRate = float64(todayPresent) * 100 / float64(totalStudents)
	}

	return &model.AdminAnalytics{
		TotalSchools:          totalSchools,
		TotalClasses:          totalClasses,
		TotalTeachers:         totalTeachers,
		TotalStudents:         totalStudents,
		TotalParents:          totalParents,
		IsSuperAdmin:          isSuperAdmin,
		SchoolName:            schoolName,
		TodayAttendanceRate:   todayAttendanceRate,
		RecentHealthAlerts24h: recentHealthAlerts24h,
	}, nil
}

// GetTeacherAnalytics trả về các chỉ số thống kê cho trang Dashboard của Giáo viên.
func (s *AnalyticsService) GetTeacherAnalytics(ctx context.Context, teacherUserID uuid.UUID) (*model.TeacherAnalytics, error) {
	var err error
	var totalClasses, totalStudents, totalPosts int

	classes, err := s.repos.TeacherScopeRepo.ListMyClasses(ctx, teacherUserID)
	if err != nil {
		return nil, err
	}
	totalClasses = len(classes)

	totalStudents, err = s.repos.TeacherScopeRepo.CountMyStudents(ctx, teacherUserID)
	if err != nil {
		return nil, err
	}

	totalPosts, err = s.repos.TeacherScopeRepo.CountMyPosts(ctx, teacherUserID)
	if err != nil {
		return nil, err
	}

	todayMarked, err := s.repos.TeacherScopeRepo.CountTodayAttendanceMarked(ctx, teacherUserID)
	if err != nil {
		return nil, err
	}

	pendingAppointments, err := s.repos.AppointmentRepo.CountTeacherPendingAppointments(ctx, teacherUserID)
	if err != nil {
		return nil, err
	}

	recentHealthAlerts24h, err := s.repos.TeacherScopeRepo.CountRecentHealthAlerts(ctx, teacherUserID, time.Now().Add(-24*time.Hour))
	if err != nil {
		return nil, err
	}

	todayPending := totalStudents - todayMarked
	if todayPending < 0 {
		todayPending = 0
	}

	return &model.TeacherAnalytics{
		TotalClasses:                totalClasses,
		TotalStudents:               totalStudents,
		TotalPosts:                  totalPosts,
		TodayAttendanceMarkedCount:  todayMarked,
		TodayAttendancePendingCount: todayPending,
		PendingAppointments:         pendingAppointments,
		RecentHealthAlerts24h:       recentHealthAlerts24h,
	}, nil
}
