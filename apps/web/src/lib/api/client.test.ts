import { afterEach, describe, expect, it } from 'vitest';

import { authHelpers } from '@/lib/api/client';

describe('authHelpers', () => {
  afterEach(() => {
    localStorage.clear();
    sessionStorage.clear();
  });

  it('sets and gets auth token', () => {
    authHelpers.setToken('token-123');
    expect(authHelpers.getToken()).toBe('token-123');
    expect(authHelpers.isAuthenticated()).toBe(true);
  });

  it('sets and gets user role', () => {
    authHelpers.setUserRole('TEACHER');
    expect(authHelpers.getUserRole()).toBe('TEACHER');
  });

  it('removes auth token and role', () => {
    authHelpers.setToken('token-123');
    authHelpers.setUserRole('PARENT');

    authHelpers.removeToken();

    expect(authHelpers.getToken()).toBeNull();
    expect(authHelpers.getUserRole()).toBeNull();
    expect(authHelpers.isAuthenticated()).toBe(false);
  });
});
