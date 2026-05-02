# Iris Web Frontend

Vue 3 + Vite frontend for the Iris School Management Platform.

## Tech Stack

- Vue 3.5 (Composition API with `<script setup>`)
- Vite 8
- Vue Router 5
- Pinia 3
- lucide-vue-next (icons)
- Plain CSS (custom properties)

## Architecture

```
View → Composable → Service → httpClient → API
```

## Structure

```
apps/web/src/
├── components/          # Reusable UI components
├── composables/         # Reusable business logic
│   └── usePostInteractions.js
├── helpers/
│   ├── auth/
│   │   ├── tokenStorage.js    # Centralized token management
│   │   └── index.js
│   ├── dateFormatter.js
│   └── errorHandler.js
├── layouts/
│   ├── AuthLayout.vue
│   └── DashboardLayout.vue
├── router/
│   ├── index.js         # Router setup (34 lines, down from 234)
│   ├── guards.js        # Auth, role, guest guards
│   └── routes/          # Route modules by role
│       ├── auth.js
│       ├── admin.js
│       ├── teacher.js
│       └── parent.js
├── services/
│   ├── httpClient.js    # Axios wrapper with JWT
│   ├── postService.js   # Shared post API (factory pattern)
│   ├── teacherService.js
│   ├── parentService.js
│   └── authService.js
├── stores/
│   └── authStore.js     # Auth state + user profile
└── views/               # Page components
    ├── admin/
    ├── teacher/
    ├── parent/
    └── auth/
```

## Key Patterns

### 1. Composables for Reusable Logic

Extract business logic from components:

```js
// composables/usePostInteractions.js
export function usePostInteractions(role) {
  const processing = ref(false)
  const toggleLike = async (postId) => { /* ... */ }
  return { processing, toggleLike }
}
```

### 2. Service Factory Pattern

Avoid duplicate API methods:

```js
// services/postService.js
export function createPostService(rolePrefix) {
  return {
    async togglePostLike(postId) { /* ... */ },
    async getPostComments(postId, params) { /* ... */ },
  }
}

export const teacherPostService = createPostService('teacher')
export const parentPostService = createPostService('parent')
```

### 3. Centralized Token Storage

Single source of truth:

```js
// helpers/auth/tokenStorage.js
export const tokenStorage = {
  get: () => localStorage.getItem('iris_token'),
  set: (token) => localStorage.setItem('iris_token', token),
  remove: () => localStorage.removeItem('iris_token'),
}
```

### 4. Router Modularization

Guards and routes separated by role:

```js
// router/index.js - 34 lines (was 234)
import { navigationGuard } from './guards'
import { authRoutes } from './routes/auth'
import { adminRoutes } from './routes/admin'
// ...
```

## Development

```bash
# Setup
cd apps/web
npm install

# Create .env from .env.example
VITE_API_URL=http://localhost:8080/api/v1

# Dev server
npm run dev              # http://localhost:5173

# Build
npm run build

# Lint
npm run lint
```

## Code Style

- Composition API with `<script setup>`
- Setup-function stores in Pinia
- Async/await over promises
- Named exports (except Vue components)
- JSDoc for public functions

## Related Docs

- [API Documentation](../api/README.md)
- [Frontend Structure Audit](../../docs/FRONTEND_STRUCTURE_AUDIT_2026-05-01.md)
- [Main README](../../README.md)
