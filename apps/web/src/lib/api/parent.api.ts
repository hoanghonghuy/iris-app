/**
 * Parent API Service
 * Quản lý các endpoint dành riêng cho phụ huynh (Parent Scope)
 * Tương ứng với: apps/api/internal/api/v1/handlers/parent_scope_handler.go
 */
import { apiClient } from './client';
import { Student, Post, ApiResponse, PaginationParams, PostComment, PostLikeResponse, PostShareResponse, CreatePostCommentRequest, CreatePostCommentResponse } from '@/types';

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
  },

  /**
   * Toggle like cho bài đăng
   * POST /api/v1/parent/posts/:post_id/like
   */
  togglePostLike: async (postId: string) => {
    const res = await apiClient.post<ApiResponse<PostLikeResponse>>(`/parent/posts/${postId}/like`);
    return res.data.data;
  },

  /**
   * Lấy danh sách bình luận của bài đăng
   * GET /api/v1/parent/posts/:post_id/comments
   */
  getPostComments: async (postId: string, params?: PaginationParams) => {
    const res = await apiClient.get<ApiResponse<PostComment[]>>(`/parent/posts/${postId}/comments`, { params });
    return res.data;
  },

  /**
   * Tạo bình luận bài đăng
   * POST /api/v1/parent/posts/:post_id/comments
   */
  createPostComment: async (postId: string, data: CreatePostCommentRequest) => {
    const res = await apiClient.post<ApiResponse<CreatePostCommentResponse>>(`/parent/posts/${postId}/comments`, data);
    return res.data.data;
  },

  /**
   * Ghi nhận chia sẻ bài đăng
   * POST /api/v1/parent/posts/:post_id/share
   */
  sharePost: async (postId: string) => {
    const res = await apiClient.post<ApiResponse<PostShareResponse>>(`/parent/posts/${postId}/share`);
    return res.data.data;
  },
};