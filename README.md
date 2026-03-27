# Iris School Management API

A RESTful API for school management systems built with Go, designed to handle student attendance, health records, and teacher-class assignments with role-based access control.

> **Note**: This project is currently under active development as part of an undergraduate thesis. Some features may be incomplete or subject to change.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Database Setup](#database-setup)
  - [Running the Application](#running-the-application)
- [API Reference](#api-reference)
  - [Authentication](#authentication)
  - [Public Endpoints](#public-endpoints)
  - [Protected Endpoints](#protected-endpoints)
  - [Teacher Endpoints](#teacher-endpoints)
  - [Admin Endpoints](#admin-endpoints)
- [Database Schema](#database-schema)
- [Development](#development)
- [Testing](#testing)
- [License](#license)
- [Acknowledgments](#acknowledgments)

## Overview

Iris is a school management system API that provides functionality for:

- Managing schools, classes, students, teachers, and parents
- Student attendance tracking with multiple status types
- Health log management for students
- Role-based access control (Admin, Teacher, Parent)
- JWT-based authentication

## Features

- **Authentication and Authorization**
  - JWT-based authentication
  - Role-based access control (RBAC)
  - Account activation workflow

- **Attendance Management**
  - Support for multiple attendance statuses (present, absent, late, excused)
  - Optional check-in and check-out time tracking
  - Teacher-scoped access (teachers can only manage their assigned classes)

- **Health Log Management**
  - Record student health information
  - Temperature tracking
  - Symptom documentation with severity levels

- **User Management**
  - Admin-controlled user creation
  - Account locking and unlocking
  - Role assignment

## Architecture

The project follows a **Layered Architecture** pattern:

```
Handler Layer (HTTP)
    |
    v
Service Layer (Business Logic)
    |
    v
Repository Layer (Data Access)
    |
    v
Database (PostgreSQL)
```

Each layer has distinct responsibilities:

- **Handler Layer**: Request parsing, input validation, HTTP response formatting
- **Service Layer**: Business logic, authorization rules, data transformation
- **Repository Layer**: Database queries, data persistence, low-level access control

## Technology Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.25 |
| Web Framework | Gin |
| Database | PostgreSQL |
| Database Driver | pgx/v5 (with connection pooling) |
| Authentication | JWT (golang-jwt/jwt/v5) |
| Password Hashing | bcrypt |
| Migration Tool | golang-migrate |
| Configuration | Environment variables (godotenv) |

## Project Structure

```
iris-app/
├── apps/
│   └── api/
│       ├── cmd/
│       │   └── api/
│       │       └── main.go              # Application entry point
│       ├── internal/
│       │   ├── api/
│       │   │   └── v1/
│       │   │       └── handlers/        # HTTP handlers
│       │   ├── auth/                    # Authentication utilities
│       │   ├── config/                  # Configuration management
│       │   ├── db/                      # Database connection
│       │   ├── http/                    # Router setup
│       │   ├── middleware/              # HTTP middleware
│       │   ├── model/                   # Domain models
│       │   ├── repo/                    # Repository layer
│       │   ├── response/                # Response utilities
│       │   └── service/                 # Business logic layer
│       └── migrations/                  # Database migrations
├── infra/
│   └── docker/
│       └── docker-compose.yml           # Docker Compose configuration
├── scripts/
│   └── db/
│       └── seed_demo.sql                # Demo data seeding script
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Docker and Docker Compose (optional, for containerized database)
- golang-migrate CLI tool

### Installation

1. Clone the repository:

```bash
git clone https://github.com/hoanghonghuy/iris-app.git
cd iris-app
```

2. Install dependencies:

```bash
go mod download
```

3. Install golang-migrate:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Configuration

Create a `.env` file in the project root:

```env
# Server
PORT=8080

# Database
DATABASE_URL=postgres://postgres:iris@localhost:5433/iris_db?sslmode=disable

# JWT
JWT_SECRET=your-secret-key-here
JWT_TTL_MINUTES=60

# Google Sign-In (phase 1)
GOOGLE_LOGIN_ENABLED=false
GOOGLE_CLIENT_ID=your-google-web-client-id.apps.googleusercontent.com
# Optional: restrict sign-in to a Google Workspace domain
GOOGLE_HOSTED_DOMAIN=
```

### Database Setup

1. Start the PostgreSQL database using Docker:

```bash
cd infra/docker
docker-compose up -d
```

2. Run database migrations:

```bash
migrate -path apps/api/migrations \
        -database "postgres://postgres:iris@localhost:5433/iris_db?sslmode=disable" \
        up
```

3. (Optional) Seed demo data:

```bash
docker exec -i iris-postgres psql -U postgres -d iris_db < scripts/db/seed_demo.sql
```

### Running the Application

```bash
go run apps/api/cmd/api/main.go
```

The API will be available at `http://localhost:8080`.

## API Reference

### Authentication

All protected endpoints require a valid JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check |
| POST | `/api/v1/auth/login` | User login |
| POST | `/api/v1/auth/login/google` | Google login with ID token |
| POST | `/api/v1/users/activate` | Activate user account |

> Google login phase 1 policy:
> - Existing local users only (no auto-provision)
> - First-time Google linking requires password confirmation
> - Google One Tap is not enabled in this phase

### Protected Endpoints

These endpoints require authentication (valid JWT):

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/me` | Get current user info |
| PUT | `/api/v1/me/password` | Update own password |
| DELETE | `/api/v1/me` | Delete own account |

### Teacher Endpoints

These endpoints require authentication and TEACHER role:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/teacher/classes` | List assigned classes |
| GET | `/api/v1/teacher/classes/:class_id/students` | List students in class |
| POST | `/api/v1/teacher/attendance` | Mark student attendance |
| GET | `/api/v1/teacher/students/:student_id/attendance` | List student attendance history |
| POST | `/api/v1/teacher/health` | Create health log |
| GET | `/api/v1/teacher/students/:student_id/health` | List student health logs |
| PUT | `/api/v1/teacher/profile` | Update own profile |
| POST | `/api/v1/teacher/posts` | Create post (class or student scope) |
| GET | `/api/v1/teacher/classes/:class_id/posts` | List class posts |
| GET | `/api/v1/teacher/students/:student_id/posts` | List student posts |

### Parent Endpoints

These endpoints require authentication and PARENT role:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/parent/children` | List own children |
| GET | `/api/v1/parent/feed` | Get aggregated feed of all children's posts |
| GET | `/api/v1/parent/children/:student_id/class-posts` | List class posts for child |
| GET | `/api/v1/parent/children/:student_id/student-posts` | List student-specific posts for child |
| GET | `/api/v1/parent/children/:student_id/posts` | List all posts related to child |

### Admin Endpoints

These endpoints require authentication and ADMIN role:

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/admin/schools` | Create school |
| GET | `/api/v1/admin/schools` | List schools |
| POST | `/api/v1/admin/classes/school` | Create class |
| GET | `/api/v1/admin/classes/school/:school_id` | List classes by school |
| POST | `/api/v1/admin/students/student` | Create student |
| GET | `/api/v1/admin/students/student/:current_class_id` | List students by class |
| POST | `/api/v1/admin/users` | Create user |
| GET | `/api/v1/admin/users` | List users |
| GET | `/api/v1/admin/users/:userid` | Get user by ID |
| POST | `/api/v1/admin/users/:userid/lock` | Lock user account |
| POST | `/api/v1/admin/users/:userid/unlock` | Unlock user account |
| POST | `/api/v1/admin/users/:userid/roles` | Assign role to user |
| GET | `/api/v1/admin/teachers` | List teachers |
| GET | `/api/v1/admin/teachers/:teacher_id` | Get teacher by ID |
| PUT | `/api/v1/admin/teachers/:teacher_id` | Update teacher |
| GET | `/api/v1/admin/teachers/class/:class_id` | List teachers of class |
| POST | `/api/v1/admin/teachers/:teacher_id/classes/:class_id` | Assign teacher to class |
| DELETE | `/api/v1/admin/teachers/:teacher_id/classes/:class_id` | Unassign teacher from class |
| GET | `/api/v1/admin/parents` | List parents |
| GET | `/api/v1/admin/parents/:parent_id` | Get parent by ID |
| POST | `/api/v1/admin/parents/:parent_id/students/:student_id` | Assign parent to student |
| DELETE | `/api/v1/admin/parents/:parent_id/students/:student_id` | Unassign parent from student |

## Database Schema

The database consists of the following main entities:

- **users**: User accounts with authentication credentials
- **roles**: Available roles (ADMIN, TEACHER, PARENT)
- **user_roles**: User-role assignments
- **schools**: School information
- **classes**: Class information linked to schools
- **students**: Student records
- **teachers**: Teacher profiles linked to users
- **teacher_classes**: Teacher-class assignments
- **parents**: Parent profiles
- **student_parents**: Student-parent relationships
- **attendance**: Daily attendance records
- **health_logs**: Student health records

For detailed schema information, refer to the migration files in `apps/api/migrations/`.

## Development

### Code Style

This project follows standard Go conventions and uses:

- Constructor pattern with unexported fields for services and handlers
- Explicit dependency injection
- Error wrapping with context
- Table-driven tests (planned)

### Adding New Features

1. Define models in `internal/model/`
2. Create repository methods in `internal/repo/`
3. Implement business logic in `internal/service/`
4. Add HTTP handlers in `internal/api/v1/handlers/`
5. Register routes in `internal/http/router.go`
6. Update `main.go` to wire dependencies

### Database Migrations

Create new migrations using golang-migrate:

```bash
migrate create -ext sql -dir apps/api/migrations -seq <migration_name>
```

For migration guidelines, see `apps/api/migrations/README.md`.

## Testing

Testing infrastructure is planned but not yet implemented.

Planned test coverage:

- Unit tests for services
- Integration tests for repositories
- End-to-end tests for API endpoints

## License

This project is developed as part of an undergraduate thesis. License terms to be determined.

## Acknowledgments

- Gin Web Framework
- pgx PostgreSQL Driver
- golang-migrate
- golang-jwt

---

**Status**: Work in Progress