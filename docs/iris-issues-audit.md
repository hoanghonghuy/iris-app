# Iris App - Issues / Audit Notes (current scan)

## Meta
- Date: 2026-03-23
- Scope: backend Go (apps/api) + frontend Next.js (apps/web)
- Basis: code review + OWASP/RFC references for authz/token/CORS/WS semantics

---

## [SECURITY] 1. CORS quá permissive (reflect any Origin + credentials=true)
- Severity: High
- Evidence:
  - `apps/api/internal/http/router.go` (CORS middleware: `Access-Control-Allow-Origin = origin`, `Access-Control-Allow-Credentials = true`)
- Impact:
  - Dễ mở rộng bề mặt tấn công cross-site (tùy threat model, nhất là nếu sau này chuyển sang cookie-based auth hoặc có endpoint stateful).
- Notes:
  - Best practice: allowlist origin cố định.

## [SECURITY] 2. WebSocket: CheckOrigin luôn true
- Severity: High
- Evidence:
  - `apps/api/internal/api/v1/handlers/chat_handler.go` (`upgrader.CheckOrigin` => `return true // TODO: restrict origins`)
- Impact:
  - Cho phép cross-origin WS handshake (phụ thuộc browser + cách token được truyền).
  
## [SECURITY] 3. WebSocket: Truy JWT qua query string (?token=...)
- Severity: High
- Evidence:
  - `apps/api/internal/http/router.go` (comment & route: `/chat/ws`)
  - `apps/api/internal/api/v1/handlers/chat_handler.go` (`tokenStr := c.Query("token")`)
  - Frontend: `apps/web/src/hooks/useChatWebSocket.ts` (`.../chat/ws?token=${token}`)
- Impact:
  - Token có thể bị log trong lịch sử/analytics/proxy (rủi ro disclosure).
- Reference (best practice):
  - RFC 6750: không nên pass bearer token qua page URL/query.

## [SECURITY] 4. Password reset link chứa token trong URL query (?token=...)
- Severity: Medium
- Evidence:
  - `apps/api/internal/service/user_service.go` (`resetLink := fmt.Sprintf("%s/reset-password?token=%s", ...)`)
- Impact:
  - Reset token vẫn là “secrets”. URL query có thể bị log/forward/ghi lịch sử browser.
- Notes:
  - Reset token không phải OAuth bearer, nhưng về mặt rủi ro disclosure thì nên cân nhắc giảm thiểu.

## [AUTHZ] 5. Teacher LIST endpoints thường trả 200 + mảng rỗng thay vì 403 khi không được phân công
- Severity: Medium
- Evidence:
  - `apps/api/internal/repo/teacher_scope_repo.go`:
    - Các query LIST dựa trên inner join (`teacher_classes`, `students.current_class_id`), “không match” => 0 rows => `err=nil`
  - `apps/api/internal/api/v1/handlers/teacher_scope_handler.go` có nhánh xử lý `ErrForbidden`, nhưng LIST repo không thực tế trả ra `ErrForbidden`.
- Impact:
  - Inconsistent behavior so với Parent scope (Parent có check quan hệ và trả 403 khi không phải parent).
  - Có thể làm logic “deny-by-default/authorization semantics” không rõ ràng.

## [DATA SCOPE] 6. Permission dữ liệu TEACHER/PARENT phụ thuộc `students.current_class_id` (ảnh hưởng lịch sử)
- Severity: Medium (tùy business requirement)
- Evidence:
  - `apps/api/internal/repo/teacher_scope_repo.go` (attendance/health/posts đều join qua `s.current_class_id`)
  - `apps/api/internal/repo/parent_scope_repo.go` (class posts/feed dựa vào `s.current_class_id`)
- Impact:
  - Nếu bé đổi lớp, giáo viên/phụ huynh có thể không truy được dữ liệu/bài đăng “trước đó” dù record vẫn tồn tại.
- Question to confirm:
  - Business có yêu cầu “xem lịch sử theo phân công tại thời điểm record tạo” không?

## [SECURITY] 7. Frontend lưu JWT token + user_role trong localStorage
- Severity: Medium
- Evidence:
  - `apps/web/src/lib/api/client.ts` (`localStorage.getItem('auth_token')`)
  - `apps/web/src/lib/api/client.ts` (`localStorage.setItem('user_role', role)`)
- Impact:
  - Bị đánh cắp nếu có XSS.
- Notes:
  - Nếu project hướng demo/thesis OK, nhưng best practice long-term: cân nhắc httpOnly cookies + CSRF protection.

## [AUTHZ/UX] 8. ProtectedRoute dựa trên role từ localStorage (có thể lệch với server)
- Severity: Low/Medium
- Evidence:
  - `apps/web/src/providers/AuthProvider.tsx` (khởi tạo role từ localStorage)
  - `apps/web/src/components/layout/ProtectedRoute.tsx` (guard theo `role`)
- Impact:
  - UI có thể hiển thị trang không đúng role thực (server vẫn chặn, nhưng UX và logic test kém).

## [CRYPTO/RACE] 9. Parent code usage_count: VerifyCode + IncrementUsage không atomic (race)
- Severity: Medium
- Evidence:
  - `apps/api/internal/service/parent_code_service.go`
    - `VerifyCode`: check `UsageCount >= MaxUsage`
    - sau đó `RegisterParent`: gọi `parentCodeRepo.IncrementUsage`
  - `apps/api/internal/repo/parent_code_repo.go`
    - `IncrementUsage` chỉ `usage_count = usage_count + 1` (không kiểm điều kiện)
- Impact:
  - Concurrent requests có thể vượt max_usage nếu chạy đồng thời.

## [CRYPTO/RACE] 10. Password reset token: có TOCTOU race giữa FindByTokenHash và MarkUsed
- Severity: Medium
- Evidence:
  - `apps/api/internal/service/user_service.go` (`ResetPasswordWithToken`):
    - tìm token theo hash với điều kiện `used_at IS NULL`
    - sau đó Update password và gọi `resetTokenRepo.MarkUsed`
  - `apps/api/internal/repo/reset_token_repo.go`:
    - `FindByTokenHash` có `used_at IS NULL`
    - `MarkUsed` chỉ update theo id, không check `used_at IS NULL`
- Impact:
  - Hai request song song có thể cùng pass check “unused” trước khi MarkUsed kịp update.

## [MAINTAINABILITY] 11. health_log_repo.go không được wire vào flow chính cho Teacher
- Severity: Low
- Evidence:
  - `apps/api/cmd/api/main.go` wire repositories vào `repo.Repositories` nhưng không thấy `HealthLogRepo`
  - Health của TEACHER thực tế dùng `teacher_scope_repo.go`
- Impact:
  - Code dư dễ gây nhầm lẫn khi bảo trì/audit.

## [ROBUSTNESS] 12. panic khi thiếu env trong config
- Severity: Low/Medium
- Evidence:
  - `apps/api/internal/config/config.go` (`panic("missing env: " + k)`)
- Impact:
  - Bản production thiếu env sẽ crash thay vì fail-safe.

---

## Status
- [ ] Chốt yêu cầu business cho “xem lịch sử” (current_class_id vs snapshot theo thời điểm)
- [ ] Quyết định chuẩn status code cho “không đủ quyền” ở LIST endpoints
- [ ] Xem lại threat model cho WS/CORS & token-in-URL
- [ ] Kiểm tra race condition cho parent_code và password_reset_tokens (atomic update/transaction)

---

## Frontend Follow-up Audit (2026-03-26)

### Progress update (2026-03-26)
- ✅ Hoàn tất P1-1: Chuẩn hóa API/WS base URL helper dùng chung (`runtime-config`).
- ✅ Hoàn tất P1-2: Giảm hardcode chuỗi `if pathname` trong `Header` bằng mapping declarative.
- ✅ Re-check sau refactor: `npx tsc --noEmit` và eslint scoped cho các file thay đổi đều sạch.

### 1) Cấu hình frontend chính
- `apps/web/package.json`
  - Next.js `^16.2.1`, React `19.2.3`, TypeScript `^5`, ESLint `^9`.
  - Có `lucide-react`, `class-variance-authority`, `clsx`, `tailwind-merge`, `next-themes`.
  - Có CLI `shadcn` trong devDependencies.
- `apps/web/tsconfig.json`
  - `strict: true`, `moduleResolution: bundler`, alias `@/*`.
- `apps/web/components.json`
  - `iconLibrary: "lucide"`, alias `ui: "@/components/ui"`, style `new-york`.
- `apps/web/next.config.ts`
  - Đang để mặc định, chưa có custom hardening/perf tuning riêng.

### 2) Đánh giá readability (điểm nóng)
- Các file TS/TSX lớn nhất hiện tại (line count):
  - `apps/web/src/app/teacher/attendance/useTeacherAttendancePage.ts` (~494)
  - `apps/web/src/types/index.ts` (~487)
  - `apps/web/src/app/admin/students/page.tsx` (~348)
  - `apps/web/src/app/admin/parents/page.tsx` (~314)
  - `apps/web/src/app/admin/teachers/page.tsx` (~303)
- Đánh giá nhanh:
  - Batch refactor trước đó đã tách tốt attendance/admin users theo hooks/components.
  - Vẫn còn cụm màn admin lớn (students/parents/teachers) nên tiếp tục tách dần để giảm cognitive load cho người mới.

### 3) Rà soát hardcode (web)
- Hardcode URL fallback local:
  - `apps/web/src/lib/api/client.ts`
    - `process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'`
  - `apps/web/src/hooks/useChatWebSocket.ts`
    - `process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'`
- Hardcode route/title mapping dài theo `pathname`:
  - `apps/web/src/components/layout/Header.tsx`
  - Chuỗi `if (pathname.startsWith(...))` hiện khá dài, nên gom thành config map để dễ bảo trì.
- Hardcode menu routes/role gate:
  - `apps/web/src/components/layout/Sidebar.tsx`
  - Chấp nhận được ở mức hiện tại, nhưng có thể tách config khỏi component render để gọn hơn.

### 4) Xác minh shadcn/ui + lucide-react
- `shadcn/ui` đang được dùng rộng trong app:
  - Có nhiều import từ `@/components/ui/*` trên toàn bộ màn auth/admin/teacher/parent.
- `lucide-react` đang dùng nhất quán cho icon:
  - Có import ở layout/shared/pages/components và cả UI primitives (`dropdown-menu`, `select`, `dialog`, `sheet`, `sonner`).
- Kết luận:
  - Stack UI hiện tại đã đi đúng hướng shadcn + lucide và đã được áp dụng nhất quán.

### 5) Smoke test checklist thủ công (ưu tiên chạy sau refactor)
- Auth / Redirect
  - Đăng nhập từng role (`SUPER_ADMIN`, `SCHOOL_ADMIN`, `TEACHER`, `PARENT`) và xác nhận redirect đúng dashboard.
  - Truy cập route trái quyền để xác nhận guard + redirect.
- Admin users
  - Tạo user mới.
  - Lọc/sort/search list.
  - Lock/unlock user và kiểm tra trạng thái hiển thị.
- Teacher attendance
  - Take mode: lọc theo trạng thái, search học sinh, mark + ghi chú + save.
  - Revert/cancel thay đổi.
  - History mode: lọc theo ngày/trạng thái/từ khóa, phân trang.
- Feed interactions (teacher/parent)
  - Like/unlike post, tạo comment, share.
  - Reload trang để xác nhận dữ liệu persisted từ API.

### 5.1) Kết quả smoke thực thi (API runtime) — 2026-03-26
- Môi trường test local:
  - Postgres qua Docker (`iris-postgres`, port `5433`), migration lên `10`.
  - Seed demo data từ `scripts/db/seed_demo.sql`.
  - API chạy local tại `http://localhost:8080`.
- Script chạy smoke:
  - `scripts/smoke/api-smoke.ps1`
- Kết quả tổng quan:
  - ✅ PASS `GET /api/v1/health`
  - ✅ PASS login: teacher1 / parent1 / super admin
  - ✅ PASS teacher flow: classes list, students list, mark attendance, create post, like/comment/share
  - ✅ PASS parent flow: feed list, like/comment/share cùng post teacher vừa tạo
  - ✅ PASS admin flow: users list
- Evidence ID từ phiên chạy:
  - `class_id = b56478a9-622b-4f02-ba06-8ba0d6ab5ba7`
  - `student_id = 0ee05609-199d-49dc-a835-de72c06487bb`
  - `post_id = 94f929d6-a1dc-4844-9e4d-cdcb0e22318f`
  - Parent feed chứa post smoke: `YES`

### 5.2) Kết quả smoke thực thi (Browser UI) — 2026-03-26
- Môi trường test local:
  - Frontend: `http://localhost:3000` (Next.js dev)
  - Backend API: `http://localhost:8080` (Go/Gin)
  - DB: Postgres local đã migrate + seed demo
- Script chạy smoke:
  - `scripts/smoke/ui-smoke.mjs` (Playwright, headless Chromium)
- Checklist đã chạy:
  - ✅ Guard redirect: truy cập `/admin/users` khi chưa login chuyển về `/login`
  - ✅ Login + redirect đúng role:
    - `admin@iris.local` → `/admin`
    - `teacher1@iris.local` → `/teacher`
    - `parent1@iris.local` → `/parent`
  - ✅ Admin users screen hiển thị các thành phần cốt lõi:
    - URL `/admin/users`
    - Ô tìm kiếm `Tìm theo email...`
    - Nút `Tạo user`
  - ✅ Teacher attendance screen hiển thị các thành phần cốt lõi:
    - Nút `Điểm danh hôm nay`
    - Nút `Lịch sử lớp`
    - Ô tìm kiếm `Tìm học sinh theo tên`
- Kết quả tổng quan:
  - `passed=6`, `failed=0`

### 6) Ưu tiên kỹ thuật đề xuất tiếp theo
- DONE: Tách `Header` title resolver thành cấu hình route metadata.
- DONE: Chuẩn hóa `NEXT_PUBLIC_API_URL` thành helper chung để tránh lặp fallback.
- P2: Tiếp tục decomposition các page admin lớn (`students`, `parents`, `teachers`) theo mẫu đã áp dụng ở `admin/users`.