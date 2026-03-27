# Iris School Management Platform

![Go](https://img.shields.io/badge/Go-1.25.5-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.11-009688?logo=gin&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)
![Next.js](https://img.shields.io/badge/Next.js-16-000000?logo=next.js&logoColor=white)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=0A0A0A)
![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6?logo=typescript&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-4-06B6D4?logo=tailwindcss&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)

Iris is a full-stack school management platform built for an undergraduate thesis, with a Go (Gin + PostgreSQL) backend and a Next.js frontend.

It supports role-based workflows for super admins, school admins, teachers, and parents, including attendance, health logs, posts/feed interactions, chat, and account management.

> [!IMPORTANT]
> This repository is actively developed. APIs, UI flows, and deployment defaults may evolve.

## What Is Included

- Backend API (`apps/api`)
  - Layered architecture: handler -> service -> repository
  - PostgreSQL persistence with SQL migrations
  - JWT authentication and role-based authorization
  - Google sign-in (phase 1) for existing accounts
  - WebSocket chat with origin checks and subprotocol auth
- Frontend Web App (`apps/web`)
  - Next.js App Router + React 19 + TypeScript
  - Role-based dashboards and guarded routes
  - Teacher/parent post interactions (like/comment/share)
  - Attendance and admin management screens
- Infrastructure and scripts
  - Local PostgreSQL with Docker Compose
  - DB migration and demo seed helpers
  - API and UI smoke scripts

## Tech Stack

### Backend
- Go `1.25.5`
- Gin
- PostgreSQL + `pgx/v5`
- JWT (`golang-jwt/jwt/v5`)
- Google ID token verification (`google.golang.org/api/idtoken`)

### Frontend
- Next.js `16.x`
- React `19`
- TypeScript `5`
- Tailwind CSS `4`
- shadcn/ui + lucide-react

## Repository Structure

```text
iris-app/
├── apps/
│   ├── api/
│   │   ├── cmd/api/                 # API entrypoint
│   │   ├── internal/                # Core backend modules
│   │   └── migrations/              # SQL migrations
│   └── web/
│       ├── src/app/                 # Next.js routes
│       ├── src/components/          # UI components
│       ├── src/lib/                 # API clients/utilities
│       └── src/hooks/               # Frontend hooks
├── infra/docker/                    # Docker compose + deploy env example
├── scripts/db/                      # DB seed/cleanup scripts
├── scripts/smoke/                   # API/UI smoke tests
└── docs/                            # Audit notes and implementation docs
```

## Core Features

- Authentication and Authorization
  - Email/password login
  - Google login endpoint: `POST /api/v1/auth/login/google`
  - Role-scoped access control for admin/teacher/parent
- School Domain Management
  - Schools, classes, students, teachers, parents
  - Teacher-class assignment and parent-student linkage
- Attendance and Health
  - Teacher attendance marking + history
  - Health log recording and listing
- Posts and Feed
  - Teacher posts (class/student scope)
  - Parent feed aggregation
  - Persisted like/comment/share interactions
- Chat
  - Conversation and messages endpoints
  - WebSocket delivery

## Current Google Auth Status

- Implemented (phase 1)
  - Google login for existing users
  - Optional hosted-domain restriction
  - First-time Google linking requires password confirmation
- Not implemented yet
  - Parent Google sign-up flow (`/register/parent/google`) is still planned (see `docs/google-signup-integration-proposal.md`)

## Prerequisites

- Go `>= 1.25`
- Node.js `>= 20`
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

Useful optional values:

```env
DB_MAX_CONNS=50
JWT_TTL_MINUTES=1440
ALLOWED_ORIGINS=http://localhost:3000
FRONTEND_URL=http://localhost:3000

# Google login (phase 1)
GOOGLE_LOGIN_ENABLED=false
GOOGLE_CLIENT_ID=
GOOGLE_HOSTED_DOMAIN=

# Compatibility mode (keep false in production)
WS_ALLOW_QUERY_TOKEN_FALLBACK=false

# Optional SMTP
SMTP_HOST=
SMTP_PORT=
SMTP_USER=
SMTP_PASS=
```

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

Create `apps/web/.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_GOOGLE_CLIENT_ID=
NEXT_PUBLIC_WS_QUERY_TOKEN_FALLBACK=false
```

### 7) Run frontend

```bash
cd apps/web
npm install
npm run dev
```

Frontend URL: `http://localhost:3000`

## Smoke and Validation

### Backend and frontend type/build checks

```bash
# From repository root
go test ./...

# From apps/web
npx tsc --noEmit
npx eslint
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

## API Surface (High Level)

- Public
  - `GET /api/v1/health`
  - `POST /api/v1/auth/login`
  - `POST /api/v1/auth/login/google`
  - `POST /api/v1/auth/forgot-password`
  - `POST /api/v1/auth/reset-password`
  - `POST /api/v1/users/activate-token`
  - `POST /api/v1/register/parent`
- Protected
  - `GET /api/v1/me`
  - `PUT /api/v1/me/password`
  - `DELETE /api/v1/me`
  - `/api/v1/teacher/*`
  - `/api/v1/parent/*`
  - `/api/v1/admin/*`
  - `/api/v1/chat/*` and `/api/v1/chat/ws`

For detailed endpoint behavior and open issues, see:
- `docs/iris-issues-audit.md`
- `docs/remaining-issues.md`

## Security Notes

- CORS uses origin allowlist, not wildcard reflection.
- WebSocket validates origin; token via subprotocol is preferred.
- Query-token WebSocket fallback exists for compatibility and should remain disabled in production.
- Password reset race-condition guards are implemented; token-in-URL hardening is still tracked in audit.

## Development Notes

- The backend is organized with explicit service/repository boundaries to keep business logic testable.
- The frontend uses domain-driven route sections (`admin`, `teacher`, `parent`) and shared typed API clients.
- Migration files are incremental and live in `apps/api/migrations`.

---

If you are preparing a thesis demo, start with:
1. `infra/docker/docker-compose.yml`
2. `apps/api/cmd/api/.env`
3. `scripts/db/seed_demo.sql`
4. `scripts/smoke/api-smoke.ps1` and `scripts/smoke/ui-smoke.mjs`
