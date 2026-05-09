# Changelog

All notable changes to this project are documented in this file.

## [Unreleased]
- No changes yet.

## [v0.4.0] - 2026-05-10

### Features
- **Frontend Migration:** Migrated web application from React/Next.js to Vue.js with improved component architecture and composables pattern.
- **Group Chat:** Added group conversation creation with participant validation, group management endpoints (rename, add/remove members), and dedicated UI in ChatPage.
- **Chat Enhancements:** Implemented unread message tracking with `last_read_at` timestamps, last message preview in conversation list, and mark-as-read functionality.
- **Authentication:** Implemented refresh token flow with opaque tokens, automatic rotation, and secure token storage.
- **Deployment:** Added Docker Compose configuration for full-stack local deployment (PostgreSQL, API, web frontend) with comprehensive environment variable documentation.
- **Design System:** Redesigned UI with new typography (Fredoka and Nunito fonts), updated color palette, improved accessibility, and light/dark mode support.
- **Audit Logging:** Added audit logging middleware for write operations with school-scoped coverage.
- **Admin Features:** Enhanced admin flows for creating teacher profiles from existing users, improved people management components with reusable modals.
- **Parent Features:** Added parent profile update functionality with validation, enhanced feed filtering by mode and child selection.
- **Database:** Expanded seed scripts with comprehensive demo data across schools, classes, users, attendance, health logs, posts, appointments, and chat.

### Fixes
- Fixed chat unread badge incorrectly resetting during message pagination by only marking read on initial fetch.
- Fixed group chat metadata leak when actor removes themselves from group conversation.
- Fixed partial success issue in `AddConversationParticipants` by wrapping operations in transaction.
- Fixed inconsistent response normalization in `teacherPostService` update/delete methods.
- Fixed API healthcheck by switching from `wget` to `curl` and installing curl in runtime image.
- Fixed CI backend build command to align with monorepo structure and Dockerfile.
- Fixed role checking in `CreateUserWithoutPassword` to use `RolesOfUser` for existing user validation.

### Refactoring
- Extracted composables for better code reusability across admin, teacher, parent, and chat features.
- Centralized auth token management with dedicated `tokenStorage` helper.
- Extracted post interaction logic into reusable `usePostInteractions` composable.
- Moved `useParentFeedPage` to `composables/parent/` for better organization.
- Extracted admin people form configuration and shared management components.
- Centralized date helper functions and query parameter builders.

### Infrastructure
- Added CI workflow triggers for `.env.example` and `docker-compose.yml` changes.
- Updated frontend CI to build Vue app instead of Next.js.
- Enhanced security scan workflow to handle base-head equality safely.

### Breaking Changes
- None identified from v0.3.0..HEAD. Backend API contracts remain backward-compatible.

### Migration Notes
- **Frontend:** If you have custom frontend integrations, note the migration from React/Next.js to Vue.js. API contracts are unchanged.
- **Database:** Run migration `000016_refresh_tokens.up.sql` and `000017_chat_participant_last_read.up.sql` before deploying.
- **Environment:** Review updated `.env.example` files for new variables (`REFRESH_TOKEN_SECRET`, `REFRESH_TOKEN_EXPIRY`, optional `VITE_WS_URL`).
- **Docker:** Use `docker compose up -d` to start the full stack locally. Existing PostgreSQL data is preserved via volume `docker_iris_pg_data`.

## [v0.3.0] - 2026-04-16

### Features
- Added comprehensive dashboard analytics for admin, teacher, and parent roles with real-time metrics.
- Expanded appointments with stronger slot handling and cancellation-window behavior.
- Introduced and enhanced audit logging capabilities, including school-scoped coverage and super-admin-only access policy.
- Added pagination and filtering controls for admin audit logs.
- Added verification route support for parent registration flow.
- Improved chat permission checks and participant retrieval by conversation IDs.
- Added CI workflows for Docker build, backend testing, frontend testing, and security scanning.

### Fixes
- Fixed Google Sign-In iframe visual artifacts by excluding iframes from global border styling.
- Fixed CI configuration issues around Go version handling, lint tool installation, and vulnerability scan paths.

### Breaking Changes
- None identified from v0.2.0..HEAD.

## [v0.2.0] - 2026-04-03

### Features
- Strengthened authentication security with configurable fixed-window rate limiting and reset-password rate-limit flow, plus test coverage for middleware/config behavior.
- Added governance rule for privileged role assignment by restricting SUPER_ADMIN promotion path.
- Expanded admin management capabilities with update/delete flows for classes, schools, students, and teachers.
- Improved school admin creation flow with optional full name support and user search enhancements.
- Introduced reusable admin table components and UI consistency improvements across dashboards and form controls.
- Upgraded auth UX with stronger client-side validation, error handling, refined layout/branding, and enhanced Google sign-in behavior including one-tap support.
- Added runtime hardening and observability improvements, including auth-route rate-limit middleware usage and WebSocket development-only logging.
- Improved local/demo operability with richer seed data and updated environment templates/README guidance.

### Fixes
- Added stricter validation for pagination and password-related request rules in API flows.
- Improved login page accessibility semantics and Google sign-in usability details.

### Breaking Changes
- None identified from v0.1.0..HEAD.

## [v0.1.0] - 2026-04-01

Initial release.

### Features
- Backend API with JWT auth and role-based access (admin, teacher, parent).
- Core school flows: classes, students, attendance, health logs, and posts.
- Basic chat and parent-facing feed/dashboard capabilities.
- Next.js web app for admin/teacher/parent workflows.

## Link References
[Unreleased]: https://github.com/hoanghonghuy/iris-app/compare/v0.4.0...HEAD
[v0.4.0]: https://github.com/hoanghonghuy/iris-app/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/hoanghonghuy/iris-app/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/hoanghonghuy/iris-app/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/hoanghonghuy/iris-app/releases/tag/v0.1.0
