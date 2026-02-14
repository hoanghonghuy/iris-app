package service

import "errors"

// User-related errors
var (
	ErrInvalidUserID                = errors.New("invalid user ID")
	ErrUserNotFound                 = errors.New("user not found")
	ErrEmailCannotBeEmpty           = errors.New("email cannot be empty")
	ErrPasswordCannotBeEmpty        = errors.New("password cannot be empty")
	ErrInvalidPassword              = errors.New("password cannot be empty")
	ErrFailedToHashPassword         = errors.New("failed to hash password")
	ErrFailedToUpdatePassword       = errors.New("failed to update password")
	ErrFailedToActivateUser         = errors.New("failed to activate user")
	ErrFailedToGenerateTempPassword = errors.New("failed to generate temporary password")
	ErrFailedToCreateUser           = errors.New("failed to create user")
	ErrEmailAlreadyExists           = errors.New("email already exists")
)

// Role-related errors
var (
	ErrRolesCannotBeEmpty = errors.New("roles cannot be empty")
	ErrFailedToAssignRole = errors.New("failed to assign role")
	ErrInvalidRoleName    = errors.New("invalid role name")
)

// Activation-related errors
var (
	ErrActivationTokenRequired = errors.New("activation token is required")
	ErrInvalidActivationToken  = errors.New("invalid activation token")
	ErrActivationTokenExpired  = errors.New("activation token has expired")
)

// Parent-related errors
var (
	ErrInvalidParentCode           = errors.New("invalid parent code")
	ErrParentCodeExpired           = errors.New("parent code has expired")
	ErrParentCodeMaxUsageReached   = errors.New("parent code has reached maximum usage")
	ErrFailedToCreateParent        = errors.New("failed to create parent")
	ErrFailedToLinkParentToStudent = errors.New("failed to link parent to student")
	ErrFailedToGetStudent          = errors.New("failed to get student")
)

// School admin-related errors
var (
	ErrSchoolAdminNotFound = errors.New("school admin not found")
	ErrCannotAssignRole    = errors.New("insufficient permissions to assign this role")
	ErrSchoolAccessDenied  = errors.New("access denied: resource does not belong to your school")
)

// Business logic errors
var (
	ErrInvalidClassID     = errors.New("invalid class ID")
	ErrInvalidDate        = errors.New("invalid date format, use YYYY-MM-DD")
	ErrInvalidStatus      = errors.New("invalid attendance status")
	ErrTeacherNotFound    = errors.New("teacher not found")
	ErrForbidden          = errors.New("forbidden action")
	ErrTeacherNotAssigned = errors.New("teacher is not assigned to this class")
	ErrInvalidValue       = errors.New("invalid value")
)
