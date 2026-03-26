/**
 * ProtectedRoute Component
 * Bọc các trang yêu cầu đăng nhập.
 * Kiểm tra trạng thái auth và role trước khi hiển thị nội dung.
 */
"use client";

import React, { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/providers/AuthProvider';
import { getDashboardRouteByRole } from '@/lib/auth-config';
import { UserRole } from '@/types';

interface ProtectedRouteProps {
  children: React.ReactNode;
  allowedRoles?: UserRole[];
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children, allowedRoles }) => {
  const { user, role, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading) {
      // 1. Nếu chưa đăng nhập -> redirect về login
      if (!user || !role) {
        router.push('/login');
        return;
      }

      // 2. Nếu có yêu cầu role cụ thể mà user không có -> redirect về dashboard tương ứng
      if (allowedRoles && !allowedRoles.includes(role)) {
        console.warn(`Access denied for role: ${role}. Allowed roles: ${allowedRoles}`);

        router.push(getDashboardRouteByRole(role));
      }
    }
  }, [user, role, isLoading, allowedRoles, router]);

  // Hiển thị loading state trong khi kiểm tra auth
  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
      </div>
    );
  }

  // Nếu đã qua các bước kiểm tra -> hiển thị nội dung
  if (user && role && (!allowedRoles || allowedRoles.includes(role))) {
    return <>{children}</>;
  }

  // Mặc định trả về null để tránh flash content
  return null;
};