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
   * @returns {string|null} Token string or null if not found
   */
  getToken() {
    return sessionStorage.getItem(TOKEN_KEY)
  },

  /**
   * Set authentication token in session storage
   * @param {string|null} token - Token to store, or null to remove
   */
  setToken(token) {
    if (token) {
      sessionStorage.setItem(TOKEN_KEY, token)
    } else {
      sessionStorage.removeItem(TOKEN_KEY)
    }
  },

  /**
   * Get refresh token from session storage
   * @returns {string|null} Refresh token string or null if not found
   */
  getRefreshToken() {
    return sessionStorage.getItem(REFRESH_TOKEN_KEY)
  },

  /**
   * Set refresh token in session storage
   * @param {string|null} refreshToken - Refresh token to store, or null to remove
   */
  setRefreshToken(refreshToken) {
    if (refreshToken) {
      sessionStorage.setItem(REFRESH_TOKEN_KEY, refreshToken)
    } else {
      sessionStorage.removeItem(REFRESH_TOKEN_KEY)
    }
  },

  /**
   * Get user role from local storage
   * @returns {string|null} Role string or null if not found
   */
  getRole() {
    return localStorage.getItem(ROLE_KEY)
  },

  /**
   * Set user role in local storage
   * @param {string|null} role - Role to store, or null to remove
   */
  setRole(role) {
    if (role) {
      localStorage.setItem(ROLE_KEY, role)
    } else {
      localStorage.removeItem(ROLE_KEY)
    }
  },

  /**
   * Clear all authentication data
   */
  clear() {
    sessionStorage.removeItem(TOKEN_KEY)
    sessionStorage.removeItem(REFRESH_TOKEN_KEY)
    localStorage.removeItem(ROLE_KEY)
  },
}
