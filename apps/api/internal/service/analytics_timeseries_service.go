package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
	"github.com/hoanghonghuy/iris-app/apps/api/internal/repo"
)

// Lỗi validate cho GET analytics/timeseries (HTTP 400).
var (
	ErrAnalyticsInvalidRange    = errors.New("invalid analytics range")
	ErrAnalyticsInvalidInterval = errors.New("invalid analytics interval (supported: day)")
)

func normalizeAnalyticsRange(rangeStr string) (days int, canonical string, err error) {
	switch rangeStr {
	case "", "14d":
		return 14, "14d", nil
	case "7d":
		return 7, "7d", nil
	case "30d":
		return 30, "30d", nil
	case "90d":
		return 90, "90d", nil
	default:
		return 0, "", ErrAnalyticsInvalidRange
	}
}

func analyticsInclusiveDayBoundsUTC(days int) (fromDate, toDate time.Time) {
	now := time.Now().UTC()
	toDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	fromDate = toDate.AddDate(0, 0, -(days - 1))
	return fromDate, toDate
}

func floatPtr(v float64) *float64 {
	return &v
}

func mergeStatusByDay(rows []repo.StatusDayCount) map[time.Time]map[string]int {
	out := make(map[time.Time]map[string]int)
	for _, r := range rows {
		d := r.Day.UTC()
		d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
		if out[d] == nil {
			out[d] = make(map[string]int)
		}
		out[d][r.Status] += r.Count
	}
	return out
}

func statusPointsForDayRange(fromDate, toDate time.Time, merged map[time.Time]map[string]int, allStatuses []string) []model.TimeseriesPoint {
	var pts []model.TimeseriesPoint
	for d := fromDate; !d.After(toDate); d = d.AddDate(0, 0, 1) {
		comp := make(map[string]int)
		for _, st := range allStatuses {
			comp[st] = 0
		}
		if m, ok := merged[d]; ok {
			for k, v := range m {
				comp[k] += v
			}
		}
		pts = append(pts, model.TimeseriesPoint{BucketStart: d, Components: comp})
	}
	return pts
}

// GetAdminAnalyticsTimeseries series cho dashboard admin (additive API).
func (s *AnalyticsService) GetAdminAnalyticsTimeseries(ctx context.Context, schoolID *uuid.UUID, rangeStr, intervalStr string) (*model.AnalyticsTimeseriesResponse, error) {
	days, canonRange, err := normalizeAnalyticsRange(rangeStr)
	if err != nil {
		return nil, err
	}
	if intervalStr != "" && intervalStr != "day" {
		return nil, ErrAnalyticsInvalidInterval
	}
	fromDate, toDate := analyticsInclusiveDayBoundsUTC(days)

	ts := s.repos.AnalyticsTimeseriesRepo
	if ts == nil {
		return nil, fmt.Errorf("analytics timeseries repo not configured")
	}

	rateRows, err := ts.AdminAttendancePresentRateDaily(ctx, schoolID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	var ratePoints []model.TimeseriesPoint
	for _, row := range rateRows {
		v := row.Value
		ratePoints = append(ratePoints, model.TimeseriesPoint{BucketStart: row.Day, Value: floatPtr(v)})
	}

	healthRows, err := ts.AdminHealthAlertsDaily(ctx, schoolID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	var healthPoints []model.TimeseriesPoint
	for _, row := range healthRows {
		v := float64(row.Value)
		healthPoints = append(healthPoints, model.TimeseriesPoint{BucketStart: row.Day, Value: floatPtr(v)})
	}

	teachers, err := s.repos.UserRepo.CountUsersByRoleAndSchool(ctx, "TEACHER", schoolID)
	if err != nil {
		return nil, err
	}
	parents, err := s.repos.UserRepo.CountUsersByRoleAndSchool(ctx, "PARENT", schoolID)
	if err != nil {
		return nil, err
	}
	students, err := s.repos.StudentRepo.CountStudentsBySchool(ctx, schoolID)
	if err != nil {
		return nil, err
	}
	popPoint := model.TimeseriesPoint{
		BucketStart: toDate,
		Components: map[string]int{
			"teacher": teachers,
			"parent":  parents,
			"student": students,
		},
	}

	meta := model.AnalyticsTimeseriesMeta{
		Range:          canonRange,
		Interval:       "day",
		BucketTimezone: "UTC",
		GeneratedAt:    time.Now().UTC(),
	}
	if schoolID != nil {
		s := schoolID.String()
		meta.SchoolID = &s
	}

	return &model.AnalyticsTimeseriesResponse{
		Meta: meta,
		Series: []model.AnalyticsTimeseriesSeries{
			{ID: "attendance_rate", Label: "Tỷ lệ có mặt", Unit: "percent", Points: ratePoints},
			{ID: "health_alerts", Label: "Cảnh báo sức khỏe", Unit: "count", Points: healthPoints},
			{ID: "population_by_role", Label: "Quy mô (GV / PH / HS)", Unit: "count", Points: []model.TimeseriesPoint{popPoint}},
		},
	}, nil
}

// GetTeacherAnalyticsTimeseries series cho dashboard giáo viên.
func (s *AnalyticsService) GetTeacherAnalyticsTimeseries(ctx context.Context, teacherUserID uuid.UUID, rangeStr, intervalStr string) (*model.AnalyticsTimeseriesResponse, error) {
	days, canonRange, err := normalizeAnalyticsRange(rangeStr)
	if err != nil {
		return nil, err
	}
	if intervalStr != "" && intervalStr != "day" {
		return nil, ErrAnalyticsInvalidInterval
	}
	fromDate, toDate := analyticsInclusiveDayBoundsUTC(days)

	ts := s.repos.AnalyticsTimeseriesRepo
	if ts == nil {
		return nil, fmt.Errorf("analytics timeseries repo not configured")
	}

	totalStudents, err := s.repos.TeacherScopeRepo.CountMyStudents(ctx, teacherUserID)
	if err != nil {
		return nil, err
	}

	markedRows, err := ts.TeacherAttendanceMarkedDaily(ctx, teacherUserID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	var mpPoints []model.TimeseriesPoint
	for _, row := range markedRows {
		marked := row.Value
		pending := totalStudents - marked
		if pending < 0 {
			pending = 0
		}
		mpPoints = append(mpPoints, model.TimeseriesPoint{
			BucketStart: row.Day,
			Components: map[string]int{
				"marked":  marked,
				"pending": pending,
			},
		})
	}

	healthRows, err := ts.TeacherHealthAlertsDaily(ctx, teacherUserID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	var healthPoints []model.TimeseriesPoint
	for _, row := range healthRows {
		v := float64(row.Value)
		healthPoints = append(healthPoints, model.TimeseriesPoint{BucketStart: row.Day, Value: floatPtr(v)})
	}

	apptRows, err := ts.TeacherAppointmentsByStatusDay(ctx, teacherUserID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	merged := mergeStatusByDay(apptRows)
	apptStatuses := []string{"pending", "confirmed", "cancelled", "completed", "no_show"}
	apptPoints := statusPointsForDayRange(fromDate, toDate, merged, apptStatuses)

	return &model.AnalyticsTimeseriesResponse{
		Meta: model.AnalyticsTimeseriesMeta{
			Range:          canonRange,
			Interval:       "day",
			BucketTimezone: "UTC",
			GeneratedAt:    time.Now().UTC(),
		},
		Series: []model.AnalyticsTimeseriesSeries{
			{ID: "attendance_marked_vs_pending", Label: "Điểm danh: đã ghi / chưa ghi", Unit: "count", Points: mpPoints},
			{ID: "health_alerts", Label: "Cảnh báo sức khỏe", Unit: "count", Points: healthPoints},
			{ID: "appointments_by_status", Label: "Lịch hẹn theo trạng thái", Unit: "count", Points: apptPoints},
		},
	}, nil
}

// GetParentAnalyticsTimeseries series cho dashboard phụ huynh (theo một học sinh).
func (s *AnalyticsService) GetParentAnalyticsTimeseries(ctx context.Context, parentUserID, studentID uuid.UUID, rangeStr, intervalStr string) (*model.AnalyticsTimeseriesResponse, error) {
	if studentID == uuid.Nil {
		return nil, ErrInvalidUserID
	}
	ok, err := s.repos.ParentScopeRepo.IsParentOfStudent(ctx, parentUserID, studentID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrForbidden
	}

	days, canonRange, err := normalizeAnalyticsRange(rangeStr)
	if err != nil {
		return nil, err
	}
	if intervalStr != "" && intervalStr != "day" {
		return nil, ErrAnalyticsInvalidInterval
	}
	fromDate, toDate := analyticsInclusiveDayBoundsUTC(days)

	ts := s.repos.AnalyticsTimeseriesRepo
	if ts == nil {
		return nil, fmt.Errorf("analytics timeseries repo not configured")
	}

	attRows, err := ts.ParentAttendanceByStatusDaily(ctx, studentID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	attMerged := mergeStatusByDay(attRows)
	attStatuses := []string{"present", "absent", "late", "excused"}
	attPoints := statusPointsForDayRange(fromDate, toDate, attMerged, attStatuses)

	healthRows, err := ts.ParentHealthAlertsDaily(ctx, studentID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	var healthPoints []model.TimeseriesPoint
	for _, row := range healthRows {
		v := float64(row.Value)
		healthPoints = append(healthPoints, model.TimeseriesPoint{BucketStart: row.Day, Value: floatPtr(v)})
	}

	apptRows, err := ts.ParentAppointmentsByStatusDay(ctx, parentUserID, studentID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	apptMerged := mergeStatusByDay(apptRows)
	apptStatuses := []string{"pending", "confirmed", "cancelled", "completed", "no_show"}
	apptPoints := statusPointsForDayRange(fromDate, toDate, apptMerged, apptStatuses)

	sid := studentID.String()
	return &model.AnalyticsTimeseriesResponse{
		Meta: model.AnalyticsTimeseriesMeta{
			Range:          canonRange,
			Interval:       "day",
			BucketTimezone: "UTC",
			GeneratedAt:    time.Now().UTC(),
			StudentID:      &sid,
		},
		Series: []model.AnalyticsTimeseriesSeries{
			{ID: "child_attendance", Label: "Điểm danh theo trạng thái", Unit: "count", Points: attPoints},
			{ID: "health_alerts", Label: "Cảnh báo sức khỏe", Unit: "count", Points: healthPoints},
			{ID: "appointments_by_status", Label: "Lịch hẹn theo trạng thái", Unit: "count", Points: apptPoints},
		},
	}, nil
}
