package analyticshandlers

import (
	"context"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/service"
)

type analyticsService interface {
	GetAdminAnalytics(ctx context.Context, schoolID *uuid.UUID) (*model.AdminAnalytics, error)
	GetTeacherAnalytics(ctx context.Context, teacherUserID uuid.UUID) (*model.TeacherAnalytics, error)
}

type AnalyticsHandler struct {
	analyticsService analyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}
