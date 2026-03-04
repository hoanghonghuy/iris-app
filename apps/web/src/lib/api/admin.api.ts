/**
 * Admin API Service
 * Quản lý các endpoint dành cho Admin (SUPER_ADMIN & SCHOOL_ADMIN)
 * Tương ứng với: apps/api/internal/api/v1/handlers/ (admin_helpers.go, school_handler.go, class_handler.go...)
 */
import { apiClient } from './client';
import { 
  School, Class, Student, UserInfo, Teacher, Parent, 
  ApiResponse, CreateSchoolRequest, CreateClassRequest, CreateStudentRequest 
} from '@/types';

export const adminApi = {
  // --- SCHOOLS ---
  getSchools: async () => {
    const res = await apiClient.get<ApiResponse<School[]>>('/admin/schools');
    return res.data.data;
  },
  createSchool: async (data: CreateSchoolRequest) => {
    const res = await apiClient.post('/admin/schools', data);
    return res.data;
  },

  // --- CLASSES ---
  getClassesBySchool: async (schoolId: string) => {
    const res = await apiClient.get<ApiResponse<Class[]>>(`/admin/classes/by-school/${schoolId}`);
    return res.data.data;
  },
  createClass: async (data: CreateClassRequest) => {
    const res = await apiClient.post('/admin/classes', data);
    return res.data;
  },

  // --- STUDENTS ---
  getStudentsByClass: async (classId: string) => {
    const res = await apiClient.get<ApiResponse<Student[]>>(`/admin/students/by-class/${classId}`);
    return res.data.data;
  },
  createStudent: async (data: CreateStudentRequest) => {
    const res = await apiClient.post('/admin/students', data);
    return res.data;
  },
  generateParentCode: async (studentId: string) => {
    const res = await apiClient.post(`/admin/students/${studentId}/generate-parent-code`);
    return res.data;
  },

  // --- USERS ---
  getUsers: async (params?: { limit?: number; offset?: number }) => {
    const res = await apiClient.get<ApiResponse<UserInfo[]>>('/admin/users', { params });
    return res.data;
  },
  lockUser: async (userId: string) => {
    const res = await apiClient.post(`/admin/users/${userId}/lock`);
    return res.data;
  },
  unlockUser: async (userId: string) => {
    const res = await apiClient.post(`/admin/users/${userId}/unlock`);
    return res.data;
  },

  // --- TEACHERS ---
  getTeachers: async () => {
    const res = await apiClient.get<ApiResponse<Teacher[]>>('/admin/teachers');
    return res.data.data;
  },
  assignTeacherToClass: async (teacherId: string, classId: string) => {
    const res = await apiClient.post(`/admin/teachers/${teacherId}/classes/${classId}`);
    return res.data;
  },

  // --- PARENTS ---
  getParents: async () => {
    const res = await apiClient.get<ApiResponse<Parent[]>>('/admin/parents');
    return res.data.data;
  },
  assignParentToStudent: async (parentId: string, studentId: string) => {
    const res = await apiClient.post(`/admin/parents/${parentId}/students/${studentId}`);
    return res.data;
  }
};