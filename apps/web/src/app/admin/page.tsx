/**
 * Admin Dashboard
 * Trang chủ dành cho Admin sau khi đăng nhập.
 */
"use client";

import React, { useEffect, useState } from 'react';
import { useAuth } from '@/providers/AuthProvider';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { adminApi } from '@/lib/api/admin.api';
import { AdminAnalytics } from '@/types';
import { School, BookOpen, Users, GraduationCap, UsersRound, Loader2 } from 'lucide-react';
import Link from 'next/link';

export default function AdminDashboard() {
  const { user, role } = useAuth();
  const [stats, setStats] = useState<AdminAnalytics | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const data = await adminApi.getAnalytics();
        setStats(data);
      } catch (error) {
        console.error("Lỗi khi tải thống kê", error);
      } finally {
        setLoading(false);
      }
    };
    fetchStats();
  }, []);

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold tracking-tight">Bảng điều khiển Quản trị</h1>

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

      {loading ? (
        <div className="flex items-center justify-center py-8">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : stats ? (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-5">
          <Link href="/admin/schools">
            <Card className="hover:bg-zinc-50 dark:hover:bg-zinc-800/50 transition-colors">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">Tổng số Trường</CardTitle>
                <School className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.total_schools}</div>
              </CardContent>
            </Card>
          </Link>
          <Link href="/admin/classes">
            <Card className="hover:bg-zinc-50 dark:hover:bg-zinc-800/50 transition-colors">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">Tổng số Lớp học</CardTitle>
                <BookOpen className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.total_classes}</div>
              </CardContent>
            </Card>
          </Link>
          <Link href="/admin/teachers">
            <Card className="hover:bg-zinc-50 dark:hover:bg-zinc-800/50 transition-colors">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">Giáo viên</CardTitle>
                <Users className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.total_teachers}</div>
              </CardContent>
            </Card>
          </Link>
          <Link href="/admin/students">
            <Card className="hover:bg-zinc-50 dark:hover:bg-zinc-800/50 transition-colors">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">Học sinh</CardTitle>
                <GraduationCap className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.total_students}</div>
              </CardContent>
            </Card>
          </Link>
          <Link href="/admin/parents">
            <Card className="hover:bg-zinc-50 dark:hover:bg-zinc-800/50 transition-colors">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">Phụ huynh</CardTitle>
                <UsersRound className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.total_parents}</div>
              </CardContent>
            </Card>
          </Link>
        </div>
      ) : null}

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