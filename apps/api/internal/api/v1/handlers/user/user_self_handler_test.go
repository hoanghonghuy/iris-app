package userhandlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

func TestUpdateMyPassword_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()

	t.Run("invalid body", func(t *testing.T) {
		h := &UserHandler{}
		r := gin.New()
		r.PUT("/me/password", h.UpdateMyPassword)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/me/password", strings.NewReader("{bad-json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing claims", func(t *testing.T) {
		h := &UserHandler{}
		r := gin.New()
		r.PUT("/me/password", h.UpdateMyPassword)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/me/password", strings.NewReader(`{"password":"newpass123"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &UserHandler{userService: &fakeUserService{updateMyPasswordFn: func(_ context.Context, gotUserID uuid.UUID, password string) error {
			if gotUserID != userID || password != "newpass123" {
				t.Fatalf("unexpected payload")
			}
			return nil
		}}}
		r := gin.New()
		r.Use(withUserClaims(userID.String()))
		r.PUT("/me/password", h.UpdateMyPassword)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/me/password", strings.NewReader(`{"password":"newpass123"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantError  string
	}{
		{name: "invalid user id", err: service.ErrInvalidUserID, wantStatus: http.StatusBadRequest, wantError: "invalid user ID"},
		{name: "not found", err: service.ErrUserNotFound, wantStatus: http.StatusNotFound, wantError: "user not found"},
		{name: "empty", err: service.ErrPasswordCannotBeEmpty, wantStatus: http.StatusBadRequest, wantError: "password cannot be empty"},
		{name: "hash", err: service.ErrFailedToHashPassword, wantStatus: http.StatusInternalServerError, wantError: "failed to hash password"},
		{name: "update", err: service.ErrFailedToUpdatePassword, wantStatus: http.StatusInternalServerError, wantError: "failed to update password"},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to update password"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserHandler{userService: &fakeUserService{updateMyPasswordFn: func(context.Context, uuid.UUID, string) error {
				return tt.err
			}}}
			r := gin.New()
			r.Use(withUserClaims(userID.String()))
			r.PUT("/me/password", h.UpdateMyPassword)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/me/password", strings.NewReader(`{"password":"newpass123"}`))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if got := decodeUserError(t, rec); got != tt.wantError {
				t.Fatalf("error = %q, want %q", got, tt.wantError)
			}
		})
	}
}

func TestDelete_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()

	t.Run("missing claims", func(t *testing.T) {
		h := &UserHandler{}
		r := gin.New()
		r.DELETE("/me", h.Delete)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/me", nil)
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &UserHandler{userService: &fakeUserService{deleteFn: func(_ context.Context, gotUserID uuid.UUID) error {
			if gotUserID != userID {
				t.Fatalf("unexpected user id")
			}
			return nil
		}}}
		r := gin.New()
		r.Use(withUserClaims(userID.String()))
		r.DELETE("/me", h.Delete)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/me", nil)
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("internal", func(t *testing.T) {
		h := &UserHandler{userService: &fakeUserService{deleteFn: func(context.Context, uuid.UUID) error {
			return errors.New("boom")
		}}}
		r := gin.New()
		r.Use(withUserClaims(userID.String()))
		r.DELETE("/me", h.Delete)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/me", nil)
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
		if got := decodeUserError(t, rec); got != "failed to delete user" {
			t.Fatalf("error = %q, want %q", got, "failed to delete user")
		}
	})
}
