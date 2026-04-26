<script setup>
import LoadingSpinner from '../../../components/common/LoadingSpinner.vue'

defineProps({
  formatDateRange: {
    type: Function,
    required: true,
  },
  children: {
    type: Array,
    required: true,
  },
  selectedChildId: {
    type: String,
    default: '',
  },
  bookingNote: {
    type: String,
    default: '',
  },
  availableSlots: {
    type: Array,
    required: true,
  },
  isLoadingChildren: {
    type: Boolean,
    default: false,
  },
  isLoadingSlots: {
    type: Boolean,
    default: false,
  },
  isLoadingAppointments: {
    type: Boolean,
    default: false,
  },
  isSubmittingBooking: {
    type: Boolean,
    default: false,
  },
  timezoneDisplay: {
    type: String,
    default: '',
  },
})

const emit = defineEmits([
  'update:selected-child-id',
  'update:booking-note',
  'change-child',
  'refresh',
  'book-slot',
])

function handleChildChange(event) {
  const nextChildId = event.target.value
  emit('update:selected-child-id', nextChildId)
  emit('change-child', nextChildId)
}

function handleNoteChange(event) {
  emit('update:booking-note', event.target.value)
}
</script>

<template>
  <section class="card panel">
    <div class="panel-head">
      <h3 class="font-bold m-0">Đặt lịch hẹn mới</h3>
      <button
        class="btn btn--outline btn--sm"
        type="button"
        :disabled="isLoadingSlots || isLoadingAppointments"
        @click="emit('refresh')"
      >
        Làm mới
      </button>
    </div>

    <div class="form-group">
      <label class="form-label" for="parentAppointmentChild">Chọn con</label>
      <select
        id="parentAppointmentChild"
        class="form-input"
        :value="selectedChildId"
        :disabled="isLoadingChildren"
        @change="handleChildChange"
      >
        <option value="" disabled>
          {{ children.length === 0 ? '-- Không có dữ liệu --' : '-- Chọn học sinh --' }}
        </option>
        <option v-for="child in children" :key="child.student_id" :value="child.student_id">
          {{ child.full_name }}
        </option>
      </select>
    </div>

    <div class="form-group">
      <label class="form-label" for="parentAppointmentNote">Ghi chú khi đặt lịch</label>
      <textarea
        id="parentAppointmentNote"
        class="form-input"
        rows="4"
        maxlength="500"
        :value="bookingNote"
        placeholder="Ví dụ: Mong muốn trao đổi về tình hình học tập tuần này..."
        @input="handleNoteChange"
      />
      <p class="text-xs text-muted m-0 mt-1">Tối đa 500 ký tự.</p>
    </div>

    <div class="slots-head">
      <p class="font-medium text-sm uppercase m-0">Khung giờ còn trống</p>
      <p class="text-xs text-muted m-0">Múi giờ: {{ timezoneDisplay }}</p>
    </div>

    <LoadingSpinner v-if="isLoadingSlots" message="Đang tải lịch trống..." />

    <div v-else-if="!selectedChildId" class="empty-slot">
      Vui lòng chọn học sinh để xem khung giờ có thể đặt.
    </div>

    <div v-else-if="availableSlots.length === 0" class="empty-slot">
      Hiện chưa có khung giờ phù hợp cho học sinh này.
    </div>

    <div v-else class="slots-list">
      <article v-for="slot in availableSlots" :key="slot.slot_id" class="slot-item">
        <div class="slot-copy">
          <p class="font-bold text-sm m-0">{{ formatDateRange(slot.start_time, slot.end_time) }}</p>
          <p class="text-xs text-muted m-0 mt-1">Giáo viên: {{ slot.teacher_name || 'N/A' }}</p>
          <p v-if="slot.class_name" class="text-xs text-muted m-0">Lớp: {{ slot.class_name }}</p>
          <p v-if="slot.note" class="text-xs text-muted m-0">Ghi chú slot: {{ slot.note }}</p>
        </div>
        <button
          class="btn btn--primary btn--sm shrink-0"
          type="button"
          :disabled="isSubmittingBooking"
          @click="emit('book-slot', slot)"
        >
          {{ isSubmittingBooking ? 'Đang đặt...' : 'Đặt lịch' }}
        </button>
      </article>
    </div>
  </section>
</template>

<style scoped>
.panel {
  padding: var(--spacing-4);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-3);
}

.slots-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-3);
  flex-wrap: wrap;
}

.slots-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-2);
  max-height: 26rem;
  overflow-y: auto;
  padding-right: 0.2rem;
}

.slot-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--spacing-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  padding: var(--spacing-3);
}

.slot-copy {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.empty-slot {
  border: 1px dashed var(--color-border);
  border-radius: var(--radius);
  padding: var(--spacing-4);
  text-align: center;
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
}

@media (max-width: 767px) {
  .panel-head {
    align-items: stretch;
    flex-direction: column;
  }

  .slot-item {
    align-items: stretch;
    flex-direction: column;
  }
}
</style>
