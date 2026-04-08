package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5"
)

type fakeTeacherServiceTeacherRepo struct {
	listCalls    int
	listSchoolID *uuid.UUID
	listLimit    int
	listOffset   int
	listResult   []model.Teacher
	listTotal    int
	listErr      error

	getCalls  int
	getArg    uuid.UUID
	getResult *model.Teacher
	getErr    error

	updateCalls   int
	updateTeacher uuid.UUID
	updateName    string
	updatePhone   string
	updateSchool  uuid.UUID
	updateErr     error

	deleteCalls   int
	deleteTeacher uuid.UUID
	deleteErr     error
}

func (f *fakeTeacherServiceTeacherRepo) List(_ context.Context, schoolID *uuid.UUID, limit, offset int) ([]model.Teacher, int, error) {
	f.listCalls++
	f.listSchoolID = schoolID
	f.listLimit = limit
	f.listOffset = offset
	return f.listResult, f.listTotal, f.listErr
}

func (f *fakeTeacherServiceTeacherRepo) GetByTeacherID(_ context.Context, teacherID uuid.UUID) (*model.Teacher, error) {
	f.getCalls++
	f.getArg = teacherID
	return f.getResult, f.getErr
}

func (f *fakeTeacherServiceTeacherRepo) Update(_ context.Context, teacherID uuid.UUID, fullName, phone string, schoolID uuid.UUID) error {
	f.updateCalls++
	f.updateTeacher = teacherID
	f.updateName = fullName
	f.updatePhone = phone
	f.updateSchool = schoolID
	return f.updateErr
}

func (f *fakeTeacherServiceTeacherRepo) Delete(_ context.Context, teacherID uuid.UUID) error {
	f.deleteCalls++
	f.deleteTeacher = teacherID
	return f.deleteErr
}

type fakeTeacherServiceTeacherClassRepo struct {
	assignCalls   int
	assignTeacher uuid.UUID
	assignClass   uuid.UUID
	assignErr     error

	listDetailsCalls int
	listDetailsClass uuid.UUID
	listDetailsRes   []model.Teacher
	listDetailsErr   error

	isAssignedCalls   int
	isAssignedTeacher uuid.UUID
	isAssignedClass   uuid.UUID
	isAssignedRes     bool
	isAssignedErr     error

	unassignCalls   int
	unassignTeacher uuid.UUID
	unassignClass   uuid.UUID
	unassignErr     error
}

func (f *fakeTeacherServiceTeacherClassRepo) Assign(_ context.Context, teacherID, classID uuid.UUID) error {
	f.assignCalls++
	f.assignTeacher = teacherID
	f.assignClass = classID
	return f.assignErr
}

func (f *fakeTeacherServiceTeacherClassRepo) ListTeacherDetailsOfClass(_ context.Context, classID uuid.UUID) ([]model.Teacher, error) {
	f.listDetailsCalls++
	f.listDetailsClass = classID
	return f.listDetailsRes, f.listDetailsErr
}

func (f *fakeTeacherServiceTeacherClassRepo) IsTeacherAssignedToClass(_ context.Context, teacherID, classID uuid.UUID) (bool, error) {
	f.isAssignedCalls++
	f.isAssignedTeacher = teacherID
	f.isAssignedClass = classID
	return f.isAssignedRes, f.isAssignedErr
}

func (f *fakeTeacherServiceTeacherClassRepo) Unassign(_ context.Context, teacherID, classID uuid.UUID) error {
	f.unassignCalls++
	f.unassignTeacher = teacherID
	f.unassignClass = classID
	return f.unassignErr
}

type fakeTeacherServiceClassRepo struct {
	getCalls  int
	getArg    uuid.UUID
	getResult *model.Class
	getErr    error
}

func (f *fakeTeacherServiceClassRepo) GetByClassID(_ context.Context, classID uuid.UUID) (*model.Class, error) {
	f.getCalls++
	f.getArg = classID
	return f.getResult, f.getErr
}

func TestTeacherServiceList(t *testing.T) {
	schoolID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name           string
		adminSchoolID  *uuid.UUID
		limit          int
		offset         int
		listErr        error
		wantErr        error
		wantLimit      int
		wantOffset     int
		wantListCalls  int
		wantResultSize int
		wantTotal      int
	}{
		{name: "default normalize", adminSchoolID: nil, limit: 0, offset: -2, wantLimit: 20, wantOffset: 0, wantListCalls: 1, wantResultSize: 1, wantTotal: 1},
		{name: "clamp max", adminSchoolID: &schoolID, limit: 999, offset: 3, wantLimit: 100, wantOffset: 3, wantListCalls: 1, wantResultSize: 1, wantTotal: 1},
		{name: "repo error", adminSchoolID: nil, limit: 30, offset: 1, listErr: sentinelErr, wantErr: sentinelErr, wantLimit: 30, wantOffset: 1, wantListCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			teacherRepo := &fakeTeacherServiceTeacherRepo{listResult: []model.Teacher{{TeacherID: uuid.New()}}, listTotal: 1, listErr: tc.listErr}
			svc := &TeacherService{teacherRepo: teacherRepo, teacherClassRepo: &fakeTeacherServiceTeacherClassRepo{}, classRepo: &fakeTeacherServiceClassRepo{}}

			items, total, err := svc.List(context.Background(), tc.adminSchoolID, tc.limit, tc.offset)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("List() error = %v", err)
			}

			if teacherRepo.listCalls != tc.wantListCalls {
				t.Fatalf("list calls = %d, want %d", teacherRepo.listCalls, tc.wantListCalls)
			}
			if teacherRepo.listLimit != tc.wantLimit || teacherRepo.listOffset != tc.wantOffset {
				t.Fatalf("limit/offset = %d/%d, want %d/%d", teacherRepo.listLimit, teacherRepo.listOffset, tc.wantLimit, tc.wantOffset)
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

func TestTeacherServiceAssign(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	teacherID := uuid.New()
	classID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		adminSchoolID   *uuid.UUID
		teacherResult   *model.Teacher
		teacherErr      error
		classResult     *model.Class
		classErr        error
		assignErr       error
		wantErr         error
		wantGetTeacher  int
		wantGetClass    int
		wantAssignCalls int
	}{
		{name: "teacher not found passthrough", adminSchoolID: nil, teacherErr: sentinelErr, wantErr: sentinelErr, wantGetTeacher: 1},
		{name: "school admin teacher cross school denied", adminSchoolID: &schoolA, teacherResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolB}, wantErr: ErrSchoolAccessDenied, wantGetTeacher: 1},
		{name: "school admin invalid class id", adminSchoolID: &schoolA, teacherResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolA}, classErr: sentinelErr, wantErr: ErrInvalidClassID, wantGetTeacher: 1, wantGetClass: 1},
		{name: "school admin class cross school denied", adminSchoolID: &schoolA, teacherResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolA}, classResult: &model.Class{ClassID: classID, SchoolID: schoolB}, wantErr: ErrSchoolAccessDenied, wantGetTeacher: 1, wantGetClass: 1},
		{name: "assign error passthrough", adminSchoolID: nil, teacherResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolA}, assignErr: sentinelErr, wantErr: sentinelErr, wantGetTeacher: 1, wantAssignCalls: 1},
		{name: "success for school admin", adminSchoolID: &schoolA, teacherResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolA}, classResult: &model.Class{ClassID: classID, SchoolID: schoolA}, wantGetTeacher: 1, wantGetClass: 1, wantAssignCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			teacherRepo := &fakeTeacherServiceTeacherRepo{getResult: tc.teacherResult, getErr: tc.teacherErr}
			teacherClassRepo := &fakeTeacherServiceTeacherClassRepo{assignErr: tc.assignErr}
			classRepo := &fakeTeacherServiceClassRepo{getResult: tc.classResult, getErr: tc.classErr}
			svc := &TeacherService{teacherRepo: teacherRepo, teacherClassRepo: teacherClassRepo, classRepo: classRepo}

			err := svc.Assign(context.Background(), tc.adminSchoolID, teacherID, classID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("Assign() error = %v", err)
			}
			if teacherRepo.getCalls != tc.wantGetTeacher || classRepo.getCalls != tc.wantGetClass || teacherClassRepo.assignCalls != tc.wantAssignCalls {
				t.Fatalf("calls teacher/class/assign = %d/%d/%d, want %d/%d/%d", teacherRepo.getCalls, classRepo.getCalls, teacherClassRepo.assignCalls, tc.wantGetTeacher, tc.wantGetClass, tc.wantAssignCalls)
			}
		})
	}
}

func TestTeacherServiceListTeachersOfClass(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	classID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name           string
		adminSchoolID  *uuid.UUID
		classResult    *model.Class
		classErr       error
		listErr        error
		wantErr        error
		wantGetClass   int
		wantListCalls  int
		wantResultSize int
	}{
		{name: "school admin invalid class id", adminSchoolID: &schoolA, classErr: sentinelErr, wantErr: ErrInvalidClassID, wantGetClass: 1},
		{name: "school admin cross school denied", adminSchoolID: &schoolA, classResult: &model.Class{ClassID: classID, SchoolID: schoolB}, wantErr: ErrSchoolAccessDenied, wantGetClass: 1},
		{name: "list error passthrough", adminSchoolID: nil, listErr: sentinelErr, wantErr: sentinelErr, wantListCalls: 1},
		{name: "success", adminSchoolID: &schoolA, classResult: &model.Class{ClassID: classID, SchoolID: schoolA}, wantGetClass: 1, wantListCalls: 1, wantResultSize: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			teacherClassRepo := &fakeTeacherServiceTeacherClassRepo{listDetailsRes: []model.Teacher{{TeacherID: uuid.New()}}, listDetailsErr: tc.listErr}
			classRepo := &fakeTeacherServiceClassRepo{getResult: tc.classResult, getErr: tc.classErr}
			svc := &TeacherService{teacherRepo: &fakeTeacherServiceTeacherRepo{}, teacherClassRepo: teacherClassRepo, classRepo: classRepo}

			items, err := svc.ListTeachersOfClass(context.Background(), tc.adminSchoolID, classID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("ListTeachersOfClass() error = %v", err)
			}
			if classRepo.getCalls != tc.wantGetClass || teacherClassRepo.listDetailsCalls != tc.wantListCalls {
				t.Fatalf("calls class/list = %d/%d, want %d/%d", classRepo.getCalls, teacherClassRepo.listDetailsCalls, tc.wantGetClass, tc.wantListCalls)
			}
			if tc.wantErr == nil && len(items) != tc.wantResultSize {
				t.Fatalf("items len = %d, want %d", len(items), tc.wantResultSize)
			}
		})
	}
}

func TestTeacherServiceUnassign(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	teacherID := uuid.New()
	classID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name             string
		adminSchoolID    *uuid.UUID
		teacherResult    *model.Teacher
		teacherErr       error
		classResult      *model.Class
		classErr         error
		isAssigned       bool
		isAssignedErr    error
		unassignErr      error
		wantErr          error
		wantGetTeacher   int
		wantGetClass     int
		wantIsAssigned   int
		wantUnassignCall int
	}{
		{name: "school admin teacher lookup error", adminSchoolID: &schoolA, teacherErr: sentinelErr, wantErr: sentinelErr, wantGetTeacher: 1},
		{name: "school admin teacher cross school denied", adminSchoolID: &schoolA, teacherResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolB}, wantErr: ErrSchoolAccessDenied, wantGetTeacher: 1},
		{name: "school admin invalid class id", adminSchoolID: &schoolA, teacherResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolA}, classErr: sentinelErr, wantErr: ErrInvalidClassID, wantGetTeacher: 1, wantGetClass: 1},
		{name: "school admin class cross school denied", adminSchoolID: &schoolA, teacherResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolA}, classResult: &model.Class{ClassID: classID, SchoolID: schoolB}, wantErr: ErrSchoolAccessDenied, wantGetTeacher: 1, wantGetClass: 1},
		{name: "isAssigned error", adminSchoolID: nil, isAssignedErr: sentinelErr, wantErr: sentinelErr, wantIsAssigned: 1},
		{name: "not assigned", adminSchoolID: nil, isAssigned: false, wantErr: ErrTeacherNotAssigned, wantIsAssigned: 1},
		{name: "unassign error", adminSchoolID: nil, isAssigned: true, unassignErr: sentinelErr, wantErr: sentinelErr, wantIsAssigned: 1, wantUnassignCall: 1},
		{name: "success", adminSchoolID: &schoolA, teacherResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolA}, classResult: &model.Class{ClassID: classID, SchoolID: schoolA}, isAssigned: true, wantGetTeacher: 1, wantGetClass: 1, wantIsAssigned: 1, wantUnassignCall: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			teacherRepo := &fakeTeacherServiceTeacherRepo{getResult: tc.teacherResult, getErr: tc.teacherErr}
			teacherClassRepo := &fakeTeacherServiceTeacherClassRepo{isAssignedRes: tc.isAssigned, isAssignedErr: tc.isAssignedErr, unassignErr: tc.unassignErr}
			classRepo := &fakeTeacherServiceClassRepo{getResult: tc.classResult, getErr: tc.classErr}
			svc := &TeacherService{teacherRepo: teacherRepo, teacherClassRepo: teacherClassRepo, classRepo: classRepo}

			err := svc.Unassign(context.Background(), tc.adminSchoolID, teacherID, classID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("Unassign() error = %v", err)
			}
			if teacherRepo.getCalls != tc.wantGetTeacher || classRepo.getCalls != tc.wantGetClass || teacherClassRepo.isAssignedCalls != tc.wantIsAssigned || teacherClassRepo.unassignCalls != tc.wantUnassignCall {
				t.Fatalf("calls teacher/class/isAssigned/unassign = %d/%d/%d/%d, want %d/%d/%d/%d", teacherRepo.getCalls, classRepo.getCalls, teacherClassRepo.isAssignedCalls, teacherClassRepo.unassignCalls, tc.wantGetTeacher, tc.wantGetClass, tc.wantIsAssigned, tc.wantUnassignCall)
			}
		})
	}
}

func TestTeacherServiceGetByTeacherID(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	teacherID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name          string
		adminSchoolID *uuid.UUID
		getResult     *model.Teacher
		getErr        error
		wantErr       error
	}{
		{name: "teacher not found mapped", adminSchoolID: nil, getErr: pgx.ErrNoRows, wantErr: ErrTeacherNotFound},
		{name: "repo error passthrough", adminSchoolID: nil, getErr: sentinelErr, wantErr: sentinelErr},
		{name: "school admin cross school denied", adminSchoolID: &schoolA, getResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolB}, wantErr: ErrSchoolAccessDenied},
		{name: "success", adminSchoolID: &schoolA, getResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolA}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			teacherRepo := &fakeTeacherServiceTeacherRepo{getResult: tc.getResult, getErr: tc.getErr}
			svc := &TeacherService{teacherRepo: teacherRepo, teacherClassRepo: &fakeTeacherServiceTeacherClassRepo{}, classRepo: &fakeTeacherServiceClassRepo{}}

			got, err := svc.GetByTeacherID(context.Background(), tc.adminSchoolID, teacherID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("GetByTeacherID() error = %v", err)
			}
			if tc.wantErr == nil && (got == nil || got.TeacherID != teacherID) {
				t.Fatalf("unexpected teacher = %#v", got)
			}
		})
	}
}

func TestTeacherServiceUpdate(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	teacherID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		adminSchoolID   *uuid.UUID
		schoolID        uuid.UUID
		getResult       *model.Teacher
		getErr          error
		updateErr       error
		wantErr         error
		wantGetCalls    int
		wantUpdateCalls int
	}{
		{name: "teacher not found", adminSchoolID: nil, schoolID: schoolA, getErr: pgx.ErrNoRows, wantErr: ErrTeacherNotFound, wantGetCalls: 1},
		{name: "repo get error", adminSchoolID: nil, schoolID: schoolA, getErr: sentinelErr, wantErr: sentinelErr, wantGetCalls: 1},
		{name: "school admin cannot update foreign teacher", adminSchoolID: &schoolA, schoolID: schoolA, getResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolB}, wantErr: ErrSchoolAccessDenied, wantGetCalls: 1},
		{name: "school admin cannot move teacher to another school", adminSchoolID: &schoolA, schoolID: schoolB, getResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolA}, wantErr: ErrSchoolAccessDenied, wantGetCalls: 1},
		{name: "update error passthrough", adminSchoolID: nil, schoolID: schoolA, getResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolA}, updateErr: sentinelErr, wantErr: sentinelErr, wantGetCalls: 1, wantUpdateCalls: 1},
		{name: "success", adminSchoolID: &schoolA, schoolID: schoolA, getResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolA}, wantGetCalls: 1, wantUpdateCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			teacherRepo := &fakeTeacherServiceTeacherRepo{getResult: tc.getResult, getErr: tc.getErr, updateErr: tc.updateErr}
			svc := &TeacherService{teacherRepo: teacherRepo, teacherClassRepo: &fakeTeacherServiceTeacherClassRepo{}, classRepo: &fakeTeacherServiceClassRepo{}}

			err := svc.Update(context.Background(), tc.adminSchoolID, teacherID, "Teacher A", "0909", tc.schoolID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("Update() error = %v", err)
			}

			if teacherRepo.getCalls != tc.wantGetCalls || teacherRepo.updateCalls != tc.wantUpdateCalls {
				t.Fatalf("calls get/update = %d/%d, want %d/%d", teacherRepo.getCalls, teacherRepo.updateCalls, tc.wantGetCalls, tc.wantUpdateCalls)
			}
		})
	}
}

func TestTeacherServiceDelete(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	teacherID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		adminSchoolID   *uuid.UUID
		getResult       *model.Teacher
		getErr          error
		deleteErr       error
		wantErr         error
		wantGetCalls    int
		wantDeleteCalls int
	}{
		{name: "get no rows maps teacher not found", adminSchoolID: nil, getErr: pgx.ErrNoRows, wantErr: ErrTeacherNotFound, wantGetCalls: 1},
		{name: "get other error still maps teacher not found", adminSchoolID: nil, getErr: sentinelErr, wantErr: ErrTeacherNotFound, wantGetCalls: 1},
		{name: "school admin cross school denied", adminSchoolID: &schoolA, getResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolB}, wantErr: ErrSchoolAccessDenied, wantGetCalls: 1},
		{name: "delete no rows maps teacher not found", adminSchoolID: nil, getResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolA}, deleteErr: pgx.ErrNoRows, wantErr: ErrTeacherNotFound, wantGetCalls: 1, wantDeleteCalls: 1},
		{name: "delete error passthrough", adminSchoolID: nil, getResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolA}, deleteErr: sentinelErr, wantErr: sentinelErr, wantGetCalls: 1, wantDeleteCalls: 1},
		{name: "success", adminSchoolID: &schoolA, getResult: &model.Teacher{TeacherID: teacherID, SchoolID: schoolA}, wantGetCalls: 1, wantDeleteCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			teacherRepo := &fakeTeacherServiceTeacherRepo{getResult: tc.getResult, getErr: tc.getErr, deleteErr: tc.deleteErr}
			svc := &TeacherService{teacherRepo: teacherRepo, teacherClassRepo: &fakeTeacherServiceTeacherClassRepo{}, classRepo: &fakeTeacherServiceClassRepo{}}

			err := svc.Delete(context.Background(), tc.adminSchoolID, teacherID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("Delete() error = %v", err)
			}

			if teacherRepo.getCalls != tc.wantGetCalls || teacherRepo.deleteCalls != tc.wantDeleteCalls {
				t.Fatalf("calls get/delete = %d/%d, want %d/%d", teacherRepo.getCalls, teacherRepo.deleteCalls, tc.wantGetCalls, tc.wantDeleteCalls)
			}
		})
	}
}
