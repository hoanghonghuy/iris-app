package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5"
)

type fakeStudentServiceStudentRepo struct {
	createCalls  int
	createSchool uuid.UUID
	createClass  uuid.UUID
	createName   string
	createDOB    time.Time
	createGender string
	createResult uuid.UUID
	createErr    error

	listByClassCalls  int
	listByClassClass  uuid.UUID
	listByClassLimit  int
	listByClassOffset int
	listByClassResult []model.Student
	listByClassTotal  int
	listByClassErr    error

	getProfileCalls  int
	getProfileArg    uuid.UUID
	getProfileResult *model.StudentProfile
	getProfileErr    error

	getByStudentCalls  int
	getByStudentArg    uuid.UUID
	getByStudentResult *model.Student
	getByStudentErr    error

	updateCalls  int
	updateArg    uuid.UUID
	updateName   string
	updateDOB    time.Time
	updateGender string
	updateErr    error

	deleteCalls int
	deleteArg   uuid.UUID
	deleteErr   error
}

func (f *fakeStudentServiceStudentRepo) Create(_ context.Context, schoolID, classID uuid.UUID, fullName string, dob time.Time, gender string) (uuid.UUID, error) {
	f.createCalls++
	f.createSchool = schoolID
	f.createClass = classID
	f.createName = fullName
	f.createDOB = dob
	f.createGender = gender
	return f.createResult, f.createErr
}

func (f *fakeStudentServiceStudentRepo) ListByClass(_ context.Context, classID uuid.UUID, limit, offset int) ([]model.Student, int, error) {
	f.listByClassCalls++
	f.listByClassClass = classID
	f.listByClassLimit = limit
	f.listByClassOffset = offset
	return f.listByClassResult, f.listByClassTotal, f.listByClassErr
}

func (f *fakeStudentServiceStudentRepo) GetStudentProfile(_ context.Context, studentID uuid.UUID) (*model.StudentProfile, error) {
	f.getProfileCalls++
	f.getProfileArg = studentID
	return f.getProfileResult, f.getProfileErr
}

func (f *fakeStudentServiceStudentRepo) GetByStudentID(_ context.Context, studentID uuid.UUID) (*model.Student, error) {
	f.getByStudentCalls++
	f.getByStudentArg = studentID
	return f.getByStudentResult, f.getByStudentErr
}

func (f *fakeStudentServiceStudentRepo) Update(_ context.Context, studentID uuid.UUID, fullName string, dob time.Time, gender string) error {
	f.updateCalls++
	f.updateArg = studentID
	f.updateName = fullName
	f.updateDOB = dob
	f.updateGender = gender
	return f.updateErr
}

func (f *fakeStudentServiceStudentRepo) Delete(_ context.Context, studentID uuid.UUID) error {
	f.deleteCalls++
	f.deleteArg = studentID
	return f.deleteErr
}

type fakeStudentServiceClassRepo struct {
	getByClassCalls  int
	getByClassArg    uuid.UUID
	getByClassResult *model.Class
	getByClassErr    error
}

func (f *fakeStudentServiceClassRepo) GetByClassID(_ context.Context, classID uuid.UUID) (*model.Class, error) {
	f.getByClassCalls++
	f.getByClassArg = classID
	return f.getByClassResult, f.getByClassErr
}

func TestStudentServiceCreate(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	classID := uuid.New()
	studentID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		adminSchoolID   *uuid.UUID
		schoolID        uuid.UUID
		dob             string
		classResult     *model.Class
		classErr        error
		createErr       error
		wantErr         error
		wantGetClass    int
		wantCreateCalls int
	}{
		{name: "school admin cross school denied", adminSchoolID: &schoolA, schoolID: schoolB, dob: "2010-01-01", wantErr: ErrSchoolAccessDenied},
		{name: "school admin invalid class", adminSchoolID: &schoolA, schoolID: schoolA, dob: "2010-01-01", classErr: sentinelErr, wantErr: ErrInvalidClassID, wantGetClass: 1},
		{name: "school admin class cross school denied", adminSchoolID: &schoolA, schoolID: schoolA, dob: "2010-01-01", classResult: &model.Class{ClassID: classID, SchoolID: schoolB}, wantErr: ErrSchoolAccessDenied, wantGetClass: 1},
		{name: "invalid dob", adminSchoolID: nil, schoolID: schoolA, dob: "invalid", wantErr: ErrInvalidValue},
		{name: "create error passthrough", adminSchoolID: nil, schoolID: schoolA, dob: "2010-01-01", createErr: sentinelErr, wantErr: sentinelErr, wantCreateCalls: 1},
		{name: "success", adminSchoolID: &schoolA, schoolID: schoolA, dob: "2010-01-01", classResult: &model.Class{ClassID: classID, SchoolID: schoolA}, wantGetClass: 1, wantCreateCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			studentRepo := &fakeStudentServiceStudentRepo{createResult: studentID, createErr: tc.createErr}
			classRepo := &fakeStudentServiceClassRepo{getByClassResult: tc.classResult, getByClassErr: tc.classErr}
			svc := &StudentService{studentRepo: studentRepo, classRepo: classRepo}

			got, err := svc.Create(context.Background(), tc.adminSchoolID, tc.schoolID, classID, "Student A", tc.dob, "male")
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("Create() error = %v", err)
			}

			if classRepo.getByClassCalls != tc.wantGetClass || studentRepo.createCalls != tc.wantCreateCalls {
				t.Fatalf("calls class/create = %d/%d, want %d/%d", classRepo.getByClassCalls, studentRepo.createCalls, tc.wantGetClass, tc.wantCreateCalls)
			}
			if tc.wantCreateCalls > 0 {
				if studentRepo.createSchool != tc.schoolID || studentRepo.createClass != classID || studentRepo.createName != "Student A" || studentRepo.createGender != "male" {
					t.Fatalf("unexpected create args: school=%v class=%v name=%q gender=%q", studentRepo.createSchool, studentRepo.createClass, studentRepo.createName, studentRepo.createGender)
				}
			}
			if tc.wantErr == nil {
				if got == nil || got.StudentID != studentID || got.SchoolID != tc.schoolID || got.CurrentClassID != classID {
					t.Fatalf("unexpected student = %#v", got)
				}
			}
		})
	}
}

func TestStudentServiceListByClass(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	classID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name           string
		adminSchoolID  *uuid.UUID
		limit          int
		offset         int
		classResult    *model.Class
		classErr       error
		listErr        error
		wantErr        error
		wantGetClass   int
		wantListCalls  int
		wantLimit      int
		wantOffset     int
		wantResultSize int
		wantTotal      int
	}{
		{name: "school admin invalid class", adminSchoolID: &schoolA, limit: 20, offset: 0, classErr: sentinelErr, wantErr: ErrInvalidClassID, wantGetClass: 1},
		{name: "school admin class cross school denied", adminSchoolID: &schoolA, limit: 20, offset: 0, classResult: &model.Class{ClassID: classID, SchoolID: schoolB}, wantErr: ErrSchoolAccessDenied, wantGetClass: 1},
		{name: "normalize default", adminSchoolID: nil, limit: 0, offset: -5, wantListCalls: 1, wantLimit: 20, wantOffset: 0, wantResultSize: 1, wantTotal: 1},
		{name: "clamp max", adminSchoolID: &schoolA, limit: 999, offset: 3, classResult: &model.Class{ClassID: classID, SchoolID: schoolA}, wantGetClass: 1, wantListCalls: 1, wantLimit: 100, wantOffset: 3, wantResultSize: 1, wantTotal: 1},
		{name: "list error passthrough", adminSchoolID: nil, limit: 50, offset: 2, listErr: sentinelErr, wantErr: sentinelErr, wantListCalls: 1, wantLimit: 50, wantOffset: 2},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			studentRepo := &fakeStudentServiceStudentRepo{listByClassResult: []model.Student{{StudentID: uuid.New()}}, listByClassTotal: 1, listByClassErr: tc.listErr}
			classRepo := &fakeStudentServiceClassRepo{getByClassResult: tc.classResult, getByClassErr: tc.classErr}
			svc := &StudentService{studentRepo: studentRepo, classRepo: classRepo}

			items, total, err := svc.ListByClass(context.Background(), tc.adminSchoolID, classID, tc.limit, tc.offset)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("ListByClass() error = %v", err)
			}
			if classRepo.getByClassCalls != tc.wantGetClass || studentRepo.listByClassCalls != tc.wantListCalls {
				t.Fatalf("calls class/list = %d/%d, want %d/%d", classRepo.getByClassCalls, studentRepo.listByClassCalls, tc.wantGetClass, tc.wantListCalls)
			}
			if tc.wantListCalls > 0 {
				if studentRepo.listByClassLimit != tc.wantLimit || studentRepo.listByClassOffset != tc.wantOffset {
					t.Fatalf("limit/offset = %d/%d, want %d/%d", studentRepo.listByClassLimit, studentRepo.listByClassOffset, tc.wantLimit, tc.wantOffset)
				}
			}
			if tc.wantErr != nil {
				return
			}
			if len(items) != tc.wantResultSize || total != tc.wantTotal {
				t.Fatalf("items/total = %d/%d, want %d/%d", len(items), total, tc.wantResultSize, tc.wantTotal)
			}
		})
	}
}

func TestStudentServiceGetProfile(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	studentID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name          string
		adminSchoolID *uuid.UUID
		profileResult *model.StudentProfile
		profileErr    error
		wantErr       error
	}{
		{name: "profile error maps failed to get student", adminSchoolID: nil, profileErr: sentinelErr, wantErr: ErrFailedToGetStudent},
		{name: "school admin cross school denied", adminSchoolID: &schoolA, profileResult: &model.StudentProfile{Student: model.Student{StudentID: studentID, SchoolID: schoolB}}, wantErr: ErrSchoolAccessDenied},
		{name: "success", adminSchoolID: &schoolA, profileResult: &model.StudentProfile{Student: model.Student{StudentID: studentID, SchoolID: schoolA}}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			studentRepo := &fakeStudentServiceStudentRepo{getProfileResult: tc.profileResult, getProfileErr: tc.profileErr}
			svc := &StudentService{studentRepo: studentRepo, classRepo: &fakeStudentServiceClassRepo{}}

			got, err := svc.GetProfile(context.Background(), tc.adminSchoolID, studentID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("GetProfile() error = %v", err)
			}
			if tc.wantErr == nil && (got == nil || got.StudentID != studentID) {
				t.Fatalf("unexpected profile = %#v", got)
			}
		})
	}
}

func TestStudentServiceUpdate(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	studentID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		adminSchoolID   *uuid.UUID
		dob             string
		getResult       *model.Student
		getErr          error
		updateErr       error
		wantErr         error
		wantGetCalls    int
		wantUpdateCalls int
	}{
		{name: "student not found", adminSchoolID: nil, dob: "2010-01-01", getErr: pgx.ErrNoRows, wantErr: ErrStudentNotFound, wantGetCalls: 1},
		{name: "get error maps failed to get student", adminSchoolID: nil, dob: "2010-01-01", getErr: sentinelErr, wantErr: ErrFailedToGetStudent, wantGetCalls: 1},
		{name: "school admin cross school denied", adminSchoolID: &schoolA, dob: "2010-01-01", getResult: &model.Student{StudentID: studentID, SchoolID: schoolB}, wantErr: ErrSchoolAccessDenied, wantGetCalls: 1},
		{name: "invalid dob", adminSchoolID: nil, dob: "invalid", getResult: &model.Student{StudentID: studentID, SchoolID: schoolA}, wantErr: ErrInvalidValue, wantGetCalls: 1},
		{name: "update no rows maps student not found", adminSchoolID: nil, dob: "2010-01-01", getResult: &model.Student{StudentID: studentID, SchoolID: schoolA}, updateErr: pgx.ErrNoRows, wantErr: ErrStudentNotFound, wantGetCalls: 1, wantUpdateCalls: 1},
		{name: "update error passthrough", adminSchoolID: nil, dob: "2010-01-01", getResult: &model.Student{StudentID: studentID, SchoolID: schoolA}, updateErr: sentinelErr, wantErr: sentinelErr, wantGetCalls: 1, wantUpdateCalls: 1},
		{name: "success", adminSchoolID: &schoolA, dob: "2010-01-01", getResult: &model.Student{StudentID: studentID, SchoolID: schoolA}, wantGetCalls: 1, wantUpdateCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			studentRepo := &fakeStudentServiceStudentRepo{getByStudentResult: tc.getResult, getByStudentErr: tc.getErr, updateErr: tc.updateErr}
			svc := &StudentService{studentRepo: studentRepo, classRepo: &fakeStudentServiceClassRepo{}}

			err := svc.Update(context.Background(), tc.adminSchoolID, studentID, "Student A", tc.dob, "female")
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("Update() error = %v", err)
			}
			if studentRepo.getByStudentCalls != tc.wantGetCalls || studentRepo.updateCalls != tc.wantUpdateCalls {
				t.Fatalf("calls get/update = %d/%d, want %d/%d", studentRepo.getByStudentCalls, studentRepo.updateCalls, tc.wantGetCalls, tc.wantUpdateCalls)
			}
		})
	}
}

func TestStudentServiceDelete(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	studentID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		adminSchoolID   *uuid.UUID
		getResult       *model.Student
		getErr          error
		deleteErr       error
		wantErr         error
		wantGetCalls    int
		wantDeleteCalls int
	}{
		{name: "student not found", adminSchoolID: nil, getErr: pgx.ErrNoRows, wantErr: ErrStudentNotFound, wantGetCalls: 1},
		{name: "get error maps failed to get student", adminSchoolID: nil, getErr: sentinelErr, wantErr: ErrFailedToGetStudent, wantGetCalls: 1},
		{name: "school admin cross school denied", adminSchoolID: &schoolA, getResult: &model.Student{StudentID: studentID, SchoolID: schoolB}, wantErr: ErrSchoolAccessDenied, wantGetCalls: 1},
		{name: "delete no rows maps student not found", adminSchoolID: nil, getResult: &model.Student{StudentID: studentID, SchoolID: schoolA}, deleteErr: pgx.ErrNoRows, wantErr: ErrStudentNotFound, wantGetCalls: 1, wantDeleteCalls: 1},
		{name: "delete error passthrough", adminSchoolID: nil, getResult: &model.Student{StudentID: studentID, SchoolID: schoolA}, deleteErr: sentinelErr, wantErr: sentinelErr, wantGetCalls: 1, wantDeleteCalls: 1},
		{name: "success", adminSchoolID: &schoolA, getResult: &model.Student{StudentID: studentID, SchoolID: schoolA}, wantGetCalls: 1, wantDeleteCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			studentRepo := &fakeStudentServiceStudentRepo{getByStudentResult: tc.getResult, getByStudentErr: tc.getErr, deleteErr: tc.deleteErr}
			svc := &StudentService{studentRepo: studentRepo, classRepo: &fakeStudentServiceClassRepo{}}

			err := svc.Delete(context.Background(), tc.adminSchoolID, studentID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("Delete() error = %v", err)
			}
			if studentRepo.getByStudentCalls != tc.wantGetCalls || studentRepo.deleteCalls != tc.wantDeleteCalls {
				t.Fatalf("calls get/delete = %d/%d, want %d/%d", studentRepo.getByStudentCalls, studentRepo.deleteCalls, tc.wantGetCalls, tc.wantDeleteCalls)
			}
		})
	}
}
