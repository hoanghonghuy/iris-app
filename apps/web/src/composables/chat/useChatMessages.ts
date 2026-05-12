import { ref, nextTick } from 'vue'
import { chatService } from '../../services/chatService'
import { normalizeListResponse } from '../../helpers/collectionUtils'

export function useChatMessages() {
  const messages = ref([])
  const loadingMessages = ref(false)
  const loadingMore = ref(false)
  const hasMore = ref(false)
  const nextCursor = ref(null)

  async function loadMessages(conversationId) {
    if (!conversationId) return

    messages.value = []
    nextCursor.value = null
    hasMore.value = false
    loadingMessages.value = true

    try {
      const response = await chatService.listMessages(conversationId)
      messages.value = normalizeListResponse(response).reverse()
      nextCursor.value = response?.next_cursor ?? null
      hasMore.value = Boolean(response?.has_more)
    } catch (error) {
      console.error('[chat] cannot load messages', error)
    } finally {
      loadingMessages.value = false
    }
  }

  async function loadOlderMessages(conversationId, messagesContainer) {
    if (
      !messagesContainer ||
      !conversationId ||
      !hasMore.value ||
      loadingMore.value ||
      !nextCursor.value
    ) {
      return
    }
    if (messagesContainer.scrollTop >= 80) return

    loadingMore.value = true
    const previousHeight = messagesContainer.scrollHeight

    try {
      const response = await chatService.listMessages(conversationId, 50, nextCursor.value)
      messages.value = normalizeListResponse(response).reverse().concat(messages.value)
      nextCursor.value = response?.next_cursor ?? null
      hasMore.value = Boolean(response?.has_more)
      await nextTick()
      messagesContainer.scrollTop = messagesContainer.scrollHeight - previousHeight
    } catch (error) {
      console.error('[chat] cannot load older messages', error)
    } finally {
      loadingMore.value = false
    }
  }

  function addMessage(message) {
    if (messages.value.some((item) => item.message_id === message.message_id)) return
    messages.value = [...messages.value, message]
  }

  function clearMessages() {
    messages.value = []
    nextCursor.value = null
    hasMore.value = false
  }

  return {
    messages,
    loadingMessages,
    loadingMore,
    hasMore,
    loadMessages,
    loadOlderMessages,
    addMessage,
    clearMessages,
  }
}
