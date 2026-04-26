import { ref, watch } from 'vue'
import { teacherService } from '../../services/teacherService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { normalizeListResponse } from '../../helpers/collectionUtils'

/**
 * Shared composable for teacher class and student selection
 * Used across TeacherClasses, TeacherHealth, TeacherClassDetail
 */
export function useTeacherClassSelection() {
  const classes = ref([])
  const selectedClassId = ref('')
  const students = ref([])
  const isLoadingClasses = ref(true)
  const isLoadingStudents = ref(false)
  const errorMessage = ref('')

  async function fetchClasses() {
    isLoadingClasses.value = true
    errorMessage.value = ''

    try {
      classes.value = normalizeListResponse(await teacherService.getMyClasses())
      
      // Auto-select first class if current selection is invalid
      if (!classes.value.some((cls) => cls.class_id === selectedClassId.value)) {
        selectedClassId.value = classes.value[0]?.class_id || ''
      }
    } catch (error) {
      errorMessage.value = extractErrorMessage(error) || 'Không thể tải danh sách lớp'
    } finally {
      isLoadingClasses.value = false
    }
  }

  async function fetchStudents() {
    if (!selectedClassId.value) {
      students.value = []
      return
    }

    isLoadingStudents.value = true
    errorMessage.value = ''

    try {
      students.value = normalizeListResponse(
        await teacherService.getStudentsInClass(selectedClassId.value)
      )
    } catch (error) {
      students.value = []
      errorMessage.value = extractErrorMessage(error) || 'Không thể tải danh sách học sinh'
    } finally {
      isLoadingStudents.value = false
    }
  }

  // Auto-fetch students when class changes
  watch(selectedClassId, fetchStudents)

  return {
    classes,
    selectedClassId,
    students,
    isLoadingClasses,
    isLoadingStudents,
    errorMessage,
    fetchClasses,
    fetchStudents,
  }
}
