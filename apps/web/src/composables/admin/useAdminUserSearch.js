import { ref } from 'vue'
import { adminService } from '../../services/adminService'
import { ADMIN_SELECTOR_FETCH_LIMIT } from '../../helpers/adminConfig'

export function useAdminUserSearch(options = {}) {
  const {
    minLength = 2,
    resultLimit = 6,
  } = options

  const searchQuery = ref('')
  const searchResults = ref([])
  const searchLoading = ref(false)
  const selectedUser = ref(null)

  let searchTimerId = null
  let activeSearchController = null

  function resetSearchState() {
    searchResults.value = []
    searchLoading.value = false
  }

  function cancelActiveSearch() {
    if (activeSearchController) {
      activeSearchController.abort()
      activeSearchController = null
    }
  }

  async function searchUsers(query) {
    const normalizedQuery = query.trim().toLowerCase()

    if (normalizedQuery.length < minLength) {
      cancelActiveSearch()
      resetSearchState()
      return
    }

    cancelActiveSearch()
    const controller = new AbortController()
    activeSearchController = controller

    searchLoading.value = true

    try {
      const matches = []
      let offset = 0
      let hasMore = true

      while (hasMore && matches.length < resultLimit) {
        const response = await adminService.getUsers({
          limit: ADMIN_SELECTOR_FETCH_LIMIT,
          offset,
        }, {
          signal: controller.signal,
        })

        if (controller.signal.aborted) {
          return
        }

        const users = response?.data ?? []
        if (!Array.isArray(users) || users.length === 0) {
          hasMore = false
          break
        }

        for (const user of users) {
          if (matches.length >= resultLimit) {
            hasMore = false
            break
          }

          const email = (user.email || '').toLowerCase()
          const fullName = (user.full_name || '').toLowerCase()
          const userId = (user.user_id || '').toLowerCase()

          if (
            email.includes(normalizedQuery) ||
            fullName.includes(normalizedQuery) ||
            userId.includes(normalizedQuery)
          ) {
            matches.push(user)
          }
        }

        offset += ADMIN_SELECTOR_FETCH_LIMIT
        hasMore = hasMore && Boolean(response?.pagination?.has_more)
      }

      if (!controller.signal.aborted) {
        searchResults.value = matches
      }
    } catch (error) {
      if (error.name !== 'AbortError') {
        console.error('User search error:', error)
        searchResults.value = []
      }
    } finally {
      if (!controller.signal.aborted) {
        searchLoading.value = false
      }
    }
  }

  function selectUser(user) {
    selectedUser.value = user
    searchQuery.value = user.email || user.full_name || user.user_id
  }

  function clearSelectedUser() {
    selectedUser.value = null
    searchQuery.value = ''
    resetSearchState()
  }

  function cleanup() {
    clearTimeout(searchTimerId)
    cancelActiveSearch()
  }

  return {
    searchQuery,
    searchResults,
    searchLoading,
    selectedUser,
    searchUsers,
    selectUser,
    clearSelectedUser,
    resetSearchState,
    cancelActiveSearch,
    cleanup,
  }
}
