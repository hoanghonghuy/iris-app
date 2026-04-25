import { httpClient } from './httpClient'

export const adminService = {
  // --- ANALYTICS ---
  async getAnalytics() {
    return await httpClient.get('/admin/analytics')
  },

  async getAuditLogs(params) {
    return await httpClient.get('/admin/audit-logs', params)
  },

  // --- SCHOOLS ---
  async getSchools(params) {
    return await httpClient.get('/admin/schools', params)
  },
  async createSchool(data) {
    return await httpClient.post('/admin/schools', data)
  },
  async updateSchool(schoolId, data) {
    return await httpClient.put(`/admin/schools/${schoolId}`, data)
  },
  async deleteSchool(schoolId) {
    return await httpClient.del(`/admin/schools/${schoolId}`)
  },

  // --- CLASSES ---
  async getClassesBySchool(schoolId, params) {
    return await httpClient.get(`/admin/classes/by-school/${schoolId}`, params)
  },
  async createClass(data) {
    return await httpClient.post('/admin/classes', data)
  },
  async updateClass(classId, data) {
    return await httpClient.put(`/admin/classes/${classId}`, data)
  },
  async deleteClass(classId) {
    return await httpClient.del(`/admin/classes/${classId}`)
  },

  // --- STUDENTS ---
  async getStudentsByClass(classId, params) {
    return await httpClient.get(`/admin/students/by-class/${classId}`, params)
  },
  async getStudentProfile(studentId) {
    return await httpClient.get(`/admin/students/${studentId}`)
  },
  async createStudent(data) {
    return await httpClient.post('/admin/students', data)
  },
  async updateStudent(studentId, data) {
    return await httpClient.put(`/admin/students/${studentId}`, data)
  },
  async deleteStudent(studentId) {
    return await httpClient.del(`/admin/students/${studentId}`)
  },
  async generateParentCode(studentId) {
    return await httpClient.post(`/admin/students/${studentId}/generate-parent-code`)
  },
  async revokeParentCode(studentId) {
    return await httpClient.del(`/admin/students/${studentId}/parent-code`)
  },

  // --- USERS ---
  async getUsers(params) {
    return await httpClient.get('/admin/users', params)
  },
  async createUser(data) {
    return await httpClient.post('/admin/users', data)
  },
  async getUserById(userId) {
    return await httpClient.get(`/admin/users/${userId}`)
  },
  async lockUser(userId) {
    return await httpClient.post(`/admin/users/${userId}/lock`)
  },
  async unlockUser(userId) {
    return await httpClient.post(`/admin/users/${userId}/unlock`)
  },
  async assignRole(userId, roleName) {
    return await httpClient.post(`/admin/users/${userId}/roles`, { role_name: roleName })
  },

  // --- TEACHERS ---
  async getTeachers(params) {
    return await httpClient.get('/admin/teachers', params)
  },
  async getTeacherById(teacherId) {
    return await httpClient.get(`/admin/teachers/${teacherId}`)
  },
  async updateTeacher(teacherId, data) {
    return await httpClient.put(`/admin/teachers/${teacherId}`, data)
  },
  async getTeachersOfClass(classId) {
    return await httpClient.get(`/admin/teachers/class/${classId}`)
  },
  async assignTeacherToClass(teacherId, classId) {
    return await httpClient.post(`/admin/teachers/${teacherId}/classes/${classId}`)
  },
  async unassignTeacherFromClass(teacherId, classId) {
    return await httpClient.del(`/admin/teachers/${teacherId}/classes/${classId}`)
  },
  async deleteTeacher(teacherId) {
    return await httpClient.del(`/admin/teachers/${teacherId}`)
  },

  // --- PARENTS ---
  async getParents(params) {
    return await httpClient.get('/admin/parents', params)
  },
  async getParentById(parentId) {
    return await httpClient.get(`/admin/parents/${parentId}`)
  },
  async updateParent(parentId, data) {
    return await httpClient.put(`/admin/parents/${parentId}`, data)
  },
  async assignParentToStudent(parentId, studentId) {
    return await httpClient.post(`/admin/parents/${parentId}/students/${studentId}`)
  },
  async unassignParentFromStudent(parentId, studentId) {
    return await httpClient.del(`/admin/parents/${parentId}/students/${studentId}`)
  },

  // --- SCHOOL ADMINS ---
  async getSchoolAdmins(params) {
    return await httpClient.get('/admin/school-admins', params)
  },
  async createSchoolAdmin(data) {
    return await httpClient.post('/admin/school-admins', data)
  },
  async deleteSchoolAdmin(adminId) {
    return await httpClient.del(`/admin/school-admins/${adminId}`)
  },
}
