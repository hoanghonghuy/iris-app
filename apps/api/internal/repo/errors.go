package repo

import "errors"

// Các lỗi kỹ thuật dùng trong repository layer.
// service layer sẽ chuyển đổi các lỗi này thành lỗi nghiệp vụ phù hợp.

// ErrNoRowsUpdated được trả về khi câu UPDATE không ảnh hưởng đến hàng nào.
// => điều kiện WHERE không thỏa mãn (không phải lỗi DB),
var ErrNoRowsUpdated = errors.New("no rows updated")

// ErrGoogleAlreadyLinkedDifferent trả về khi user đã liên kết với một Google account khác.
var ErrGoogleAlreadyLinkedDifferent = errors.New("user already linked with a different google account")

// ErrRoleAssignmentFailed trả về khi gán role thất bại trong flow tạo user.
var ErrRoleAssignmentFailed = errors.New("failed to assign role")

// ErrAppointmentSlotUnavailable trả về khi slot đã có lịch hẹn active.
var ErrAppointmentSlotUnavailable = errors.New("appointment slot is unavailable")

// ErrAppointmentCancellationWindowPassed trả về khi đã quá thời gian cho phép hủy lịch.
var ErrAppointmentCancellationWindowPassed = errors.New("appointment cancellation window has passed")
