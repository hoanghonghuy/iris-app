package service

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
	"github.com/jackc/pgx/v5"
)

type fakeParentCodeServiceParentCodeRepo struct {
	createCalls      int
	createStudentID  uuid.UUID
	createCode       string
	createMaxUsage   int
	createExpiresAt  time.Time
	createErr        error
	deleteCalls      int
	deleteStudentID  uuid.UUID
	deleteErr        error
	findByCodeCalls  int
	findByCodeArg    string
	findByCodeRes    *model.StudentParentCode
	findByCodeErr    error
	incrementCalls   int
	incrementCodeArg string
	incrementErr     error
	registerTxCalls  int
	registerTxParams repo.RegisterParentTxParams
	registerTxRes    uuid.UUID
	registerTxErr    error
}

func (f *fakeParentCodeServiceParentCodeRepo) Create(_ context.Context, studentID uuid.UUID, code string, maxUsage int, expiresAt time.Time) error {
	f.createCalls++
	f.createStudentID = studentID
	f.createCode = code
	f.createMaxUsage = maxUsage
	f.createExpiresAt = expiresAt
	return f.createErr
}

func (f *fakeParentCodeServiceParentCodeRepo) DeleteByStudentID(_ context.Context, studentID uuid.UUID) error {
	f.deleteCalls++
	f.deleteStudentID = studentID
	return f.deleteErr
}

func (f *fakeParentCodeServiceParentCodeRepo) FindByCode(_ context.Context, code string) (*model.StudentParentCode, error) {
	f.findByCodeCalls++
	f.findByCodeArg = code
	return f.findByCodeRes, f.findByCodeErr
}

func (f *fakeParentCodeServiceParentCodeRepo) IncrementUsageIfNotMaxed(_ context.Context, code string) error {
	f.incrementCalls++
	f.incrementCodeArg = code
	return f.incrementErr
}

func (f *fakeParentCodeServiceParentCodeRepo) RegisterParentTx(_ context.Context, p repo.RegisterParentTxParams) (uuid.UUID, error) {
	f.registerTxCalls++
	f.registerTxParams = p
	return f.registerTxRes, f.registerTxErr
}

type fakeParentCodeServiceUserRepo struct {
	findByEmailCalls int
	findByEmailArg   string
	findByEmailRes   *model.User
	findByEmailErr   error

	createActiveCalls  int
	createActiveEmail  string
	createActiveResult uuid.UUID
	createActiveErr    error

	assignRoleCalls int
	assignRoleUser  uuid.UUID
	assignRoleName  string
	assignRoleErr   error
}

func (f *fakeParentCodeServiceUserRepo) FindByEmail(_ context.Context, email string) (*model.User, error) {
	f.findByEmailCalls++
	f.findByEmailArg = email
	return f.findByEmailRes, f.findByEmailErr
}

func (f *fakeParentCodeServiceUserRepo) CreateActive(_ context.Context, email, _ string) (uuid.UUID, error) {
	f.createActiveCalls++
	f.createActiveEmail = email
	return f.createActiveResult, f.createActiveErr
}

func (f *fakeParentCodeServiceUserRepo) AssignRole(_ context.Context, userID uuid.UUID, roleName string) error {
	f.assignRoleCalls++
	f.assignRoleUser = userID
	f.assignRoleName = roleName
	return f.assignRoleErr
}

type fakeParentCodeServiceParentRepo struct {
	createCalls    int
	createUserID   uuid.UUID
	createSchoolID uuid.UUID
	createFullName string
	createPhone    string
	createResult   uuid.UUID
	createErr      error
}

func (f *fakeParentCodeServiceParentRepo) Create(_ context.Context, userID, schoolID uuid.UUID, fullName, phone string) (uuid.UUID, error) {
	f.createCalls++
	f.createUserID = userID
	f.createSchoolID = schoolID
	f.createFullName = fullName
	f.createPhone = phone
	return f.createResult, f.createErr
}

type fakeParentCodeServiceStudentParentRepo struct {
	assignCalls int
	assignErr   error

	assignStudentID    uuid.UUID
	assignParentID     uuid.UUID
	assignRelationship string
}

func (f *fakeParentCodeServiceStudentParentRepo) Assign(_ context.Context, studentID, parentID uuid.UUID, relationship string) error {
	f.assignCalls++
	f.assignStudentID = studentID
	f.assignParentID = parentID
	f.assignRelationship = relationship
	return f.assignErr
}

type fakeParentCodeServiceStudentRepo struct {
	getByStudentCalls int
	getByStudentArg   uuid.UUID
	getByStudentRes   *model.Student
	getByStudentErr   error

	getSchoolIDCalls int
	getSchoolIDArg   uuid.UUID
	getSchoolIDRes   uuid.UUID
	getSchoolIDErr   error
}

func (f *fakeParentCodeServiceStudentRepo) GetByStudentID(_ context.Context, studentID uuid.UUID) (*model.Student, error) {
	f.getByStudentCalls++
	f.getByStudentArg = studentID
	return f.getByStudentRes, f.getByStudentErr
}

func (f *fakeParentCodeServiceStudentRepo) GetSchoolIDByStudentID(_ context.Context, studentID uuid.UUID) (uuid.UUID, error) {
	f.getSchoolIDCalls++
	f.getSchoolIDArg = studentID
	return f.getSchoolIDRes, f.getSchoolIDErr
}

type fakeGoogleTokenVerifier struct {
	verifyCalls int
	verifyToken string
	verifyRes   *auth.GoogleIdentity
	verifyErr   error
}

func (f *fakeGoogleTokenVerifier) Verify(_ context.Context, rawIDToken string) (*auth.GoogleIdentity, error) {
	f.verifyCalls++
	f.verifyToken = rawIDToken
	if f.verifyErr != nil {
		return nil, f.verifyErr
	}
	return f.verifyRes, nil
}

func TestParentCodeServiceGenerateCodeForStudent(t *testing.T) {
	adminSchoolID := uuid.New()
	otherSchoolID := uuid.New()
	studentID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		adminSchoolID   *uuid.UUID
		studentRes      *model.Student
		studentErr      error
		createErr       error
		wantErr         error
		wantGetCalls    int
		wantDeleteCalls int
		wantCreateCalls int
	}{
		{name: "school-admin student lookup failure", adminSchoolID: &adminSchoolID, studentErr: sentinelErr, wantErr: ErrFailedToGetStudent, wantGetCalls: 1},
		{name: "school-admin cross-school denied", adminSchoolID: &adminSchoolID, studentRes: &model.Student{StudentID: studentID, SchoolID: otherSchoolID}, wantErr: ErrSchoolAccessDenied, wantGetCalls: 1},
		{name: "create error passthrough", adminSchoolID: &adminSchoolID, studentRes: &model.Student{StudentID: studentID, SchoolID: adminSchoolID}, createErr: sentinelErr, wantErr: sentinelErr, wantGetCalls: 1, wantDeleteCalls: 1, wantCreateCalls: 1},
		{name: "success for school-admin", adminSchoolID: &adminSchoolID, studentRes: &model.Student{StudentID: studentID, SchoolID: adminSchoolID}, wantGetCalls: 1, wantDeleteCalls: 1, wantCreateCalls: 1},
		{name: "success for super-admin", adminSchoolID: nil, wantGetCalls: 0, wantDeleteCalls: 1, wantCreateCalls: 1},
	}

	codePattern := regexp.MustCompile(`^[A-Z0-9]{8}$`)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parentCodeRepo := &fakeParentCodeServiceParentCodeRepo{createErr: tc.createErr}
			studentRepo := &fakeParentCodeServiceStudentRepo{getByStudentRes: tc.studentRes, getByStudentErr: tc.studentErr}
			svc := &ParentCodeService{
				parentCodeRepo:    parentCodeRepo,
				userRepo:          &fakeParentCodeServiceUserRepo{},
				parentRepo:        &fakeParentCodeServiceParentRepo{},
				studentParentRepo: &fakeParentCodeServiceStudentParentRepo{},
				studentRepo:       studentRepo,
			}

			before := time.Now().UTC()
			code, err := svc.GenerateCodeForStudent(context.Background(), tc.adminSchoolID, studentID)
			after := time.Now().UTC()

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("GenerateCodeForStudent() error = %v", err)
			}

			if studentRepo.getByStudentCalls != tc.wantGetCalls || parentCodeRepo.deleteCalls != tc.wantDeleteCalls || parentCodeRepo.createCalls != tc.wantCreateCalls {
				t.Fatalf("calls get/delete/create = %d/%d/%d, want %d/%d/%d", studentRepo.getByStudentCalls, parentCodeRepo.deleteCalls, parentCodeRepo.createCalls, tc.wantGetCalls, tc.wantDeleteCalls, tc.wantCreateCalls)
			}

			if tc.wantCreateCalls > 0 {
				if parentCodeRepo.createStudentID != studentID {
					t.Fatalf("create studentID = %v, want %v", parentCodeRepo.createStudentID, studentID)
				}
				if parentCodeRepo.createMaxUsage != 4 {
					t.Fatalf("create maxUsage = %d, want 4", parentCodeRepo.createMaxUsage)
				}
				expiresMin := before.AddDate(0, 0, 7).Add(-2 * time.Second)
				expiresMax := after.AddDate(0, 0, 7).Add(2 * time.Second)
				if parentCodeRepo.createExpiresAt.Before(expiresMin) || parentCodeRepo.createExpiresAt.After(expiresMax) {
					t.Fatalf("expiresAt out of range: %v", parentCodeRepo.createExpiresAt)
				}
				if !codePattern.MatchString(parentCodeRepo.createCode) {
					t.Fatalf("generated code = %q does not match pattern", parentCodeRepo.createCode)
				}
			}

			if tc.wantErr == nil {
				if !codePattern.MatchString(code) {
					t.Fatalf("returned code = %q does not match pattern", code)
				}
			}
		})
	}
}

func TestParentCodeServiceRevokeCode(t *testing.T) {
	adminSchoolID := uuid.New()
	otherSchoolID := uuid.New()
	studentID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		adminSchoolID   *uuid.UUID
		studentRes      *model.Student
		studentErr      error
		deleteErr       error
		wantErr         error
		wantGetCalls    int
		wantDeleteCalls int
	}{
		{name: "school-admin student lookup failure", adminSchoolID: &adminSchoolID, studentErr: sentinelErr, wantErr: ErrFailedToGetStudent, wantGetCalls: 1},
		{name: "school-admin cross-school denied", adminSchoolID: &adminSchoolID, studentRes: &model.Student{StudentID: studentID, SchoolID: otherSchoolID}, wantErr: ErrSchoolAccessDenied, wantGetCalls: 1},
		{name: "delete error passthrough", adminSchoolID: nil, deleteErr: sentinelErr, wantErr: sentinelErr, wantDeleteCalls: 1},
		{name: "success", adminSchoolID: &adminSchoolID, studentRes: &model.Student{StudentID: studentID, SchoolID: adminSchoolID}, wantGetCalls: 1, wantDeleteCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parentCodeRepo := &fakeParentCodeServiceParentCodeRepo{deleteErr: tc.deleteErr}
			studentRepo := &fakeParentCodeServiceStudentRepo{getByStudentRes: tc.studentRes, getByStudentErr: tc.studentErr}
			svc := &ParentCodeService{
				parentCodeRepo:    parentCodeRepo,
				userRepo:          &fakeParentCodeServiceUserRepo{},
				parentRepo:        &fakeParentCodeServiceParentRepo{},
				studentParentRepo: &fakeParentCodeServiceStudentParentRepo{},
				studentRepo:       studentRepo,
			}

			err := svc.RevokeCode(context.Background(), tc.adminSchoolID, studentID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("RevokeCode() error = %v", err)
			}

			if studentRepo.getByStudentCalls != tc.wantGetCalls || parentCodeRepo.deleteCalls != tc.wantDeleteCalls {
				t.Fatalf("calls get/delete = %d/%d, want %d/%d", studentRepo.getByStudentCalls, parentCodeRepo.deleteCalls, tc.wantGetCalls, tc.wantDeleteCalls)
			}
		})
	}
}

func TestParentCodeServiceVerifyCode(t *testing.T) {
	now := time.Now().UTC()
	studentID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name          string
		findRes       *model.StudentParentCode
		findErr       error
		wantErr       error
		wantFindCalls int
	}{
		{name: "invalid code maps to invalid-parent-code", findErr: sentinelErr, wantErr: ErrInvalidParentCode, wantFindCalls: 1},
		{name: "max usage reached", findRes: &model.StudentParentCode{StudentID: studentID, UsageCount: 2, MaxUsage: 2, ExpiresAt: now.Add(time.Hour)}, wantErr: ErrParentCodeMaxUsageReached, wantFindCalls: 1},
		{name: "expired", findRes: &model.StudentParentCode{StudentID: studentID, UsageCount: 1, MaxUsage: 2, ExpiresAt: now.Add(-time.Minute)}, wantErr: ErrParentCodeExpired, wantFindCalls: 1},
		{name: "success", findRes: &model.StudentParentCode{StudentID: studentID, UsageCount: 1, MaxUsage: 2, ExpiresAt: now.Add(time.Hour)}, wantFindCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parentCodeRepo := &fakeParentCodeServiceParentCodeRepo{findByCodeRes: tc.findRes, findByCodeErr: tc.findErr}
			svc := &ParentCodeService{parentCodeRepo: parentCodeRepo}

			got, err := svc.VerifyCode(context.Background(), "ABC12345")
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("VerifyCode() error = %v", err)
			}
			if parentCodeRepo.findByCodeCalls != tc.wantFindCalls {
				t.Fatalf("findByCode calls = %d, want %d", parentCodeRepo.findByCodeCalls, tc.wantFindCalls)
			}
			if tc.wantErr == nil && got == nil {
				t.Fatal("expected code info on success")
			}
		})
	}
}

func TestParentCodeServiceRegisterParent(t *testing.T) {
	now := time.Now().UTC()
	studentID := uuid.New()
	schoolID := uuid.New()
	userID := uuid.New()
	parentID := uuid.New()
	sentinelErr := errors.New("repo failed")

	baseCodeInfo := &model.StudentParentCode{StudentID: studentID, ExpiresAt: now.Add(time.Hour), UsageCount: 0, MaxUsage: 4}

	tests := []struct {
		name                 string
		findCodeRes          *model.StudentParentCode
		findCodeErr          error
		findByEmailErr       error
		createActiveErr      error
		assignRoleErr        error
		getSchoolIDErr       error
		createParentErr      error
		assignStudentErr     error
		incrementErr         error
		wantErr              error
		wantContains         string
		wantCreateUserCalls  int
		wantAssignRoleCalls  int
		wantCreateParentCall int
		wantAssignLinkCalls  int
		wantIncrementCalls   int
	}{
		{name: "invalid code", findCodeErr: sentinelErr, wantErr: ErrInvalidParentCode},
		{name: "expired code", findCodeRes: &model.StudentParentCode{StudentID: studentID, ExpiresAt: now.Add(-time.Minute)}, wantErr: ErrParentCodeExpired},
		{name: "email exists", findCodeRes: baseCodeInfo, findByEmailErr: nil, wantErr: ErrEmailAlreadyExists},
		{name: "find email unexpected error", findCodeRes: baseCodeInfo, findByEmailErr: sentinelErr, wantErr: sentinelErr},
		{name: "create user fails", findCodeRes: baseCodeInfo, findByEmailErr: pgx.ErrNoRows, createActiveErr: sentinelErr, wantErr: ErrFailedToCreateUser, wantCreateUserCalls: 1},
		{name: "assign role fails", findCodeRes: baseCodeInfo, findByEmailErr: pgx.ErrNoRows, assignRoleErr: sentinelErr, wantErr: ErrFailedToAssignRole, wantCreateUserCalls: 1, wantAssignRoleCalls: 1},
		{name: "get school id fails", findCodeRes: baseCodeInfo, findByEmailErr: pgx.ErrNoRows, getSchoolIDErr: sentinelErr, wantErr: ErrFailedToGetStudent, wantCreateUserCalls: 1, wantAssignRoleCalls: 1},
		{name: "create parent fails", findCodeRes: baseCodeInfo, findByEmailErr: pgx.ErrNoRows, createParentErr: sentinelErr, wantErr: ErrFailedToCreateParent, wantCreateUserCalls: 1, wantAssignRoleCalls: 1, wantCreateParentCall: 1},
		{name: "link parent student fails", findCodeRes: baseCodeInfo, findByEmailErr: pgx.ErrNoRows, assignStudentErr: sentinelErr, wantErr: ErrFailedToLinkParentToStudent, wantCreateUserCalls: 1, wantAssignRoleCalls: 1, wantCreateParentCall: 1, wantAssignLinkCalls: 1},
		{name: "increment usage no rows maps max usage reached", findCodeRes: baseCodeInfo, findByEmailErr: pgx.ErrNoRows, incrementErr: repo.ErrNoRowsUpdated, wantErr: ErrParentCodeMaxUsageReached, wantCreateUserCalls: 1, wantAssignRoleCalls: 1, wantCreateParentCall: 1, wantAssignLinkCalls: 1, wantIncrementCalls: 1},
		{name: "increment usage generic error passthrough", findCodeRes: baseCodeInfo, findByEmailErr: pgx.ErrNoRows, incrementErr: sentinelErr, wantErr: sentinelErr, wantCreateUserCalls: 1, wantAssignRoleCalls: 1, wantCreateParentCall: 1, wantAssignLinkCalls: 1, wantIncrementCalls: 1},
		{name: "success", findCodeRes: baseCodeInfo, findByEmailErr: pgx.ErrNoRows, wantCreateUserCalls: 1, wantAssignRoleCalls: 1, wantCreateParentCall: 1, wantAssignLinkCalls: 1, wantIncrementCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parentCodeRepo := &fakeParentCodeServiceParentCodeRepo{findByCodeRes: tc.findCodeRes, findByCodeErr: tc.findCodeErr, incrementErr: tc.incrementErr}
			userRepo := &fakeParentCodeServiceUserRepo{findByEmailRes: &model.User{UserID: uuid.New()}, findByEmailErr: tc.findByEmailErr, createActiveResult: userID, createActiveErr: tc.createActiveErr, assignRoleErr: tc.assignRoleErr}
			parentRepo := &fakeParentCodeServiceParentRepo{createResult: parentID, createErr: tc.createParentErr}
			studentParentRepo := &fakeParentCodeServiceStudentParentRepo{assignErr: tc.assignStudentErr}
			studentRepo := &fakeParentCodeServiceStudentRepo{getSchoolIDRes: schoolID, getSchoolIDErr: tc.getSchoolIDErr}
			svc := &ParentCodeService{
				parentCodeRepo:    parentCodeRepo,
				userRepo:          userRepo,
				parentRepo:        parentRepo,
				studentParentRepo: studentParentRepo,
				studentRepo:       studentRepo,
				jwtAuth:           &auth.Authenticator{Secret: "test-secret", TTLSeconds: 3600},
			}

			resp, err := svc.RegisterParent(context.Background(), "parent@example.com", "Pass1234", "ABC12345")
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if tc.wantContains != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantContains) {
					t.Fatalf("error = %v, want contains %q", err, tc.wantContains)
				}
			} else if err != nil {
				t.Fatalf("RegisterParent() error = %v", err)
			}

			if userRepo.createActiveCalls != tc.wantCreateUserCalls || userRepo.assignRoleCalls != tc.wantAssignRoleCalls || parentRepo.createCalls != tc.wantCreateParentCall || studentParentRepo.assignCalls != tc.wantAssignLinkCalls || parentCodeRepo.incrementCalls != tc.wantIncrementCalls {
				t.Fatalf("calls createUser/assignRole/createParent/link/increment = %d/%d/%d/%d/%d, want %d/%d/%d/%d/%d", userRepo.createActiveCalls, userRepo.assignRoleCalls, parentRepo.createCalls, studentParentRepo.assignCalls, parentCodeRepo.incrementCalls, tc.wantCreateUserCalls, tc.wantAssignRoleCalls, tc.wantCreateParentCall, tc.wantAssignLinkCalls, tc.wantIncrementCalls)
			}

			if tc.wantCreateUserCalls > 0 && userRepo.createActiveEmail != "parent@example.com" {
				t.Fatalf("create user email = %q, want %q", userRepo.createActiveEmail, "parent@example.com")
			}
			if tc.wantAssignRoleCalls > 0 && (userRepo.assignRoleUser != userID || userRepo.assignRoleName != "PARENT") {
				t.Fatalf("assign role args mismatch: user=%v role=%q", userRepo.assignRoleUser, userRepo.assignRoleName)
			}
			if tc.wantCreateParentCall > 0 {
				if parentRepo.createUserID != userID || parentRepo.createSchoolID != schoolID {
					t.Fatalf("create parent args mismatch: user=%v school=%v", parentRepo.createUserID, parentRepo.createSchoolID)
				}
				if parentRepo.createFullName != "" || parentRepo.createPhone != "" {
					t.Fatalf("create parent name/phone should be empty, got %q/%q", parentRepo.createFullName, parentRepo.createPhone)
				}
			}
			if tc.wantAssignLinkCalls > 0 {
				if studentParentRepo.assignStudentID != studentID || studentParentRepo.assignParentID != parentID || studentParentRepo.assignRelationship != "parent" {
					t.Fatalf("assign link args mismatch: student=%v parent=%v rel=%q", studentParentRepo.assignStudentID, studentParentRepo.assignParentID, studentParentRepo.assignRelationship)
				}
			}
			if tc.wantIncrementCalls > 0 && parentCodeRepo.incrementCodeArg != "ABC12345" {
				t.Fatalf("increment code = %q, want %q", parentCodeRepo.incrementCodeArg, "ABC12345")
			}

			if tc.wantErr == nil {
				if resp == nil || resp.AccessToken == "" || resp.TokenType != "Bearer" || resp.ExpiresIn != 3600 {
					t.Fatalf("unexpected response = %#v", resp)
				}
			}
		})
	}
}

func TestParentCodeServiceRegisterParentWithGoogle(t *testing.T) {
	now := time.Now().UTC()
	studentID := uuid.New()
	schoolID := uuid.New()
	parentID := uuid.New()
	sentinelErr := errors.New("repo failed")

	baseIdentity := &auth.GoogleIdentity{Sub: "google-sub", Email: "parent@example.com", EmailVerified: true, HostedDomain: "iris.edu.vn", Name: "Parent Name"}
	baseCodeInfo := &model.StudentParentCode{StudentID: studentID, ExpiresAt: now.Add(time.Hour), UsageCount: 0, MaxUsage: 4}

	tests := []struct {
		name               string
		googleEnabled      bool
		googleHD           string
		verifyRes          *auth.GoogleIdentity
		verifyErr          error
		findCodeRes        *model.StudentParentCode
		findCodeErr        error
		findByEmailErr     error
		getSchoolIDErr     error
		registerTxErr      error
		wantErr            error
		wantErrContains    string
		wantRegisterTxCall int
		wantFallbackName   string
	}{
		{name: "google disabled", googleEnabled: false, verifyRes: baseIdentity, wantErr: ErrGoogleLoginDisabled},
		{name: "verify token error", googleEnabled: true, verifyErr: sentinelErr, wantErr: auth.ErrInvalidCredentials},
		{name: "email not verified", googleEnabled: true, verifyRes: &auth.GoogleIdentity{Sub: "s", Email: "e@example.com", EmailVerified: false, HostedDomain: "iris.edu.vn"}, wantErr: auth.ErrInvalidCredentials},
		{name: "hosted domain not allowed", googleEnabled: true, googleHD: "iris.edu.vn", verifyRes: &auth.GoogleIdentity{Sub: "s", Email: "e@example.com", EmailVerified: true, HostedDomain: "other.edu.vn"}, wantErr: ErrGoogleDomainNotAllowed},
		{name: "invalid parent code", googleEnabled: true, verifyRes: baseIdentity, findCodeErr: sentinelErr, wantErr: ErrInvalidParentCode},
		{name: "expired parent code", googleEnabled: true, verifyRes: baseIdentity, findCodeRes: &model.StudentParentCode{StudentID: studentID, ExpiresAt: now.Add(-time.Minute)}, wantErr: ErrParentCodeExpired},
		{name: "email exists", googleEnabled: true, verifyRes: baseIdentity, findCodeRes: baseCodeInfo, findByEmailErr: nil, wantErr: ErrEmailAlreadyExists},
		{name: "find email unexpected error", googleEnabled: true, verifyRes: baseIdentity, findCodeRes: baseCodeInfo, findByEmailErr: sentinelErr, wantErr: sentinelErr},
		{name: "get school error", googleEnabled: true, verifyRes: baseIdentity, findCodeRes: baseCodeInfo, findByEmailErr: pgx.ErrNoRows, getSchoolIDErr: sentinelErr, wantErr: ErrFailedToGetStudent},
		{name: "register tx max usage", googleEnabled: true, verifyRes: baseIdentity, findCodeRes: baseCodeInfo, findByEmailErr: pgx.ErrNoRows, registerTxErr: repo.ErrNoRowsUpdated, wantErr: ErrParentCodeMaxUsageReached, wantRegisterTxCall: 1},
		{name: "register tx generic error", googleEnabled: true, verifyRes: baseIdentity, findCodeRes: baseCodeInfo, findByEmailErr: pgx.ErrNoRows, registerTxErr: sentinelErr, wantErrContains: "failed to register parent", wantRegisterTxCall: 1},
		{name: "success with fallback full name", googleEnabled: true, verifyRes: &auth.GoogleIdentity{Sub: "google-sub", Email: "parent@example.com", EmailVerified: true, HostedDomain: "iris.edu.vn", Name: ""}, findCodeRes: baseCodeInfo, findByEmailErr: pgx.ErrNoRows, wantRegisterTxCall: 1, wantFallbackName: "Google Parent"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parentCodeRepo := &fakeParentCodeServiceParentCodeRepo{findByCodeRes: tc.findCodeRes, findByCodeErr: tc.findCodeErr, registerTxRes: parentID, registerTxErr: tc.registerTxErr}
			userRepo := &fakeParentCodeServiceUserRepo{findByEmailRes: &model.User{UserID: uuid.New()}, findByEmailErr: tc.findByEmailErr}
			studentRepo := &fakeParentCodeServiceStudentRepo{getSchoolIDRes: schoolID, getSchoolIDErr: tc.getSchoolIDErr}
			verifier := &fakeGoogleTokenVerifier{verifyRes: tc.verifyRes, verifyErr: tc.verifyErr}
			svc := &ParentCodeService{
				parentCodeRepo:    parentCodeRepo,
				userRepo:          userRepo,
				parentRepo:        &fakeParentCodeServiceParentRepo{},
				studentParentRepo: &fakeParentCodeServiceStudentParentRepo{},
				studentRepo:       studentRepo,
				jwtAuth:           &auth.Authenticator{Secret: "test-secret", TTLSeconds: 3600},
				googleVerifier:    verifier,
				googleEnabled:     tc.googleEnabled,
				googleHD:          tc.googleHD,
			}

			resp, err := svc.RegisterParentWithGoogle(context.Background(), "google-token", "ABC12345")
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if tc.wantErrContains != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErrContains) {
					t.Fatalf("error = %v, want contains %q", err, tc.wantErrContains)
				}
			} else if err != nil {
				t.Fatalf("RegisterParentWithGoogle() error = %v", err)
			}

			if parentCodeRepo.registerTxCalls != tc.wantRegisterTxCall {
				t.Fatalf("register tx calls = %d, want %d", parentCodeRepo.registerTxCalls, tc.wantRegisterTxCall)
			}
			if tc.wantRegisterTxCall > 0 {
				if parentCodeRepo.registerTxParams.Email != "parent@example.com" || parentCodeRepo.registerTxParams.StudentID != studentID || parentCodeRepo.registerTxParams.SchoolID != schoolID || parentCodeRepo.registerTxParams.Code != "ABC12345" || parentCodeRepo.registerTxParams.GoogleSub == "" {
					t.Fatalf("unexpected register tx params: %#v", parentCodeRepo.registerTxParams)
				}
				if tc.wantFallbackName != "" && parentCodeRepo.registerTxParams.FullName != tc.wantFallbackName {
					t.Fatalf("fullName = %q, want %q", parentCodeRepo.registerTxParams.FullName, tc.wantFallbackName)
				}
			}

			if tc.wantErr == nil && tc.wantErrContains == "" {
				if resp == nil || resp.AccessToken == "" || resp.TokenType != "Bearer" || resp.ExpiresIn != 3600 {
					t.Fatalf("unexpected response = %#v", resp)
				}
			}
		})
	}
}

func TestParentCodeServiceGetStudentInfo(t *testing.T) {
	studentID := uuid.New()
	sentinelErr := errors.New("repo failed")

	t.Run("repo error passthrough", func(t *testing.T) {
		studentRepo := &fakeParentCodeServiceStudentRepo{getByStudentErr: sentinelErr}
		svc := &ParentCodeService{studentRepo: studentRepo}
		_, err := svc.GetStudentInfo(context.Background(), studentID)
		if !errors.Is(err, sentinelErr) {
			t.Fatalf("error = %v, want %v", err, sentinelErr)
		}
	})

	t.Run("success", func(t *testing.T) {
		studentRepo := &fakeParentCodeServiceStudentRepo{getByStudentRes: &model.Student{StudentID: studentID, FullName: "Student A"}}
		svc := &ParentCodeService{studentRepo: studentRepo}
		got, err := svc.GetStudentInfo(context.Background(), studentID)
		if err != nil {
			t.Fatalf("GetStudentInfo() error = %v", err)
		}
		if got == nil || got.StudentID != studentID || got.FullName != "Student A" {
			t.Fatalf("unexpected student = %#v", got)
		}
	})
}
