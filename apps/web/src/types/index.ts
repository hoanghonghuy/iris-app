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
  full_name: string;
  dob: string;
  gender: StudentGender;
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
  created_at: string;
  updated_at: string;
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
 * Backend: { "error": "message" }
 */
export interface ApiError {
  error: string;
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
