import { ref, computed } from 'vue'
import { teacherService } from '../../services/teacherService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { formatDateTime } from '../../helpers/appointmentConfig'

export function useAppointmentSlotCreation(classes, fetchAllAppointments) {
  const showCreateForm = ref(false)
  const submitting = ref(false)
  const errorMessage = ref('')

  const classId = ref('')
  const startTime = ref('')
  const durationMinutes = ref(30)
  const bufferMinutes = ref(10)
  const maxBookingsPerDay = ref(12)
  const note = ref('')

  const minStartTime = computed(() => {
    const local = new Date(Date.now() - new Date().getTimezoneOffset() * 60000)
    return local.toISOString().slice(0, 16)
  })

  async function createSlot(onSuccess) {
    if (!classId.value || !startTime.value) return

    submitting.value = true
    errorMessage.value = ''

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

      const activeAppointments = (
        await fetchAllAppointments({
          from: dayStart.toISOString(),
          to: dayEnd.toISOString(),
        })
      ).filter((item) => item.status !== 'cancelled')

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
        return !(
          proposedEndMs + bufferMs <= existingStart || proposedStartMs >= existingEnd + bufferMs
        )
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

      if (onSuccess) await onSuccess()
    } catch (error) {
      errorMessage.value = extractErrorMessage(error) || 'Không thể tạo khung giờ.'
    } finally {
      submitting.value = false
    }
  }

  function initializeClassId() {
    if (!classId.value && classes.value.length) {
      classId.value = classes.value[0].class_id
    }
  }

  return {
    showCreateForm,
    submitting,
    errorMessage,
    classId,
    startTime,
    durationMinutes,
    bufferMinutes,
    maxBookingsPerDay,
    note,
    minStartTime,
    createSlot,
    initializeClassId,
  }
}
