import { defineStore } from 'pinia'
import { authService } from '../services/authService'
import { tokenStorage } from '@/helpers/auth'
import type { User, UserRole } from '@/types'

interface AuthState {
  token: string | null
  refreshToken: string | null
  currentUser: User | null
  currentUserRole: UserRole | null
  isLoading: boolean
  errorMessage: string | null
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: tokenStorage.getToken() || null,
    refreshToken: tokenStorage.getRefreshToken() || null,
    currentUser: null,
    currentUserRole: (tokenStorage.getRole() as UserRole) || null,
    isLoading: false,
    errorMessage: null,
  }),

  getters: {
    isAuthenticated: (state): boolean => !!state.token,
    isAdmin: (state): boolean =>
      state.currentUserRole === 'SUPER_ADMIN' || state.currentUserRole === 'SCHOOL_ADMIN',
    isTeacher: (state): boolean => state.currentUserRole === 'TEACHER',
    isParent: (state): boolean => state.currentUserRole === 'PARENT',
  },

  actions: {
    setToken(newToken: string | null): void {
      this.token = newToken
      tokenStorage.setToken(newToken)
    },

    setRefreshToken(newRefreshToken: string | null): void {
      this.refreshToken = newRefreshToken
      tokenStorage.setRefreshToken(newRefreshToken)
    },

    setRole(role: UserRole | null): void {
      this.currentUserRole = role
      tokenStorage.setRole(role)
    },

    setUser(user: User | null): void {
      this.currentUser = user
    },

    clearAuth(): void {
      this.token = null
      this.refreshToken = null
      this.currentUserRole = null
      this.currentUser = null
      tokenStorage.clear()
    },

    async fetchCurrentUser(): Promise<User | null> {
      if (!this.token) return null

      this.isLoading = true
      this.errorMessage = null
      try {
        const user = await authService.getMe()
        this.setUser(user)

        // Single role per user - lấy role đầu tiên (và duy nhất)
        if (user && Array.isArray(user.roles) && user.roles.length > 0) {
          this.setRole(user.roles[0] as UserRole)
        } else if (user && (user as any).role) {
          this.setRole((user as any).role as UserRole)
        }
        return user
      } catch (error) {
        // Nếu lỗi 401, httpClient đã tự xử lý redirect
        // Ở đây chỉ clear state nếu fetch lỗi
        console.error('Failed to fetch user:', error)
        this.clearAuth()
        return null
      } finally {
        this.isLoading = false
      }
    },
  },
})
