# Iris App - Remaining Issues (pending)

**Cập nhật lần cuối:** 2026-03-28
**Trạng thái:** Đã đối chiếu với code hiện tại, chỉ giữ lại các mục thực sự còn mở.

---

## 1) Business policy đã chốt
- [x] Chính sách dữ liệu lịch sử (CHỐT): dùng `students.current_class_id` để kiểm soát phạm vi truy cập của teacher/parent.
- [x] Authz semantics cho LIST endpoints (CHỐT):
	- Trả `200 + []` khi user đã authenticated, đúng role endpoint nhưng không có bản ghi trong phạm vi được phép xem.
	- Trả `403` khi user không có quyền truy cập endpoint theo role/scope policy ở mức endpoint.

## 2) Feature gap
- [ ] Triển khai `POST /api/v1/register/parent/google` theo đề xuất `docs/google-signup-integration-proposal.md`.

## 3) Frontend follow-up (P2)
- [ ] Decompose các admin pages lớn (`students`, `parents`, `teachers`) thành hooks/services/components riêng.
- [ ] Chuyển config sidebar route/role cứng thành data-driven config.

## 4) Kiểm tra & ghi nhận
- [ ] Tổ chức PR review với checklist (tsc/eslint/api-smoke/ui-smoke).
- [x] Cập nhật docs sau khi chốt business policy (2026-03-28).

---

## Ghi chú
- Các mục security (CORS, WS CheckOrigin, race condition, config panic) đã xử lý xong.
- Reset token đã bỏ khỏi URL query (2026-03-28).
- WS query-token fallback đã loại bỏ hoàn toàn (2026-03-28).
- Smoke test đã bổ sung Google login negative cases (2026-03-28).
- Business policy đã chốt (2026-03-28):
	- Historical scope theo `students.current_class_id`.
	- LIST authz dùng `200 + []` cho out-of-scope data, `403` cho endpoint-level forbidden.
