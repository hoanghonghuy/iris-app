<script setup>
import { onMounted } from 'vue'
import {
  AlertTriangle,
  CalendarDays,
  Clock3,
  Download,
  LoaderCircle,
  RefreshCcw,
  UserRound,
} from 'lucide-vue-next'
import { downloadCsv } from '../../helpers/csvExport'
import { toDateInputValue } from '../../helpers/dateHelpers'
import {
  APPOINTMENT_STATUS_CONFIG,
  APPOINTMENT_STATUS_OPTIONS,
  getCancelReasonText,
  getUtcOffsetLabel,
  formatDayHeading,
  formatDateTime,
  getLocalDateKey,
} from '../../helpers/appointmentConfig'
import { useAppointmentsList, useAppointmentSlotCreation } from '../../composables/teacher'

const timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone || 'Local'
const utcOffsetLabel = getUtcOffsetLabel()
const timezoneDisplay = `${timeZone} (${utcOffsetLabel})`

const {
  classes,
  appointments,
  loading,
  errorMessage: listErrorMessage,
  updatingAppointmentId,
  statusFilter,
  filterFromDate,
  filterToDate,
  stats,
  groupedAppointments,
  loadData,
  updateStatus,
  resetLastSevenDays,
  fetchAllAppointments,
} = useAppointmentsList()

const {
  showCreateForm,
  submitting,
  errorMessage: createErrorMessage,
  classId,
  startTime,
  durationMinutes,
  bufferMinutes,
  maxBookingsPerDay,
  note,
  minStartTime,
  createSlot,
  initializeClassId,
} = useAppointmentSlotCreation(classes, fetchAllAppointments)

const errorMessage = listErrorMessage

async function handleCreateSlot() {
  await createSlot(async () => {
    await loadData()
    initializeClassId()
  })
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
      APPOINTMENT_STATUS_CONFIG[appointment.status]?.label || appointment.status,
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

onMounted(async () => {
  await loadData()
  initializeClassId()
})

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
      <button type="button" class="btn btn--outline btn--sm" @click="loadData">Thử lại</button>
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
        <button type="button" class="btn btn--primary" @click="showCreateForm = !showCreateForm">
          {{ showCreateForm ? 'Đóng biểu mẫu' : 'Tạo khung giờ mới' }}
        </button>
        <button type="button" class="btn btn--outline" :disabled="loading" @click="loadData">
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

        <div v-if="createErrorMessage" class="alert alert--error alert-row">
          <AlertTriangle :size="16" />
          {{ createErrorMessage }}
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn--outline" @click="showCreateForm = false">Đóng</button>
          <button type="button" class="btn btn--primary" :disabled="submitting || !classId || !startTime" @click="handleCreateSlot">
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
            <option v-for="option in APPOINTMENT_STATUS_OPTIONS" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>

          <div class="date-range">
            <input v-model="filterFromDate" class="form-input" type="date" @change="loadData" />
            <span>đến</span>
            <input v-model="filterToDate" class="form-input" type="date" @change="loadData" />
          </div>

          <button type="button" class="btn btn--outline btn--sm" @click="resetLastSevenDays">7 ngày gần nhất</button>
          <button type="button" class="btn btn--outline btn--sm" @click="exportAppointmentsCsv">
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
              <button type="button" class="btn btn--sm btn--primary" :disabled="appointment.status !== 'pending' || updatingAppointmentId === appointment.appointment_id" @click="updateStatus(appointment.appointment_id, 'confirmed')">
                Xác nhận
              </button>
              <button type="button" class="btn btn--sm btn--outline" :disabled="appointment.status !== 'confirmed' || updatingAppointmentId === appointment.appointment_id" @click="updateStatus(appointment.appointment_id, 'completed')">
                Hoàn tất
              </button>
              <button type="button" class="btn btn--sm btn--outline" :disabled="appointment.status !== 'confirmed' || updatingAppointmentId === appointment.appointment_id" @click="updateStatus(appointment.appointment_id, 'no_show')">
                Vắng mặt
              </button>
              <button type="button" class="btn btn--sm btn--danger" :disabled="['cancelled', 'completed'].includes(appointment.status) || updatingAppointmentId === appointment.appointment_id" @click="updateStatus(appointment.appointment_id, 'cancelled')">
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
