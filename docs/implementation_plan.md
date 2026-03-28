# Implementation Plan: Tích hợp Google Sign-Up cho Phụ huynh (Phase 5)

Dựa trên tài liệu [docs/google-signup-integration-proposal.md](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/docs/google-signup-integration-proposal.md), kế hoạch này bổ sung luồng đăng ký trực tiếp bằng tài khoản Google cho phụ huynh (sử dụng Parent Code).

## Mục tiêu luồng nghiệp vụ
1. User nhập mã **Parent Code** trên form.
2. User bấm nút **"Đăng ký với Google"**.
3. Frontend gửi `id_token` và `parent_code` xuống API mới.
4. Backend xác thực Google ID Token, kiểm tra Parent Code, và xử lý đăng ký (tạo User -> gán role PARENT -> tạo Parent record -> Link Parent-Student -> Lưu link Google -> Tăng usage count).
5. Frontend nhận JWT trả về và auto-login.

---

## Các thay đổi chi tiết

### 1. API: Request & Handler (Backend)
Tạo API endpoint mới để nhận request đăng ký bằng Google.

#### [MODIFY] [apps/api/internal/api/v1/handlers/parent_code_handler.go](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/api/v1/handlers/parent_code_handler.go)
- Thêm struct `RegisterParentWithGoogleRequest`:
  ```go
  type RegisterParentWithGoogleRequest struct {
      IDToken    string `json:"id_token" binding:"required"`
      ParentCode string `json:"parent_code" binding:"required"`
  }
  ```
- Thêm method `RegisterParentWithGoogle(c *gin.Context)` xử lý request, gọi service `parentCodeService.RegisterParentWithGoogle(ctx, req.IDToken, req.ParentCode)`, và trả về JWT.

#### [MODIFY] [apps/api/internal/http/router.go](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/http/router.go)
- Ánh xạ route endpoint public mới:
  `v1.POST("/register/parent/google", parentCodeHandler.RegisterParentWithGoogle)`

### 2. Service Logic (Backend)
Nhúng logic xác thực Google ID token vào class [ParentCodeService](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/service/parent_code_service.go#18-26).

#### [MODIFY] [apps/api/internal/service/parent_code_service.go](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/service/parent_code_service.go)
- Thêm dependency `googleVerifier auth.GoogleTokenVerifier` vào struct [ParentCodeService](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/service/parent_code_service.go#18-26) và constructor [NewParentCodeService](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/service/parent_code_service.go#27-38).
- Thêm hàm `RegisterParentWithGoogle(ctx context.Context, idToken, parentCode string) (*LoginResponse, error)`.
  Logic hàm:
  1. Xác thực `idToken` bằng `googleVerifier.Verify`. Validate email đã verified (thông tin này có trong jwt claims).
  2. Parse `parentCode` để lấy và kiểm tra codeInfo giống luồng register thường (chưa expire).
  3. Kiểm tra email trong hệ thống:
     - Nếu đã tồn tại email: Trả lỗi `ErrEmailAlreadyExists`. (Policy: Không tự động link tài khoản cũ nếu luồng đi từ màn hình "Đăng ký mới". User phải về trang "Đăng nhập" để link).
  4. Bắt đầu DB Transaction (thông qua hàm orchestration ở Tầng Repo để giữ Service không phụ thuộc trực tiếp vào package `pgx` của DB).
  5. Trong Transaction:
     - Tạo User active (generate random password fallback/empty).
     - Assign ROLE "PARENT".
     - Mở rộng hàm lấy schoolID từ student.
     - Tạo Parent record.
     - Liên kết student_parents.
     - Lưu Google sub vào user.
     - Atomic increment Parent Code usage.
  6. Commit Transaction nếu không lỗi, ngược lại Rollback.
  7. Ký `jwtAuth.SignToken` và trả về `service.LoginResponse` (không trỏ trực tiếp đến type FE).

#### [NEW/MODIFY] Cơ chế Database Transaction (Tầng Repo)
- Các Repo hiện tại nhận `*pgxpool.Pool`. Ta cần 1 cách để Service có thể yêu cầu Transaction mà không rò rỉ thư viện `pgx` lên Service.
- **Giải pháp dứt điểm:** Tạo một hàm orchestration chung ngay tại tầng Repo (ví dụ: `RegisterParentTx` trong [parent_code_repo.go](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/repo/parent_code_repo.go) hoặc tạo `registry_repo.go` mới).
- Hàm này sẽ nhận toàn bộ data (user, parent, student_id, google_sub) và tự gọi `pool.Begin(ctx)`. Bên trong tx, nó thực thi chuỗi SQL tuần tự. Nếu bất kỳ lệnh nào fail, nó tự `tx.Rollback()`. Nếu thành công tất cả thì `tx.Commit()`.
- Bằng cách này, Service chỉ cần gọi 1 hàm duy nhất chặn hết exception, đảm bảo tính toàn vẹn dữ liệu (ACID) mà không cần refactor lớn các file repo khác.

#### [MODIFY] [apps/api/cmd/api/main.go](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/cmd/api/main.go)
- Cập nhật chỗ khởi tạo [NewParentCodeService](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/api/internal/service/parent_code_service.go#27-38) để inject thêm `googleVerifier`.

### 3. Cấu trúc Data Frontend (Frontend Types)
Khai báo Request Type và API Client method tương ứng.

#### [MODIFY] [apps/web/src/types/index.ts](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/web/src/types/index.ts)
- Thêm type `RegisterParentWithGoogleRequest`:
  ```typescript
  export interface RegisterParentWithGoogleRequest {
    id_token: string;
    parent_code: string;
  }
  ```

#### [MODIFY] [apps/web/src/lib/api/auth.api.ts](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/web/src/lib/api/auth.api.ts)
- Thêm method `registerParentWithGoogle(data: RegisterParentWithGoogleRequest)` gọi endpoint mới POST `/api/v1/register/parent/google`.

### 4. Giao diện (Frontend Component)
Thêm nút Google vào màn hình `/register`, yêu cầu có Parent Code trước khi gọi.

#### [MODIFY] [apps/web/src/app/(auth)/register/page.tsx](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/web/src/app/%28auth%29/register/page.tsx)
- Import component [GoogleSignInButton](file:///c:/Users/yuhh/data/workspace/project/personal/web-apps/iris-app/apps/web/src/components/auth/GoogleSignInButton.tsx#39-168) (như bên trang `/login`).
- Tạo hàm `onGoogleSignUp` để handle click bắt buộc kiểm tra state `parentCode` trước khi xử lý (nếu chưa có thì báo lỗi UI).
- Khi Submit thành công nhận được JWT:
  - Gọi API `authApi.getMe()` bằng token vừa nhận để lấy profile & role.
  - Sau đó gọi tiếp `useAuth().login(token, primaryRole)` tương tự luồng `finalizeLogin` bên trang Login, qua đó tự động redirect sang giao diện Parent Dashboard.
- Bắt và map các mã lỗi (error_code/status) tương tự flow login để show message rõ ràng (VD: Email đã được sử dụng -> báo user quay lại trang đăng nhập).

### 5. Error Contract Matrix (Backend -> Frontend)
- Google Login bị tắt (`googleEnabled = false`): `403 Forbidden` (ErrGoogleLoginDisabled) -> Màn hình báo "Tính năng đăng nhập Google đang tắt".
- Hosted Domain không cho phép: `403 Forbidden` (ErrGoogleDomainNotAllowed) -> Màn hình đăng ký báo "Tên miền Google không hợp lệ".
- Token Google không hợp lệ: `401 Unauthorized` (auth.ErrInvalidCredentials) -> Màn hình đăng ký báo "Xác thực Google thất bại".
- Mã Parent Code sai/hết hạn/max usage: `400 Bad Request` -> Màn hình đăng ký báo lỗi "Mã phụ huynh không hợp lệ/hết hạn/hết lượt dùng".
- Email đã tồn tại (conflict): `409 Conflict` (ErrEmailAlreadyExists) -> Màn hình báo "Email này đã được đăng ký. Vui lòng quay lại trang Đăng nhập".

### 6. Verification Plan
- Chạy `go vet` và `go test` ở API.
- Chạy `tsc` ở Frontend.
- Test thủ công luồng đăng ký bằng Google:
  1. Tạo Parent Code mới từ Admin panel.
  2. Dùng Parent Code đó để đăng ký tài khoản Parent bằng tính năng "Register with Google" ở Frontend.
  3. Đăng nhập thành công và verify link giữa student-parent được tạo chuẩn.
  4. Test trường hợp fail: Google account bị trùng email, Parent Code sai/hết hạn.
