/**
 * Login Page
 * Trang đăng nhập chính của ứng dụng.
 * Sử dụng AuthProvider để xử lý logic đăng nhập và redirect.
 */
"use client";

import React, { useState } from 'react';
import { useAuth } from '@/providers/AuthProvider';
import { authApi } from '@/lib/api/auth.api';
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
      // 1. Gọi API login
      const response = await authApi.login({ email, password });
      
      // 2. Giải mã role từ token (trong thực tế backend nên trả về role hoặc dùng JWT decode)
      // Ở đây ta giả định backend trả về role trong response hoặc ta lấy từ /me sau
      // Để đơn giản cho bước này, ta sẽ fetch /me để lấy role chính xác
      const userData = await authApi.getMe();
      const primaryRole = userData.roles[0] as UserRole;

      // 3. Lưu vào AuthProvider
      login(response.access_token, primaryRole);
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
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <label className="text-sm font-medium" htmlFor="password">Mật khẩu</label>
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
          <CardFooter>
            <Button className="w-full" type="submit" disabled={isSubmitting}>
              {isSubmitting ? 'Đang đăng nhập...' : 'Đăng nhập'}
            </Button>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}