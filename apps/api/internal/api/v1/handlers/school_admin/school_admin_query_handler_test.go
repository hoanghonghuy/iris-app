package schooladminhandlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

func TestList_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("invalid pagination", func(t *testing.T) {
		h := &SchoolAdminHandler{}
		r := gin.New()
		r.GET("/school-admins", h.List)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/school-admins?limit=0&offset=-1", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &SchoolAdminHandler{schoolAdminService: &fakeSchoolAdminService{listFn: func(_ context.Context, limit, offset int) ([]model.SchoolAdmin, int, error) {
			if limit != 20 || offset != 0 {
				t.Fatalf("unexpected pagination forwarded")
			}
			return []model.SchoolAdmin{{AdminID: uuid.New(), UserID: uuid.New(), SchoolID: uuid.New()}}, 1, nil
		}}}
		r := gin.New()
		r.GET("/school-admins", h.List)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/school-admins?limit=20&offset=0", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("internal", func(t *testing.T) {
		h := &SchoolAdminHandler{schoolAdminService: &fakeSchoolAdminService{listFn: func(context.Context, int, int) ([]model.SchoolAdmin, int, error) {
			return nil, 0, errors.New("boom")
		}}}
		r := gin.New()
		r.GET("/school-admins", h.List)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/school-admins?limit=20&offset=0", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
		if got := decodeSchoolAdminError(t, rec); got != "failed to fetch school admins" {
			t.Fatalf("error = %q, want %q", got, "failed to fetch school admins")
		}
	})
}
