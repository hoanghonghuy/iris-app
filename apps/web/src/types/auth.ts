// Auth & User types
export interface User {
  id: string
  email: string
  roles: UserRole[]
  name?: string
  phone?: string
  school_id?: string
  created_at?: string
  updated_at?: string
}

export type UserRole = 'SUPER_ADMIN' | 'SCHOOL_ADMIN' | 'TEACHER' | 'PARENT'

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
}

export interface RegisterParentRequest {
  email: string
  password: string
  parent_code: string
}

export interface GoogleLoginRequest {
  id_token: string
  password?: string
}

export interface ActivateAccountRequest {
  token: string
  password: string
}

export interface RefreshTokenRequest {
  refresh_token: string
}

export interface ChangePasswordRequest {
  old_password: string
  new_password: string
}

export interface ResetPasswordRequest {
  token: string
  new_password: string
}
