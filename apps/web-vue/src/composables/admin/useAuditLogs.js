import { ref, computed } from 'vue'
import { adminService } from '../../services/adminService'
import { extractErrorMessage } from '../../helpers/errorHandler'
import { normalizePaginatedResponse } from '../../helpers/collectionUtils'
import { ADMIN_DEFAULT_PAGE_SIZE } from '../../helpers/adminConfig'

export function useAuditLogs(isSuperAdmin) {
  const auditLogs = ref([])
  const totalPages = ref(0)
  const currentPage = ref(1)
  const totalItems = ref(0)
  const pageSize = ref(ADMIN_DEFAULT_PAGE_SIZE)
  const isLoading = ref(true)
  const errorMessage = ref('')

  const searchQuery = ref('')
  const actionQuery = ref('')
  const selectedEntityType = ref('')
  const fromDate = ref('')
  const toDate = ref('')

  const appliedSearchQuery = ref('')
  const appliedActionQuery = ref('')
  const appliedEntityType = ref('')
  const appliedFromDate = ref('')
  const appliedToDate = ref('')

  async function fetchAuditLogs(page = 1) {
    if (!isSuperAdmin.value) {
      auditLogs.value = []
      totalPages.value = 0
      totalItems.value = 0
      currentPage.value = 1
      isLoading.value = false
      errorMessage.value = 'Bạn không có quyền truy cập nhật ký hệ thống.'
      return
    }

    isLoading.value = true
    errorMessage.value = ''
    currentPage.value = page

    try {
      const params = {
        limit: pageSize.value,
        offset: (page - 1) * pageSize.value,
        q: appliedSearchQuery.value || undefined,
        action: appliedActionQuery.value || undefined,
        entity_type: appliedEntityType.value || undefined,
        from: appliedFromDate.value ? new Date(appliedFromDate.value).toISOString() : undefined,
        to: appliedToDate.value ? new Date(appliedToDate.value).toISOString() : undefined,
      }

      const data = await adminService.getAuditLogs(params)
      const normalized = normalizePaginatedResponse(data, pageSize.value)
      const hasPagination = Boolean(data?.pagination)

      auditLogs.value = normalized.items

      if (hasPagination) {
        totalItems.value = normalized.pagination.total || 0
        totalPages.value = Math.ceil(totalItems.value / Math.max(normalized.pagination.limit || pageSize.value, 1)) || 1
        pageSize.value = normalized.pagination.limit || pageSize.value
      } else {
        totalItems.value = auditLogs.value.length
        totalPages.value = auditLogs.value.length > 0 ? 1 : 0
      }
    } catch (error) {
      errorMessage.value = extractErrorMessage(error) || 'Không thể tải dữ liệu nhật ký hệ thống'
    } finally {
      isLoading.value = false
    }
  }

  function applyFilters() {
    appliedSearchQuery.value = searchQuery.value.trim()
    appliedActionQuery.value = actionQuery.value.trim()
    appliedEntityType.value = selectedEntityType.value
    appliedFromDate.value = fromDate.value
    appliedToDate.value = toDate.value
    fetchAuditLogs(1)
  }

  function resetFilters() {
    searchQuery.value = ''
    actionQuery.value = ''
    selectedEntityType.value = ''
    fromDate.value = ''
    toDate.value = ''
    appliedSearchQuery.value = ''
    appliedActionQuery.value = ''
    appliedEntityType.value = ''
    appliedFromDate.value = ''
    appliedToDate.value = ''
    fetchAuditLogs(1)
  }

  function handlePageChange(page) {
    fetchAuditLogs(page)
  }

  function handlePageSizeChange(newSize) {
    pageSize.value = newSize
    fetchAuditLogs(1)
  }

  return {
    auditLogs,
    totalPages,
    currentPage,
    totalItems,
    pageSize,
    isLoading,
    errorMessage,
    searchQuery,
    actionQuery,
    selectedEntityType,
    fromDate,
    toDate,
    appliedSearchQuery,
    appliedActionQuery,
    appliedEntityType,
    appliedFromDate,
    appliedToDate,
    fetchAuditLogs,
    applyFilters,
    resetFilters,
    handlePageChange,
    handlePageSizeChange,
  }
}
