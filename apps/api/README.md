# Iris API

Backend service for the Iris School Management Platform — Go + Gin + PostgreSQL.

Base URL: `http://localhost:8080/api/v1`

## Structure

```
apps/api/
├── cmd/api/main.go              # Entry point, wire up dependencies, start server
├── internal/
│   ├── api/v1/handlers/         # HTTP handlers grouped by domain/role
│   │   ├── analytics/
│   │   ├── audit_log/
│   │   ├── auth/
│   │   ├── chat/
│   │   ├── class/
│   │   ├── parent/
│   │   ├── parent_code/
│   │   ├── parent_scope/
│   │   ├── school/
│   │   ├── school_admin/
│   │   ├── shared/
│   │   ├── student/
│   │   ├── teacher/
│   │   ├── teacher_scope/
│   │   └── user/
│   ├── auth/                    # JWT helper & Google ID token verification
│   ├── config/                  # Env-based configuration
│   ├── db/                      # PostgreSQL connection pool
│   ├── http/                    # Router setup (Gin)
│   ├── middleware/               # Auth JWT, role guard, rate limiting, audit logging
│   ├── model/                   # Domain models (structs matching DB tables)
│   ├── repo/                    # Database access layer (pgx)
│   ├── response/                # Standardized JSON response helpers
│   ├── service/                 # Business logic layer
│   └── ws/                      # WebSocket hub for real-time chat
└── migrations/                  # PostgreSQL migrations (golang-migrate)
```

## Packages

| Package | Responsibility |
|---------|---------------|
| `api/v1/handlers/` | HTTP handlers — parse request, call service, return response |
| `auth/` | JWT generation/validation, Google ID token verification |
| `config/` | Load & validate `.env` configuration |
| `db/` | PostgreSQL connection pool with `pgx/v5` |
| `http/` | Router setup — register all routes, apply middleware |
| `middleware/` | `AuthJWT` (JWT validation), `RequireRole` / `RequireAnyRole` (RBAC), `InjectAdminScope` (school-level scoping), rate limiting (auth endpoints), `AuditLogger` (auto-log all protected requests) |
| `model/` | Domain structs: `User`, `Student`, `Teacher`, `Class`, `AttendanceRecord`, `HealthLog`, `Post`, `Appointment`, `AuditLog`, `Chat` (Conversation, Message), etc. |
| `repo/` | SQL queries via `pgx` — one repo per domain |
| `response/` | `OK()`, `Created()`, `Fail()`, `FailWithCode()`, `OKPaginated()` — standardized JSON envelopes |
| `service/` | Business logic — validation, orchestration, calls repo(s) |
| `ws/` | WebSocket `Hub` — manage client connections, broadcast messages |

## Role System

| Role | Scope |
|------|-------|
| `SUPER_ADMIN` | Full system access — all schools, user management, role assignment, audit logs |
| `SCHOOL_ADMIN` | Scoped to assigned school — manage classes, students, teachers, parents |
| `TEACHER` | Scoped to assigned classes — attendance, health logs, posts, appointments |
| `PARENT` | Scoped to linked children — view feed, health, appointments, book slots |

## API Endpoints

All endpoints are prefixed with `/api/v1`.

### Public

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Service health check |
| `POST` | `/auth/login` | Email/password login (rate-limited) |
| `POST` | `/auth/login/google` | Google OAuth login (rate-limited) |
| `POST` | `/auth/refresh` | Rotate refresh token and issue a new token pair |
| `POST` | `/auth/forgot-password` | Request password reset email (rate-limited) |
| `POST` | `/auth/reset-password` | Reset password with token (rate-limited) |
| `POST` | `/users/activate-token` | Activate teacher account via token |
| `GET` | `/register/parent/verify` | Verify parent registration code |
| `POST` | `/register/parent` | Register parent account with code |
| `POST` | `/register/parent/google` | Register parent via Google + code |
| `GET` | `/chat/ws` | WebSocket upgrade for real-time messaging |

### Protected — All Authenticated Users

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/me` | Current user profile (from JWT claims) |
| `PUT` | `/me/password` | Change own password |
| `DELETE` | `/me` | Delete own account |

### Chat (`/chat/*`)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/users/search` | Search users to start conversation |
| `POST` | `/conversations/direct` | Create or return existing 1-on-1 conversation (`201` if newly created, `200` if already existed) |
| `POST` | `/conversations/group` | Create group conversation (creator + `participant_user_ids`) |
| `GET` | `/conversations` | List user's conversations (sorted by latest activity; each item may include `last_message`, `unread_count`) |
| `PATCH` | `/conversations/:conversation_id/group` | Rename group (`{ "name": "..." }`, empty clears name) |
| `POST` | `/conversations/:conversation_id/participants` | Add members (`{ "user_ids": ["uuid", ...] }`) |
| `DELETE` | `/conversations/:conversation_id/participants/:user_id` | Remove a member (group must stay ≥ 2 members) |
| `POST` | `/conversations/:conversation_id/read` | Mark conversation read up to latest message (sidebar unread) |
| `GET` | `/conversations/:conversation_id/messages` | Get messages in conversation (also marks read up to latest message) |

### Teacher Scope (`/teacher/*`) — requires `TEACHER` role

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/classes` | List assigned classes |
| `GET` | `/classes/:class_id/students` | List students in class |
| `POST` | `/attendance` | Mark attendance for students |
| `DELETE` | `/attendance` | Cancel attendance record |
| `GET` | `/students/:student_id/attendance` | View attendance history |
| `GET` | `/students/:student_id/attendance-changes` | View attendance change log (audit trail) |
| `GET` | `/classes/:class_id/attendance-changes` | View class-wide change log |
| `POST` | `/health` | Create health log |
| `GET` | `/students/:student_id/health` | View health log history |
| `POST` | `/posts` | Create post |
| `PUT` | `/posts/:post_id` | Edit post |
| `DELETE` | `/posts/:post_id` | Delete post |
| `POST` | `/posts/:post_id/like` | Toggle like on post |
| `GET` | `/posts/:post_id/comments` | List comments on post |
| `POST` | `/posts/:post_id/comments` | Add comment on post |
| `POST` | `/posts/:post_id/share` | Share post |
| `GET` | `/classes/:class_id/posts` | List class-scoped posts |
| `GET` | `/students/:student_id/posts` | List student-scoped posts |
| `POST` | `/appointments/slots` | Create appointment time slot |
| `GET` | `/appointments` | List own appointments |
| `PATCH` | `/appointments/:appointment_id/status` | Update appointment status |
| `GET` | `/analytics` | Dashboard stats |
| `GET` | `/analytics/timeseries` | Dashboard chart series (`range`, `interval=day`) |
| `PUT` | `/profile` | Update own profile |

### Parent Scope (`/parent/*`) — requires `PARENT` role

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/children` | List linked children |
| `GET` | `/feed` | Aggregated feed across all children |
| `POST` | `/posts/:post_id/like` | Toggle like on post |
| `GET` | `/posts/:post_id/comments` | List comments on post |
| `POST` | `/posts/:post_id/comments` | Add comment on post |
| `POST` | `/posts/:post_id/share` | Share post |
| `GET` | `/children/:student_id/class-posts` | View child's class posts |
| `GET` | `/children/:student_id/student-posts` | View child's individual posts |
| `GET` | `/children/:student_id/posts` | View all posts related to child |
| `GET` | `/appointments/slots` | List available appointment slots |
| `POST` | `/appointments` | Book appointment |
| `GET` | `/appointments` | List own appointments |
| `PATCH` | `/appointments/:appointment_id/cancel` | Cancel appointment |
| `GET` | `/analytics` | Dashboard stats |
| `GET` | `/analytics/timeseries` | Dashboard chart series (`student_id` required, `range`, `interval=day`) |
| `PUT` | `/profile` | Update own profile |

### Admin Scope (`/admin/*`) — requires `SUPER_ADMIN` or `SCHOOL_ADMIN`

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/ping` | Admin health check |
| `GET` | `/analytics` | Dashboard stats |
| `GET` | `/analytics/timeseries` | Dashboard chart series (optional `school_id` for Super Admin, `range`, `interval=day`) |
| `GET` | `/audit-logs` | Query audit logs (**Super Admin only**) |
| `GET` | `/schools` | List schools |
| `POST` | `/schools` | Create school (**Super Admin only**) |
| `PUT` | `/schools/:school_id` | Update school (**Super Admin only**) |
| `DELETE` | `/schools/:school_id` | Delete school (**Super Admin only**) |
| `POST` | `/classes` | Create class |
| `GET` | `/classes/by-school/:school_id` | List classes by school |
| `PUT` | `/classes/:class_id` | Update class |
| `DELETE` | `/classes/:class_id` | Delete class |
| `POST` | `/students` | Create student |
| `GET` | `/students/by-class/:class_id` | List students by class |
| `GET` | `/students/:student_id` | Get student profile (includes parents) |
| `PUT` | `/students/:student_id` | Update student |
| `DELETE` | `/students/:student_id` | Delete student |
| `POST` | `/students/:student_id/generate-parent-code` | Generate parent registration code |
| `DELETE` | `/students/:student_id/parent-code` | Revoke parent code |
| `POST` | `/users` | Create user |
| `GET` | `/users` | List users |
| `GET` | `/users/:user_id` | Get user by ID |
| `POST` | `/users/:user_id/lock` | Lock user account |
| `POST` | `/users/:user_id/unlock` | Unlock user account |
| `POST` | `/users/:user_id/roles` | Assign role (**Super Admin only**) |
| `POST` | `/teachers` | Create teacher profile |
| `GET` | `/teachers` | List all teachers |
| `GET` | `/teachers/:teacher_id` | Get teacher by ID |
| `PUT` | `/teachers/:teacher_id` | Update teacher |
| `GET` | `/teachers/class/:class_id` | List teachers of a class |
| `POST` | `/teachers/:teacher_id/classes/:class_id` | Assign teacher to class |
| `DELETE` | `/teachers/:teacher_id/classes/:class_id` | Unassign teacher from class |
| `DELETE` | `/teachers/:teacher_id` | Delete teacher |
| `GET` | `/parents` | List all parents |
| `GET` | `/parents/:parent_id` | Get parent by ID |
| `PUT` | `/parents/:parent_id` | Update parent |
| `POST` | `/parents/:parent_id/students/:student_id` | Link parent to student |
| `DELETE` | `/parents/:parent_id/students/:student_id` | Unlink parent from student |
| `POST` | `/school-admins` | Create school admin (**Super Admin only**) |
| `GET` | `/school-admins` | List school admins (**Super Admin only**) |
| `DELETE` | `/school-admins/:admin_id` | Delete school admin (**Super Admin only**) |

## Migrations

Run with [golang-migrate](https://github.com/golang-migrate/migrate):

```bash
migrate -path apps/api/migrations -database "$DATABASE_URL" up
```

| # | Name | Purpose |
|---|------|---------|
| 000001 | `core` | Base schema: schools, classes, students, teachers, users, profiles |
| 000002 | `attendance_health_posts` | Attendance records, health logs, posts |
| 000003 | `activation_parent_codes` | Teacher activation tokens, student parent codes |
| 000004 | `role_system` | Role-based access control tables |
| 000005 | `password_reset_tokens` | Password reset token support |
| 000006 | `chat` | Conversations and messages for real-time chat |
| 000007 | `post_edit_history` | Audit trail for post content edits |
| 000008 | `attendance_change_logs` | Audit trail for attendance changes (create/update/delete) |
| 000009 | `attendance_change_logs_remove_attendance_fk` | Drop FK constraint to retain logs after attendance deletion |
| 000010 | `post_interactions` | Likes, comments, shares on posts |
| 000011 | `google_identity_linking` | Link Google identities to existing user accounts |
| 000012 | `appointments` | Teacher appointment slots and parent bookings |
| 000013 | `audit_logs` | Immutable activity audit trail for all protected operations |
| 000014 | `appointments_drop_slot_unique` | Relax slot uniqueness constraint |
| 000015 | `audit_logs_school_scope` | Add school-level scoping to audit logs |

## Key Design Decisions

- **Handler → Service → Repo** layering: handlers parse HTTP, services contain business logic, repos handle SQL.
- **Standardized JSON responses**: all endpoints use `response.OK()`, `response.Fail()`, etc. (`{"data": ...}` for success, `{"error": "..."}` for errors).
- **JWT-based auth**: token contains `user_id`, `email`, `roles` (array), `school_id` (for admin scoping).
- **Rate limiting**: in-memory fixed-window limiter on auth endpoints; configurable via env vars.
- **Audit logging**: middleware auto-logs every protected request to `audit_logs` table.
- **WebSocket auth**: token sent via `Sec-WebSocket-Protocol` subprotocol (not query string) for security.
- **Unit tests**: colocated with source files (`*_test.go`), using standard `testing` package.
