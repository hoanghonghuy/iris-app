package parentcodehandlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

func TestRegisterParentWithGoogle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("invalid body", func(t *testing.T) {
		h := &ParentCodeHandler{}
		r := gin.New()
		r.POST("/register/parent/google", h.RegisterParentWithGoogle)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/register/parent/google", strings.NewReader(`{bad-json`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &ParentCodeHandler{parentCodeService: &fakeParentCodeService{registerParentWithGoogleFn: func(_ context.Context, idToken, parentCode string) (*service.LoginResponse, error) {
			if idToken != "id-token" || parentCode != "PARENT001" {
				t.Fatalf("unexpected payload")
			}
			return &service.LoginResponse{AccessToken: "token", TokenType: "Bearer", ExpiresIn: 3600}, nil
		}}}
		r := gin.New()
		r.POST("/register/parent/google", h.RegisterParentWithGoogle)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/register/parent/google", strings.NewReader(`{"id_token":"id-token","parent_code":"PARENT001"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("error mappings", func(t *testing.T) {
		cases := []struct {
			name       string
			err        error
			wantStatus int
			wantError  string
		}{
			{name: "google disabled", err: service.ErrGoogleLoginDisabled, wantStatus: http.StatusForbidden, wantError: "tính năng đăng nhập Google đang tắt"},
			{name: "domain not allowed", err: service.ErrGoogleDomainNotAllowed, wantStatus: http.StatusForbidden, wantError: "tên miền Google không hợp lệ"},
			{name: "invalid credentials", err: auth.ErrInvalidCredentials, wantStatus: http.StatusUnauthorized, wantError: "xác thực Google thất bại"},
			{name: "invalid parent code", err: service.ErrInvalidParentCode, wantStatus: http.StatusBadRequest, wantError: "mã phụ huynh không hợp lệ/hết hạn/hết lượt dùng"},
			{name: "expired parent code", err: service.ErrParentCodeExpired, wantStatus: http.StatusBadRequest, wantError: "mã phụ huynh không hợp lệ/hết hạn/hết lượt dùng"},
			{name: "max usage", err: service.ErrParentCodeMaxUsageReached, wantStatus: http.StatusBadRequest, wantError: "mã phụ huynh không hợp lệ/hết hạn/hết lượt dùng"},
			{name: "email exists", err: service.ErrEmailAlreadyExists, wantStatus: http.StatusConflict, wantError: "Email này đã được đăng ký. Vui lòng quay lại trang Đăng nhập"},
			{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "server error"},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				h := &ParentCodeHandler{parentCodeService: &fakeParentCodeService{registerParentWithGoogleFn: func(context.Context, string, string) (*service.LoginResponse, error) {
					return nil, tc.err
				}}}
				r := gin.New()
				r.POST("/register/parent/google", h.RegisterParentWithGoogle)

				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodPost, "/register/parent/google", strings.NewReader(`{"id_token":"id-token","parent_code":"PARENT001"}`))
				req.Header.Set("Content-Type", "application/json")
				r.ServeHTTP(rec, req)

				if rec.Code != tc.wantStatus {
					t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
				}
				if got := decodeParentCodeError(t, rec); got != tc.wantError {
					t.Fatalf("error = %q, want %q", got, tc.wantError)
				}
			})
		}
	})
}

func TestRegisterParent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("invalid body", func(t *testing.T) {
		h := &ParentCodeHandler{}
		r := gin.New()
		r.POST("/register/parent", h.RegisterParent)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/register/parent", strings.NewReader(`{bad-json`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &ParentCodeHandler{parentCodeService: &fakeParentCodeService{registerParentFn: func(_ context.Context, email, password, parentCode string) (*service.LoginResponse, error) {
			if email != "parent@example.com" || password != "secret123" || parentCode != "PARENT001" {
				t.Fatalf("unexpected payload")
			}
			return &service.LoginResponse{AccessToken: "token", TokenType: "Bearer", ExpiresIn: 3600}, nil
		}}}
		r := gin.New()
		r.POST("/register/parent", h.RegisterParent)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/register/parent", strings.NewReader(`{"email":"parent@example.com","password":"secret123","parent_code":"PARENT001"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
		}
	})

	t.Run("error mappings", func(t *testing.T) {
		cases := []struct {
			name       string
			err        error
			wantStatus int
			wantError  string
		}{
			{name: "invalid parent code", err: service.ErrInvalidParentCode, wantStatus: http.StatusBadRequest, wantError: "invalid parent code"},
			{name: "expired", err: service.ErrParentCodeExpired, wantStatus: http.StatusBadRequest, wantError: "parent code has expired"},
			{name: "max usage", err: service.ErrParentCodeMaxUsageReached, wantStatus: http.StatusBadRequest, wantError: "parent code has reached maximum usage"},
			{name: "email exists", err: service.ErrEmailAlreadyExists, wantStatus: http.StatusConflict, wantError: "email already exists"},
			{name: "empty password", err: service.ErrPasswordCannotBeEmpty, wantStatus: http.StatusBadRequest, wantError: "password cannot be empty"},
			{name: "failed hash", err: service.ErrFailedToHashPassword, wantStatus: http.StatusInternalServerError, wantError: "failed to hash password"},
			{name: "failed create user", err: service.ErrFailedToCreateUser, wantStatus: http.StatusInternalServerError, wantError: "failed to create user"},
			{name: "failed create parent", err: service.ErrFailedToCreateParent, wantStatus: http.StatusInternalServerError, wantError: "failed to create parent"},
			{name: "failed link", err: service.ErrFailedToLinkParentToStudent, wantStatus: http.StatusInternalServerError, wantError: "failed to link parent to student"},
			{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to register parent"},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				h := &ParentCodeHandler{parentCodeService: &fakeParentCodeService{registerParentFn: func(context.Context, string, string, string) (*service.LoginResponse, error) {
					return nil, tc.err
				}}}
				r := gin.New()
				r.POST("/register/parent", h.RegisterParent)

				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodPost, "/register/parent", strings.NewReader(`{"email":"parent@example.com","password":"secret123","parent_code":"PARENT001"}`))
				req.Header.Set("Content-Type", "application/json")
				r.ServeHTTP(rec, req)

				if rec.Code != tc.wantStatus {
					t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
				}
				if got := decodeParentCodeError(t, rec); got != tc.wantError {
					t.Fatalf("error = %q, want %q", got, tc.wantError)
				}
			})
		}
	})
}

func TestVerifyCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	studentID := uuid.New()
	schoolID := uuid.New()
	classID := uuid.New()

	t.Run("code required", func(t *testing.T) {
		h := &ParentCodeHandler{}
		r := gin.New()
		r.GET("/parent-codes/verify", h.VerifyCode)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/parent-codes/verify", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("verify code mappings", func(t *testing.T) {
		cases := []struct {
			name       string
			err        error
			wantStatus int
			wantError  string
		}{
			{name: "invalid", err: service.ErrInvalidParentCode, wantStatus: http.StatusBadRequest, wantError: "invalid parent code"},
			{name: "expired", err: service.ErrParentCodeExpired, wantStatus: http.StatusBadRequest, wantError: "parent code has expired"},
			{name: "max usage", err: service.ErrParentCodeMaxUsageReached, wantStatus: http.StatusBadRequest, wantError: "parent code has reached maximum usage"},
			{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to verify code"},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				h := &ParentCodeHandler{parentCodeService: &fakeParentCodeService{verifyCodeFn: func(context.Context, string) (*model.StudentParentCode, error) {
					return nil, tc.err
				}}}
				r := gin.New()
				r.GET("/parent-codes/verify", h.VerifyCode)

				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/parent-codes/verify?code=PARENT001", nil)
				r.ServeHTTP(rec, req)

				if rec.Code != tc.wantStatus {
					t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
				}
				if got := decodeParentCodeError(t, rec); got != tc.wantError {
					t.Fatalf("error = %q, want %q", got, tc.wantError)
				}
			})
		}
	})

	t.Run("student info error", func(t *testing.T) {
		h := &ParentCodeHandler{parentCodeService: &fakeParentCodeService{
			verifyCodeFn: func(context.Context, string) (*model.StudentParentCode, error) {
				return &model.StudentParentCode{StudentID: studentID, Code: "PARENT001", UsageCount: 1, MaxUsage: 4, ExpiresAt: time.Now().Add(24 * time.Hour)}, nil
			},
			getStudentInfoFn: func(context.Context, uuid.UUID) (*model.Student, error) {
				return nil, errors.New("boom")
			},
		}}
		r := gin.New()
		r.GET("/parent-codes/verify", h.VerifyCode)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/parent-codes/verify?code=PARENT001", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &ParentCodeHandler{parentCodeService: &fakeParentCodeService{
			verifyCodeFn: func(context.Context, string) (*model.StudentParentCode, error) {
				return &model.StudentParentCode{StudentID: studentID, Code: "PARENT001", UsageCount: 1, MaxUsage: 4, ExpiresAt: time.Now().Add(24 * time.Hour)}, nil
			},
			getStudentInfoFn: func(context.Context, uuid.UUID) (*model.Student, error) {
				return &model.Student{StudentID: studentID, FullName: "Student A", DOB: time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC), Gender: "female", SchoolID: schoolID, CurrentClassID: classID}, nil
			},
		}}
		r := gin.New()
		r.GET("/parent-codes/verify", h.VerifyCode)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/parent-codes/verify?code=PARENT001", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if !strings.Contains(rec.Body.String(), studentID.String()) {
			t.Fatalf("response does not contain student id: %s", rec.Body.String())
		}
	})
}
