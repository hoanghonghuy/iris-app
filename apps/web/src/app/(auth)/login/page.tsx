/**
 * Login Page
 * Trang đăng nhập chính của ứng dụng.
 * Sử dụng AuthProvider để xử lý logic đăng nhập và redirect.
 */
"use client";

import React, { useCallback, useState } from 'react';
import Link from 'next/link';
import { Loader2 } from 'lucide-react';
import { useAuth } from '@/providers/AuthProvider';
import { authApi } from '@/lib/api/auth.api';
import { authHelpers } from '@/lib/api/client';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardHeader, CardDescription, CardContent, CardFooter } from '@/components/ui/card';
import { GoogleSignInButton } from '@/components/auth/GoogleSignInButton';
import { UserRole } from '@/types';
import { extractApiErrorRawMessage } from '@/lib/api-error';

type LoginResponse = {
  data?: { access_token?: string };
  access_token?: string;
};

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [emailError, setEmailError] = useState('');
  const [passwordError, setPasswordError] = useState('');
  const [error, setError] = useState('');
  const [errorCode, setErrorCode] = useState<string | undefined>(undefined);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { login } = useAuth();

  const finalizeLogin = useCallback(async (token: string) => {
    authHelpers.setToken(token);
    const userData = await authApi.getMe();
    const primaryRole = userData.roles[0] as UserRole;
    login(token, primaryRole);
  }, [login]);

  const clearGoogleError = useCallback(() => {
    setError('');
    setErrorCode(undefined);
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setErrorCode(undefined);
    setEmailError('');
    setPasswordError('');
    setIsSubmitting(true);

    let hasLocalErr = false;
    if (!email.trim() || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      setEmailError('Vui lòng nhập Email hợp lệ');
      hasLocalErr = true;
    }
    if (!password) {
      setPasswordError('Vui lòng nhập Mật khẩu');
      hasLocalErr = true;
    }

    if (hasLocalErr) {
      setIsSubmitting(false);
      return;
    }

    try {
      // 1. Gọi API login → backend trả về { data: { access_token, ... } }
      const response = await authApi.login({ email, password }) as LoginResponse;
      const token = response.data?.access_token || response.access_token;

      if (!token) {
        setError('Không nhận được token từ server');
        return;
      }

      // 2. Hoàn tất lifecycle login theo flow hiện tại
      await finalizeLogin(token);
    } catch (err: unknown) {
      setError(extractApiErrorRawMessage(err) || 'Đăng nhập thất bại. Vui lòng kiểm tra lại.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleGoogleSubmit = useCallback(async ({ idToken, password }: { idToken: string; password?: string }) => {
    const response = await authApi.loginWithGoogle({ id_token: idToken, password }) as LoginResponse;
    const token = response.data?.access_token || response.access_token;
    if (!token) {
      throw new Error('Không nhận được token từ server');
    }
    await finalizeLogin(token);
  }, [finalizeLogin]);

  const onGoogleSignIn = useCallback(async ({ idToken, password }: { idToken: string; password?: string }) => {
    try {
      setErrorCode(undefined);
      await handleGoogleSubmit({ idToken, password });
    } catch (err: unknown) {
      const axiosErr = err as { response?: { data?: { error?: string; error_code?: string } } };
      setError(axiosErr.response?.data?.error || extractApiErrorRawMessage(err) || 'Đăng nhập Google thất bại.');
      setErrorCode(axiosErr.response?.data?.error_code);
      throw err;
    }
  }, [handleGoogleSubmit]);

  return (
    <div className="flex w-full max-w-screen-xl items-center justify-center">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <h1 className="text-2xl font-bold text-center">Iris School</h1>
          <CardDescription className="text-center">
            Đăng nhập để quản lý thông tin trường học
          </CardDescription>
        </CardHeader>
        <form noValidate onSubmit={handleSubmit}>
          <CardContent className="space-y-4">
            {error && (
              <div
                role="alert"
                aria-live="assertive"
                className="p-3 text-sm text-destructive-foreground bg-destructive rounded-md"
              >
                {error}
              </div>
            )}

            <div className="space-y-2">
              <label className="text-sm font-medium" htmlFor="email">Email</label>
              <Input
                id="email"
                type="email"
                placeholder="name@example.com"
                value={email}
                onChange={(e) => {
                  setEmail(e.target.value);
                  if (emailError) setEmailError('');
                }}
                aria-invalid={!!emailError}
                required
              />
              {emailError && <p className="text-[0.8rem] font-medium text-destructive mt-1">{emailError}</p>}
            </div>
            <div className="space-y-2 mb-3">
              <div className="flex items-center justify-between">
                <label className="text-sm font-medium" htmlFor="password">Mật khẩu</label>
                <Link
                  href="/forgot-password"
                  className="text-xs text-muted-foreground hover:text-primary transition-colors"
                >
                  Quên mật khẩu?
                </Link>
              </div>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => {
                  setPassword(e.target.value);
                  if (passwordError) setPasswordError('');
                }}
                aria-invalid={!!passwordError}
                required
              />
              {passwordError && <p className="text-[0.8rem] font-medium text-destructive mt-1">{passwordError}</p>}
            </div>
          </CardContent>
          <CardFooter className="flex flex-col gap-4">
            <Button className="w-full" type="submit" disabled={isSubmitting}>
              {isSubmitting ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Đang đăng nhập...
                </>
              ) : (
                'Đăng nhập'
              )}
            </Button>

            <div className="w-full space-y-3">
              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <span className="w-full border-t" />
                </div>
                <div className="relative flex justify-center text-xs">
                  <span className="bg-card px-2 text-muted-foreground">phương thức khác</span>
                </div>
              </div>

              <GoogleSignInButton
                onSubmitGoogle={onGoogleSignIn}
                errorCode={errorCode}
                clearError={clearGoogleError}
                disabled={isSubmitting}
              />
            </div>

            <p className="text-sm text-center text-muted-foreground">
              Phụ huynh chưa có tài khoản?{' '}
              <Link href="/register" className="font-medium text-primary hover:underline">
                Đăng ký tại đây
              </Link>
            </p>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}