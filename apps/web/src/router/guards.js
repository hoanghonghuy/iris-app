import { useAuthStore } from '@/stores/authStore'
import { DASHBOARD_ROUTE_BY_ROLE } from '@/helpers/authConfig'

/**
 * Authentication guard
 * Ensures user is authenticated before accessing protected routes
 */
export async function authGuard(to) {
  const authStore = useAuthStore()

  // Fetch user data if authenticated but role not loaded
  if (authStore.isAuthenticated && !authStore.currentUserRole) {
    await authStore.fetchCurrentUser()
  }

  // Redirect to login if auth required but not authenticated
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  return true
}

/**
 * Guest-only guard
 * Redirects authenticated users away from guest-only pages (login, register)
 */
export function guestOnlyGuard(to) {
  const authStore = useAuthStore()

  if (to.meta.guestOnly && authStore.isAuthenticated) {
    return DASHBOARD_ROUTE_BY_ROLE[authStore.currentUserRole] || '/'
  }

  return true
}

/**
 * Role-based access control guard
 * Ensures user has required role to access route
 */
export function roleGuard(to) {
  const authStore = useAuthStore()

  if (to.meta.roles?.length) {
    const role = authStore.currentUserRole
    if (!role || !to.meta.roles.includes(role)) {
      return DASHBOARD_ROUTE_BY_ROLE[role] || { name: 'login' }
    }
  }

  return true
}

/**
 * Combined navigation guard
 * Runs all guards in sequence
 */
export async function navigationGuard(to) {
  // Run guards in order
  const authResult = await authGuard(to)
  if (authResult !== true) return authResult

  const guestResult = guestOnlyGuard(to)
  if (guestResult !== true) return guestResult

  const roleResult = roleGuard(to)
  if (roleResult !== true) return roleResult

  return true
}
