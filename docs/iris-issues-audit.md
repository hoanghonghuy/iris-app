# Iris App - Issues Audit (Re-validated)

## Meta
- Date: 2026-03-27
- Reviewer: Copilot re-check theo yêu cầu
- Scope: backend `apps/api`, frontend `apps/web`, scripts/docs liên quan
- Baseline branch check: `git diff --name-status develop...HEAD` không có khác biệt (đang đồng bộ với `develop`)

---

## 1) Kết luận tổng quan
Tài liệu audit cũ có một số điểm đã lỗi thời. Sau khi rà soát lại mã nguồn hiện tại, các vấn đề sau đã **được xử lý**, và các mục còn lại bên dưới là danh sách **chưa triển khai xong / còn mở**.

---

## 2) Đã xử lý (đã xác nhận bằng code)

### 2.1 CORS đã chuyển sang allowlist
- Trạng thái: RESOLVED
- Evidence:
  - `apps/api/internal/http/router.go`
  - Có `originSet` + chỉ set `Access-Control-Allow-Origin` khi origin nằm trong allowlist.
  - Preflight origin không hợp lệ bị `403`.

### 2.2 WebSocket `CheckOrigin` không còn `return true`
- Trạng thái: RESOLVED
- Evidence:
  - `apps/api/internal/api/v1/handlers/chat_handler.go`
  - `CheckOrigin` đối chiếu origin với allowlist đã inject vào handler.

### 2.3 Config không còn panic khi thiếu env
- Trạng thái: RESOLVED
- Evidence:
  - `apps/api/internal/config/config.go`
  - Hàm `must(...)` trả `error` thay vì `panic(...)`.

### 2.4 Race của parent code usage đã có atomic guard
- Trạng thái: RESOLVED
- Evidence:
  - `apps/api/internal/repo/parent_code_repo.go`: `IncrementUsageIfNotMaxed(...)`
  - `apps/api/internal/service/parent_code_service.go`: map `ErrNoRowsUpdated` về business error.

### 2.5 Race của reset token đã có atomic guard
- Trạng thái: RESOLVED
- Evidence:
  - `apps/api/internal/repo/reset_token_repo.go`: `MarkUsed(...)` với điều kiện `used_at IS NULL`
  - `apps/api/internal/service/user_service.go`: xử lý `ErrNoRowsUpdated`.

---

## 3) Còn mở / chưa triển khai xong

## [P0] Security/Hardening

### 3.1 WebSocket vẫn còn fallback token qua query string (bật bằng flag)
- Severity: High (khi bật fallback), Low (khi giữ mặc định tắt)
- Evidence:
  - Backend: `apps/api/internal/api/v1/handlers/chat_handler.go`
  - Backend config flag: `apps/api/internal/config/config.go` (`WS_ALLOW_QUERY_TOKEN_FALLBACK`)
  - Frontend fallback URL: `apps/web/src/hooks/useChatWebSocket.ts`
- Impact:
  - Nếu bật fallback, token có nguy cơ lộ qua logs/proxy/history.
- Khuyến nghị:
  - Giữ fallback tắt ở mọi môi trường production.
  - Kế hoạch dài hạn: bỏ hẳn code path query token.

### 3.2 Password reset vẫn đưa token vào URL query
- Severity: Medium
- Evidence:
  - `apps/api/internal/service/user_service.go`
  - Đang set `resetURL.Path = "/reset-password"` và `query.Set("token", plainToken)`.
- Impact:
  - Token reset có thể bị lộ qua logs/history/referrer.
- Khuyến nghị:
  - Chuyển sang flow chỉ dùng reset code nhập tay (không append token vào URL).

## [P1] Feature gap

### 3.3 Chưa có Google sign-up cho parent (mới có Google login)
- Severity: Medium
- Evidence:
  - Đã có: `POST /api/v1/auth/login/google` trong `apps/api/internal/http/router.go`
  - Chưa có: `POST /api/v1/register/parent/google`
  - Proposal đang có: `docs/google-signup-integration-proposal.md`
- Impact:
  - Chưa đáp ứng use case đăng ký mới bằng Google.

## [P1] Business policy đã chốt (2026-03-28)

### 3.4 Authz semantics cho LIST endpoint
- Severity: Medium
- Quy định đã chốt:
  - `200 + []` khi user đã authenticated, đúng role endpoint nhưng không có bản ghi trong phạm vi dữ liệu được phép xem.
  - `403` khi user không có quyền truy cập endpoint ở mức role/scope policy.
- Evidence:
  - Pattern LIST hiện tại ở scope repo thiên về empty result theo join filter.
- Impact:
  - Giảm không nhất quán contract giữa các LIST endpoint.

### 3.5 Data scope lịch sử theo `current_class_id` hay snapshot
- Severity: Medium
- Quy định đã chốt:
  - Dùng `students.current_class_id` làm nguồn kiểm soát truy cập cho teacher/parent ở thời điểm truy vấn.
  - Không dùng snapshot assignment-at-time trong phase hiện tại.
- Evidence:
  - Nhiều query teacher/parent scope dựa vào `students.current_class_id`.
- Impact:
  - Hành vi truy cập lịch sử nhất quán với implementation hiện tại và không cần thay đổi schema.

## [P2] Testing completeness

### 3.6 Smoke test chưa bao phủ Google auth flow
- Severity: Medium
- Evidence:
  - `scripts/smoke/api-smoke.ps1` chưa có case cho `/auth/login/google`.
- Impact:
  - Thiếu regression guard cho luồng Google login.

### 3.7 Chưa thấy bằng chứng chạy rehearsal migration `000011` up/down trong chuỗi smoke hiện tại
- Severity: Low/Medium
- Evidence:
  - Có migration file:
    - `apps/api/migrations/000011_google_identity_linking.up.sql`
    - `apps/api/migrations/000011_google_identity_linking.down.sql`
  - Chưa có bước xác nhận rõ trong script smoke hiện tại.
- Impact:
  - Rủi ro nhỏ khi rollout DB nếu chưa rehearse lên/xuống.

## [P3] Maintainability

### 3.8 Còn TODO comment trong repo
- Severity: Low
- Evidence (cập nhật 2026-03-28):
  - `apps/api/cmd/api/main.go` — 3 TODO: tách helper functions (giữ lại, cải tiến dài hạn)
  - `apps/api/internal/ws/hub.go` — TODO: Horizontal Scaling (giữ lại, thiết kế tương lai)
  - `apps/api/internal/service/parent_code_service.go` — TODO: Parent name (giữ lại, cần business decision)
  - `apps/api/internal/api/v1/handlers/student_handler.go` — TODO: school_id source (giữ lại, cần business decision)
  - Đã dọn dẹp: `response.go`, `teacher_scope_repo.go`, `jwt.go`, `student_handler.go:61` (2026-03-28)
- Impact:
  - Không ảnh hưởng runtime, các TODO còn lại đều có giá trị ghi chú cho tương lai.

---

## 4) Trạng thái checklist triển khai (re-check)

- [x] Implement backend Google config
- [x] Add Google login endpoint flow
- [x] Add migration and user fields
- [x] Integrate frontend Google button
- [ ] Run tests and smoke checks cho Google flow (negative + positive)

Ghi chú:
- Các mục đầu checklist đã có trong code.
- Mục cuối cùng chưa đầy đủ vì smoke hiện tại chưa cover luồng Google.

---

## 5) Đề xuất ưu tiên thực thi tiếp theo

1. P0: Tắt tuyệt đối WS query-token fallback ở production và lập kế hoạch bỏ code path fallback.
2. P0: Chuyển reset-password sang cơ chế không để token trong URL query.
3. P1: Triển khai endpoint `POST /api/v1/register/parent/google` theo proposal đã lưu.
4. P1: Theo dõi tính phù hợp của policy đã chốt sau UAT; nếu cần snapshot lịch sử, lập migration riêng.
5. P2: Mở rộng `scripts/smoke/api-smoke.ps1` cho Google login + migration 000011 rehearsal.

---

## 6) Ghi chú xác nhận lần re-check này
- File này đã được viết mới hoàn toàn theo yêu cầu (xóa nội dung cũ, ghi lại từ đầu).
- Nội dung phản ánh trạng thái mã nguồn hiện tại, không giữ lại các kết luận đã lỗi thời từ bản audit trước.
