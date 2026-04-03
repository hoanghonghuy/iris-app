# Changelog

All notable changes to this project are documented in this file.

## [Unreleased]
- No changes yet.

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
[Unreleased]: https://github.com/hoanghonghuy/iris-app/compare/v0.2.0...HEAD
[v0.2.0]: https://github.com/hoanghonghuy/iris-app/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/hoanghonghuy/iris-app/releases/tag/v0.1.0
