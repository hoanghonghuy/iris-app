import { ref } from 'vue'
import { teacherService } from '../../services/teacherService'
import { normalizeListResponse } from '../../helpers/collectionUtils'
import { extractErrorMessage } from '../../helpers/errorHandler'

export function useAttendanceClasses(routeQuery) {
  const classes = ref([])
  const selectedClassId = ref('')
  const isLoadingClasses = ref(true)
  const errorMessage = ref('')

  async function fetchClasses() {
    isLoadingClasses.value = true
    errorMessage.value = ''

    try {
      classes.value = normalizeListResponse(await teacherService.getMyClasses())

      const requestedClassId = typeof routeQuery.classId === 'string' ? routeQuery.classId : ''
      const nextClassId = classes.value.some((item) => item.class_id === requestedClassId)
        ? requestedClassId
        : classes.value[0]?.class_id || ''

      if (!classes.value.some((item) => item.class_id === selectedClassId.value)) {
        selectedClassId.value = nextClassId
      }
    } catch (error) {
      errorMessage.value = extractErrorMessage(error) || 'Không thể tải lớp'
    } finally {
      isLoadingClasses.value = false
    }
  }

  return {
    classes,
    selectedClassId,
    isLoadingClasses,
    errorMessage,
    fetchClasses,
  }
}
