/**
 * Teacher routes
 * Routes for TEACHER role
 */
export const teacherRoutes = {
  path: '/teacher',
  component: () => import('@/layouts/DashboardLayout.vue'),
  meta: { requiresAuth: true, roles: ['TEACHER'] },
  children: [
    {
      path: '',
      name: 'teacher-dashboard',
      component: () => import('@/views/teacher/TeacherDashboard.vue'),
    },
    {
      path: 'classes',
      name: 'teacher-classes',
      component: () => import('@/views/teacher/TeacherClasses.vue'),
    },
    {
      path: 'classes/:classId',
      name: 'teacher-class-detail',
      component: () => import('@/views/teacher/TeacherClassDetail.vue'),
    },
    {
      path: 'attendance',
      name: 'teacher-attendance',
      component: () => import('@/views/teacher/TeacherAttendance.vue'),
    },
    {
      path: 'health',
      name: 'teacher-health',
      component: () => import('@/views/teacher/TeacherHealth.vue'),
    },
    {
      path: 'posts',
      name: 'teacher-posts',
      component: () => import('@/views/teacher/TeacherPosts.vue'),
    },
    {
      path: 'appointments',
      name: 'teacher-appointments',
      component: () => import('@/views/teacher/TeacherAppointments.vue'),
    },
    {
      path: 'chat',
      name: 'teacher-chat',
      component: () => import('@/views/ChatPage.vue'),
    },
    {
      path: 'profile',
      name: 'teacher-profile',
      component: () => import('@/views/teacher/TeacherProfile.vue'),
    },
    // Legacy redirects
    { path: 'activities', redirect: '/teacher/posts' },
    { path: 'schedule', redirect: '/teacher/appointments' },
    { path: 'menu', redirect: '/teacher' },
  ],
}
