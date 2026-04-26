<script setup>
import { onMounted } from 'vue'
import ConfirmDialog from '../../components/common/ConfirmDialog.vue'
import ParentAppointmentBookingPanel from './appointments/ParentAppointmentBookingPanel.vue'
import ParentAppointmentHistoryPanel from './appointments/ParentAppointmentHistoryPanel.vue'
import ParentAppointmentsSummaryCard from './appointments/ParentAppointmentsSummaryCard.vue'
import { parentService } from '../../services/parentService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { formatDateTime } from '../../helpers/dateFormatter'
import { downloadCsv } from '../../helpers/csvExport'
import { getDateInputValue } from '../../helpers/dateHelpers'
import {
  getStatusText,
  getCancelReasonText,
  getTimezoneDisplay,
  formatDateRange as formatDateRangeHelper,
} from '../../helpers/appointmentConfig'
import { useParentAppointments, useParentAppointmentActions } from '../../composables/parent'

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
  activeAppointmentsCount,
  filteredAppointments,
  totalHistoryPages,
  pagedAppointments,
  historySummary,
  clearMessages,
  fetchAvailableSlots,
  fetchAppointments,
  fetchAnalytics,
  initializePage,
  refreshPage,
  applyHistoryFilters,
  resetHistoryFilters,
  switchHistoryView,
  changeHistoryPage,
} = useParentAppointments()

const {
  isCancelConfirmOpen,
  appointmentToCancel,
  openCancelConfirm,
  closeCancelConfirm,
  handleCancelAppointment: handleCancelAppointmentAction,
} = useParentAppointmentActions()

const timezoneDisplay = getTimezoneDisplay()

function formatDateRange(startTime, endTime) {
  return formatDateRangeHelper(startTime, endTime, formatDateTime)
}

async function handleChildChange(nextChildId) {
  selectedChildId.value = nextChildId
  clearMessages()
  await fetchAvailableSlots(nextChildId)
}

async function handleBookSlot(slot) {
  if (!selectedChildId.value || !slot?.slot_id) {
    actionError.value = 'Vui lòng chọn học sinh trước khi đặt lịch.'
    return
  }

  isSubmittingBooking.value = true
  clearMessages()

  try {
    await parentService.createAppointment({
      slot_id: slot.slot_id,
      student_id: selectedChildId.value,
      note: bookingNote.value.trim() || undefined,
    })

    bookingNote.value = ''
    successMessage.value = 'Đặt lịch hẹn thành công. Vui lòng chờ giáo viên xác nhận.'
    await Promise.all([fetchAvailableSlots(), fetchAppointments(), fetchAnalytics()])
  } catch (error) {
    actionError.value = extractErrorMessage(error) || 'Không thể đặt lịch hẹn.'
  } finally {
    isSubmittingBooking.value = false
  }
}

async function handleCancelAppointment() {
  await handleCancelAppointmentAction(
    cancellingAppointmentId,
    fetchAppointments,
    fetchAvailableSlots,
    fetchAnalytics,
    clearMessages,
    actionError,
    successMessage,
  )
}

function exportHistoryCsv() {
  clearMessages()

  if (!Array.isArray(filteredAppointments.value) || filteredAppointments.value.length === 0) {
    actionError.value = 'Không có dữ liệu để xuất CSV.'
    return
  }

  const rows = [...filteredAppointments.value]
    .sort(
      (left, right) => new Date(left.start_time).getTime() - new Date(right.start_time).getTime(),
    )
    .map((item) => [
      item.student_name || item.student_id || 'N/A',
      item.teacher_name || item.teacher_id || 'N/A',
      getStatusText(item.status),
      formatDateTime(item.start_time),
      formatDateTime(item.end_time),
      timezoneDisplay,
      item.note || '',
      getCancelReasonText(item.cancel_reason),
    ])

  downloadCsv(
    `parent-appointments-${historyView.value || 'history'}-${getDateInputValue(new Date())}.csv`,
    ['HocSinh', 'GiaoVien', 'TrangThai', 'BatDau', 'KetThuc', 'MuiGio', 'GhiChu', 'LyDoHuy'],
    rows,
  )

  successMessage.value = 'Đã xuất CSV theo bộ lọc hiện tại.'
}

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

onMounted(async () => {
  await initializePage()
})
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
        :format-date-range="formatDateRange"
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
        :format-date-range="formatDateRange"
        :get-cancel-reason-text="getCancelReasonText"
        :get-status-badge="getStatusBadge"
        :get-status-text="getStatusText"
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
