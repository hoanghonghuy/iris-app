// src/services/httpClient.js

// base URL mặc định, đọc từ biến môi trường
const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

// thời gian chờ tối đa cho mỗi request (30 giây)
const REQUEST_TIMEOUT = 30000

// hàm gọi API chung — tất cả service đều dùng hàm này
async function request(method, path, body = null, options = {}) {
  const token = sessionStorage.getItem('auth_token')

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
  const timeoutId = setTimeout(() => controller.abort(), REQUEST_TIMEOUT)

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
      const error = new Error('Yêu cầu quá thời gian, vui lòng thử lại')
      error.status = 408
      throw error
    }
    const error = new Error('Không thể kết nối đến máy chủ')
    error.status = 0
    throw error
  } finally {
    clearTimeout(timeoutId)
  }

  // xử lý 401 — token hết hạn → về trang login
  if (response.status === 401) {
    sessionStorage.removeItem('auth_token')
    localStorage.removeItem('user_role')
    if (!window.location.pathname.includes('/login')) {
      window.location.href = '/login'
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

// các hàm tiện ích — service chỉ cần gọi httpClient.get('/path')
export const httpClient = {
  get: (path, params) => {
    // tự nối query string từ object params
    if (params) {
      const query = new URLSearchParams(
        Object.fromEntries(
          Object.entries(params).filter(([, value]) => value != null)
        )
      ).toString()
      if (query) path = `${path}?${query}`
    }
    return request('GET', path)
  },
  post: (path, body) => request('POST', path, body),
  put: (path, body) => request('PUT', path, body),
  patch: (path, body) => request('PATCH', path, body),
  del: (path) => request('DELETE', path),
}
