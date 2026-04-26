# Iris School Management Platform

![Go](https://img.shields.io/badge/Go-1.25.5-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.11-009688?logo=gin&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)
![Vue.js](https://img.shields.io/badge/Vue.js-3.5-4FC08D?logo=vue.js&logoColor=white)
![Vite](https://img.shields.io/badge/Vite-8-646CFF?logo=vite&logoColor=white)
![Pinia](https://img.shields.io/badge/Pinia-3-FFD859?logo=pinia&logoColor=0A0A0A)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)

Iris is a full-stack school management platform built for an undergraduate thesis, with a Go (Gin + PostgreSQL) backend and a Vue 3 (Vite) frontend.

The platform supports role-based workflows for super admins, school admins, teachers, and parents.

> [!IMPORTANT]
> This repository is actively developed. APIs and UI flows may evolve.

## Highlights

- Role-based auth and authorization (admin, teacher, parent)
- School domain management (schools, classes, students, teachers, parents)
- Attendance and health logs
- Posts/feed interactions (like, comment, share)
- Real-time chat via WebSocket
- Local dev stack with Docker, migrations, and smoke scripts

## Tech Stack

### Backend
- Go `1.25.5`
- Gin
- PostgreSQL + `pgx/v5`
- JWT (`golang-jwt/jwt/v5`)
- Google ID token verification (`google.golang.org/api/idtoken`)

### Frontend
- Vue `3.5`
- Vite `8`
- Vue Router `5`
- Pinia `3` (state management)
- lucide-vue-next (icons)
- Plain CSS (custom properties, no UI framework)

## Repository Structure

```text
iris-app/
├── apps/
│   ├── api/                        # Go API (cmd, internal, migrations)
│   └── web/                        # Vue 3 + Vite frontend app
├── infra/docker/                   # Docker Compose and deploy env example
├── scripts/db/                     # Seed and cleanup scripts
├── scripts/smoke/                  # API/UI smoke checks
└── docs/                           # Audit notes and design docs
```

## Google Auth Status

- Implemented (phase 1): Google login for existing users
- Optional hosted-domain restriction is supported
- Planned: parent Google sign-up flow (`/register/parent/google`)

## Prerequisites

- Go `>= 1.25`
- Node.js `>= 20.19` or `>= 22.12`
- npm
- Docker + Docker Compose
- PostgreSQL migration CLI (`migrate`)

Install `migrate`:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Quick Start (Local)

### 1) Start local PostgreSQL

```bash
cd infra/docker
docker-compose up -d
```

### 2) Configure backend environment

Create `apps/api/cmd/api/.env` from `apps/api/cmd/api/.env.example`.

Minimum required values:

```env
DATABASE_URL=postgres://postgres:iris@localhost:5433/iris_db?sslmode=disable
JWT_SECRET=replace-with-strong-secret
PORT=8080
```

Optional values:

```env
DB_MAX_CONNS=50
JWT_TTL_MINUTES=1440
ALLOWED_ORIGINS=http://localhost:5173
FRONTEND_URL=http://localhost:5173

# Google login (phase 1)
GOOGLE_LOGIN_ENABLED=false
GOOGLE_CLIENT_ID=
GOOGLE_HOSTED_DOMAIN=

# Optional SMTP
SMTP_HOST=
SMTP_PORT=
SMTP_USER=
SMTP_PASS=
```

Auth rate-limit values (defaults):

- `AUTH_LOGIN_RATE_LIMIT=10`: maximum 10 requests per minute for each login `IP + route` (`/auth/login`, `/auth/login/google`).
- `AUTH_FORGOT_PASSWORD_RATE_LIMIT=5`: maximum 5 requests per minute for each forgot-password `IP + route`.
- `AUTH_RESET_PASSWORD_RATE_LIMIT=5`: maximum 5 requests per minute for each reset-password `IP + route` (`/auth/reset-password`).
- `AUTH_RATE_LIMIT_WINDOW_SECONDS=60`: fixed-window duration in seconds.
- `AUTH_RATE_LIMIT_CLEANUP_EVERY=256`: number of requests between limiter map cleanup runs to avoid unbounded growth.
- `AUTH_RATE_LIMIT_STALE_TTL_MULTIPLIER=5`: stale-key TTL multiplier; effective TTL = `multiplier * window`.

### 3) Run migrations

```bash
migrate -path apps/api/migrations -database "postgres://postgres:iris@localhost:5433/iris_db?sslmode=disable" up
```

### 4) (Optional) Seed demo data

```bash
docker exec -i iris-postgres psql -U postgres -d iris_db < scripts/db/seed_demo.sql
```

### 5) Run backend

```bash
cd apps/api/cmd/api
go run main.go
```

Backend base URL: `http://localhost:8080/api/v1`

### 6) Configure frontend environment

Create `apps/web/.env` with:

```env
VITE_API_URL=http://localhost:8080/api/v1
VITE_GOOGLE_CLIENT_ID=
```

### 7) Run frontend

```bash
cd apps/web
npm install
npm run dev
```

Frontend URL: `http://localhost:5173`

## Smoke and Validation

### Type/build checks

```bash
# From repository root
go test ./...

# From apps/web
npm run lint
```

### API smoke script

```bash
powershell -File scripts/smoke/api-smoke.ps1
```

### UI smoke script

```bash
node scripts/smoke/ui-smoke.mjs
```

> [!NOTE]
> Current smoke scripts validate core auth/admin/teacher/parent flows. Google-login-specific smoke coverage is planned.

## API Summary

- Public: health, login, Google login, forgot/reset password, account activation, parent registration
- Protected: profile (`/me`), role-scoped routes (`/admin/*`, `/teacher/*`, `/parent/*`), chat (`/chat/*`, `/chat/ws`)

For detailed endpoint behavior and open issues, see:
- `docs/iris-issues-audit.md`
- `docs/remaining-issues.md`

## Security Notes

- CORS uses origin allowlist, not wildcard reflection.
- WebSocket validates origin; token is sent exclusively via Sec-WebSocket-Protocol subprotocol.
- Password reset tokens are no longer embedded in URL query strings; users enter the reset code manually.

## Development Notes

- The backend keeps unit tests colocated next to the code they cover, using Go's standard `*_test.go` convention.
- The frontend uses Vue 3 Composition API with `<script setup>` syntax throughout.
- State management is handled via Pinia stores with setup-function syntax.
- Composables are organized by feature/domain under `src/composables/` with barrel exports.
- API calls go through a centralized `httpClient.js` with automatic JWT token injection and request timeout.
- The frontend uses domain-driven route sections (`admin`, `teacher`, `parent`) with lazy-loaded views.
- Layouts (`AuthLayout`, `DashboardLayout`) wrap route groups for consistent UI structure.
- Styling uses plain CSS with custom properties (no Tailwind or UI framework).
- Linting uses ESLint + oxlint + Prettier.

---

If you are preparing a thesis demo, start with:
1. `infra/docker/docker-compose.yml`
2. `apps/api/cmd/api/.env`
3. `scripts/db/seed_demo.sql`
4. `scripts/smoke/api-smoke.ps1` and `scripts/smoke/ui-smoke.mjs`
