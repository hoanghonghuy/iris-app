/**
 * Auth API Service
 * Quản lý các endpoint liên quan đến xác thực: login, register, activate...
 * Tương ứng với: apps/api/internal/api/v1/handlers/auth_handler.go + user_handler.go
 */
import { apiClient } from './client';
import { GoogleLoginRequest, LoginRequest, LoginResponse, RegisterParentRequest } from '@/types';

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
   * Đăng nhập bằng Google ID token
   * POST /api/v1/auth/login/google
   */
  loginWithGoogle: async (data: GoogleLoginRequest) => {
    const res = await apiClient.post<LoginResponse>('/auth/login/google', data);
    return res.data;
  },

  /**
   * Lấy thông tin người dùng hiện tại từ token
   * GET /api/v1/me
   */
  getMe: async () => {
    const res = await apiClient.get<{data: any}>('/me');
    return res.data.data;
  },

  /**
   * Kích hoạt tài khoản bằng token
   * POST /api/v1/users/activate-token
   */
  activateWithToken: async (data: { token: string; password: string }) => {
    const res = await apiClient.post('/users/activate-token', data);
    return res.data;
  },

  /**
   * Cập nhật mật khẩu (self-service)
   * PUT /api/v1/me/password
   */
  updateMyPassword: async (password: string) => {
    const res = await apiClient.put('/me/password', { password });
    return res.data;
  },

  /**
   * Phụ huynh tự đăng ký bằng mã code
   * POST /api/v1/register/parent
   */
  registerParent: async (data: RegisterParentRequest) => {
    const res = await apiClient.post('/register/parent', data);
    return res.data;
  },

  /**
   * Yêu cầu đặt lại mật khẩu (gửi email)
   * POST /api/v1/auth/forgot-password
   */
  forgotPassword: async (email: string) => {
    const res = await apiClient.post('/auth/forgot-password', { email });
    return res.data;
  },

  /**
   * Đặt lại mật khẩu bằng token
   * POST /api/v1/auth/reset-password
   */
  resetPassword: async (email: string, token: string, password: string) => {
    const res = await apiClient.post('/auth/reset-password', { email, token, password });
    return res.data;
  },
};