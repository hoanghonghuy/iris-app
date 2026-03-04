/**
 * Parent API Service
 * Quản lý các endpoint dành riêng cho phụ huynh (Parent Scope)
 * Tương ứng với: apps/api/internal/api/v1/handlers/parent_scope_handler.go
 */
import { apiClient } from './client';
import { Student, Post, ApiResponse } from '@/types';

export const parentApi = {
  /**
   * Lấy danh sách con của phụ huynh hiện tại
   * GET /api/v1/parent/children
   */
  getMyChildren: async () => {
    const res = await apiClient.get<ApiResponse<Student[]>>('/parent/children');
    return res.data.data;
  },

  /**
   * Lấy feed tổng hợp của tất cả con
   * GET /api/v1/parent/feed
   */
  getMyFeed: async (params?: { limit?: number; offset?: number }) => {
    const res = await apiClient.get<ApiResponse<Post[]>>('/parent/feed', { params });
    return res.data;
  },

  /**
   * Lấy tất cả bài đăng liên quan đến một đứa con cụ thể
   * GET /api/v1/parent/children/:student_id/posts
   */
  getChildPosts: async (studentId: string, params?: { limit?: number; offset?: number }) => {
    const res = await apiClient.get<ApiResponse<Post[]>>(`/parent/children/${studentId}/posts`, { params });
    return res.data;
  }
};