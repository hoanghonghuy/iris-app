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

  /** Ưu tiên object từ list sau refresh để có participants đầy đủ cho UI */
  function pickConversationAfterCreate(payload) {
    const id = getConversationId(payload)
    if (!id) return payload
    return (
      conversations.value.find((c) => String(getConversationId(c)) === String(id)) ?? payload
    )
  }

  async function createDirectConversation(userId) {
    if (!userId) return null
    try {
      const conversation = await chatService.createDirectConversation(userId)
      await fetchConversations()
      return pickConversationAfterCreate(conversation)
    } catch (error) {
      console.error('[chat] cannot start conversation', error)
      throw error
    }
  }

  async function createGroupConversation(name, participantUserIds) {
    if (!participantUserIds?.length) return null
    try {
      const conversation = await chatService.createGroupConversation({
        name,
        participantUserIds,
      })
      await fetchConversations()
      return pickConversationAfterCreate(conversation)
    } catch (error) {
      console.error('[chat] cannot create group', error)
      throw error
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
    createGroupConversation,
    getSelectedConversationId,
  }
}
