/**
 * Admin API Service
 * Quản lý các endpoint dành cho Admin (SUPER_ADMIN & SCHOOL_ADMIN)
 */
import { apiClient } from './client';
import {
  School, Class, Student, UserInfo, Teacher, Parent,
  ApiResponse, PaginationParams, CreateSchoolRequest, CreateClassRequest, CreateStudentRequest, AdminAnalytics
} from '@/types';

export const adminApi = {
  // --- ANALYTICS ---
  getAnalytics: async () => {
    const res = await apiClient.get<ApiResponse<AdminAnalytics>>('/admin/analytics');
    return res.data.data;
  },

  // --- SCHOOLS ---
  getSchools: async (params?: PaginationParams) => {
    const res = await apiClient.get<ApiResponse<School[]>>('/admin/schools', { params });
    return res.data;
  },
  createSchool: async (data: CreateSchoolRequest) => {
    const res = await apiClient.post('/admin/schools', data);
    return res.data;
  },

  // --- CLASSES ---
  getClassesBySchool: async (schoolId: string, params?: PaginationParams) => {
    const res = await apiClient.get<ApiResponse<Class[]>>(`/admin/classes/by-school/${schoolId}`, { params });
    return res.data;
  },
  createClass: async (data: CreateClassRequest) => {
    const res = await apiClient.post('/admin/classes', data);
    return res.data;
  },

  // --- STUDENTS ---
  getStudentsByClass: async (classId: string, params?: PaginationParams) => {
    const res = await apiClient.get<ApiResponse<Student[]>>(`/admin/students/by-class/${classId}`, { params });
    return res.data;
  },
  createStudent: async (data: CreateStudentRequest) => {
    const res = await apiClient.post('/admin/students', data);
    return res.data;
  },
  generateParentCode: async (studentId: string) => {
    const res = await apiClient.post(`/admin/students/${studentId}/generate-parent-code`);
    return res.data;
  },
  revokeParentCode: async (studentId: string) => {
    const res = await apiClient.delete(`/admin/students/${studentId}/parent-code`);
    return res.data;
  },

  // --- USERS ---
  getUsers: async (params?: PaginationParams) => {
    const res = await apiClient.get<ApiResponse<UserInfo[]>>('/admin/users', { params });
    return res.data;
  },
  createUser: async (data: { email: string; roles: string[] }) => {
    const res = await apiClient.post('/admin/users', data);
    return res.data;
  },
  getUserById: async (userId: string) => {
    const res = await apiClient.get<ApiResponse<UserInfo>>(`/admin/users/${userId}`);
    return res.data.data;
  },
  lockUser: async (userId: string) => {
    const res = await apiClient.post(`/admin/users/${userId}/lock`);
    return res.data;
  },
  unlockUser: async (userId: string) => {
    const res = await apiClient.post(`/admin/users/${userId}/unlock`);
    return res.data;
  },
  assignRole: async (userId: string, roleName: string) => {
    const res = await apiClient.post(`/admin/users/${userId}/roles`, { role_name: roleName });
    return res.data;
  },

  // --- TEACHERS ---
  getTeachers: async (params?: PaginationParams) => {
    const res = await apiClient.get<ApiResponse<Teacher[]>>('/admin/teachers', { params });
    return res.data;
  },
  getTeacherById: async (teacherId: string) => {
    const res = await apiClient.get<ApiResponse<Teacher>>(`/admin/teachers/${teacherId}`);
    return res.data.data;
  },
  updateTeacher: async (teacherId: string, data: { full_name?: string; phone?: string }) => {
    const res = await apiClient.put(`/admin/teachers/${teacherId}`, data);
    return res.data;
  },
  getTeachersOfClass: async (classId: string) => {
    const res = await apiClient.get<ApiResponse<Teacher[]>>(`/admin/teachers/class/${classId}`);
    return res.data.data;
  },
  assignTeacherToClass: async (teacherId: string, classId: string) => {
    const res = await apiClient.post(`/admin/teachers/${teacherId}/classes/${classId}`);
    return res.data;
  },
  unassignTeacherFromClass: async (teacherId: string, classId: string) => {
    const res = await apiClient.delete(`/admin/teachers/${teacherId}/classes/${classId}`);
    return res.data;
  },

  // --- PARENTS ---
  getParents: async (params?: PaginationParams) => {
    const res = await apiClient.get<ApiResponse<Parent[]>>('/admin/parents', { params });
    return res.data;
  },
  getParentById: async (parentId: string) => {
    const res = await apiClient.get<ApiResponse<Parent>>(`/admin/parents/${parentId}`);
    return res.data.data;
  },
  assignParentToStudent: async (parentId: string, studentId: string) => {
    const res = await apiClient.post(`/admin/parents/${parentId}/students/${studentId}`);
    return res.data;
  },
  unassignParentFromStudent: async (parentId: string, studentId: string) => {
    const res = await apiClient.delete(`/admin/parents/${parentId}/students/${studentId}`);
    return res.data;
  },

  // --- SCHOOL ADMINS ---
  getSchoolAdmins: async (params?: PaginationParams) => {
    const res = await apiClient.get<ApiResponse<any[]>>('/admin/school-admins', { params });
    return res.data;
  },
  createSchoolAdmin: async (data: { user_id: string; school_id: string }) => {
    const res = await apiClient.post('/admin/school-admins', data);
    return res.data;
  },
  deleteSchoolAdmin: async (adminId: string) => {
    const res = await apiClient.delete(`/admin/school-admins/${adminId}`);
    return res.data;
  },
};