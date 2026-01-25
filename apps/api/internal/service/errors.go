package service

import "errors"

// Sentinel errors for the service layer
var (
	ErrInvalidUserID                = errors.New("invalid user ID")
	ErrInvalidClassID               = errors.New("invalid class ID")
	ErrInvalidDate                  = errors.New("invalid date format, use YYYY-MM-DD")
	ErrInvalidStatus                = errors.New("invalid attendance status")
	ErrUserNotFound                 = errors.New("user not found")
	ErrTeacherNotFound              = errors.New("teacher not found")
	ErrInvalidPassword              = errors.New("password cannot be empty")
	ErrEmailCannotBeEmpty           = errors.New("email cannot be empty")
	ErrPasswordCannotBeEmpty        = errors.New("password cannot be empty")
	ErrFailedToHashPassword         = errors.New("failed to hash password")
	ErrFailedToUpdatePassword       = errors.New("failed to update password")
	ErrFailedToActivateUser         = errors.New("failed to activate user")
	ErrRolesCannotBeEmpty           = errors.New("roles cannot be empty")
	ErrFailedToGenerateTempPassword = errors.New("failed to generate temporary password")
	ErrFailedToCreateUser           = errors.New("failed to create user")
	ErrFailedToAssignRole           = errors.New("failed to assign role")
	ErrInvalidRoleName              = errors.New("invalid role name")
	ErrForbidden                    = errors.New("forbidden action")
)
