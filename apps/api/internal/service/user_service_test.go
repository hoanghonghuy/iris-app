package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestCreateUserWithoutPasswordRejectsSuperAdminRole(t *testing.T) {
	svc := &UserService{}

	_, err := svc.CreateUserWithoutPassword(context.Background(), nil, "admin@example.com", []string{"SUPER_ADMIN"})
	if !errors.Is(err, ErrCannotAssignRoleSuperAdmin) {
		t.Fatalf("err = %v, want %v", err, ErrCannotAssignRoleSuperAdmin)
	}
}

func TestAssignRoleRejectsSuperAdminRole(t *testing.T) {
	svc := &UserService{}

	err := svc.AssignRole(context.Background(), mustUUID(t, "f7f1d4cb-9708-4fa2-b1ab-e7f58d2bb1ee"), "SUPER_ADMIN")
	if !errors.Is(err, ErrCannotAssignRoleSuperAdmin) {
		t.Fatalf("err = %v, want %v", err, ErrCannotAssignRoleSuperAdmin)
	}
}

func TestAssignRoleRejectsInvalidRoleName(t *testing.T) {
	svc := &UserService{}

	err := svc.AssignRole(context.Background(), mustUUID(t, "f7f1d4cb-9708-4fa2-b1ab-e7f58d2bb1ee"), "INVALID_ROLE")
	if !errors.Is(err, ErrInvalidRoleName) {
		t.Fatalf("err = %v, want %v", err, ErrInvalidRoleName)
	}
}

func mustUUID(t *testing.T, value string) uuid.UUID {
	t.Helper()
	id, err := uuid.Parse(value)
	if err != nil {
		t.Fatalf("invalid uuid in test: %v", err)
	}
	return id
}
