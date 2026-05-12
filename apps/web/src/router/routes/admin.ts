import type { RouteRecordRaw } from 'vue-router'

/**
 * Admin routes
 * Protected routes for SUPER_ADMIN and SCHOOL_ADMIN roles
 */
export const adminRoutes: RouteRecordRaw = {
  path: '/admin',
  component: () => import('@/layouts/DashboardLayout.vue'),
  meta: { requiresAuth: true, roles: ['SUPER_ADMIN', 'SCHOOL_ADMIN'] },
  children: [
    {
      path: '',
      name: 'admin-dashboard',
      component: () => import('@/views/admin/AdminDashboard.vue'),
    },
    {
      path: 'schools',
      name: 'admin-schools',
      component: () => import('@/views/admin/AdminSchools.vue'),
    },
    {
      path: 'classes',
      name: 'admin-classes',
      component: () => import('@/views/admin/AdminClasses.vue'),
    },
    {
      path: 'students',
      name: 'admin-students',
      component: () => import('@/views/admin/AdminStudents.vue'),
    },
    {
      path: 'students/:id',
      name: 'admin-student-detail',
      component: () => import('@/views/admin/AdminStudentDetail.vue'),
    },
    {
      path: 'users',
      name: 'admin-users',
      component: () => import('@/views/admin/AdminUsers.vue'),
    },
    {
      path: 'teachers',
      name: 'admin-teachers',
      component: () => import('@/views/admin/AdminTeachers.vue'),
    },
    {
      path: 'parents',
      name: 'admin-parents',
      component: () => import('@/views/admin/AdminParents.vue'),
    },
    {
      path: 'school-admins',
      name: 'admin-school-admins',
      component: () => import('@/views/admin/AdminSchoolAdmins.vue'),
      meta: { requiresAuth: true, roles: ['SUPER_ADMIN'] },
    },
    {
      path: 'audit-logs',
      name: 'admin-audit-logs',
      component: () => import('@/views/admin/AdminAuditLogs.vue'),
      meta: { requiresAuth: true, roles: ['SUPER_ADMIN'] },
    },
    {
      path: 'chat',
      name: 'admin-chat',
      component: () => import('@/views/ChatPage.vue'),
      meta: { requiresAuth: true, roles: ['SCHOOL_ADMIN'] },
    },
  ],
}
