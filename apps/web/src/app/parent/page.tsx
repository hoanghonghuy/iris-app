/**
 * Parent Dashboard
 * Trang chủ dành cho phụ huynh sau khi đăng nhập.
 */
"use client";

import React from 'react';
import { useAuth } from '@/providers/AuthProvider';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

export default function ParentDashboard() {
  const { user, logout } = useAuth();

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold tracking-tight">Bảng điều khiển Phụ huynh</h1>
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
              Vai trò: Phụ huynh
            </p>
          </CardContent>
        </Card>
      </div>

      <Card className="col-span-4">
        <CardHeader>
          <CardTitle>Con của tôi</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            Thông tin về con và các bài đăng từ lớp học sẽ được hiển thị tại đây.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}