import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useAuthStore } from '../../stores/authStore'
import { adminService } from '../../services/adminService'
import { normalizeListResponse } from '../../helpers/collectionUtils'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { ADMIN_SELECTOR_FETCH_LIMIT } from '../../helpers/adminConfig'

export function useAdminPeopleManagement(options) {
  const {
    pageSize = 20,
    searchFields = [],
    fetchList,
    createInitialEditForm,
    toEditForm,
    validateEditForm,
    updateItem,
    updateErrorMessage,
    assignItem,
    assignErrorMessage = '',
    enableStudentSelector = false,
    toUnassignTarget,
    unassignItem,
    unassignErrorPrefix = 'Lỗi hủy gán',
  } = options

  const authStore = useAuthStore()

  const items = ref([])
  const totalPages = ref(0)
  const currentPage = ref(1)
  const totalItems = ref(0)
  const isLoading = ref(true)
  const errorMessage = ref('')

  const searchQuery = ref('')
  const isSuperAdmin = computed(() => authStore.currentUserRole === 'SUPER_ADMIN')

  const isEditModalOpen = ref(false)
  const editLoading = ref(false)
  const editError = ref('')
  const editForm = ref(createInitialEditForm())

  const isAssignModalOpen = ref(false)
  const assignTarget = ref(null)
  const assignLoading = ref(false)
  const assignError = ref('')

  const isUnassignOpen = ref(false)
  const unassignTarget = ref(null)
  const unassignLoading = ref(false)

  const schools = ref([])
  const classes = ref([])
  const students = ref([])
  const selectedSchoolId = ref('')
  const selectedClassId = ref('')
  const selectedStudentId = ref('')

  let activeSchoolsController = null
  let activeClassesController = null
  let activeStudentsController = null

  const filteredItems = computed(() => {
    const query = searchQuery.value.trim().toLowerCase()
    if (!query) {
      return items.value
    }

    return items.value.filter((item) =>
      searchFields.some((field) => String(item?.[field] || '').toLowerCase().includes(query)),
    )
  })

  async function fetchItems(page = 1) {
    isLoading.value = true
    errorMessage.value = ''
    currentPage.value = page

    try {
      const data = await fetchList({
        limit: pageSize,
        offset: (page - 1) * pageSize,
      })

      items.value = normalizeListResponse(data)

      if (data?.pagination) {
        const limit = Math.max(data.pagination.limit || pageSize, 1)
        totalItems.value = data.pagination.total || 0
        totalPages.value = Math.ceil(totalItems.value / limit) || 1
      } else {
        totalItems.value = items.value.length
        totalPages.value = items.value.length > 0 ? 1 : 0
      }
    } catch (error) {
      errorMessage.value = extractErrorMessage(error)
    } finally {
      isLoading.value = false
    }
  }

  function resetStudentSelectorState() {
    students.value = []
    selectedStudentId.value = ''
  }

  function resetClassSelectorState() {
    classes.value = []
    selectedClassId.value = ''
    if (enableStudentSelector) {
      resetStudentSelectorState()
    }
  }

  function cancelActiveSchoolsRequest() {
    if (activeSchoolsController) {
      activeSchoolsController.abort()
      activeSchoolsController = null
    }
  }

  function cancelActiveClassesRequest() {
    if (activeClassesController) {
      activeClassesController.abort()
      activeClassesController = null
    }
  }

  function cancelActiveStudentsRequest() {
    if (activeStudentsController) {
      activeStudentsController.abort()
      activeStudentsController = null
    }
  }

  function cancelActiveSelectorRequests() {
    cancelActiveSchoolsRequest()
    cancelActiveClassesRequest()
    cancelActiveStudentsRequest()
  }

  async function fetchSchoolsForSelector() {
    cancelActiveSchoolsRequest()
    const controller = new AbortController()
    activeSchoolsController = controller

    try {
      const data = await adminService.getSchools(
        { limit: ADMIN_SELECTOR_FETCH_LIMIT, offset: 0 },
        { signal: controller.signal },
      )

      if (controller.signal.aborted) {
        return
      }

      schools.value = normalizeListResponse(data)
      selectedSchoolId.value = schools.value[0]?.school_id || ''
    } catch (error) {
      if (controller.signal.aborted || error?.name === 'AbortError') {
        return
      }

      schools.value = []
      selectedSchoolId.value = ''
    } finally {
      if (activeSchoolsController === controller) {
        activeSchoolsController = null
      }
    }
  }

  async function fetchClassesForSelector() {
    cancelActiveClassesRequest()

    if (!selectedSchoolId.value) {
      resetClassSelectorState()
      return
    }

    resetClassSelectorState()
    const controller = new AbortController()
    activeClassesController = controller

    try {
      const data = await adminService.getClassesBySchool(
        selectedSchoolId.value,
        {
          limit: ADMIN_SELECTOR_FETCH_LIMIT,
          offset: 0,
        },
        { signal: controller.signal },
      )

      if (controller.signal.aborted) {
        return
      }

      classes.value = normalizeListResponse(data)
      selectedClassId.value = classes.value[0]?.class_id || ''
    } catch (error) {
      if (controller.signal.aborted || error?.name === 'AbortError') {
        return
      }

      classes.value = []
    } finally {
      if (activeClassesController === controller) {
        activeClassesController = null
      }
    }
  }

  async function fetchStudentsForSelector() {
    if (!enableStudentSelector) return

    cancelActiveStudentsRequest()

    if (!selectedClassId.value) {
      resetStudentSelectorState()
      return
    }

    resetStudentSelectorState()
    const controller = new AbortController()
    activeStudentsController = controller

    try {
      const data = await adminService.getStudentsByClass(
        selectedClassId.value,
        {
          limit: ADMIN_SELECTOR_FETCH_LIMIT,
          offset: 0,
        },
        { signal: controller.signal },
      )

      if (controller.signal.aborted) {
        return
      }

      students.value = normalizeListResponse(data)
      selectedStudentId.value = students.value[0]?.student_id || ''
    } catch (error) {
      if (controller.signal.aborted || error?.name === 'AbortError') {
        return
      }

      students.value = []
    } finally {
      if (activeStudentsController === controller) {
        activeStudentsController = null
      }
    }
  }

  watch(selectedSchoolId, () => {
    fetchClassesForSelector()
  })

  if (enableStudentSelector) {
    watch(selectedClassId, () => {
      fetchStudentsForSelector()
    })
  }

  onMounted(() => {
    fetchItems()
    fetchSchoolsForSelector()
  })

  onBeforeUnmount(() => {
    cancelActiveSelectorRequests()
  })

  function openAssignModal(item) {
    assignTarget.value = item
    assignError.value = ''
    isAssignModalOpen.value = true
  }

  function closeAssignModal() {
    isAssignModalOpen.value = false
    assignError.value = ''
  }

  function openEditModal(item) {
    editError.value = ''
    editForm.value = toEditForm(item, {
      selectedSchoolId: selectedSchoolId.value,
    })
    isEditModalOpen.value = true
  }

  function closeEditModal() {
    isEditModalOpen.value = false
    editError.value = ''
  }

  async function handleEdit() {
    const validationMessage = validateEditForm(editForm.value)
    if (validationMessage) {
      editError.value = validationMessage
      return
    }

    editLoading.value = true
    editError.value = ''
    try {
      await updateItem(editForm.value)
      closeEditModal()
      await fetchItems(currentPage.value)
    } catch (error) {
      editError.value = extractErrorMessage(error) || updateErrorMessage
    } finally {
      editLoading.value = false
    }
  }

  async function handleAssign() {
    const selectedId = enableStudentSelector ? selectedStudentId.value : selectedClassId.value
    if (!selectedId || !assignTarget.value) {
      return
    }

    assignLoading.value = true
    assignError.value = ''
    try {
      await assignItem({
        target: assignTarget.value,
        selectedSchoolId: selectedSchoolId.value,
        selectedClassId: selectedClassId.value,
        selectedStudentId: selectedStudentId.value,
      })
      closeAssignModal()
      await fetchItems(currentPage.value)
    } catch (error) {
      assignError.value = extractErrorMessage(error) || assignErrorMessage
    } finally {
      assignLoading.value = false
    }
  }

  function openUnassignDialog(item, relation) {
    unassignTarget.value = toUnassignTarget(item, relation)
    isUnassignOpen.value = true
  }

  function closeUnassignDialog() {
    isUnassignOpen.value = false
    unassignTarget.value = null
  }

  async function handleUnassign() {
    if (!unassignTarget.value) {
      return
    }

    unassignLoading.value = true
    try {
      await unassignItem(unassignTarget.value)
      closeUnassignDialog()
      await fetchItems(currentPage.value)
    } catch (error) {
      errorMessage.value = `${unassignErrorPrefix}: ${extractErrorMessage(error)}`
      closeUnassignDialog()
    } finally {
      unassignLoading.value = false
    }
  }

  return {
    items,
    totalPages,
    currentPage,
    totalItems,
    isLoading,
    errorMessage,
    searchQuery,
    isSuperAdmin,
    isEditModalOpen,
    editLoading,
    editError,
    editForm,
    isAssignModalOpen,
    assignTarget,
    assignLoading,
    assignError,
    isUnassignOpen,
    unassignTarget,
    unassignLoading,
    schools,
    classes,
    students,
    selectedSchoolId,
    selectedClassId,
    selectedStudentId,
    filteredItems,
    fetchItems,
    openAssignModal,
    closeAssignModal,
    openEditModal,
    closeEditModal,
    handleEdit,
    handleAssign,
    openUnassignDialog,
    closeUnassignDialog,
    handleUnassign,
  }
}
