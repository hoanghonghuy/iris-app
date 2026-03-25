/**
 * Teacher API Service
 * Quản lý các endpoint dành riêng cho giáo viên (Teacher Scope)
 * Tương ứng với: apps/api/internal/api/v1/handlers/teacher_scope_handler.go
 */
import { apiClient } from './client';
import { Class, Student, ApiResponse, PaginationParams, MarkAttendanceRequest, CreateHealthLogRequest, CreatePostRequest, AttendanceRecord, AttendanceChangeLog, HealthLog, Post, TeacherAnalytics, PostComment, PostLikeResponse, PostShareResponse, CreatePostCommentRequest, CreatePostCommentResponse } from '@/types';

export const teacherApi = {
  /**
   * Lấy thống kê Dashboard
   * GET /api/v1/teacher/analytics
   */
  getAnalytics: async () => {
    const res = await apiClient.get<ApiResponse<TeacherAnalytics>>('/teacher/analytics');
    return res.data.data;
  },

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
   * Huỷ điểm danh đã lưu trong ngày
   * DELETE /api/v1/teacher/attendance?student_id=...&date=...
   */
  cancelAttendance: async (studentId: string, date: string) => {
    const res = await apiClient.delete('/teacher/attendance', {
      params: {
        student_id: studentId,
        date,
      },
    });
    return res.data;
  },

  /**
   * Lấy lịch sử điểm danh của học sinh
   * GET /api/v1/teacher/students/:student_id/attendance
   */
  getStudentAttendance: async (studentId: string, from?: string, to?: string) => {
    const params = new URLSearchParams();
    if (from) params.set('from', from);
    if (to) params.set('to', to);
    const query = params.toString() ? `?${params.toString()}` : '';
    const res = await apiClient.get<ApiResponse<AttendanceRecord[]>>(`/teacher/students/${studentId}/attendance${query}`);
    return res.data.data;
  },

  /**
   * Lấy lịch sử chỉnh sửa điểm danh của học sinh
   * GET /api/v1/teacher/students/:student_id/attendance-changes
   */
  getStudentAttendanceChanges: async (studentId: string, from?: string, to?: string) => {
    const params = new URLSearchParams();
    if (from) params.set('from', from);
    if (to) params.set('to', to);
    const query = params.toString() ? `?${params.toString()}` : '';
    const res = await apiClient.get<ApiResponse<AttendanceChangeLog[]>>(`/teacher/students/${studentId}/attendance-changes${query}`);
    return res.data.data;
  },

  /**
   * Lấy lịch sử chỉnh sửa điểm danh theo lớp (phân trang)
   * GET /api/v1/teacher/classes/:class_id/attendance-changes
   */
  getClassAttendanceChanges: async (
    classId: string,
    params?: {
      from?: string;
      to?: string;
      student_id?: string;
      status?: string;
      limit?: number;
      offset?: number;
    }
  ) => {
    const res = await apiClient.get<ApiResponse<AttendanceChangeLog[]>>(`/teacher/classes/${classId}/attendance-changes`, { params });
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
   * Lấy nhật ký sức khỏe của học sinh
   * GET /api/v1/teacher/students/:student_id/health
   */
  getStudentHealth: async (studentId: string, from?: string, to?: string) => {
    const params = new URLSearchParams();
    if (from) params.set('from', from);
    if (to) params.set('to', to);
    const query = params.toString() ? `?${params.toString()}` : '';
    const res = await apiClient.get<ApiResponse<HealthLog[]>>(`/teacher/students/${studentId}/health${query}`);
    return res.data.data;
  },

  /**
   * Cập nhật hồ sơ cá nhân (chỉ phone)
   * PUT /api/v1/teacher/profile
   */
  updateMyProfile: async (phone: string) => {
    const res = await apiClient.put('/teacher/profile', { phone });
    return res.data;
  },

  /**
   * Tạo bài đăng mới (cho lớp hoặc học sinh)
   * POST /api/v1/teacher/posts
   */
  createPost: async (data: CreatePostRequest) => {
    const res = await apiClient.post('/teacher/posts', data);
    return res.data;
  },

  /**
   * Lấy bài đăng của lớp
   * GET /api/v1/teacher/classes/:class_id/posts
   */
  getClassPosts: async (classId: string, params?: PaginationParams) => {
    const res = await apiClient.get<ApiResponse<Post[]>>(`/teacher/classes/${classId}/posts`, { params });
    return res.data;
  },

  /**
   * Lấy bài đăng của học sinh
   * GET /api/v1/teacher/students/:student_id/posts
   */
  getStudentPosts: async (studentId: string, params?: PaginationParams) => {
    const res = await apiClient.get<ApiResponse<Post[]>>(`/teacher/students/${studentId}/posts`, { params });
    return res.data;
  },

  /**
   * Toggle like cho bài đăng
   * POST /api/v1/teacher/posts/:post_id/like
   */
  togglePostLike: async (postId: string) => {
    const res = await apiClient.post<ApiResponse<PostLikeResponse>>(`/teacher/posts/${postId}/like`);
    return res.data.data;
  },

  /**
   * Lấy danh sách bình luận của bài đăng
   * GET /api/v1/teacher/posts/:post_id/comments
   */
  getPostComments: async (postId: string, params?: PaginationParams) => {
    const res = await apiClient.get<ApiResponse<PostComment[]>>(`/teacher/posts/${postId}/comments`, { params });
    return res.data;
  },

  /**
   * Tạo bình luận bài đăng
   * POST /api/v1/teacher/posts/:post_id/comments
   */
  createPostComment: async (postId: string, data: CreatePostCommentRequest) => {
    const res = await apiClient.post<ApiResponse<CreatePostCommentResponse>>(`/teacher/posts/${postId}/comments`, data);
    return res.data.data;
  },

  /**
   * Ghi nhận chia sẻ bài đăng
   * POST /api/v1/teacher/posts/:post_id/share
   */
  sharePost: async (postId: string) => {
    const res = await apiClient.post<ApiResponse<PostShareResponse>>(`/teacher/posts/${postId}/share`);
    return res.data.data;
  },
};