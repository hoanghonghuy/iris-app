# Đề xuất tích hợp Google Sign-Up cho Iris

Ngày cập nhật: 2026-03-27
Trạng thái: Đề xuất, chưa triển khai

## 1) Mục tiêu

Tích hợp đăng ký bằng Google theo cách không phá vỡ nghiệp vụ hiện tại của Iris.

Mục tiêu cụ thể:
- Bổ sung đăng ký Google cho phụ huynh.
- Vẫn bắt buộc xác thực parent code để liên kết với học sinh.
- Giữ nguyên luồng tạo tài khoản teacher/admin hiện có.

## 2) Phạm vi và phi phạm vi

Trong phạm vi:
- Google Sign-Up cho vai trò parent.
- Xác thực id_token phía backend.
- Tạo user parent, gán role PARENT, link parent-student theo parent code.
- Trả JWT để auto-login sau đăng ký.

Ngoài phạm vi:
- Không auto-provision cho teacher/admin.
- Không thay đổi luồng activate token của teacher.
- Không bật Google One Tap trong phase này.

## 3) Hiện trạng hệ thống

Backend đã có:
- Đăng nhập Google: POST /api/v1/auth/login/google.
- Lưu liên kết Google: users.google_sub, users.google_linked_at.
- Verify Google ID token phía server.

Frontend đã có:
- Nút Google trên trang login.
- Đăng ký parent hiện tại sử dụng parent code và email/password.

## 4) Đề xuất nghiệp vụ

Đăng ký Google cho parent sẽ theo quy trình:
1. User nhập parent code trên trang register.
2. User bấm Đăng ký với Google.
3. Frontend lấy id_token từ Google GIS.
4. Backend verify id_token và kiểm tra email_verified.
5. Backend verify parent code (chưa hết hạn, chưa vượt max usage).
6. Backend kiểm tra email chưa tồn tại.
7. Backend tạo user active, gán role PARENT, tạo parent record, gán student-parent.
8. Backend lưu google_sub vào users.
9. Backend tăng usage_count atomically và trả JWT.

## 5) API đề xuất mới

Endpoint mới:
- POST /api/v1/register/parent/google

Request:
- id_token: string
- parent_code: string

Response thành công:
- access_token
- token_type
- expires_in

Lỗi chính:
- invalid parent code
- parent code expired
- parent code max usage reached
- email already exists
- invalid google token
- google email not verified

## 6) Vùng ảnh hưởng cần sửa

Backend:
- apps/api/internal/http/router.go
  - Thêm route public register parent google.

- apps/api/internal/api/v1/handlers/parent_code_handler.go
  - Thêm request struct RegisterParentWithGoogleRequest.
  - Thêm handler RegisterParentWithGoogle.

- apps/api/internal/service/parent_code_service.go
  - Thêm hàm RegisterParentWithGoogle.
  - Bổ sung dependency GoogleTokenVerifier vào service constructor.

- apps/api/cmd/api/main.go
  - Truyền googleVerifier vào ParentCodeService như AuthService.

- apps/api/internal/repo/user_repo.go
  - Tái sử dụng FindByEmail và LinkGoogleSub.

Frontend:
- apps/web/src/lib/api/auth.api.ts
  - Thêm method registerParentWithGoogle.

- apps/web/src/types/index.ts
  - Thêm type RegisterParentWithGoogleRequest.

- apps/web/src/app/(auth)/register/page.tsx
  - Thêm section Đăng ký bằng phương thức khác.
  - Nút Google đăng ký, bắt buộc có parent_code trước khi submit.

## 7) Nguyên tắc bảo mật

- Chỉ tin id_token đã verify server-side.
- Bắt buộc email_verified = true.
- Không tạo tài khoản nếu email đã tồn tại.
- Không cho phép bypass parent code.
- Không đưa token nhạy cảm vào URL query.

## 8) Kế hoạch triển khai đề xuất

Bước 1: Backend API và service
- Tạo endpoint register parent google.
- Implement service flow tạo parent từ Google identity + parent code.

Bước 2: Frontend register
- Bổ sung nút Đăng ký với Google ở trang register.
- Hiển thị lỗi nghiệp vụ rõ ràng cho parent code.

Bước 3: Kiểm thử
- go test ./... ở apps/api.
- npx tsc --noEmit ở apps/web.
- Smoke test các nhánh thành công và lỗi.

Bước 4: Rollout
- Bật flag theo môi trường.
- Theo dõi log đăng ký, conflict email, lỗi verify token.

## 9) Tiêu chí hoàn tất

Hoàn tất khi:
- Parent có thể đăng ký bằng Google + parent code thành công.
- Luồng register email/password hiện tại vẫn hoạt động bình thường.
- Luồng teacher activation không bị ảnh hưởng.
- Build backend/frontend xanh và smoke test pass.

## 10) Rủi ro và cách giảm thiểu

Rủi ro:
- User đã có email local nhưng chưa link Google.
- Parent code race condition khi usage gần max.
- Sai cấu hình Google client ID giữa frontend và backend.

Giảm thiểu:
- Trả error rõ ràng cho email conflict.
- Dùng IncrementUsageIfNotMaxed để tránh race.
- Kiểm tra env và startup validation theo từng môi trường.
