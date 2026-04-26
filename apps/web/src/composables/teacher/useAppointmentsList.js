import { ref, computed } from 'vue'
import { teacherService } from '../../services/teacherService'
import { fetchAllPaginated, normalizeListResponse } from '../../helpers/collectionUtils'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { getLocalDateKey } from '../../helpers/appointmentConfig'
import { offsetDate, toDateInputValue } from '../../helpers/dateHelpers'

const FETCH_LIMIT = 100

export function useAppointmentsList() {
  const classes = ref([])
  const appointments = ref([])
  const loading = ref(true)
  const errorMessage = ref('')
  const updatingAppointmentId = ref(null)

  const statusFilter = ref('')
  const filterFromDate = ref(toDateInputValue(offsetDate(-6)))
  const filterToDate = ref(toDateInputValue(new Date()))

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

  async function fetchAllAppointments(params = {}) {
    const { items } = await fetchAllPaginated(
      ({ limit, offset }) =>
        teacherService.getAppointments({
          ...params,
          limit,
          offset,
        }),
      { limit: FETCH_LIMIT },
    )

    return items
  }

  async function loadData() {
    loading.value = true
    errorMessage.value = ''

    try {
      const from = filterFromDate.value
        ? new Date(`${filterFromDate.value}T00:00:00`).toISOString()
        : undefined
      const to = filterToDate.value
        ? new Date(`${filterToDate.value}T23:59:59.999`).toISOString()
        : undefined

      if (from && to && new Date(from).getTime() > new Date(to).getTime()) {
        errorMessage.value =
          'Khoảng ngày lọc không hợp lệ: Từ ngày phải nhỏ hơn hoặc bằng Đến ngày.'
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
    } catch (error) {
      errorMessage.value = extractErrorMessage(error) || 'Không thể tải dữ liệu lịch hẹn.'
    } finally {
      loading.value = false
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

  return {
    classes,
    appointments,
    loading,
    errorMessage,
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
  }
}
