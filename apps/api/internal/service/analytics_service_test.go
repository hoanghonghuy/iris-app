package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

type fakeAnalyticsServiceSchoolRepo struct {
	countAllCalls int
	countAllRes   int
	countAllErr   error
}

func (f *fakeAnalyticsServiceSchoolRepo) CountAll(_ context.Context) (int, error) {
	f.countAllCalls++
	return f.countAllRes, f.countAllErr
}

type fakeAnalyticsServiceClassRepo struct {
	countBySchoolCalls int
	countBySchoolArg   *uuid.UUID
	countBySchoolRes   int
	countBySchoolErr   error
}

func (f *fakeAnalyticsServiceClassRepo) CountBySchool(_ context.Context, schoolID *uuid.UUID) (int, error) {
	f.countBySchoolCalls++
	f.countBySchoolArg = schoolID
	return f.countBySchoolRes, f.countBySchoolErr
}

type fakeAnalyticsServiceUserRepo struct {
	countCalls []struct {
		role     string
		schoolID *uuid.UUID
	}
	countRes map[string]int
	countErr map[string]error
}

func (f *fakeAnalyticsServiceUserRepo) CountUsersByRoleAndSchool(_ context.Context, role string, schoolID *uuid.UUID) (int, error) {
	f.countCalls = append(f.countCalls, struct {
		role     string
		schoolID *uuid.UUID
	}{role: role, schoolID: schoolID})
	if err, ok := f.countErr[role]; ok {
		if err != nil {
			return 0, err
		}
	}
	if v, ok := f.countRes[role]; ok {
		return v, nil
	}
	return 0, nil
}

type fakeAnalyticsServiceStudentRepo struct {
	countCalls int
	countArg   *uuid.UUID
	countRes   int
	countErr   error
}

func (f *fakeAnalyticsServiceStudentRepo) CountStudentsBySchool(_ context.Context, schoolID *uuid.UUID) (int, error) {
	f.countCalls++
	f.countArg = schoolID
	return f.countRes, f.countErr
}

type fakeAnalyticsServiceTeacherScopeRepo struct {
	listMyClassesCalls int
	listMyClassesArg   uuid.UUID
	listMyClassesRes   []model.Class
	listMyClassesErr   error

	countMyStudentsCalls int
	countMyStudentsArg   uuid.UUID
	countMyStudentsRes   int
	countMyStudentsErr   error

	countMyPostsCalls int
	countMyPostsArg   uuid.UUID
	countMyPostsRes   int
	countMyPostsErr   error
}

func (f *fakeAnalyticsServiceTeacherScopeRepo) ListMyClasses(_ context.Context, teacherUserID uuid.UUID) ([]model.Class, error) {
	f.listMyClassesCalls++
	f.listMyClassesArg = teacherUserID
	return f.listMyClassesRes, f.listMyClassesErr
}

func (f *fakeAnalyticsServiceTeacherScopeRepo) CountMyStudents(_ context.Context, teacherUserID uuid.UUID) (int, error) {
	f.countMyStudentsCalls++
	f.countMyStudentsArg = teacherUserID
	return f.countMyStudentsRes, f.countMyStudentsErr
}

func (f *fakeAnalyticsServiceTeacherScopeRepo) CountMyPosts(_ context.Context, teacherUserID uuid.UUID) (int, error) {
	f.countMyPostsCalls++
	f.countMyPostsArg = teacherUserID
	return f.countMyPostsRes, f.countMyPostsErr
}

func TestAnalyticsServiceGetAdminAnalytics(t *testing.T) {
	schoolID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name                 string
		schoolID             *uuid.UUID
		schoolCountErr       error
		classCountErr        error
		teacherCountErr      error
		parentCountErr       error
		studentCountErr      error
		wantErr              error
		wantSchoolCountCalls int
		wantTotalSchools     int
	}{
		{name: "super-admin school count error", schoolID: nil, schoolCountErr: sentinelErr, wantErr: sentinelErr, wantSchoolCountCalls: 1},
		{name: "class count error", schoolID: nil, classCountErr: sentinelErr, wantErr: sentinelErr, wantSchoolCountCalls: 1},
		{name: "teacher count error", schoolID: nil, teacherCountErr: sentinelErr, wantErr: sentinelErr, wantSchoolCountCalls: 1},
		{name: "parent count error", schoolID: nil, parentCountErr: sentinelErr, wantErr: sentinelErr, wantSchoolCountCalls: 1},
		{name: "student count error", schoolID: nil, studentCountErr: sentinelErr, wantErr: sentinelErr, wantSchoolCountCalls: 1},
		{name: "success super-admin", schoolID: nil, wantSchoolCountCalls: 1, wantTotalSchools: 5},
		{name: "success school-admin bypass school count", schoolID: &schoolID, wantSchoolCountCalls: 0, wantTotalSchools: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			schoolRepo := &fakeAnalyticsServiceSchoolRepo{countAllRes: 5, countAllErr: tc.schoolCountErr}
			classRepo := &fakeAnalyticsServiceClassRepo{countBySchoolRes: 7, countBySchoolErr: tc.classCountErr}
			userRepo := &fakeAnalyticsServiceUserRepo{
				countRes: map[string]int{"TEACHER": 11, "PARENT": 13},
				countErr: map[string]error{"TEACHER": tc.teacherCountErr, "PARENT": tc.parentCountErr},
			}
			studentRepo := &fakeAnalyticsServiceStudentRepo{countRes: 17, countErr: tc.studentCountErr}
			svc := &AnalyticsService{schoolRepo: schoolRepo, classRepo: classRepo, userRepo: userRepo, studentRepo: studentRepo, teacherScopeRepo: &fakeAnalyticsServiceTeacherScopeRepo{}}

			got, err := svc.GetAdminAnalytics(context.Background(), tc.schoolID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
				if got != nil {
					t.Fatalf("expected nil analytics when error, got %#v", got)
				}
			} else if err != nil {
				t.Fatalf("GetAdminAnalytics() error = %v", err)
			}

			if schoolRepo.countAllCalls != tc.wantSchoolCountCalls {
				t.Fatalf("school count calls = %d, want %d", schoolRepo.countAllCalls, tc.wantSchoolCountCalls)
			}

			if classRepo.countBySchoolCalls > 0 && classRepo.countBySchoolArg != tc.schoolID {
				t.Fatalf("class count schoolID pointer mismatch")
			}
			if studentRepo.countCalls > 0 && studentRepo.countArg != tc.schoolID {
				t.Fatalf("student count schoolID pointer mismatch")
			}

			if tc.wantErr == nil {
				if got == nil {
					t.Fatal("expected analytics on success")
				}
				if got.TotalSchools != tc.wantTotalSchools || got.TotalClasses != 7 || got.TotalTeachers != 11 || got.TotalParents != 13 || got.TotalStudents != 17 {
					t.Fatalf("unexpected analytics = %#v", got)
				}
				if len(userRepo.countCalls) != 2 {
					t.Fatalf("user role count calls = %d, want 2", len(userRepo.countCalls))
				}
				if userRepo.countCalls[0].role != "TEACHER" || userRepo.countCalls[1].role != "PARENT" {
					t.Fatalf("role query order mismatch: %#v", userRepo.countCalls)
				}
			}
		})
	}
}

func TestAnalyticsServiceGetTeacherAnalytics(t *testing.T) {
	teacherUserID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name                 string
		listClassesErr       error
		countStudentsErr     error
		countPostsErr        error
		wantErr              error
		wantListCalls        int
		wantCountStudentCall int
		wantCountPostCall    int
	}{
		{name: "list classes error", listClassesErr: sentinelErr, wantErr: sentinelErr, wantListCalls: 1},
		{name: "count students error", countStudentsErr: sentinelErr, wantErr: sentinelErr, wantListCalls: 1, wantCountStudentCall: 1},
		{name: "count posts error", countPostsErr: sentinelErr, wantErr: sentinelErr, wantListCalls: 1, wantCountStudentCall: 1, wantCountPostCall: 1},
		{name: "success", wantListCalls: 1, wantCountStudentCall: 1, wantCountPostCall: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			teacherScopeRepo := &fakeAnalyticsServiceTeacherScopeRepo{
				listMyClassesRes:   []model.Class{{ClassID: uuid.New()}, {ClassID: uuid.New()}},
				listMyClassesErr:   tc.listClassesErr,
				countMyStudentsRes: 23,
				countMyStudentsErr: tc.countStudentsErr,
				countMyPostsRes:    9,
				countMyPostsErr:    tc.countPostsErr,
			}
			svc := &AnalyticsService{schoolRepo: &fakeAnalyticsServiceSchoolRepo{}, classRepo: &fakeAnalyticsServiceClassRepo{}, userRepo: &fakeAnalyticsServiceUserRepo{}, studentRepo: &fakeAnalyticsServiceStudentRepo{}, teacherScopeRepo: teacherScopeRepo}

			got, err := svc.GetTeacherAnalytics(context.Background(), teacherUserID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
				if got != nil {
					t.Fatalf("expected nil analytics on error, got %#v", got)
				}
			} else if err != nil {
				t.Fatalf("GetTeacherAnalytics() error = %v", err)
			}

			if teacherScopeRepo.listMyClassesCalls != tc.wantListCalls || teacherScopeRepo.countMyStudentsCalls != tc.wantCountStudentCall || teacherScopeRepo.countMyPostsCalls != tc.wantCountPostCall {
				t.Fatalf("calls list/students/posts = %d/%d/%d, want %d/%d/%d", teacherScopeRepo.listMyClassesCalls, teacherScopeRepo.countMyStudentsCalls, teacherScopeRepo.countMyPostsCalls, tc.wantListCalls, tc.wantCountStudentCall, tc.wantCountPostCall)
			}
			if teacherScopeRepo.listMyClassesCalls > 0 && teacherScopeRepo.listMyClassesArg != teacherUserID {
				t.Fatalf("teacher userID not forwarded to list classes")
			}
			if teacherScopeRepo.countMyStudentsCalls > 0 && teacherScopeRepo.countMyStudentsArg != teacherUserID {
				t.Fatalf("teacher userID not forwarded to count students")
			}
			if teacherScopeRepo.countMyPostsCalls > 0 && teacherScopeRepo.countMyPostsArg != teacherUserID {
				t.Fatalf("teacher userID not forwarded to count posts")
			}

			if tc.wantErr == nil {
				if got == nil {
					t.Fatal("expected analytics on success")
				}
				if got.TotalClasses != 2 || got.TotalStudents != 23 || got.TotalPosts != 9 {
					t.Fatalf("unexpected teacher analytics = %#v", got)
				}
			}
		})
	}
}
