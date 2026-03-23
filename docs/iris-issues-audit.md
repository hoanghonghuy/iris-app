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