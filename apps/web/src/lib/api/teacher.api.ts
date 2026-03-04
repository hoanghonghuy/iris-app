/**
 * Teacher API Service
 * Quản lý các endpoint dành riêng cho giáo viên (Teacher Scope)
 * Tương ứng với: apps/api/internal/api/v1/handlers/teacher_scope_handler.go
 */
import { apiClient } from './client';
import { Class, Student, ApiResponse, MarkAttendanceRequest, CreateHealthLogRequest, CreatePostRequest } from '@/types';

export const teacherApi = {
  /**
   * Lấy danh sách lớp học được phân công cho giáo viên hiện tại
   * GET /api/v1/teacher/classes
   */
  getMyClasses: async () => {
    const res = await apiClient.get<ApiResponse<Class[]>>('/teacher/classes');
    return res.data.data;
  },

  /**
   * Lấy danh sách học sinh trong một lớp cụ thể
   * GET /api/v1/teacher/classes/:class_id/students
   */
  getStudentsInClass: async (classId: string) => {
    const res = await apiClient.get<ApiResponse<Student[]>>(`/teacher/classes/${classId}/students`);
    return res.data.data;
  },

  /**
   * Điểm danh cho học sinh
   * POST /api/v1/teacher/attendance
   */
  markAttendance: async (data: MarkAttendanceRequest) => {
    const res = await apiClient.post('/teacher/attendance', data);
    return res.data;
  },

  /**
   * Tạo nhật ký sức khỏe cho học sinh
   * POST /api/v1/teacher/health
   */
  createHealthLog: async (data: CreateHealthLogRequest) => {
    const res = await apiClient.post('/teacher/health', data);
    return res.data;
  },

  /**
   * Tạo bài đăng mới (cho lớp hoặc học sinh)
   * POST /api/v1/teacher/posts
   */
  createPost: async (data: CreatePostRequest) => {
    const res = await apiClient.post('/teacher/posts', data);
    return res.data;
  }
};