package analyticshandlers

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

func TestAdminDashboardStats_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success with school scope", func(t *testing.T) {
		schoolID := uuid.New()
		h := &AnalyticsHandler{analyticsService: &fakeAnalyticsService{getAdminAnalyticsFn: func(_ context.Context, gotSchoolID *uuid.UUID) (*model.AdminAnalytics, error) {
			if gotSchoolID == nil || *gotSchoolID != schoolID {
				t.Fatalf("expected school scope to be forwarded")
			}
			return &model.AdminAnalytics{TotalSchools: 1, TotalClasses: 2, TotalTeachers: 3, TotalStudents: 4, TotalParents: 5}, nil
		}}}

		r := gin.New()
		r.Use(withAnalyticsAdminScope(schoolID))
		r.GET("/admin/analytics", h.AdminDashboardStats)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/analytics", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("success with super admin scope", func(t *testing.T) {
		h := &AnalyticsHandler{analyticsService: &fakeAnalyticsService{getAdminAnalyticsFn: func(_ context.Context, gotSchoolID *uuid.UUID) (*model.AdminAnalytics, error) {
			if gotSchoolID != nil {
				t.Fatalf("expected nil school scope for super admin")
			}
			return &model.AdminAnalytics{TotalSchools: 10}, nil
		}}}

		r := gin.New()
		r.GET("/admin/analytics", h.AdminDashboardStats)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/analytics", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("internal error", func(t *testing.T) {
		h := &AnalyticsHandler{analyticsService: &fakeAnalyticsService{getAdminAnalyticsFn: func(context.Context, *uuid.UUID) (*model.AdminAnalytics, error) {
			return nil, errors.New("boom")
		}}}

		r := gin.New()
		r.GET("/admin/analytics", h.AdminDashboardStats)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/admin/analytics", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
		if got := decodeAnalyticsError(t, rec); got != "Lỗi khi thống kê dữ liệu: boom" {
			t.Fatalf("error = %q, want %q", got, "Lỗi khi thống kê dữ liệu: boom")
		}
	})
}
