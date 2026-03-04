/**
 * Admin Dashboard
 * Trang chủ dành cho Admin sau khi đăng nhập.
 */
"use client";

import React from 'react';
import { useAuth } from '@/providers/AuthProvider';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

export default function AdminDashboard() {
  const { user, role, logout } = useAuth();

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold tracking-tight">Bảng điều khiển Quản trị</h1>
        <Button variant="outline" onClick={logout}>Đăng xuất</Button>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Xin chào,</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{user?.email}</div>
            <p className="text-xs text-muted-foreground">
              Vai trò: {role === 'SUPER_ADMIN' ? 'Quản trị viên cấp cao' : 'Quản trị viên trường'}
            </p>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle>Quản lý Trường học</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Thêm, sửa, xóa thông tin các trường học trong hệ thống.
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Quản lý Người dùng</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Quản lý tài khoản giáo viên, phụ huynh và nhân viên.
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader>
            <CardTitle>Quản lý Lớp học</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Thiết lập danh sách lớp và phân công giáo viên.
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}