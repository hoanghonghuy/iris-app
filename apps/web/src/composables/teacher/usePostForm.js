import { ref } from 'vue'
import { teacherService } from '../../services/teacherService'
import { extractErrorMessage } from '../../helpers/errorHandler'

export function usePostForm() {
  const showForm = ref(false)
  const scopeType = ref('class')
  const formStudentId = ref('')
  const postType = ref('announcement')
  const content = ref('')
  const submitting = ref(false)
  const formError = ref('')

  function openForm() {
    showForm.value = true
    formError.value = ''
  }

  function closeForm() {
    showForm.value = false
    content.value = ''
    formError.value = ''
  }

  function resetForm() {
    scopeType.value = 'class'
    formStudentId.value = ''
    postType.value = 'announcement'
    content.value = ''
    formError.value = ''
  }

  async function submitPost(classId, onSuccess) {
    formError.value = ''

    if (!classId) {
      formError.value = 'Vui lòng chọn lớp học'
      return false
    }

    if (!content.value.trim()) {
      formError.value = 'Vui lòng nhập nội dung bài đăng'
      return false
    }

    if (scopeType.value === 'student' && !formStudentId.value) {
      formError.value = 'Vui lòng chọn học sinh'
      return false
    }

    submitting.value = true
    try {
      await teacherService.createPost({
        scope_type: scopeType.value,
        class_id: classId,
        student_id: scopeType.value === 'student' ? formStudentId.value : undefined,
        type: postType.value,
        content: content.value.trim(),
      })
      resetForm()
      closeForm()
      if (onSuccess) await onSuccess()
      return true
    } catch (error) {
      formError.value = extractErrorMessage(error) || 'Không thể tạo bài đăng'
      return false
    } finally {
      submitting.value = false
    }
  }

  return {
    showForm,
    scopeType,
    formStudentId,
    postType,
    content,
    submitting,
    formError,
    openForm,
    closeForm,
    resetForm,
    submitPost,
  }
}
