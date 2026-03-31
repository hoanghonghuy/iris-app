# Frontend Duplicate Code Audit (Re-Verification)

- Date: 2026-03-30
- Scope: `apps/web/src` (pages, hooks, API client utilities)
- Goal: Re-verify duplication status and list all confirmed duplicate hotspots with concrete references.

## Executive Summary

After the 2026-03-31 implementation passes (N1 + N2 + N3 slice), frontend duplication is reduced significantly. The largest remaining maintenance risk is residual load orchestration wrappers around shared utilities in admin/teacher hooks.

## Method

- Pattern scan over frontend files for exact/near-duplicate signatures.
- Direct file read to confirm each cluster is real duplication (not false positives).
- Group findings by impact on maintainability and risk of inconsistent fixes.

## Implementation Update (2026-03-31)

Applied refactors from this audit in current codebase:

- Resolved cluster #1 (auth helper duplication): extracted shared parser via `extractApiErrorRawMessage` in `apps/web/src/lib/api-error.ts` and removed per-page helper redeclarations in auth pages.
- Resolved cluster #2 (role chat wrapper duplication): extracted shared route wrapper `apps/web/src/components/chat/ChatRoutePage.tsx` and reused it in admin/teacher/parent chat pages.
- Resolved cluster #6 (teacher API date-range query duplication): extracted `buildDateRangeQuery` into `apps/web/src/lib/api/query.ts` and reused in `teacher.api.ts`.

Batch N1 implementation (completed):

- Introduced shared loading utilities in `apps/web/src/lib/list-loaders.ts`:
  - `fetchCollectionWithState` for fetch lifecycle (`loading/error/fetch/set/pagination/finally`).
  - `loadListWithDefaultSelection` for dependent list loading + default selection.
- Applied `fetchCollectionWithState` to all admin list callbacks in:
  - `apps/web/src/app/admin/parents/useAdminParentsPage.ts`
  - `apps/web/src/app/admin/teachers/useAdminTeachersPage.ts`
  - `apps/web/src/app/admin/students/useAdminStudentsPage.ts`
  - `apps/web/src/app/admin/school-admins/useAdminSchoolAdminsPage.ts`
- Applied `loadListWithDefaultSelection` across admin/teacher hooks in:
  - `apps/web/src/app/admin/parents/useAdminParentsPage.ts`
  - `apps/web/src/app/admin/teachers/useAdminTeachersPage.ts`
  - `apps/web/src/app/admin/students/useAdminStudentsPage.ts`
  - `apps/web/src/app/teacher/posts/useTeacherPostsPage.ts`
  - `apps/web/src/app/teacher/health/useTeacherHealthPage.ts`
  - `apps/web/src/app/teacher/attendance/useTeacherAttendancePage.ts`

Batch N2 implementation (completed):

- Resolved cluster #3 (auth container class duplication): extracted shared auth layout constants in `apps/web/src/components/auth/auth-layout.ts` and reused them in auth pages (`register`, `activate`, `forgot-password`).
- Resolved cluster #7 (modal/alert close-reset duplication): extracted reusable close handlers in admin hooks/pages and replaced inline `onClose={() => setX({ ...reset... })}` patterns.
- Updated admin close-reset flow in:
  - `apps/web/src/app/admin/parents/useAdminParentsPage.ts`
  - `apps/web/src/app/admin/parents/page.tsx`
  - `apps/web/src/app/admin/teachers/useAdminTeachersPage.ts`
  - `apps/web/src/app/admin/teachers/page.tsx`
  - `apps/web/src/app/admin/students/useAdminStudentsPage.ts`
  - `apps/web/src/app/admin/students/page.tsx`
  - `apps/web/src/app/admin/school-admins/useAdminSchoolAdminsPage.ts`
  - `apps/web/src/app/admin/school-admins/page.tsx`
  - `apps/web/src/app/admin/users/page.tsx`

Remaining active clusters after implementation pass:

- None in tracked top-8 clusters.

Measured delta after this implementation pass:

- Resolved clusters: 8/8 (#1, #2, #3, #4, #5, #6, #7, #8)
- Partially resolved clusters: 0/8
- Remaining active clusters: 0/8

Residual #4 cleanup (completed):

- Added `loadListEffect` helper in `apps/web/src/lib/list-loaders.ts` to remove repeated `const loadX = async ...; void loadX();` effect wrappers while preserving guard and reset behaviors via `enabled` and `beforeLoad`.
- Applied to admin/teacher hooks:
  - `apps/web/src/app/admin/parents/useAdminParentsPage.ts`
  - `apps/web/src/app/admin/teachers/useAdminTeachersPage.ts`
  - `apps/web/src/app/admin/students/useAdminStudentsPage.ts`
  - `apps/web/src/app/teacher/posts/useTeacherPostsPage.ts`
  - `apps/web/src/app/teacher/health/useTeacherHealthPage.ts`
  - `apps/web/src/app/teacher/attendance/useTeacherAttendancePage.ts`
- Verification signal: no remaining `const loadSchools/loadClasses/loadStudents` wrapper pattern in scoped hooks.
- Validation: `npx tsc --noEmit` pass, targeted eslint pass, diagnostics clean on touched files.

Batch N3 implementation (completed):

- Added shared responsive primitive `apps/web/src/components/shared/ResponsiveSplitView.tsx`.
- Replaced repeated desktop/mobile conditional wrappers in admin pages:
  - `apps/web/src/app/admin/parents/page.tsx`
  - `apps/web/src/app/admin/teachers/page.tsx`
  - `apps/web/src/app/admin/students/page.tsx`
  - `apps/web/src/app/admin/school-admins/page.tsx`
  - `apps/web/src/app/admin/schools/page.tsx`
  - `apps/web/src/app/admin/classes/page.tsx`
  - `apps/web/src/app/admin/users/page.tsx`
- Validation: `npx tsc --noEmit` pass and targeted eslint pass.
- Structural duplication signal: admin page-level split markers reduced from 19 to 7 (remaining occurrences are shared primitive usage parameters and one non-list text visibility case on admin landing page).

## Current Status Matrix (2026-03-31)

1. #1 Auth helper duplication: Resolved (5 duplicated helpers removed).
2. #2 Role chat wrapper duplication: Resolved (3 route wrappers unified).
3. #3 Auth container wrapper class duplication: Resolved (5 duplicated wrapper occurrences replaced by shared constants).
4. #4 Repeated load orchestration across hooks: Resolved (wrapper pattern unified via `loadListEffect`; hook-specific side effects remain explicit by design).
5. #5 Repeated list-fetch callback pattern in admin hooks: Resolved (shared lifecycle helper applied).
6. #6 URL query-builder duplication in teacher API: Resolved (3 repeated blocks unified).
7. #7 Modal/alert close-reset pattern duplication: Resolved (inline close-reset objects replaced with reusable close handlers).
8. #8 Admin desktop/mobile render structure duplication: Resolved (shared primitive adopted across target admin pages).

## Post-Resolution Notes

### Cluster #4 (Load orchestration duplication)

- Former footprint was 11 call sites across 6 hooks. Wrapper duplication is now unified via `loadListEffect`, while hook-specific side effects remain explicit by design.
- Primary files: `apps/web/src/app/admin/parents/useAdminParentsPage.ts`, `apps/web/src/app/admin/students/useAdminStudentsPage.ts`, `apps/web/src/app/admin/teachers/useAdminTeachersPage.ts`, `apps/web/src/app/teacher/posts/useTeacherPostsPage.ts`, `apps/web/src/app/teacher/health/useTeacherHealthPage.ts`, `apps/web/src/app/teacher/attendance/useTeacherAttendancePage.ts`.
- Validation completed with type/lint/diagnostics clean.

### Cluster #5 (Admin fetch callback skeleton duplication)

- Status: Resolved in Batch N1.
- Delta: 4 duplicated callback skeletons replaced with shared `fetchCollectionWithState`.
- Watchout retained: preserve role-specific guard behavior in `useAdminSchoolAdminsPage.ts` during future feature additions.

## Next Batch Plan (Impact and Expected Effect)

Top-8 tracked duplication clusters are now resolved. Next work should shift from structural dedup to regression confidence:

- Optional UI smoke re-run on admin/teacher list flows.
- Keep future feature work using `ResponsiveSplitView` and `loadListEffect` to prevent duplication reintroduction.

## Historical Duplication Snapshot (Before 2026-03-31 Implementation)

This section preserves the original duplication snapshot before the implementation pass. Current status and active scope are tracked in "Current Status Matrix", "Remaining Impact Zones", and "Next Batch Plan" above.

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

After re-verification and implementation passes (N1 + N2 + N3 + residual #4 cleanup, 2026-03-31), the codebase has reduced/cleared tracked duplication in 8/8 clusters (#1 through #8). Future focus should be guardrails and smoke checks to prevent reintroduction.
