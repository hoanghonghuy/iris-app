<script setup>
import LoadingSpinner from '../../../components/LoadingSpinner.vue'
import EmptyState from '../../../components/EmptyState.vue'
import PaginationBar from '../../../components/PaginationBar.vue'
import {
  formatDateRange,
  getCancelReasonText,
  getStatusBadge,
  getStatusText,
} from './appointmentsPresentation'

defineProps({
  historyView: {
    type: String,
    required: true,
  },
  historyFromDate: {
    type: String,
    default: '',
  },
  historyToDate: {
    type: String,
    default: '',
  },
  isLoadingAppointments: {
    type: Boolean,
    default: false,
  },
  filteredAppointments: {
    type: Array,
    required: true,
  },
  pagedAppointments: {
    type: Array,
    required: true,
  },
  totalHistoryPages: {
    type: Number,
    required: true,
  },
  historyCurrentPage: {
    type: Number,
    required: true,
  },
  historySummary: {
    type: String,
    default: '',
  },
  totalAppointmentCount: {
    type: Number,
    default: 0,
  },
  fetchedAppointmentCount: {
    type: Number,
    default: 0,
  },
  timezoneDisplay: {
    type: String,
    default: '',
  },
  cancellingAppointmentId: {
    type: String,
    default: '',
  },
})

const emit = defineEmits([
  'update:history-view',
  'update:history-from-date',
  'update:history-to-date',
  'apply-filters',
  'reset-filters',
  'export-csv',
  'change-page',
  'sync-child',
  'cancel-appointment',
])

function switchView(view) {
  emit('update:history-view', view)
}
</script>

<template>
  <section class="history-panel">
    <div class="panel-head">
      <h3 class="font-bold text-lg m-0">Lịch sử lịch hẹn</h3>
      <div class="history-toggle">
        <button
          class="history-toggle__btn"
          :class="{ 'history-toggle__btn--active': historyView === 'active' }"
          type="button"
          @click="switchView('active')"
        >
          Đang hoạt động
        </button>
        <button
          class="history-toggle__btn"
          :class="{ 'history-toggle__btn--active': historyView === 'cancelled' }"
          type="button"
          @click="switchView('cancelled')"
        >
          Đã hủy
        </button>
      </div>
    </div>

    <div class="card p-4">
      <div class="history-filters">
        <div class="form-group mb-0">
          <label class="form-label">Từ ngày</label>
          <input
            :value="historyFromDate"
            type="date"
            class="form-input"
            @input="emit('update:history-from-date', $event.target.value)"
          >
        </div>
        <div class="form-group mb-0">
          <label class="form-label">Đến ngày</label>
          <input
            :value="historyToDate"
            type="date"
            class="form-input"
            @input="emit('update:history-to-date', $event.target.value)"
          >
        </div>
        <button class="btn btn--outline btn--sm" type="button" @click="emit('reset-filters')">
          7 ngày gần nhất
        </button>
        <button class="btn btn--primary btn--sm" type="button" :disabled="isLoadingAppointments" @click="emit('apply-filters')">
          Lọc lịch sử
        </button>
      </div>

      <div class="history-actions">
        <p class="text-xs text-muted m-0">
          {{ historySummary }}
          <span v-if="totalAppointmentCount > fetchedAppointmentCount">
            Đã tải {{ fetchedAppointmentCount }}/{{ totalAppointmentCount }} lịch hẹn.
          </span>
        </p>
        <button class="btn btn--outline btn--sm" type="button" @click="emit('export-csv')">
          Xuất CSV
        </button>
      </div>
    </div>

    <LoadingSpinner v-if="isLoadingAppointments" message="Đang tải lịch hẹn..." />

    <div v-else-if="filteredAppointments.length === 0" class="card">
      <EmptyState
        :title="historyView === 'cancelled' ? 'Chưa có lịch hẹn đã hủy' : 'Chưa có lịch hẹn hoạt động'"
        message="Bạn chưa có lịch hẹn phù hợp với bộ lọc hiện tại."
      />
    </div>

    <div v-else class="appointments-list">
      <article
        v-for="appointment in pagedAppointments"
        :key="appointment.appointment_id"
        class="card p-4 appointment-card"
      >
        <div class="appointment-top">
          <div>
            <p class="font-bold m-0">{{ formatDateRange(appointment.start_time, appointment.end_time) }}</p>
            <p class="text-sm text-muted m-0 mt-1">Giáo viên: {{ appointment.teacher_name || 'N/A' }}</p>
            <p class="text-sm text-muted m-0">Học sinh: {{ appointment.student_name || 'N/A' }}</p>
            <p class="text-xs text-muted m-0 mt-1">Múi giờ: {{ timezoneDisplay }}</p>
          </div>
          <span :class="getStatusBadge(appointment.status)">
            {{ getStatusText(appointment.status) }}
          </span>
        </div>

        <p v-if="appointment.note" class="appointment-note">
          {{ appointment.note }}
        </p>

        <p v-if="appointment.cancel_reason" class="appointment-cancel-reason">
          Lý do hủy: {{ getCancelReasonText(appointment.cancel_reason) }}
        </p>

        <div class="appointment-actions">
          <button class="btn btn--outline btn--sm" type="button" @click="emit('sync-child', appointment)">
            Chọn học sinh
          </button>
          <button
            v-if="appointment.status === 'pending' || appointment.status === 'confirmed'"
            class="btn btn--danger btn--sm"
            type="button"
            :disabled="cancellingAppointmentId === appointment.appointment_id"
            @click="emit('cancel-appointment', appointment)"
          >
            {{ cancellingAppointmentId === appointment.appointment_id ? 'Đang hủy...' : 'Hủy lịch' }}
          </button>
        </div>
      </article>

      <PaginationBar
        :current-page="historyCurrentPage"
        :total-pages="totalHistoryPages"
        :total-items="filteredAppointments.length"
        :limit="8"
        @page-change="emit('change-page', $event)"
      />
    </div>
  </section>
</template>

<style scoped>
.history-panel {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.panel-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--spacing-3);
}

.history-toggle {
  display: inline-flex;
  gap: var(--spacing-1);
  padding: 0.25rem;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-surface);
}

.history-toggle__btn {
  border: 0;
  background: transparent;
  color: var(--color-text-muted);
  border-radius: var(--radius-sm);
  padding: 0.5rem 0.75rem;
  font-size: var(--font-size-sm);
  font-weight: 600;
}

.history-toggle__btn--active {
  background: var(--color-background);
  color: var(--color-primary);
}

.history-filters {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: var(--spacing-3);
  align-items: end;
}

.history-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--spacing-3);
  margin-top: var(--spacing-3);
}

.appointments-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.appointment-card {
  border-left: 4px solid var(--color-primary);
}

.appointment-top {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--spacing-3);
  flex-wrap: wrap;
}

.appointment-note {
  margin: var(--spacing-3) 0 0;
  padding: var(--spacing-3);
  border-radius: var(--radius);
  background: var(--color-background);
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

.appointment-cancel-reason {
  margin: var(--spacing-3) 0 0;
  padding: var(--spacing-3);
  border-radius: var(--radius);
  border: 1px solid var(--color-danger-soft-border);
  background: var(--color-danger-soft-bg);
  color: var(--color-danger);
  font-size: var(--font-size-sm);
}

.appointment-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  flex-wrap: wrap;
  gap: var(--spacing-2);
  margin-top: var(--spacing-3);
}

@media (max-width: 767px) {
  .panel-head {
    align-items: stretch;
    flex-direction: column;
  }

  .history-toggle {
    width: 100%;
  }

  .history-toggle__btn {
    flex: 1;
  }

  .history-filters {
    grid-template-columns: 1fr;
  }

  .history-actions {
    align-items: stretch;
    flex-direction: column;
  }
}
</style>
