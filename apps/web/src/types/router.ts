import type { UserRole } from './auth'

// Router meta types
export interface RouteMeta {
  requiresAuth?: boolean
  roles?: UserRole[]
  title?: string
  layout?: 'dashboard' | 'auth' | 'default'
}

// Route location types (extend vue-router types)
export interface RouteLocationNormalized {
  path: string
  name?: string | symbol
  params: Record<string, string | string[]>
  query: Record<string, string | string[]>
  hash: string
  meta: RouteMeta
  matched: any[]
}

// Navigation guard types
export type NavigationGuardNext = (to?: string | false | void | { name?: string; path?: string; replace?: boolean }) => void

export interface NavigationGuardContext {
  to: RouteLocationNormalized
  from: RouteLocationNormalized
  next: NavigationGuardNext
}
