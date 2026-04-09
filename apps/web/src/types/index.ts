// ============================================================================
// AUTH TYPES
// ============================================================================

/**
 * LoginRequest - Gửi lên backend khi đăng nhập
 * Backend: apps/api/internal/api/v1/handlers/auth_handler.go
 */
export interface LoginRequest {
  email: string;
  password: string;
}

/**
 * GoogleLoginRequest - Gửi Google ID token lên backend.
 * password là optional để xác nhận link account local lần đầu.
 */
export interface GoogleLoginRequest {
  id_token: string;
  password?: string;
}

/**
 * LoginResponse - Backend trả về sau khi đăng nhập thành công
 * Backend: apps/api/internal/service/auth_service.go
 */
export interface LoginResponse {
  access_token: string;
  token_type: string;
  expires_in: number;
}

/**
 * ActivateAccountRequest - Kích hoạt tài khoản (PUBLIC endpoint)
 * Backend: apps/api/internal/api/v1/handlers/user_handler.go
 */
export interface ActivateAccountRequest {
  email: string;
  password: string;
}

/**
 * RegisterParentRequest - Phụ huynh tự đăng ký với parent code
 * Backend: apps/api/internal/api/v1/handlers/auth_handler.go
 */
export interface RegisterParentRequest {
  email: string;
  password: string;
  parent_code: string;
}

/**
 * RegisterParentWithGoogleRequest - Phụ huynh tự đăng ký bằng Google với parent code
 * Backend: apps/api/internal/api/v1/handlers/parent_code_handler.go
 */
export interface RegisterParentWithGoogleRequest {
  id_token: string;
  parent_code: string;
}

// ============================================================================
// USER TYPES
// ============================================================================

/**
 * UserInfo - Thông tin user đầy đủ (từ endpoint /me)
 * Backend: apps/api/internal/api/v1/handlers/auth_handler.go
 */
export interface UserInfo {
  user_id: string;
  email: string;
  full_name?: string;
  status: UserStatus;
  roles: UserRole[];
  school_id?: string;
}

/**
 * UserStatus - Trạng thái tài khoản
 */
export type UserStatus = "pending" | "active" | "locked";

/**
 * UserRole - Các role trong hệ thống
 */
export type UserRole = "SUPER_ADMIN" | "SCHOOL_ADMIN" | "TEACHER" | "PARENT";

// ============================================================================
// SCHOOL TYPES
// ============================================================================

/**
 * School - Thông tin trường học
 * Backend: apps/api/internal/model/school.go
 */
export interface School {
  school_id: string;
  name: string;
  address: string;
}

/**
 * CreateSchoolRequest - Tạo trường mới
 */
export interface CreateSchoolRequest {
  name: string;
  address: string;
}

// ============================================================================
// CLASS TYPES
// ============================================================================

/**
 * Class - Thông tin lớp học
 * Backend: apps/api/internal/model/class.go
 */
export interface Class {
  class_id: string;
  school_id: string;
  name: string;
  school_year: string;
  created_at: string;
  updated_at: string;
}

/**
 * CreateClassRequest - Tạo lớp mới
 */
export interface CreateClassRequest {
  school_id: string;
  name: string;
  school_year: string;
}

// ============================================================================
// STUDENT TYPES
// ============================================================================

/**
 * Student - Thông tin học sinh
 * Backend: apps/api/internal/model/student.go
 */
export interface Student {
  student_id: string;
  school_id: string;
  current_class_id: string;
  current_class_name?: string;
  full_name: string;
  dob: string;
  gender: StudentGender;
  active_parent_code?: string;
  code_expires_at?: string;
  code_usage_count?: number;
  code_max_usage?: number;
}

/**
 * StudentParentInfo - Thông tin cha mẹ (rút gọn) trả về cùng hồ sơ học sinh
 */
export interface StudentParentInfo {
  parent_id: string;
  full_name: string;
  phone: string;
  email: string;
}

/**
 * StudentProfile - Hồ sơ chi tiết học sinh gồm cả cha mẹ
 */
export interface StudentProfile extends Student {
  parents: StudentParentInfo[];
}

/**
 * StudentGender - Giới tính học sinh
 */
export type StudentGender = "male" | "female" | "other";

/**
 * CreateStudentRequest - Tạo học sinh mới
 */
export interface CreateStudentRequest {
  school_id: string;
  class_id: string;
  full_name: string;
  dob: string;
  gender: StudentGender;
}

// ============================================================================
// TEACHER TYPES
// ============================================================================

/**
 * Teacher - Thông tin giáo viên
 * Backend: apps/api/internal/model/teacher.go
 */
export interface Teacher {
  teacher_id: string;
  user_id: string;
  email: string;
  full_name: string;
  phone: string;
  school_id: string;
  classes?: {
    class_id: string;
    name: string;
  }[];
}

/**
 * CreateTeacherRequest - Tạo giáo viên mới
 */
export interface CreateTeacherRequest {
  email: string;
  full_name: string;
  phone?: string;
  role: "TEACHER";
}

/**
 * UpdateTeacherRequest - Cập nhật thông tin giáo viên
 */
export interface UpdateTeacherRequest {
  full_name?: string;
  phone?: string;
}

// ============================================================================
// PARENT TYPES
// ============================================================================

/**
 * Parent - Thông tin phụ huynh
 * Backend: apps/api/internal/model/parent.go
 */
export interface Parent {
  parent_id: string;
  user_id: string;
  email: string;
  full_name: string;
  phone: string;
  school_id: string;
  children?: {
    student_id: string;
    full_name: string;
  }[];
}

/**
 * CreateParentRequest - Tạo phụ huynh mới
 */
export interface CreateParentRequest {
  email: string;
  full_name: string;
  phone?: string;
  role: "PARENT";
}

/**
 * SchoolAdmin - Thông tin quản trị trường
 */
export interface SchoolAdmin {
  admin_id: string;
  user_id: string;
  email?: string;
  full_name?: string;
  phone?: string;
  school_id: string;
  school_name?: string;
}

/**
 * ParentCodeResponse - Dữ liệu trả về sau khi tạo parent code
 */
export interface ParentCodeResponse {
  student_id: string;
  parent_code: string;
  message?: string;
  max_usage?: number;
  expires_at: string;
}

// ============================================================================
// ATTENDANCE TYPES
// ============================================================================

/**
 * AttendanceRecord - Bản ghi điểm danh
 * Backend: apps/api/internal/model/attendance.go
 */
export interface AttendanceRecord {
  attendance_id: string;
  student_id: string;
  date: string;
  status: AttendanceStatus;
  check_in_at?: string;
  check_out_at?: string;
  note?: string;
  recorded_at: string;
  recorded_by: string;
}

export interface AttendanceChangeLog {
  change_id: string;
  attendance_id: string;
  student_id: string;
  student_name?: string;
  date: string;
  change_type: "create" | "update" | "delete";
  old_status?: AttendanceStatus;
  new_status?: AttendanceStatus;
  old_note?: string;
  new_note?: string;
  changed_by: string;
  changed_at: string;
}

/**
 * AttendanceStatus - Trạng thái điểm danh
 */
export type AttendanceStatus = "present" | "absent" | "late" | "excused";

/**
 * MarkAttendanceRequest - Điểm danh học sinh
 */
export interface MarkAttendanceRequest {
  student_id: string;
  date: string;
  status: AttendanceStatus;
  check_in_at?: string;
  check_out_at?: string;
  note?: string;
}

// ============================================================================
// HEALTH TYPES
// ============================================================================

/**
 * HealthLog - Bản ghi sức khỏe
 * Backend: apps/api/internal/model/health.go
 */
export interface HealthLog {
  health_log_id: string;
  student_id: string;
  recorded_at: string;
  temperature?: number;
  symptoms?: string;
  severity: HealthSeverity;
  note?: string;
  recorded_by: string;
}

/**
 * HealthSeverity - Mức độ nghiêm trọng sức khỏe
 */
export type HealthSeverity = "normal" | "watch" | "urgent";

/**
 * CreateHealthLogRequest - Tạo bản ghi sức khỏe
 */
export interface CreateHealthLogRequest {
  student_id: string;
  temperature?: number;
  symptoms?: string;
  severity: HealthSeverity;
  note?: string;
}

// ============================================================================
// POST TYPES
// ============================================================================

/**
 * Post - Bài đăng (Newsfeed)
 * Backend: apps/api/internal/model/post.go
 */
export interface Post {
  post_id: string;
  author_user_id: string;
  scope_type: PostScope;
  school_id?: string;
  class_id?: string;
  student_id?: string;
  type: PostType;
  content: string;
  like_count: number;
  comment_count: number;
  share_count: number;
  liked_by_me: boolean;
  created_at: string;
  updated_at: string;
}

/**
 * PostComment - Bình luận của bài đăng
 */
export interface PostComment {
  comment_id: string;
  post_id: string;
  author_user_id: string;
  author_display: string;
  content: string;
  created_at: string;
}

/**
 * PostLikeResponse - Kết quả thao tác like/unlike
 */
export interface PostLikeResponse {
  post_id: string;
  liked_by_me: boolean;
  like_count: number;
}

/**
 * PostShareResponse - Kết quả thao tác share
 */
export interface PostShareResponse {
  post_id: string;
  share_count: number;
}

/**
 * CreatePostCommentRequest - Tạo bình luận bài đăng
 */
export interface CreatePostCommentRequest {
  content: string;
}

/**
 * CreatePostCommentResponse - Kết quả tạo bình luận
 */
export interface CreatePostCommentResponse {
  post_id: string;
  comment_count: number;
  comment: PostComment;
}

/**
 * PostScope - Phạm vi bài đăng
 */
export type PostScope = "school" | "class" | "student";

/**
 * PostType - Loại bài đăng
 */
export type PostType =
  | "announcement"
  | "activity"
  | "daily_note"
  | "health_note";

/**
 * CreatePostRequest - Tạo bài đăng
 */
export interface CreatePostRequest {
  scope_type: PostScope;
  school_id?: string;
  class_id?: string;
  student_id?: string;
  type: PostType;
  content: string;
}

// ============================================================================
// PAGINATION & RESPONSE WRAPPERS
// ============================================================================

/**
 * Pagination - Phân trang
 * Backend: apps/api/internal/response/response.go
 */
export interface Pagination {
  total: number;
  limit: number;
  offset: number;
  has_more: boolean;
}

/**
 * ApiResponse - Response thành công
 * Backend: { "data": {...}, "pagination": {...} }
 */
export interface ApiResponse<T = unknown> {
  data: T;
  pagination?: Pagination;
}

/**
 * ApiError - Response lỗi
 * Backend: { "error": "message", "error_code": "MACHINE_CODE" }
 */
export interface ApiError {
  error: string;
  error_code?: string;
}

// ============================================================================
// REQUEST TYPES (Helpers)
// ============================================================================

/**
 * PaginationParams - Params cho phân trang
 */
export interface PaginationParams {
  limit?: number;
  offset?: number;
}

/**
 * AssignTeacherToClassRequest - Gán giáo viên vào lớp
 */
export interface AssignTeacherToClassRequest {
  teacher_id: string;
  class_id: string;
}

/**
 * AssignParentToStudentRequest - Gán phụ huynh vào học sinh
 */
export interface AssignParentToStudentRequest {
  parent_id: string;
  student_id: string;
  relationship: string;
}

// ============================================================================
// ANALYTICS TYPES
// ============================================================================

/**
 * AdminAnalytics - Thống kê Dashboard Admin
 */
export interface AdminAnalytics {
  total_schools: number;
  total_classes: number;
  total_teachers: number;
  total_students: number;
  total_parents: number;
  is_super_admin: boolean;
  school_name: string;
  today_attendance_rate: number;
  today_pending_appointments: number;
  recent_health_alerts_24h: number;
}

/**
 * TeacherAnalytics - Thống kê Dashboard Giáo viên
 */
export interface TeacherAnalytics {
  total_classes: number;
  total_students: number;
  total_posts: number;
  today_attendance_marked_count: number;
  today_attendance_pending_count: number;
  pending_appointments: number;
  recent_health_alerts_24h: number;
}

/**
 * ParentAnalytics - Thống kê Dashboard Phụ huynh
 */
export interface ParentAnalytics {
  total_children: number;
  upcoming_appointments: number;
  recent_posts_7d: number;
  recent_health_alerts_7d: number;
  today_attendance_present_count: number;
  today_attendance_pending_count: number;
  recent_health_alerts_24h: number;
}

// ============================================================================
// APPOINTMENT TYPES
// ============================================================================

export type AppointmentStatus = "pending" | "confirmed" | "cancelled" | "completed" | "no_show";

export interface AppointmentSlot {
  slot_id: string;
  teacher_id: string;
  teacher_name?: string;
  class_id: string;
  class_name?: string;
  start_time: string;
  end_time: string;
  note?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Appointment {
  appointment_id: string;
  slot_id: string;
  parent_id: string;
  parent_name?: string;
  student_id: string;
  student_name?: string;
  teacher_id?: string;
  teacher_name?: string;
  class_id?: string;
  class_name?: string;
  status: AppointmentStatus;
  note?: string;
  cancel_reason?: string;
  start_time: string;
  end_time: string;
  confirmed_at?: string;
  completed_at?: string;
  cancelled_at?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateAppointmentSlotRequest {
  class_id: string;
  start_time: string;
  end_time?: string;
  duration_minutes?: number;
  note?: string;
}

export interface CreateAppointmentRequest {
  student_id: string;
  slot_id: string;
  note?: string;
}

// ============================================================================
// AUDIT LOG TYPES
// ============================================================================

export interface AuditLog {
  audit_log_id: string;
  actor_user_id: string;
  actor_role?: string;
  action: string;
  entity_type: string;
  entity_id?: string;
  details?: Record<string, unknown>;
  created_at: string;
}

// ============================================================================
// CHAT TYPES
// ============================================================================

/**
 * ParticipantInfo - Thông tin cơ bản của thành viên trong cuộc hội thoại
 * Backend: apps/api/internal/model/chat.go
 */
export interface ParticipantInfo {
  user_id: string;
  email: string;
  full_name?: string;
}

/**
 * Conversation - Cuộc hội thoại
 * Backend: apps/api/internal/model/chat.go
 */
export interface Conversation {
  conversation_id: string;
  type: "direct" | "group";
  name?: string;
  created_at: string;
  participants: ParticipantInfo[];
}

/**
 * Message - Tin nhắn trong cuộc hội thoại
 * Backend: apps/api/internal/model/chat.go
 */
export interface Message {
  message_id: string;
  conversation_id: string;
  sender_id: string;
  sender_email: string;
  content: string;
  created_at: string;
}

/**
 * WSEvent - Sự kiện WebSocket realtime
 * Backend: apps/api/internal/ws/hub.go
 */
export interface WSEvent {
  type: "new_message" | "conversation_created";
  data: unknown;
}

/**
 * WSSendMessage - Cấu trúc JSON gửi tin nhắn qua WebSocket
 */
export interface WSSendMessage {
  conversation_id: string;
  content: string;
}
