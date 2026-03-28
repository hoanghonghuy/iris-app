# Báo Cáo Phân Tích Staged Changes — Google Sign-In Phase 1

**Ngày:** 2026-03-27  
**Tổng số file:** 20 (3 NEW, 17 MODIFIED)

## Tổng Quan

Changeset triển khai **Google Sign-In Phase 1** cho Iris App, bao gồm:
- Backend: Google ID Token verification, login endpoint, DB linking
- Frontend: Google Sign-In button component, login page integration
- Database: Migration thêm `google_sub` + `google_linked_at`
- Config: Env variables cho Google OAuth

### Chính sách Phase 1
- Chỉ cho phép user đã tồn tại local (không auto-provision)
- Lần đầu liên kết Google → yêu cầu xác nhận mật khẩu local
- Không bật Google One Tap

---

## Phân Tích Theo Tầng

### 1. Backend — Go API

#### ✅ Điểm tốt

| File | Nhận xét |
|------|----------|
| [google_id_token.go](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/api/internal/auth/google_id_token.go) | Clean separation: [GoogleTokenVerifier](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/api/internal/auth/google_id_token.go#20-23) interface dễ mock cho testing. Claims parser an toàn, validate đầy đủ (sub, email, email_verified) |
| [auth_service.go](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/api/internal/service/auth_service.go) | Business logic rõ ràng. Refactor [buildLoginResponse()](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/api/internal/service/auth_service.go#133-164) giảm code duplication giữa [Login](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/api/internal/service/auth_service.go#55-76) và [LoginWithGoogle](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/api/internal/service/auth_service.go#77-132). Sentinel errors đúng pattern codebase |
| [auth_handler.go](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/api/internal/api/v1/handlers/auth_handler.go) | Error switch đầy đủ, HTTP status codes hợp lý. Timeout 8s cho context phù hợp |
| [user_repo.go](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/api/internal/repo/user_repo.go) | `COALESCE(google_sub, '')` xử lý NULL safety tốt. `LinkGoogleSub` có guard condition chặt chẽ (chỉ link khi chưa có hoặc trùng sub) |
| [router.go](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/api/internal/http/router.go) | Route đặt đúng nhóm public (không cần JWT) |

#### ⚠️ Vấn đề cần xem xét

**1. [main.go](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/api/cmd/api/main.go) — `log.Fatal` khi [NewGoogleIDTokenVerifier](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/api/internal/auth/google_id_token.go#29-36) fail**

```go
googleVerifier, err := auth.NewGoogleIDTokenVerifier(cfg.GoogleClientID)
if err != nil {
    log.Fatal(err)
}
```

> [!WARNING]
> Khi `GOOGLE_LOGIN_ENABLED=false` và `GOOGLE_CLIENT_ID=""`, `idtoken.NewValidator()` vẫn được gọi, phụ thuộc vào Google's discovery endpoint. Nếu server không có internet (air-gapped env hoặc CI), app sẽ **crash lúc khởi động** dù Google login bị disable.
> 
> **Đề xuất:** Chỉ khởi tạo verifier khi `cfg.GoogleLoginEnabled == true`:
> ```go
> var googleVerifier auth.GoogleTokenVerifier
> if cfg.GoogleLoginEnabled {
>     googleVerifier, err = auth.NewGoogleIDTokenVerifier(cfg.GoogleClientID)
>     if err != nil {
>         log.Fatal(err)
>     }
> }
> ```

**2. `GoogleLoginRequest` — field `password` dùng `omitempty` nhưng không có `binding`**

```go
type GoogleLoginRequest struct {
    IDToken  string `json:"id_token" binding:"required"`
    Password string `json:"password,omitempty"`
}
```

Điều này OK về functional, nhưng `omitempty` trên input request không có tác dụng thực tế (nó chỉ ảnh hưởng khi **marshal** JSON ra, không phải khi unmarshal vào). Không gây bug nhưng hơi misleading.

**3. `LinkGoogleSub` — Error string thay vì sentinel error**

```go
if res.RowsAffected() == 0 {
    return errors.New("user already linked with a different google account")
}
```

> [!NOTE]
> Codebase hiện tại dùng sentinel errors ở service layer (ví dụ `ErrGoogleLoginDisabled`, `ErrGoogleAccountNotProvisioned`). Error string ở repo layer là inline `errors.New(...)` — không thể match bằng `errors.Is()` ở handler. Tuy nhiên, business logic flow hiện tại thì case này khó xảy ra (vì service đã check `user.GoogleSub == ""` trước). **Rủi ro thấp nhưng nên cải thiện** nếu cần handle error cụ thể ở tương lai.

---

### 2. Database Migration

#### ✅ Điểm tốt

| File | Nhận xét |
|------|----------|
| [000011...up.sql](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/api/migrations/000011_google_identity_linking.up.sql) | Partial unique index `WHERE google_sub IS NOT NULL` — đúng pattern, tránh nhiều NULL violate unique |
| [000011...down.sql](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/api/migrations/000011_google_identity_linking.down.sql) | Rollback chính xác: drop index trước, drop columns sau. Dùng `IF EXISTS` an toàn |

#### ⚠️ Lưu ý

- Migration number `000011` — cần chắc chắn migration `000010` đã tồn tại và đã chạy thành công. (Dựa vào pattern tên file trong project thì OK).

---

### 3. Frontend — Next.js

#### ✅ Điểm tốt

| File | Nhận xét |
|------|----------|
| [GoogleSignInButton.tsx](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/web/src/components/auth/GoogleSignInButton.tsx) | Component tự quản lý script loading, idempotent (check `existing` script). Password-link step UX rõ ràng |
| [login/page.tsx](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/web/src/app/(auth)/login/page.tsx) | Refactor [finalizeLogin()](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/web/src/app/%28auth%29/login/page.tsx#40-46) giảm duplication. Divider "hoặc đăng nhập bằng email" UX chuẩn |
| [auth.api.ts](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/web/src/lib/api/auth.api.ts) | API function đúng pattern codebase |
| [types/index.ts](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/web/src/types/index.ts) | Type definition rõ ràng, JSDoc đầy đủ |

#### ⚠️ Vấn đề cần xem xét

**4. `GoogleSignInButton` — `onSubmitGoogle` trong dependency array của `useEffect`**

```tsx
useEffect(() => {
    // ... google.accounts.id.initialize({ callback: async (response) => {
    //   await onSubmitGoogle({ idToken: response.credential });
    // }})
    // ... google.accounts.id.renderButton(...)
}, [clearError, clientId, onSubmitGoogle, scriptReady]);
```

> [!WARNING]
> `onSubmitGoogle` là một **inline arrow function** được tạo mới mỗi lần [LoginPage](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/web/src/app/%28auth%29/login/page.tsx#33-173) render:
> ```tsx
> <GoogleSignInButton onSubmitGoogle={async ({ idToken, password }) => { ... }} />
> ```
> Điều này sẽ gây re-run `useEffect` liên tục → **re-initialize Google button và clearhtml rofui** mỗi khi parent component re-render (ví dụ khi user nhập email/password). Có thể gây **flicker UI** hoặc mất trạng thái button.
> 
> **Đề xuất:** Wrap `onSubmitGoogle` trong `useCallback` ở [LoginPage](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/web/src/app/%28auth%29/login/page.tsx#33-173), hoặc dùng `useRef` trong component con để tránh re-render.

**5. [login/page.tsx](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/web/src/app/%28auth%29/login/page.tsx) — Duplicate className**

```tsx
<div className="flex w-full items-center justify-center w-full max-w-screen-xl flex justify-center">
```

`w-full` và `flex justify-center` bị duplicate. Không gây lỗi nhưng nên cleanup (đây là bug **có sẵn**, không phải do changeset mới gây ra).

**6. Password link step detection dùng string matching**

```tsx
const showPasswordLinkStep = Boolean(
  pendingCredential && errorMessage?.toLowerCase().includes("password confirmation required")
);
```

> [!NOTE]
> Coupling chặt với message string từ backend. Nếu backend đổi message, UI sẽ break. **Đề xuất:** Backend trả error code riêng (ví dụ `GOOGLE_LINK_PASSWORD_REQUIRED`) thay vì chỉ dựa vào text message.

---

### 4. Config & Documentation

#### ✅ Hoàn chỉnh

| File | Nhận xét |
|------|----------|
| [config.go](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/api/internal/config/config.go) | 3 config fields mới, dùng `parseBoolEnv` đúng pattern |
| [.env.example](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/api/cmd/api/.env.example) | Đầy đủ, default `false` |
| [deploy.env.example](file:///c:/Users/huyhh1/workspace/personal/iris-app/infra/docker/deploy.env.example) | Thêm cả backend + frontend env vars |
| [README.md](file:///c:/Users/huyhh1/workspace/personal/iris-app/README.md) | API docs cập nhật, policy phase 1 documented rõ ràng |
| [apps/web/README.md](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/web/README.md) | Env vars documented |

### 5. Dependencies (`go.mod` / `go.sum`)

- Thêm `google.golang.org/api v0.273.0` — dependency chính cho `idtoken`
- Kéo theo nhiều transitive deps (GCP auth, OpenTelemetry, gRPC, etc.)
- Các existing dependencies được upgrade minor versions (crypto, net, sync, etc.)
- Một số deps chuyển từ `// indirect` sang `require` trực tiếp (gin, jwt, pgx, etc.) — **đây là cleanup tốt**, phản ánh đúng actual usage

> [!NOTE]
> `google.golang.org/api` pulls in ~20 transitive dependencies (gRPC, OpenTelemetry, protobuf, etc.). Đây là trade-off hợp lý cho chức năng Google ID token verification, nhưng sẽ tăng binary size đáng kể.

---

## Tổng Kết

### Verdict: ✅ Changeset chất lượng tốt, bám sát codebase patterns

| Tiêu chí | Đánh giá |
|-----------|----------|
| Kiến trúc layers (handler → service → repo) | ✅ Đúng pattern |
| Error handling (sentinel errors + switch) | ✅ Nhất quán |
| Config management | ✅ Đúng pattern `parseBoolEnv` / `os.Getenv` |
| Migration naming & rollback | ✅ Đúng convention |
| Frontend component structure | ✅ Shadcn UI + existing patterns |
| Documentation | ✅ Cập nhật đầy đủ README + env examples |
| Type safety (TS types) | ✅ Đúng pattern |

### Các item nên fix trước khi commit

| # | Mức độ | Vấn đề | File |
|---|--------|--------|------|
| 1 | 🔴 **Cao** | [NewGoogleIDTokenVerifier](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/api/internal/auth/google_id_token.go#29-36) gọi kể cả khi Google login disabled → crash khi không có internet | [main.go](file:///c:/Users/huyhh1/workspace/personal/iris-app/apps/api/cmd/api/main.go) |
| 2 | 🟡 **Trung bình** | `onSubmitGoogle` inline function gây re-render loop trong `useEffect` | `GoogleSignInButton.tsx` |
| 3 | 🟡 **Trung bình** | Password link detection dùng string matching → fragile coupling | `GoogleSignInButton.tsx` |
| 4 | 🟢 **Thấp** | `LinkGoogleSub` trả inline error thay vì sentinel error | `user_repo.go` |
| 5 | 🟢 **Thấp** | `omitempty` trên request struct không có tác dụng khi unmarshal | `auth_handler.go` |

---

## Post-review Status (2026-03-28)

Sau khi commit, các mục cần fix đã được xử lý như sau:

| # | Vấn đề | Trạng thái | Evidence |
|---|--------|-----------|---------|
| 1 | Google verifier crash khi disabled | ✅ **RESOLVED** | `main.go`: conditional init `if cfg.GoogleLoginEnabled` |
| 2 | `onSubmitGoogle` inline → re-render | ✅ **RESOLVED** | `login/page.tsx`: `handleGoogleSubmit` wrapped in `useCallback` |
| 3 | Password link string matching | ✅ **RESOLVED** | FE dùng `errorCode === "GOOGLE_LINK_PASSWORD_REQUIRED"`, BE trả `FailWithCode` |
| 4 | `LinkGoogleSub` inline error | ✅ **RESOLVED** | `repo/errors.go`: sentinel `ErrGoogleAlreadyLinkedDifferent` (2026-03-28) |
| 5 | `omitempty` vô tác dụng | ✅ **RESOLVED** | `auth_handler.go`: đã bỏ `omitempty` khỏi `Password` field |
