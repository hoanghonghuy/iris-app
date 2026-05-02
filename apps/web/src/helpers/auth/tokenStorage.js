/**
 * Token Storage Helper
 * Centralized token and role management for authentication
 */

const TOKEN_KEY = 'auth_token'
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
    localStorage.removeItem(ROLE_KEY)
  },
}
