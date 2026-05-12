// School types
export interface School {
  id: string
  name: string
  address?: string
  phone?: string
  email?: string
  created_at: string
  updated_at: string
}

// Class types
export interface Class {
  id: string
  name: string
  school_id: string
  school_name?: string
  grade_level?: string
  academic_year?: string
  created_at: string
  updated_at: string
}

// Student types
export interface Student {
  id: string
  name: string
  date_of_birth?: string
  gender?: 'male' | 'female' | 'other'
  class_id: string
  class_name?: string
  school_id: string
  school_name?: string
  parent_ids?: string[]
  created_at: string
  updated_at: string
}

// Teacher types
export interface Teacher {
  id: string
  user_id: string
  name: string
  email: string
  phone?: string
  school_id: string
  school_name?: string
  subject?: string
  class_ids?: string[]
  created_at: string
  updated_at: string
}

// Parent types
export interface Parent {
  id: string
  user_id: string
  name: string
  email: string
  phone?: string
  student_ids?: string[]
  created_at: string
  updated_at: string
}

// Attendance types
export type AttendanceStatus = 'present' | 'absent' | 'late' | 'excused'

export interface AttendanceRecord {
  id: string
  student_id: string
  student_name?: string
  class_id: string
  date: string
  status: AttendanceStatus
  notes?: string
  marked_by?: string
  created_at: string
  updated_at: string
}

export interface AttendanceBulkSaveRequest {
  class_id: string
  date: string
  records: Array<{
    student_id: string
    status: AttendanceStatus
    notes?: string
  }>
}

// Health log types
export type HealthStatus = 'normal' | 'watch' | 'urgent'

export interface HealthLog {
  id: string
  student_id: string
  student_name?: string
  date: string
  status: HealthStatus
  symptoms?: string
  notes?: string
  recorded_by?: string
  created_at: string
  updated_at: string
}

// Post types
export interface Post {
  id: string
  title: string
  content: string
  author_id: string
  author_name?: string
  author_role?: string
  class_id?: string
  class_name?: string
  school_id: string
  attachments?: PostAttachment[]
  comments_count?: number
  interactions_count?: number
  created_at: string
  updated_at: string
}

export interface PostAttachment {
  id: string
  post_id: string
  file_url: string
  file_name: string
  file_type: string
  file_size?: number
  created_at: string
}

export interface PostComment {
  id: string
  post_id: string
  user_id: string
  user_name?: string
  content: string
  created_at: string
  updated_at: string
}

export type PostInteractionType = 'like' | 'love' | 'care'

export interface PostInteraction {
  id: string
  post_id: string
  user_id: string
  interaction_type: PostInteractionType
  created_at: string
}

// Appointment types
export type AppointmentStatus = 'pending' | 'confirmed' | 'completed' | 'cancelled' | 'no_show'

export interface AppointmentSlot {
  id: string
  teacher_id: string
  teacher_name?: string
  start_time: string
  end_time: string
  is_active: boolean
  max_appointments?: number
  created_at: string
  updated_at: string
}

export interface Appointment {
  id: string
  slot_id: string
  parent_id: string
  parent_name?: string
  student_id: string
  student_name?: string
  teacher_id: string
  teacher_name?: string
  start_time: string
  end_time: string
  status: AppointmentStatus
  notes?: string
  created_at: string
  updated_at: string
}

export interface AppointmentBookRequest {
  slot_id: string
  student_id: string
  notes?: string
}

// Chat types
export interface Conversation {
  id: string
  name?: string
  type: 'direct' | 'group'
  participant_ids: string[]
  participants?: ConversationParticipant[]
  last_message?: Message
  unread_count?: number
  created_at: string
  updated_at: string
}

export interface ConversationParticipant {
  user_id: string
  user_name: string
  user_role: string
  joined_at: string
}

export interface Message {
  id: string
  conversation_id: string
  sender_id: string
  sender_name?: string
  content: string
  created_at: string
  read_at?: string
}

export interface SendMessageRequest {
  conversation_id: string
  content: string
}

// WebSocket event types
export interface WebSocketMessage {
  type: 'message' | 'conversation_created' | 'conversation_updated' | 'error'
  data: any
}

// Analytics types
export interface AnalyticsSnapshot {
  total_students?: number
  total_teachers?: number
  total_parents?: number
  total_classes?: number
  attendance_rate?: number
  health_alerts?: number
  pending_appointments?: number
  unread_messages?: number
}

export interface TimeseriesDataPoint {
  timestamp: string
  value: number
  components?: Record<string, number>
}

export interface TimeseriesSeries {
  id: string
  name: string
  data: TimeseriesDataPoint[]
}

export interface AnalyticsTimeseriesResponse {
  series: TimeseriesSeries[]
  range: string
  interval: string
}

// Audit log types
export interface AuditLog {
  id: string
  user_id?: string
  user_email?: string
  action: string
  resource_type?: string
  resource_id?: string
  details?: Record<string, any>
  ip_address?: string
  user_agent?: string
  created_at: string
}
