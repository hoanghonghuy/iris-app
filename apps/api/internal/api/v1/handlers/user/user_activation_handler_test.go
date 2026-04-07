package userhandlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

func TestActivateUserWithToken_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("invalid body", func(t *testing.T) {
		h := &UserHandler{}
		r := gin.New()
		r.POST("/users/activate-token", h.ActivateUserWithToken)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/users/activate-token", strings.NewReader("{bad-json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &UserHandler{userService: &fakeUserService{activateUserWithTokenFn: func(_ context.Context, token, password string) error {
			if token != "tok-1" || password != "secret123" {
				t.Fatalf("unexpected activation payload")
			}
			return nil
		}}}
		r := gin.New()
		r.POST("/users/activate-token", h.ActivateUserWithToken)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/users/activate-token", strings.NewReader(`{"token":"tok-1","password":"secret123"}`))
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
		{name: "token required", err: service.ErrActivationTokenRequired, wantStatus: http.StatusBadRequest, wantError: "activation token is required"},
		{name: "invalid token", err: service.ErrInvalidActivationToken, wantStatus: http.StatusBadRequest, wantError: "invalid activation token"},
		{name: "token expired", err: service.ErrActivationTokenExpired, wantStatus: http.StatusBadRequest, wantError: "activation token has expired"},
		{name: "password empty", err: service.ErrPasswordCannotBeEmpty, wantStatus: http.StatusBadRequest, wantError: "password cannot be empty"},
		{name: "hash error", err: service.ErrFailedToHashPassword, wantStatus: http.StatusInternalServerError, wantError: "failed to hash password"},
		{name: "activate error", err: service.ErrFailedToActivateUser, wantStatus: http.StatusInternalServerError, wantError: "failed to activate user"},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to activate user"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserHandler{userService: &fakeUserService{activateUserWithTokenFn: func(context.Context, string, string) error {
				return tt.err
			}}}
			r := gin.New()
			r.POST("/users/activate-token", h.ActivateUserWithToken)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/users/activate-token", strings.NewReader(`{"token":"tok-1","password":"secret123"}`))
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
