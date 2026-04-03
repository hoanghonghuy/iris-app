import axios from 'axios';
import { getApiBaseUrl } from '@/lib/runtime-config';

// base URL
const API_BASE_URL = getApiBaseUrl();

// axios instance
export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 30000, // 30 seconds
});

// Request Interceptor: chạy trước mỗi request, tự động lấy JWT token từ sessionStorage gắn vào Authorization header
apiClient.interceptors.request.use(
  (config) => {
    // Lấy token từ sessionStorage (chỉ chạy client-side)
    if (typeof window !== 'undefined') {
      const token = sessionStorage.getItem('auth_token');
      
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response Interceptor: chạy sau mỗi response, xử lý lỗi tự động
apiClient.interceptors.response.use(
  (response) => {
    // Request thành công, trả về response
    return response;
  },
  (error) => {
    // Xử lý lỗi
    if (error.response) {
      const status = error.response.status;
      
      switch (status) {
        case 401:
          // Unauthorized - Token hết hạn hoặc không hợp lệ
          if (typeof window !== 'undefined') {
            sessionStorage.removeItem('auth_token');
            localStorage.removeItem('user_role');
            
            // Chỉ redirect nếu không phải trang login
            if (!window.location.pathname.includes('/login')) {
              window.location.href = '/login';
            }
          }
          break;
          
        case 403:
          // Forbidden - Không có quyền
          console.error('Access forbidden:', error.response.data?.error);
          break;
          
        case 404:
          // Not Found
          console.error('Resource not found:', error.response.data?.error);
          break;
          
        case 500:
          // Server Error
          console.error('Server error:', error.response.data?.error);
          break;
          
        default:
          console.error(`HTTP ${status}:`, error.response.data?.error);
      }
    } else if (error.request) {
      // Request được gửi nhưng không nhận được response
      console.error('Network error: No response from server');
    } else {
      // lỗi khác
      console.error('Request error:', error.message);
    }
    
    return Promise.reject(error);
  }
);

// Helper functions
export const authHelpers = {
  setToken: (token: string) => {
    if (typeof window !== 'undefined') {
      sessionStorage.setItem('auth_token', token);
    }
  },
  
  getToken: (): string | null => {
    if (typeof window !== 'undefined') {
      return sessionStorage.getItem('auth_token');
    }
    return null;
  },
  
  removeToken: () => {
    if (typeof window !== 'undefined') {
      sessionStorage.removeItem('auth_token');
      localStorage.removeItem('user_role');
    }
  },
  
  setUserRole: (role: string) => {
    if (typeof window !== 'undefined') {
      localStorage.setItem('user_role', role);
    }
  },
  
  getUserRole: (): string | null => {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('user_role');
    }
    return null;
  },
  
  isAuthenticated: (): boolean => {
    return !!authHelpers.getToken(); // ép kiểu boolean
  },
};
