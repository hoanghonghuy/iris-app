// src/services/httpClient.js
import { tokenStorage } from '@/helpers/auth'

// base URL mặc định, đọc từ biến môi trường
const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

// thời gian chờ tối đa cho mỗi request (30 giây)
const REQUEST_TIMEOUT = 30000

// flag để tránh refresh loop
let isRefreshing = false
let refreshPromise = null

// hàm gọi API chung — tất cả service đều dùng hàm này
async function request(method, path, body = null, options = {}) {
  const token = tokenStorage.getToken()
  const { signal: externalSignal, timeout = REQUEST_TIMEOUT, skipRefresh = false } = options

  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  }

  // tự gắn token nếu có
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  // tạo AbortController để hủy request khi quá thời gian chờ
  const controller = new AbortController()
  const timeoutId = setTimeout(() => controller.abort(), timeout)

  const handleExternalAbort = () => controller.abort()
  if (externalSignal) {
    if (externalSignal.aborted) {
      controller.abort()
    } else {
      externalSignal.addEventListener('abort', handleExternalAbort, { once: true })
    }
  }

  const fetchOptions = {
    method,
    headers,
    signal: controller.signal,
  }

  // chỉ gắn body khi có dữ liệu gửi lên (POST, PUT, PATCH)
  if (body !== null) {
    fetchOptions.body = JSON.stringify(body)
  }

  let response
  try {
    response = await fetch(`${BASE_URL}${path}`, fetchOptions)
  } catch (err) {
    // xử lý lỗi mạng hoặc timeout
    if (err.name === 'AbortError') {
      const isCancelledByCaller = Boolean(externalSignal?.aborted)
      const error = new Error(
        isCancelledByCaller ? 'Yêu cầu đã bị hủy' : 'Yêu cầu quá thời gian, vui lòng thử lại',
      )
      error.name = 'AbortError'
      error.status = isCancelledByCaller ? 499 : 408
      throw error
    }
    const error = new Error('Không thể kết nối đến máy chủ')
    error.status = 0
    throw error
  } finally {
    clearTimeout(timeoutId)
    if (externalSignal) {
      externalSignal.removeEventListener('abort', handleExternalAbort)
    }
  }

  // xử lý 401 — token hết hạn → thử refresh token
  if (response.status === 401 && !skipRefresh) {
    const refreshToken = tokenStorage.getRefreshToken()
    
    // Nếu có refresh token, thử làm mới access token
    if (refreshToken && path !== '/auth/refresh') {
      try {
        // Nếu đang refresh, chờ promise hiện tại
        if (isRefreshing && refreshPromise) {
          await refreshPromise
        } else {
          // Bắt đầu refresh mới
          isRefreshing = true
          refreshPromise = refreshAccessToken(refreshToken)
          await refreshPromise
        }

        // Sau khi refresh thành công, retry request ban đầu với token mới
        return await request(method, path, body, { ...options, skipRefresh: true })
      } catch (refreshError) {
        // Refresh thất bại → clear tokens và redirect login
        tokenStorage.clear()
        if (!window.location.pathname.includes('/login')) {
          window.location.href = '/login'
        }
        throw refreshError
      } finally {
        isRefreshing = false
        refreshPromise = null
      }
    } else {
      // Không có refresh token hoặc đang gọi /auth/refresh → logout
      tokenStorage.clear()
      if (!window.location.pathname.includes('/login')) {
        window.location.href = '/login'
      }
    }
  }

  // parse body JSON (có thể rỗng với 204 No Content)
  let data = null
  const contentType = response.headers.get('content-type')
  if (contentType && contentType.includes('application/json')) {
    data = await response.json()
  }

  // nếu HTTP lỗi (4xx, 5xx) → throw để service/component bắt được
  if (!response.ok) {
    const error = new Error(data?.error || `HTTP ${response.status}`)
    error.status = response.status
    error.data = data
    throw error
  }

  return data
}

// Hàm refresh access token
async function refreshAccessToken(refreshToken) {
  const response = await fetch(`${BASE_URL}/auth/refresh`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ refresh_token: refreshToken }),
  })

  if (!response.ok) {
    throw new Error('Refresh token failed')
  }

  const data = await response.json()
  const newAccessToken = data?.data?.access_token || data?.access_token
  const newRefreshToken = data?.data?.refresh_token || data?.refresh_token

  if (!newAccessToken || !newRefreshToken) {
    throw new Error('Invalid refresh response')
  }

  // Lưu tokens mới
  tokenStorage.setToken(newAccessToken)
  tokenStorage.setRefreshToken(newRefreshToken)

  return data
}

// các hàm tiện ích — service chỉ cần gọi httpClient.get('/path')
export const httpClient = {
  get: (path, params, options = {}) => {
    // tự nối query string từ object params
    if (params) {
      const query = new URLSearchParams(
        Object.fromEntries(Object.entries(params).filter(([, value]) => value != null)),
      ).toString()
      if (query) path = `${path}?${query}`
    }
    return request('GET', path, null, options)
  },
  post: (path, body, options = {}) => request('POST', path, body, options),
  put: (path, body, options = {}) => request('PUT', path, body, options),
  patch: (path, body, options = {}) => request('PATCH', path, body, options),
  del: (path, options = {}) => request('DELETE', path, null, options),
}
