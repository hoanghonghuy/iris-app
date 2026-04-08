package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

type fakeParentRepo struct {
	listSchoolID *uuid.UUID
	listLimit    int
	listOffset   int
	listResult   []model.Parent
	listTotal    int
	listErr      error

	createID  uuid.UUID
	createErr error

	getByParentIDCalls  int
	getByParentIDArg    uuid.UUID
	getByParentIDResult *model.Parent
	getByParentIDErr    error
}

func (f *fakeParentRepo) List(_ context.Context, schoolID *uuid.UUID, limit, offset int) ([]model.Parent, int, error) {
	f.listSchoolID = schoolID
	f.listLimit = limit
	f.listOffset = offset
	return f.listResult, f.listTotal, f.listErr
}

func (f *fakeParentRepo) Create(_ context.Context, userID, schoolID uuid.UUID, fullName, phone string) (uuid.UUID, error) {
	return f.createID, f.createErr
}

func (f *fakeParentRepo) GetByParentID(_ context.Context, parentID uuid.UUID) (*model.Parent, error) {
	f.getByParentIDCalls++
	f.getByParentIDArg = parentID
	return f.getByParentIDResult, f.getByParentIDErr
}

type fakeStudentParentRepo struct {
	assignCalls int
	assignArgs  struct {
		studentID    uuid.UUID
		parentID     uuid.UUID
		relationship string
	}
	assignErr error

	unassignCalls int
	unassignArgs  struct {
		studentID uuid.UUID
		parentID  uuid.UUID
	}
	unassignErr error
}

func (f *fakeStudentParentRepo) Assign(_ context.Context, studentID, parentID uuid.UUID, relationship string) error {
	f.assignCalls++
	f.assignArgs.studentID = studentID
	f.assignArgs.parentID = parentID
	f.assignArgs.relationship = relationship
	return f.assignErr
}

func (f *fakeStudentParentRepo) Unassign(_ context.Context, studentID, parentID uuid.UUID) error {
	f.unassignCalls++
	f.unassignArgs.studentID = studentID
	f.unassignArgs.parentID = parentID
	return f.unassignErr
}

type fakeStudentRepo struct {
	getByStudentIDCalls  int
	getByStudentIDArg    uuid.UUID
	getByStudentIDResult *model.Student
	getByStudentIDErr    error
}

func (f *fakeStudentRepo) GetByStudentID(_ context.Context, studentID uuid.UUID) (*model.Student, error) {
	f.getByStudentIDCalls++
	f.getByStudentIDArg = studentID
	return f.getByStudentIDResult, f.getByStudentIDErr
}

func TestParentServiceListNormalizesPagination(t *testing.T) {
	schoolID := uuid.New()
	parentRepo := &fakeParentRepo{
		listResult: []model.Parent{{ParentID: uuid.New()}},
		listTotal:  1,
	}
	svc := &ParentService{parentRepo: parentRepo}

	tests := []struct {
		name       string
		limit      int
		offset     int
		wantLimit  int
		wantOffset int
	}{
		{name: "default limit and clamp negative offset", limit: 0, offset: -4, wantLimit: 20, wantOffset: 0},
		{name: "clamp max limit", limit: 999, offset: 3, wantLimit: 100, wantOffset: 3},
		{name: "keep valid values", limit: 35, offset: 7, wantLimit: 35, wantOffset: 7},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			items, total, err := svc.List(context.Background(), &schoolID, tc.limit, tc.offset)
			if err != nil {
				t.Fatalf("List() error = %v", err)
			}
			if len(items) != 1 || total != 1 {
				t.Fatalf("List() returned unexpected result: len=%d total=%d", len(items), total)
			}
			if parentRepo.listLimit != tc.wantLimit {
				t.Fatalf("limit = %d, want %d", parentRepo.listLimit, tc.wantLimit)
			}
			if parentRepo.listOffset != tc.wantOffset {
				t.Fatalf("offset = %d, want %d", parentRepo.listOffset, tc.wantOffset)
			}
			if parentRepo.listSchoolID == nil || *parentRepo.listSchoolID != schoolID {
				t.Fatalf("schoolID was not forwarded correctly")
			}
		})
	}
}

func TestParentServiceGetByParentID(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	parentID := uuid.New()
	repoErr := errors.New("boom")

	tests := []struct {
		name          string
		adminSchoolID *uuid.UUID
		parentResult  *model.Parent
		parentErr     error
		wantErr       error
	}{
		{
			name:      "repo error",
			parentErr: repoErr,
			wantErr:   repoErr,
		},
		{
			name:         "super admin bypasses school check",
			parentResult: &model.Parent{ParentID: parentID, SchoolID: schoolA},
		},
		{
			name:          "school admin blocked by different school",
			adminSchoolID: uuidPtr(schoolB),
			parentResult:  &model.Parent{ParentID: parentID, SchoolID: schoolA},
			wantErr:       ErrSchoolAccessDenied,
		},
		{
			name:          "school admin can access same school",
			adminSchoolID: uuidPtr(schoolA),
			parentResult:  &model.Parent{ParentID: parentID, SchoolID: schoolA},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parentRepo := &fakeParentRepo{
				getByParentIDResult: tc.parentResult,
				getByParentIDErr:    tc.parentErr,
			}
			svc := &ParentService{parentRepo: parentRepo}

			got, err := svc.GetByParentID(context.Background(), tc.adminSchoolID, parentID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetByParentID() error = %v", err)
			}
			if got == nil || got.ParentID != parentID {
				t.Fatalf("unexpected parent returned: %#v", got)
			}
		})
	}
}

func TestParentServiceAssignStudent(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	parentID := uuid.New()
	studentID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name               string
		adminSchoolID      *uuid.UUID
		parentID           uuid.UUID
		studentID          uuid.UUID
		relationship       string
		parentResult       *model.Parent
		parentErr          error
		studentResult      *model.Student
		studentErr         error
		assignErr          error
		wantErr            error
		wantAssignCalls    int
		wantParentLookups  int
		wantStudentLookups int
	}{
		{
			name:         "invalid parent id",
			parentID:     uuid.Nil,
			studentID:    studentID,
			relationship: "mother",
			wantErr:      ErrInvalidUserID,
		},
		{
			name:         "invalid student id",
			parentID:     parentID,
			studentID:    uuid.Nil,
			relationship: "mother",
			wantErr:      ErrInvalidUserID,
		},
		{
			name:            "super admin assigns without school lookups",
			adminSchoolID:   nil,
			parentID:        parentID,
			studentID:       studentID,
			relationship:    "mother",
			wantAssignCalls: 1,
		},
		{
			name:              "parent lookup error",
			adminSchoolID:     uuidPtr(schoolA),
			parentID:          parentID,
			studentID:         studentID,
			relationship:      "mother",
			parentErr:         sentinelErr,
			wantErr:           sentinelErr,
			wantParentLookups: 1,
		},
		{
			name:              "parent school mismatch",
			adminSchoolID:     uuidPtr(schoolA),
			parentID:          parentID,
			studentID:         studentID,
			relationship:      "mother",
			parentResult:      &model.Parent{ParentID: parentID, SchoolID: schoolB},
			wantErr:           ErrSchoolAccessDenied,
			wantParentLookups: 1,
		},
		{
			name:               "student lookup error",
			adminSchoolID:      uuidPtr(schoolA),
			parentID:           parentID,
			studentID:          studentID,
			relationship:       "mother",
			parentResult:       &model.Parent{ParentID: parentID, SchoolID: schoolA},
			studentErr:         sentinelErr,
			wantErr:            sentinelErr,
			wantParentLookups:  1,
			wantStudentLookups: 1,
		},
		{
			name:               "student school mismatch",
			adminSchoolID:      uuidPtr(schoolA),
			parentID:           parentID,
			studentID:          studentID,
			relationship:       "mother",
			parentResult:       &model.Parent{ParentID: parentID, SchoolID: schoolA},
			studentResult:      &model.Student{StudentID: studentID, SchoolID: schoolB},
			wantErr:            ErrSchoolAccessDenied,
			wantParentLookups:  1,
			wantStudentLookups: 1,
		},
		{
			name:               "assign error",
			adminSchoolID:      uuidPtr(schoolA),
			parentID:           parentID,
			studentID:          studentID,
			relationship:       "mother",
			parentResult:       &model.Parent{ParentID: parentID, SchoolID: schoolA},
			studentResult:      &model.Student{StudentID: studentID, SchoolID: schoolA},
			assignErr:          sentinelErr,
			wantErr:            sentinelErr,
			wantAssignCalls:    1,
			wantParentLookups:  1,
			wantStudentLookups: 1,
		},
		{
			name:               "success",
			adminSchoolID:      uuidPtr(schoolA),
			parentID:           parentID,
			studentID:          studentID,
			relationship:       "mother",
			parentResult:       &model.Parent{ParentID: parentID, SchoolID: schoolA},
			studentResult:      &model.Student{StudentID: studentID, SchoolID: schoolA},
			wantAssignCalls:    1,
			wantParentLookups:  1,
			wantStudentLookups: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parentRepo := &fakeParentRepo{getByParentIDResult: tc.parentResult, getByParentIDErr: tc.parentErr}
			studentParentRepo := &fakeStudentParentRepo{assignErr: tc.assignErr}
			studentRepo := &fakeStudentRepo{getByStudentIDResult: tc.studentResult, getByStudentIDErr: tc.studentErr}
			svc := &ParentService{parentRepo: parentRepo, studentParentRepo: studentParentRepo, studentRepo: studentRepo}

			err := svc.AssignStudent(context.Background(), tc.adminSchoolID, tc.parentID, tc.studentID, tc.relationship)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("AssignStudent() error = %v", err)
			}

			if parentRepo.getByParentIDCalls != tc.wantParentLookups {
				t.Fatalf("parent lookups = %d, want %d", parentRepo.getByParentIDCalls, tc.wantParentLookups)
			}
			if studentRepo.getByStudentIDCalls != tc.wantStudentLookups {
				t.Fatalf("student lookups = %d, want %d", studentRepo.getByStudentIDCalls, tc.wantStudentLookups)
			}
			if studentParentRepo.assignCalls != tc.wantAssignCalls {
				t.Fatalf("assign calls = %d, want %d", studentParentRepo.assignCalls, tc.wantAssignCalls)
			}
			if tc.wantAssignCalls > 0 {
				if studentParentRepo.assignArgs.parentID != tc.parentID || studentParentRepo.assignArgs.studentID != tc.studentID || studentParentRepo.assignArgs.relationship != tc.relationship {
					t.Fatalf("assign args mismatch: %+v", studentParentRepo.assignArgs)
				}
			}
		})
	}
}

func TestParentServiceUnassignStudent(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	parentID := uuid.New()
	studentID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name               string
		adminSchoolID      *uuid.UUID
		parentID           uuid.UUID
		studentID          uuid.UUID
		parentResult       *model.Parent
		parentErr          error
		studentResult      *model.Student
		studentErr         error
		unassignErr        error
		wantErr            error
		wantUnassignCalls  int
		wantParentLookups  int
		wantStudentLookups int
	}{
		{
			name:      "invalid parent id",
			parentID:  uuid.Nil,
			studentID: studentID,
			wantErr:   ErrInvalidUserID,
		},
		{
			name:      "invalid student id",
			parentID:  parentID,
			studentID: uuid.Nil,
			wantErr:   ErrInvalidUserID,
		},
		{
			name:              "super admin unassigns without school lookups",
			parentID:          parentID,
			studentID:         studentID,
			wantUnassignCalls: 1,
		},
		{
			name:              "parent school mismatch",
			adminSchoolID:     uuidPtr(schoolA),
			parentID:          parentID,
			studentID:         studentID,
			parentResult:      &model.Parent{ParentID: parentID, SchoolID: schoolB},
			wantErr:           ErrSchoolAccessDenied,
			wantParentLookups: 1,
		},
		{
			name:               "student school mismatch",
			adminSchoolID:      uuidPtr(schoolA),
			parentID:           parentID,
			studentID:          studentID,
			parentResult:       &model.Parent{ParentID: parentID, SchoolID: schoolA},
			studentResult:      &model.Student{StudentID: studentID, SchoolID: schoolB},
			wantErr:            ErrSchoolAccessDenied,
			wantParentLookups:  1,
			wantStudentLookups: 1,
		},
		{
			name:               "repo unassign error",
			adminSchoolID:      uuidPtr(schoolA),
			parentID:           parentID,
			studentID:          studentID,
			parentResult:       &model.Parent{ParentID: parentID, SchoolID: schoolA},
			studentResult:      &model.Student{StudentID: studentID, SchoolID: schoolA},
			unassignErr:        sentinelErr,
			wantErr:            sentinelErr,
			wantUnassignCalls:  1,
			wantParentLookups:  1,
			wantStudentLookups: 1,
		},
		{
			name:               "success",
			adminSchoolID:      uuidPtr(schoolA),
			parentID:           parentID,
			studentID:          studentID,
			parentResult:       &model.Parent{ParentID: parentID, SchoolID: schoolA},
			studentResult:      &model.Student{StudentID: studentID, SchoolID: schoolA},
			wantUnassignCalls:  1,
			wantParentLookups:  1,
			wantStudentLookups: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parentRepo := &fakeParentRepo{getByParentIDResult: tc.parentResult, getByParentIDErr: tc.parentErr}
			studentParentRepo := &fakeStudentParentRepo{unassignErr: tc.unassignErr}
			studentRepo := &fakeStudentRepo{getByStudentIDResult: tc.studentResult, getByStudentIDErr: tc.studentErr}
			svc := &ParentService{parentRepo: parentRepo, studentParentRepo: studentParentRepo, studentRepo: studentRepo}

			err := svc.UnassignStudent(context.Background(), tc.adminSchoolID, tc.parentID, tc.studentID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("UnassignStudent() error = %v", err)
			}

			if parentRepo.getByParentIDCalls != tc.wantParentLookups {
				t.Fatalf("parent lookups = %d, want %d", parentRepo.getByParentIDCalls, tc.wantParentLookups)
			}
			if studentRepo.getByStudentIDCalls != tc.wantStudentLookups {
				t.Fatalf("student lookups = %d, want %d", studentRepo.getByStudentIDCalls, tc.wantStudentLookups)
			}
			if studentParentRepo.unassignCalls != tc.wantUnassignCalls {
				t.Fatalf("unassign calls = %d, want %d", studentParentRepo.unassignCalls, tc.wantUnassignCalls)
			}
			if tc.wantUnassignCalls > 0 {
				if studentParentRepo.unassignArgs.parentID != tc.parentID || studentParentRepo.unassignArgs.studentID != tc.studentID {
					t.Fatalf("unassign args mismatch: %+v", studentParentRepo.unassignArgs)
				}
			}
		})
	}
}

func uuidPtr(v uuid.UUID) *uuid.UUID {
	return &v
}
