package analyticshandlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/auth"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/middleware"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

type fakeAnalyticsService struct {
	getAdminAnalyticsFn   func(context.Context, *uuid.UUID) (*model.AdminAnalytics, error)
	getTeacherAnalyticsFn func(context.Context, uuid.UUID) (*model.TeacherAnalytics, error)
}

func (f *fakeAnalyticsService) GetAdminAnalytics(ctx context.Context, schoolID *uuid.UUID) (*model.AdminAnalytics, error) {
	if f.getAdminAnalyticsFn == nil {
		return nil, errors.New("unexpected GetAdminAnalytics call")
	}
	return f.getAdminAnalyticsFn(ctx, schoolID)
}

func (f *fakeAnalyticsService) GetTeacherAnalytics(ctx context.Context, teacherUserID uuid.UUID) (*model.TeacherAnalytics, error) {
	if f.getTeacherAnalyticsFn == nil {
		return nil, errors.New("unexpected GetTeacherAnalytics call")
	}
	return f.getTeacherAnalyticsFn(ctx, teacherUserID)
}

func withAnalyticsAdminScope(schoolID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.CtxAdminSchoolID, schoolID.String())
		c.Next()
	}
}

func withAnalyticsTeacherClaims(userID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middleware.CtxClaims, &auth.Claims{UserID: userID, Roles: []string{"TEACHER"}})
		c.Next()
	}
}

func decodeAnalyticsError(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	return body["error"]
}
