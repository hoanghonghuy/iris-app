package authhandlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type fakeAuthService struct {
	loginFn                func(context.Context, string, string) (*service.LoginResponse, error)
	loginWithGoogleTokenFn func(context.Context, string, string) (*service.LoginResponse, error)
}

func (f *fakeAuthService) Login(ctx context.Context, email, password string) (*service.LoginResponse, error) {
	if f.loginFn == nil {
		return nil, errors.New("unexpected Login call")
	}
	return f.loginFn(ctx, email, password)
}

func (f *fakeAuthService) LoginWithGoogleToken(ctx context.Context, googleIDToken, password string) (*service.LoginResponse, error) {
	if f.loginWithGoogleTokenFn == nil {
		return nil, errors.New("unexpected LoginWithGoogleToken call")
	}
	return f.loginWithGoogleTokenFn(ctx, googleIDToken, password)
}

type fakeUserService struct {
	requestPasswordResetFn   func(context.Context, string) error
	resetPasswordWithTokenFn func(context.Context, string, string, string) error
	findByIDFn               func(context.Context, *uuid.UUID, uuid.UUID) (*model.UserInfo, error)
}

func (f *fakeUserService) RequestPasswordReset(ctx context.Context, email string) error {
	if f.requestPasswordResetFn == nil {
		return errors.New("unexpected RequestPasswordReset call")
	}
	return f.requestPasswordResetFn(ctx, email)
}

func (f *fakeUserService) ResetPasswordWithToken(ctx context.Context, email, plainToken, newPassword string) error {
	if f.resetPasswordWithTokenFn == nil {
		return errors.New("unexpected ResetPasswordWithToken call")
	}
	return f.resetPasswordWithTokenFn(ctx, email, plainToken, newPassword)
}

func (f *fakeUserService) FindByID(ctx context.Context, adminSchoolID *uuid.UUID, userID uuid.UUID) (*model.UserInfo, error) {
	if f.findByIDFn == nil {
		return nil, errors.New("unexpected FindByID call")
	}
	return f.findByIDFn(ctx, adminSchoolID, userID)
}

func decodeAuthError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if v, ok := body["error"].(string); ok {
		return v
	}
	return ""
}

func decodeAuthData(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	data, _ := body["data"].(map[string]any)
	return data
}

func TestLogin_InvalidRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &AuthHandler{}
	r := gin.New()
	r.POST("/auth/login", h.Login)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader("{bad-json"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestLogin_ErrorMappingAndSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		err       error
		wantCode  int
		wantError string
	}{
		{name: "no rows", err: pgx.ErrNoRows, wantCode: http.StatusUnauthorized, wantError: "invalid credentials"},
		{name: "invalid credentials", err: auth.ErrInvalidCredentials, wantCode: http.StatusUnauthorized, wantError: "email or password incorrect"},
		{name: "user locked", err: auth.ErrUserLocked, wantCode: http.StatusForbidden, wantError: "user account locked"},
		{name: "internal", err: errors.New("boom"), wantCode: http.StatusInternalServerError, wantError: "server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandler{authService: &fakeAuthService{loginFn: func(context.Context, string, string) (*service.LoginResponse, error) {
				return nil, tt.err
			}}}
			r := gin.New()
			r.POST("/auth/login", h.Login)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"user@example.com","password":"secret123"}`))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantCode)
			}
			if got := decodeAuthError(t, rec); got != tt.wantError {
				t.Fatalf("error = %q, want %q", got, tt.wantError)
			}
		})
	}

	t.Run("success", func(t *testing.T) {
		h := &AuthHandler{authService: &fakeAuthService{loginFn: func(_ context.Context, email, password string) (*service.LoginResponse, error) {
			if email != "user@example.com" || password != "secret123" {
				t.Fatalf("unexpected credentials forwarded")
			}
			return &service.LoginResponse{AccessToken: "token-abc", TokenType: "Bearer", ExpiresIn: 3600}, nil
		}}}
		r := gin.New()
		r.POST("/auth/login", h.Login)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"user@example.com","password":"secret123"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		data := decodeAuthData(t, rec)
		if data["access_token"] != "token-abc" {
			t.Fatalf("access_token = %v, want %v", data["access_token"], "token-abc")
		}
	})
}

func TestLoginWithGoogle_ErrorMappingAndSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		err           error
		wantCode      int
		wantError     string
		wantErrorCode string
	}{
		{name: "google disabled", err: service.ErrGoogleLoginDisabled, wantCode: http.StatusForbidden, wantError: "google login is disabled"},
		{name: "domain not allowed", err: service.ErrGoogleDomainNotAllowed, wantCode: http.StatusForbidden, wantError: "google account domain not allowed"},
		{name: "not provisioned", err: service.ErrGoogleAccountNotProvisioned, wantCode: http.StatusUnauthorized, wantError: "google account is not provisioned"},
		{name: "link password required", err: service.ErrGoogleLinkPasswordRequired, wantCode: http.StatusForbidden, wantError: "password confirmation required to link google account", wantErrorCode: "GOOGLE_LINK_PASSWORD_REQUIRED"},
		{name: "invalid credentials", err: auth.ErrInvalidCredentials, wantCode: http.StatusUnauthorized, wantError: "invalid credentials"},
		{name: "user locked", err: auth.ErrUserLocked, wantCode: http.StatusForbidden, wantError: "user account locked"},
		{name: "internal", err: errors.New("boom"), wantCode: http.StatusInternalServerError, wantError: "server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandler{authService: &fakeAuthService{loginWithGoogleTokenFn: func(context.Context, string, string) (*service.LoginResponse, error) {
				return nil, tt.err
			}}}
			r := gin.New()
			r.POST("/auth/login/google", h.LoginWithGoogle)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/login/google", strings.NewReader(`{"id_token":"id-token","password":"secret123"}`))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantCode)
			}

			var body map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if got, _ := body["error"].(string); got != tt.wantError {
				t.Fatalf("error = %q, want %q", got, tt.wantError)
			}
			if tt.wantErrorCode != "" {
				if got, _ := body["error_code"].(string); got != tt.wantErrorCode {
					t.Fatalf("error_code = %q, want %q", got, tt.wantErrorCode)
				}
			}
		})
	}

	t.Run("success", func(t *testing.T) {
		h := &AuthHandler{authService: &fakeAuthService{loginWithGoogleTokenFn: func(_ context.Context, idToken, password string) (*service.LoginResponse, error) {
			if idToken != "id-token" || password != "secret123" {
				t.Fatalf("unexpected google login params")
			}
			return &service.LoginResponse{AccessToken: "google-token", TokenType: "Bearer", ExpiresIn: 3600}, nil
		}}}
		r := gin.New()
		r.POST("/auth/login/google", h.LoginWithGoogle)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/login/google", strings.NewReader(`{"id_token":"id-token","password":"secret123"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		data := decodeAuthData(t, rec)
		if data["access_token"] != "google-token" {
			t.Fatalf("access_token = %v, want %v", data["access_token"], "google-token")
		}
	})
}

func TestForgotPassword_Behavior(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("invalid request body", func(t *testing.T) {
		h := &AuthHandler{}
		r := gin.New()
		r.POST("/auth/forgot-password", h.ForgotPassword)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/forgot-password", strings.NewReader(`{"email":"bad-email"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("always returns ok even if service fails", func(t *testing.T) {
		called := false
		h := &AuthHandler{userService: &fakeUserService{requestPasswordResetFn: func(_ context.Context, email string) error {
			called = true
			if email != "user@example.com" {
				t.Fatalf("unexpected email forwarded")
			}
			return errors.New("smtp down")
		}}}
		r := gin.New()
		r.POST("/auth/forgot-password", h.ForgotPassword)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/forgot-password", strings.NewReader(`{"email":"user@example.com"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if !called {
			t.Fatalf("expected RequestPasswordReset to be called")
		}
	})
}

func TestResetPassword_ErrorMappingAndSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("invalid request body", func(t *testing.T) {
		h := &AuthHandler{}
		r := gin.New()
		r.POST("/auth/reset-password", h.ResetPassword)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/reset-password", strings.NewReader(`{"email":"invalid"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	tests := []struct {
		name      string
		err       error
		wantCode  int
		wantError string
	}{
		{name: "invalid token", err: service.ErrResetTokenInvalid, wantCode: http.StatusBadRequest, wantError: "token không hợp lệ hoặc đã hết hạn"},
		{name: "empty password", err: service.ErrPasswordCannotBeEmpty, wantCode: http.StatusBadRequest, wantError: "mật khẩu không được để trống"},
		{name: "internal", err: errors.New("boom"), wantCode: http.StatusInternalServerError, wantError: "server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &AuthHandler{userService: &fakeUserService{resetPasswordWithTokenFn: func(context.Context, string, string, string) error {
				return tt.err
			}}}
			r := gin.New()
			r.POST("/auth/reset-password", h.ResetPassword)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/reset-password", strings.NewReader(`{"email":"user@example.com","token":"abc","password":"secret123"}`))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantCode {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantCode)
			}
			if got := decodeAuthError(t, rec); got != tt.wantError {
				t.Fatalf("error = %q, want %q", got, tt.wantError)
			}
		})
	}

	t.Run("success", func(t *testing.T) {
		called := false
		h := &AuthHandler{userService: &fakeUserService{resetPasswordWithTokenFn: func(_ context.Context, email, token, password string) error {
			called = true
			if email != "user@example.com" || token != "abc" || password != "secret123" {
				t.Fatalf("unexpected reset password params")
			}
			return nil
		}}}
		r := gin.New()
		r.POST("/auth/reset-password", h.ResetPassword)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/auth/reset-password", strings.NewReader(`{"email":"user@example.com","token":"abc","password":"secret123"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if !called {
			t.Fatalf("expected ResetPasswordWithToken to be called")
		}
	})
}

func TestMe_ClaimsAndProfileMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("missing claims", func(t *testing.T) {
		h := &AuthHandler{}
		r := gin.New()
		r.GET("/me", h.Me)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("returns claims with school_id and full_name when available", func(t *testing.T) {
		userID := uuid.New()
		h := &AuthHandler{userService: &fakeUserService{findByIDFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotUserID uuid.UUID) (*model.UserInfo, error) {
			if adminSchoolID != nil {
				t.Fatalf("expected adminSchoolID to be nil for Me")
			}
			if gotUserID != userID {
				t.Fatalf("unexpected user id forwarded")
			}
			return &model.UserInfo{UserID: userID, FullName: "Test User"}, nil
		}}}
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set(middleware.CtxClaims, &auth.Claims{UserID: userID.String(), Email: "user@example.com", Roles: []string{"PARENT"}, SchoolID: "school-1"})
			c.Next()
		})
		r.GET("/me", h.Me)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		data := decodeAuthData(t, rec)
		if data["user_id"] != userID.String() {
			t.Fatalf("user_id = %v, want %v", data["user_id"], userID.String())
		}
		if data["full_name"] != "Test User" {
			t.Fatalf("full_name = %v, want %v", data["full_name"], "Test User")
		}
		if data["school_id"] != "school-1" {
			t.Fatalf("school_id = %v, want %v", data["school_id"], "school-1")
		}
	})

	t.Run("invalid claims user_id still returns claim payload", func(t *testing.T) {
		h := &AuthHandler{userService: &fakeUserService{findByIDFn: func(context.Context, *uuid.UUID, uuid.UUID) (*model.UserInfo, error) {
			return nil, errors.New("must not be called")
		}}}
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set(middleware.CtxClaims, &auth.Claims{UserID: "not-a-uuid", Email: "user@example.com", Roles: []string{"PARENT"}})
			c.Next()
		})
		r.GET("/me", h.Me)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		data := decodeAuthData(t, rec)
		if _, ok := data["full_name"]; ok {
			t.Fatalf("did not expect full_name when user_id is invalid")
		}
	})
}
