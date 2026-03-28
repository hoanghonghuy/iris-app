# Iris App - Remaining Issues (pending)

**Cập nhật lần cuối:** 2026-03-28
**Trạng thái:** Đã đối chiếu với code hiện tại, chỉ giữ lại các mục thực sự còn mở.

---

## 1) Business policy chưa chốt
- [ ] Xác định chính sách dữ liệu lịch sử: teacher/parent truy dữ liệu theo `students.current_class_id` hay theo snapshot thời điểm record tạo.
- [ ] Quy định rõ status code cho LIST endpoints khi không đủ quyền: `403` hay `200 + []`.

## 2) Feature gap
- [ ] Triển khai `POST /api/v1/register/parent/google` theo đề xuất `docs/google-signup-integration-proposal.md`.

## 3) Frontend follow-up (P2)
- [ ] Decompose các admin pages lớn (`students`, `parents`, `teachers`) thành hooks/services/components riêng.
- [ ] Chuyển config sidebar route/role cứng thành data-driven config.

## 4) Kiểm tra & ghi nhận
- [ ] Tổ chức PR review với checklist (tsc/eslint/api-smoke/ui-smoke).
- [ ] Cập nhật docs khi có quyết định business policy.

---

## Ghi chú
- Các mục security (CORS, WS CheckOrigin, race condition, config panic) đã xử lý xong.
- Reset token đã bỏ khỏi URL query (2026-03-28).
- WS query-token fallback đã loại bỏ hoàn toàn (2026-03-28).
- Smoke test đã bổ sung Google login negative cases (2026-03-28).
