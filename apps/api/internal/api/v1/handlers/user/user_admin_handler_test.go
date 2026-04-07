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
	"github.com/jackc/pgx/v5"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

func TestCreateUser_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()

	t.Run("invalid body", func(t *testing.T) {
		h := &UserHandler{}
		r := gin.New()
		r.POST("/users", h.CreateUser)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader("{bad-json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		h := &UserHandler{userService: &fakeUserService{createUserWithoutPasswordFn: func(_ context.Context, adminSchoolID *uuid.UUID, email string, roles []string) (*model.UserInfo, error) {
			if adminSchoolID == nil || *adminSchoolID != schoolID {
				t.Fatalf("expected admin school scope")
			}
			if email != "u@example.com" || len(roles) != 1 || roles[0] != "PARENT" {
				t.Fatalf("unexpected payload")
			}
			return &model.UserInfo{UserID: userID, Email: email}, nil
		}}}
		r := gin.New()
		r.Use(withUserAdminScope(schoolID))
		r.POST("/users", h.CreateUser)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"email":"u@example.com","roles":["PARENT"]}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
		}
	})

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantError  string
	}{
		{name: "cannot assign role", err: service.ErrCannotAssignRole, wantStatus: http.StatusForbidden, wantError: "insufficient permissions to assign this role"},
		{name: "super admin forbidden", err: service.ErrCannotAssignRoleSuperAdmin, wantStatus: http.StatusForbidden, wantError: "SUPER_ADMIN role requires dedicated promote flow"},
		{name: "assign role failed", err: service.ErrFailedToAssignRole, wantStatus: http.StatusBadRequest, wantError: "failed to assign role"},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to create user"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserHandler{userService: &fakeUserService{createUserWithoutPasswordFn: func(context.Context, *uuid.UUID, string, []string) (*model.UserInfo, error) {
				return nil, tt.err
			}}}
			r := gin.New()
			r.Use(withUserAdminScope(schoolID))
			r.POST("/users", h.CreateUser)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"email":"u@example.com","roles":["PARENT"]}`))
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

func TestGetByID_List_Lock_Unlock_AssignRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	schoolID := uuid.New()
	userID := uuid.New()

	t.Run("get by id invalid user id", func(t *testing.T) {
		h := &UserHandler{}
		r := gin.New()
		r.GET("/users/:user_id", h.GetByID)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/not-a-uuid", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("get by id success", func(t *testing.T) {
		h := &UserHandler{userService: &fakeUserService{findByIDFn: func(_ context.Context, adminSchoolID *uuid.UUID, gotUserID uuid.UUID) (*model.UserInfo, error) {
			if adminSchoolID == nil || *adminSchoolID != schoolID || gotUserID != userID {
				t.Fatalf("unexpected find params")
			}
			return &model.UserInfo{UserID: userID, Email: "u@example.com"}, nil
		}}}
		r := gin.New()
		r.Use(withUserAdminScope(schoolID))
		r.GET("/users/:user_id", h.GetByID)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/"+userID.String(), nil)
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("get by id mappings", func(t *testing.T) {
		cases := []struct {
			name       string
			err        error
			wantStatus int
			wantError  string
		}{
			{name: "access denied", err: service.ErrSchoolAccessDenied, wantStatus: http.StatusForbidden, wantError: "access denied"},
			{name: "not found", err: pgx.ErrNoRows, wantStatus: http.StatusNotFound, wantError: "user not found"},
			{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError, wantError: "failed to get user"},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				h := &UserHandler{userService: &fakeUserService{findByIDFn: func(context.Context, *uuid.UUID, uuid.UUID) (*model.UserInfo, error) {
					return nil, tc.err
				}}}
				r := gin.New()
				r.Use(withUserAdminScope(schoolID))
				r.GET("/users/:user_id", h.GetByID)
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/users/"+userID.String(), nil)
				r.ServeHTTP(rec, req)
				if rec.Code != tc.wantStatus {
					t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
				}
				if got := decodeUserError(t, rec); got != tc.wantError {
					t.Fatalf("error = %q, want %q", got, tc.wantError)
				}
			})
		}
	})

	t.Run("list success", func(t *testing.T) {
		h := &UserHandler{userService: &fakeUserService{listFn: func(_ context.Context, adminSchoolID *uuid.UUID, role string, limit, offset int) ([]model.UserInfo, int, error) {
			if adminSchoolID == nil || *adminSchoolID != schoolID || role != "PARENT" || limit != 20 || offset != 0 {
				t.Fatalf("unexpected list params")
			}
			return []model.UserInfo{{UserID: userID}}, 1, nil
		}}}
		r := gin.New()
		r.Use(withUserAdminScope(schoolID))
		r.GET("/users", h.List)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users?role=PARENT&limit=20&offset=0", nil)
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("list invalid pagination", func(t *testing.T) {
		h := &UserHandler{}
		r := gin.New()
		r.GET("/users", h.List)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users?limit=0&offset=-1", nil)
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("lock/unlock/assign role mappings", func(t *testing.T) {
		h := &UserHandler{userService: &fakeUserService{
			lockFn:       func(context.Context, *uuid.UUID, uuid.UUID) error { return service.ErrSchoolAccessDenied },
			unlockFn:     func(context.Context, *uuid.UUID, uuid.UUID) error { return errors.New("boom") },
			assignRoleFn: func(context.Context, uuid.UUID, string) error { return service.ErrInvalidRoleName },
		}}
		r := gin.New()
		r.Use(withUserAdminScope(schoolID))
		r.POST("/users/:user_id/lock", h.Lock)
		r.POST("/users/:user_id/unlock", h.Unlock)
		r.POST("/users/:user_id/roles", h.AssignRole)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/users/"+userID.String()+"/lock", nil)
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("lock status = %d, want %d", rec.Code, http.StatusForbidden)
		}

		rec = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodPost, "/users/"+userID.String()+"/unlock", nil)
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("unlock status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}

		rec = httptest.NewRecorder()
		req = httptest.NewRequest(http.MethodPost, "/users/"+userID.String()+"/roles", strings.NewReader(`{"role_name":"BAD_ROLE"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("assign role status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})
}
