package model

// AdminAnalytics chứa các chỉ số thống kê tổng quan cho Admin Dashboard
type AdminAnalytics struct {
	TotalSchools          int     `json:"total_schools"`
	TotalClasses          int     `json:"total_classes"`
	TotalTeachers         int     `json:"total_teachers"`
	TotalStudents         int     `json:"total_students"`
	TotalParents          int     `json:"total_parents"`
	IsSuperAdmin          bool    `json:"is_super_admin"`
	SchoolName            string  `json:"school_name"`
	TodayAttendanceRate   float64 `json:"today_attendance_rate"`
	RecentHealthAlerts24h int     `json:"recent_health_alerts_24h"`
}

// TeacherAnalytics chứa các chỉ số thống kê cho Teacher Dashboard
type TeacherAnalytics struct {
	TotalClasses                int `json:"total_classes"`
	TotalStudents               int `json:"total_students"`
	TotalPosts                  int `json:"total_posts"`
	TodayAttendanceMarkedCount  int `json:"today_attendance_marked_count"`
	TodayAttendancePendingCount int `json:"today_attendance_pending_count"`
	PendingAppointments         int `json:"pending_appointments"`
	RecentHealthAlerts24h       int `json:"recent_health_alerts_24h"`
}

// ParentAnalytics chứa các chỉ số thống kê cho Parent Dashboard
type ParentAnalytics struct {
	TotalChildren               int `json:"total_children"`
	UpcomingAppointments        int `json:"upcoming_appointments"`
	RecentPosts7d               int `json:"recent_posts_7d"`
	RecentHealthAlerts7d        int `json:"recent_health_alerts_7d"`
	TodayAttendancePresentCount int `json:"today_attendance_present_count"`
	TodayAttendancePendingCount int `json:"today_attendance_pending_count"`
	RecentHealthAlerts24h       int `json:"recent_health_alerts_24h"`
}
