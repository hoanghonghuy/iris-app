import { httpClient } from './httpClient'

export const authService = {
  // Đăng nhập người dùng
  // POST /api/v1/auth/login
  async login(email, password) {
    return await httpClient.post('/auth/login', { email, password })
  },

  // Đăng nhập bằng Google ID token
  // POST /api/v1/auth/login/google
  async loginWithGoogle(idToken) {
    return await httpClient.post('/auth/login/google', { id_token: idToken })
  },

  // Liên kết mật khẩu với Google ID token (khi GOOGLE_LINK_PASSWORD_REQUIRED)
  // POST /api/v1/auth/login/google
  async linkGooglePassword(idToken, password) {
    return await httpClient.post('/auth/login/google', { id_token: idToken, password })
  },

  // Lấy thông tin người dùng hiện tại từ token
  // GET /api/v1/me
  // BE trả về { data: { id, email, roles: [...] } }
  async getMe() {
    const response = await httpClient.get('/me')
    // httpClient trả toàn bộ JSON body, BE wrap trong { data: ... }
    return response?.data ?? response
  },

  // Kích hoạt tài khoản bằng token
  // POST /api/v1/users/activate-token
  async activateWithToken(token, password) {
    return await httpClient.post('/users/activate-token', { token, password })
  },

  // Phụ huynh tự đăng ký bằng mã code
  // POST /api/v1/register/parent
  async registerParent(email, password, parentCode) {
    return await httpClient.post('/register/parent', {
      email,
      password,
      parent_code: parentCode,
    })
  },

  // Phụ huynh tự đăng ký bằng Google ID token và mã code
  // POST /api/v1/register/parent/google
  async registerParentWithGoogle(idToken, parentCode) {
    return await httpClient.post('/register/parent/google', {
      id_token: idToken,
      parent_code: parentCode,
    })
  },

  // Yêu cầu đặt lại mật khẩu (gửi email)
  // POST /api/v1/auth/forgot-password
  async forgotPassword(email) {
    return await httpClient.post('/auth/forgot-password', { email })
  },

  // Đặt lại mật khẩu bằng token
  // POST /api/v1/auth/reset-password
  async resetPassword(email, token, password) {
    return await httpClient.post('/auth/reset-password', { email, token, password })
  },

  async updateMyPassword(password) {
    return await httpClient.put('/me/password', { password })
  },
}
