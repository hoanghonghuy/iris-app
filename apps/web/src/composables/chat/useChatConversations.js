import { ref } from 'vue'
import { chatService } from '../../services/chatService'
import { normalizeListResponse } from '../../helpers/collectionUtils'
import { getConversationId } from '../../helpers/chatHelpers'

export function useChatConversations() {
  const conversations = ref([])
  const selectedConversation = ref(null)
  const loading = ref(true)

  async function fetchConversations() {
    loading.value = true
    try {
      conversations.value = normalizeListResponse(await chatService.listConversations())
    } catch (error) {
      console.error('[chat] cannot load conversations', error)
      conversations.value = []
    } finally {
      loading.value = false
    }
  }

  function selectConversation(conversation) {
    selectedConversation.value = conversation
  }

  function clearSelection() {
    selectedConversation.value = null
  }

  async function createDirectConversation(userId) {
    if (!userId) return null
    try {
      const conversation = await chatService.createDirectConversation(userId)
      await fetchConversations()
      return conversation
    } catch (error) {
      console.error('[chat] cannot start conversation', error)
      return null
    }
  }

  function getSelectedConversationId() {
    return getConversationId(selectedConversation.value)
  }

  return {
    conversations,
    selectedConversation,
    loading,
    fetchConversations,
    selectConversation,
    clearSelection,
    createDirectConversation,
    getSelectedConversationId,
  }
}
