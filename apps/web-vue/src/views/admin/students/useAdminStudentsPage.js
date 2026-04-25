import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useAuthStore } from '../../../stores/authStore'
import { adminService } from '../../../services/adminService'
import { extractErrorMessage } from '../../../helpers/errorHandler'
import { normalizeGender } from './studentPresentation'

const FETCH_LIMIT = 100

export function useAdminStudentsPage() {
  const authStore = useAuthStore()

  const schools = ref([])
  const classes = ref([])
  const selectedSchoolId = ref('')
  const selectedClassId = ref('')

  const students = ref([])
  const searchQuery = ref('')
  const isBootstrapping = ref(true)
  const isLoadingStudents = ref(false)
  const errorMessage = ref('')
  const codeError = ref('')

  const isFormModalOpen = ref(false)
  const isFormSubmitting = ref(false)
  const formError = ref('')
  const formMode = ref('add')
  const formData = ref({
    id: '',
    full_name: '',
    dob: '',
    gender: 'male',
  })

  const generatingCodeStudentId = ref('')
  const revokingCodeStudentId = ref('')
  const copiedStudentId = ref('')

  const isRevokeConfirmOpen = ref(false)
  const revokeTarget = ref(null)

  const isDeleteConfirmOpen = ref(false)
  const deleteTarget = ref(null)
  const isDeleteLoading = ref(false)

  const hasInitialized = ref(false)
  let copyTimeoutId = null

  const isSuperAdmin = computed(() => authStore.currentUserRole === 'SUPER_ADMIN')
  const selectedSchoolName = computed(() => {
    return schools.value.find((school) => school.school_id === selectedSchoolId.value)?.name || ''
  })
  const selectedClassName = computed(() => {
    return classes.value.find((classItem) => classItem.class_id === selectedClassId.value)?.name || ''
  })
  const filteredStudents = computed(() => {
    const normalizedQuery = searchQuery.value.trim().toLowerCase()
    if (!normalizedQuery) {
      return students.value
    }

    return students.value.filter((student) => student.full_name?.toLowerCase().includes(normalizedQuery))
  })

  function resetForm() {
    formData.value = {
      id: '',
      full_name: '',
      dob: '',
      gender: 'male',
    }
    formError.value = ''
  }

  async function ensureCurrentUser() {
    if (!authStore.currentUser && authStore.isAuthenticated) {
      await authStore.fetchCurrentUser()
    }
  }

  function getPreferredSchoolId(items) {
    if (!Array.isArray(items) || items.length === 0) {
      return ''
    }

    const currentSchoolId = authStore.currentUser?.school_id
    if (currentSchoolId && items.some((school) => school.school_id === currentSchoolId)) {
      return currentSchoolId
    }

    return items[0].school_id
  }

  async function loadSchools() {
    const response = await adminService.getSchools({ limit: 100, offset: 0 })
    const items = Array.isArray(response?.data) ? response.data : []
    schools.value = items

    if (!items.some((school) => school.school_id === selectedSchoolId.value)) {
      selectedSchoolId.value = getPreferredSchoolId(items)
    }
  }

  async function loadClasses() {
    if (!selectedSchoolId.value) {
      classes.value = []
      selectedClassId.value = ''
      students.value = []
      return
    }

    const response = await adminService.getClassesBySchool(selectedSchoolId.value, {
      limit: 100,
      offset: 0,
    })
    const items = Array.isArray(response?.data) ? response.data : []
    classes.value = items

    if (!items.some((classItem) => classItem.class_id === selectedClassId.value)) {
      selectedClassId.value = items[0]?.class_id || ''
    }

    if (items.length === 0) {
      students.value = []
    }
  }

  async function loadStudents() {
    if (!selectedClassId.value) {
      students.value = []
      return
    }

    isLoadingStudents.value = true
    errorMessage.value = ''
    codeError.value = ''

    try {
      let offset = 0
      let hasMore = true
      const allStudents = []

      while (hasMore) {
        const response = await adminService.getStudentsByClass(selectedClassId.value, {
          limit: FETCH_LIMIT,
          offset,
        })

        const items = Array.isArray(response?.data) ? response.data : []
        allStudents.push(...items)

        const pagination = response?.pagination
        hasMore = Boolean(pagination?.has_more) && items.length > 0
        offset += pagination?.limit || FETCH_LIMIT
      }

      students.value = allStudents
    } catch (error) {
      errorMessage.value = extractErrorMessage(error) || 'Không thể tải danh sách học sinh'
    } finally {
      isLoadingStudents.value = false
    }
  }

  function openAddModal() {
    formMode.value = 'add'
    resetForm()
    isFormModalOpen.value = true
  }

  function openEditModal(student) {
    formMode.value = 'edit'
    formData.value = {
      id: student.student_id,
      full_name: student.full_name || '',
      dob: student.dob?.includes('T') ? student.dob.split('T')[0] : student.dob || '',
      gender: normalizeGender(student.gender),
    }
    formError.value = ''
    isFormModalOpen.value = true
  }

  function closeFormModal() {
    isFormModalOpen.value = false
    formError.value = ''
  }

  function updateFormData(nextValue) {
    formData.value = nextValue
  }

  async function submitForm() {
    if (!formData.value.full_name.trim()) {
      formError.value = 'Họ tên không được để trống'
      return
    }

    if (!formData.value.dob) {
      formError.value = 'Ngày sinh không được để trống'
      return
    }

    isFormSubmitting.value = true
    formError.value = ''

    try {
      if (formMode.value === 'add') {
        await adminService.createStudent({
          school_id: selectedSchoolId.value,
          class_id: selectedClassId.value,
          full_name: formData.value.full_name.trim(),
          dob: formData.value.dob,
          gender: formData.value.gender,
        })
      } else {
        await adminService.updateStudent(formData.value.id, {
          full_name: formData.value.full_name.trim(),
          dob: formData.value.dob,
          gender: formData.value.gender,
        })
      }

      closeFormModal()
      await loadStudents()
    } catch (error) {
      formError.value = extractErrorMessage(error) || 'Không thể lưu học sinh'
    } finally {
      isFormSubmitting.value = false
    }
  }

  function confirmDelete(student) {
    deleteTarget.value = student
    isDeleteConfirmOpen.value = true
  }

  function closeDeleteConfirm() {
    isDeleteConfirmOpen.value = false
    deleteTarget.value = null
  }

  async function handleDelete() {
    if (!deleteTarget.value) {
      return
    }

    isDeleteLoading.value = true

    try {
      await adminService.deleteStudent(deleteTarget.value.student_id)
      closeDeleteConfirm()
      await loadStudents()
    } catch (error) {
      errorMessage.value = `Lỗi xóa: ${extractErrorMessage(error)}`
      isDeleteConfirmOpen.value = false
    } finally {
      isDeleteLoading.value = false
    }
  }

  async function handleGenerateCode(student) {
    generatingCodeStudentId.value = student.student_id
    codeError.value = ''

    try {
      const response = await adminService.generateParentCode(student.student_id)
      const payload = response?.data || {}

      students.value = students.value.map((item) => {
        if (item.student_id !== student.student_id) {
          return item
        }

        return {
          ...item,
          active_parent_code: payload.parent_code,
          code_expires_at: payload.expires_at,
        }
      })
    } catch (error) {
      codeError.value = extractErrorMessage(error) || 'Không thể tạo mã phụ huynh'
    } finally {
      generatingCodeStudentId.value = ''
    }
  }

  function confirmRevokeCode(student) {
    revokeTarget.value = student
    isRevokeConfirmOpen.value = true
  }

  function closeRevokeConfirm() {
    isRevokeConfirmOpen.value = false
    revokeTarget.value = null
  }

  async function handleRevokeCode() {
    if (!revokeTarget.value) {
      return
    }

    revokingCodeStudentId.value = revokeTarget.value.student_id
    codeError.value = ''

    try {
      await adminService.revokeParentCode(revokeTarget.value.student_id)
      students.value = students.value.map((item) => {
        if (item.student_id !== revokeTarget.value.student_id) {
          return item
        }

        return {
          ...item,
          active_parent_code: undefined,
          code_expires_at: undefined,
        }
      })
      closeRevokeConfirm()
    } catch (error) {
      codeError.value = extractErrorMessage(error) || 'Không thể thu hồi mã phụ huynh'
    } finally {
      revokingCodeStudentId.value = ''
    }
  }

  async function handleCopyCode(code, studentId) {
    if (!code) {
      return
    }

    try {
      await navigator.clipboard.writeText(code)
      copiedStudentId.value = studentId
      if (copyTimeoutId) {
        clearTimeout(copyTimeoutId)
      }
      copyTimeoutId = window.setTimeout(() => {
        copiedStudentId.value = ''
      }, 2000)
    } catch {
      codeError.value = 'Không thể sao chép mã phụ huynh'
    }
  }

  watch(selectedSchoolId, async (newValue, oldValue) => {
    if (!hasInitialized.value || newValue === oldValue) {
      return
    }

    await loadClasses()
  })

  watch(selectedClassId, async (newValue, oldValue) => {
    if (!hasInitialized.value || newValue === oldValue) {
      return
    }

    await loadStudents()
  })

  onMounted(async () => {
    try {
      await ensureCurrentUser()
      await loadSchools()
      await loadClasses()
      await loadStudents()
      hasInitialized.value = true
    } catch (error) {
      errorMessage.value = extractErrorMessage(error) || 'Không thể tải dữ liệu học sinh'
    } finally {
      isBootstrapping.value = false
    }
  })

  onUnmounted(() => {
    if (copyTimeoutId) {
      clearTimeout(copyTimeoutId)
    }
  })

  return {
    schools,
    classes,
    selectedSchoolId,
    selectedClassId,
    students,
    searchQuery,
    isBootstrapping,
    isLoadingStudents,
    errorMessage,
    codeError,
    isFormModalOpen,
    isFormSubmitting,
    formError,
    formMode,
    formData,
    generatingCodeStudentId,
    revokingCodeStudentId,
    copiedStudentId,
    isRevokeConfirmOpen,
    revokeTarget,
    isDeleteConfirmOpen,
    deleteTarget,
    isDeleteLoading,
    isSuperAdmin,
    selectedSchoolName,
    selectedClassName,
    filteredStudents,
    hasInitialized,
    openAddModal,
    openEditModal,
    closeFormModal,
    updateFormData,
    submitForm,
    confirmDelete,
    closeDeleteConfirm,
    handleDelete,
    handleGenerateCode,
    confirmRevokeCode,
    closeRevokeConfirm,
    handleRevokeCode,
    handleCopyCode,
    loadStudents,
  }
}
