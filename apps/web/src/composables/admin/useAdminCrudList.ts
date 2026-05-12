import { ref } from 'vue'
import { normalizeListResponse } from '../../helpers/collectionUtils'
import { extractErrorMessage } from '../../helpers/errorHandler'

interface FetchPageParams {
  page: number
  pageSize: number
}

interface PaginationPayload {
  total?: number
  limit?: number
}

interface ListResponse<TItem> {
  data?: TItem[]
  pagination?: PaginationPayload
}

interface AfterSaveContext<TForm> {
  mode: 'add' | 'edit'
  form: TForm
  currentPage: number
  fetchItems: (page?: number) => Promise<void>
}

interface AfterDeleteContext<TItem> {
  item: TItem
  currentPage: number
  currentItemsCount: number
  fetchItems: (page?: number) => Promise<void>
}

interface AdminCrudListOptions<TItem, TForm> {
  pageSize?: number
  fetchPage: (params: FetchPageParams) => Promise<ListResponse<TItem> | TItem[] | unknown>
  createEmptyForm: () => TForm
  toEditForm: (item: TItem) => TForm
  validateForm: (form: TForm, mode: 'add' | 'edit') => string
  createItem?: (form: TForm) => Promise<unknown>
  updateItem?: (form: TForm) => Promise<unknown>
  deleteItem?: (item: TItem) => Promise<unknown>
  saveErrorMessage?: string
  deleteErrorPrefix?: string
  onAfterSave?: (ctx: AfterSaveContext<TForm>) => Promise<boolean> | boolean
  onAfterDelete?: (ctx: AfterDeleteContext<TItem>) => Promise<boolean> | boolean
}

function hasPagination<TItem>(
  data: ListResponse<TItem> | TItem[] | unknown,
): data is ListResponse<TItem> {
  return Boolean(data) && typeof data === 'object' && 'pagination' in data
}

export function useAdminCrudList<TItem, TForm>(options: AdminCrudListOptions<TItem, TForm>) {
  const {
    pageSize = 10,
    fetchPage,
    createEmptyForm,
    toEditForm,
    validateForm,
    createItem,
    updateItem,
    deleteItem,
    saveErrorMessage = 'Không thể lưu dữ liệu',
    deleteErrorPrefix = 'Lỗi xóa',
    onAfterSave,
    onAfterDelete,
  } = options

  const items = ref<TItem[]>([])
  const totalPages = ref(0)
  const currentPage = ref(1)
  const totalItems = ref(0)
  const isLoading = ref(true)
  const errorMessage = ref('')

  const isModalOpen = ref(false)
  const isSubmitting = ref(false)
  const modalError = ref('')
  const formMode = ref<'add' | 'edit'>('add')
  const formData = ref<TForm>(createEmptyForm())

  const isConfirmOpen = ref(false)
  const itemToDelete = ref<TItem | null>(null)

  async function fetchItems(page = 1) {
    isLoading.value = true
    errorMessage.value = ''
    currentPage.value = page

    try {
      const data = await fetchPage({
        page,
        pageSize,
      })

      items.value = normalizeListResponse(data)

      if (hasPagination<TItem>(data) && data.pagination) {
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

  function openAddModal() {
    formMode.value = 'add'
    formData.value = createEmptyForm()
    modalError.value = ''
    isModalOpen.value = true
  }

  function closeModal() {
    isModalOpen.value = false
    modalError.value = ''
  }

  function openEditModal(item: TItem) {
    formMode.value = 'edit'
    formData.value = toEditForm(item)
    modalError.value = ''
    isModalOpen.value = true
  }

  async function handleSave() {
    const validationMessage = validateForm(formData.value, formMode.value)
    if (validationMessage) {
      modalError.value = validationMessage
      return
    }

    isSubmitting.value = true
    modalError.value = ''

    try {
      if (formMode.value === 'add') {
        if (typeof createItem !== 'function') {
          throw new Error('Thiếu handler createItem')
        }

        await createItem(formData.value)
      } else {
        if (typeof updateItem !== 'function') {
          throw new Error('Chức năng chỉnh sửa chưa được hỗ trợ')
        }

        await updateItem(formData.value)
      }

      closeModal()

      const handled = onAfterSave
        ? await onAfterSave({
            mode: formMode.value,
            form: formData.value,
            currentPage: currentPage.value,
            fetchItems,
          })
        : false

      if (!handled) {
        await fetchItems(currentPage.value)
      }
    } catch (error) {
      modalError.value = extractErrorMessage(error) || saveErrorMessage
    } finally {
      isSubmitting.value = false
    }
  }

  function confirmDelete(item: TItem) {
    itemToDelete.value = item
    isConfirmOpen.value = true
  }

  function closeDeleteConfirm() {
    isConfirmOpen.value = false
    itemToDelete.value = null
  }

  async function handleDelete() {
    if (!itemToDelete.value) {
      return
    }

    if (typeof deleteItem !== 'function') {
      errorMessage.value = `${deleteErrorPrefix}: Chức năng xóa chưa được hỗ trợ`
      closeDeleteConfirm()
      return
    }

    isSubmitting.value = true
    const targetItem = itemToDelete.value

    try {
      await deleteItem(targetItem)
      closeDeleteConfirm()

      const handled = onAfterDelete
        ? await onAfterDelete({
            item: targetItem,
            currentPage: currentPage.value,
            currentItemsCount: items.value.length,
            fetchItems,
          })
        : false

      if (!handled) {
        await fetchItems(currentPage.value)
      }
    } catch (error) {
      errorMessage.value = `${deleteErrorPrefix}: ${extractErrorMessage(error)}`
      closeDeleteConfirm()
    } finally {
      isSubmitting.value = false
    }
  }

  return {
    items,
    totalPages,
    currentPage,
    totalItems,
    isLoading,
    errorMessage,
    isModalOpen,
    isSubmitting,
    modalError,
    formMode,
    formData,
    isConfirmOpen,
    itemToDelete,
    fetchItems,
    openAddModal,
    closeModal,
    openEditModal,
    handleSave,
    confirmDelete,
    closeDeleteConfirm,
    handleDelete,
  }
}
