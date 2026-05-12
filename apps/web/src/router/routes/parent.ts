import type { RouteRecordRaw } from 'vue-router'

/**
 * Parent routes
 * Routes for PARENT role
 */
export const parentRoutes: RouteRecordRaw = {
  path: '/parent',
  component: () => import('@/layouts/DashboardLayout.vue'),
  meta: { requiresAuth: true, roles: ['PARENT'] },
  children: [
    {
      path: '',
      name: 'parent-dashboard',
      component: () => import('@/views/parent/ParentDashboard.vue'),
    },
    {
      path: 'children',
      name: 'parent-children',
      component: () => import('@/views/parent/ParentChildren.vue'),
    },
    {
      path: 'children/:studentId',
      name: 'parent-child-detail',
      component: () => import('@/views/parent/ParentChildDetail.vue'),
    },
    {
      path: 'feed',
      name: 'parent-feed',
      component: () => import('@/views/parent/ParentFeed.vue'),
    },
    {
      path: 'posts',
      name: 'parent-posts',
      component: () => import('@/views/parent/ParentFeed.vue'),
    },
    {
      path: 'appointments',
      name: 'parent-appointments',
      component: () => import('@/views/parent/ParentAppointments.vue'),
    },
    {
      path: 'chat',
      name: 'parent-chat',
      component: () => import('@/views/ChatPage.vue'),
    },
    {
      path: 'profile',
      name: 'parent-profile',
      component: () => import('@/views/parent/ParentProfile.vue'),
    },
    // Legacy redirects
    { path: 'newsfeed', redirect: { name: 'parent-feed' } },
  ],
}
