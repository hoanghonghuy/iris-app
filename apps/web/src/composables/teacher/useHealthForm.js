import { ref, computed } from 'vue'
import { teacherService } from '../../services/teacherService'
import { extractErrorMessage } from '../../helpers/errorHandler'

export function useHealthForm(students) {
  const isModalOpen = ref(false)
  const isSubmitting = ref(false)
  const formError = ref('')
  const successMessage = ref('')

  const formStudentId = ref('')
  const temperature = ref('')
  const symptoms = ref('')
  const severity = ref('normal')
  const note = ref('')

  const selectedStudent = computed(() => {
    return students.value.find((student) => student.student_id === formStudentId.value) || null
  })

  function resetForm(studentId = '') {
    formStudentId.value = studentId || students.value[0]?.student_id || ''
    temperature.value = ''
    symptoms.value = ''
    severity.value = 'normal'
    note.value = ''
    formError.value = ''
  }

  function openHealthModal(studentId = '') {
    resetForm(studentId)
    successMessage.value = ''
    isModalOpen.value = true
  }

  function closeModal() {
    isModalOpen.value = false
    formError.value = ''
  }

  async function handleSave() {
    if (!formStudentId.value) {
      formError.value = 'Vui lòng chọn học sinh'
      return
    }

    isSubmitting.value = true
    formError.value = ''
    successMessage.value = ''

    try {
      await teacherService.createHealthLog({
        student_id: formStudentId.value,
        temperature: temperature.value ? Number(temperature.value) : undefined,
        symptoms: symptoms.value.trim() || undefined,
        severity: severity.value,
        note: note.value.trim() || undefined,
      })

      successMessage.value = 'Đã ghi nhận sức khỏe thành công!'
      closeModal()

      return formStudentId.value // Return student ID for history refresh
    } catch (error) {
      formError.value = extractErrorMessage(error) || 'Không thể ghi nhận sức khỏe'
      return null
    } finally {
      isSubmitting.value = false
    }
  }

  return {
    isModalOpen,
    isSubmitting,
    formError,
    successMessage,
    formStudentId,
    temperature,
    symptoms,
    severity,
    note,
    selectedStudent,
    openHealthModal,
    closeModal,
    handleSave,
  }
}
