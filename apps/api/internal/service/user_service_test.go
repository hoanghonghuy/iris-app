package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

type fakeUserServiceUserRepo struct {
	createWithRolesCalls  int
	createWithRolesEmail  string
	createWithRolesStatus string
	createWithRolesRoles  []string
	createWithRolesResult uuid.UUID
	createWithRolesErr    error

	findByActivationTokenCalls  int
	findByActivationTokenResult *model.UserWithToken
	findByActivationTokenErr    error

	activateWithPasswordCalls int
	activateWithPasswordUser  uuid.UUID
	activateWithPasswordErr   error

	assignRoleCalls int
	assignRoleUser  uuid.UUID
	assignRoleName  string
	assignRoleErr   error

	findByEmailCalls  int
	findByEmailArg    string
	findByEmailResult *model.User
	findByEmailErr    error

	findByIDCalls  int
	findByIDArg    uuid.UUID
	findByIDResult *model.UserInfo
	findByIDErr    error

	isUserInSchoolCalls  int
	isUserInSchoolResult bool
	isUserInSchoolErr    error

	listCalls      int
	listRoleFilter string
	listLimit      int
	listOffset     int
	listResult     []model.UserInfo
	listTotal      int
	listErr        error

	updateEmailCalls int
	updateEmailUser  uuid.UUID
	updateEmailEmail string
	updateEmailErr   error

	updatePasswordCalls int
	updatePasswordUser  uuid.UUID
	updatePasswordEmail string
	updatePasswordErr   error

	deleteCalls int
	deleteUser  uuid.UUID
	deleteErr   error

	lockCalls int
	lockUser  uuid.UUID
	lockErr   error

	unlockCalls int
	unlockUser  uuid.UUID
	unlockErr   error
}

func (f *fakeUserServiceUserRepo) CreateWithRolesTx(_ context.Context, email, _ string, status string, roles []string) (uuid.UUID, error) {
	f.createWithRolesCalls++
	f.createWithRolesEmail = email
	f.createWithRolesStatus = status
	f.createWithRolesRoles = append([]string{}, roles...)
	return f.createWithRolesResult, f.createWithRolesErr
}

func (f *fakeUserServiceUserRepo) FindByActivationToken(_ context.Context, _ string) (*model.UserWithToken, error) {
	f.findByActivationTokenCalls++
	return f.findByActivationTokenResult, f.findByActivationTokenErr
}

func (f *fakeUserServiceUserRepo) ActivateWithPassword(_ context.Context, userID uuid.UUID, _ string) error {
	f.activateWithPasswordCalls++
	f.activateWithPasswordUser = userID
	return f.activateWithPasswordErr
}

func (f *fakeUserServiceUserRepo) AssignRole(_ context.Context, userID uuid.UUID, roleName string) error {
	f.assignRoleCalls++
	f.assignRoleUser = userID
	f.assignRoleName = roleName
	return f.assignRoleErr
}

func (f *fakeUserServiceUserRepo) FindByEmail(_ context.Context, email string) (*model.User, error) {
	f.findByEmailCalls++
	f.findByEmailArg = email
	return f.findByEmailResult, f.findByEmailErr
}

func (f *fakeUserServiceUserRepo) RolesOfUser(_ context.Context, _ uuid.UUID) ([]string, error) {
	return nil, nil
}

func (f *fakeUserServiceUserRepo) FindByID(_ context.Context, userID uuid.UUID) (*model.UserInfo, error) {
	f.findByIDCalls++
	f.findByIDArg = userID
	return f.findByIDResult, f.findByIDErr
}

func (f *fakeUserServiceUserRepo) IsUserInSchool(_ context.Context, _, _ uuid.UUID) (bool, error) {
	f.isUserInSchoolCalls++
	return f.isUserInSchoolResult, f.isUserInSchoolErr
}

func (f *fakeUserServiceUserRepo) List(_ context.Context, _ *uuid.UUID, roleFilter string, limit, offset int) ([]model.UserInfo, int, error) {
	f.listCalls++
	f.listRoleFilter = roleFilter
	f.listLimit = limit
	f.listOffset = offset
	return f.listResult, f.listTotal, f.listErr
}

func (f *fakeUserServiceUserRepo) UpdateEmail(_ context.Context, userID uuid.UUID, email string) error {
	f.updateEmailCalls++
	f.updateEmailUser = userID
	f.updateEmailEmail = email
	return f.updateEmailErr
}

func (f *fakeUserServiceUserRepo) UpdatePassword(_ context.Context, userID uuid.UUID, email, _ string) error {
	f.updatePasswordCalls++
	f.updatePasswordUser = userID
	f.updatePasswordEmail = email
	return f.updatePasswordErr
}

func (f *fakeUserServiceUserRepo) Delete(_ context.Context, userID uuid.UUID) error {
	f.deleteCalls++
	f.deleteUser = userID
	return f.deleteErr
}

func (f *fakeUserServiceUserRepo) Lock(_ context.Context, userID uuid.UUID) error {
	f.lockCalls++
	f.lockUser = userID
	return f.lockErr
}

func (f *fakeUserServiceUserRepo) Unlock(_ context.Context, userID uuid.UUID) error {
	f.unlockCalls++
	f.unlockUser = userID
	return f.unlockErr
}

type fakeUserServiceResetTokenRepo struct {
	createCalls     int
	createErr       error
	findByHashCalls int
	findByHashRes   *model.ResetToken
	findByHashErr   error
	markUsedCalls   int
	markUsedErr     error
}

func (f *fakeUserServiceResetTokenRepo) Create(_ context.Context, _ uuid.UUID, _ string, _ time.Time) error {
	f.createCalls++
	return f.createErr
}

func (f *fakeUserServiceResetTokenRepo) FindByTokenHash(_ context.Context, _ string) (*model.ResetToken, error) {
	f.findByHashCalls++
	return f.findByHashRes, f.findByHashErr
}

func (f *fakeUserServiceResetTokenRepo) MarkUsed(_ context.Context, _ uuid.UUID) error {
	f.markUsedCalls++
	return f.markUsedErr
}

func TestUserServiceCreateUserWithoutPassword(t *testing.T) {
	schoolID := uuid.New()
	createdUserID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		adminSchoolID   *uuid.UUID
		email           string
		roles           []string
		createErr       error
		findByIDErr     error
		wantErr         error
		wantCreateCalls int
		wantFindByID    int
		wantStatus      string
		wantRoles       []string
	}{
		{name: "empty email", email: "", roles: []string{"PARENT"}, wantErr: ErrEmailCannotBeEmpty},
		{name: "empty roles", email: "a@example.com", roles: nil, wantErr: ErrRolesCannotBeEmpty},
		{name: "invalid role", email: "a@example.com", roles: []string{"INVALID"}, wantErr: ErrInvalidRoleName},
		{name: "reject super admin", email: "a@example.com", roles: []string{"SUPER_ADMIN"}, wantErr: ErrCannotAssignRoleSuperAdmin},
		{name: "school admin cannot assign school admin", adminSchoolID: &schoolID, email: "a@example.com", roles: []string{"SCHOOL_ADMIN"}, wantErr: ErrCannotAssignRole},
		{name: "map role assignment failure", email: "a@example.com", roles: []string{"TEACHER"}, createErr: fmt.Errorf("%w: TEACHER", repo.ErrRoleAssignmentFailed), wantErr: ErrFailedToAssignRole, wantCreateCalls: 1, wantStatus: "pending", wantRoles: []string{"TEACHER"}},
		{name: "map generic create failure", email: "a@example.com", roles: []string{"PARENT"}, createErr: sentinelErr, wantErr: ErrFailedToCreateUser, wantCreateCalls: 1, wantStatus: "active", wantRoles: []string{"PARENT"}},
		{name: "find by id error propagated", email: "a@example.com", roles: []string{"PARENT"}, findByIDErr: sentinelErr, wantErr: sentinelErr, wantCreateCalls: 1, wantFindByID: 1, wantStatus: "active", wantRoles: []string{"PARENT"}},
		{name: "success deduplicate roles and pending for teacher", email: "a@example.com", roles: []string{"TEACHER", "TEACHER", "PARENT"}, wantCreateCalls: 1, wantFindByID: 1, wantStatus: "pending", wantRoles: []string{"TEACHER", "PARENT"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := &fakeUserServiceUserRepo{
				createWithRolesResult: createdUserID,
				createWithRolesErr:    tc.createErr,
				findByIDResult:        &model.UserInfo{UserID: createdUserID, Email: tc.email},
				findByIDErr:           tc.findByIDErr,
			}
			svc := &UserService{userRepo: userRepo, resetTokenRepo: &fakeUserServiceResetTokenRepo{}}

			got, err := svc.CreateUserWithoutPassword(context.Background(), tc.adminSchoolID, tc.email, tc.roles)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("CreateUserWithoutPassword() error = %v", err)
			}

			if userRepo.createWithRolesCalls != tc.wantCreateCalls {
				t.Fatalf("create calls = %d, want %d", userRepo.createWithRolesCalls, tc.wantCreateCalls)
			}
			if userRepo.findByIDCalls != tc.wantFindByID {
				t.Fatalf("findByID calls = %d, want %d", userRepo.findByIDCalls, tc.wantFindByID)
			}
			if tc.wantCreateCalls > 0 {
				if userRepo.createWithRolesEmail != tc.email {
					t.Fatalf("create email = %q, want %q", userRepo.createWithRolesEmail, tc.email)
				}
				if userRepo.createWithRolesStatus != tc.wantStatus {
					t.Fatalf("create status = %q, want %q", userRepo.createWithRolesStatus, tc.wantStatus)
				}
				if !reflect.DeepEqual(userRepo.createWithRolesRoles, tc.wantRoles) {
					t.Fatalf("create roles = %#v, want %#v", userRepo.createWithRolesRoles, tc.wantRoles)
				}
			}
			if tc.wantErr == nil && (got == nil || got.UserID != createdUserID) {
				t.Fatalf("unexpected user info = %#v", got)
			}
		})
	}
}

func TestUserServiceAssignRole(t *testing.T) {
	userID := mustUUID(t, "f7f1d4cb-9708-4fa2-b1ab-e7f58d2bb1ee")

	tests := []struct {
		name            string
		roleName        string
		wantErr         error
		wantAssignCalls int
	}{
		{name: "reject super admin", roleName: "SUPER_ADMIN", wantErr: ErrCannotAssignRoleSuperAdmin},
		{name: "reject invalid role", roleName: "INVALID", wantErr: ErrInvalidRoleName},
		{name: "success", roleName: "PARENT", wantAssignCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := &fakeUserServiceUserRepo{}
			svc := &UserService{userRepo: userRepo, resetTokenRepo: &fakeUserServiceResetTokenRepo{}}

			err := svc.AssignRole(context.Background(), userID, tc.roleName)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("AssignRole() error = %v", err)
			}
			if userRepo.assignRoleCalls != tc.wantAssignCalls {
				t.Fatalf("assign calls = %d, want %d", userRepo.assignRoleCalls, tc.wantAssignCalls)
			}
		})
	}
}

func TestUserServiceFindByID(t *testing.T) {
	schoolID := uuid.New()
	userID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name              string
		adminSchoolID     *uuid.UUID
		findByIDErr       error
		isInSchoolResult  bool
		isInSchoolErr     error
		wantErr           error
		wantFindByIDCalls int
		wantInSchoolCalls int
	}{
		{name: "find user error", adminSchoolID: nil, findByIDErr: sentinelErr, wantErr: sentinelErr, wantFindByIDCalls: 1},
		{name: "school scope check error", adminSchoolID: &schoolID, isInSchoolErr: sentinelErr, wantErr: sentinelErr, wantFindByIDCalls: 1, wantInSchoolCalls: 1},
		{name: "school scope denied", adminSchoolID: &schoolID, isInSchoolResult: false, wantErr: ErrSchoolAccessDenied, wantFindByIDCalls: 1, wantInSchoolCalls: 1},
		{name: "success super admin", adminSchoolID: nil, wantFindByIDCalls: 1},
		{name: "success school admin", adminSchoolID: &schoolID, isInSchoolResult: true, wantFindByIDCalls: 1, wantInSchoolCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := &fakeUserServiceUserRepo{findByIDResult: &model.UserInfo{UserID: userID}, findByIDErr: tc.findByIDErr, isUserInSchoolResult: tc.isInSchoolResult, isUserInSchoolErr: tc.isInSchoolErr}
			svc := &UserService{userRepo: userRepo, resetTokenRepo: &fakeUserServiceResetTokenRepo{}}

			got, err := svc.FindByID(context.Background(), tc.adminSchoolID, userID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("FindByID() error = %v", err)
			}
			if userRepo.findByIDCalls != tc.wantFindByIDCalls || userRepo.isUserInSchoolCalls != tc.wantInSchoolCalls {
				t.Fatalf("find/scope calls = %d/%d, want %d/%d", userRepo.findByIDCalls, userRepo.isUserInSchoolCalls, tc.wantFindByIDCalls, tc.wantInSchoolCalls)
			}
			if tc.wantErr == nil && (got == nil || got.UserID != userID) {
				t.Fatalf("unexpected user info = %#v", got)
			}
		})
	}
}

func TestUserServiceList(t *testing.T) {
	schoolID := uuid.New()

	tests := []struct {
		name          string
		roleFilter    string
		limit         int
		offset        int
		wantErr       error
		wantListCalls int
		wantLimit     int
		wantOffset    int
	}{
		{name: "invalid role", roleFilter: "INVALID", limit: 20, offset: 0, wantErr: ErrInvalidRoleName},
		{name: "default normalize", roleFilter: "", limit: 0, offset: -10, wantListCalls: 1, wantLimit: 20, wantOffset: 0},
		{name: "clamp max", roleFilter: "", limit: 999, offset: 4, wantListCalls: 1, wantLimit: 100, wantOffset: 4},
		{name: "preserve valid", roleFilter: "PARENT", limit: 30, offset: 2, wantListCalls: 1, wantLimit: 30, wantOffset: 2},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := &fakeUserServiceUserRepo{listResult: []model.UserInfo{{UserID: uuid.New()}}, listTotal: 1}
			svc := &UserService{userRepo: userRepo, resetTokenRepo: &fakeUserServiceResetTokenRepo{}}

			_, _, err := svc.List(context.Background(), &schoolID, tc.roleFilter, tc.limit, tc.offset)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("List() error = %v", err)
			}
			if userRepo.listCalls != tc.wantListCalls {
				t.Fatalf("list calls = %d, want %d", userRepo.listCalls, tc.wantListCalls)
			}
			if tc.wantListCalls > 0 && (userRepo.listLimit != tc.wantLimit || userRepo.listOffset != tc.wantOffset) {
				t.Fatalf("limit/offset = %d/%d, want %d/%d", userRepo.listLimit, userRepo.listOffset, tc.wantLimit, tc.wantOffset)
			}
		})
	}
}

func TestUserServiceUpdateEmail(t *testing.T) {
	userID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		email           string
		findByIDErr     error
		updateEmailErr  error
		wantErr         error
		wantFindByID    int
		wantUpdateEmail int
	}{
		{name: "empty email", email: "", wantErr: ErrEmailCannotBeEmpty},
		{name: "user not found", email: "new@example.com", findByIDErr: sentinelErr, wantErr: ErrUserNotFound, wantFindByID: 1},
		{name: "update error", email: "new@example.com", updateEmailErr: sentinelErr, wantErr: sentinelErr, wantFindByID: 1, wantUpdateEmail: 1},
		{name: "success", email: "new@example.com", wantFindByID: 1, wantUpdateEmail: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := &fakeUserServiceUserRepo{findByIDResult: &model.UserInfo{UserID: userID, Email: "old@example.com"}, findByIDErr: tc.findByIDErr, updateEmailErr: tc.updateEmailErr}
			svc := &UserService{userRepo: userRepo, resetTokenRepo: &fakeUserServiceResetTokenRepo{}}

			err := svc.UpdateEmail(context.Background(), userID, tc.email)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("UpdateEmail() error = %v", err)
			}
			if userRepo.findByIDCalls != tc.wantFindByID || userRepo.updateEmailCalls != tc.wantUpdateEmail {
				t.Fatalf("find/update calls = %d/%d, want %d/%d", userRepo.findByIDCalls, userRepo.updateEmailCalls, tc.wantFindByID, tc.wantUpdateEmail)
			}
		})
	}
}

func TestUserServiceUpdateMyPassword(t *testing.T) {
	userID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name               string
		userID             uuid.UUID
		password           string
		findByIDErr        error
		updatePasswordErr  error
		wantErr            error
		wantFindByIDCalls  int
		wantUpdatePwdCalls int
	}{
		{name: "invalid user id", userID: uuid.Nil, password: "abc", wantErr: ErrInvalidUserID},
		{name: "empty password", userID: userID, password: "", wantErr: ErrPasswordCannotBeEmpty},
		{name: "user not found", userID: userID, password: "abc", findByIDErr: sentinelErr, wantErr: ErrUserNotFound, wantFindByIDCalls: 1},
		{name: "update password failed", userID: userID, password: "abc", updatePasswordErr: sentinelErr, wantErr: ErrFailedToUpdatePassword, wantFindByIDCalls: 1, wantUpdatePwdCalls: 1},
		{name: "success", userID: userID, password: "abc", wantFindByIDCalls: 1, wantUpdatePwdCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := &fakeUserServiceUserRepo{findByIDResult: &model.UserInfo{UserID: userID, Email: "u@example.com"}, findByIDErr: tc.findByIDErr, updatePasswordErr: tc.updatePasswordErr}
			svc := &UserService{userRepo: userRepo, resetTokenRepo: &fakeUserServiceResetTokenRepo{}}

			err := svc.UpdateMyPassword(context.Background(), tc.userID, tc.password)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("UpdateMyPassword() error = %v", err)
			}
			if userRepo.findByIDCalls != tc.wantFindByIDCalls || userRepo.updatePasswordCalls != tc.wantUpdatePwdCalls {
				t.Fatalf("find/update calls = %d/%d, want %d/%d", userRepo.findByIDCalls, userRepo.updatePasswordCalls, tc.wantFindByIDCalls, tc.wantUpdatePwdCalls)
			}
			if tc.wantUpdatePwdCalls > 0 && userRepo.updatePasswordEmail != "u@example.com" {
				t.Fatalf("update email = %q, want %q", userRepo.updatePasswordEmail, "u@example.com")
			}
		})
	}
}

func TestUserServiceActivateUserWithToken(t *testing.T) {
	userID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name               string
		token              string
		password           string
		findTokenResult    *model.UserWithToken
		findTokenErr       error
		activateErr        error
		wantErr            error
		wantFindTokenCalls int
		wantActivateCalls  int
	}{
		{name: "missing token", token: "", password: "abc", wantErr: ErrActivationTokenRequired},
		{name: "missing password", token: "token", password: "", wantErr: ErrPasswordCannotBeEmpty},
		{name: "invalid token", token: "token", password: "abc", findTokenErr: sentinelErr, wantErr: ErrInvalidActivationToken, wantFindTokenCalls: 1},
		{name: "expired token", token: "token", password: "abc", findTokenResult: &model.UserWithToken{UserID: userID, TokenExpiresAt: time.Now().Add(-time.Minute)}, wantErr: ErrActivationTokenExpired, wantFindTokenCalls: 1},
		{name: "activate fails", token: "token", password: "abc", findTokenResult: &model.UserWithToken{UserID: userID, TokenExpiresAt: time.Now().Add(time.Minute)}, activateErr: sentinelErr, wantErr: ErrFailedToActivateUser, wantFindTokenCalls: 1, wantActivateCalls: 1},
		{name: "success", token: "token", password: "abc", findTokenResult: &model.UserWithToken{UserID: userID, TokenExpiresAt: time.Now().Add(time.Minute)}, wantFindTokenCalls: 1, wantActivateCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := &fakeUserServiceUserRepo{findByActivationTokenResult: tc.findTokenResult, findByActivationTokenErr: tc.findTokenErr, activateWithPasswordErr: tc.activateErr}
			svc := &UserService{userRepo: userRepo, resetTokenRepo: &fakeUserServiceResetTokenRepo{}}

			err := svc.ActivateUserWithToken(context.Background(), tc.token, tc.password)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("ActivateUserWithToken() error = %v", err)
			}
			if userRepo.findByActivationTokenCalls != tc.wantFindTokenCalls || userRepo.activateWithPasswordCalls != tc.wantActivateCalls {
				t.Fatalf("find/activate calls = %d/%d, want %d/%d", userRepo.findByActivationTokenCalls, userRepo.activateWithPasswordCalls, tc.wantFindTokenCalls, tc.wantActivateCalls)
			}
		})
	}
}

func TestUserServiceResetPasswordWithToken(t *testing.T) {
	userID := uuid.New()
	tokenID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name               string
		email              string
		plainToken         string
		newPassword        string
		findEmailErr       error
		findByHashResult   *model.ResetToken
		findByHashErr      error
		markUsedErr        error
		updatePasswordErr  error
		wantErr            error
		wantFindEmailCalls int
		wantFindHashCalls  int
		wantMarkUsedCalls  int
		wantUpdateCalls    int
	}{
		{name: "missing email", email: "", plainToken: "token", newPassword: "new", wantErr: ErrResetTokenInvalid},
		{name: "missing token", email: "u@example.com", plainToken: "", newPassword: "new", wantErr: ErrResetTokenInvalid},
		{name: "missing password", email: "u@example.com", plainToken: "token", newPassword: "", wantErr: ErrPasswordCannotBeEmpty},
		{name: "email not found", email: "u@example.com", plainToken: "token", newPassword: "new", findEmailErr: sentinelErr, wantErr: ErrResetTokenInvalid, wantFindEmailCalls: 1},
		{name: "token hash not found", email: "u@example.com", plainToken: "token", newPassword: "new", findByHashErr: sentinelErr, wantErr: ErrResetTokenInvalid, wantFindEmailCalls: 1, wantFindHashCalls: 1},
		{name: "token user mismatch", email: "u@example.com", plainToken: "token", newPassword: "new", findByHashResult: &model.ResetToken{ID: tokenID, UserID: uuid.New(), ExpiresAt: time.Now().Add(time.Minute)}, wantErr: ErrResetTokenInvalid, wantFindEmailCalls: 1, wantFindHashCalls: 1},
		{name: "token expired", email: "u@example.com", plainToken: "token", newPassword: "new", findByHashResult: &model.ResetToken{ID: tokenID, UserID: userID, ExpiresAt: time.Now().Add(-time.Minute)}, wantErr: ErrResetTokenInvalid, wantFindEmailCalls: 1, wantFindHashCalls: 1},
		{name: "mark used no rows", email: "u@example.com", plainToken: "token", newPassword: "new", findByHashResult: &model.ResetToken{ID: tokenID, UserID: userID, ExpiresAt: time.Now().Add(time.Minute)}, markUsedErr: repo.ErrNoRowsUpdated, wantErr: ErrResetTokenInvalid, wantFindEmailCalls: 1, wantFindHashCalls: 1, wantMarkUsedCalls: 1},
		{name: "mark used generic error", email: "u@example.com", plainToken: "token", newPassword: "new", findByHashResult: &model.ResetToken{ID: tokenID, UserID: userID, ExpiresAt: time.Now().Add(time.Minute)}, markUsedErr: sentinelErr, wantErr: ErrResetTokenInvalid, wantFindEmailCalls: 1, wantFindHashCalls: 1, wantMarkUsedCalls: 1},
		{name: "update password failure", email: "u@example.com", plainToken: "token", newPassword: "new", findByHashResult: &model.ResetToken{ID: tokenID, UserID: userID, ExpiresAt: time.Now().Add(time.Minute)}, updatePasswordErr: sentinelErr, wantErr: ErrFailedToUpdatePassword, wantFindEmailCalls: 1, wantFindHashCalls: 1, wantMarkUsedCalls: 1, wantUpdateCalls: 1},
		{name: "success", email: "u@example.com", plainToken: "token", newPassword: "new", findByHashResult: &model.ResetToken{ID: tokenID, UserID: userID, ExpiresAt: time.Now().Add(time.Minute)}, wantFindEmailCalls: 1, wantFindHashCalls: 1, wantMarkUsedCalls: 1, wantUpdateCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := &fakeUserServiceUserRepo{findByEmailResult: &model.User{UserID: userID, Email: tc.email}, findByEmailErr: tc.findEmailErr, updatePasswordErr: tc.updatePasswordErr}
			resetRepo := &fakeUserServiceResetTokenRepo{findByHashRes: tc.findByHashResult, findByHashErr: tc.findByHashErr, markUsedErr: tc.markUsedErr}
			svc := &UserService{userRepo: userRepo, resetTokenRepo: resetRepo}

			err := svc.ResetPasswordWithToken(context.Background(), tc.email, tc.plainToken, tc.newPassword)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("ResetPasswordWithToken() error = %v", err)
			}
			if userRepo.findByEmailCalls != tc.wantFindEmailCalls || resetRepo.findByHashCalls != tc.wantFindHashCalls || resetRepo.markUsedCalls != tc.wantMarkUsedCalls || userRepo.updatePasswordCalls != tc.wantUpdateCalls {
				t.Fatalf("calls email/hash/mark/update = %d/%d/%d/%d, want %d/%d/%d/%d", userRepo.findByEmailCalls, resetRepo.findByHashCalls, resetRepo.markUsedCalls, userRepo.updatePasswordCalls, tc.wantFindEmailCalls, tc.wantFindHashCalls, tc.wantMarkUsedCalls, tc.wantUpdateCalls)
			}
		})
	}
}

func TestUserServiceRequestPasswordReset(t *testing.T) {
	userID := uuid.New()
	sentinelErr := errors.New("repo failed")

	tests := []struct {
		name            string
		email           string
		findByEmailErr  error
		createTokenErr  error
		wantErr         error
		wantFindCalls   int
		wantCreateCalls int
	}{
		{name: "empty email", email: "", wantErr: ErrEmailCannotBeEmpty},
		{name: "non-existent email should not leak", email: "u@example.com", findByEmailErr: sentinelErr, wantFindCalls: 1},
		{name: "create token failure swallowed", email: "u@example.com", createTokenErr: sentinelErr, wantFindCalls: 1, wantCreateCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := &fakeUserServiceUserRepo{findByEmailResult: &model.User{UserID: userID, Email: tc.email}, findByEmailErr: tc.findByEmailErr}
			resetRepo := &fakeUserServiceResetTokenRepo{createErr: tc.createTokenErr}
			svc := &UserService{userRepo: userRepo, resetTokenRepo: resetRepo, frontendURL: "http://localhost:3000"}

			err := svc.RequestPasswordReset(context.Background(), tc.email)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("RequestPasswordReset() error = %v", err)
			}
			if userRepo.findByEmailCalls != tc.wantFindCalls || resetRepo.createCalls != tc.wantCreateCalls {
				t.Fatalf("find/create calls = %d/%d, want %d/%d", userRepo.findByEmailCalls, resetRepo.createCalls, tc.wantFindCalls, tc.wantCreateCalls)
			}
		})
	}
}

func TestUserServiceLockUnlock(t *testing.T) {
	schoolID := uuid.New()
	userID := uuid.New()
	sentinelErr := errors.New("repo failed")

	t.Run("lock school access denied", func(t *testing.T) {
		userRepo := &fakeUserServiceUserRepo{isUserInSchoolResult: false}
		svc := &UserService{userRepo: userRepo, resetTokenRepo: &fakeUserServiceResetTokenRepo{}}
		err := svc.Lock(context.Background(), &schoolID, userID)
		if !errors.Is(err, ErrSchoolAccessDenied) {
			t.Fatalf("error = %v, want %v", err, ErrSchoolAccessDenied)
		}
		if userRepo.lockCalls != 0 {
			t.Fatalf("lock should not be called when access denied")
		}
	})

	t.Run("lock scope check error", func(t *testing.T) {
		userRepo := &fakeUserServiceUserRepo{isUserInSchoolErr: sentinelErr}
		svc := &UserService{userRepo: userRepo, resetTokenRepo: &fakeUserServiceResetTokenRepo{}}
		err := svc.Lock(context.Background(), &schoolID, userID)
		if !errors.Is(err, sentinelErr) {
			t.Fatalf("error = %v, want %v", err, sentinelErr)
		}
	})

	t.Run("lock success for super admin", func(t *testing.T) {
		userRepo := &fakeUserServiceUserRepo{}
		svc := &UserService{userRepo: userRepo, resetTokenRepo: &fakeUserServiceResetTokenRepo{}}
		err := svc.Lock(context.Background(), nil, userID)
		if err != nil {
			t.Fatalf("Lock() error = %v", err)
		}
		if userRepo.lockCalls != 1 {
			t.Fatalf("lock calls = %d, want 1", userRepo.lockCalls)
		}
	})

	t.Run("unlock school access denied", func(t *testing.T) {
		userRepo := &fakeUserServiceUserRepo{isUserInSchoolResult: false}
		svc := &UserService{userRepo: userRepo, resetTokenRepo: &fakeUserServiceResetTokenRepo{}}
		err := svc.Unlock(context.Background(), &schoolID, userID)
		if !errors.Is(err, ErrSchoolAccessDenied) {
			t.Fatalf("error = %v, want %v", err, ErrSchoolAccessDenied)
		}
		if userRepo.unlockCalls != 0 {
			t.Fatalf("unlock should not be called when access denied")
		}
	})

	t.Run("unlock success for school admin", func(t *testing.T) {
		userRepo := &fakeUserServiceUserRepo{isUserInSchoolResult: true}
		svc := &UserService{userRepo: userRepo, resetTokenRepo: &fakeUserServiceResetTokenRepo{}}
		err := svc.Unlock(context.Background(), &schoolID, userID)
		if err != nil {
			t.Fatalf("Unlock() error = %v", err)
		}
		if userRepo.isUserInSchoolCalls != 1 || userRepo.unlockCalls != 1 {
			t.Fatalf("scope/unlock calls = %d/%d, want 1/1", userRepo.isUserInSchoolCalls, userRepo.unlockCalls)
		}
	})
}

func mustUUID(t *testing.T, value string) uuid.UUID {
	t.Helper()
	id, err := uuid.Parse(value)
	if err != nil {
		t.Fatalf("invalid uuid in test: %v", err)
	}
	return id
}
