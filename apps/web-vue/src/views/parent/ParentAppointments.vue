<script setup>
import ConfirmDialog from '../../components/ConfirmDialog.vue'
import ParentAppointmentBookingPanel from './appointments/ParentAppointmentBookingPanel.vue'
import ParentAppointmentHistoryPanel from './appointments/ParentAppointmentHistoryPanel.vue'
import ParentAppointmentsSummaryCard from './appointments/ParentAppointmentsSummaryCard.vue'
import { formatDateRange } from './appointments/appointmentsPresentation'
import { useParentAppointmentsPage } from './appointments/useParentAppointmentsPage'

const {
  children,
  selectedChildId,
  availableSlots,
  analytics,
  bookingNote,
  historyView,
  historyCurrentPage,
  historyFromDate,
  historyToDate,
  isBootstrapping,
  isLoadingChildren,
  isLoadingSlots,
  isLoadingAppointments,
  isSubmittingBooking,
  cancellingAppointmentId,
  errorMessage,
  actionError,
  successMessage,
  fetchedAppointmentCount,
  totalAppointmentCount,
  isCancelConfirmOpen,
  appointmentToCancel,
  timezoneDisplay,
  activeAppointmentsCount,
  filteredAppointments,
  totalHistoryPages,
  pagedAppointments,
  historySummary,
  refreshPage,
  applyHistoryFilters,
  resetHistoryFilters,
  handleChildChange,
  syncChildFromAppointment,
  handleBookSlot,
  openCancelConfirm,
  closeCancelConfirm,
  handleCancelAppointment,
  switchHistoryView,
  changeHistoryPage,
  exportHistoryCsv,
} = useParentAppointmentsPage()

function updateSelectedChildId(value) {
  selectedChildId.value = value
}

function updateBookingNote(value) {
  bookingNote.value = value
}

function updateHistoryFromDate(value) {
  historyFromDate.value = value
}

function updateHistoryToDate(value) {
  historyToDate.value = value
}
</script>

<template>
  <div class="parent-appointments page-stack">
    <header class="page-head">
      <h2 class="font-bold m-0">Lịch hẹn với giáo viên</h2>
      <p class="text-sm text-muted m-0 mt-1">
        Chọn học sinh, xem khung giờ trống và đặt lịch ngay trên một màn hình.
      </p>
      <p class="text-xs text-muted m-0 mt-1">Múi giờ hiển thị: {{ timezoneDisplay }}</p>
    </header>

    <div v-if="errorMessage" class="alert alert--error">
      <p class="font-bold m-0">Lỗi tải dữ liệu</p>
      <p class="m-0 mt-1">{{ errorMessage }}</p>
      <button class="btn btn--outline btn--sm mt-2" type="button" @click="refreshPage">
        Thử lại
      </button>
    </div>

    <div v-if="actionError" class="alert alert--error">
      {{ actionError }}
    </div>

    <div v-if="successMessage" class="alert alert--success">
      {{ successMessage }}
    </div>

    <ParentAppointmentsSummaryCard
      :analytics="analytics"
      :active-appointments-count="activeAppointmentsCount"
    />

    <div class="layout-grid">
      <ParentAppointmentBookingPanel
        :children="children"
        :selected-child-id="selectedChildId"
        :booking-note="bookingNote"
        :available-slots="availableSlots"
        :is-loading-children="isLoadingChildren || isBootstrapping"
        :is-loading-slots="isLoadingSlots || isBootstrapping"
        :is-loading-appointments="isLoadingAppointments || isBootstrapping"
        :is-submitting-booking="isSubmittingBooking"
        :timezone-display="timezoneDisplay"
        @update:selected-child-id="updateSelectedChildId"
        @update:booking-note="updateBookingNote"
        @change-child="handleChildChange"
        @refresh="refreshPage"
        @book-slot="handleBookSlot"
      />

      <ParentAppointmentHistoryPanel
        :history-view="historyView"
        :history-from-date="historyFromDate"
        :history-to-date="historyToDate"
        :is-loading-appointments="isLoadingAppointments || isBootstrapping"
        :filtered-appointments="filteredAppointments"
        :paged-appointments="pagedAppointments"
        :total-history-pages="totalHistoryPages"
        :history-current-page="historyCurrentPage"
        :history-summary="historySummary"
        :total-appointment-count="totalAppointmentCount"
        :fetched-appointment-count="fetchedAppointmentCount"
        :timezone-display="timezoneDisplay"
        :cancelling-appointment-id="cancellingAppointmentId"
        @update:history-view="switchHistoryView"
        @update:history-from-date="updateHistoryFromDate"
        @update:history-to-date="updateHistoryToDate"
        @apply-filters="applyHistoryFilters"
        @reset-filters="resetHistoryFilters"
        @export-csv="exportHistoryCsv"
        @change-page="changeHistoryPage"
        @sync-child="syncChildFromAppointment"
        @cancel-appointment="openCancelConfirm"
      />
    </div>

    <ConfirmDialog
      :is-open="isCancelConfirmOpen"
      title="Xác nhận hủy lịch"
      :message="`Bạn có chắc muốn hủy lịch hẹn vào ${appointmentToCancel ? formatDateRange(appointmentToCancel.start_time, appointmentToCancel.end_time) : ''}?`"
      confirm-text="Hủy lịch hẹn"
      is-danger
      :is-loading="Boolean(cancellingAppointmentId)"
      @confirm="handleCancelAppointment"
      @cancel="closeCancelConfirm"
    />
  </div>
</template>

<style scoped>
.page-stack {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-4);
}

.page-head {
  display: flex;
  flex-direction: column;
}

.layout-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--spacing-4);
  align-items: start;
}

@media (min-width: 1200px) {
  .layout-grid {
    grid-template-columns: 340px minmax(0, 1fr);
  }
}
</style>
