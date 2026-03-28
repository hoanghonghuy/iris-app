# Phân Tích Tài Liệu `docs/` vs Codebase Thực Tế

**Ngày phân tích:** 2026-03-28
**Phạm vi:** 4 file trong `docs/`, đối chiếu code thực tế trong `apps/api` và `apps/web`

---

## Tổng Quan Kết Quả

| Tài liệu | Trạng thái | Mức độ chính xác |
|-----------|------------|-----------------|
| `remaining-issues.md` | 🔴 **LỖI THỜI NẶNG** — 6/10 mục đã được fix rồi | ~40% |
| `iris-issues-audit.md` | 🟢 **Phần lớn chính xác** — vài chi tiết nhỏ sai | ~90% |
| `staged_changes_review.md` | 🟡 **Hầu hết đúng** — 1 issue quan trọng đã được fix | ~85% |
| `google-signup-integration-proposal.md` | 🟢 **Chính xác** — proposal chưa triển khai, đúng hiện trạng | ~95% |

---

## 1. `remaining-issues.md` — 🔴 LỖI THỜI NẶNG

> [!CAUTION]
> File này gây hiểu lầm **nghiêm trọng**. Nhiều mục được liệt kê "chưa xử lý" nhưng thực tế **đã có trong code**.

### Các mục **ĐÃ SỬA** nhưng doc vẫn ghi `[ ]`:

| # | Mục trong doc | Trạng thái thực tế | Evidence |
|---|--------------|-------------------|---------|
| 2.1 | CORS permissive (allow all Origin) | ✅ **ĐÃ FIX** — dùng `originSet` allowlist | [router.go:36-64](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/http/router.go#L36-L64) |
| 2.2 | WS `CheckOrigin` return true | ✅ **ĐÃ FIX** — check origin vs `allowedOrigins` map | [chat_handler.go:212-220](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/api/v1/handlers/chat_handler.go#L212-L220) |
| 2.3 | JWT qua query string | ✅ **ĐÃ GIẢM** — chuyển sang `Sec-WebSocket-Protocol`, query token là fallback controllable bởi flag | [chat_handler.go:226-241](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/api/v1/handlers/chat_handler.go#L226-L241) |
| 3.1 | Parent code race condition | ✅ **ĐÃ FIX** — `IncrementUsageIfNotMaxed` atomic guard | [parent_code_repo.go:78-94](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/repo/parent_code_repo.go#L78-L94) |
| 3.2 | Reset token TOCTOU | ✅ **ĐÃ FIX** — `MarkUsed` với `used_at IS NULL` | [reset_token_repo.go:51-64](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/repo/reset_token_repo.go#L51-L64) |
| 4.1 | Config panic khi thiếu env | ✅ **ĐÃ FIX** — `must()` trả `error` | [config.go:88-94](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/config/config.go#L88-L94) |

### Các mục **CÒN MỞ** (doc ghi đúng):

| # | Mục | Nhận xét |
|---|-----|---------|
| 2.4 | Password reset token trong URL query | ✅ **Đúng**, vẫn còn tại [user_service.go:334-338](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/service/user_service.go#L334-L338) |
| 1.1 | Chính sách dữ liệu lịch sử | ✅ **Đúng**, đã chốt (theo `students.current_class_id`) |
| 1.2 | LIST endpoint 403 vs empty | ✅ **Đúng**, đã chốt (`200 + []` cho out-of-scope data, `403` cho endpoint forbidden) |
| 5.1 | Decompose frontend pages | ✅ **Đúng**, chưa thực hiện |

> **Khuyến nghị:** Nên viết lại hoặc xóa file này vì gây nhầm lẫn nghiêm trọng.

---

## 2. `iris-issues-audit.md` — 🟢 Chính xác

### Phần "Đã xử lý" (§2) — **TẤT CẢ ĐÚNG** ✅

| Mục | Kết luận audit | Code xác nhận |
|-----|---------------|--------------|
| §2.1 CORS allowlist | RESOLVED | ✅ `originSet` + allowlist trong [router.go:36-64](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/http/router.go#L36-L64) |
| §2.2 WS CheckOrigin | RESOLVED | ✅ Allowlist inject trong [chat_handler.go:30-44](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/api/v1/handlers/chat_handler.go#L30-L44) |
| §2.3 Config no panic | RESOLVED | ✅ `must()` trả `error` trong [config.go:88-94](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/config/config.go#L88-L94) |
| §2.4 Parent code atomic | RESOLVED | ✅ `IncrementUsageIfNotMaxed` trong [parent_code_repo.go:78-94](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/repo/parent_code_repo.go#L78-L94) |
| §2.5 Reset token atomic | RESOLVED | ✅ `MarkUsed` atomic trong [reset_token_repo.go:51-64](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/repo/reset_token_repo.go#L51-L64) |

### Phần "Còn mở" (§3) — **GẦN NHƯ ĐÚNG**, 1 chi tiết nhỏ sai:

| Mục | Kết luận audit | Kiểm tra | Nhận xét |
|-----|---------------|---------|---------|
| §3.1 WS query-token fallback | Còn mở, default tắt | ✅ **Đúng** | Flag `WSAllowQueryTokenFallback` default `false` |
| §3.2 Reset token trong URL | Còn mở | ✅ **Đúng** | `query.Set("token", plainToken)` tại `user_service.go:336` |
| §3.3 Chưa có Google sign-up parent | Còn mở | ✅ **Đúng** | Không có route `register/parent/google` trong `router.go` |
| §3.4 Authz semantics LIST | Đã chốt | ✅ **Đúng** | `200 + []` cho out-of-scope data, `403` cho endpoint forbidden |
| §3.5 Data scope lịch sử | Đã chốt | ✅ **Đúng** | Dùng `students.current_class_id` |
| §3.6 Smoke test Google flow | Chưa cover | ✅ **Đúng** | `api-smoke.ps1` không có case Google |
| §3.7 Migration 000011 rehearsal | Chưa xác nhận | ✅ **Đúng** | File migration tồn tại nhưng chưa có smoke test |
| §3.8 TODO trong `teacher_scope_repo.go` | Còn TODO | ⚠️ **THIẾU** | Doc chỉ nói `teacher_scope_repo.go` nhưng thực tế có **12 TODO** rải ở nhiều file |

### Chi tiết về TODO comments (§3.8 — audit ghi thiếu):

```
hub.go:14              — TODO: Horizontal Scaling
parent_code_service.go:154 — TODO: Parent name
response.go:14,59       — TODO: RFC 7807, TODO: CreatedWithLocation
teacher_scope_repo.go:353,536 — TODO: comment giai thich
jwt.go:53              — TODO: kiểm tra code
student_handler.go:27,61 — TODO: school_id, TODO: Location header
main.go:37,83,101       — TODO: tách helper functions
```

---

## 3. `staged_changes_review.md` — 🟡 Hầu hết đúng

### Issue #1 (🔴 Cao) — **ĐÃ ĐƯỢC FIX** ⚠️

Doc nói `NewGoogleIDTokenVerifier` gọi kể cả khi Google login disabled → crash.

**Thực tế:** Code hiện tại tại [main.go:59-65](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/cmd/api/main.go#L59-L65) **đã có conditional init** đúng như đề xuất của doc:

```go
var googleVerifier auth.GoogleTokenVerifier
if cfg.GoogleLoginEnabled {
    googleVerifier, err = auth.NewGoogleIDTokenVerifier(cfg.GoogleClientID)
    if err != nil {
        log.Fatal(err)
    }
}
```

> **Kết luận:** Issue này đã được fix. Doc ghi vấn đề nhưng code đã áp dụng đề xuất rồi.

### Issue #2-5 — Chưa kiểm chứng được đầy đủ (cần xem frontend code chi tiết)

| # | Issue | Trạng thái ước lượng |
|---|-------|---------------------|
| #2 | `onSubmitGoogle` inline → re-render | ❓ Cần check `GoogleSignInButton.tsx` chi tiết |
| #3 | Password link string matching | ❓ Cần check frontend code |
| #4 | `LinkGoogleSub` inline error | ✅ **Đúng** — thấy `errors.New(...)` tại [user_repo.go:70](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/repo/user_repo.go#L70) |
| #5 | `omitempty` vô tác dụng trên request | ✅ Đúng về mặt kỹ thuật |

### Các đánh giá "Điểm tốt" — **ĐÚNG** ✅

Verdict tổng thể, bảng đánh giá kiến trúc, migration, config patterns đều chính xác.

> [!NOTE]
> File paths trong doc dùng `huyhh1` thay vì `yuhh` (username hiện tại), nhưng đây chỉ là issue hiển thị, nội dung code reference vẫn chính xác.

---

## 4. `google-signup-integration-proposal.md` — 🟢 Chính xác

### Hiện trạng hệ thống (§3) — **ĐÚNG** ✅

| Claim | Kiểm tra |
|-------|---------|
| Đã có `POST /api/v1/auth/login/google` | ✅ [router.go:74](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/http/router.go#L74) |
| Đã có `google_sub`, `google_linked_at` trong `users` | ✅ Migration `000011` tồn tại, `FindByEmail` scan `google_sub` |
| Frontend đã có nút Google trên login | ✅ `GoogleSignInButton.tsx` tồn tại |
| Đã có RegisterParent endpoint | ✅ [router.go:78](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/http/router.go#L78) |
| Chưa có `POST /api/v1/register/parent/google` | ✅ Đúng — không có trong `router.go` |

### Vùng ảnh hưởng (§6) — **ĐÚNG** ✅

Các file path được liệt kê đều tồn tại và đúng chức năng.

### Kế hoạch triển khai (§8) — Hợp lý, chưa triển khai

---

## Tóm Tắt & Khuyến Nghị

### Hành động cần làm ngay:

| Ưu tiên | Hành động |
|---------|----------|
| 🔴 P0 | **Xóa hoặc viết lại `remaining-issues.md`** — gây hiểu lầm nghiêm trọng |
| 🟡 P1 | Cập nhật `staged_changes_review.md` — đánh dấu Issue #1 là ĐÃ FIX |
| 🟢 P2 | Bổ sung `iris-issues-audit.md` §3.8 — liệt kê đầy đủ 12 TODO, không chỉ 1 file |

### Các vấn đề thực sự còn mở (confirmed bằng code):

1. **Password reset token trong URL** (`user_service.go:336`) — Security concern
2. **WS query-token fallback** code path vẫn tồn tại (dù default tắt)
3. **Google sign-up cho parent** chưa triển khai
4. **Smoke test** chưa cover Google auth flow
5. **12 TODO comments** rải trong codebase
6. **Business policy** LIST authz + data history scope đã chốt và đã cập nhật docs
