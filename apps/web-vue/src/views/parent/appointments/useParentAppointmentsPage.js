import { computed, onMounted, ref } from 'vue'
import { parentService } from '../../../services/parentService'
import { extractErrorMessage } from '../../../helpers/errorHandler'
import {
  daysAgo,
  exportAppointmentsToCsv,
  getDateInputValue,
  getTimezoneDisplay,
} from './appointmentsPresentation'

const FETCH_LIMIT = 100
const HISTORY_PAGE_SIZE = 8
const ACTIVE_STATUSES = ['pending', 'confirmed']

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

export function useParentAppointmentsPage() {
  const children = ref([])
  const selectedChildId = ref('')
  const availableSlots = ref([])
  const appointments = ref([])
  const analytics = ref(null)
  const bookingNote = ref('')

  const historyView = ref('active')
  const historyCurrentPage = ref(1)
  const historyFromDate = ref(getDateInputValue(daysAgo(6)))
  const historyToDate = ref(getDateInputValue(new Date()))

  const isBootstrapping = ref(true)
  const isLoadingChildren = ref(false)
  const isLoadingSlots = ref(false)
  const isLoadingAppointments = ref(false)
  const isSubmittingBooking = ref(false)
  const cancellingAppointmentId = ref('')

  const errorMessage = ref('')
  const actionError = ref('')
  const successMessage = ref('')

  const fetchedAppointmentCount = ref(0)
  const totalAppointmentCount = ref(0)

  const isCancelConfirmOpen = ref(false)
  const appointmentToCancel = ref(null)
  const timezoneDisplay = ref(getTimezoneDisplay())

  const activeAppointmentsCount = computed(() =>
    appointments.value.filter((item) => ACTIVE_STATUSES.includes(item.status)).length,
  )

  const filteredAppointments = computed(() => {
    const sorted = [...appointments.value].sort(
      (left, right) => new Date(right.start_time).getTime() - new Date(left.start_time).getTime(),
    )

    if (historyView.value === 'cancelled') {
      return sorted.filter((item) => item.status === 'cancelled')
    }

    return sorted.filter((item) => item.status !== 'cancelled')
  })

  const totalHistoryPages = computed(() =>
    Math.ceil(filteredAppointments.value.length / HISTORY_PAGE_SIZE) || 0,
  )

  const pagedAppointments = computed(() => {
    const start = (historyCurrentPage.value - 1) * HISTORY_PAGE_SIZE
    return filteredAppointments.value.slice(start, start + HISTORY_PAGE_SIZE)
  })

  const historySummary = computed(() => {
    if (filteredAppointments.value.length === 0) {
      return 'Không có lịch hẹn phù hợp với bộ lọc hiện tại.'
    }

    return `Đang hiển thị ${filteredAppointments.value.length} lịch hẹn trong khoảng thời gian đã chọn.`
  })

  function clearMessages() {
    actionError.value = ''
    successMessage.value = ''
  }

  function toRangeParams() {
    if (!historyFromDate.value && !historyToDate.value) return {}

    const from = historyFromDate.value
      ? new Date(`${historyFromDate.value}T00:00:00`).toISOString()
      : undefined
    const to = historyToDate.value
      ? new Date(`${historyToDate.value}T23:59:59.999`).toISOString()
      : undefined

    return { from, to }
  }

  function resolveSyncedChildId(nextChildren, nextAppointments, explicitChildId = '') {
    const hasChild = (childId) => nextChildren.some((item) => item.student_id === childId)

    if (explicitChildId && hasChild(explicitChildId)) {
      return explicitChildId
    }

    if (selectedChildId.value && hasChild(selectedChildId.value)) {
      return selectedChildId.value
    }

    const activeAppointment = nextAppointments.find((item) => ACTIVE_STATUSES.includes(item.status))
    if (activeAppointment?.student_id && hasChild(activeAppointment.student_id)) {
      return activeAppointment.student_id
    }

    const fallbackFromHistory = nextAppointments.find((item) => hasChild(item.student_id))
    if (fallbackFromHistory?.student_id) {
      return fallbackFromHistory.student_id
    }

    return nextChildren[0]?.student_id || ''
  }

  async function fetchChildren() {
    isLoadingChildren.value = true
    try {
      children.value = normalizeListResponse(await parentService.getMyChildren())
    } finally {
      isLoadingChildren.value = false
    }
  }

  async function fetchAnalytics() {
    try {
      const response = await parentService.getAnalytics()
      analytics.value = response?.data ?? response ?? null
    } catch {
      analytics.value = null
    }
  }

  async function fetchAppointments() {
    isLoadingAppointments.value = true
    errorMessage.value = ''

    try {
      const { from, to } = toRangeParams()
      if (from && to && new Date(from).getTime() > new Date(to).getTime()) {
        appointments.value = []
        fetchedAppointmentCount.value = 0
        totalAppointmentCount.value = 0
        historyCurrentPage.value = 1
        errorMessage.value = 'Khoảng thời gian không hợp lệ: Từ ngày phải nhỏ hơn hoặc bằng Đến ngày.'
        return
      }

      let offset = 0
      let hasMore = true
      const combinedAppointments = []
      let total = 0

      while (hasMore) {
        const response = await parentService.getAppointments({
          limit: FETCH_LIMIT,
          offset,
          from,
          to,
        })

        const { items, pagination } = normalizePaginatedResponse(response)
        combinedAppointments.push(...items)
        total = pagination.total ?? combinedAppointments.length

        hasMore = Boolean(pagination.has_more) && items.length > 0
        offset += pagination.limit || FETCH_LIMIT

        if (!hasMore || offset >= total) {
          break
        }
      }

      appointments.value = combinedAppointments
      fetchedAppointmentCount.value = combinedAppointments.length
      totalAppointmentCount.value = total
      historyCurrentPage.value = 1
    } catch (error) {
      errorMessage.value = extractErrorMessage(error) || 'Không thể tải dữ liệu lịch hẹn.'
    } finally {
      isLoadingAppointments.value = false
    }
  }

  async function fetchAvailableSlots(childId = selectedChildId.value) {
    if (!childId) {
      availableSlots.value = []
      return
    }

    isLoadingSlots.value = true
    try {
      availableSlots.value = normalizeListResponse(
        await parentService.getAvailableSlots({
          student_id: childId,
          limit: 50,
          offset: 0,
        }),
      )
    } catch (error) {
      actionError.value = extractErrorMessage(error) || 'Không thể tải khung giờ khả dụng.'
    } finally {
      isLoadingSlots.value = false
    }
  }

  async function initializePage(preferredChildId = '') {
    isBootstrapping.value = true
    clearMessages()

    try {
      await Promise.all([fetchChildren(), fetchAppointments(), fetchAnalytics()])

      const syncedChildId = resolveSyncedChildId(
        children.value,
        appointments.value,
        preferredChildId,
      )
      selectedChildId.value = syncedChildId
      await fetchAvailableSlots(syncedChildId)
    } catch (error) {
      errorMessage.value = extractErrorMessage(error) || 'Không thể tải dữ liệu lịch hẹn.'
    } finally {
      isBootstrapping.value = false
    }
  }

  async function refreshPage() {
    await initializePage(selectedChildId.value)
  }

  async function applyHistoryFilters() {
    clearMessages()
    await fetchAppointments()
  }

  async function resetHistoryFilters() {
    historyFromDate.value = getDateInputValue(daysAgo(6))
    historyToDate.value = getDateInputValue(new Date())
    await applyHistoryFilters()
  }

  async function handleChildChange(nextChildId) {
    selectedChildId.value = nextChildId
    clearMessages()
    await fetchAvailableSlots(nextChildId)
  }

  async function syncChildFromAppointment(appointment) {
    if (!appointment?.student_id) return
    await handleChildChange(appointment.student_id)
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

  function openCancelConfirm(appointment) {
    appointmentToCancel.value = appointment
    isCancelConfirmOpen.value = true
  }

  function closeCancelConfirm() {
    appointmentToCancel.value = null
    isCancelConfirmOpen.value = false
  }

  async function handleCancelAppointment() {
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

  function switchHistoryView(view) {
    historyView.value = view
    historyCurrentPage.value = 1
  }

  function changeHistoryPage(page) {
    historyCurrentPage.value = page
  }

  function exportHistoryCsv() {
    clearMessages()

    const exported = exportAppointmentsToCsv({
      appointments: filteredAppointments.value,
      historyView: historyView.value,
      timezoneDisplay: timezoneDisplay.value,
    })

    if (!exported) {
      actionError.value = 'Không có dữ liệu để xuất CSV.'
      return
    }

    successMessage.value = 'Đã xuất CSV theo bộ lọc hiện tại.'
  }

  onMounted(async () => {
    await initializePage()
  })

  return {
    children,
    selectedChildId,
    availableSlots,
    appointments,
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
  }
}
