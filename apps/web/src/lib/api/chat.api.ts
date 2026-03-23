/**
 * Chat API Service
 * Quản lý các endpoint cho hệ thống chat realtime
 * Tương ứng với: apps/api/internal/api/v1/handlers/chat_handler.go
 */
import { apiClient } from './client';
import { ApiResponse, Conversation, Message } from '@/types';

export const chatApi = {
  /**
   * Tạo hoặc tìm cuộc hội thoại direct (1-1) với user khác
   * POST /api/v1/chat/conversations/direct
   */
  createDirectConversation: async (targetUserId: string) => {
    const res = await apiClient.post<ApiResponse<Conversation>>('/chat/conversations/direct', {
      target_user_id: targetUserId,
    });
    return res.data.data;
  },

  /**
   * Lấy danh sách cuộc hội thoại của user hiện tại
   * GET /api/v1/chat/conversations
   */
  listConversations: async () => {
    const res = await apiClient.get<ApiResponse<Conversation[]>>('/chat/conversations');
    return res.data.data;
  },

  /**
   * Lấy danh sách tin nhắn trong cuộc hội thoại
   * GET /api/v1/chat/conversations/:conversation_id/messages
   */
  listMessages: async (conversationId: string, limit = 50, offset = 0) => {
    const res = await apiClient.get<ApiResponse<Message[]>>(
      `/chat/conversations/${conversationId}/messages`,
      { params: { limit, offset } }
    );
    return res.data.data;
  },

  /**
   * Tìm kiếm user theo email để bắt đầu cuộc trò chuyện
   * GET /api/v1/chat/users/search?q=keyword
   */
  searchUsers: async (keyword: string) => {
    const res = await apiClient.get<ApiResponse<{ user_id: string; email: string; full_name: string }[]>>(
      `/chat/users/search`,
      { params: { q: keyword } }
    );
    return res.data.data;
  },
};
