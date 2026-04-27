export const ATTENDANCE_STATUS_OPTIONS = [
  { value: 'present', label: 'Có mặt', badge: 'badge badge--primary' },
  { value: 'absent', label: 'Vắng', badge: 'badge badge--danger' },
  { value: 'late', label: 'Muộn', badge: 'badge badge--warning' },
  { value: 'excused', label: 'Có phép', badge: 'badge badge--outline' },
]

function getStatusOption(status) {
  return ATTENDANCE_STATUS_OPTIONS.find((option) => option.value === status)
}

export function getStatusLabel(status) {
  return getStatusOption(status)?.label || status || '-'
}

export function getDefaultAttendanceValue() {
  return { status: 'present', note: '' }
}

export function isSameAttendance(left, right) {
  return left?.status === right?.status && (left?.note || '') === (right?.note || '')
}
