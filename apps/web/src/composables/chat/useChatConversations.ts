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

  async function renameGroup(conversationId, name) {
    const updated = await chatService.patchGroupConversation(conversationId, { name })
    await fetchConversations()
    return pickConversationAfterCreate(updated)
  }

  async function addGroupParticipants(conversationId, userIds) {
    const updated = await chatService.addConversationParticipants(conversationId, userIds)
    await fetchConversations()
    return pickConversationAfterCreate(updated)
  }

  /** Cập nhật sidebar khi có tin realtime (WebSocket); giữ hội thoại đang chọn đồng bộ. */
  function handleIncomingWsMessage(payload, currentUserId) {
    if (!payload?.conversation_id) return

    const cid = String(payload.conversation_id)
    const fromSelf = String(payload.sender_id) === String(currentUserId)
    const selId = getSelectedConversationId()
    const isSelected = selId && String(selId) === cid

    const lastMessage = {
      message_id: payload.message_id,
      sender_id: payload.sender_id,
      sender_email: payload.sender_email || '',
      content: payload.content,
      created_at: payload.created_at,
    }

    const idx = conversations.value.findIndex((c) => String(getConversationId(c)) === cid)
    if (idx === -1) {
      fetchConversations()
      return
    }

    const row = { ...conversations.value[idx], last_message: lastMessage }
    if (isSelected) {
      row.unread_count = 0
    } else if (!fromSelf) {
      row.unread_count = (conversations.value[idx].unread_count || 0) + 1
    }

    const next = [...conversations.value]
    next.splice(idx, 1)
    next.unshift(row)
    conversations.value = next

    if (isSelected && selectedConversation.value) {
      selectedConversation.value = {
        ...selectedConversation.value,
        last_message: lastMessage,
        unread_count: 0,
      }
    }

    if (isSelected && !fromSelf) {
      chatService.markConversationRead(cid).catch(() => {})
    }
  }

  /** Sau khi xóa thành viên: nếu user hiện tại không còn trong danh sách hội thoại, trả null. */
  async function removeGroupParticipant(conversationId, userId) {
    const updated = await chatService.removeConversationParticipant(conversationId, userId)
    await fetchConversations()
    const id = getConversationId(updated)
    const stillListed = conversations.value.some(
      (c) => String(getConversationId(c)) === String(id),
    )
    if (!stillListed) return null
    return pickConversationAfterCreate(updated)
  }

  function getSelectedConversationId() {
    return getConversationId(selectedConversation.value)
  }

  /** Sau khi mở hội thoại và tải tin (GET messages đã mark read trên server). */
  function syncUnreadAfterOpen(conversationId) {
    if (!conversationId) return
    const sid = String(conversationId)
    conversations.value = conversations.value.map((c) =>
      String(getConversationId(c)) === sid ? { ...c, unread_count: 0 } : c,
    )
    if (
      selectedConversation.value &&
      String(getConversationId(selectedConversation.value)) === sid
    ) {
      selectedConversation.value = { ...selectedConversation.value, unread_count: 0 }
    }
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
    renameGroup,
    addGroupParticipants,
    removeGroupParticipant,
    getSelectedConversationId,
    handleIncomingWsMessage,
    syncUnreadAfterOpen,
  }
}
