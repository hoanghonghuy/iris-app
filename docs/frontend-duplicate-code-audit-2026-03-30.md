# Frontend Duplicate Code Audit (Re-Verification)

- Date: 2026-03-30
- Scope: `apps/web/src` (pages, hooks, API client utilities)
- Goal: Re-verify duplication status and list all confirmed duplicate hotspots with concrete references.

## Executive Summary

Frontend still contains meaningful duplication. The largest maintenance risks are:

1. Repeated auth error helper across auth pages.
2. Repeated role-chat wrappers.
3. Repeated class/student loading orchestration across admin and teacher hooks (with school bootstrap mainly in admin hooks).
4. Repeated admin desktop/mobile rendering structures.
5. Repeated URL query-building snippets.
6. Repeated modal close/reset state patterns.

## Method

- Pattern scan over frontend files for exact/near-duplicate signatures.
- Direct file read to confirm each cluster is real duplication (not false positives).
- Group findings by impact on maintainability and risk of inconsistent fixes.

## Confirmed Duplication Clusters

### 1) Auth helper duplication: `extractErrorMessage` repeated 5 times (High)

Same helper logic is declared in multiple auth pages.

- `apps/web/src/app/(auth)/activate/page.tsx:16`
- `apps/web/src/app/(auth)/reset-password/page.tsx:19`
- `apps/web/src/app/(auth)/login/page.tsx:25`
- `apps/web/src/app/(auth)/register/page.tsx:20`
- `apps/web/src/app/(auth)/forgot-password/page.tsx:18`

Risk:
- Future changes to auth error parsing can drift between pages.

### 2) Role chat wrapper duplication: identical page wrappers (High)

The 3 role chat pages are effectively copy-paste wrappers around `ChatPage` with identical container classes.

- `apps/web/src/app/admin/chat/page.tsx:9`
- `apps/web/src/app/teacher/chat/page.tsx:9`
- `apps/web/src/app/parent/chat/page.tsx:9`

Risk:
- Layout or behavior tweaks require touching 3 files for one concern.

### 3) Auth container wrapper class duplication (Medium)

Repeated auth page container class string:

`flex w-full items-center justify-center w-full max-w-screen-xl flex justify-center`

Matches:

- `apps/web/src/app/(auth)/forgot-password/page.tsx:48`
- `apps/web/src/app/(auth)/register/page.tsx:93`
- `apps/web/src/app/(auth)/register/page.tsx:109`
- `apps/web/src/app/(auth)/activate/page.tsx:49`
- `apps/web/src/app/(auth)/activate/page.tsx:65`

Risk:
- Styling/layout adjustments are scattered and easy to miss.

### 4) Repeated load orchestration across hooks (High)

Highly similar effect/data-loading orchestration is repeated across multiple domains, but with two variants:

- Admin hooks commonly use `loadSchools` + `loadClasses` (and in some cases `loadStudents`).
- Teacher hooks commonly use `loadClasses` and/or `loadStudents` (without the full admin-style school bootstrap chain).

- `apps/web/src/app/admin/parents/useAdminParentsPage.ts:65`
- `apps/web/src/app/admin/parents/useAdminParentsPage.ts:85`
- `apps/web/src/app/admin/parents/useAdminParentsPage.ts:110`
- `apps/web/src/app/admin/students/useAdminStudentsPage.ts:47`
- `apps/web/src/app/admin/students/useAdminStudentsPage.ts:70`
- `apps/web/src/app/admin/teachers/useAdminTeachersPage.ts:69`
- `apps/web/src/app/admin/teachers/useAdminTeachersPage.ts:89`

Teacher-side related references (class/student loading orchestration):

- `apps/web/src/app/teacher/posts/useTeacherPostsPage.ts:45`
- `apps/web/src/app/teacher/posts/useTeacherPostsPage.ts:69`
- `apps/web/src/app/teacher/health/useTeacherHealthPage.ts:39`
- `apps/web/src/app/teacher/attendance/useTeacherAttendancePage.ts:160`

Risk:
- Bug fixes and edge-case handling become inconsistent across pages.

### 5) Repeated list-fetch callback pattern in admin hooks (Medium)

Same callback skeleton (`loading/error/fetch/set/finally`) is repeated.

- `apps/web/src/app/admin/teachers/useAdminTeachersPage.ts:46` (`fetchTeachers`)
- `apps/web/src/app/admin/parents/useAdminParentsPage.ts:44` (`fetchParents`)
- `apps/web/src/app/admin/school-admins/useAdminSchoolAdminsPage.ts:34` (`fetchAdmins`)
- `apps/web/src/app/admin/students/useAdminStudentsPage.ts:90` (`fetchStudents`)

Risk:
- Error handling and UX behavior can drift by page.

### 6) URL query-builder duplication in teacher API (Medium)

Repeated `URLSearchParams` range query build logic.

- `apps/web/src/lib/api/teacher.api.ts:65`
- `apps/web/src/lib/api/teacher.api.ts:78`
- `apps/web/src/lib/api/teacher.api.ts:119`

Risk:
- Inconsistent query behavior when extending date filters.

### 7) Modal/alert close-reset pattern duplication (Low-Medium)

Repeated reset state object shape on modal close.

- `apps/web/src/app/admin/parents/page.tsx:197`
- `apps/web/src/app/admin/parents/page.tsx:240`
- `apps/web/src/app/admin/teachers/page.tsx:237`
- `apps/web/src/app/admin/teachers/page.tsx:273`
- `apps/web/src/app/admin/students/page.tsx:284`
- `apps/web/src/app/admin/school-admins/page.tsx:164`
- `apps/web/src/app/admin/users/page.tsx:228`

Risk:
- Verbose handlers and repeated object literals increase maintenance noise.

### 8) Admin desktop/mobile render structure duplication (Medium)

Repeated `hidden md:block` / `md:hidden` bifurcation with very similar table/card composition.

Representative references:

- `apps/web/src/app/admin/teachers/page.tsx:128`
- `apps/web/src/app/admin/teachers/page.tsx:184`
- `apps/web/src/app/admin/parents/page.tsx:88`
- `apps/web/src/app/admin/parents/page.tsx:144`
- `apps/web/src/app/admin/students/page.tsx:188`
- `apps/web/src/app/admin/students/page.tsx:239`
- `apps/web/src/app/admin/school-admins/page.tsx:107`
- `apps/web/src/app/admin/school-admins/page.tsx:137`
- `apps/web/src/app/admin/classes/page.tsx:150`
- `apps/web/src/app/admin/classes/page.tsx:174`
- `apps/web/src/app/admin/schools/page.tsx:179`
- `apps/web/src/app/admin/schools/page.tsx:207`

Risk:
- UI behavior consistency and styling updates require many repeated edits.

## Severity Overview

- High:
  - Auth helper duplication
  - Role-chat wrapper duplication
  - Repeated admin/teacher data-loading orchestration (school/class/student variants)
- Medium:
  - List-fetch callback pattern
  - URL query-builder duplication
  - Desktop/mobile render duplication
- Low-Medium:
  - Modal close/reset duplication
  - Auth wrapper class duplication

## Recommended Refactor Order (Highest Impact First)

1. Extract one shared auth error parser utility for all auth pages.
2. Consolidate role chat wrappers into one shared route wrapper/layout.
3. Introduce reusable data-loading hooks for class/student loading, with an admin-only extension for school bootstrap.
4. Extract shared admin list-view primitives (desktop/mobile pair abstractions).
5. Extract common date-range query builder utility for API layer.
6. Normalize modal close/reset helper utilities to reduce repetitive handlers.

## Final Verdict

After re-verification, duplication is still present and should be treated as active technical debt. The codebase is partially improved compared to earlier state, but not yet "clean" regarding repeated logic and repeated UI structures.
