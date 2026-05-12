import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { navigationGuard } from './guards'
import { authRoutes } from './routes/auth'
import { adminRoutes } from './routes/admin'
import { teacherRoutes } from './routes/teacher'
import { parentRoutes } from './routes/parent'

/**
 * Application routes
 * Organized by role and functionality
 */
const routes: RouteRecordRaw[] = [
  ...authRoutes,
  adminRoutes,
  teacherRoutes,
  parentRoutes,
  // Catch-all redirect
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

/**
 * Vue Router instance
 */
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

/**
 * Global navigation guard
 * Handles authentication, role-based access, and guest-only routes
 */
router.beforeEach(navigationGuard)

export default router
