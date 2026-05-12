import type { RouteRecordRaw } from 'vue-router'

/**
 * Authentication routes
 * Public routes for login, register, password reset, etc.
 */
export const authRoutes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'landing',
    component: () => import('@/views/LandingPage.vue'),
    meta: { guestOnly: true },
  },
  {
    path: '/',
    component: () => import('@/layouts/AuthLayout.vue'),
    children: [
      {
        path: 'login',
        name: 'login',
        component: () => import('@/views/auth/LoginPage.vue'),
        meta: { guestOnly: true },
      },
      {
        path: 'register',
        name: 'register',
        component: () => import('@/views/auth/RegisterPage.vue'),
        meta: { guestOnly: true },
      },
      {
        path: 'forgot-password',
        name: 'forgot-password',
        component: () => import('@/views/auth/ForgotPasswordPage.vue'),
        meta: { guestOnly: true },
      },
      {
        path: 'reset-password',
        name: 'reset-password',
        component: () => import('@/views/auth/ResetPasswordPage.vue'),
        meta: { guestOnly: true },
      },
      {
        path: 'activate',
        name: 'activate',
        component: () => import('@/views/auth/ActivatePage.vue'),
        meta: { guestOnly: true },
      },
    ],
  },
]
