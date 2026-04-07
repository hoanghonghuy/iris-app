package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/hoanghonghuy/iris-app/apps/api/internal/model"
)

func TestAuditLogServiceCreateValidation(t *testing.T) {
	svc := &AuditLogService{}

	err := svc.Create(context.Background(), model.AuditLogCreate{})
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidUserID)
	}

	err = svc.Create(context.Background(), model.AuditLogCreate{ActorUserID: uuid.New()})
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidValue)
	}

	err = svc.Create(context.Background(), model.AuditLogCreate{ActorUserID: uuid.New(), Action: "user.created"})
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidValue)
	}
}

func TestAuditLogServiceParseTimeRange(t *testing.T) {
	svc := &AuditLogService{}

	from, to, err := svc.ParseTimeRange("2026-04-07T10:00:00Z", "2026-04-07T12:00:00Z")
	if err != nil {
		t.Fatalf("ParseTimeRange() error = %v", err)
	}
	if from == nil || to == nil {
		t.Fatalf("expected both from and to to be parsed")
	}

	from, to, err = svc.ParseTimeRange("", "")
	if err != nil {
		t.Fatalf("ParseTimeRange() error = %v", err)
	}
	if from != nil || to != nil {
		t.Fatalf("expected nil pointers for empty inputs")
	}

	_, _, err = svc.ParseTimeRange("bad-from", "")
	if !errors.Is(err, ErrInvalidDate) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidDate)
	}

	_, _, err = svc.ParseTimeRange("", "bad-to")
	if !errors.Is(err, ErrInvalidDate) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidDate)
	}

	_, _, err = svc.ParseTimeRange("2026-04-07T12:00:00Z", "2026-04-07T11:00:00Z")
	if !errors.Is(err, ErrInvalidDate) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidDate)
	}

	fromRaw := "2026-04-07T12:00:00+07:00"
	toRaw := "2026-04-07T06:00:00+01:00"
	expectedFrom, _ := time.Parse(time.RFC3339, fromRaw)
	expectedTo, _ := time.Parse(time.RFC3339, toRaw)

	from, to, err = svc.ParseTimeRange(fromRaw, toRaw)
	if err != nil {
		t.Fatalf("ParseTimeRange() timezone offset error = %v", err)
	}
	if from == nil || to == nil {
		t.Fatalf("expected non-nil from/to for timezone inputs")
	}
	if !from.Equal(expectedFrom) || !to.Equal(expectedTo) {
		t.Fatalf("timezone parse mismatch: got from=%v to=%v", *from, *to)
	}
}

func TestAuditLogServiceListValidation(t *testing.T) {
	svc := &AuditLogService{}

	from := time.Date(2026, time.April, 7, 12, 0, 0, 0, time.UTC)
	to := from.Add(-time.Hour)

	_, _, err := svc.List(context.Background(), model.AuditLogFilter{From: &from, To: &to, Limit: 20, Offset: 0})
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidValue)
	}
}
