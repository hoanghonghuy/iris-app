<script setup>
import { onMounted, watch } from 'vue'
import { AlertCircle, CheckCircle2, LoaderCircle, Plus, RefreshCw, X } from 'lucide-vue-next'
import {
  useTeacherClassSelection,
  useHealthForm,
  useHealthHistory,
} from '../../composables/teacher'
import {
  HEALTH_SEVERITY_OPTIONS,
  getSeverityLabel,
  getSeverityBadge,
} from '../../helpers/healthConfig'
import { formatDateTimeVN } from '../../helpers/dateFormatter'
import LoadingSpinner from '../../components/common/LoadingSpinner.vue'
import EmptyState from '../../components/common/EmptyState.vue'
import ActionModal from '../../components/ActionModal.vue'

const {
  classes,
  selectedClassId,
  students,
  isLoadingClasses,
  isLoadingStudents,
  errorMessage,
  fetchClasses,
  fetchStudents,
} = useTeacherClassSelection()

const {
  isModalOpen,
  isSubmitting,
  formError,
  successMessage,
  formStudentId,
  temperature,
  symptoms,
  severity,
  note,
  selectedStudent,
  openHealthModal,
  closeModal,
  handleSave,
} = useHealthForm(students)

const {
  historyStudentId,
  historyFrom,
  historyTo,
  historyLogs,
  isLoadingHistory,
  historyError,
  fetchHistory,
} = useHealthHistory(students)

async function handleSaveAndRefresh() {
  const savedStudentId = await handleSave()
  if (savedStudentId) {
    if (historyStudentId.value === savedStudentId) {
      await fetchHistory()
    } else {
      historyStudentId.value = savedStudentId
    }
  }
}

watch(selectedClassId, () => {
  successMessage.value = ''
})

onMounted(async () => {
  await fetchClasses()
  if (selectedClassId.value) {
    await fetchStudents()
  }
  if (historyStudentId.value) {
    await fetchHistory()
  }
})
</script>

<template>
  <div class="teacher-health page-stack">
    <div class="toolbar">
      <div class="toolbar-row">
        <div class="form-group mb-0 class-filter">
          <label class="form-label" for="classFilter">Chọn lớp học</label>
          <select
            id="classFilter"
            v-model="selectedClassId"
            class="form-input"
            :disabled="isLoadingClasses"
          >
            <option value="" disabled v-if="classes.length === 0">-- Không có lớp --</option>
            <option v-for="cls in classes" :key="cls.class_id" :value="cls.class_id">
              {{ cls.name }} ({{ cls.school_year }})
            </option>
          </select>
        </div>

        <button
          class="btn btn--primary"
          type="button"
          :disabled="students.length === 0"
          @click="openHealthModal()"
        >
          <Plus :size="16" />
          Ghi nhận
        </button>
      </div>
    </div>

    <div v-if="errorMessage" class="alert alert--error alert-row">
      <AlertCircle :size="16" />
      {{ errorMessage }}
    </div>

    <div v-if="successMessage" class="alert alert--success alert-row">
      <CheckCircle2 :size="16" />
      {{ successMessage }}
    </div>

    <LoadingSpinner
      v-if="isLoadingClasses || isLoadingStudents"
      message="Đang tải danh sách học sinh..."
    />

    <div v-else class="page-stack">
      <div class="card">
        <EmptyState
          v-if="classes.length === 0"
          title="Chưa có lớp học"
          message="Bạn chưa được phân công phụ trách lớp học nào."
        />

        <EmptyState
          v-else-if="students.length === 0"
          title="Không có học sinh"
          message="Lớp này hiện chưa có học sinh nào."
          icon="heart"
        />

        <div v-else class="table-responsive">
          <table class="table">
            <thead>
              <tr>
                <th width="50">STT</th>
                <th>Họ và tên</th>
                <th class="text-right">Thao tác</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(student, index) in students" :key="student.student_id">
                <td class="text-center">{{ index + 1 }}</td>
                <td class="font-medium">{{ student.full_name }}</td>
                <td class="text-right">
                  <button
                    class="btn btn--sm btn--primary"
                    type="button"
                    @click="openHealthModal(student.student_id)"
                  >
                    Ghi nhận sức khỏe
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div v-if="students.length > 0" class="card history-card">
        <div class="history-head">
          <div>
            <h3>Lịch sử sức khỏe</h3>
            <p class="history-copy">
              Theo dõi nhật ký sức khỏe đã ghi theo từng học sinh và khoảng ngày.
            </p>
          </div>
          <button
            class="btn btn--outline btn--sm"
            type="button"
            :disabled="isLoadingHistory || !historyStudentId"
            @click="fetchHistory"
          >
            <RefreshCw :size="14" :class="{ spin: isLoadingHistory }" />
            Làm mới
          </button>
        </div>

        <div class="history-filters">
          <div class="form-group mb-0">
            <label class="form-label">Học sinh</label>
            <select v-model="historyStudentId" class="form-input">
              <option
                v-for="student in students"
                :key="student.student_id"
                :value="student.student_id"
              >
                {{ student.full_name }}
              </option>
            </select>
          </div>

          <div class="form-group mb-0">
            <label class="form-label">Từ ngày</label>
            <input v-model="historyFrom" class="form-input" type="date" />
          </div>

          <div class="form-group mb-0">
            <label class="form-label">Đến ngày</label>
            <input v-model="historyTo" class="form-input" type="date" />
          </div>
        </div>

        <div v-if="historyError" class="alert alert--error mt-3">{{ historyError }}</div>

        <div v-if="isLoadingHistory" class="loading-inline">
          <LoaderCircle class="spin text-muted" :size="20" />
          Đang tải lịch sử...
        </div>

        <EmptyState
          v-else-if="historyLogs.length === 0"
          title="Chưa có nhật ký sức khỏe"
          message="Nhật ký mới sẽ hiển thị sau khi giáo viên ghi nhận cho học sinh."
          icon="heart"
        />

        <div v-else class="history-list">
          <article v-for="log in historyLogs" :key="log.health_log_id" class="history-item">
            <div class="history-item__head">
              <p class="history-date">{{ formatDateTimeVN(log.recorded_at) }}</p>
              <span :class="getSeverityBadge(log.severity)">{{
                getSeverityLabel(log.severity)
              }}</span>
            </div>

            <div class="history-item__body">
              <p>
                <span class="label">Nhiệt độ:</span>
                {{ typeof log.temperature === 'number' ? `${log.temperature}°C` : 'Không ghi' }}
              </p>
              <p><span class="label">Triệu chứng:</span> {{ log.symptoms || 'Không ghi' }}</p>
              <p><span class="label">Ghi chú:</span> {{ log.note || 'Không ghi' }}</p>
            </div>
          </article>
        </div>
      </div>
    </div>

    <ActionModal
      :is-open="isModalOpen"
      :title="`Ghi nhận sức khỏe${selectedStudent ? `: ${selectedStudent.full_name}` : ''}`"
      @close="closeModal"
    >
      <form class="form-stack" @submit.prevent="handleSave">
        <div v-if="formError" class="alert alert--error">{{ formError }}</div>

        <div class="form-group mb-0">
          <label class="form-label">Học sinh</label>
          <select v-model="formStudentId" class="form-input" :disabled="isSubmitting">
            <option
              v-for="student in students"
              :key="student.student_id"
              :value="student.student_id"
            >
              {{ student.full_name }}
            </option>
          </select>
        </div>

        <div class="form-grid">
          <div class="form-group mb-0">
            <label class="form-label">Nhiệt độ (°C)</label>
            <input
              v-model="temperature"
              type="number"
              step="0.1"
              class="form-input"
              placeholder="36.5"
              :disabled="isSubmitting"
            />
          </div>

          <div class="form-group mb-0">
            <label class="form-label">Mức độ</label>
            <div class="severity-options">
              <button
                v-for="option in HEALTH_SEVERITY_OPTIONS"
                :key="option.value"
                type="button"
                class="severity-option"
                :class="{ 'severity-option--active': severity === option.value }"
                :disabled="isSubmitting"
                @click="severity = option.value"
              >
                {{ option.label }}
              </button>
            </div>
          </div>
        </div>

        <div class="form-group mb-0">
          <label class="form-label">Triệu chứng</label>
          <input
            v-model="symptoms"
            class="form-input"
            type="text"
            placeholder="VD: ho nhẹ, sổ mũi..."
            :disabled="isSubmitting"
          />
        </div>

        <div class="form-group mb-0">
          <label class="form-label">Ghi chú</label>
          <textarea
            v-model="note"
            class="form-input"
            rows="3"
            placeholder="Ghi chú thêm..."
            :disabled="isSubmitting"
          ></textarea>
        </div>

        <div class="modal-actions">
          <button
            type="button"
            class="btn btn--outline"
            :disabled="isSubmitting"
            @click="closeModal"
          >
            <X :size="16" />
            Đóng
          </button>
          <button
            type="submit"
            class="btn btn--primary"
            :disabled="isSubmitting || !formStudentId"
            @click="handleSaveAndRefresh"
          >
            {{ isSubmitting ? 'Đang lưu...' : 'Lưu ghi nhận' }}
          </button>
        </div>
      </form>
    </ActionModal>
  </div>
</template>

<style scoped>
.page-stack,
.form-stack,
.history-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.toolbar-row,
.history-head,
.alert-row,
.loading-inline,
.modal-actions {
  display: flex;
  align-items: center;
}

.toolbar-row,
.history-head,
.modal-actions {
  justify-content: space-between;
  gap: var(--spacing-3);
}

.class-filter {
  min-width: min(100%, 18rem);
}

.history-head h3,
.history-head p,
.history-date,
.history-item__body p {
  margin: 0;
}

.history-copy {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  margin-top: var(--spacing-1);
}

.history-filters {
  display: grid;
  gap: var(--spacing-3);
  grid-template-columns: 1fr;
}

.history-card {
  padding: var(--spacing-4);
}

.history-item {
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--spacing-3);
}

.history-item__head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--spacing-2);
  margin-bottom: var(--spacing-2);
}

.history-date {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.history-item__body {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-1);
  font-size: var(--font-size-sm);
}

.label {
  font-weight: 700;
  color: var(--color-text);
}

.loading-inline {
  gap: var(--spacing-2);
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.form-grid {
  display: grid;
  gap: var(--spacing-3);
  grid-template-columns: 1fr;
}

.severity-options {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: var(--spacing-2);
}

.severity-option {
  border: 1px solid var(--color-border);
  background: var(--color-surface);
  color: var(--color-text-muted);
  border-radius: var(--radius-md);
  min-height: 2.5rem;
  font-size: var(--font-size-sm);
  font-weight: 700;
}

.severity-option--active {
  border-color: var(--color-primary);
  color: var(--color-primary);
  background: color-mix(in srgb, var(--color-primary) 10%, transparent);
}

.spin {
  animation: spin 1s linear infinite;
}

@media (min-width: 768px) {
  .toolbar-row {
    justify-content: space-between;
  }

  .history-filters,
  .form-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .form-grid > :first-child {
    grid-column: span 1;
  }
}

@media (max-width: 767px) {
  .toolbar-row,
  .history-head,
  .modal-actions {
    flex-direction: column;
    align-items: stretch;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
