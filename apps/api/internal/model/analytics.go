package model

// AdminAnalytics chứa các chỉ số thống kê tổng quan cho Admin Dashboard
type AdminAnalytics struct {
	TotalSchools  int `json:"total_schools"`
	TotalClasses  int `json:"total_classes"`
	TotalTeachers int `json:"total_teachers"`
	TotalStudents int `json:"total_students"`
	TotalParents  int `json:"total_parents"`
}

// TeacherAnalytics chứa các chỉ số thống kê cho Teacher Dashboard
type TeacherAnalytics struct {
	TotalClasses  int `json:"total_classes"`
	TotalStudents int `json:"total_students"`
	TotalPosts    int `json:"total_posts"`
}
