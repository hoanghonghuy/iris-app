import { formatDateTime } from '../../../helpers/dateFormatter'

const APPOINTMENT_STATUS_MAP = {
  pending: { text: 'Chờ xác nhận', badgeClass: 'badge badge--outline' },
  confirmed: { text: 'Đã xác nhận', badgeClass: 'badge badge--success' },
  cancelled: { text: 'Đã hủy', badgeClass: 'badge badge--danger' },
  completed: { text: 'Đã hoàn thành', badgeClass: 'badge badge--primary' },
  no_show: { text: 'Vắng mặt', badgeClass: 'badge badge--outline' },
}

const CANCEL_REASON_MAP = {
  parent_cancelled: 'Phụ huynh đã hủy lịch',
  teacher_cancelled: 'Giáo viên đã hủy lịch',
  system_cancelled: 'Hệ thống đã hủy lịch',
}

function csvEscape(value) {
  const text = String(value ?? '')
  if (/[",\n]/.test(text)) {
    return `"${text.replace(/"/g, '""')}"`
  }
  return text
}

function triggerCsvDownload(filename, headers, rows) {
  const lines = [
    headers.map(csvEscape).join(','),
    ...rows.map((row) => row.map(csvEscape).join(',')),
  ]

  const blob = new Blob([`\uFEFF${lines.join('\n')}`], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

export function daysAgo(days) {
  const date = new Date()
  date.setDate(date.getDate() - days)
  return date
}

export function getDateInputValue(date) {
  const local = new Date(date.getTime() - date.getTimezoneOffset() * 60000)
  return local.toISOString().slice(0, 10)
}

export function getTimezoneDisplay() {
  const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone || 'Local'
  const totalMinutes = -new Date().getTimezoneOffset()
  const sign = totalMinutes >= 0 ? '+' : '-'
  const abs = Math.abs(totalMinutes)
  const hours = String(Math.floor(abs / 60)).padStart(2, '0')
  const minutes = String(abs % 60).padStart(2, '0')
  return `${timezone} (UTC${sign}${hours}:${minutes})`
}

function formatShortTime(value) {
  if (!value) return ''
  return new Date(value).toLocaleTimeString('vi-VN', {
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function formatDateRange(startTime, endTime) {
  if (!startTime) return ''
  if (!endTime) return formatDateTime(startTime)

  const startDate = new Date(startTime)
  const endDate = new Date(endTime)
  const sameDay =
    startDate.getFullYear() === endDate.getFullYear()
    && startDate.getMonth() === endDate.getMonth()
    && startDate.getDate() === endDate.getDate()

  if (sameDay) {
    return `${formatDateTime(startTime)} - ${formatShortTime(endTime)}`
  }

  return `${formatDateTime(startTime)} - ${formatDateTime(endTime)}`
}

export function getStatusBadge(status) {
  return APPOINTMENT_STATUS_MAP[status]?.badgeClass || 'badge badge--outline'
}

export function getStatusText(status) {
  return APPOINTMENT_STATUS_MAP[status]?.text || status || 'Không xác định'
}

export function getCancelReasonText(reason) {
  if (!reason) return ''
  return CANCEL_REASON_MAP[reason] || 'Lý do khác'
}

export function exportAppointmentsToCsv({ appointments, historyView, timezoneDisplay }) {
  if (!Array.isArray(appointments) || appointments.length === 0) {
    return false
  }

  const rows = [...appointments]
    .sort((left, right) => new Date(left.start_time).getTime() - new Date(right.start_time).getTime())
    .map((item) => [
      item.student_name || item.student_id || 'N/A',
      item.teacher_name || item.teacher_id || 'N/A',
      getStatusText(item.status),
      formatDateTime(item.start_time),
      formatDateTime(item.end_time),
      timezoneDisplay || getTimezoneDisplay(),
      item.note || '',
      getCancelReasonText(item.cancel_reason),
    ])

  triggerCsvDownload(
    `parent-appointments-${historyView || 'history'}-${getDateInputValue(new Date())}.csv`,
    ['HocSinh', 'GiaoVien', 'TrangThai', 'BatDau', 'KetThuc', 'MuiGio', 'GhiChu', 'LyDoHuy'],
    rows,
  )

  return true
}
