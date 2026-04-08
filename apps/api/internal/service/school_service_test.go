package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5"
)

type fakeSchoolRepo struct {
	createCalls   int
	createName    string
	createAddress string
	createResult  uuid.UUID
	createErr     error

	listCalls  int
	listLimit  int
	listOffset int
	listResult []model.School
	listTotal  int
	listErr    error

	getByIDCalls  int
	getByIDArg    uuid.UUID
	getByIDResult *model.School
	getByIDErr    error

	updateCalls   int
	updateSchool  uuid.UUID
	updateName    string
	updateAddress string
	updateErr     error

	deleteCalls  int
	deleteSchool uuid.UUID
	deleteErr    error
}

func (f *fakeSchoolRepo) Create(_ context.Context, name, address string) (uuid.UUID, error) {
	f.createCalls++
	f.createName = name
	f.createAddress = address
	return f.createResult, f.createErr
}

func (f *fakeSchoolRepo) List(_ context.Context, limit, offset int) ([]model.School, int, error) {
	f.listCalls++
	f.listLimit = limit
	f.listOffset = offset
	return f.listResult, f.listTotal, f.listErr
}

func (f *fakeSchoolRepo) GetByID(_ context.Context, schoolID uuid.UUID) (*model.School, error) {
	f.getByIDCalls++
	f.getByIDArg = schoolID
	return f.getByIDResult, f.getByIDErr
}

func (f *fakeSchoolRepo) Update(_ context.Context, schoolID uuid.UUID, name, address string) error {
	f.updateCalls++
	f.updateSchool = schoolID
	f.updateName = name
	f.updateAddress = address
	return f.updateErr
}

func (f *fakeSchoolRepo) Delete(_ context.Context, schoolID uuid.UUID) error {
	f.deleteCalls++
	f.deleteSchool = schoolID
	return f.deleteErr
}

func TestSchoolServiceCreate(t *testing.T) {
	createdID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		createErr       error
		wantErr         error
		wantCreateCalls int
	}{
		{
			name:            "repo error",
			createErr:       sentinelErr,
			wantErr:         sentinelErr,
			wantCreateCalls: 1,
		},
		{
			name:            "success",
			wantCreateCalls: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeSchoolRepo{createResult: createdID, createErr: tc.createErr}
			svc := &SchoolService{schoolRepo: repo}

			got, err := svc.Create(context.Background(), "Iris School", "District 1")
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
					t.Fatal("expected school on success")
				}
				if got.SchoolID != createdID || got.Name != "Iris School" || got.Address != "District 1" {
					t.Fatalf("unexpected school = %#v", got)
				}
			}
		})
	}
}

func TestSchoolServiceList(t *testing.T) {
	adminSchoolID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name           string
		adminSchoolID  *uuid.UUID
		limit          int
		offset         int
		listErr        error
		getByIDErr     error
		wantErr        error
		wantListCalls  int
		wantGetByID    int
		wantLimit      int
		wantOffset     int
		wantResultSize int
		wantTotal      int
	}{
		{
			name:           "normalize default limit and offset for super admin",
			adminSchoolID:  nil,
			limit:          0,
			offset:         -5,
			wantListCalls:  1,
			wantGetByID:    0,
			wantLimit:      20,
			wantOffset:     0,
			wantResultSize: 1,
			wantTotal:      1,
		},
		{
			name:           "clamp max limit for super admin",
			adminSchoolID:  nil,
			limit:          999,
			offset:         3,
			wantListCalls:  1,
			wantGetByID:    0,
			wantLimit:      100,
			wantOffset:     3,
			wantResultSize: 1,
			wantTotal:      1,
		},
		{
			name:          "list repo error for super admin",
			adminSchoolID: nil,
			limit:         30,
			offset:        0,
			listErr:       sentinelErr,
			wantErr:       sentinelErr,
			wantListCalls: 1,
			wantGetByID:   0,
			wantLimit:     30,
			wantOffset:    0,
		},
		{
			name:           "school admin returns own school only",
			adminSchoolID:  &adminSchoolID,
			limit:          50,
			offset:         9,
			wantListCalls:  0,
			wantGetByID:    1,
			wantResultSize: 1,
			wantTotal:      1,
		},
		{
			name:          "school admin get own school error",
			adminSchoolID: &adminSchoolID,
			limit:         50,
			offset:        9,
			getByIDErr:    sentinelErr,
			wantErr:       sentinelErr,
			wantListCalls: 0,
			wantGetByID:   1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeSchoolRepo{
				listResult: []model.School{{SchoolID: uuid.New(), Name: "Iris", Address: "A"}},
				listTotal:  1,
				listErr:    tc.listErr,
				getByIDResult: &model.School{
					SchoolID: adminSchoolID,
					Name:     "Admin School",
					Address:  "B",
				},
				getByIDErr: tc.getByIDErr,
			}
			svc := &SchoolService{schoolRepo: repo}

			items, total, err := svc.List(context.Background(), tc.adminSchoolID, tc.limit, tc.offset)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("List() error = %v", err)
			}

			if repo.listCalls != tc.wantListCalls {
				t.Fatalf("list calls = %d, want %d", repo.listCalls, tc.wantListCalls)
			}
			if repo.getByIDCalls != tc.wantGetByID {
				t.Fatalf("getByID calls = %d, want %d", repo.getByIDCalls, tc.wantGetByID)
			}
			if tc.wantListCalls > 0 {
				if repo.listLimit != tc.wantLimit || repo.listOffset != tc.wantOffset {
					t.Fatalf("limit/offset = %d/%d, want %d/%d", repo.listLimit, repo.listOffset, tc.wantLimit, tc.wantOffset)
				}
			}
			if tc.wantGetByID > 0 {
				if repo.getByIDArg != adminSchoolID {
					t.Fatalf("getByID arg = %v, want %v", repo.getByIDArg, adminSchoolID)
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

func TestSchoolServiceUpdate(t *testing.T) {
	schoolID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		updateErr       error
		wantErr         error
		wantUpdateCalls int
	}{
		{
			name:            "not found",
			updateErr:       pgx.ErrNoRows,
			wantErr:         ErrSchoolNotFound,
			wantUpdateCalls: 1,
		},
		{
			name:            "repo error",
			updateErr:       sentinelErr,
			wantErr:         sentinelErr,
			wantUpdateCalls: 1,
		},
		{
			name:            "success",
			wantUpdateCalls: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeSchoolRepo{updateErr: tc.updateErr}
			svc := &SchoolService{schoolRepo: repo}

			err := svc.Update(context.Background(), schoolID, "Iris", "District 1")
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("Update() error = %v", err)
			}

			if repo.updateCalls != tc.wantUpdateCalls {
				t.Fatalf("update calls = %d, want %d", repo.updateCalls, tc.wantUpdateCalls)
			}
			if tc.wantUpdateCalls > 0 {
				if repo.updateSchool != schoolID || repo.updateName != "Iris" || repo.updateAddress != "District 1" {
					t.Fatalf("unexpected update args: school=%v name=%q address=%q", repo.updateSchool, repo.updateName, repo.updateAddress)
				}
			}
		})
	}
}

func TestSchoolServiceDelete(t *testing.T) {
	schoolID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		deleteErr       error
		wantErr         error
		wantDeleteCalls int
	}{
		{
			name:            "not found",
			deleteErr:       pgx.ErrNoRows,
			wantErr:         ErrSchoolNotFound,
			wantDeleteCalls: 1,
		},
		{
			name:            "repo error",
			deleteErr:       sentinelErr,
			wantErr:         sentinelErr,
			wantDeleteCalls: 1,
		},
		{
			name:            "success",
			wantDeleteCalls: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeSchoolRepo{deleteErr: tc.deleteErr}
			svc := &SchoolService{schoolRepo: repo}

			err := svc.Delete(context.Background(), schoolID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("Delete() error = %v", err)
			}

			if repo.deleteCalls != tc.wantDeleteCalls {
				t.Fatalf("delete calls = %d, want %d", repo.deleteCalls, tc.wantDeleteCalls)
			}
			if tc.wantDeleteCalls > 0 && repo.deleteSchool != schoolID {
				t.Fatalf("delete school = %v, want %v", repo.deleteSchool, schoolID)
			}
		})
	}
}
