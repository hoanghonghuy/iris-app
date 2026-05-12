# Frontend Migration Baseline Snapshot
**Created:** 2026-05-11 22:37 (Asia/Bangkok)  
**Purpose:** Đóng băng hành vi hiện tại trước khi migrate JS → TypeScript

## Phạm vi migrate
- **Location:** `apps/web/src`
- **JS files:** 61 file
- **Vue SFC:** 55 file (100% dùng `<script setup>`, 0% dùng `lang="ts"`)
- **TS files hiện tại:** 0

## Stack hiện tại
- Vue 3.5.32
- Vite 8
- Pinia 3.0.4
- Vue Router 5.0.6
- ESLint + Oxlint + Prettier
- Chart.js + vue-chartjs

## Cấu trúc thư mục nguồn
```
apps/web/src/
├── assets/          # CSS, images
├── components/      # Vue components (admin, charts, common)
├── composables/     # Composables theo domain (admin, teacher, parent, chat)
├── helpers/         # Utility functions (auth, chat, error, format, validation)
├── layouts/         # Layout components
├── router/          # Router config + guards + routes
├── services/        # API services (httpClient, authService, adminService, etc.)
├── stores/          # Pinia stores (authStore)
├── views/           # Page components theo role
├── App.vue
└── main.js
```

## Chức năng chính theo role (baseline behavior)

### Auth Flow
- **Login:** POST `/api/v1/auth/login` → lưu access_token + refresh_token → GET `/api/v1/me` → redirect theo role
- **Register Parent:** POST `/api/v1/parent-code/register` với parent_code
- **Refresh Token:** 401 trigger → POST `/api/v1/auth/refresh` → retry request
- **Logout:** Clear tokens + redirect `/login`

### Admin Role
- **Dashboard:** Snapshot analytics + timeseries charts (7d/14d/30d)
- **CRUD:** Schools, Classes, Students, Teachers, Parents, School Admins
- **Attendance:** View/edit attendance records
- **Posts:** Manage posts + comments
- **Appointments:** View/manage appointment slots
- **Audit Logs:** View system audit logs (friendly/raw mode)

### Teacher Role
- **Dashboard:** Snapshot analytics + timeseries charts (7d/14d/30d)
- **Attendance Taking:** Bulk attendance entry với status (present/absent/late/excused) + health notes
- **Classes:** View assigned classes + students
- **Posts:** Create/edit/delete posts với attachments
- **Appointments:** View appointment schedule
- **Chat:** Real-time chat với parents (WebSocket)

### Parent Role
- **Dashboard:** Snapshot analytics + timeseries charts (7d/14d/30d) cho từng con
- **Student Selector:** Switch giữa các con
- **Attendance:** View attendance history của con
- **Health Logs:** View health records
- **Posts/Newsfeed:** View posts từ teachers, comment, interact
- **Appointments:** Book/cancel appointments với teachers
- **Chat:** Real-time chat với teachers (WebSocket)
- **Profile:** Update phone number

## API Integration Pattern
```
View/Component
  ↓ (import)
Composable (useXxx)
  ↓ (call)
Service (xxxService)
  ↓ (call)
httpClient
  ↓ (HTTP)
Backend API
```

## Điểm quan trọng cần giữ nguyên

### 1. Router Guards
- `requireAuth`: Check token, redirect `/login` nếu chưa auth
- `requireRole`: Check user role, redirect `/unauthorized` nếu không đủ quyền
- Route meta typing: `{ requiresAuth: boolean, roles?: string[] }`

### 2. HTTP Client Features
- Auto attach `Authorization: Bearer <token>`
- 401 handling → refresh token flow (single-flight)
- Timeout 30s
- Query params serialization
- Error normalization

### 3. WebSocket Chat
- Connect: `wss://...` với token query param
- Events: `message`, `conversation_created`, `conversation_updated`
- Reconnect logic khi disconnect
- Message history load qua REST, realtime qua WS

### 4. State Management
- Auth state: `authStore` (user, token, isAuthenticated)
- Local state: `ref/computed` trong composables
- No global event bus

### 5. Error Handling
- `extractErrorMessage()` normalize API errors
- Display error qua `errorMessage` ref trong composables
- Toast/alert cho user feedback

### 6. Form Validation
- Client-side validation trước submit
- Server error display sau submit fail
- Loading state trong submit

## Test Checklist (manual smoke test sau migrate)

### Auth
- [ ] Login với admin/teacher/parent credentials
- [ ] Logout và verify redirect
- [ ] Register parent với parent_code hợp lệ
- [ ] Token refresh tự động khi 401

### Admin Dashboard
- [ ] View snapshot metrics
- [ ] Switch timeseries range (7d/14d/30d)
- [ ] Charts render đúng data
- [ ] Tab switch "Tổng quan" ↔ "Biểu đồ"

### Teacher Dashboard
- [ ] View snapshot metrics
- [ ] Timeseries charts hiển thị
- [ ] Attendance taking: select class → mark students → save
- [ ] Create post với attachment
- [ ] Chat với parent: send/receive messages

### Parent Dashboard
- [ ] Student selector hoạt động
- [ ] View attendance/health logs của con
- [ ] View newsfeed + comment
- [ ] Book appointment
- [ ] Chat với teacher
- [ ] Update profile phone

### Responsive
- [ ] Mobile view (≤768px): hamburger menu, stacked layout
- [ ] Tablet view (768-1024px): adaptive grid
- [ ] Desktop view (>1024px): full sidebar

## File quan trọng cần ưu tiên test sau migrate

### Core Infrastructure (rủi ro cao nếu lỗi)
1. `services/httpClient.js` - HTTP client + refresh token
2. `services/authService.js` - Auth API calls
3. `stores/authStore.js` - Auth state
4. `router/index.js` + `router/guards.js` - Routing + guards
5. `helpers/errorHandler.js` - Error normalization

### High-Traffic Composables
6. `composables/teacher/useAttendanceTaking.js` - Bulk attendance
7. `composables/chat/useChatWebSocket.js` - WebSocket chat
8. `composables/admin/useAdminCrudList.js` - Admin CRUD operations

### Critical Views
9. `views/auth/LoginPage.vue` - Entry point
10. `views/admin/AdminDashboard.vue` - Admin landing
11. `views/teacher/TeacherDashboard.vue` - Teacher landing
12. `views/parent/ParentDashboard.vue` - Parent landing

## Convention Code (phải bám sát)

### Naming
- Components/Views: `PascalCase.vue`
- Composables: `useXxx.js` → `useXxx.ts`
- Services: `xxxService.js` → `xxxService.ts`
- Stores: `xxxStore.js` → `xxxStore.ts`

### Import Alias
- `@/` → `apps/web/src/`
- Relative imports cho module gần nhau

### Code Style
- No semicolons (`semi: false`)
- Single quotes (`singleQuote: true`)
- Print width 100 (`printWidth: 100`)
- 2 spaces indent

### Vue SFC Pattern
```vue
<script setup>
// imports
// props/emits
// composables
// refs/computed
// methods
// lifecycle hooks
</script>

<template>
  <!-- markup -->
</template>

<style scoped>
  /* styles */
</style>
```

## Build/Lint Commands (baseline pass)
```bash
cd apps/web
npm run dev          # Vite dev server
npm run build        # Production build
npm run preview      # Preview build
npm run lint         # ESLint + Oxlint
npm run format       # Prettier
```

## Baseline Verification Status
- [x] Lint pass: `npm run lint` clean
- [x] Build pass: `npm run build` success
- [x] Dev server: `npm run dev` starts without errors
- [x] Manual smoke test: Admin/Teacher/Parent flows functional

## Notes
- Ledger line 84 ghi nhận FE migration batches V1-V3 đã hoàn thành trước đây với `<script setup lang="ts">`, nhưng audit hiện tại cho thấy **0 file Vue có `lang="ts"`** → có thể đã rollback hoặc ledger cũ không chính xác. Baseline này phản ánh trạng thái thực tế hiện tại: **100% JS, 0% TS**.
