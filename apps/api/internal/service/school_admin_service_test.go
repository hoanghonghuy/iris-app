package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

type fakeSchoolAdminRepo struct {
	createCalls    int
	createUserID   uuid.UUID
	createSchoolID uuid.UUID
	createFullName string
	createPhone    string
	createResultID uuid.UUID
	createErr      error

	getByAdminCalls  int
	getByAdminArg    uuid.UUID
	getByAdminResult *model.SchoolAdmin
	getByAdminErr    error

	listCalls  int
	listLimit  int
	listOffset int
	listResult []model.SchoolAdmin
	listTotal  int
	listErr    error

	deleteCalls int
	deleteAdmin uuid.UUID
	deleteErr   error
}

func (f *fakeSchoolAdminRepo) Create(_ context.Context, userID, schoolID uuid.UUID, fullName, phone string) (uuid.UUID, error) {
	f.createCalls++
	f.createUserID = userID
	f.createSchoolID = schoolID
	f.createFullName = fullName
	f.createPhone = phone
	return f.createResultID, f.createErr
}

func (f *fakeSchoolAdminRepo) GetByAdminID(_ context.Context, adminID uuid.UUID) (*model.SchoolAdmin, error) {
	f.getByAdminCalls++
	f.getByAdminArg = adminID
	return f.getByAdminResult, f.getByAdminErr
}

func (f *fakeSchoolAdminRepo) List(_ context.Context, limit, offset int) ([]model.SchoolAdmin, int, error) {
	f.listCalls++
	f.listLimit = limit
	f.listOffset = offset
	return f.listResult, f.listTotal, f.listErr
}

func (f *fakeSchoolAdminRepo) Delete(_ context.Context, adminID uuid.UUID) error {
	f.deleteCalls++
	f.deleteAdmin = adminID
	return f.deleteErr
}

type fakeSchoolAdminUserRepo struct {
	findByIDCalls  int
	findByIDArg    uuid.UUID
	findByIDResult *model.UserInfo
	findByIDErr    error

	assignRoleCalls int
	assignRoleUser  uuid.UUID
	assignRoleName  string
	assignRoleErr   error
}

func (f *fakeSchoolAdminUserRepo) FindByID(_ context.Context, userID uuid.UUID) (*model.UserInfo, error) {
	f.findByIDCalls++
	f.findByIDArg = userID
	return f.findByIDResult, f.findByIDErr
}

func (f *fakeSchoolAdminUserRepo) AssignRole(_ context.Context, userID uuid.UUID, roleName string) error {
	f.assignRoleCalls++
	f.assignRoleUser = userID
	f.assignRoleName = roleName
	return f.assignRoleErr
}

func TestSchoolAdminServiceCreate(t *testing.T) {
	userID := uuid.New()
	schoolID := uuid.New()
	adminID := uuid.New()
	sentinelErr := errors.New("repo failed")

	admin := &model.SchoolAdmin{
		AdminID:  adminID,
		UserID:   userID,
		Email:    "admin@example.com",
		FullName: "Admin Name",
		Phone:    "0909",
		SchoolID: schoolID,
	}

	tests := []struct {
		name                string
		fullName            string
		findUserResult      *model.UserInfo
		findUserErr         error
		assignRoleErr       error
		createErr           error
		getByAdminErr       error
		wantErr             error
		wantFindByIDCalls   int
		wantAssignRoleCalls int
		wantCreateCalls     int
		wantGetByAdminCalls int
		wantForwardedName   string
	}{
		{
			name:                "provided fullName skips user lookup",
			fullName:            "Provided Name",
			wantFindByIDCalls:   0,
			wantAssignRoleCalls: 1,
			wantCreateCalls:     1,
			wantGetByAdminCalls: 1,
			wantForwardedName:   "Provided Name",
		},
		{
			name:              "missing fullName and user lookup fails",
			fullName:          "",
			findUserErr:       sentinelErr,
			wantErr:           ErrUserNotFound,
			wantFindByIDCalls: 1,
		},
		{
			name:                "missing fullName uses user fullName",
			fullName:            "",
			findUserResult:      &model.UserInfo{UserID: userID, Email: "u@example.com", FullName: "Derived Name"},
			wantFindByIDCalls:   1,
			wantAssignRoleCalls: 1,
			wantCreateCalls:     1,
			wantGetByAdminCalls: 1,
			wantForwardedName:   "Derived Name",
		},
		{
			name:                "missing fullName falls back to email",
			fullName:            "",
			findUserResult:      &model.UserInfo{UserID: userID, Email: "fallback@example.com", FullName: ""},
			wantFindByIDCalls:   1,
			wantAssignRoleCalls: 1,
			wantCreateCalls:     1,
			wantGetByAdminCalls: 1,
			wantForwardedName:   "fallback@example.com",
		},
		{
			name:                "assign role fails",
			fullName:            "Provided Name",
			assignRoleErr:       sentinelErr,
			wantErr:             ErrFailedToAssignRole,
			wantAssignRoleCalls: 1,
		},
		{
			name:                "create fails",
			fullName:            "Provided Name",
			createErr:           sentinelErr,
			wantErr:             sentinelErr,
			wantAssignRoleCalls: 1,
			wantCreateCalls:     1,
			wantForwardedName:   "Provided Name",
		},
		{
			name:                "getByAdmin after create fails",
			fullName:            "Provided Name",
			getByAdminErr:       sentinelErr,
			wantErr:             sentinelErr,
			wantAssignRoleCalls: 1,
			wantCreateCalls:     1,
			wantGetByAdminCalls: 1,
			wantForwardedName:   "Provided Name",
		},
		{
			name:                "success",
			fullName:            "Provided Name",
			wantAssignRoleCalls: 1,
			wantCreateCalls:     1,
			wantGetByAdminCalls: 1,
			wantForwardedName:   "Provided Name",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			adminRepo := &fakeSchoolAdminRepo{
				createResultID:   adminID,
				createErr:        tc.createErr,
				getByAdminResult: admin,
				getByAdminErr:    tc.getByAdminErr,
			}
			userRepo := &fakeSchoolAdminUserRepo{
				findByIDResult: tc.findUserResult,
				findByIDErr:    tc.findUserErr,
				assignRoleErr:  tc.assignRoleErr,
			}
			svc := &SchoolAdminService{schoolAdminRepo: adminRepo, userRepo: userRepo}

			got, err := svc.Create(context.Background(), userID, schoolID, tc.fullName, "0909")
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("Create() error = %v", err)
			}

			if userRepo.findByIDCalls != tc.wantFindByIDCalls {
				t.Fatalf("findByID calls = %d, want %d", userRepo.findByIDCalls, tc.wantFindByIDCalls)
			}
			if userRepo.assignRoleCalls != tc.wantAssignRoleCalls {
				t.Fatalf("assignRole calls = %d, want %d", userRepo.assignRoleCalls, tc.wantAssignRoleCalls)
			}
			if adminRepo.createCalls != tc.wantCreateCalls {
				t.Fatalf("create calls = %d, want %d", adminRepo.createCalls, tc.wantCreateCalls)
			}
			if adminRepo.getByAdminCalls != tc.wantGetByAdminCalls {
				t.Fatalf("getByAdmin calls = %d, want %d", adminRepo.getByAdminCalls, tc.wantGetByAdminCalls)
			}

			if tc.wantCreateCalls > 0 {
				if adminRepo.createUserID != userID || adminRepo.createSchoolID != schoolID || adminRepo.createPhone != "0909" {
					t.Fatalf("unexpected create args: user=%v school=%v phone=%q", adminRepo.createUserID, adminRepo.createSchoolID, adminRepo.createPhone)
				}
				if adminRepo.createFullName != tc.wantForwardedName {
					t.Fatalf("create fullName = %q, want %q", adminRepo.createFullName, tc.wantForwardedName)
				}
			}
			if tc.wantAssignRoleCalls > 0 {
				if userRepo.assignRoleUser != userID || userRepo.assignRoleName != "SCHOOL_ADMIN" {
					t.Fatalf("assign role args mismatch: user=%v role=%q", userRepo.assignRoleUser, userRepo.assignRoleName)
				}
			}

			if tc.wantErr == nil {
				if got == nil || got.AdminID != adminID {
					t.Fatalf("unexpected school admin = %#v", got)
				}
			}
		})
	}
}

func TestSchoolAdminServiceGetByAdminID(t *testing.T) {
	adminID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name         string
		repoErr      error
		wantErr      error
		wantRepoCall int
	}{
		{name: "not found mapping", repoErr: sentinelErr, wantErr: ErrSchoolAdminNotFound, wantRepoCall: 1},
		{name: "success", repoErr: nil, wantRepoCall: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeSchoolAdminRepo{getByAdminResult: &model.SchoolAdmin{AdminID: adminID}, getByAdminErr: tc.repoErr}
			svc := &SchoolAdminService{schoolAdminRepo: repo, userRepo: &fakeSchoolAdminUserRepo{}}

			got, err := svc.GetByAdminID(context.Background(), adminID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("GetByAdminID() error = %v", err)
			}
			if repo.getByAdminCalls != tc.wantRepoCall {
				t.Fatalf("repo calls = %d, want %d", repo.getByAdminCalls, tc.wantRepoCall)
			}
			if tc.wantErr == nil && (got == nil || got.AdminID != adminID) {
				t.Fatalf("unexpected school admin = %#v", got)
			}
		})
	}
}

func TestSchoolAdminServiceList(t *testing.T) {
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name           string
		limit          int
		offset         int
		listErr        error
		wantErr        error
		wantListCalls  int
		wantLimit      int
		wantOffset     int
		wantResultSize int
		wantTotal      int
	}{
		{
			name:           "default limit and offset",
			limit:          0,
			offset:         -8,
			wantListCalls:  1,
			wantLimit:      20,
			wantOffset:     0,
			wantResultSize: 1,
			wantTotal:      1,
		},
		{
			name:           "clamp max limit",
			limit:          999,
			offset:         3,
			wantListCalls:  1,
			wantLimit:      100,
			wantOffset:     3,
			wantResultSize: 1,
			wantTotal:      1,
		},
		{
			name:          "repo error",
			limit:         40,
			offset:        2,
			listErr:       sentinelErr,
			wantErr:       sentinelErr,
			wantListCalls: 1,
			wantLimit:     40,
			wantOffset:    2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeSchoolAdminRepo{
				listResult: []model.SchoolAdmin{{AdminID: uuid.New(), FullName: "Admin"}},
				listTotal:  1,
				listErr:    tc.listErr,
			}
			svc := &SchoolAdminService{schoolAdminRepo: repo, userRepo: &fakeSchoolAdminUserRepo{}}

			items, total, err := svc.List(context.Background(), tc.limit, tc.offset)
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
			if repo.listLimit != tc.wantLimit || repo.listOffset != tc.wantOffset {
				t.Fatalf("limit/offset = %d/%d, want %d/%d", repo.listLimit, repo.listOffset, tc.wantLimit, tc.wantOffset)
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

func TestSchoolAdminServiceDelete(t *testing.T) {
	adminID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		deleteErr       error
		wantErr         error
		wantDeleteCalls int
	}{
		{name: "repo error passthrough", deleteErr: sentinelErr, wantErr: sentinelErr, wantDeleteCalls: 1},
		{name: "success", wantDeleteCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeSchoolAdminRepo{deleteErr: tc.deleteErr}
			svc := &SchoolAdminService{schoolAdminRepo: repo, userRepo: &fakeSchoolAdminUserRepo{}}

			err := svc.Delete(context.Background(), adminID)
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
			if tc.wantDeleteCalls > 0 && repo.deleteAdmin != adminID {
				t.Fatalf("delete admin = %v, want %v", repo.deleteAdmin, adminID)
			}
		})
	}
}
