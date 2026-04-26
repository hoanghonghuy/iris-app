const ROLE_LABELS = {
  SUPER_ADMIN: 'Super Admin',
  SCHOOL_ADMIN: 'School Admin',
  TEACHER: 'Giáo viên',
  PARENT: 'Phụ huynh',
}

const ENTITY_LABELS = {
  admin: 'Quản trị hệ thống',
  schools: 'Trường học',
  classes: 'Lớp học',
  students: 'Học sinh',
  teachers: 'Giáo viên',
  parents: 'Phụ huynh',
  users: 'Người dùng',
  school_admins: 'Quản trị trường',
  appointments: 'Lịch hẹn',
  appointment_slots: 'Khung giờ hẹn',
  posts: 'Bài đăng',
  audit_logs: 'Nhật ký hệ thống',
  messages: 'Tin nhắn',
  conversations: 'Cuộc trò chuyện',
  health_records: 'Sổ sức khỏe',
  attendance: 'Điểm danh',
}

const DIRECT_ACTION_LABELS = {
  'appointments.book': 'Đặt lịch hẹn',
  'appointments.cancel': 'Hủy lịch hẹn',
  'appointments.confirm': 'Xác nhận lịch hẹn',
  'appointments.complete': 'Hoàn tất lịch hẹn',
  'appointments.no_show': 'Đánh dấu không đến',
  'appointments.slot.create': 'Tạo khung giờ hẹn',
  'appointment_slots.create': 'Tạo khung giờ hẹn',
  'appointment_slots.delete': 'Xóa khung giờ hẹn',
  'posts.create': 'Tạo bài đăng',
  'posts.update': 'Cập nhật bài đăng',
  'posts.delete': 'Xóa bài đăng',
  'audit_logs.list': 'Xem nhật ký hệ thống',
}

const STATUS_LABELS = {
  pending: 'Chờ xác nhận',
  confirmed: 'Đã xác nhận',
  completed: 'Đã hoàn tất',
  cancelled: 'Đã hủy',
  canceled: 'Đã hủy',
  no_show: 'Không đến',
  active: 'Đang hoạt động',
  locked: 'Đã khóa',
  unlocked: 'Đã mở khóa',
}

export const AUDIT_ENTITY_OPTIONS = [
  { value: '', label: 'Tất cả đối tượng' },
  { value: 'schools', label: 'Trường học' },
  { value: 'classes', label: 'Lớp học' },
  { value: 'students', label: 'Học sinh' },
  { value: 'teachers', label: 'Giáo viên' },
  { value: 'parents', label: 'Phụ huynh' },
  { value: 'users', label: 'Người dùng' },
  { value: 'school-admins', label: 'Quản trị trường' },
  { value: 'appointments', label: 'Lịch hẹn' },
  { value: 'appointment_slots', label: 'Khung giờ hẹn' },
  { value: 'posts', label: 'Bài đăng' },
  { value: 'audit_logs', label: 'Nhật ký hệ thống' },
]

function normalizeEntityKey(value) {
  return String(value || '')
    .trim()
    .replace(/-/g, '_')
}

function normalizeText(value) {
  return String(value || '').trim()
}

function shortenId(value) {
  const text = normalizeText(value)
  if (!text) {
    return ''
  }

  return text.length > 8 ? text.slice(0, 8) : text
}

function getEntityLabel(entityType) {
  const key = normalizeEntityKey(entityType)
  return ENTITY_LABELS[key] || entityType || 'Đối tượng khác'
}

function parseHttpAction(action) {
  const match = normalizeText(action).match(/^(GET|POST|PUT|PATCH|DELETE)\s+(.+)$/)
  if (!match) {
    return null
  }

  return {
    method: match[1],
    path: match[2],
  }
}

function getFriendlyHttpAction(method, path, entityType) {
  if (
    path.includes('/students/:student_id/generate-parent-code') ||
    path.endsWith('/generate-parent-code')
  ) {
    return 'Tạo mã phụ huynh'
  }

  if (path.includes('/students/:student_id/parent-code') || path.endsWith('/parent-code')) {
    return 'Thu hồi mã phụ huynh'
  }

  if (path.includes('/teachers/:teacher_id/classes/:class_id')) {
    return method === 'DELETE' ? 'Gỡ lớp khỏi giáo viên' : 'Gán lớp cho giáo viên'
  }

  if (path.includes('/parents/:parent_id/students/:student_id')) {
    return method === 'DELETE' ? 'Gỡ học sinh khỏi phụ huynh' : 'Gán học sinh cho phụ huynh'
  }

  if (path.includes('/users/:user_id/lock') || path.endsWith('/lock')) {
    return 'Khóa tài khoản'
  }

  if (path.includes('/users/:user_id/unlock') || path.endsWith('/unlock')) {
    return 'Mở khóa tài khoản'
  }

  const entityLabel = getEntityLabel(entityType).toLowerCase()

  if (method === 'POST') {
    return `Tạo ${entityLabel}`
  }

  if (method === 'PUT' || method === 'PATCH') {
    return `Cập nhật ${entityLabel}`
  }

  if (method === 'DELETE') {
    return `Xóa ${entityLabel}`
  }

  return normalizeText(path) || 'Hoạt động hệ thống'
}

function getEventEntityLabel(prefix) {
  const normalizedPrefix = normalizeEntityKey(prefix)

  if (normalizedPrefix === 'appointments.slot') {
    return 'khung giờ hẹn'
  }

  return getEntityLabel(normalizedPrefix).toLowerCase()
}

function getFriendlyEventAction(action) {
  if (DIRECT_ACTION_LABELS[action]) {
    return DIRECT_ACTION_LABELS[action]
  }

  const parts = action.split('.').filter(Boolean)
  if (parts.length < 2) {
    return action
  }

  const verb = parts.at(-1)
  const entityPrefix = parts.slice(0, -1).join('.')
  const entityLabel = getEventEntityLabel(entityPrefix)

  if (verb === 'create') {
    return `Tạo ${entityLabel}`
  }

  if (verb === 'update') {
    return `Cập nhật ${entityLabel}`
  }

  if (verb === 'delete') {
    return `Xóa ${entityLabel}`
  }

  if (verb === 'list') {
    return `Xem danh sách ${entityLabel}`
  }

  if (verb === 'book') {
    return `Đặt ${entityLabel}`
  }

  if (verb === 'cancel') {
    return `Hủy ${entityLabel}`
  }

  if (verb === 'confirm') {
    return `Xác nhận ${entityLabel}`
  }

  if (verb === 'complete') {
    return `Hoàn tất ${entityLabel}`
  }

  if (verb === 'lock') {
    return `Khóa ${entityLabel}`
  }

  if (verb === 'unlock') {
    return `Mở khóa ${entityLabel}`
  }

  return action
}

function formatInlineValue(value) {
  if (value === null || value === undefined || value === '') {
    return ''
  }

  if (Array.isArray(value)) {
    return value
      .map((item) => formatInlineValue(item))
      .filter(Boolean)
      .join(', ')
  }

  if (typeof value === 'object') {
    return Object.entries(value)
      .filter(([, itemValue]) => itemValue !== null && itemValue !== undefined && itemValue !== '')
      .map(([key, itemValue]) => `${key}=${formatInlineValue(itemValue)}`)
      .join(', ')
  }

  return String(value)
}

function formatQueryValue(query) {
  if (!query) {
    return ''
  }

  if (typeof query === 'string') {
    const queryText = query.trim()
    if (!queryText) {
      return ''
    }

    const params = new URLSearchParams(queryText)
    const parts = []
    params.forEach((value, key) => {
      parts.push(`${key}=${value}`)
    })

    return parts.length > 0 ? parts.join(', ') : queryText
  }

  return formatInlineValue(query)
}

function formatStatusLine(status) {
  const value = normalizeText(status)
  if (!value) {
    return ''
  }

  if (/^\d{3}$/.test(value)) {
    return `Mã phản hồi: HTTP ${value}`
  }

  return `Trạng thái: ${STATUS_LABELS[value] || value}`
}

export function getAuditActorLabel(log) {
  if (!log?.actor_user_id) {
    return 'Hệ thống'
  }

  const roleLabel = ROLE_LABELS[log.actor_role] || log.actor_role || 'Người dùng'
  return roleLabel
}

export function getAuditFriendlyAction(log) {
  const action = normalizeText(log?.action)
  if (!action) {
    return 'Hoạt động hệ thống'
  }

  if (DIRECT_ACTION_LABELS[action]) {
    return DIRECT_ACTION_LABELS[action]
  }

  const httpAction = parseHttpAction(action)
  if (httpAction) {
    return getFriendlyHttpAction(httpAction.method, httpAction.path, log?.entity_type)
  }

  if (action.includes('.')) {
    return getFriendlyEventAction(action)
  }

  return action
}

export function getAuditActionTone(log) {
  const action = normalizeText(log?.action).toLowerCase()

  if (
    (action.includes('delete') ||
      action.includes('cancel') ||
      action.includes('lock') ||
      action.includes('revoke')) &&
    !action.includes('unlock')
  ) {
    return 'badge badge--danger'
  }

  if (
    action.includes('create') ||
    action.includes('book') ||
    action.includes('confirm') ||
    action.includes('complete') ||
    action.includes('unlock') ||
    action.includes('generate')
  ) {
    return 'badge badge--success'
  }

  if (action.includes('update') || action.includes('patch') || action.includes('edit')) {
    return 'badge badge--primary'
  }

  return 'badge badge--outline'
}

export function getAuditEntitySummary(log) {
  const label = getEntityLabel(log?.entity_type)
  return label
}

export function getAuditDetailLines(log) {
  const details = log?.details && typeof log.details === 'object' ? log.details : {}
  const lines = []

  const statusLine = formatStatusLine(details.status)
  if (statusLine) {
    lines.push(statusLine)
  }

  const requestPath = normalizeText(details.request_path)
  if (requestPath) {
    lines.push(`Đường dẫn: ${requestPath}`)
  } else {
    const route = normalizeText(details.route)
    if (route) {
      lines.push(`Tuyến xử lý: ${route}`)
    }
  }

  const queryText = formatQueryValue(details.query)
  if (queryText) {
    lines.push(`Tham số: ${queryText}`)
  }

  const schoolId = shortenId(details.school_id)
  if (schoolId) {
    lines.push(`Trường áp dụng: ${schoolId}`)
  }

  return lines.length > 0 ? lines : ['Không có chi tiết bổ sung']
}
