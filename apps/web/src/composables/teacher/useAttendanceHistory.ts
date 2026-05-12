import { ref } from 'vue'
import { teacherService } from '../../services/teacherService'
import { normalizeListResponse } from '../../helpers/collectionUtils'
import { daysAgo, getDateInputValue, getTodayDateString } from '../../helpers/dateHelpers'

export function useAttendanceHistory() {
  const historyOpen = ref(new Set())
  const historyLoading = ref(new Set())
  const historyByStudent = ref({})

  const historyFrom = ref(getDateInputValue(daysAgo(7)))
  const historyTo = ref(getTodayDateString())
  const historyStudentId = ref('all')
  const historyStatus = ref('all')
  const historyListLoading = ref(false)
  const historyList = ref([])
  const historyOffset = ref(0)
  const historyLimit = ref(20)
  const historyTotal = ref(0)
  const historyHasMore = ref(false)

  function resetStudentHistoryState() {
    historyOpen.value = new Set()
    historyLoading.value = new Set()
    historyByStudent.value = {}
  }

  async function toggleHistory(studentId) {
    const today = getTodayDateString()
    const nextOpen = new Set(historyOpen.value)
    const shouldOpen = !nextOpen.has(studentId)

    if (shouldOpen) {
      nextOpen.add(studentId)
    } else {
      nextOpen.delete(studentId)
    }
    historyOpen.value = nextOpen

    if (!shouldOpen || historyByStudent.value[studentId]) {
      return
    }

    const nextLoading = new Set(historyLoading.value)
    nextLoading.add(studentId)
    historyLoading.value = nextLoading

    try {
      const formattedFromDate = getDateInputValue(daysAgo(30))
      const records = normalizeListResponse(
        await teacherService.getStudentAttendanceChanges(studentId, formattedFromDate, today),
      )
      historyByStudent.value = {
        ...historyByStudent.value,
        [studentId]: records,
      }
    } catch {
      historyByStudent.value = {
        ...historyByStudent.value,
        [studentId]: [],
      }
    } finally {
      const finishedLoading = new Set(historyLoading.value)
      finishedLoading.delete(studentId)
      historyLoading.value = finishedLoading
    }
  }

  async function loadClassHistory(classId, students, offset = 0) {
    if (!classId || students.length === 0) {
      historyList.value = []
      historyTotal.value = 0
      historyHasMore.value = false
      return
    }

    historyListLoading.value = true

    try {
      const response = await teacherService.getClassAttendanceChanges(classId, {
        from: historyFrom.value || undefined,
        to: historyTo.value || undefined,
        student_id: historyStudentId.value === 'all' ? undefined : historyStudentId.value,
        status: historyStatus.value === 'all' ? undefined : historyStatus.value,
        limit: historyLimit.value,
        offset,
      })

      const items = normalizeListResponse(response).map((item) => ({
        ...item,
        student_name:
          item.student_name ||
          students.find((student) => student.student_id === item.student_id)?.full_name ||
          'Không rõ',
      }))

      historyList.value = items
      historyTotal.value = response?.pagination?.total || items.length
      historyHasMore.value = Boolean(response?.pagination?.has_more)
      historyOffset.value = offset
    } catch (error) {
      historyList.value = []
      historyTotal.value = 0
      historyHasMore.value = false
      throw error
    } finally {
      historyListLoading.value = false
    }
  }

  function handleHistoryPrev(classId, students) {
    if (historyOffset.value <= 0) return
    loadClassHistory(classId, students, Math.max(0, historyOffset.value - historyLimit.value))
  }

  function handleHistoryNext(classId, students) {
    if (!historyHasMore.value) return
    loadClassHistory(classId, students, historyOffset.value + historyLimit.value)
  }

  return {
    historyOpen,
    historyLoading,
    historyByStudent,
    historyFrom,
    historyTo,
    historyStudentId,
    historyStatus,
    historyListLoading,
    historyList,
    historyOffset,
    historyLimit,
    historyTotal,
    historyHasMore,
    resetStudentHistoryState,
    toggleHistory,
    loadClassHistory,
    handleHistoryPrev,
    handleHistoryNext,
  }
}
