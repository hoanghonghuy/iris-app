import { defineStore } from 'pinia'
import { authService } from '../services/authService'
import { tokenStorage } from '@/helpers/auth'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: tokenStorage.getToken() || null,
    currentUser: null,
    currentUserRole: tokenStorage.getRole() || null,
    isLoading: false,
    errorMessage: null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) =>
      state.currentUserRole === 'SUPER_ADMIN' || state.currentUserRole === 'SCHOOL_ADMIN',
    isTeacher: (state) => state.currentUserRole === 'TEACHER',
    isParent: (state) => state.currentUserRole === 'PARENT',
  },

  actions: {
    setToken(newToken) {
      this.token = newToken
      tokenStorage.setToken(newToken)
    },

    setRole(role) {
      this.currentUserRole = role
      tokenStorage.setRole(role)
    },

    setUser(user) {
      this.currentUser = user
    },

    clearAuth() {
      this.token = null
      this.currentUserRole = null
      this.currentUser = null
      tokenStorage.clear()
    },

    async fetchCurrentUser() {
      if (!this.token) return null

      this.isLoading = true
      this.errorMessage = null
      try {
        const user = await authService.getMe()
        this.setUser(user)

        // Single role per user - lấy role đầu tiên (và duy nhất)
        if (user && Array.isArray(user.roles) && user.roles.length > 0) {
          this.setRole(user.roles[0])
        } else if (user && user.role) {
          this.setRole(user.role)
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
