// API response envelope types
export interface ApiResponse<T = any> {
  data?: T
  message?: string
  error?: string
  errors?: Record<string, string[]>
}

export interface ApiError {
  status: number
  message: string
  name?: string
  data?: any
}

// Pagination types
export interface PaginationParams {
  page?: number
  limit?: number
  sort?: string
  order?: 'asc' | 'desc'
}

export interface PaginatedResponse<T> {
  data: T[]
  pagination: {
    page: number
    limit: number
    total: number
    total_pages: number
  }
}

// Common query params
export interface SearchParams {
  q?: string
  school_id?: string
  class_id?: string
  student_id?: string
  teacher_id?: string
  parent_id?: string
  status?: string
  from_date?: string
  to_date?: string
}

// HTTP client types
export interface HttpClientConfig {
  baseURL: string
  timeout: number
  headers?: Record<string, string>
}

export interface RequestConfig {
  headers?: Record<string, string>
  params?: Record<string, any>
  timeout?: number
}
