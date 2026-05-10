package model

import "time"

// AnalyticsTimeseriesResponse payload cho GET .../analytics/timeseries (additive, không đổi /analytics).
type AnalyticsTimeseriesResponse struct {
	Meta   AnalyticsTimeseriesMeta     `json:"meta"`
	Series []AnalyticsTimeseriesSeries `json:"series"`
}

// AnalyticsTimeseriesMeta mô tả bucket và thời điểm sinh response.
type AnalyticsTimeseriesMeta struct {
	Range          string    `json:"range"`
	Interval       string    `json:"interval"`
	BucketTimezone string    `json:"bucket_timezone"`
	GeneratedAt    time.Time `json:"generated_at"`
	SchoolID       *string   `json:"school_id,omitempty"`
	StudentID      *string   `json:"student_id,omitempty"`
}

// AnalyticsTimeseriesSeries một chuỗi dữ liệu (line/bar/stacked).
type AnalyticsTimeseriesSeries struct {
	ID     string            `json:"id"`
	Label  string            `json:"label"`
	Unit   string            `json:"unit"`
	Points []TimeseriesPoint `json:"points"`
}

// TimeseriesPoint một điểm theo bucket thời gian (UTC).
type TimeseriesPoint struct {
	BucketStart time.Time      `json:"bucket_start"`
	Value       *float64       `json:"value,omitempty"`
	Components  map[string]int `json:"components,omitempty"`
}
