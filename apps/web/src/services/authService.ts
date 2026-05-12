import { httpClient } from './httpClient'
import type { LoginResponse, User, ApiResponse } from '@/types'

export const authService = {
  // Đăng nhập người dùng
  // POST /api/v1/auth/login
  async login(email: string, password: string): Promise<ApiResponse<LoginResponse>> {
    return await httpClient.post<ApiResponse<LoginResponse>>('/auth/login', { email, password })
  },

  // Đăng nhập bằng Google ID token
  // POST /api/v1/auth/login/google
  async loginWithGoogle(idToken: string): Promise<ApiResponse<LoginResponse>> {
    return await httpClient.post<ApiResponse<LoginResponse>>('/auth/login/google', {
      id_token: idToken,
    })
  },

  // Liên kết mật khẩu với Google ID token (khi GOOGLE_LINK_PASSWORD_REQUIRED)
  // POST /api/v1/auth/login/google
  async linkGooglePassword(
    idToken: string,
    password: string,
  ): Promise<ApiResponse<LoginResponse>> {
    return await httpClient.post<ApiResponse<LoginResponse>>('/auth/login/google', {
      id_token: idToken,
      password,
    })
  },

  // Lấy thông tin người dùng hiện tại từ token
  // GET /api/v1/me
  // BE trả về { data: { id, email, roles: [...] } }
  async getMe(): Promise<User> {
    const response = await httpClient.get<ApiResponse<User>>('/me')
    // httpClient trả toàn bộ JSON body, BE wrap trong { data: ... }
    return response?.data ?? response
  },

  // Kích hoạt tài khoản bằng token
  // POST /api/v1/users/activate-token
  async activateWithToken(token: string, password: string): Promise<ApiResponse> {
    return await httpClient.post<ApiResponse>('/users/activate-token', { token, password })
  },

  // Phụ huynh tự đăng ký bằng mã code
  // POST /api/v1/register/parent
  async registerParent(
    email: string,
    password: string,
    parentCode: string,
  ): Promise<ApiResponse<LoginResponse>> {
    return await httpClient.post<ApiResponse<LoginResponse>>('/register/parent', {
      email,
      password,
      parent_code: parentCode,
    })
  },

  // Phụ huynh tự đăng ký bằng Google ID token và mã code
  // POST /api/v1/register/parent/google
  async registerParentWithGoogle(
    idToken: string,
    parentCode: string,
  ): Promise<ApiResponse<LoginResponse>> {
    return await httpClient.post<ApiResponse<LoginResponse>>('/register/parent/google', {
      id_token: idToken,
      parent_code: parentCode,
    })
  },

  // Yêu cầu đặt lại mật khẩu (gửi email)
  // POST /api/v1/auth/forgot-password
  async forgotPassword(email: string): Promise<ApiResponse> {
    return await httpClient.post<ApiResponse>('/auth/forgot-password', { email })
  },

  // Đặt lại mật khẩu bằng token
  // POST /api/v1/auth/reset-password
  async resetPassword(email: string, token: string, password: string): Promise<ApiResponse> {
    return await httpClient.post<ApiResponse>('/auth/reset-password', { email, token, password })
  },

  async updateMyPassword(password: string): Promise<ApiResponse> {
    return await httpClient.put<ApiResponse>('/me/password', { password })
  },

  // Làm mới access token bằng refresh token
  // POST /api/v1/auth/refresh
  async refresh(refreshToken: string): Promise<ApiResponse<LoginResponse>> {
    return await httpClient.post<ApiResponse<LoginResponse>>('/auth/refresh', {
      refresh_token: refreshToken,
    })
  },
}
