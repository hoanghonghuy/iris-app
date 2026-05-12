import type { AppointmentStatus } from '@/types'

interface StatusConfig {
  label: string
  badge: string
  text: string
  badgeClass: string
}

export const APPOINTMENT_STATUS_CONFIG: Record<AppointmentStatus, StatusConfig> = {
  pending: {
    label: 'Chờ xác nhận',
    badge: 'badge--outline',
    text: 'Chờ xác nhận',
    badgeClass: 'badge badge--outline',
  },
  confirmed: {
    label: 'Đã xác nhận',
    badge: 'badge--primary',
    text: 'Đã xác nhận',
    badgeClass: 'badge badge--success',
  },
  cancelled: {
    label: 'Đã hủy',
    badge: 'badge--danger',
    text: 'Đã hủy',
    badgeClass: 'badge badge--danger',
  },
  completed: {
    label: 'Hoàn tất',
    badge: 'badge--outline',
    text: 'Đã hoàn thành',
    badgeClass: 'badge badge--primary',
  },
  no_show: {
    label: 'Vắng mặt',
    badge: 'badge--outline',
    text: 'Vắng mặt',
    badgeClass: 'badge badge--outline',
  },
}

export const APPOINTMENT_STATUS_OPTIONS = [
  { value: 'pending' as AppointmentStatus, label: 'Chờ xác nhận' },
  { value: 'confirmed' as AppointmentStatus, label: 'Đã xác nhận' },
  { value: 'cancelled' as AppointmentStatus, label: 'Đã hủy' },
  { value: 'completed' as AppointmentStatus, label: 'Hoàn tất' },
  { value: 'no_show' as AppointmentStatus, label: 'Vắng mặt' },
]

const CANCEL_REASON_MAP: Record<string, string> = {
  parent_cancelled: 'Phụ huynh đã hủy lịch',
  teacher_cancelled: 'Giáo viên đã hủy lịch',
  system_cancelled: 'Hệ thống đã hủy lịch',
}

export function getStatusBadge(status: AppointmentStatus | null | undefined): string {
  return APPOINTMENT_STATUS_CONFIG[status as AppointmentStatus]?.badgeClass || 'badge badge--outline'
}

export function getStatusText(status: AppointmentStatus | null | undefined): string {
  return APPOINTMENT_STATUS_CONFIG[status as AppointmentStatus]?.text || status || 'Không xác định'
}

export function getCancelReasonText(reason: string | null | undefined): string {
  return CANCEL_REASON_MAP[reason || ''] || reason || ''
}

export function getUtcOffsetLabel(): string {
  const totalMinutes = -new Date().getTimezoneOffset()
  const sign = totalMinutes >= 0 ? '+' : '-'
  const abs = Math.abs(totalMinutes)
  const hours = String(Math.floor(abs / 60)).padStart(2, '0')
  const mins = String(abs % 60).padStart(2, '0')
  return `UTC${sign}${hours}:${mins}`
}

export function getTimezoneDisplay(): string {
  const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone || 'Local'
  const totalMinutes = -new Date().getTimezoneOffset()
  const sign = totalMinutes >= 0 ? '+' : '-'
  const abs = Math.abs(totalMinutes)
  const hours = String(Math.floor(abs / 60)).padStart(2, '0')
  const minutes = String(abs % 60).padStart(2, '0')
  return `${timezone} (UTC${sign}${hours}:${minutes})`
}

function formatShortTime(value: string | null | undefined): string {
  if (!value) return ''
  return new Date(value).toLocaleTimeString('vi-VN', {
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function formatDateRange(
  startTime: string,
  endTime: string | null | undefined,
  formatDateTime: (date: string) => string,
): string {
  if (!startTime) return ''
  if (!endTime) return formatDateTime(startTime)

  const startDate = new Date(startTime)
  const endDate = new Date(endTime)
  const sameDay =
    startDate.getFullYear() === endDate.getFullYear() &&
    startDate.getMonth() === endDate.getMonth() &&
    startDate.getDate() === endDate.getDate()

  if (sameDay) {
    return `${formatDateTime(startTime)} - ${formatShortTime(endTime)}`
  }

  return `${formatDateTime(startTime)} - ${formatDateTime(endTime)}`
}

export function getLocalDateKey(value: string): string {
  const date = new Date(value)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

export function formatDayHeading(dateKey: string): string {
  return new Date(`${dateKey}T00:00:00`).toLocaleDateString('vi-VN', {
    weekday: 'long',
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  })
}
