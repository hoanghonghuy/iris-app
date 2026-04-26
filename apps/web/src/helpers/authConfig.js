export const DASHBOARD_ROUTE_BY_ROLE = {
  SUPER_ADMIN: '/admin',
  SCHOOL_ADMIN: '/admin',
  TEACHER: '/teacher',
  PARENT: '/parent',
}

export const PROFILE_ROUTE_BY_ROLE = {
  SUPER_ADMIN: null,
  SCHOOL_ADMIN: null,
  TEACHER: '/teacher/profile',
  PARENT: '/parent/profile',
}

export const ROLE_LABELS = {
  SUPER_ADMIN: 'Quản trị viên cấp cao',
  SCHOOL_ADMIN: 'Quản trị viên trường',
  TEACHER: 'Giáo viên',
  PARENT: 'Phụ huynh',
}

export const adminMenuItems = [
  { label: 'Tổng quan', path: '/admin', icon: 'dashboard' },
  { label: 'Trường học', path: '/admin/schools', icon: 'school' },
  { label: 'Lớp học', path: '/admin/classes', icon: 'class' },
  { label: 'Học sinh', path: '/admin/students', icon: 'students' },
  { label: 'Người dùng', path: '/admin/users', icon: 'users' },
  { label: 'Giáo viên', path: '/admin/teachers', icon: 'teacher' },
  { label: 'Phụ huynh', path: '/admin/parents', icon: 'parent' },
  { label: 'School Admin', path: '/admin/school-admins', icon: 'shield', roles: ['SUPER_ADMIN'] },
  { label: 'Audit Logs', path: '/admin/audit-logs', icon: 'logs', roles: ['SUPER_ADMIN'] },
  { label: 'Tin nhắn', path: '/admin/chat', icon: 'message', roles: ['SCHOOL_ADMIN'] },
]

export const teacherMenuItems = [
  { label: 'Tổng quan', path: '/teacher', icon: 'dashboard' },
  { label: 'Lớp của tôi', path: '/teacher/classes', icon: 'class' },
  { label: 'Điểm danh', path: '/teacher/attendance', icon: 'attendance' },
  { label: 'Sức khỏe', path: '/teacher/health', icon: 'health' },
  { label: 'Bài đăng', path: '/teacher/posts', icon: 'post' },
  { label: 'Lịch hẹn', path: '/teacher/appointments', icon: 'calendar' },
  { label: 'Tin nhắn', path: '/teacher/chat', icon: 'message' },
]

export const parentMenuItems = [
  { label: 'Tổng quan', path: '/parent', icon: 'dashboard' },
  { label: 'Con của tôi', path: '/parent/children', icon: 'child' },
  { label: 'Bảng tin', path: '/parent/feed', icon: 'feed' },
  { label: 'Lịch hẹn', path: '/parent/appointments', icon: 'calendar' },
  { label: 'Tin nhắn', path: '/parent/chat', icon: 'message' },
]
