# Iris App - Remaining Issues (pending)

Dựa trên audit hiện tại (`docs/iris-issues-audit.md`) và tiến độ đã làm tới 2026-03-26, các điểm sau vẫn chưa triển khai xong:

## 1) Business decision / policy
- [ ] Xác định chính sách dữ liệu lịch sử: `teacher/parent` nên truy dữ liệu theo `students.current_class_id` hay theo snapshot thời điểm record tạo. (Issue #6 in audit) 
- [ ] Quy định rõ status code cho LIST endpoints khi không đủ quyền: 403 (forbidden) hay 200 + empty array. (Issue #5 in audit)

## 2) Security hardening (API)
- [ ] CORS hiện đang permissive (allow all Origin + credentials). Triển khai allowlist origin cụ thể. (Issue #1)
- [ ] WebSocket `CheckOrigin` đang return true; thay bằng logic tương ứng origin hợp lệ. (Issue #2)
- [ ] Tránh chuyển JWT qua query string `?token=`; dùng header `Authorization: Bearer ...` hoặc mô hình secure websocket (Issue #3).
- [ ] Xem xét `password reset token` không nên nằm trong URL query. (Issue #4)

## 3) Concurrency/Race conditions
- [ ] Parent code `usage_count` cần yêu cầu atomic update/transaction với điều kiện `max_usage` (Issue #9)
- [ ] Password reset token `FindByTokenHash` + `MarkUsed` cần lock hoặc atomic update to avoid TOCTOU (Issue #10)

## 4) Backend robustness
- [ ] Thay hàm panic trong config bằng error trả về để app khởi động an toàn hơn (Issue #12)
- [ ] Kiểm tra wire repo `HealthLogRepo`/scope route nếu cần (Issue #11)

## 5) Frontend follow-up (P2)
- [ ] Decompose các page admin lớn:
  - `apps/web/src/app/admin/students/page.tsx`
  - `apps/web/src/app/admin/parents/page.tsx`
  - `apps/web/src/app/admin/teachers/page.tsx`
  - Giảm độ lớn file, tách thành hooks/service/component riêng (mẫu đã làm với `admin/users`)
- [ ] Chuyển khẩu bị hardcode trong `Sidebar` route/role config thành data-driven config (next step sau mapping Header đã xong).

## 6) Kiểm tra & ghi nhận
- [ ] Tổ chức pull request review với checklist đã chạy (tsc/eslint/api-smoke/ui-smoke) và bổ sung các test còn thiếu cho kịch bản bảo mật/cross-role.
- [ ] Cập nhật lại file docs/audit khi có quyết định business (khoảng policy cấp quyền, data life history rules).

---

## Ghi chú
- Đã xong P1 (runtime-config + Header route mapping), đã có script Smoke API/UI chạy pass.
- File điều hướng hiện chưa sửa `Sidebar` + chính sách role-state sync với server.
- Sau khi hoàn tất trên, ưu tiên đóng nốt phần security & race issues trước khi release.
