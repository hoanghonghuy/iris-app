/**
 * Auth API Service
 * Quản lý các endpoint liên quan đến xác thực: login, register, activate...
 * Tương ứng với: apps/api/internal/api/v1/handlers/auth_handler.go
 */
import { apiClient } from './client';
import { LoginRequest, LoginResponse, RegisterParentRequest } from '@/types';

export const authApi = {
  /**
   * Đăng nhập người dùng
   * POST /api/v1/auth/login
   */
  login: async (data: LoginRequest) => {
    const res = await apiClient.post<LoginResponse>('/auth/login', data);
    return res.data;
  },

  /**
   * Lấy thông tin người dùng hiện tại từ token
   * GET /api/v1/me
   */
  getMe: async () => {
    const res = await apiClient.get('/me');
    return res.data;
  },

  /**
   * Phụ huynh tự đăng ký bằng mã code
   * POST /api/v1/register/parent
   */
  registerParent: async (data: RegisterParentRequest) => {
    const res = await apiClient.post('/register/parent', data);
    return res.data;
  }
};