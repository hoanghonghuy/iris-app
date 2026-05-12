import type { AttendanceStatus } from '@/types'

interface StatusOption {
  value: AttendanceStatus
  label: string
  badge: string
}

export const ATTENDANCE_STATUS_OPTIONS: StatusOption[] = [
  { value: 'present', label: 'Có mặt', badge: 'badge badge--primary' },
  { value: 'absent', label: 'Vắng', badge: 'badge badge--danger' },
  { value: 'late', label: 'Muộn', badge: 'badge badge--warning' },
  { value: 'excused', label: 'Có phép', badge: 'badge badge--outline' },
]

function getStatusOption(status: AttendanceStatus): StatusOption | undefined {
  return ATTENDANCE_STATUS_OPTIONS.find((option) => option.value === status)
}

export function getStatusLabel(status: AttendanceStatus | null | undefined): string {
  return getStatusOption(status as AttendanceStatus)?.label || status || '-'
}

export function getDefaultAttendanceValue(): { status: AttendanceStatus; note: string } {
  return { status: 'present', note: '' }
}

export function isSameAttendance(
  left: { status?: AttendanceStatus; note?: string } | null | undefined,
  right: { status?: AttendanceStatus; note?: string } | null | undefined,
): boolean {
  return left?.status === right?.status && (left?.note || '') === (right?.note || '')
}
