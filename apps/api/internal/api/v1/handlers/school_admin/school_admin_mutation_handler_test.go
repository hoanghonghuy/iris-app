package schooladminhandlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

func TestCreate_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	schoolID := uuid.New()

	t.Run("invalid body", func(t *testing.T) {
		h := &SchoolAdminHandler{}
		r := gin.New()
		r.POST("/school-admins", h.Create)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/school-admins", strings.NewReader("{bad-json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		adminID := uuid.New()
		h := &SchoolAdminHandler{schoolAdminService: &fakeSchoolAdminService{createFn: func(_ context.Context, gotUserID, gotSchoolID uuid.UUID, fullName, phone string) (*model.SchoolAdmin, error) {
			if gotUserID != userID || gotSchoolID != schoolID || fullName != "Admin A" || phone != "0909" {
				t.Fatalf("unexpected create payload")
			}
			return &model.SchoolAdmin{AdminID: adminID}, nil
		}}}
		r := gin.New()
		r.POST("/school-admins", h.Create)

		body := `{"user_id":"` + userID.String() + `","school_id":"` + schoolID.String() + `","full_name":"Admin A","phone":"0909"}`
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/school-admins", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
		}
		if rec.Header().Get("Location") != "/school-admins/"+adminID.String() {
			t.Fatalf("location = %q, want %q", rec.Header().Get("Location"), "/school-admins/"+adminID.String())
		}
	})

	t.Run("internal", func(t *testing.T) {
		h := &SchoolAdminHandler{schoolAdminService: &fakeSchoolAdminService{createFn: func(context.Context, uuid.UUID, uuid.UUID, string, string) (*model.SchoolAdmin, error) {
			return nil, errors.New("boom")
		}}}
		r := gin.New()
		r.POST("/school-admins", h.Create)

		body := `{"user_id":"` + userID.String() + `","school_id":"` + schoolID.String() + `"}`
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/school-admins", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
		if got := decodeSchoolAdminError(t, rec); got != "failed to create school admin" {
			t.Fatalf("error = %q, want %q", got, "failed to create school admin")
		}
	})
}

func TestDelete_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminID := uuid.New()

	t.Run("invalid admin id", func(t *testing.T) {
		h := &SchoolAdminHandler{}
		r := gin.New()
		r.DELETE("/school-admins/:admin_id", h.Delete)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/school-admins/not-a-uuid", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &SchoolAdminHandler{schoolAdminService: &fakeSchoolAdminService{deleteFn: func(_ context.Context, gotAdminID uuid.UUID) error {
			if gotAdminID != adminID {
				t.Fatalf("unexpected admin id forwarded")
			}
			return nil
		}}}
		r := gin.New()
		r.DELETE("/school-admins/:admin_id", h.Delete)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/school-admins/"+adminID.String(), nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("internal", func(t *testing.T) {
		h := &SchoolAdminHandler{schoolAdminService: &fakeSchoolAdminService{deleteFn: func(context.Context, uuid.UUID) error {
			return errors.New("boom")
		}}}
		r := gin.New()
		r.DELETE("/school-admins/:admin_id", h.Delete)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/school-admins/"+adminID.String(), nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
		if got := decodeSchoolAdminError(t, rec); got != "failed to delete school admin" {
			t.Fatalf("error = %q, want %q", got, "failed to delete school admin")
		}
	})
}
