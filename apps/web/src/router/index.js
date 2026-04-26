import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/authStore'
import { DASHBOARD_ROUTE_BY_ROLE } from '../helpers/authConfig'

const routes = [
  {
    path: '/',
    name: 'landing',
    component: () => import('../views/LandingPage.vue'),
    meta: { guestOnly: true },
  },
  {
    path: '/',
    component: () => import('../layouts/AuthLayout.vue'),
    children: [
      {
        path: 'login',
        name: 'login',
        component: () => import('../views/auth/LoginPage.vue'),
        meta: { guestOnly: true },
      },
      {
        path: 'register',
        name: 'register',
        component: () => import('../views/auth/RegisterPage.vue'),
        meta: { guestOnly: true },
      },
      {
        path: 'forgot-password',
        name: 'forgot-password',
        component: () => import('../views/auth/ForgotPasswordPage.vue'),
        meta: { guestOnly: true },
      },
      {
        path: 'reset-password',
        name: 'reset-password',
        component: () => import('../views/auth/ResetPasswordPage.vue'),
        meta: { guestOnly: true },
      },
      {
        path: 'activate',
        name: 'activate',
        component: () => import('../views/auth/ActivatePage.vue'),
        meta: { guestOnly: true },
      },
    ],
  },
  {
    path: '/admin',
    component: () => import('../layouts/DashboardLayout.vue'),
    meta: { requiresAuth: true, roles: ['SUPER_ADMIN', 'SCHOOL_ADMIN'] },
    children: [
      {
        path: '',
        name: 'admin-dashboard',
        component: () => import('../views/admin/AdminDashboard.vue'),
      },
      {
        path: 'schools',
        name: 'admin-schools',
        component: () => import('../views/admin/AdminSchools.vue'),
      },
      {
        path: 'classes',
        name: 'admin-classes',
        component: () => import('../views/admin/AdminClasses.vue'),
      },
      {
        path: 'students',
        name: 'admin-students',
        component: () => import('../views/admin/AdminStudents.vue'),
      },
      {
        path: 'students/:id',
        name: 'admin-student-detail',
        component: () => import('../views/admin/AdminStudentDetail.vue'),
      },
      {
        path: 'teachers',
        name: 'admin-teachers',
        component: () => import('../views/admin/AdminTeachers.vue'),
      },
      {
        path: 'parents',
        name: 'admin-parents',
        component: () => import('../views/admin/AdminParents.vue'),
      },
      {
        path: 'users',
        name: 'admin-users',
        component: () => import('../views/admin/AdminUsers.vue'),
      },
      {
        path: 'school-admins',
        name: 'admin-school-admins',
        component: () => import('../views/admin/AdminSchoolAdmins.vue'),
      },
      {
        path: 'audit-logs',
        name: 'admin-audit-logs',
        component: () => import('../views/admin/AdminAuditLogs.vue'),
      },
      { path: 'chat', name: 'admin-chat', component: () => import('../views/ChatPage.vue') },
    ],
  },
  {
    path: '/teacher',
    component: () => import('../layouts/DashboardLayout.vue'),
    meta: { requiresAuth: true, roles: ['TEACHER'] },
    children: [
      {
        path: '',
        name: 'teacher-dashboard',
        component: () => import('../views/teacher/TeacherDashboard.vue'),
      },
      {
        path: 'classes',
        name: 'teacher-classes',
        component: () => import('../views/teacher/TeacherClasses.vue'),
      },
      {
        path: 'classes/:classId',
        name: 'teacher-class-detail',
        component: () => import('../views/teacher/TeacherClassDetail.vue'),
      },
      {
        path: 'attendance',
        name: 'teacher-attendance',
        component: () => import('../views/teacher/TeacherAttendance.vue'),
      },
      {
        path: 'health',
        name: 'teacher-health',
        component: () => import('../views/teacher/TeacherHealth.vue'),
      },
      {
        path: 'posts',
        name: 'teacher-posts',
        component: () => import('../views/teacher/TeacherPosts.vue'),
      },
      {
        path: 'appointments',
        name: 'teacher-appointments',
        component: () => import('../views/teacher/TeacherAppointments.vue'),
      },
      { path: 'chat', name: 'teacher-chat', component: () => import('../views/ChatPage.vue') },
      {
        path: 'profile',
        name: 'teacher-profile',
        component: () => import('../views/teacher/TeacherProfile.vue'),
      },
      { path: 'activities', redirect: '/teacher/posts' },
      { path: 'schedule', redirect: '/teacher/appointments' },
      { path: 'menu', redirect: '/teacher' },
    ],
  },
  {
    path: '/parent',
    component: () => import('../layouts/DashboardLayout.vue'),
    meta: { requiresAuth: true, roles: ['PARENT'] },
    children: [
      {
        path: '',
        name: 'parent-dashboard',
        component: () => import('../views/parent/ParentDashboard.vue'),
      },
      {
        path: 'children',
        name: 'parent-children',
        component: () => import('../views/parent/ParentChildren.vue'),
      },
      {
        path: 'children/:studentId',
        name: 'parent-child-detail',
        component: () => import('../views/parent/ParentChildDetail.vue'),
      },
      {
        path: 'feed',
        name: 'parent-feed',
        component: () => import('../views/parent/ParentFeed.vue'),
      },
      {
        path: 'posts',
        name: 'parent-posts',
        component: () => import('../views/parent/ParentFeed.vue'),
      },
      {
        path: 'appointments',
        name: 'parent-appointments',
        component: () => import('../views/parent/ParentAppointments.vue'),
      },
      { path: 'chat', name: 'parent-chat', component: () => import('../views/ChatPage.vue') },
      {
        path: 'profile',
        name: 'parent-profile',
        component: () => import('../views/parent/ParentProfile.vue'),
      },
      { path: 'notifications', redirect: '/parent/feed' },
    ],
  },
  { path: '/:pathMatch(.*)*', redirect: '/' },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach(async (to) => {
  const authStore = useAuthStore()

  if (authStore.isAuthenticated && !authStore.currentUserRole) {
    await authStore.fetchCurrentUser()
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  if (to.meta.guestOnly && authStore.isAuthenticated) {
    return DASHBOARD_ROUTE_BY_ROLE[authStore.currentUserRole] || '/'
  }

  if (to.meta.roles?.length) {
    const role = authStore.currentUserRole
    if (!role || !to.meta.roles.includes(role)) {
      return DASHBOARD_ROUTE_BY_ROLE[role] || { name: 'login' }
    }
  }

  return true
})

export default router
