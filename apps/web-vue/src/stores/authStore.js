import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { authService } from '../services/authService'

export const useAuthStore = defineStore('auth', () => {
  // State
  const token = ref(sessionStorage.getItem('auth_token') || null)
  const currentUser = ref(null)
  const currentUserRole = ref(localStorage.getItem('user_role') || null)
  const isLoading = ref(false)
  const errorMessage = ref(null)

  // Getters
  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => currentUserRole.value === 'SUPER_ADMIN' || currentUserRole.value === 'SCHOOL_ADMIN')
  const isTeacher = computed(() => currentUserRole.value === 'TEACHER')
  const isParent = computed(() => currentUserRole.value === 'PARENT')

  // Actions
  function setToken(newToken) {
    token.value = newToken
    if (newToken) {
      sessionStorage.setItem('auth_token', newToken)
    } else {
      sessionStorage.removeItem('auth_token')
    }
  }

  function setRole(role) {
    currentUserRole.value = role
    if (role) {
      localStorage.setItem('user_role', role)
    } else {
      localStorage.removeItem('user_role')
    }
  }

  function setUser(user) {
    currentUser.value = user
  }

  function clearAuth() {
    setToken(null)
    setRole(null)
    setUser(null)
  }

  async function fetchCurrentUser() {
    if (!token.value) return null

    isLoading.value = true
    errorMessage.value = null
    try {
      const user = await authService.getMe()
      setUser(user)
      // BE trả roles là array: ['SUPER_ADMIN'], lấy phần tử đầu tiên
      if (user && Array.isArray(user.roles) && user.roles.length > 0) {
        setRole(user.roles[0])
      } else if (user && user.role) {
        setRole(user.role)
      }
      return user
    } catch (error) {
      // Nếu lỗi 401, httpClient đã tự xử lý redirect
      // Ở đây chỉ clear state nếu fetch lỗi
      console.error('Failed to fetch user:', error)
      clearAuth()
      return null
    } finally {
      isLoading.value = false
    }
  }

  return {
    // State
    token,
    currentUser,
    currentUserRole,
    isLoading,
    errorMessage,
    
    // Getters
    isAuthenticated,
    isAdmin,
    isTeacher,
    isParent,
    
    // Actions
    setToken,
    setRole,
    setUser,
    clearAuth,
    fetchCurrentUser
  }
})
