import { ref, watch, onUnmounted } from 'vue'
import { chatService } from '../../services/chatService'
import { normalizeListResponse } from '../../helpers/collectionUtils'

export function useChatSearch() {
  const searchQuery = ref('')
  const searchResults = ref([])
  const showNewConversation = ref(false)
  let searchTimer = null

  function toggleNewConversation() {
    showNewConversation.value = !showNewConversation.value
    if (!showNewConversation.value) {
      searchQuery.value = ''
      searchResults.value = []
    }
  }

  function clearSearch() {
    searchQuery.value = ''
    searchResults.value = []
  }

  watch(searchQuery, (query) => {
    clearTimeout(searchTimer)
    if (!query.trim()) {
      searchResults.value = []
      return
    }

    searchTimer = setTimeout(async () => {
      try {
        searchResults.value = normalizeListResponse(await chatService.searchUsers(query.trim()))
      } catch {
        searchResults.value = []
      }
    }, 500)
  })

  onUnmounted(() => {
    clearTimeout(searchTimer)
  })

  return {
    searchQuery,
    searchResults,
    showNewConversation,
    toggleNewConversation,
    clearSearch,
  }
}
