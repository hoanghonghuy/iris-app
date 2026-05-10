import { httpClient } from './httpClient'
import { teacherPostService } from './postService'
import { buildDateRangeParams } from '../helpers/queryParams'

export const teacherService = {
  async getAnalytics() {
    return await httpClient.get('/teacher/analytics')
  },

  async getAnalyticsTimeseries(params) {
    return await httpClient.get('/teacher/analytics/timeseries', params)
  },

  async getMyClasses() {
    return await httpClient.get('/teacher/classes')
  },

  async getStudentsInClass(classId) {
    return await httpClient.get(`/teacher/classes/${classId}/students`)
  },

  async markAttendance(data) {
    return await httpClient.post('/teacher/attendance', data)
  },

  async cancelAttendance(studentId, date) {
    return await httpClient.del(
      `/teacher/attendance?${new URLSearchParams({ student_id: studentId, date })}`,
    )
  },

  async getStudentAttendance(studentId, from, to) {
    return await httpClient.get(
      `/teacher/students/${studentId}/attendance`,
      buildDateRangeParams(from, to),
    )
  },

  async getStudentAttendanceChanges(studentId, from, to) {
    return await httpClient.get(
      `/teacher/students/${studentId}/attendance-changes`,
      buildDateRangeParams(from, to),
    )
  },

  async getClassAttendanceChanges(classId, params) {
    return await httpClient.get(`/teacher/classes/${classId}/attendance-changes`, params)
  },

  async createHealthLog(data) {
    return await httpClient.post('/teacher/health', data)
  },

  async getStudentHealth(studentId, from, to) {
    return await httpClient.get(
      `/teacher/students/${studentId}/health`,
      buildDateRangeParams(from, to),
    )
  },

  async updateMyProfile(phone) {
    return await httpClient.put('/teacher/profile', { phone })
  },

  async createPost(data) {
    return await httpClient.post('/teacher/posts', data)
  },

  async getClassPosts(classId, params) {
    return await httpClient.get(`/teacher/classes/${classId}/posts`, params)
  },

  async getStudentPosts(studentId, params) {
    return await httpClient.get(`/teacher/students/${studentId}/posts`, params)
  },

  // Post interactions delegated to postService
  ...teacherPostService,

  async getAppointments(params) {
    return await httpClient.get('/teacher/appointments', params)
  },

  async createAppointmentSlot(data) {
    return await httpClient.post('/teacher/appointments/slots', data)
  },

  async updateAppointmentStatus(appointmentId, status, cancelReason) {
    return await httpClient.patch(`/teacher/appointments/${appointmentId}/status`, {
      status,
      cancel_reason: cancelReason,
    })
  },
}
