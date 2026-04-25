<script setup>
import { computed, onMounted, ref } from 'vue'
import {
  AlertTriangle,
  CalendarDays,
  Clock3,
  Download,
  LoaderCircle,
  RefreshCcw,
  UserRound,
} from 'lucide-vue-next'
import { teacherService } from '../../services/teacherService'
import { extractErrorMessage } from '../../helpers/errorHandler'

const FETCH_LIMIT = 100

const statusConfig = {
  pending: { label: 'Chờ xác nhận', badge: 'badge--outline' },
  confirmed: { label: 'Đã xác nhận', badge: 'badge--primary' },
  cancelled: { label: 'Đã hủy', badge: 'badge--danger' },
  completed: { label: 'Hoàn tất', badge: 'badge--outline' },
  no_show: { label: 'Vắng mặt', badge: 'badge--outline' },
}

const statusOptions = [
  { value: 'pending', label: 'Chờ xác nhận' },
  { value: 'confirmed', label: 'Đã xác nhận' },
  { value: 'cancelled', label: 'Đã hủy' },
  { value: 'completed', label: 'Hoàn tất' },
  { value: 'no_show', label: 'Vắng mặt' },
]

const classes = ref([])
const appointments = ref([])
const loading = ref(true)
const errorMessage = ref('')
const submitting = ref(false)
const updatingAppointmentId = ref(null)
const showCreateForm = ref(false)

const classId = ref('')
const startTime = ref('')
const durationMinutes = ref(30)
const bufferMinutes = ref(10)
const maxBookingsPerDay = ref(12)
const note = ref('')

const statusFilter = ref('')
const filterFromDate = ref(toDateInputValue(offsetDate(-6)))
const filterToDate = ref(toDateInputValue(new Date()))

const timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone || 'Local'
const utcOffsetLabel = getUtcOffsetLabel()
const timezoneDisplay = `${timeZone} (${utcOffsetLabel})`

const minStartTime = computed(() => {
  const local = new Date(Date.now() - new Date().getTimezoneOffset() * 60000)
  return local.toISOString().slice(0, 16)
})

const stats = computed(() => ({
  totalClasses: classes.value.length,
  totalAppointments: appointments.value.length,
  pendingCount: appointments.value.filter((item) => item.status === 'pending').length,
  confirmedCount: appointments.value.filter((item) => item.status === 'confirmed').length,
}))

const groupedAppointments = computed(() => {
  const groups = new Map()
  const sorted = [...appointments.value].sort(
    (a, b) => new Date(a.start_time).getTime() - new Date(b.start_time).getTime(),
  )

  for (const appointment of sorted) {
    const key = getLocalDateKey(appointment.start_time)
    if (!groups.has(key)) groups.set(key, [])
    groups.get(key).push(appointment)
  }

  return Array.from(groups.entries()).map(([dateKey, items]) => ({ dateKey, items }))
})

function offsetDate(days) {
  const date = new Date()
  date.setDate(date.getDate() + days)
  return date
}

function toDateInputValue(date) {
  const local = new Date(date.getTime() - date.getTimezoneOffset() * 60000)
  return local.toISOString().slice(0, 10)
}

function getUtcOffsetLabel() {
  const totalMinutes = -new Date().getTimezoneOffset()
  const sign = totalMinutes >= 0 ? '+' : '-'
  const abs = Math.abs(totalMinutes)
  const hours = String(Math.floor(abs / 60)).padStart(2, '0')
  const mins = String(abs % 60).padStart(2, '0')
  return `UTC${sign}${hours}:${mins}`
}

function getLocalDateKey(value) {
  const date = new Date(value)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function formatDayHeading(dateKey) {
  return new Date(`${dateKey}T00:00:00`).toLocaleDateString('vi-VN', {
    weekday: 'long',
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  })
}

function formatDateTime(value) {
  return new Date(value).toLocaleString('vi-VN', {
    hour: '2-digit',
    minute: '2-digit',
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    timeZoneName: 'short',
  })
}

function csvEscape(value) {
  const text = String(value ?? '')
  if (/[",\n]/.test(text)) {
    return `"${text.replace(/"/g, '""')}"`
  }
  return text
}

function downloadCsv(filename, headers, rows) {
  const lines = [headers.map(csvEscape).join(','), ...rows.map((row) => row.map(csvEscape).join(','))]
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

function normalizeListResponse(value) {
  if (Array.isArray(value?.data)) return value.data.filter(Boolean)
  if (Array.isArray(value)) return value.filter(Boolean)
  return []
}

function normalizePaginatedResponse(value) {
  return {
    items: normalizeListResponse(value),
    pagination: value?.pagination ?? {
      total: normalizeListResponse(value).length,
      limit: FETCH_LIMIT,
      offset: 0,
      has_more: false,
    },
  }
}

async function fetchAllAppointments(params = {}) {
  let offset = 0
  let hasMore = true
  let combined = []

  while (hasMore) {
    const response = await teacherService.getAppointments({
      ...params,
      limit: FETCH_LIMIT,
      offset,
    })

    const { items, pagination } = normalizePaginatedResponse(response)
    combined = combined.concat(items)
    hasMore = Boolean(pagination.has_more) && items.length > 0

    const nextOffset = offset + (pagination.limit || FETCH_LIMIT)
    if (!hasMore || nextOffset >= (pagination.total ?? combined.length)) {
      break
    }

    offset = nextOffset
  }

  return combined
}

function getCancelReasonText(reason) {
  const map = {
    parent_cancelled: 'Phụ huynh đã hủy lịch',
    teacher_cancelled: 'Giáo viên đã hủy lịch',
    system_cancelled: 'Hệ thống đã hủy lịch',
  }

  return map[reason] || reason || ''
}

async function loadData() {
  loading.value = true
  errorMessage.value = ''

  try {
    const from = filterFromDate.value ? new Date(`${filterFromDate.value}T00:00:00`).toISOString() : undefined
    const to = filterToDate.value ? new Date(`${filterToDate.value}T23:59:59.999`).toISOString() : undefined

    if (from && to && new Date(from).getTime() > new Date(to).getTime()) {
      errorMessage.value = 'Khoảng ngày lọc không hợp lệ: Từ ngày phải nhỏ hơn hoặc bằng Đến ngày.'
      appointments.value = []
      return
    }

    const [classResponse, appointmentResponse] = await Promise.all([
      teacherService.getMyClasses(),
      fetchAllAppointments({
        status: statusFilter.value || undefined,
        from,
        to,
      }),
    ])

    classes.value = normalizeListResponse(classResponse)
    appointments.value = appointmentResponse
    if (!classId.value && classes.value.length) {
      classId.value = classes.value[0].class_id
    }
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể tải dữ liệu lịch hẹn.'
  } finally {
    loading.value = false
  }
}

async function createSlot() {
  if (!classId.value || !startTime.value) return

  submitting.value = true
  try {
    const startDate = new Date(startTime.value)
    if (Number.isNaN(startDate.getTime())) {
      errorMessage.value = 'Thời gian bắt đầu không hợp lệ.'
      return
    }

    const dayStart = new Date(startDate)
    dayStart.setHours(0, 0, 0, 0)
    const dayEnd = new Date(startDate)
    dayEnd.setHours(23, 59, 59, 999)

    const activeAppointments = (await fetchAllAppointments({
      from: dayStart.toISOString(),
      to: dayEnd.toISOString(),
    })).filter((item) => item.status !== 'cancelled')

    if (activeAppointments.length >= maxBookingsPerDay.value) {
      errorMessage.value = `Đã đạt giới hạn ${maxBookingsPerDay.value} lịch trong ngày này.`
      return
    }

    const proposedStartMs = startDate.getTime()
    const proposedEndMs = proposedStartMs + Number(durationMinutes.value) * 60000
    const bufferMs = Math.max(0, Number(bufferMinutes.value)) * 60000
    const conflicting = activeAppointments.find((appointment) => {
      const existingStart = new Date(appointment.start_time).getTime()
      const existingEnd = new Date(appointment.end_time).getTime()
      return !(proposedEndMs + bufferMs <= existingStart || proposedStartMs >= existingEnd + bufferMs)
    })

    if (conflicting) {
      errorMessage.value = `Khung giờ mới chưa đảm bảo khoảng nghỉ ${bufferMinutes.value} phút với lịch ${formatDateTime(conflicting.start_time)}.`
      return
    }

    await teacherService.createAppointmentSlot({
      class_id: classId.value,
      start_time: startDate.toISOString(),
      duration_minutes: Number(durationMinutes.value),
      buffer_minutes: Number(bufferMinutes.value),
      max_bookings_per_day: Number(maxBookingsPerDay.value),
      note: note.value.trim() || undefined,
    })

    startTime.value = ''
    note.value = ''
    showCreateForm.value = false
    await loadData()
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể tạo khung giờ.'
  } finally {
    submitting.value = false
  }
}

async function updateStatus(appointmentId, status) {
  updatingAppointmentId.value = appointmentId
  try {
    await teacherService.updateAppointmentStatus(
      appointmentId,
      status,
      status === 'cancelled' ? 'teacher_cancelled' : undefined,
    )
    await loadData()
  } catch (error) {
    errorMessage.value = extractErrorMessage(error) || 'Không thể cập nhật trạng thái lịch hẹn.'
  } finally {
    updatingAppointmentId.value = null
  }
}

async function resetLastSevenDays() {
  filterFromDate.value = toDateInputValue(offsetDate(-6))
  filterToDate.value = toDateInputValue(new Date())
  await loadData()
}

function exportAppointmentsCsv() {
  if (appointments.value.length === 0) {
    errorMessage.value = 'Không có dữ liệu để xuất CSV.'
    return
  }

  const rows = [...appointments.value]
    .sort((a, b) => new Date(a.start_time).getTime() - new Date(b.start_time).getTime())
    .map((appointment) => [
      formatDayHeading(getLocalDateKey(appointment.start_time)),
      appointment.student_name || appointment.student_id,
      appointment.class_name || appointment.class_id,
      appointment.parent_name || appointment.parent_id,
      statusConfig[appointment.status]?.label || appointment.status,
      formatDateTime(appointment.start_time),
      formatDateTime(appointment.end_time),
      timezoneDisplay,
      appointment.note || '',
      appointment.cancel_reason ? getCancelReasonText(appointment.cancel_reason) : '',
    ])

  downloadCsv(
    `teacher-appointments-${toDateInputValue(new Date())}.csv`,
    ['Ngay', 'HocSinh', 'Lop', 'PhuHuynh', 'TrangThai', 'BatDau', 'KetThuc', 'MuiGio', 'GhiChu', 'LyDoHuy'],
    rows,
  )
}

onMounted(loadData)
</script>

<template>
  <div class="teacher-appointments page-stack">
    <section class="page-heading">
      <p class="timezone">Múi giờ hiển thị: {{ timezoneDisplay }}</p>
    </section>

    <div v-if="errorMessage" class="alert alert--error alert-row">
      <AlertTriangle :size="16" />
      <div>
        <strong>Tải dữ liệu thất bại</strong>
        <p>{{ errorMessage }}</p>
      </div>
      <button class="btn btn--outline btn--sm" @click="loadData">Thử lại</button>
    </div>

    <section class="card stats-grid">
      <div><b>Số lớp đang phụ trách:</b> {{ stats.totalClasses }}</div>
      <div><b>Tổng lịch hẹn:</b> {{ stats.totalAppointments }}</div>
      <div><b>Đang chờ xác nhận:</b> {{ stats.pendingCount }}</div>
      <div><b>Đã xác nhận:</b> {{ stats.confirmedCount }}</div>
    </section>

    <section class="card create-summary">
      <div>
        <h2>Tạo khung giờ mới</h2>
        <p>Mở biểu mẫu để thêm khung giờ và kiểm tra khoảng nghỉ trước khi tạo.</p>
      </div>
      <div class="action-row">
        <button class="btn btn--primary" @click="showCreateForm = !showCreateForm">
          {{ showCreateForm ? 'Đóng biểu mẫu' : 'Tạo khung giờ mới' }}
        </button>
        <button class="btn btn--outline" :disabled="loading" @click="loadData">
          <RefreshCcw :size="16" />
          Làm mới
        </button>
      </div>
    </section>

    <section v-if="showCreateForm" class="card form-card">
      <div class="card__header">
        <h2 class="section-title">Tạo khung giờ mới</h2>
      </div>
      <div class="card__body form-stack">
        <div class="form-grid four">
          <div class="form-group mb-0">
            <label class="form-label">Lớp học</label>
            <select v-model="classId" class="form-input">
              <option v-for="classInfo in classes" :key="classInfo.class_id" :value="classInfo.class_id">
                {{ classInfo.name }}
              </option>
            </select>
          </div>

          <div class="form-group mb-0">
            <label class="form-label">Thời gian bắt đầu</label>
            <input v-model="startTime" class="form-input" type="datetime-local" :min="minStartTime" />
          </div>

          <div class="form-group mb-0">
            <label class="form-label">Thời lượng</label>
            <select v-model="durationMinutes" class="form-input">
              <option :value="15">15 phút</option>
              <option :value="20">20 phút</option>
              <option :value="30">30 phút</option>
              <option :value="45">45 phút</option>
              <option :value="60">60 phút</option>
            </select>
          </div>

          <div class="form-group mb-0">
            <label class="form-label">Kết thúc dự kiến</label>
            <div class="readonly-field">
              <Clock3 :size="16" />
              {{ startTime ? formatDateTime(new Date(new Date(startTime).getTime() + durationMinutes * 60000).toISOString()) : 'Chưa xác định' }}
            </div>
          </div>
        </div>

        <div class="form-grid two">
          <div class="form-group mb-0">
            <label class="form-label">Khoảng nghỉ giữa 2 lịch (phút)</label>
            <input v-model.number="bufferMinutes" class="form-input" type="number" min="0" step="5" />
            <p class="hint">Hệ thống sẽ giữ khoảng nghỉ này trước hoặc sau mỗi lịch.</p>
          </div>

          <div class="form-group mb-0">
            <label class="form-label">Số lịch tối đa trong ngày</label>
            <input v-model.number="maxBookingsPerDay" class="form-input" type="number" min="1" step="1" />
          </div>
        </div>

        <div class="form-group mb-0">
          <label class="form-label">Ghi chú cho phụ huynh</label>
          <textarea v-model="note" class="form-input" rows="3" maxlength="500" placeholder="Ví dụ: Vui lòng chuẩn bị các câu hỏi liên quan tới tiến độ học tập của bé..."></textarea>
          <p class="hint">Tối đa 500 ký tự.</p>
        </div>

        <div class="form-actions">
          <button class="btn btn--outline" @click="showCreateForm = false">Đóng</button>
          <button class="btn btn--primary" :disabled="submitting || !classId || !startTime" @click="createSlot">
            {{ submitting ? 'Đang tạo...' : 'Tạo khung giờ' }}
          </button>
        </div>
      </div>
    </section>

    <section class="card list-card">
      <div class="list-header">
        <h2>Danh sách lịch hẹn</h2>
        <div class="filters">
          <select v-model="statusFilter" class="form-input" @change="loadData">
            <option value="">Tất cả trạng thái</option>
            <option v-for="option in statusOptions" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>

          <div class="date-range">
            <input v-model="filterFromDate" class="form-input" type="date" @change="loadData" />
            <span>đến</span>
            <input v-model="filterToDate" class="form-input" type="date" @change="loadData" />
          </div>

          <button class="btn btn--outline btn--sm" @click="resetLastSevenDays">7 ngày gần nhất</button>
          <button class="btn btn--outline btn--sm" @click="exportAppointmentsCsv">
            <Download :size="14" />
            Xuất CSV
          </button>
        </div>
      </div>

      <div v-if="loading" class="loading-block">
        <LoaderCircle class="spin text-muted" :size="28" />
      </div>

      <p v-else-if="appointments.length === 0" class="text-sm text-muted">Chưa có lịch hẹn.</p>

      <div v-else class="appointment-list">
        <div v-for="group in groupedAppointments" :key="group.dateKey" class="date-group">
          <div class="date-heading">
            <CalendarDays :size="14" />
            {{ formatDayHeading(group.dateKey) }}
          </div>

          <article v-for="appointment in group.items" :key="appointment.appointment_id" class="appointment-item">
            <div class="appointment-copy">
              <p class="appointment-title">
                {{ appointment.student_name || appointment.student_id }} - {{ appointment.class_name || appointment.class_id }}
              </p>
              <p>{{ formatDateTime(appointment.start_time) }} - {{ formatDateTime(appointment.end_time) }}</p>
              <p class="text-xs">Múi giờ: {{ timezoneDisplay }}</p>
              <p class="inline-meta">
                <UserRound :size="14" />
                Phụ huynh: {{ appointment.parent_name || appointment.parent_id }}
              </p>
              <p v-if="appointment.note">Ghi chú: {{ appointment.note }}</p>
              <span class="badge" :class="statusConfig[appointment.status]?.badge || 'badge--outline'">
                {{ statusConfig[appointment.status]?.label || appointment.status }}
              </span>
            </div>

            <div class="status-actions">
              <button class="btn btn--sm btn--primary" :disabled="appointment.status !== 'pending' || updatingAppointmentId === appointment.appointment_id" @click="updateStatus(appointment.appointment_id, 'confirmed')">
                Xác nhận
              </button>
              <button class="btn btn--sm btn--outline" :disabled="appointment.status !== 'confirmed' || updatingAppointmentId === appointment.appointment_id" @click="updateStatus(appointment.appointment_id, 'completed')">
                Hoàn tất
              </button>
              <button class="btn btn--sm btn--outline" :disabled="appointment.status !== 'confirmed' || updatingAppointmentId === appointment.appointment_id" @click="updateStatus(appointment.appointment_id, 'no_show')">
                Vắng mặt
              </button>
              <button class="btn btn--sm btn--danger" :disabled="['cancelled', 'completed'].includes(appointment.status) || updatingAppointmentId === appointment.appointment_id" @click="updateStatus(appointment.appointment_id, 'cancelled')">
                {{ updatingAppointmentId === appointment.appointment_id ? 'Đang xử lý...' : 'Hủy' }}
              </button>
            </div>
          </article>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.page-stack,
.form-stack,
.appointment-list,
.date-group {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.page-heading h1,
.page-heading p,
.create-summary h2,
.create-summary p,
.list-header h2,
.appointment-title {
  margin: 0;
}

.page-heading h1 {
  font-size: var(--font-size-2xl);
}

.page-heading p,
.create-summary p,
.appointment-copy {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.timezone {
  font-size: var(--font-size-xs);
}

.alert-row,
.create-summary,
.action-row,
.filters,
.date-range,
.inline-meta,
.date-heading,
.status-actions {
  display: flex;
  align-items: center;
}

.alert-row {
  gap: var(--spacing-3);
}

.alert-row > div {
  flex: 1;
}

.stats-grid {
  display: grid;
  gap: var(--spacing-2);
  padding: var(--spacing-4);
  font-size: var(--font-size-sm);
}

.create-summary {
  justify-content: space-between;
  gap: var(--spacing-4);
  padding: var(--spacing-4);
}

.action-row,
.filters,
.status-actions {
  flex-wrap: wrap;
  gap: var(--spacing-2);
}

.section-title,
.list-header h2 {
  font-size: var(--font-size-lg);
}

.form-grid {
  display: grid;
  gap: var(--spacing-3);
}

.readonly-field {
  min-height: 2.25rem;
  display: flex;
  align-items: center;
  gap: var(--spacing-2);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-background);
  color: var(--color-text-muted);
  padding: var(--spacing-2) var(--spacing-3);
  font-size: var(--font-size-sm);
}

.hint {
  margin-top: var(--spacing-1);
  color: var(--color-text-muted);
  font-size: var(--font-size-xs);
}

.form-actions,
.list-header {
  display: flex;
  justify-content: space-between;
  gap: var(--spacing-3);
}

.form-actions {
  justify-content: flex-end;
}

.list-card {
  padding: var(--spacing-4);
}

.list-header {
  flex-direction: column;
  align-items: stretch;
  margin-bottom: var(--spacing-3);
}

.date-range {
  gap: var(--spacing-2);
}

.date-heading {
  position: sticky;
  top: 0;
  z-index: 1;
  gap: var(--spacing-2);
  border-radius: var(--radius-md);
  background: color-mix(in srgb, var(--color-background) 95%, transparent);
  color: var(--color-text-muted);
  padding: var(--spacing-1);
  font-size: var(--font-size-xs);
  font-weight: 700;
}

.appointment-item {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--spacing-3);
}

.appointment-copy {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-1);
}

.appointment-title {
  color: var(--color-text);
  font-weight: 700;
}

.inline-meta {
  gap: var(--spacing-1);
}

.loading-block {
  display: flex;
  justify-content: center;
  padding: 2rem 0;
}

.spin {
  animation: spin 1s linear infinite;
}

@media (min-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }

  .form-grid.two {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .form-grid.four {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }

  .list-header {
    flex-direction: row;
    align-items: flex-start;
  }

  .appointment-item {
    flex-direction: row;
    justify-content: space-between;
  }

  .status-actions {
    justify-content: flex-end;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
