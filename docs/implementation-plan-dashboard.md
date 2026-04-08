# Implementation Plan – Dashboard Redesign V2 (Admin, Teacher, Parent)

## 🎯 Goal
Nâng cấp dashboard cho 3 nhóm người dùng chính nhưng vẫn bám chặt codebase hiện tại:
1. Super Admin
2. School Admin
3. Teacher
4. Parent

Mục tiêu chính:
- Không đổi endpoint hiện có, chỉ mở rộng payload trả về để tránh break client.
- Tuân thủ naming snake_case cho JSON/API type.
- Tuân thủ auth scope hiện tại trong router.
- Tuân thủ stack test hiện tại (Go test + Vitest).
- Tuân thủ checklist tại docs/checklist-review-ai-code.md.

## ✅ Nguyên tắc tương thích với code hiện tại

### API/Contract
- Giữ nguyên endpoint hiện có:
	- GET /api/v1/admin/analytics
	- GET /api/v1/admin/audit-logs
	- GET /api/v1/teacher/analytics
	- GET /api/v1/parent/analytics
- Không đổi cấu trúc response wrapper hiện có: data + pagination (nếu có).
- Chỉ thêm field mới theo snake_case, không đổi hoặc xóa field cũ.

### Authorization
- Giữ đúng behavior hiện tại:
	- Admin routes: SUPER_ADMIN và SCHOOL_ADMIN.
	- Audit logs: cả SUPER_ADMIN và SCHOOL_ADMIN đều được truy cập; SCHOOL_ADMIN bị scope theo school_id.
- Nếu muốn đổi sang Super Admin only cho audit logs, phải tạo PR riêng vì là thay đổi behavior.

### Pagination
- Giữ default limit cho audit logs là 20 (đúng với handler hiện tại).
- UI có thể chọn limit nhỏ hơn (ví dụ 10/20/50), nhưng không đổi default backend trong scope này.

### Design System
- Giữ NurturedLayer hiện có (palette pastel + token trong globals.css).
- Giữ font Geist hiện có, không đổi sang font khác.

## 📦 Scope chi tiết theo layer

| Layer | Thay đổi | Ghi chú tương thích |
|------|----------|---------------------|
| Backend (Go) | Mở rộng AdminAnalytics, TeacherAnalytics, ParentAnalytics bằng field mới cho insight hôm nay/24h. | Không đổi endpoint, không đổi field cũ.
| Backend (Repo) | Bổ sung các hàm COUNT chuyên dụng theo vai trò và scope (global/school/teacher/parent). | Ưu tiên COUNT + WHERE, không SELECT *.
| Frontend (Web) | Cập nhật type tại src/types/index.ts; cập nhật trang /admin, /teacher, /parent để render thêm khối insight. | Không thay đổi flow auth/router.
| Frontend (Components) | Thêm component tái sử dụng TodayCard, CompactList, ChildCard (nếu cần). | Tách component để giảm trùng lặp.
| Testing | Thêm unit test Go cho service/repo analytics, thêm Vitest component tests cho dashboard. | Không dùng Jest mới; dùng Vitest theo package hiện tại.

## 🧩 Contract mở rộng đề xuất (snake_case)

### 1) AdminAnalytics
Giữ nguyên:
- total_schools
- total_classes
- total_teachers
- total_students
- total_parents

Thêm mới:
- is_super_admin: boolean
- school_name: string | null
- today_attendance_rate: number
- today_pending_appointments: number
- recent_health_alerts_24h: number

Ghi chú:
- UI nên ưu tiên phân nhánh bằng role từ AuthProvider; is_super_admin là field hỗ trợ hiển thị.

### 2) TeacherAnalytics
Giữ nguyên:
- total_classes
- total_students
- total_posts

Thêm mới:
- today_attendance_marked_count: number
- today_attendance_pending_count: number
- pending_appointments: number
- recent_health_alerts_24h: number

### 3) ParentAnalytics
Giữ nguyên:
- total_children
- upcoming_appointments
- recent_posts_7d
- recent_health_alerts_7d

Thêm mới:
- today_attendance_present_count: number
- today_attendance_pending_count: number
- recent_health_alerts_24h: number

## 🖥️ UI Plan theo role

### A. Admin Dashboard
- Hero header giữ hiện tại, thêm subtitle theo role:
	- SUPER_ADMIN: toàn hệ thống
	- SCHOOL_ADMIN: theo trường đang quản trị
- Giữ các stat card hiện có, thêm nhóm Today Insights gồm:
	- Attendance rate hôm nay
	- Pending appointments
	- Health alerts 24h
- Audit logs page:
	- Giữ filter hiện tại
	- Thêm pagination controls dùng limit/offset + pagination metadata trả về
	- Không hardcode limit 100

### B. Teacher Dashboard
- Hero header: chào giáo viên + lớp phụ trách.
- Stat cards compact: tổng lớp, tổng trẻ, tổng bài đăng.
- Today Insights:
	- today_attendance_marked_count
	- today_attendance_pending_count
	- pending_appointments
	- recent_health_alerts_24h
- Compact list lớp của tôi (mỗi dòng có Quản lý).
- Recent posts: lấy từ API hiện có theo lớp giáo viên phụ trách.

### C. Parent Dashboard
- Hero header: chào phụ huynh + danh sách con.
- Stat cards: upcoming_appointments, recent_health_alerts_24h.
- Child cards:
	- tên + lớp
	- attendance status hôm nay
	- health alert gần nhất (nếu có)
	- nút xem chi tiết
- Feed compact 3-5 bài từ parent feed hiện có.

## 🧪 Testing Plan (đồng bộ stack hiện tại)

### Backend
- Unit test cho AnalyticsService:
	- Super Admin path và School Admin path
	- Teacher path
	- Parent path
	- Edge cases: 0 data, nil school_id, repo error propagation
- Repo tests (nếu có test DB harness):
	- COUNT query theo today/24h trả đúng kết quả
	- Scope đúng theo school/teacher/parent

### Frontend
- Vitest + Testing Library:
	- Render dashboard với payload cũ (backward compatibility)
	- Render dashboard với payload mở rộng
	- Empty state và loading state
	- Pagination audit logs (offset thay đổi đúng)

### Smoke
- Chạy lại script smoke hiện có cho API/UI sau khi merge.

## 🗓️ Timeline đề xuất
| Day | Activity |
|-----|----------|
| Day 1 | Chốt contract field mới + cập nhật types + mock data.
| Day 2 | Backend repo/service cho analytics mở rộng + unit tests.
| Day 3 | UI cập nhật admin/teacher/parent dashboard + vitest components.
| Day 4 | Audit-log pagination UI + integration validation + bug fix.
| Day 5 | Review theo checklist + docs update + merge.

## ⚠️ Risks & Watchouts
- Timezone: phải chốt rõ today dùng timezone nào (DB/server) và nhất quán giữa attendance/appointment/health.
- Query cost: với COUNT theo today/24h cần kiểm tra index coverage trước khi merge.
- Scope leakage: mọi query mới phải filter đúng theo role scope để tránh lộ dữ liệu chéo trường.
- Backward compatibility: FE phải hoạt động cả khi backend chưa deploy field mới.

## 📌 Rollout Checklist
- [ ] Cập nhật API types bằng snake_case.
- [ ] Không đổi endpoint path hiện có.
- [ ] Không xóa field JSON cũ.
- [ ] Bổ sung unit tests backend cho analytics mở rộng.
- [ ] Bổ sung vitest frontend cho dashboard.
- [ ] Kiểm tra pagination audit-log theo metadata response.
- [ ] Chạy lint + typecheck + go test trước PR.
- [ ] Review theo docs/checklist-review-ai-code.md.

## ❓Open Questions
- Có cần đổi policy audit logs thành Super Admin only không, hay giữ đúng behavior hiện tại?
- Today attendance ở Parent dashboard tính theo từng con hay tổng hợp theo tất cả con?
- Có cần hiển thị school_name trên Admin dashboard cho SUPER_ADMIN (ví dụ “Toàn hệ thống”) hay để null?

---
