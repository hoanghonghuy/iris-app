/**
 * Login Page
 * Trang đăng nhập chính của ứng dụng.
 * Sử dụng AuthProvider để xử lý logic đăng nhập và redirect.
 */
"use client";

import React, { useState } from 'react';
import Link from 'next/link';
import { Loader2 } from 'lucide-react';
import { useAuth } from '@/providers/AuthProvider';
import { authApi } from '@/lib/api/auth.api';
import { authHelpers } from '@/lib/api/client';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card';
import { UserRole } from '@/types';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const { login } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsSubmitting(true);

    try {
      // 1. Gọi API login → backend trả về { data: { access_token, ... } }
      const response = await authApi.login({ email, password }) as any;
      const token = response.data?.access_token || response.access_token;

      if (!token) {
        setError('Không nhận được token từ server');
        return;
      }

      // 2. Lưu token vào localStorage TRƯỚC khi gọi /me
      authHelpers.setToken(token);

      // 3. Gọi /me để lấy role (giờ interceptor sẽ gắn token vào header)
      const meResponse = await authApi.getMe();
      const userData = meResponse.data; // backend wrap trong { data: {...} }
      const primaryRole = userData.roles[0] as UserRole;

      // 4. Lưu vào AuthProvider (sẽ redirect dựa theo role)
      login(token, primaryRole);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Đăng nhập thất bại. Vui lòng kiểm tra lại.');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-50 p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold text-center">Iris School</CardTitle>
          <CardDescription className="text-center">
            Đăng nhập để quản lý thông tin trường học
          </CardDescription>
        </CardHeader>
        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-4">
            {error && (
              <div className="p-3 text-sm text-white bg-destructive rounded-md">
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
                onChange={(e) => setEmail(e.target.value)}
                required
              />
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
                onChange={(e) => setPassword(e.target.value)}
                required
              />
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