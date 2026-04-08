package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5"
)

type fakeClassRepo struct {
	createCalls      int
	createSchoolID   uuid.UUID
	createName       string
	createSchoolYear string
	createResultID   uuid.UUID
	createErr        error

	listCalls    int
	listSchoolID uuid.UUID
	listLimit    int
	listOffset   int
	listResult   []model.Class
	listTotal    int
	listErr      error

	getByClassIDCalls  int
	getByClassIDArg    uuid.UUID
	getByClassIDResult *model.Class
	getByClassIDErr    error

	updateCalls      int
	updateClassID    uuid.UUID
	updateName       string
	updateSchoolYear string
	updateErr        error

	deleteCalls   int
	deleteClassID uuid.UUID
	deleteErr     error
}

func (f *fakeClassRepo) Create(_ context.Context, schoolID uuid.UUID, name, schoolYear string) (uuid.UUID, error) {
	f.createCalls++
	f.createSchoolID = schoolID
	f.createName = name
	f.createSchoolYear = schoolYear
	return f.createResultID, f.createErr
}

func (f *fakeClassRepo) List(_ context.Context, schoolID uuid.UUID, limit, offset int) ([]model.Class, int, error) {
	f.listCalls++
	f.listSchoolID = schoolID
	f.listLimit = limit
	f.listOffset = offset
	return f.listResult, f.listTotal, f.listErr
}

func (f *fakeClassRepo) GetByClassID(_ context.Context, classID uuid.UUID) (*model.Class, error) {
	f.getByClassIDCalls++
	f.getByClassIDArg = classID
	return f.getByClassIDResult, f.getByClassIDErr
}

func (f *fakeClassRepo) Update(_ context.Context, classID uuid.UUID, name, schoolYear string) error {
	f.updateCalls++
	f.updateClassID = classID
	f.updateName = name
	f.updateSchoolYear = schoolYear
	return f.updateErr
}

func (f *fakeClassRepo) Delete(_ context.Context, classID uuid.UUID) error {
	f.deleteCalls++
	f.deleteClassID = classID
	return f.deleteErr
}

func TestClassServiceCreate(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	createdID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		adminSchoolID   *uuid.UUID
		schoolID        uuid.UUID
		createResultID  uuid.UUID
		createErr       error
		wantErr         error
		wantCreateCalls int
	}{
		{
			name:          "school admin cannot create class for another school",
			adminSchoolID: &schoolA,
			schoolID:      schoolB,
			wantErr:       ErrSchoolAccessDenied,
		},
		{
			name:            "create repo error",
			adminSchoolID:   nil,
			schoolID:        schoolA,
			createErr:       sentinelErr,
			wantErr:         sentinelErr,
			wantCreateCalls: 1,
		},
		{
			name:            "success",
			adminSchoolID:   &schoolA,
			schoolID:        schoolA,
			createResultID:  createdID,
			wantCreateCalls: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeClassRepo{createResultID: tc.createResultID, createErr: tc.createErr}
			svc := &ClassService{classRepo: repo}

			got, err := svc.Create(context.Background(), tc.adminSchoolID, tc.schoolID, "Class A", "2026-2027")
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("Create() error = %v", err)
			}

			if repo.createCalls != tc.wantCreateCalls {
				t.Fatalf("create calls = %d, want %d", repo.createCalls, tc.wantCreateCalls)
			}

			if tc.wantErr == nil {
				if got == nil {
					t.Fatal("expected class on success")
				}
				if got.ClassID != tc.createResultID || got.SchoolID != tc.schoolID || got.Name != "Class A" || got.SchoolYear != "2026-2027" {
					t.Fatalf("unexpected class = %#v", got)
				}
			}
		})
	}
}

func TestClassServiceListBySchool(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name           string
		adminSchoolID  *uuid.UUID
		schoolID       uuid.UUID
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
		{
			name:          "school admin cannot list another school",
			adminSchoolID: &schoolA,
			schoolID:      schoolB,
			wantErr:       ErrSchoolAccessDenied,
		},
		{
			name:           "normalize default limit and offset",
			adminSchoolID:  &schoolA,
			schoolID:       schoolA,
			limit:          0,
			offset:         -3,
			wantLimit:      20,
			wantOffset:     0,
			wantListCalls:  1,
			wantResultSize: 1,
			wantTotal:      1,
		},
		{
			name:           "clamp max limit",
			adminSchoolID:  nil,
			schoolID:       schoolA,
			limit:          999,
			offset:         4,
			wantLimit:      100,
			wantOffset:     4,
			wantListCalls:  1,
			wantResultSize: 1,
			wantTotal:      1,
		},
		{
			name:           "repo error",
			adminSchoolID:  nil,
			schoolID:       schoolA,
			limit:          30,
			offset:         2,
			listErr:        sentinelErr,
			wantErr:        sentinelErr,
			wantLimit:      30,
			wantOffset:     2,
			wantListCalls:  1,
			wantResultSize: 0,
			wantTotal:      0,
		},
		{
			name:           "success",
			adminSchoolID:  nil,
			schoolID:       schoolA,
			limit:          25,
			offset:         1,
			wantLimit:      25,
			wantOffset:     1,
			wantListCalls:  1,
			wantResultSize: 1,
			wantTotal:      1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeClassRepo{
				listResult: []model.Class{{ClassID: uuid.New(), SchoolID: schoolA, Name: "Class A", SchoolYear: "2026-2027"}},
				listTotal:  1,
				listErr:    tc.listErr,
			}
			svc := &ClassService{classRepo: repo}

			items, total, err := svc.ListBySchool(context.Background(), tc.adminSchoolID, tc.schoolID, tc.limit, tc.offset)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("ListBySchool() error = %v", err)
			}

			if repo.listCalls != tc.wantListCalls {
				t.Fatalf("list calls = %d, want %d", repo.listCalls, tc.wantListCalls)
			}
			if tc.wantListCalls > 0 {
				if repo.listLimit != tc.wantLimit || repo.listOffset != tc.wantOffset {
					t.Fatalf("limit/offset = %d/%d, want %d/%d", repo.listLimit, repo.listOffset, tc.wantLimit, tc.wantOffset)
				}
			}

			if tc.wantErr != nil {
				return
			}

			if len(items) != tc.wantResultSize {
				t.Fatalf("items len = %d, want %d", len(items), tc.wantResultSize)
			}
			if total != tc.wantTotal {
				t.Fatalf("total = %d, want %d", total, tc.wantTotal)
			}
		})
	}
}

func TestClassServiceUpdate(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	classID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		adminSchoolID   *uuid.UUID
		getResult       *model.Class
		getErr          error
		updateErr       error
		wantErr         error
		wantGetCalls    int
		wantUpdateCalls int
	}{
		{
			name:          "class not found from get",
			adminSchoolID: &schoolA,
			getErr:        pgx.ErrNoRows,
			wantErr:       ErrClassNotFound,
			wantGetCalls:  1,
		},
		{
			name:          "get class error",
			adminSchoolID: nil,
			getErr:        sentinelErr,
			wantErr:       sentinelErr,
			wantGetCalls:  1,
		},
		{
			name:          "school admin cannot update foreign school",
			adminSchoolID: &schoolA,
			getResult:     &model.Class{ClassID: classID, SchoolID: schoolB},
			wantErr:       ErrSchoolAccessDenied,
			wantGetCalls:  1,
		},
		{
			name:            "update reports not found",
			adminSchoolID:   &schoolA,
			getResult:       &model.Class{ClassID: classID, SchoolID: schoolA},
			updateErr:       pgx.ErrNoRows,
			wantErr:         ErrClassNotFound,
			wantGetCalls:    1,
			wantUpdateCalls: 1,
		},
		{
			name:            "update error",
			adminSchoolID:   nil,
			getResult:       &model.Class{ClassID: classID, SchoolID: schoolA},
			updateErr:       sentinelErr,
			wantErr:         sentinelErr,
			wantGetCalls:    1,
			wantUpdateCalls: 1,
		},
		{
			name:            "success",
			adminSchoolID:   &schoolA,
			getResult:       &model.Class{ClassID: classID, SchoolID: schoolA},
			wantGetCalls:    1,
			wantUpdateCalls: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeClassRepo{getByClassIDResult: tc.getResult, getByClassIDErr: tc.getErr, updateErr: tc.updateErr}
			svc := &ClassService{classRepo: repo}

			err := svc.Update(context.Background(), tc.adminSchoolID, classID, "Class B", "2027-2028")
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("Update() error = %v", err)
			}

			if repo.getByClassIDCalls != tc.wantGetCalls {
				t.Fatalf("getByClassID calls = %d, want %d", repo.getByClassIDCalls, tc.wantGetCalls)
			}
			if repo.updateCalls != tc.wantUpdateCalls {
				t.Fatalf("update calls = %d, want %d", repo.updateCalls, tc.wantUpdateCalls)
			}
			if tc.wantUpdateCalls > 0 {
				if repo.updateClassID != classID || repo.updateName != "Class B" || repo.updateSchoolYear != "2027-2028" {
					t.Fatalf("unexpected update args: classID=%v name=%q schoolYear=%q", repo.updateClassID, repo.updateName, repo.updateSchoolYear)
				}
			}
		})
	}
}

func TestClassServiceDelete(t *testing.T) {
	schoolA := uuid.New()
	schoolB := uuid.New()
	classID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		adminSchoolID   *uuid.UUID
		getResult       *model.Class
		getErr          error
		deleteErr       error
		wantErr         error
		wantGetCalls    int
		wantDeleteCalls int
	}{
		{
			name:          "class not found from get",
			adminSchoolID: &schoolA,
			getErr:        pgx.ErrNoRows,
			wantErr:       ErrClassNotFound,
			wantGetCalls:  1,
		},
		{
			name:          "get class error",
			adminSchoolID: nil,
			getErr:        sentinelErr,
			wantErr:       sentinelErr,
			wantGetCalls:  1,
		},
		{
			name:          "school admin cannot delete foreign school",
			adminSchoolID: &schoolA,
			getResult:     &model.Class{ClassID: classID, SchoolID: schoolB},
			wantErr:       ErrSchoolAccessDenied,
			wantGetCalls:  1,
		},
		{
			name:            "delete reports not found",
			adminSchoolID:   &schoolA,
			getResult:       &model.Class{ClassID: classID, SchoolID: schoolA},
			deleteErr:       pgx.ErrNoRows,
			wantErr:         ErrClassNotFound,
			wantGetCalls:    1,
			wantDeleteCalls: 1,
		},
		{
			name:            "delete error",
			adminSchoolID:   nil,
			getResult:       &model.Class{ClassID: classID, SchoolID: schoolA},
			deleteErr:       sentinelErr,
			wantErr:         sentinelErr,
			wantGetCalls:    1,
			wantDeleteCalls: 1,
		},
		{
			name:            "success",
			adminSchoolID:   &schoolA,
			getResult:       &model.Class{ClassID: classID, SchoolID: schoolA},
			wantGetCalls:    1,
			wantDeleteCalls: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeClassRepo{getByClassIDResult: tc.getResult, getByClassIDErr: tc.getErr, deleteErr: tc.deleteErr}
			svc := &ClassService{classRepo: repo}

			err := svc.Delete(context.Background(), tc.adminSchoolID, classID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("Delete() error = %v", err)
			}

			if repo.getByClassIDCalls != tc.wantGetCalls {
				t.Fatalf("getByClassID calls = %d, want %d", repo.getByClassIDCalls, tc.wantGetCalls)
			}
			if repo.deleteCalls != tc.wantDeleteCalls {
				t.Fatalf("delete calls = %d, want %d", repo.deleteCalls, tc.wantDeleteCalls)
			}
			if tc.wantDeleteCalls > 0 && repo.deleteClassID != classID {
				t.Fatalf("delete classID = %v, want %v", repo.deleteClassID, classID)
			}
		})
	}
}
