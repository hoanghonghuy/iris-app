/**
 * Token Storage Helper
 * Centralized token and role management for authentication
 */

const TOKEN_KEY = 'auth_token'
const REFRESH_TOKEN_KEY = 'refresh_token'
const ROLE_KEY = 'user_role'

export const tokenStorage = {
  /**
   * Get authentication token from session storage
   */
  getToken(): string | null {
    return sessionStorage.getItem(TOKEN_KEY)
  },

  /**
   * Set authentication token in session storage
   */
  setToken(token: string | null): void {
    if (token) {
      sessionStorage.setItem(TOKEN_KEY, token)
    } else {
      sessionStorage.removeItem(TOKEN_KEY)
    }
  },

  /**
   * Get refresh token from session storage
   */
  getRefreshToken(): string | null {
    return sessionStorage.getItem(REFRESH_TOKEN_KEY)
  },

  /**
   * Set refresh token in session storage
   */
  setRefreshToken(refreshToken: string | null): void {
    if (refreshToken) {
      sessionStorage.setItem(REFRESH_TOKEN_KEY, refreshToken)
    } else {
      sessionStorage.removeItem(REFRESH_TOKEN_KEY)
    }
  },

  /**
   * Get user role from local storage
   */
  getRole(): string | null {
    return localStorage.getItem(ROLE_KEY)
  },

  /**
   * Set user role in local storage
   */
  setRole(role: string | null): void {
    if (role) {
      localStorage.setItem(ROLE_KEY, role)
    } else {
      localStorage.removeItem(ROLE_KEY)
    }
  },

  /**
   * Clear all authentication data
   */
  clear(): void {
    sessionStorage.removeItem(TOKEN_KEY)
    sessionStorage.removeItem(REFRESH_TOKEN_KEY)
    localStorage.removeItem(ROLE_KEY)
  },
}
