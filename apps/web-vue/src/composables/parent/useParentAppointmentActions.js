import { ref } from 'vue'
import { parentService } from '../../services/parentService'
import { extractErrorMessage } from '../../helpers/errorHandler'

export function useParentAppointmentActions() {
  const isCancelConfirmOpen = ref(false)
  const appointmentToCancel = ref(null)

  async function handleChildChange(selectedChildId, fetchAvailableSlots, clearMessages) {
    clearMessages()
    await fetchAvailableSlots(selectedChildId)
  }

  async function handleBookSlot(slot, selectedChildId, bookingNote, fetchAvailableSlots, fetchAppointments, fetchAnalytics, clearMessages, actionError, successMessage, isSubmittingBooking) {
    if (!selectedChildId || !slot?.slot_id) {
      actionError.value = 'Vui lòng chọn học sinh trước khi đặt lịch.'
      return
    }

    isSubmittingBooking.value = true
    clearMessages()

    try {
      await parentService.createAppointment({
        slot_id: slot.slot_id,
        student_id: selectedChildId,
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

  function openCancelConfirm(appointment) {
    appointmentToCancel.value = appointment
    isCancelConfirmOpen.value = true
  }

  function closeCancelConfirm() {
    appointmentToCancel.value = null
    isCancelConfirmOpen.value = false
  }

  async function handleCancelAppointment(cancellingAppointmentId, fetchAppointments, fetchAvailableSlots, fetchAnalytics, clearMessages, actionError, successMessage) {
    if (!appointmentToCancel.value?.appointment_id) return

    cancellingAppointmentId.value = appointmentToCancel.value.appointment_id
    clearMessages()

    try {
      await parentService.cancelAppointment(appointmentToCancel.value.appointment_id, 'parent_cancelled')
      successMessage.value = 'Đã hủy lịch hẹn.'
      closeCancelConfirm()
      await Promise.all([fetchAppointments(), fetchAvailableSlots(), fetchAnalytics()])
    } catch (error) {
      actionError.value = extractErrorMessage(error) || 'Không thể hủy lịch hẹn.'
      closeCancelConfirm()
    } finally {
      cancellingAppointmentId.value = ''
    }
  }

  return {
    isCancelConfirmOpen,
    appointmentToCancel,
    openCancelConfirm,
    closeCancelConfirm,
    handleChildChange,
    handleBookSlot,
    handleCancelAppointment,
  }
}
