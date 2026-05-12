import { ref, watch } from 'vue'
import { teacherService } from '../../services/teacherService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { normalizeListResponse } from '../../helpers/collectionUtils'

function daysAgo(count) {
  const date = new Date()
  date.setDate(date.getDate() - count)
  return date
}

function getDateInputValue(date) {
  const local = new Date(date.getTime() - date.getTimezoneOffset() * 60000)
  return local.toISOString().slice(0, 10)
}

export function useHealthHistory(students) {
  const historyStudentId = ref('')
  const historyFrom = ref(getDateInputValue(daysAgo(6)))
  const historyTo = ref(getDateInputValue(new Date()))
  const historyLogs = ref([])
  const isLoadingHistory = ref(false)
  const historyError = ref('')

  async function fetchHistory() {
    if (!historyStudentId.value) {
      historyLogs.value = []
      historyError.value = ''
      return
    }

    isLoadingHistory.value = true
    historyError.value = ''

    try {
      historyLogs.value = normalizeListResponse(
        await teacherService.getStudentHealth(
          historyStudentId.value,
          historyFrom.value || undefined,
          historyTo.value || undefined,
        ),
      )
    } catch (error) {
      historyLogs.value = []
      historyError.value = extractErrorMessage(error) || 'Không thể tải lịch sử sức khỏe'
    } finally {
      isLoadingHistory.value = false
    }
  }

  // Auto-update history student when students list changes
  watch(students, (newStudents) => {
    if (!newStudents.some((student) => student.student_id === historyStudentId.value)) {
      historyStudentId.value = newStudents[0]?.student_id || ''
    }
  })

  // Auto-fetch when filters change
  watch([historyStudentId, historyFrom, historyTo], fetchHistory)

  return {
    historyStudentId,
    historyFrom,
    historyTo,
    historyLogs,
    isLoadingHistory,
    historyError,
    fetchHistory,
  }
}
