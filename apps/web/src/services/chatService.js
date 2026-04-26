import { httpClient } from './httpClient'

export const chatService = {
  async createDirectConversation(targetUserId) {
    const response = await httpClient.post('/chat/conversations/direct', {
      target_user_id: targetUserId,
    })
    return response?.data ?? response
  },

  async listConversations() {
    const response = await httpClient.get('/chat/conversations')
    return response?.data ?? response
  },

  async listMessages(conversationId, limit = 50, before) {
    const params = { limit }
    if (before) params.before = before
    return await httpClient.get(`/chat/conversations/${conversationId}/messages`, params)
  },

  async searchUsers(keyword) {
    const response = await httpClient.get('/chat/users/search', { q: keyword })
    return response?.data ?? response
  },
}

export function getChatWsUrl() {
  const apiBaseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'
  return `${apiBaseUrl.replace(/^http/, 'ws')}/chat/ws`
}
