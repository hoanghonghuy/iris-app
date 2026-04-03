package service

import "errors"

// User-related errors
var (
	ErrInvalidUserID                = errors.New("invalid user ID")
	ErrUserNotFound                 = errors.New("user not found")
	ErrEmailCannotBeEmpty           = errors.New("email cannot be empty")
	ErrPasswordCannotBeEmpty        = errors.New("password cannot be empty")
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

// Password reset errors
var (
	ErrResetTokenInvalid = errors.New("invalid or expired reset token")
	ErrResetTokenUsed    = errors.New("reset token has already been used")
	ErrFailedToSendEmail = errors.New("failed to send reset email")
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
	ErrSchoolAdminNotFound        = errors.New("school admin not found")
	ErrCannotAssignRole           = errors.New("insufficient permissions to assign this role")
	ErrCannotAssignRoleSuperAdmin = errors.New("SUPER_ADMIN role requires dedicated promote flow with approval")
	ErrSchoolAccessDenied         = errors.New("access denied: resource does not belong to your school")
)

// Business logic errors
var (
	ErrInvalidClassID     = errors.New("invalid class ID")
	ErrClassNotFound      = errors.New("class not found")
	ErrSchoolNotFound     = errors.New("school not found")
	ErrStudentNotFound    = errors.New("student not found")
	ErrInvalidDate        = errors.New("invalid date format, use YYYY-MM-DD")
	ErrInvalidStatus      = errors.New("invalid attendance status")
	ErrTeacherNotFound    = errors.New("teacher not found")
	ErrForbidden          = errors.New("forbidden action")
	ErrTeacherNotAssigned = errors.New("teacher is not assigned to this class")
	ErrInvalidValue       = errors.New("invalid value")
)

// Chat-related errors
var (
	ErrChatCannotMessageSelf = errors.New("cannot create conversation with yourself")
	ErrChatTargetNotAllowed  = errors.New("target user is not allowed for direct conversation")
	ErrChatGroupNeedMembers  = errors.New("group conversation needs at least 2 participants")
	ErrChatNotParticipant    = errors.New("you are not a participant of this conversation")
	ErrChatEmptyMessage      = errors.New("message content cannot be empty")
)
