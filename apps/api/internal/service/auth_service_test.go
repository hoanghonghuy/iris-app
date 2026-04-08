package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type fakeAuthServiceUserRepo struct {
	findByEmailCalls int
	findByEmailArg   string
	findByEmailRes   *model.User
	findByEmailErr   error

	findByGoogleSubCalls int
	findByGoogleSubArg   string
	findByGoogleSubRes   *model.User
	findByGoogleSubErr   error

	linkGoogleSubCalls int
	linkGoogleSubUser  uuid.UUID
	linkGoogleSubSub   string
	linkGoogleSubErr   error

	rolesOfUserCalls int
	rolesOfUserArg   uuid.UUID
	rolesOfUserRes   []string
	rolesOfUserErr   error

	findByIDCalls int
	findByIDArg   uuid.UUID
	findByIDRes   *model.UserInfo
	findByIDErr   error
}

func (f *fakeAuthServiceUserRepo) FindByEmail(_ context.Context, email string) (*model.User, error) {
	f.findByEmailCalls++
	f.findByEmailArg = email
	return f.findByEmailRes, f.findByEmailErr
}

func (f *fakeAuthServiceUserRepo) FindByGoogleSub(_ context.Context, googleSub string) (*model.User, error) {
	f.findByGoogleSubCalls++
	f.findByGoogleSubArg = googleSub
	return f.findByGoogleSubRes, f.findByGoogleSubErr
}

func (f *fakeAuthServiceUserRepo) LinkGoogleSub(_ context.Context, userID uuid.UUID, googleSub string) error {
	f.linkGoogleSubCalls++
	f.linkGoogleSubUser = userID
	f.linkGoogleSubSub = googleSub
	return f.linkGoogleSubErr
}

func (f *fakeAuthServiceUserRepo) RolesOfUser(_ context.Context, userID uuid.UUID) ([]string, error) {
	f.rolesOfUserCalls++
	f.rolesOfUserArg = userID
	return f.rolesOfUserRes, f.rolesOfUserErr
}

func (f *fakeAuthServiceUserRepo) FindByID(_ context.Context, userID uuid.UUID) (*model.UserInfo, error) {
	f.findByIDCalls++
	f.findByIDArg = userID
	return f.findByIDRes, f.findByIDErr
}

type fakeAuthServiceSchoolAdminRepo struct {
	getByUserIDCalls int
	getByUserIDArg   uuid.UUID
	getByUserIDRes   *model.SchoolAdmin
	getByUserIDErr   error
}

func (f *fakeAuthServiceSchoolAdminRepo) GetByUserID(_ context.Context, userID uuid.UUID) (*model.SchoolAdmin, error) {
	f.getByUserIDCalls++
	f.getByUserIDArg = userID
	return f.getByUserIDRes, f.getByUserIDErr
}

type fakeAuthServiceGoogleVerifier struct {
	verifyCalls int
	verifyToken string
	verifyRes   *auth.GoogleIdentity
	verifyErr   error
}

func (f *fakeAuthServiceGoogleVerifier) Verify(_ context.Context, rawIDToken string) (*auth.GoogleIdentity, error) {
	f.verifyCalls++
	f.verifyToken = rawIDToken
	if f.verifyErr != nil {
		return nil, f.verifyErr
	}
	return f.verifyRes, nil
}

func mustHashPassword(t *testing.T, plain string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}
	return string(hash)
}

func TestAuthServiceLogin(t *testing.T) {
	userID := uuid.New()
	schoolID := uuid.New()
	sentinelErr := errors.New("repo failed")
	pwdHash := mustHashPassword(t, "Pass1234")

	tests := []struct {
		name               string
		findByEmailRes     *model.User
		findByEmailErr     error
		rolesRes           []string
		rolesErr           error
		schoolAdminRes     *model.SchoolAdmin
		schoolAdminErr     error
		password           string
		wantErr            error
		wantFindCalls      int
		wantRolesCalls     int
		wantSchoolGetCalls int
		wantTokenSchoolID  string
	}{
		{name: "find user error passthrough", findByEmailErr: sentinelErr, password: "Pass1234", wantErr: sentinelErr, wantFindCalls: 1},
		{name: "locked user", findByEmailRes: &model.User{UserID: userID, Email: "u@example.com", PasswordHash: pwdHash, Status: "locked"}, password: "Pass1234", wantErr: auth.ErrUserLocked, wantFindCalls: 1},
		{name: "invalid password", findByEmailRes: &model.User{UserID: userID, Email: "u@example.com", PasswordHash: pwdHash, Status: "active"}, password: "WrongPass", wantErr: auth.ErrInvalidCredentials, wantFindCalls: 1},
		{name: "roles query error", findByEmailRes: &model.User{UserID: userID, Email: "u@example.com", PasswordHash: pwdHash, Status: "active"}, password: "Pass1234", rolesErr: sentinelErr, wantErr: sentinelErr, wantFindCalls: 1, wantRolesCalls: 1},
		{name: "school admin role uses school id", findByEmailRes: &model.User{UserID: userID, Email: "u@example.com", PasswordHash: pwdHash, Status: "active"}, password: "Pass1234", rolesRes: []string{"SCHOOL_ADMIN"}, schoolAdminRes: &model.SchoolAdmin{UserID: userID, SchoolID: schoolID}, wantFindCalls: 1, wantRolesCalls: 1, wantSchoolGetCalls: 1, wantTokenSchoolID: schoolID.String()},
		{name: "school admin role get school failed keeps token without school", findByEmailRes: &model.User{UserID: userID, Email: "u@example.com", PasswordHash: pwdHash, Status: "active"}, password: "Pass1234", rolesRes: []string{"SCHOOL_ADMIN"}, schoolAdminErr: sentinelErr, wantFindCalls: 1, wantRolesCalls: 1, wantSchoolGetCalls: 1, wantTokenSchoolID: ""},
		{name: "success non school admin", findByEmailRes: &model.User{UserID: userID, Email: "u@example.com", PasswordHash: pwdHash, Status: "active"}, password: "Pass1234", rolesRes: []string{"PARENT"}, wantFindCalls: 1, wantRolesCalls: 1, wantSchoolGetCalls: 0, wantTokenSchoolID: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := &fakeAuthServiceUserRepo{findByEmailRes: tc.findByEmailRes, findByEmailErr: tc.findByEmailErr, rolesOfUserRes: tc.rolesRes, rolesOfUserErr: tc.rolesErr}
			schoolAdminRepo := &fakeAuthServiceSchoolAdminRepo{getByUserIDRes: tc.schoolAdminRes, getByUserIDErr: tc.schoolAdminErr}
			svc := &AuthService{userRepo: userRepo, schoolAdminRepo: schoolAdminRepo, jwtAuth: &auth.Authenticator{Secret: "test-secret", TTLSeconds: 3600}}

			resp, err := svc.Login(context.Background(), "u@example.com", tc.password)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("Login() error = %v", err)
			}

			if userRepo.findByEmailCalls != tc.wantFindCalls || userRepo.rolesOfUserCalls != tc.wantRolesCalls || schoolAdminRepo.getByUserIDCalls != tc.wantSchoolGetCalls {
				t.Fatalf("calls find/roles/school = %d/%d/%d, want %d/%d/%d", userRepo.findByEmailCalls, userRepo.rolesOfUserCalls, schoolAdminRepo.getByUserIDCalls, tc.wantFindCalls, tc.wantRolesCalls, tc.wantSchoolGetCalls)
			}

			if tc.wantErr == nil {
				if resp == nil || resp.AccessToken == "" || resp.TokenType != "Bearer" || resp.ExpiresIn != 3600 {
					t.Fatalf("unexpected login response = %#v", resp)
				}
				claims, parseErr := auth.Parse("test-secret", resp.AccessToken)
				if parseErr != nil {
					t.Fatalf("Parse token error = %v", parseErr)
				}
				if claims.Email != "u@example.com" {
					t.Fatalf("claims email = %q, want %q", claims.Email, "u@example.com")
				}
				if claims.SchoolID != tc.wantTokenSchoolID {
					t.Fatalf("claims schoolID = %q, want %q", claims.SchoolID, tc.wantTokenSchoolID)
				}
			}
		})
	}
}

func TestAuthServiceLoginWithGoogleToken(t *testing.T) {
	userID := uuid.New()
	sentinelErr := errors.New("repo failed")
	pwdHash := mustHashPassword(t, "Pass1234")
	identity := &auth.GoogleIdentity{Sub: "google-sub", Email: "g@example.com", EmailVerified: true, HostedDomain: "iris.edu.vn", Name: "Google User"}

	tests := []struct {
		name                string
		googleEnabled       bool
		googleVerifier      auth.GoogleTokenVerifier
		googleHD            string
		linkedUserRes       *model.User
		linkedUserErr       error
		emailUserRes        *model.User
		emailUserErr        error
		rolesRes            []string
		rolesErr            error
		password            string
		linkErr             error
		wantErr             error
		wantFindGoogleCalls int
		wantFindEmailCalls  int
		wantLinkCalls       int
	}{
		{name: "google disabled", googleEnabled: false, googleVerifier: &fakeAuthServiceGoogleVerifier{verifyRes: identity}, wantErr: ErrGoogleLoginDisabled},
		{name: "missing verifier", googleEnabled: true, googleVerifier: nil, wantErr: ErrGoogleLoginDisabled},
		{name: "verify error maps invalid credentials", googleEnabled: true, googleVerifier: &fakeAuthServiceGoogleVerifier{verifyErr: sentinelErr}, wantErr: auth.ErrInvalidCredentials},
		{name: "hosted domain mismatch", googleEnabled: true, googleHD: "iris.edu.vn", googleVerifier: &fakeAuthServiceGoogleVerifier{verifyRes: &auth.GoogleIdentity{Sub: "google-sub", Email: "g@example.com", EmailVerified: true, HostedDomain: "other.edu.vn"}}, wantErr: ErrGoogleDomainNotAllowed},
		{name: "linked user locked", googleEnabled: true, googleVerifier: &fakeAuthServiceGoogleVerifier{verifyRes: identity}, linkedUserRes: &model.User{UserID: userID, Email: "g@example.com", PasswordHash: pwdHash, Status: "locked", GoogleSub: "google-sub"}, wantErr: auth.ErrUserLocked, wantFindGoogleCalls: 1},
		{name: "linked user success", googleEnabled: true, googleVerifier: &fakeAuthServiceGoogleVerifier{verifyRes: identity}, linkedUserRes: &model.User{UserID: userID, Email: "g@example.com", PasswordHash: pwdHash, Status: "active", GoogleSub: "google-sub"}, linkedUserErr: nil, rolesRes: []string{"PARENT"}, wantFindGoogleCalls: 1},
		{name: "linked lookup non-no-rows error", googleEnabled: true, googleVerifier: &fakeAuthServiceGoogleVerifier{verifyRes: identity}, linkedUserErr: sentinelErr, wantErr: sentinelErr, wantFindGoogleCalls: 1},
		{name: "not provisioned", googleEnabled: true, googleVerifier: &fakeAuthServiceGoogleVerifier{verifyRes: identity}, linkedUserErr: pgx.ErrNoRows, emailUserErr: pgx.ErrNoRows, wantErr: ErrGoogleAccountNotProvisioned, wantFindGoogleCalls: 1, wantFindEmailCalls: 1},
		{name: "email lookup generic error", googleEnabled: true, googleVerifier: &fakeAuthServiceGoogleVerifier{verifyRes: identity}, linkedUserErr: pgx.ErrNoRows, emailUserErr: sentinelErr, wantErr: sentinelErr, wantFindGoogleCalls: 1, wantFindEmailCalls: 1},
		{name: "email user locked", googleEnabled: true, googleVerifier: &fakeAuthServiceGoogleVerifier{verifyRes: identity}, linkedUserErr: pgx.ErrNoRows, emailUserRes: &model.User{UserID: userID, Email: "g@example.com", PasswordHash: pwdHash, Status: "locked", GoogleSub: ""}, wantErr: auth.ErrUserLocked, wantFindGoogleCalls: 1, wantFindEmailCalls: 1},
		{name: "link requires password", googleEnabled: true, googleVerifier: &fakeAuthServiceGoogleVerifier{verifyRes: identity}, linkedUserErr: pgx.ErrNoRows, emailUserRes: &model.User{UserID: userID, Email: "g@example.com", PasswordHash: pwdHash, Status: "active", GoogleSub: ""}, password: "", wantErr: ErrGoogleLinkPasswordRequired, wantFindGoogleCalls: 1, wantFindEmailCalls: 1},
		{name: "link wrong password", googleEnabled: true, googleVerifier: &fakeAuthServiceGoogleVerifier{verifyRes: identity}, linkedUserErr: pgx.ErrNoRows, emailUserRes: &model.User{UserID: userID, Email: "g@example.com", PasswordHash: pwdHash, Status: "active", GoogleSub: ""}, password: "WrongPass", wantErr: auth.ErrInvalidCredentials, wantFindGoogleCalls: 1, wantFindEmailCalls: 1},
		{name: "link call error", googleEnabled: true, googleVerifier: &fakeAuthServiceGoogleVerifier{verifyRes: identity}, linkedUserErr: pgx.ErrNoRows, emailUserRes: &model.User{UserID: userID, Email: "g@example.com", PasswordHash: pwdHash, Status: "active", GoogleSub: ""}, password: "Pass1234", linkErr: sentinelErr, wantErr: sentinelErr, wantFindGoogleCalls: 1, wantFindEmailCalls: 1, wantLinkCalls: 1},
		{name: "already linked local user success without password", googleEnabled: true, googleVerifier: &fakeAuthServiceGoogleVerifier{verifyRes: identity}, linkedUserErr: pgx.ErrNoRows, emailUserRes: &model.User{UserID: userID, Email: "g@example.com", PasswordHash: pwdHash, Status: "active", GoogleSub: "existing-sub"}, password: "", rolesRes: []string{"TEACHER"}, wantFindGoogleCalls: 1, wantFindEmailCalls: 1, wantLinkCalls: 0},
		{name: "link with correct password success", googleEnabled: true, googleVerifier: &fakeAuthServiceGoogleVerifier{verifyRes: identity}, linkedUserErr: pgx.ErrNoRows, emailUserRes: &model.User{UserID: userID, Email: "g@example.com", PasswordHash: pwdHash, Status: "active", GoogleSub: ""}, password: "Pass1234", rolesRes: []string{"PARENT"}, wantFindGoogleCalls: 1, wantFindEmailCalls: 1, wantLinkCalls: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := &fakeAuthServiceUserRepo{
				findByGoogleSubRes: tc.linkedUserRes,
				findByGoogleSubErr: tc.linkedUserErr,
				findByEmailRes:     tc.emailUserRes,
				findByEmailErr:     tc.emailUserErr,
				rolesOfUserRes:     tc.rolesRes,
				rolesOfUserErr:     tc.rolesErr,
				linkGoogleSubErr:   tc.linkErr,
			}
			svc := &AuthService{
				userRepo:        userRepo,
				schoolAdminRepo: &fakeAuthServiceSchoolAdminRepo{},
				jwtAuth:         &auth.Authenticator{Secret: "test-secret", TTLSeconds: 3600},
				googleVerifier:  tc.googleVerifier,
				googleEnabled:   tc.googleEnabled,
				googleHD:        tc.googleHD,
			}

			resp, err := svc.LoginWithGoogleToken(context.Background(), "google-token", tc.password)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("error = %v, want %v", err, tc.wantErr)
				}
			} else if err != nil {
				t.Fatalf("LoginWithGoogleToken() error = %v", err)
			}

			if userRepo.findByGoogleSubCalls != tc.wantFindGoogleCalls || userRepo.findByEmailCalls != tc.wantFindEmailCalls || userRepo.linkGoogleSubCalls != tc.wantLinkCalls {
				t.Fatalf("calls findGoogle/findEmail/link = %d/%d/%d, want %d/%d/%d", userRepo.findByGoogleSubCalls, userRepo.findByEmailCalls, userRepo.linkGoogleSubCalls, tc.wantFindGoogleCalls, tc.wantFindEmailCalls, tc.wantLinkCalls)
			}

			if tc.wantLinkCalls > 0 {
				if userRepo.linkGoogleSubSub != "google-sub" {
					t.Fatalf("linked sub = %q, want %q", userRepo.linkGoogleSubSub, "google-sub")
				}
			}

			if tc.wantErr == nil {
				if resp == nil || resp.AccessToken == "" || resp.TokenType != "Bearer" || resp.ExpiresIn != 3600 {
					t.Fatalf("unexpected login response = %#v", resp)
				}
			}
		})
	}
}

func TestAuthServiceGetUserInfo(t *testing.T) {
	userID := uuid.New()
	sentinelErr := errors.New("repo failed")

	t.Run("repo error passthrough", func(t *testing.T) {
		userRepo := &fakeAuthServiceUserRepo{findByIDErr: sentinelErr}
		svc := &AuthService{userRepo: userRepo, schoolAdminRepo: &fakeAuthServiceSchoolAdminRepo{}}

		_, err := svc.GetUserInfo(context.Background(), userID)
		if !errors.Is(err, sentinelErr) {
			t.Fatalf("error = %v, want %v", err, sentinelErr)
		}
	})

	t.Run("success", func(t *testing.T) {
		userInfo := &model.UserInfo{UserID: userID, Email: "u@example.com", Roles: []string{"PARENT"}}
		userRepo := &fakeAuthServiceUserRepo{findByIDRes: userInfo}
		svc := &AuthService{userRepo: userRepo, schoolAdminRepo: &fakeAuthServiceSchoolAdminRepo{}}

		got, err := svc.GetUserInfo(context.Background(), userID)
		if err != nil {
			t.Fatalf("GetUserInfo() error = %v", err)
		}
		if got == nil || got.UserID != userID || got.Email != "u@example.com" {
			t.Fatalf("unexpected user info = %#v", got)
		}
	})
}
