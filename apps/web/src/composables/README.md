# Composables

Thư mục chứa các Vue 3 Composition API composables được tổ chức theo feature/domain.

## Cấu trúc thư mục

```
composables/
├── admin/          # Admin-related composables
│   ├── useAdminCrudList.js           # Generic CRUD operations
│   ├── useAdminPeopleManagement.js   # People (teachers/parents) management
│   ├── useAdminUserSearch.js         # User search với pagination
│   ├── useAuditLogs.js               # Audit logs filtering & pagination
│   └── index.js                      # Barrel export
├── chat/           # Chat-related composables
│   ├── useChatWebSocket.js           # WebSocket connection & reconnect
│   ├── useChatConversations.js       # Conversations management
│   ├── useChatMessages.js            # Messages & pagination
│   ├── useChatSearch.js              # User search with debounce
│   └── index.js                      # Barrel export
├── parent/         # Parent-related composables
│   ├── useParentAppointments.js      # Appointments & booking management
│   ├── useParentAppointmentActions.js # Booking & cancellation actions
│   └── index.js                      # Barrel export
├── shared/         # Shared composables (dùng chung nhiều feature)
│   ├── useAuthForm.js                # Auth form state management
│   └── index.js                      # Barrel export
└── teacher/        # Teacher-related composables
    ├── useAppointmentsList.js        # Appointments list & filtering
    ├── useAppointmentSlotCreation.js # Slot creation & validation
    ├── useAttendanceClasses.js       # Classes loading & selection
    ├── useAttendanceHistory.js       # Student & class history
    ├── useAttendanceTaking.js        # Attendance taking & bulk operations
    ├── usePostForm.js                # Post form state & submission
    ├── useTeacherPosts.js            # Posts & data management
    └── index.js                      # Barrel export
```

## Naming Convention

- Composables: `use[Feature][Purpose].js`
- Luôn bắt đầu với prefix `use`
- Tên file mô tả rõ chức năng

## Import

Sử dụng barrel exports từ index.js:

```js
// ✅ Good
import { useChatWebSocket, useChatMessages } from '@/composables/chat'

// ❌ Avoid
import { useChatWebSocket } from '@/composables/chat/useChatWebSocket'
```
