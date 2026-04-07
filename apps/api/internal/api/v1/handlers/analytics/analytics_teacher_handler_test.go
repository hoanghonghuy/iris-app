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

func TestTeacherDashboardStats_SuccessAndErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)
	teacherID := uuid.New()

	t.Run("missing claims", func(t *testing.T) {
		h := &AnalyticsHandler{}
		r := gin.New()
		r.GET("/teacher/analytics", h.TeacherDashboardStats)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/teacher/analytics", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
		}
		if got := decodeAnalyticsError(t, rec); got != "unauthorized" {
			t.Fatalf("error = %q, want %q", got, "unauthorized")
		}
	})

	t.Run("invalid claims user id", func(t *testing.T) {
		h := &AnalyticsHandler{}
		r := gin.New()
		r.Use(withAnalyticsTeacherClaims("not-a-uuid"))
		r.GET("/teacher/analytics", h.TeacherDashboardStats)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/teacher/analytics", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		if got := decodeAnalyticsError(t, rec); got != "invalid user ID" {
			t.Fatalf("error = %q, want %q", got, "invalid user ID")
		}
	})

	t.Run("service error", func(t *testing.T) {
		h := &AnalyticsHandler{analyticsService: &fakeAnalyticsService{getTeacherAnalyticsFn: func(context.Context, uuid.UUID) (*model.TeacherAnalytics, error) {
			return nil, errors.New("boom")
		}}}
		r := gin.New()
		r.Use(withAnalyticsTeacherClaims(teacherID.String()))
		r.GET("/teacher/analytics", h.TeacherDashboardStats)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/teacher/analytics", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
		if got := decodeAnalyticsError(t, rec); got != "Lỗi khi thống kê dữ liệu giáo viên: boom" {
			t.Fatalf("error = %q, want %q", got, "Lỗi khi thống kê dữ liệu giáo viên: boom")
		}
	})

	t.Run("success", func(t *testing.T) {
		h := &AnalyticsHandler{analyticsService: &fakeAnalyticsService{getTeacherAnalyticsFn: func(_ context.Context, gotTeacherID uuid.UUID) (*model.TeacherAnalytics, error) {
			if gotTeacherID != teacherID {
				t.Fatalf("teacher_id = %s, want %s", gotTeacherID, teacherID)
			}
			return &model.TeacherAnalytics{TotalClasses: 2, TotalStudents: 30, TotalPosts: 10}, nil
		}}}
		r := gin.New()
		r.Use(withAnalyticsTeacherClaims(teacherID.String()))
		r.GET("/teacher/analytics", h.TeacherDashboardStats)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/teacher/analytics", nil)
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})
}
