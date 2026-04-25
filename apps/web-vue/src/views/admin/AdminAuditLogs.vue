<script setup>
import { computed, onMounted, ref } from 'vue'
import { useAuthStore } from '../../stores/authStore'
import { adminService } from '../../services/adminService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { formatDateTime } from '../../helpers/dateFormatter'
import LoadingSpinner from '../../components/LoadingSpinner.vue'
import EmptyState from '../../components/EmptyState.vue'
import PaginationBar from '../../components/PaginationBar.vue'

const DEFAULT_PAGE_SIZE = 20
const PAGE_SIZE_OPTIONS = [20, 50, 100]

const ENTITY_OPTIONS = [
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

const authStore = useAuthStore()

const auditLogs = ref([])
const totalPages = ref(0)
const currentPage = ref(1)
const totalItems = ref(0)
const pageSize = ref(DEFAULT_PAGE_SIZE)
const isLoading = ref(true)
const errorMessage = ref('')

const searchQuery = ref('')
const actionQuery = ref('')
const selectedEntityType = ref('')
const fromDate = ref('')
const toDate = ref('')

const appliedSearchQuery = ref('')
const appliedActionQuery = ref('')
const appliedEntityType = ref('')
const appliedFromDate = ref('')
const appliedToDate = ref('')

const isSuperAdmin = computed(() => authStore.currentUserRole === 'SUPER_ADMIN')

function normalizeEntityKey(value) {
  return String(value || '').trim().replace(/-/g, '_')
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

function getActorLabel(log) {
  if (!log?.actor_user_id) {
    return 'Hệ thống'
  }

  const roleLabel = ROLE_LABELS[log.actor_role] || log.actor_role || 'Người dùng'
  return `${roleLabel} • ${shortenId(log.actor_user_id)}`
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
  if (path.includes('/students/:student_id/generate-parent-code') || path.endsWith('/generate-parent-code')) {
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

function getFriendlyAction(log) {
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

function getActionTone(log) {
  const action = normalizeText(log?.action).toLowerCase()

  if ((action.includes('delete') || action.includes('cancel') || action.includes('lock') || action.includes('revoke')) && !action.includes('unlock')) {
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

function getEntitySummary(log) {
  const label = getEntityLabel(log?.entity_type)
  const shortId = shortenId(log?.entity_id)
  return shortId ? `${label} • ${shortId}` : label
}

function formatInlineValue(value) {
  if (value === null || value === undefined || value === '') {
    return ''
  }

  if (Array.isArray(value)) {
    return value.map((item) => formatInlineValue(item)).filter(Boolean).join(', ')
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

function getDetailLines(log) {
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

async function fetchAuditLogs(page = 1) {
  if (!isSuperAdmin.value) {
    auditLogs.value = []
    totalPages.value = 0
    totalItems.value = 0
    currentPage.value = 1
    isLoading.value = false
    errorMessage.value = 'Bạn không có quyền truy cập nhật ký hệ thống.'
    return
  }

  isLoading.value = true
  errorMessage.value = ''
  currentPage.value = page

  try {
    const params = {
      limit: pageSize.value,
      offset: (page - 1) * pageSize.value,
      q: appliedSearchQuery.value || undefined,
      action: appliedActionQuery.value || undefined,
      entity_type: appliedEntityType.value || undefined,
      from: appliedFromDate.value ? new Date(appliedFromDate.value).toISOString() : undefined,
      to: appliedToDate.value ? new Date(appliedToDate.value).toISOString() : undefined,
    }

    const data = await adminService.getAuditLogs(params)
    auditLogs.value = Array.isArray(data?.data) ? data.data : []

    if (data?.pagination) {
      totalItems.value = data.pagination.total || 0
      totalPages.value = Math.ceil(totalItems.value / Math.max(data.pagination.limit || pageSize.value, 1)) || 1
      pageSize.value = data.pagination.limit || pageSize.value
    } else {
      totalItems.value = auditLogs.value.length
      totalPages.value = auditLogs.value.length > 0 ? 1 : 0
    }
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể tải dữ liệu nhật ký hệ thống'
  } finally {
    isLoading.value = false
  }
}

function applyFilters() {
  appliedSearchQuery.value = searchQuery.value.trim()
  appliedActionQuery.value = actionQuery.value.trim()
  appliedEntityType.value = selectedEntityType.value
  appliedFromDate.value = fromDate.value
  appliedToDate.value = toDate.value
  fetchAuditLogs(1)
}

function resetFilters() {
  searchQuery.value = ''
  actionQuery.value = ''
  selectedEntityType.value = ''
  fromDate.value = ''
  toDate.value = ''
  appliedSearchQuery.value = ''
  appliedActionQuery.value = ''
  appliedEntityType.value = ''
  appliedFromDate.value = ''
  appliedToDate.value = ''
  fetchAuditLogs(1)
}

function updatePageSize(event) {
  const nextValue = Number(event.target.value)
  pageSize.value = Number.isFinite(nextValue) && nextValue > 0 ? nextValue : DEFAULT_PAGE_SIZE
  fetchAuditLogs(1)
}

onMounted(async () => {
  if (!authStore.currentUser && authStore.isAuthenticated) {
    await authStore.fetchCurrentUser()
  }

  fetchAuditLogs()
})
</script>

<template>
  <div class="admin-audit-logs page-stack">
    <div class="page-actions">
      <button class="btn btn--outline" type="button" :disabled="isLoading" @click="fetchAuditLogs(1)">
        Làm mới
      </button>
    </div>

    <div v-if="!isSuperAdmin" class="card">
      <EmptyState
        title="Không có quyền truy cập"
        message="Chỉ Super Admin mới có thể xem nhật ký hệ thống."
      />
    </div>

    <template v-else>
      <div class="card filters-card">
        <div class="filters-grid">
          <div class="form-group mb-0">
            <label class="form-label" for="auditSearch">Tìm kiếm</label>
            <input
              id="auditSearch"
              v-model="searchQuery"
              type="search"
              class="form-input"
              placeholder="Tìm trong đường dẫn hoặc chi tiết..."
            />
          </div>

          <div class="form-group mb-0">
            <label class="form-label" for="auditAction">Hành động</label>
            <input
              id="auditAction"
              v-model="actionQuery"
              type="search"
              class="form-input"
              placeholder="Ví dụ: appointments.book hoặc POST /admin/users"
            />
          </div>

          <div class="form-group mb-0">
            <label class="form-label" for="entityType">Đối tượng</label>
            <select id="entityType" v-model="selectedEntityType" class="form-input">
              <option v-for="option in ENTITY_OPTIONS" :key="option.value || 'all'" :value="option.value">
                {{ option.label }}
              </option>
            </select>
          </div>

          <div class="form-group mb-0">
            <label class="form-label" for="fromDate">Từ thời điểm</label>
            <input id="fromDate" v-model="fromDate" type="datetime-local" class="form-input" />
          </div>

          <div class="form-group mb-0">
            <label class="form-label" for="toDate">Đến thời điểm</label>
            <input id="toDate" v-model="toDate" type="datetime-local" class="form-input" />
          </div>
        </div>

        <div class="filter-actions">
          <label class="page-size">
            <span class="page-size__label">Hiển thị</span>
            <select class="form-input page-size__select" :value="pageSize" @change="updatePageSize">
              <option v-for="size in PAGE_SIZE_OPTIONS" :key="size" :value="size">
                {{ size }}
              </option>
            </select>
          </label>

          <div class="filter-actions__buttons">
            <button class="btn btn--outline" type="button" @click="resetFilters">
              Xóa lọc
            </button>
            <button class="btn btn--primary" type="button" @click="applyFilters">
              Áp dụng
            </button>
          </div>
        </div>
      </div>

      <div v-if="errorMessage" class="alert alert--error">
        <p class="font-bold">Lỗi tải dữ liệu</p>
        <p>{{ errorMessage }}</p>
        <button class="btn btn--outline mt-2" type="button" @click="fetchAuditLogs(currentPage)">
          Thử lại
        </button>
      </div>

      <LoadingSpinner v-else-if="isLoading" message="Đang tải dữ liệu nhật ký..." />

      <div v-else-if="auditLogs.length === 0" class="card">
        <EmptyState
          title="Chưa có bản ghi nào"
          message="Hệ thống chưa ghi nhận hoạt động phù hợp với bộ lọc hiện tại."
        />
      </div>

      <template v-else>
        <div class="card desktop-table">
          <div class="table-responsive">
            <table class="table">
              <thead>
                <tr>
                  <th>Thời gian</th>
                  <th>Người thực hiện</th>
                  <th>Hành động</th>
                  <th>Đối tượng</th>
                  <th>Chi tiết</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="log in auditLogs" :key="log.audit_log_id">
                  <td class="whitespace-nowrap">{{ formatDateTime(log.created_at) }}</td>
                  <td>{{ getActorLabel(log) }}</td>
                  <td>
                    <span :class="getActionTone(log)">
                      {{ getFriendlyAction(log) }}
                    </span>
                  </td>
                  <td>{{ getEntitySummary(log) }}</td>
                  <td>
                    <ul class="detail-list">
                      <li v-for="line in getDetailLines(log)" :key="line">{{ line }}</li>
                    </ul>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="mobile-list">
          <article v-for="log in auditLogs" :key="log.audit_log_id" class="card audit-card">
            <div class="audit-card__head">
              <p class="audit-card__time">{{ formatDateTime(log.created_at) }}</p>
              <span :class="getActionTone(log)">{{ getFriendlyAction(log) }}</span>
            </div>

            <p class="audit-card__meta">Người thực hiện: {{ getActorLabel(log) }}</p>
            <p class="audit-card__meta">Đối tượng: {{ getEntitySummary(log) }}</p>

            <ul class="detail-list">
              <li v-for="line in getDetailLines(log)" :key="line">{{ line }}</li>
            </ul>
          </article>
        </div>

        <PaginationBar
          :current-page="currentPage"
          :total-pages="totalPages"
          :total-items="totalItems"
          :limit="pageSize"
          @page-change="fetchAuditLogs"
        />
      </template>
    </template>
  </div>
</template>

<style scoped>
.page-stack,
.mobile-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.page-actions {
  display: flex;
  justify-content: flex-end;
}

.filters-card,
.desktop-table,
.audit-card {
  padding: var(--spacing-4);
}

.filters-grid {
  display: grid;
  gap: var(--spacing-3);
  grid-template-columns: repeat(1, minmax(0, 1fr));
}

.filter-actions {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
  margin-top: var(--spacing-3);
}

.filter-actions__buttons {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-2);
}

.page-size {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-2);
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.page-size__label {
  white-space: nowrap;
}

.page-size__select {
  width: 96px;
}

.table-responsive {
  overflow-x: auto;
}

.table {
  width: 100%;
  border-collapse: collapse;
}

.table th,
.table td {
  padding: var(--spacing-3) var(--spacing-4);
  text-align: left;
  border-bottom: 1px solid var(--color-border);
  vertical-align: top;
}

.table th {
  background-color: var(--color-background);
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  font-weight: 600;
  text-transform: uppercase;
}

.table tbody tr:hover {
  background-color: var(--color-background);
}

.detail-list {
  margin: 0;
  padding-left: 1rem;
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.detail-list li + li {
  margin-top: 0.3rem;
}

.whitespace-nowrap {
  white-space: nowrap;
}

.mobile-list {
  display: none;
}

.audit-card__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--spacing-3);
}

.audit-card__time,
.audit-card__meta {
  margin: 0;
}

.audit-card__time {
  font-weight: 700;
  color: var(--color-text);
}

.audit-card__meta {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.mt-2 {
  margin-top: var(--spacing-2);
}

@media (min-width: 768px) {
  .filters-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .filter-actions {
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
  }
}

@media (min-width: 1200px) {
  .filters-grid {
    grid-template-columns: repeat(5, minmax(0, 1fr));
  }
}

@media (max-width: 767px) {
  .desktop-table {
    display: none;
  }

  .mobile-list {
    display: flex;
  }
}
</style>
